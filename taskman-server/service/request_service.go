package service

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/go-common-lib/cipher"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/rpc"
	"io"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/exterror"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/gin-gonic/gin"
)

type RequestService struct {
	requestDao            *dao.RequestDao
	taskHandleTemplateDao *dao.TaskHandleTemplateDao
	taskHandleDao         *dao.TaskHandleDao
}

var (
	requestIdLock = new(sync.RWMutex)
)

func GetEntityData(requestId, userToken, language string) (result models.EntityQueryResult, err error) {
	var list []*models.ProcDefEntityDataObj
	requestTemplateId, tmpErr := getRequestTemplateByRequest(requestId)
	if tmpErr != nil {
		return result, tmpErr
	}
	requestTemplateObj, getTemplateErr := GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if getTemplateErr != nil {
		err = getTemplateErr
		return
	}
	if requestTemplateObj.PackageName == "" || requestTemplateObj.EntityName == "" {
		err = fmt.Errorf("RequestTemplate packageName or entityName illegal ")
		return
	}

	list, err = GetProcDefService().GetProcDefRootEntities(requestTemplateObj.ProcDefId, userToken, language)
	if err != nil {
		return
	}
	if len(list) > 0 {
		for _, v := range list {
			result.Data = append(result.Data, &models.EntityDataObj{Id: v.Id, DisplayName: v.DisplayName, PackageName: requestTemplateObj.PackageName, Entity: requestTemplateObj.EntityName})
		}
	}
	return
}

func ProcessDataPreview(requestTemplateId, entityDataId, userToken, language string) (result *models.EntityTreeData, err error) {
	requestTemplateObj, getTemplateErr := GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if getTemplateErr != nil {
		err = getTemplateErr
		return
	}
	if requestTemplateObj.ProcDefId == "" {
		err = fmt.Errorf("RequestTemplate proDefId illegal ")
		return
	}
	return GetProcDefService().GetProcDefDataPreview(requestTemplateObj.ProcDefId, entityDataId, userToken, language)
}

func ListRequest(param *models.QueryRequestParam, userRoles []string, userToken, permission, operator, language string) (pageInfo models.PageInfo, rowData []*models.RequestTable, err error) {
	rowData = []*models.RequestTable{}
	if strings.ToLower(permission) == "mgmt" {
		permission = "MGMT"
	} else {
		permission = "USE"
	}
	filterSql, _, queryParam := dao.TransFiltersToSQL(param, &models.TransFiltersParam{IsStruct: true, StructObj: models.RequestTable{}, PrimaryKey: "id"})
	userRolesFilterSql, userRolesFilterParam := dao.CreateListParams(userRoles, "")
	baseSql := fmt.Sprintf("select id,name,form,request_template,proc_instance_id,proc_instance_key,reporter,handler,report_time,emergency,status,expire_time,expect_time,confirm_time,created_by,created_time,updated_by,updated_time,rollback_desc from request where del_flag=0 and (created_by=? or request_template in (select id from request_template where id in (select request_template from request_template_role where role_type=? and `role` in ("+userRolesFilterSql+")))) %s ", filterSql)
	queryParam = append(append([]interface{}{operator, permission}, userRolesFilterParam...), queryParam...)
	if param.Paging {
		pageInfo.StartIndex = param.Pageable.StartIndex
		pageInfo.PageSize = param.Pageable.PageSize
		pageInfo.TotalRows = dao.QueryCount(baseSql, queryParam...)
		pageSql, pageParam := dao.TransPageInfoToSQL(*param.Pageable)
		baseSql += pageSql
		queryParam = append(queryParam, pageParam...)
	}
	err = dao.X.SQL(baseSql, queryParam...).Find(&rowData)
	if len(rowData) > 0 {
		var requestTemplateTable []*models.RequestTemplateTable
		dao.X.SQL("select id,name,status,version from request_template").Find(&requestTemplateTable)
		rtMap := make(map[string]string)
		for _, v := range requestTemplateTable {
			if v.Status != "confirm" {
				rtMap[v.Id] = fmt.Sprintf("%s(beta)", v.Name)
			} else {
				rtMap[v.Id] = fmt.Sprintf("%s(%s)", v.Name, v.Version)
			}
		}
		rtRoleMap := getRequestTemplateMGMTRole()
		var actions []*dao.ExecAction
		for _, v := range rowData {
			v.RequestTemplateName = rtMap[v.RequestTemplate]
			if tmpRoles, b := rtRoleMap[v.RequestTemplate]; b {
				v.HandleRoles = tmpRoles
			} else {
				v.HandleRoles = []string{}
			}
			if strings.Contains(v.Status, "InProgress") && v.ProcInstanceId != "" {
				newStatus := getInstanceStatus(v.ProcInstanceId, userToken, language)
				if newStatus == "InternallyTerminated" {
					newStatus = "Termination"
				}
				if newStatus != "" && newStatus != v.Status {
					actions = append(actions, &dao.ExecAction{Sql: "update request set status=? where id=?", Param: []interface{}{newStatus, v.Id}})
					v.Status = newStatus
				}
			}
			if v.Status == "Completed" {
				v.CompletedTime = v.UpdatedTime
			}
		}
		if len(actions) > 0 {
			updateStatusErr := dao.Transaction(actions)
			if updateStatusErr != nil {
				log.Logger.Error("Try to update request status fail", log.Error(updateStatusErr))
			}
		}
	}
	return
}

func getRequestTemplateMGMTRole() (result map[string][]string) {
	result = make(map[string][]string)
	var requestTemplateRole []*models.RequestTemplateRoleTable
	dao.X.SQL("select * from request_template_role where role_type='MGMT' order by request_template").Find(&requestTemplateRole)
	if len(requestTemplateRole) == 0 {
		return result
	}
	tmpTemplate := requestTemplateRole[0].RequestTemplate
	tmpRoles := []string{}
	for _, v := range requestTemplateRole {
		if v.RequestTemplate != tmpTemplate {
			result[tmpTemplate] = tmpRoles
			tmpTemplate = v.RequestTemplate
			tmpRoles = []string{}
		}
		tmpRoles = append(tmpRoles, v.Role)
	}
	if len(tmpRoles) > 0 {
		tmpTemplate = requestTemplateRole[len(requestTemplateRole)-1].RequestTemplate
		result[tmpTemplate] = tmpRoles
	}
	return result
}

func calcExpireObj(param *models.ExpireObj) {
	if param.ReportTime == "" || param.ExpireTime == "" {
		return
	}
	reportT, _ := time.Parse(models.DateTimeFormat, param.ReportTime)
	expireT, _ := time.Parse(models.DateTimeFormat, param.ExpireTime)
	nowT, _ := time.Parse(models.DateTimeFormat, param.NowTime)
	max := expireT.Sub(reportT).Seconds()
	use := nowT.Sub(reportT).Seconds()
	param.Percent = (use / max) * 100
	param.TotalDay = max / 86400
	param.UseDay = use / 86400
	param.LeftDay = (max - use) / 86400
	param.Percent, _ = strconv.ParseFloat(fmt.Sprintf("%.0f", param.Percent), 64)
	param.TotalDay, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", param.TotalDay), 64)
	param.LeftDay, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", param.LeftDay), 64)
	param.UseDay, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", param.UseDay), 64)
}

func calcExpireTime(reportTime string, expireDay int) (expire string) {
	t, err := time.Parse(models.DateTimeFormat, reportTime)
	if err != nil {
		return
	}
	expire = t.Add(time.Duration(expireDay*24) * time.Hour).Format(models.DateTimeFormat)
	return
}

func RevokeRequest(requestId, user string) (err error) {
	var request models.RequestTable
	var actions []*dao.ExecAction
	var latestDoingTask *models.TaskTable
	var taskSort int
	request, err = GetSimpleRequest(requestId)
	if err != nil {
		return
	}
	if request.Id == "" {
		err = fmt.Errorf("param requestId:%s not exist", requestId)
		return
	}
	// 校验是否能撤回
	if !calcShowRequestRevokeButton(requestId, request.Status) {
		err = exterror.New().RevokeRequestError
		return
	}
	if request.CreatedBy != user {
		err = fmt.Errorf("request not yours")
		return
	}
	// 添加撤回 任务
	newTaskId := "re_" + guid.CreateGuid()
	now := time.Now().Format(models.DateTimeFormat)
	// 删除正在处理的任务(只可能是定版或者审批)
	if latestDoingTask, err = GetTaskService().GetDoingTask(requestId, request.RequestTemplate, ""); err != nil {
		return
	}
	if latestDoingTask == nil {
		err = fmt.Errorf("not allow revoke")
		return
	}
	actions = append(actions, &dao.ExecAction{Sql: "update task set del_flag=1,status =? where id=?", Param: []interface{}{models.TaskStatusDone, latestDoingTask.Id}})
	actions = append(actions, &dao.ExecAction{Sql: "delete from task_handle where task= ?", Param: []interface{}{latestDoingTask.Id}})
	taskSort = GetTaskService().GenerateTaskOrderByRequestId(requestId)
	actions = append(actions, &dao.ExecAction{Sql: "insert into task (id,name,status,request,type,task_result,created_by,created_time,updated_by,updated_time,sort,request_created_time) values(?,?,?,?,?,?,?,?,?,?,?,?)",
		Param: []interface{}{newTaskId, "revoke", models.TaskStatusDone, request.Id, models.TaskTypeRevoke, models.TaskHandleResultTypeApprove, "system", now, "system", now, taskSort, request.CreatedTime}})
	actions = append(actions, &dao.ExecAction{Sql: "insert into task_handle (id,task,role,handler,handle_result,handle_status,created_time,updated_time) values(?,?,?,?,?,?,?,?)",
		Param: []interface{}{guid.CreateGuid(), newTaskId, request.Role, request.CreatedBy, models.TaskHandleResultTypeApprove, models.TaskHandleResultTypeComplete, now, now}})

	actions = append(actions, &dao.ExecAction{Sql: "update request set status ='Draft',revoke_flag=1,updated_time=? where id=?", Param: []interface{}{now, request.Id}})
	err = dao.Transaction(actions)
	return
}

func GetRequestWithRoot(requestId string) (result models.RequestTable, err error) {
	result = models.RequestTable{}
	var requestTable []*models.RequestTable
	err = dao.X.SQL("select id,name,form,request_template,proc_instance_id,proc_instance_key,reporter,report_time,emergency,status,cache,expire_time,expect_time,confirm_time,handler from request where id=?", requestId).Find(&requestTable)
	if err != nil {
		return
	}
	if len(requestTable) == 0 {
		err = fmt.Errorf("can not find any request with id:%s ", requestId)
		return
	}
	result = *requestTable[0]
	if result.Cache != "" {
		var cacheObj models.RequestPreDataDto
		err = json.Unmarshal([]byte(result.Cache), &cacheObj)
		if err != nil {
			err = fmt.Errorf("try to json unmarshal cache data fail,%s ", err.Error())
			return
		}
		result.Cache = cacheObj.RootEntityId
	}
	result.AttachFiles = GetRequestAttachFileList(requestId)
	return
}

func GetRequest(requestId string) (result models.RequestTable, err error) {
	result = models.RequestTable{}
	var requestTable []*models.RequestTable
	if err = dao.X.SQL("select id,name,form,request_template,proc_instance_id,proc_instance_key,reporter,report_time,emergency,status from request where id=?", requestId).Find(&requestTable); err != nil {
		return
	}
	if len(requestTable) == 0 {
		err = fmt.Errorf("can not find any request with id:%s ", requestId)
		return
	}
	result = *requestTable[0]
	return
}

func CreateRequest(param *models.RequestTable, operatorRoles []string, userToken, language string) error {
	var actions []*dao.ExecAction
	requestTemplateObj, err := GetRequestTemplateService().GetRequestTemplate(param.RequestTemplate)
	if err != nil {
		return err
	}
	if requestTemplateObj.ProcDefId != "" {
		err = GetRequestTemplateService().SyncProcDefId(requestTemplateObj.Id, requestTemplateObj.ProcDefId, requestTemplateObj.ProcDefName, "", userToken, language)
		if err != nil {
			return fmt.Errorf("try to sync proDefId fail,%s ", err.Error())
		}
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	param.Id = newRequestId()
	requestInsertAction := dao.ExecAction{Sql: "insert into request(id,name,request_template,reporter,emergency,report_role,status,expire_time," +
		"expect_time,handler,created_by,created_time,updated_by,updated_time,type,role) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
	requestInsertAction.Param = []interface{}{param.Id, param.Name, param.RequestTemplate, param.CreatedBy, param.Emergency,
		strings.Join(operatorRoles, ","), "Draft", "", param.ExpectTime, requestTemplateObj.Handler, param.CreatedBy, nowTime,
		param.CreatedBy, nowTime, param.Type, param.Role}
	actions = append(actions, &requestInsertAction)
	return dao.TransactionWithoutForeignCheck(actions)
}

func UpdateRequest(param *models.RequestTable) error {
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err := dao.X.Exec("update request set name=?,expect_time=?,emergency=?,handler=?,updated_by=?,updated_time=? where id=?", param.Name, param.ExpectTime, param.Emergency, param.Handler, param.UpdatedBy, nowTime, param.Id)
	return err
}

func DeleteRequest(requestId, operator string) error {
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err := dao.X.Exec("update request set del_flag=1,updated_by=?,updated_time=? where id=?", operator, nowTime, requestId)
	return err
}

func SaveRequestCacheNew(requestId, operator, userToken string, param *models.RequestPreDataDto) error {
	var err error
	if err = ValidateRequestForm(param.Data, userToken); err != nil {
		return err
	}
	var formItemNameQuery []*models.FormItemTemplateTable
	if err = dao.X.SQL("select item_group,group_concat(name,',') as name from form_item_template where in_display_name='yes' and form_template in (select form_template from request_template where id in (select request_template from request where id=?)) group by item_group", requestId).Find(&formItemNameQuery); err != nil {
		return err
	}
	itemGroupNameMap := make(map[string][]string)
	for _, v := range formItemNameQuery {
		itemGroupNameMap[v.ItemGroup] = strings.Split(v.Name, ",")
	}
	for _, v := range param.Data {
		nameList := itemGroupNameMap[v.ItemGroup]
		for _, value := range v.Value {
			if value.Id == "" {
				value.PackageName = v.PackageName
				value.EntityName = v.Entity
				value.EntityDataOp = "create"
				value.Id = fmt.Sprintf("tmp%s%s", models.SysTableIdConnector, guid.CreateGuid())
				value.DisplayName = concatItemDisplayName(value.EntityData, nameList)
			}
		}
	}
	paramBytes, err := json.Marshal(param)
	if err != nil {
		return fmt.Errorf("try to json marshal param fail,%s ", err.Error())
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	actions := UpdateRequestFormItem(requestId, operator, nowTime, param)
	actions = append(actions, &dao.ExecAction{Sql: "update request set cache=?,updated_by=?,updated_time=?,operator_obj=? where id=?", Param: []interface{}{string(paramBytes), operator, nowTime, param.EntityName, requestId}})
	return dao.Transaction(actions)
}

func SaveRequestCacheV2(requestId, operator, userToken string, param *models.RequestProDataV2Dto) error {
	var customFormCache []byte
	var taskApprovalCache string
	var err error
	if err = ValidateRequestForm(param.Data, userToken); err != nil {
		return err
	}
	var formItemNameQuery []*models.FormItemTemplateTable
	if err = dao.X.SQL("select item_group,group_concat(name,',') as name from form_item_template where in_display_name='yes' and form_template in (select request_template from form_template where id in (select request_template from request where id=?)) group by item_group", requestId).Find(&formItemNameQuery); err != nil {
		return err
	}
	itemGroupNameMap := make(map[string][]string)
	for _, v := range formItemNameQuery {
		itemGroupNameMap[v.ItemGroup] = strings.Split(v.Name, ",")
	}
	for _, v := range param.Data {
		nameList := itemGroupNameMap[v.ItemGroup]
		for _, value := range v.Value {
			if value.Id == "" {
				value.PackageName = v.PackageName
				value.EntityName = v.Entity
				value.EntityDataOp = "create"
				value.Id = fmt.Sprintf("tmp%s%s", models.SysTableIdConnector, guid.CreateGuid())
				value.DisplayName = concatItemDisplayName(value.EntityData, nameList)
			} else if value.EntityDataOp == "create" {
				if !strings.HasPrefix(value.Id, "tmp") {
					value.Id = fmt.Sprintf("tmp%s%s", models.SysTableIdConnector, value.Id)
				}
			}
		}
	}
	newParam := &models.RequestPreDataDto{
		RootEntityId: param.RootEntityId,
		EntityName:   param.EntityName,
		Data:         param.Data,
	}
	paramBytes, err := json.Marshal(newParam)
	if err != nil {
		return fmt.Errorf("try to json marshal param fail,%s ", err.Error())
	}
	if len(param.CustomForm.Value) > 0 {
		customFormCache, err = json.Marshal(param.CustomForm)
		if err != nil {
			return fmt.Errorf("try to json marshal param fail,%s ", err.Error())
		}
	}
	if len(param.ApprovalList) > 0 {
		for _, approval := range param.ApprovalList {
			if approval != nil && len(approval.HandleTemplates) > 0 {
				for _, handleTemplate := range approval.HandleTemplates {
					taskHandle, _ := GetRequestService().taskHandleTemplateDao.Get(handleTemplate.Id)
					if taskHandle != nil {
						if strings.TrimSpace(taskHandle.AssignRule) != "" {
							if err = json.Unmarshal([]byte(taskHandle.AssignRule), &handleTemplate.AssignRule); err != nil {
								return err
							}
						}
						if strings.TrimSpace(taskHandle.FilterRule) != "" {
							if err = json.Unmarshal([]byte(taskHandle.FilterRule), &handleTemplate.FilterRule); err != nil {
								return err
							}
						}
					}
				}
			}
		}
		approvalBytes, _ := json.Marshal(param.ApprovalList)
		if approvalBytes != nil {
			taskApprovalCache = string(approvalBytes)
		}
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	actions, buildActionErr := UpdateRequestFormItemNew(requestId, operator, newParam)
	if buildActionErr != nil {
		return fmt.Errorf("build update request form data action fail,%s ", buildActionErr.Error())
	}
	actions = append(actions, &dao.ExecAction{Sql: "update request set cache=?,updated_by=?,updated_time=?,name=?,description=?,expect_time=?,operator_obj=?,custom_form_cache=?,task_approval_cache=?,ref_id=?,ref_type=?" +
		" where id=?", Param: []interface{}{string(paramBytes), operator, nowTime, param.Name, param.Description, param.ExpectTime, param.EntityName, string(customFormCache), taskApprovalCache, param.RefId, param.RefType, requestId}})
	return dao.Transaction(actions)
}

func SaveRequestBindCache(requestId, operator string, param *models.RequestCacheData) error {
	cacheBytes, _ := json.Marshal(param)
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err := dao.X.Exec("update request set bind_cache=?,updated_by=?,updated_time=? where id=?", string(cacheBytes), operator, nowTime, requestId)
	return err
}

func concatItemDisplayName(rowData map[string]interface{}, nameList []string) string {
	displayNameList := []string{}
	for _, v := range nameList {
		if v == "" {
			continue
		}
		if _, b := rowData[v]; b {
			tmpValue := fmt.Sprintf("%s", rowData[v])
			if rowData[v] == nil {
				tmpValue = ""
			}
			displayNameList = append(displayNameList, tmpValue)
		}
	}
	return strings.Join(displayNameList, "__")
}

func UpdateRequestFormItem(requestId, operator, now string, param *models.RequestPreDataDto) []*dao.ExecAction {
	var actions []*dao.ExecAction
	var newFormId string
	actions = append(actions, &dao.ExecAction{Sql: "delete from form_item where request = ?", Param: []interface{}{requestId}})
	actions = append(actions, &dao.ExecAction{Sql: "delete from form where request = ?", Param: []interface{}{requestId}})

	for _, v := range param.Data {
		tmpFormIdList := guid.CreateGuidList(len(v.Value))
		for valueIndex, valueObj := range v.Value {
			// 每行记录新增一个表单
			newFormId = tmpFormIdList[valueIndex]
			actions = append(actions, &dao.ExecAction{Sql: "insert into form(id,request,form_template,data_id,created_by,created_time,updated_by," +
				"updated_time) values (?,?,?,?,?,?,?,?)", Param: []interface{}{newFormId, requestId, v.FormTemplateId, valueObj.Id, operator, now, operator, now}})
			tmpGuidList := guid.CreateGuidList(len(v.Title))
			for i, title := range v.Title {
				if strings.EqualFold(title.Multiple, models.Yes) || strings.EqualFold(title.Multiple, models.Y) {
					if tmpV, b := valueObj.EntityData[title.Name]; b {
						var tmpStringV []string
						for _, interfaceV := range tmpV.([]interface{}) {
							tmpStringV = append(tmpStringV, fmt.Sprintf("%s", interfaceV))
						}
						actions = append(actions, &dao.ExecAction{Sql: "insert into form_item(id,form,form_item_template,name,value,request,updated_time) values (?,?,?,?,?,?,?)",
							Param: []interface{}{tmpGuidList[i], newFormId, title.Id, title.Name, strings.Join(tmpStringV, ","), requestId, now}})
					}
				} else {
					actions = append(actions, &dao.ExecAction{Sql: "insert into form_item(id,form,form_item_template,name,value,request,updated_time) values (?,?,?,?,?,?,?)",
						Param: []interface{}{tmpGuidList[i], newFormId, title.Id, title.Name, valueObj.EntityData[title.Name], requestId, now}})
				}
			}
		}
	}
	return actions
}

func UpdateRequestFormItemNew(requestId, operator string, param *models.RequestPreDataDto) (actions []*dao.ExecAction, err error) {
	formParam := models.ProcessTaskFormParam{
		Operator:          operator,
		RequestPreDataDto: param,
		RequestId:         requestId,
		FormData:          param.Data,
	}
	actions, err = GetRequestService().processTaskForm(formParam)
	return
}

// 表单通用处理,根据 ProcessTaskFormParam里面 Task属性为空判断是请求保存,还是任务保存
func (s *RequestService) processTaskForm(formParam models.ProcessTaskFormParam) (actions []*dao.ExecAction, err error) {
	nowTime := time.Now()
	// 查请求数据池里的数据(里面的数据可能包含当前任务之前保存的数据)
	requestPoolRows := models.RequestPoolDataQueryRows{}
	if err = dao.X.SQL("select t1.id as form_id,t1.task,t1.form_template,t3.item_group,t3.item_group_type ,t1.data_id,t2.id as form_item_id,t2.form_item_template,t2.name,t2.value,t2.updated_time from form t1 left join form_item t2 on t1.id=t2.form left join form_template t3 on t1.form_template=t3.id where t1.request=?", formParam.RequestId).Find(&requestPoolRows); err != nil {
		err = fmt.Errorf("query request item pool data fail,%s ", err.Error())
	}
	requestPoolForms := requestPoolRows.DataParse()
	// 把form数据按itemGroup分开来
	itemGroupFormMap := make(map[string][]*models.RequestPoolForm)
	for _, poolForm := range requestPoolForms {
		if existForms, ok := itemGroupFormMap[poolForm.ItemGroup]; ok {
			itemGroupFormMap[poolForm.ItemGroup] = append(existForms, poolForm)
		} else {
			itemGroupFormMap[poolForm.ItemGroup] = []*models.RequestPoolForm{poolForm}
		}
	}
	for _, tableForm := range formParam.FormData {
		columnNameIdMap := make(map[string]string)
		isColumnMultiMap := make(map[string]int)
		passwordAttrMap := make(map[string]bool)
		for _, title := range tableForm.Title {
			columnNameIdMap[title.Name] = title.Id
			if strings.EqualFold(title.Multiple, models.Yes) || strings.EqualFold(title.Multiple, models.Y) {
				isColumnMultiMap[title.Name] = 1
			}
			if title.ElementType == string(models.FormItemElementTypePassword) {
				passwordAttrMap[title.Name] = true
			}
		}
		poolForms := itemGroupFormMap[tableForm.ItemGroup]
		for _, valueObj := range tableForm.Value {
			// 密码处理,web传递原密码,需要加密处理
			for key, value := range valueObj.EntityData {
				inputValue := ""
				if value != nil {
					inputValue = fmt.Sprintf("%+v", value)
				}
				if passwordAttrMap[key] && !strings.HasPrefix(strings.ToLower(inputValue), "{cipher_a}") {
					if inputValue, err = cipher.AesEnPasswordByGuid("", models.Config.EncryptSeed, inputValue, ""); err != nil {
						err = fmt.Errorf("try to encrypt password type column:%s value:%s fail,%s  ", key, inputValue, err.Error())
						return
					}
					valueObj.EntityData[key] = inputValue
				}
			}
			if formParam.Task == nil && valueObj.EntityDataOp == "create" && valueObj.Id != "" {
				if !strings.HasPrefix(valueObj.Id, "tmp") {
					valueObj.Id = fmt.Sprintf("tmp%s%s", models.SysTableIdConnector, valueObj.Id)
				}
			}
			if valueObj.Id == "" {
				valueObj.Id = fmt.Sprintf("tmp%s%s", models.SysTableIdConnector, guid.CreateGuid())
			}
			// 判断数据行的变化
			existForm := &models.RequestPoolForm{}
			for _, poolForm := range poolForms {
				if poolForm.DataId == valueObj.Id {
					existForm = poolForm
					break
				}
			}
			formId := existForm.FormId
			if formId == "" {
				formId = "form_" + guid.CreateGuid()
				// 数据行不存在，新增
				var taskId = sql.NullString{String: "", Valid: false}
				if formParam.Task != nil {
					taskId = sql.NullString{String: formParam.Task.Id, Valid: true}
				}
				actions = append(actions, &dao.ExecAction{Sql: "insert into form(id,request,task,form_template,data_id,created_by,updated_by,created_time,updated_time) values (?,?,?,?,?,?,?,?,?)", Param: []interface{}{
					formId, formParam.RequestId, taskId, tableForm.FormTemplateId, valueObj.Id, formParam.Operator, formParam.Operator, nowTime, nowTime,
				}})
			}
			// 判断数据行属性的变化
			for k, v := range valueObj.EntityData {
				var valueString string
				// 判断属性合不合法，是不是属性该表单的属性
				formItemTemplateId, nameLegalCheck := columnNameIdMap[k]
				if !nameLegalCheck {
					continue
				}
				// map需要特殊处理
				if reflect.ValueOf(v).Kind() == reflect.Map {
					m, ok := v.(map[string]interface{})
					if !ok {
						// 处理类型断言失败的情况
						valueString = fmt.Sprintf("%v", v)
					} else {
						byteArr, err := json.Marshal(m)
						if err != nil {
							// 处理 JSON 序列化失败的情况
							valueString = fmt.Sprintf("Error marshaling map: %v", err)
						} else {
							valueString = string(byteArr)
						}
					}
				} else {
					switch vType := v.(type) {
					case int:
						valueString = strconv.Itoa(vType)
					case int32:
						valueString = strconv.Itoa(int(vType))
					case int64:
						valueString = strconv.FormatInt(vType, 10)
					case float32:
						valueString = strconv.FormatFloat(float64(vType), 'f', -1, 32)
					case float64:
						valueString = strconv.FormatFloat(vType, 'f', -1, 64)
					default:
						// 如果 v 不是数字类型，使用 fmt.Sprintf
						valueString = fmt.Sprintf("%v", vType)
					}
				}
				if _, multipleFlag := isColumnMultiMap[k]; multipleFlag {
					// 此处需要支持 CMDB multiInt 和multiObject 类型
					if vInterfaceList, assertOk := v.([]interface{}); assertOk {
						var tmpV []string
						var multiObject bool
						for _, interfaceV := range vInterfaceList {
							switch interfaceType := interfaceV.(type) {
							// 因为 JSON 解析后数字默认是 float64,int也会变成float64
							case float64:
								tmpV = append(tmpV, fmt.Sprintf("%d", int(interfaceType)))
							case map[string]interface{}:
								multiObject = true
								break
							default:
								tmpV = append(tmpV, fmt.Sprintf("%s", interfaceV))
							}
						}
						if multiObject {
							byteArr, _ := json.Marshal(vInterfaceList)
							valueString = string(byteArr)
						} else {
							valueString = strings.Join(tmpV, ",")
						}
					} else {
						// 多选非必填情况下,valueString 为空
						valueString = ""
					}
				}
				// 从数据池里尝试查找有没有已存在的数据(同一个itemGroup，同一个数据行下的同一属性)
				latestPoolItem := getRequestPoolLatestItem(poolForms, valueObj.Id, k)
				var taskHandleId = sql.NullString{String: "", Valid: false}
				if formParam.Task != nil {
					taskHandleId = sql.NullString{String: formParam.TaskApproveParam.TaskHandleId, Valid: true}
				}
				if latestPoolItem.FormItemId == "" {
					// 没有在数据池里找到相关数据行的该属性
					actions = append(actions, &dao.ExecAction{Sql: "insert into form_item(id,form,form_item_template,name,value,request,updated_time,task_handle) values (?,?,?,?,?,?,?,?)", Param: []interface{}{
						"item_" + guid.CreateGuid(), formId, formItemTemplateId, k, valueString, formParam.RequestId, nowTime, taskHandleId,
					}})
				} else {
					if latestPoolItem.Value != valueString {
						// 数据有更新
						if formParam.Task != nil && latestPoolItem.Task == formParam.Task.Id {
							// 属于该任务，更新数据值
							actions = append(actions, &dao.ExecAction{Sql: "update form_item set value=?,updated_time=?,task_handle=? where id=?", Param: []interface{}{
								valueString, nowTime, taskHandleId, latestPoolItem.FormItemId,
							}})
						} else {
							// 不属于该任务，新增数据纪录
							actions = append(actions, &dao.ExecAction{Sql: "insert into form_item(id,form,form_item_template,name,value,request,updated_time,original_id,task_handle) values (?,?,?,?,?,?,?,?,?)", Param: []interface{}{
								"item_" + guid.CreateGuid(), formId, formItemTemplateId, k, valueString, formParam.RequestId, nowTime, latestPoolItem.FormItemId, taskHandleId,
							}})
						}
					}
				}
			}
		}
		// 如果之前是该任务保存的数据行但又没传过来了，说明已经删除行
		for _, poolForm := range poolForms {
			if formParam.Task == nil || poolForm.Task == formParam.Task.Id {
				deleteFlag := true
				for _, valueObj := range tableForm.Value {
					if poolForm.DataId == valueObj.Id {
						deleteFlag = false
						break
					}
				}
				if deleteFlag {
					actions = append(actions, &dao.ExecAction{Sql: "delete from form_item where form=?", Param: []interface{}{poolForm.FormId}})
					actions = append(actions, &dao.ExecAction{Sql: "delete from form where id=?", Param: []interface{}{poolForm.FormId}})
				}
			}
		}
	}
	return
}

func UpdateSingleRequestFormNew(requestId, operator string, param *models.RequestPreDataTableObj) (actions []*dao.ExecAction, err error) {
	updateParam := models.RequestPreDataDto{Data: []*models.RequestPreDataTableObj{param}}
	actions, err = UpdateRequestFormItemNew(requestId, operator, &updateParam)
	return
}

func GetRequestCache(requestId, cacheType string) (result interface{}, err error) {
	var requestTable []*models.RequestTable
	if cacheType == "data" {
		err = dao.X.SQL("select cache from request where id=?", requestId).Find(&requestTable)
		if err != nil {
			return
		}
		if len(requestTable) == 0 {
			err = fmt.Errorf("can not find any request with id:%s ", requestId)
			return
		}
		if requestTable[0].Cache == "" {
			return models.RequestPreDataDto{Data: []*models.RequestPreDataTableObj{}}, nil
		}
		var dataCache models.RequestPreDataDto
		err = json.Unmarshal([]byte(requestTable[0].Cache), &dataCache)
		return dataCache, err
	} else {
		err = dao.X.SQL("select bind_cache from request where id=?", requestId).Find(&requestTable)
		if err != nil {
			return
		}
		if len(requestTable) == 0 {
			err = fmt.Errorf("can not find any request with id:%s ", requestId)
			return
		}
		if requestTable[0].BindCache == "" {
			return models.RequestCacheData{TaskNodeBindInfos: []*models.RequestCacheTaskNodeBindObj{}}, nil
		}
		var bindCache models.RequestCacheData
		err = json.Unmarshal([]byte(requestTable[0].BindCache), &bindCache)
		return bindCache, err
	}
}

func getRequestTemplateByRequest(requestId string) (templateId string, err error) {
	var requestTable []*models.RequestTable
	err = dao.X.SQL("select request_template from request where id=?", requestId).Find(&requestTable)
	if err != nil {
		return
	}
	if len(requestTable) == 0 {
		err = fmt.Errorf("can not find request with id:%s ", requestId)
		return
	}
	templateId = requestTable[0].RequestTemplate
	return
}

func GetRequestRootForm(requestId string) (result models.RequestTemplateFormStruct, err error) {
	var formTemplate *models.FormTemplateTable
	var items []*models.FormItemTemplateTable
	result = models.RequestTemplateFormStruct{}
	requestTemplateId, tmpErr := getRequestTemplateByRequest(requestId)
	if tmpErr != nil {
		return result, tmpErr
	}
	requestTemplateObj, _ := GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	result.Id = requestTemplateObj.Id
	result.Name = requestTemplateObj.Name
	result.PackageName = requestTemplateObj.PackageName
	result.EntityName = requestTemplateObj.EntityName
	result.ProcDefName = requestTemplateObj.ProcDefName
	result.ProcDefId = requestTemplateObj.ProcDefId
	result.ProcDefKey = requestTemplateObj.ProcDefKey
	formTemplate, err = GetFormTemplateService().QueryRequestFormByRequestTemplateIdAndType(requestTemplateId, string(models.RequestFormTypeMessage))
	if formTemplate != nil {
		dao.X.SQL("select * from form_item_template where form_template=?", formTemplate.Id).Find(&items)
		result.FormItems = items
	}
	return
}

func GetRequestPreData(requestId, entityDataId, userToken, language string) (result []*models.RequestPreDataTableObj, previewData *models.EntityTreeData, err error) {
	var requestTables []*models.RequestTable
	var targetData []map[string]interface{}
	err = dao.X.SQL("select cache from request where id=?", requestId).Find(&requestTables)
	if err != nil {
		return
	}
	if len(requestTables) == 0 {
		return result, previewData, fmt.Errorf("can not find requestId:%s ", requestId)
	}
	if requestTables[0].Cache != "" {
		var cacheObj models.RequestPreDataDto
		err = json.Unmarshal([]byte(requestTables[0].Cache), &cacheObj)
		if err != nil {
			return result, previewData, fmt.Errorf("try to json unmarshal cache data fail,%s ", err.Error())
		}
		if cacheObj.RootEntityId == entityDataId {
			result = cacheObj.Data
			SensitiveDataEncryption(result)
			return
		}
	}
	result = []*models.RequestPreDataTableObj{}
	requestTemplateId, tmpErr := getRequestTemplateByRequest(requestId)
	if tmpErr != nil {
		return result, previewData, tmpErr
	}
	result = getRequestPreDataByTemplateId(requestTemplateId)
	if entityDataId == "" {
		return
	}
	previewData, err = ProcessDataPreview(requestTemplateId, entityDataId, userToken, language)
	if err != nil {
		return result, previewData, err
	}
	if len(previewData.EntityTreeNodes) == 0 {
		return
	}
	// 处理表单敏感数据
	for _, entityData := range previewData.EntityTreeNodes {
		if targetData, err = HandleFormSensitiveData([]map[string]interface{}{entityData.EntityData}, entityData.EntityName, userToken); err != nil {
			return
		}
		if len(targetData) > 0 {
			entityData.EntityData = targetData[0]
		}
	}

	previewDataBytes, _ := json.Marshal(previewData)
	_, err = dao.X.Exec("update request set preview_cache=? where id=?", string(previewDataBytes), requestId)
	if err != nil {
		err = fmt.Errorf("update request preview cache data fail,%s ", err.Error())
		return
	}
	for _, entity := range result {
		if entity.ItemGroupRule == "new" {
			continue
		}
		for _, tmpData := range previewData.EntityTreeNodes {
			if tmpData.EntityName == entity.Entity {
				tmpValueData := make(map[string]interface{})
				for _, title := range entity.Title {
					if title.RoutineExpression != "" {
						tmpEntityData, tmpQueryErr := rpc.QueryEntityExpressionData(title.RoutineExpression, tmpData.DataId, userToken, language)
						if tmpQueryErr != nil {
							err = tmpQueryErr
							return
						}
						tmpEntityDataBytes, _ := json.Marshal(tmpEntityData)
						tmpValueData[title.Name] = string(tmpEntityDataBytes)
					} else {
						tmpValueData[title.Name] = tmpData.EntityData[title.Name]
					}
				}
				entity.Value = append(entity.Value, &models.EntityTreeObj{Id: tmpData.Id, PackageName: tmpData.PackageName, EntityName: tmpData.EntityName, DataId: tmpData.DataId, PreviousIds: tmpData.PreviousIds, SucceedingIds: tmpData.SucceedingIds, DisplayName: tmpData.DisplayName, FullDataId: tmpData.FullDataId, EntityData: tmpValueData})
			}
		}
	}
	return
}

func getRequestPreDataByTemplateId(requestTemplateId string) []*models.RequestPreDataTableObj {
	var result []*models.RequestPreDataTableObj
	var formTemplateList []*models.FormTemplateTable
	dao.X.SQL("select * from form_template where request_template=? and request_form_type = ? order by item_group_sort", requestTemplateId, models.RequestFormTypeData).Find(&formTemplateList)
	if len(formTemplateList) > 0 {
		for _, formTemplate := range formTemplateList {
			var formItemTemplateList []*models.FormItemTemplateTable
			var packageName, entity string
			var title []*models.FormItemTemplateDto
			dao.X.SQL("select * from form_item_template where form_template=?  order by sort", formTemplate.Id).Find(&formItemTemplateList)
			if len(formItemTemplateList) > 0 {
				for _, formItem := range formItemTemplateList {
					title = append(title, models.ConvertFormItemTemplateModel2Dto(formItem, *formTemplate))
					if packageName == "" && formItem.PackageName != "" {
						packageName = formItem.PackageName
					}
					if entity == "" && formItem.Entity != "" {
						entity = formItem.Entity
					}
				}
			}
			result = append(result, &models.RequestPreDataTableObj{
				PackageName:    packageName,
				Entity:         entity,
				FormTemplateId: formTemplate.Id,
				ItemGroup:      formTemplate.ItemGroup,
				ItemGroupName:  formTemplate.ItemGroupName,
				ItemGroupType:  formTemplate.ItemGroupType,
				ItemGroupRule:  formTemplate.ItemGroupRule,
				Title:          title,
				Value:          []*models.EntityTreeObj{},
			})
		}
	}
	return result
}

func getItemTemplateTitle(items []*models.FormItemTemplateTable) []*models.RequestPreDataTableObj {
	var result []*models.RequestPreDataTableObj
	tmpPackageName := items[0].PackageName
	tmpEntity := items[0].Entity
	tmpItemGroup := items[0].ItemGroup
	tmpItemGroupName := items[0].ItemGroupName
	tmpFormTemplateId := items[0].FormTemplate
	formTemplateIdList := []string{items[0].FormTemplate}
	var tmpRefEntity []string
	var tmpItems []*models.FormItemTemplateDto
	existItemMap := make(map[string]int)
	for _, v := range items {
		tmpKey := fmt.Sprintf("%s__%s", v.ItemGroup, v.Name)
		if _, b := existItemMap[tmpKey]; b {
			continue
		} else {
			existItemMap[tmpKey] = 1
		}
		if v.ItemGroup != tmpItemGroup {
			if tmpItemGroup != "" {
				result = append(result, &models.RequestPreDataTableObj{
					PackageName:    tmpPackageName,
					Entity:         tmpEntity,
					FormTemplateId: tmpFormTemplateId,
					ItemGroup:      tmpItemGroup,
					ItemGroupName:  tmpItemGroupName,
					RefEntity:      tmpRefEntity,
					SortLevel:      0,
					Title:          tmpItems,
					Value:          []*models.EntityTreeObj{},
				})
			}
			tmpItems = []*models.FormItemTemplateDto{}
			tmpEntity = v.Entity
			tmpPackageName = v.PackageName
			tmpItemGroup = v.ItemGroup
			tmpItemGroupName = v.ItemGroupName
			tmpFormTemplateId = v.FormTemplate
			formTemplateIdList = append(formTemplateIdList, v.FormTemplate)
			tmpRefEntity = []string{}
		} else {
			if tmpEntity == "" && v.Entity != "" {
				tmpEntity = v.Entity
				tmpPackageName = v.PackageName
			}
		}
		tmpItems = append(tmpItems, models.ConvertFormItemTemplateModel2Dto(v, models.FormTemplateTable{}))
		if v.RefEntity != "" {
			existFlag := false
			for _, vv := range tmpRefEntity {
				if vv == v.RefEntity {
					existFlag = true
					break
				}
			}
			if !existFlag {
				tmpRefEntity = append(tmpRefEntity, v.RefEntity)
			}
		}
	}
	if len(tmpItems) > 0 {
		if items[len(items)-1].Entity != "" {
			tmpEntity = items[len(items)-1].Entity
		}
		if items[len(items)-1].PackageName != "" {
			tmpPackageName = items[len(items)-1].PackageName
		}
		tmpItemGroup = items[len(items)-1].ItemGroup
		tmpItemGroupName = items[len(items)-1].ItemGroupName
		result = append(result, &models.RequestPreDataTableObj{
			PackageName:    tmpPackageName,
			Entity:         tmpEntity,
			FormTemplateId: tmpFormTemplateId,
			ItemGroup:      tmpItemGroup,
			ItemGroupName:  tmpItemGroupName,
			RefEntity:      tmpRefEntity,
			SortLevel:      0,
			Title:          tmpItems,
			Value:          []*models.EntityTreeObj{},
		})
	}
	if len(formTemplateIdList) > 0 {
		var formTemplateRows []*models.FormTemplateTable
		filterSql, filterParam := dao.CreateListParams(formTemplateIdList, "")
		err := dao.X.SQL("select id,item_group_type,item_group_rule,item_group_sort from form_template where id in ("+filterSql+")", filterParam...).Find(&formTemplateRows)
		if err != nil {
			log.Logger.Error("query for template table fail", log.Error(err))
		} else {
			formTemplateMap := make(map[string]*models.FormTemplateTable)
			for _, row := range formTemplateRows {
				formTemplateMap[row.Id] = row
			}
			for _, v := range result {
				if formTemplateRow, ok := formTemplateMap[v.FormTemplateId]; ok {
					v.ItemGroupType = formTemplateRow.ItemGroupType
					v.ItemGroupRule = formTemplateRow.ItemGroupRule
					for _, item := range v.Title {
						item.ItemGroupType = formTemplateRow.ItemGroupType
						item.ItemGroupRule = formTemplateRow.ItemGroupRule
						item.ItemGroupSort = formTemplateRow.ItemGroupSort
					}
				}
			}
		}
	}
	result = sortRequestEntity(result)
	for _, v := range result {
		for _, vv := range v.Title {
			if vv.SelectList == nil {
				vv.SelectList = []*models.EntityDataObj{}
			}
		}
	}
	return result
}

func sortRequestEntity(param []*models.RequestPreDataTableObj) models.RequestPreDataSort {
	var result models.RequestPreDataSort
	var entityMap = make(map[string][]string)
	entityNumMap := make(map[string]int)
	for _, v := range param {
		for _, vv := range v.RefEntity {
			if vv == v.Entity {
				continue
			}
			if entityMap[vv] != nil {
				entityMap[vv] = append(entityMap[vv], v.Entity)
			} else {
				entityMap[vv] = []string{v.Entity}
			}
		}
	}
	countNum := 0
	levelNum := 1
	for _, v := range param {
		if v.Entity == "" {
			v.SortLevel = levelNum
			entityNumMap[v.Entity] = v.SortLevel
			countNum = countNum + 1
			continue
		}
		if _, b := entityMap[v.Entity]; !b {
			v.SortLevel = levelNum
			entityNumMap[v.Entity] = v.SortLevel
			countNum = countNum + 1
		}
	}
	for countNum < len(param) {
		levelNum = levelNum + 1
		for k, v := range entityMap {
			ready := true
			for _, vv := range v {
				if _, b := entityNumMap[vv]; !b {
					ready = false
					break
				} else {
					if entityNumMap[vv] == levelNum {
						ready = false
						break
					}
				}
			}
			if ready {
				for _, vv := range param {
					if vv.Entity == k {
						vv.SortLevel = levelNum
						entityNumMap[vv.Entity] = vv.SortLevel
						countNum = countNum + 1
					}
				}
				delete(entityMap, k)
			}
		}
	}
	for _, v := range param {
		result = append(result, v)
	}
	sort.Sort(result)
	return result
}

func CheckRequest(request models.RequestTable, task *models.TaskTable, operator, userToken, language string, cacheData models.RequestCacheData) (err error) {
	var entityDepMap map[string][]string
	var requestTemplate *models.RequestTemplateTable
	var actions, approvalActions []*dao.ExecAction
	var checkTaskHandle *models.TaskHandleTable
	var taskSort int
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(request.RequestTemplate)
	if err != nil {
		return
	}
	if requestTemplate == nil {
		return
	}
	entityDepMap, _, err = AppendUselessEntity(requestTemplate.Id, userToken, language, &cacheData, request.Id)
	if err != nil {
		err = fmt.Errorf("try to append useless entity fail,%s ", err.Error())
		return
	}
	if requestTemplate.ProcDefId != "" {
		request.AssociationWorkflow = true
	}
	fillBindingWithRequestData(request.Id, userToken, language, &cacheData, entityDepMap)
	cacheBytes, _ := json.Marshal(cacheData)
	nowTime := time.Now().Format(models.DateTimeFormat)
	expireTime := calcExpireTime(nowTime, requestTemplate.ExpireDay)
	checkTaskHandle, err = GetTaskHandleService().GetRequestCheckTaskHandle(task.Id)
	if err != nil {
		return
	}
	if checkTaskHandle == nil {
		err = fmt.Errorf("check taskhandle is empty")
		return
	}
	// 更新请求表
	actions = append(actions, &dao.ExecAction{Sql: "update request set handler=?,confirm_time=?,expire_time=?,bind_cache=?,updated_by=?,updated_time=? where id=?",
		Param: []interface{}{operator, nowTime, expireTime, string(cacheBytes), operator, nowTime, request.Id}})
	// 更新请求处理状态为完成
	actions = append(actions, &dao.ExecAction{Sql: "update task_handle set handle_result=?,handle_status=?,updated_time=? where id=?",
		Param: []interface{}{models.TaskHandleResultTypeApprove, models.TaskHandleResultTypeComplete, nowTime, checkTaskHandle.Id}})
	// 更新任务为完成
	actions = append(actions, &dao.ExecAction{Sql: "update task set status=?,task_result=?,updated_by=?,updated_time=? where id=?",
		Param: []interface{}{models.TaskStatusDone, models.TaskHandleResultTypeComplete, operator, nowTime, task.Id}})
	taskSort = GetTaskService().GenerateTaskOrderByRequestId(request.Id)
	approvalActions, err = GetRequestService().CreateRequestApproval(request, "", userToken, language, taskSort, false)
	if err != nil {
		return
	}
	if len(approvalActions) > 0 {
		actions = append(actions, approvalActions...)
	}
	if err = dao.Transaction(actions); err != nil {
		return
	}
	return GetRequestService().AutoExecTaskHandle(request, userToken, language)
}

func StartRequest(request models.RequestTable, operator, userToken, language string, cacheData models.RequestCacheData) (result *models.StartInstanceResultData, err error) {
	var requestTemplateTable []*models.RequestTemplateTable
	dao.X.SQL("select * from request_template where id in (select request_template from request where id=?)", request.Id).Find(&requestTemplateTable)
	if len(requestTemplateTable) == 0 {
		return result, fmt.Errorf("can not find requestTemplate with request:%s ", request.Id)
	}
	cacheData.ProcDefId = requestTemplateTable[0].ProcDefId
	cacheData.ProcDefKey = requestTemplateTable[0].ProcDefKey
	entityDepMap, preData, tmpErr := AppendUselessEntity(requestTemplateTable[0].Id, userToken, language, &cacheData, request.Id)
	if tmpErr != nil {
		return result, fmt.Errorf("try to append useless entity fail,%s ", tmpErr.Error())
	}
	fillBindingWithRequestData(request.Id, userToken, language, &cacheData, entityDepMap)
	cacheBytes, _ := json.Marshal(cacheData)
	log.Logger.Info("cacheByte", log.String("cacheBytes", string(cacheBytes)))
	startParam := BuildRequestProcessData(cacheData, preData)
	startParam.SimpleRequestDto = &models.SimpleRequestDto{
		Id:              request.Id,
		Name:            request.Name,
		RequestTemplate: request.RequestTemplate,
		CreatedBy:       request.CreatedBy,
		CreatedTime:     request.CreatedTime,
		Type:            request.Type,
	}
	result, err = GetProcDefService().StartProcDefInstances(startParam, userToken, language)
	if err != nil {
		return
	}
	if result == nil {
		err = fmt.Errorf("StartProcDefInstances response empty")
		return
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	expireTime := calcExpireTime(nowTime, requestTemplateTable[0].ExpireDay)
	procInstId := fmt.Sprintf("%v", result.Id)
	_, err = dao.X.Exec("update request set handler=?,proc_instance_id=?,proc_instance_key=?,confirm_time=?,expire_time=?,status=?,bind_cache=?,updated_by=?,updated_time=? where id=?", operator, procInstId, result.ProcInstKey, nowTime, expireTime, result.Status, string(cacheBytes), operator, nowTime, request.Id)
	return
}

func StartRequestNew(request models.RequestTable, userToken, language string, cacheData models.RequestCacheData) (actions []*dao.ExecAction, err error) {
	log.Logger.Debug("StartRequestNew", log.JsonObj("request", request))
	var requestTemplateTable []*models.RequestTemplateTable
	var result *models.StartInstanceResultData
	actions = []*dao.ExecAction{}
	err = dao.X.SQL("select * from request_template where id=?", request.RequestTemplate).Find(&requestTemplateTable)
	if err != nil {
		err = fmt.Errorf("query request_template with id:%s fail,%s ", request.RequestTemplate, err.Error())
		return
	}
	if len(requestTemplateTable) == 0 {
		err = fmt.Errorf("can not find requestTemplate with request:%s ,requestTemplateId:%s ", request.Id, request.RequestTemplate)
		return
	}
	if err = buildCacheDataWithPool(request.Id, &cacheData); err != nil {
		err = fmt.Errorf("build cache data with request pool fail,%s ", err.Error())
		return
	}
	cacheData.ProcDefId = requestTemplateTable[0].ProcDefId
	cacheData.ProcDefKey = requestTemplateTable[0].ProcDefKey
	entityDepMap, preData, tmpErr := AppendUselessEntity(requestTemplateTable[0].Id, userToken, language, &cacheData, request.Id)
	if tmpErr != nil {
		err = fmt.Errorf("try to append useless entity fail,%s ", tmpErr.Error())
		return
	}
	fillBindingWithRequestData(request.Id, userToken, language, &cacheData, entityDepMap)
	startParam := BuildRequestProcessData(cacheData, preData)
	startParam.SimpleRequestDto = &models.SimpleRequestDto{
		Id:              request.Id,
		Name:            request.Name,
		RequestTemplate: request.RequestTemplate,
		CreatedBy:       request.CreatedBy,
		CreatedTime:     request.CreatedTime,
		Type:            request.Type,
	}
	log.Logger.Debug("start proc instance", log.JsonObj("startParam", startParam))
	result, err = GetProcDefService().StartProcDefInstances(startParam, userToken, language)
	if err != nil {
		return
	}
	if result == nil {
		err = fmt.Errorf("StartProcDefInstances response empty")
		return
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	procInstId := fmt.Sprintf("%v", result.Id)
	actions = append(actions, &dao.ExecAction{Sql: "update request set proc_instance_id=?,proc_instance_key=?,status=?,updated_time=? where id=?", Param: []interface{}{
		procInstId, result.ProcInstKey, result.Status, nowTime, request.Id,
	}})
	return
}

func buildCacheDataWithPool(requestId string, cacheData *models.RequestCacheData) (err error) {
	// 查请求数据池里的数据(里面的数据可能包含当前任务之前保存的数据)
	requestPoolRows := models.RequestPoolDataQueryRows{}
	if err = dao.X.SQL("select t1.id as form_id,t1.task,t1.form_template,t3.item_group,t3.item_group_type ,t1.data_id,t2.id as form_item_id,t2.form_item_template,t2.name,t2.value,t2.updated_time from form t1 left join form_item t2 on t1.id=t2.form left join form_template t3 on t1.form_template=t3.id where t1.request=?", requestId).Find(&requestPoolRows); err != nil {
		err = fmt.Errorf("query request item pool data fail,%s ", err.Error())
		return
	}
	requestPoolForms := requestPoolRows.DataParse()
	rowAttrMap := make(map[string]map[string]string)
	for _, poolForm := range requestPoolForms {
		tmpKVMap := make(map[string]string)
		for _, row := range poolForm.Items {
			if _, existFlag := tmpKVMap[row.Name]; !existFlag {
				matchItem := getRequestPoolLatestItem(requestPoolForms, poolForm.DataId, row.Name)
				tmpKVMap[row.Name] = matchItem.Value
			}
		}
		rowAttrMap[poolForm.DataId] = tmpKVMap
	}
	if rowKVData, ok := rowAttrMap[cacheData.RootEntityValue.Oid]; ok {
		for _, v := range cacheData.RootEntityValue.AttrValues {
			if newValue, matchFlag := rowKVData[v.AttrName]; matchFlag {
				v.DataValue = newValue
			}
		}
	}
	for _, node := range cacheData.TaskNodeBindInfos {
		for _, boundEntity := range node.BoundEntityValues {
			if rowKVData, ok := rowAttrMap[boundEntity.Oid]; ok {
				for _, v := range boundEntity.AttrValues {
					if newValue, matchFlag := rowKVData[v.AttrName]; matchFlag {
						v.DataValue = newValue
					}
				}
			}
		}
	}
	return
}

func UpdateRequestStatus(requestId, status, operator, userToken, language, description string) error {
	var err error
	var request models.RequestTable
	var requestTemplate *models.RequestTemplateTable
	var bindCache string
	nowTime := time.Now().Format(models.DateTimeFormat)
	request, err = GetSimpleRequest(requestId)
	if err != nil {
		return err
	}
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(request.RequestTemplate)
	if err != nil {
		return err
	}
	if requestTemplate == nil {
		return fmt.Errorf("requestId:%s is invalid", requestId)
	}
	if status == "Pending" {
		if requestTemplate.ProcDefId != "" {
			bindData, bindErr := GetRequestPreBindData(request, requestTemplate, userToken, language)
			if bindErr != nil {
				return fmt.Errorf("try to build bind data fail,%s ", bindErr.Error())
			}
			bindCacheBytes, _ := json.Marshal(bindData)
			bindCache = string(bindCacheBytes)
		}
		// 设置bindCache,模版没配置定版,审批自动通过到任务场景用到
		request.BindCache = bindCache
		// 请求定版, 根据模板配置开启是否确认定版
		err = GetRequestService().CreateRequestCheck(request, operator, request.Cache, bindCache, userToken, language)
	} else if status == "Draft" {
		// 只有定版人才能处理
		var checkTask *models.TaskTable
		var taskHandleList []*models.TaskHandleTable
		var actions []*dao.ExecAction
		if checkTask, err = GetTaskService().GetLatestCheckTask(requestId); err != nil {
			return err
		}
		if checkTask == nil {
			err = fmt.Errorf("requestId:%s check task not exist", requestId)
			return err
		}

		taskHandleList, err = GetTaskHandleService().GetTaskHandleListByTaskId(checkTask.Id)
		if err != nil {
			return err
		}

		if len(taskHandleList) == 0 {
			err = fmt.Errorf("requestId:%s check task handler is empty", requestId)
			return err
		}
		if taskHandleList[0].Handler != operator {
			err = exterror.New().UpdateRequestHandlerStatusError
			return err
		}
		// 更新处理人,拒绝
		actions = append(actions, &dao.ExecAction{Sql: "update task_handle set handle_result = ?,handle_status = ?,result_desc=?,updated_time =? where id= ?", Param: []interface{}{models.TaskHandleResultTypeDeny, models.TaskHandleResultTypeComplete, description, nowTime, taskHandleList[0].Id}})
		// 更新任务到完成
		actions = append(actions, &dao.ExecAction{Sql: "update task set status = ?,task_result = ?,description = ?,updated_by =?,updated_time =? where id = ?", Param: []interface{}{models.TaskStatusDone, models.TaskHandleResultTypeRedraw, description, operator, nowTime, checkTask.Id}})
		// 更新请求
		actions = append(actions, &dao.ExecAction{Sql: "update request set status=?,rollback_desc=?,updated_by=?,handler=?,updated_time=?,confirm_time=? where id=?", Param: []interface{}{status, description, operator, operator, nowTime, nowTime, requestId}})
		// 定版退回邮件通知
		go NotifyTaskBackMail(request.Name, RequestPending, request.CreatedBy, taskHandleList[0].Handler, userToken, language)
		err = dao.Transaction(actions)
	} else {
		_, err = dao.X.Exec("update request set status=?,updated_by=?,updated_time=? where id=?", status, operator, nowTime, requestId)
	}
	return err
}

func fillBindingWithRequestData(requestId, userToken, language string, cacheData *models.RequestCacheData, existDepMap map[string][]string) {
	var items []*models.FormItemTemplateTable
	//dao.X.SQL("select * from form_item_template where form_template in (select form_template from request_template where id in (select request_template from request where id=?)) order by entity,sort", requestId).Find(&items)
	dao.X.SQL("select * from form_item_template where form_template in (select id from form_template where request_template in (select request_template from request where id=?)) order by entity,sort", requestId).Find(&items)
	itemMap := make(map[string][]string)
	for _, item := range items {
		if item.RefEntity == "" {
			continue
		}
		if nowRefList, b := itemMap[item.Entity]; b {
			existFlag := false
			for _, v := range nowRefList {
				if v == item.RefEntity {
					existFlag = true
				}
			}
			if !existFlag {
				itemMap[item.Entity] = append(itemMap[item.Entity], item.RefEntity)
			}
		} else {
			itemMap[item.Entity] = []string{item.RefEntity}
		}
	}
	// itemMap -> entity:[refEntity]
	for k, v := range itemMap {
		log.Logger.Info("itemMap", log.String("key", k), log.StringList("value", v))
	}
	entityNewMap := make(map[string][]string)
	for k, v := range existDepMap {
		log.Logger.Info("existDepMap", log.String("k", k), log.StringList("v", v))
		if k != "" && len(v) > 0 {
			entityNewMap[k] = v
		}
	}
	entityOidMap := make(map[string]int)
	dataIdOidMap := make(map[string]string)
	matchEntityRoot(requestId, userToken, language, cacheData)
	if cacheData.RootEntityValue.EntityDataOp != "create" {
		dataIdOidMap[cacheData.RootEntityValue.EntityDataId] = cacheData.RootEntityValue.Oid
	}
	for _, taskNode := range cacheData.TaskNodeBindInfos {
		for _, entityValue := range taskNode.BoundEntityValues {
			if entityValue.EntityDataId != "" {
				dataIdOidMap[entityValue.EntityDataId] = entityValue.Oid
			}
			entityOidMap[entityValue.Oid] = 1
			if _, b := entityNewMap[entityValue.Oid]; b {
				continue
			}
			if entityRefs, b := itemMap[entityValue.EntityName]; b {
				findEntityRefByItemRef(entityValue, entityRefs, entityNewMap, dataIdOidMap)
			}
		}
	}
	for k, v := range dataIdOidMap {
		log.Logger.Debug("dataIdOidMap", log.String("k", k), log.String("v", v))
	}
	for k, v := range entityNewMap {
		tmpRefOidList := []string{}
		for _, vv := range v {
			if oid, b := dataIdOidMap[vv]; b {
				tmpRefOidList = append(tmpRefOidList, oid)
			} else {
				tmpRefOidList = append(tmpRefOidList, vv)
			}
		}
		entityNewMap[k] = tmpRefOidList
		log.Logger.Info("entityNewMap", log.String("key", k), log.StringList("value", tmpRefOidList))
	}
	if len(entityNewMap) > 0 {
		rebuildEntityRefOids(&cacheData.RootEntityValue, entityNewMap, entityOidMap)
		for _, taskNode := range cacheData.TaskNodeBindInfos {
			for _, entityValue := range taskNode.BoundEntityValues {
				rebuildEntityRefOids(entityValue, entityNewMap, entityOidMap)
			}
		}
	}
}

func rebuildEntityRefOids(entityValue *models.RequestCacheEntityValue, entityNewMap map[string][]string, entityOidMap map[string]int) {
	if refOids, b := entityNewMap[entityValue.Oid]; b {
		entityValue.PreviousOids = append(entityValue.PreviousOids, refOids...)
	}
	for tmpOid, refOids := range entityNewMap {
		inRefFlag := false
		for _, refOid := range refOids {
			if entityValue.Oid == refOid {
				inRefFlag = true
			}
		}
		if inRefFlag {
			entityValue.SucceedingOids = append(entityValue.SucceedingOids, tmpOid)
		}
	}
	entityValue.PreviousOids = listToSet(entityValue.PreviousOids, entityOidMap)
	entityValue.SucceedingOids = listToSet(entityValue.SucceedingOids, entityOidMap)
}

func findEntityRefByItemRef(entityValue *models.RequestCacheEntityValue, entityRefs []string, entityNewMap map[string][]string, dataIdOidMap map[string]string) {
	if entityValue.EntityDataOp == "create" {
		log.Logger.Debug("findEntityRefByItemRef create", log.String("oid", entityValue.Oid))
		tmpRefOidList := []string{}
		for _, attrValueObj := range entityValue.AttrValues {
			tmpAttrEntity := getEntityNameFromAttrDefId(attrValueObj.AttrDefId, attrValueObj.AttrName)
			for _, entityRef := range entityRefs {
				if tmpAttrEntity == entityRef && attrValueObj.DataType == "ref" {
					tmpV := fmt.Sprintf("%s", attrValueObj.DataValue)
					if strings.Contains(tmpV, ",") {
						tmpRefOidList = append(tmpRefOidList, strings.Split(tmpV, ",")...)
					} else {
						tmpRefOidList = append(tmpRefOidList, tmpV)
					}
				}
			}
		}
		entityNewMap[entityValue.Oid] = tmpRefOidList
	} else {
		log.Logger.Debug("findEntityRefByItemRef exist", log.String("oid", entityValue.Oid), log.String("EntityDataId", entityValue.EntityDataId))
		dataIdOidMap[entityValue.EntityDataId] = entityValue.Oid
		tmpRefOidList := []string{}
		for _, attrValueObj := range entityValue.AttrValues {
			tmpAttrEntity := getEntityNameFromAttrDefId(attrValueObj.AttrDefId, attrValueObj.AttrName)
			for _, entityRef := range entityRefs {
				if tmpAttrEntity == entityRef && attrValueObj.DataType == "ref" {
					valueString := fmt.Sprintf("%s", attrValueObj.DataValue)
					log.Logger.Debug("findEntityRefByItemRef ref", log.String("oid", entityValue.Oid), log.String("valueString", valueString))
					if strings.Contains(valueString, ",") {
						for _, tmpV := range strings.Split(valueString, ",") {
							if strings.HasPrefix(tmpV, "tmp") {
								tmpRefOidList = append(tmpRefOidList, tmpV)
							}
						}
					} else {
						if strings.HasPrefix(valueString, "tmp") {
							tmpRefOidList = append(tmpRefOidList, valueString)
						}
					}
				}
			}
		}
		entityNewMap[entityValue.Oid] = tmpRefOidList
	}
}

func getEntityNameFromAttrDefId(attrDefId, attrName string) string {
	stringSplit := strings.Split(attrDefId, ":")
	if len(stringSplit) == 3 {
		return stringSplit[2]
	}
	return attrName
}

func matchEntityRoot(requestId, userToken, language string, cacheData *models.RequestCacheData) {
	for _, taskNode := range cacheData.TaskNodeBindInfos {
		existFlag := false
		for _, entityValue := range taskNode.BoundEntityValues {
			if entityValue.EntityDataId == cacheData.RootEntityValue.Oid {
				existFlag = true
				cacheData.RootEntityValue.Oid = entityValue.Oid
				cacheData.RootEntityValue.EntityName = entityValue.EntityName
				cacheData.RootEntityValue.EntityDataOp = entityValue.EntityDataOp
				cacheData.RootEntityValue.AttrValues = entityValue.AttrValues
				cacheData.RootEntityValue.PackageName = entityValue.PackageName
				cacheData.RootEntityValue.EntityDisplayName = entityValue.EntityDisplayName
				cacheData.RootEntityValue.BindFlag = entityValue.BindFlag
				cacheData.RootEntityValue.EntityDataId = entityValue.EntityDataId
				cacheData.RootEntityValue.EntityDataState = entityValue.EntityDataState
				cacheData.RootEntityValue.EntityDefId = entityValue.EntityDefId
				cacheData.RootEntityValue.FullEntityDataId = entityValue.FullEntityDataId
				cacheData.RootEntityValue.PreviousOids = entityValue.PreviousOids
				cacheData.RootEntityValue.SucceedingOids = entityValue.SucceedingOids
				break
			}
		}
		if existFlag {
			break
		}
	}
	if cacheData.RootEntityValue.Oid != "" && cacheData.RootEntityValue.EntityName == "" {
		entityQueryResult, entityQueryErr := GetEntityData(requestId, userToken, language)
		if entityQueryErr != nil {
			log.Logger.Error("Try to fill root entity data fail", log.Error(entityQueryErr))
		} else {
			for _, v := range entityQueryResult.Data {
				if cacheData.RootEntityValue.Oid == v.Id {
					cacheData.RootEntityValue.PackageName = v.PackageName
					cacheData.RootEntityValue.EntityName = v.Entity
					cacheData.RootEntityValue.BindFlag = "N"
					cacheData.RootEntityValue.EntityDataId = v.Id
					cacheData.RootEntityValue.Oid = fmt.Sprintf("%s:%s:%s", v.PackageName, v.Entity, v.Id)
					cacheData.RootEntityValue.EntityDisplayName = v.DisplayName
					cacheData.RootEntityValue.Processed = false
					cacheData.RootEntityValue.SucceedingOids = []string{}
					cacheData.RootEntityValue.PreviousOids = []string{}
					break
				}
			}
		}
	}
}

func listToSet(input []string, itemMap map[string]int) []string {
	result := []string{}
	tmpMap := make(map[string]int)
	for _, v := range input {
		if v == "" {
			continue
		}
		if _, b := tmpMap[v]; !b {
			if _, bb := itemMap[v]; bb {
				result = append(result, v)
				tmpMap[v] = 1
			}
		}
	}
	return result
}

func RequestTermination(requestId, operator, userToken, language string) error {
	requestObj, err := GetRequest(requestId)
	if err != nil {
		return err
	}
	param := models.TerminateInstanceParam{ProcInstId: requestObj.ProcInstanceId, ProcInstKey: requestObj.ProcInstanceKey}
	err = GetProcDefService().TerminationsProcDefInstance(param, userToken, language)
	if err != nil {
		return err
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err = dao.X.Exec("update request set status='Termination',updated_by=?,updated_time=? where id=?", operator, nowTime, requestId)
	return err
}

func GetCmdbReferenceData(attrId, userToken string, param models.QueryRequestParam) (result []byte, statusCode int, err error) {
	paramBytes, tmpErr := json.Marshal(param)
	if tmpErr != nil {
		err = fmt.Errorf("json marshal param data fail,%s ", tmpErr.Error())
		return
	}
	req, newReqErr := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/wecmdb/api/v1/ci-data/reference-data/query/%s", models.Config.Wecube.BaseUrl, attrId), bytes.NewReader(paramBytes))
	if newReqErr != nil {
		err = fmt.Errorf("try to new http request fail,%s ", newReqErr.Error())
		return
	}
	req.Header.Set("Authorization", userToken)
	req.Header.Set("Content-Type", "application/json")
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = fmt.Errorf("try to do http request fail,%s ", respErr.Error())
		return
	}
	statusCode = resp.StatusCode
	result, _ = io.ReadAll(resp.Body)
	resp.Body.Close()
	return
}

func GetRequestPreBindData(request models.RequestTable, requestTemplate *models.RequestTemplateTable, userToken, language string) (result models.RequestCacheData, err error) {
	if request.Cache == "" {
		return result, fmt.Errorf("can not find request cache data with id:%s ", request.Id)
	}
	processNodes, processErr := GetProcDefService().GetProcessDefineTaskNodes(requestTemplate, userToken, language, "bind")
	if processErr != nil {
		return result, processErr
	}
	entityDefIdMap := make(map[string]string)
	entityBindMap := make(map[string][]string)
	for _, v := range processNodes {
		var tmpBoundEntities []string
		for _, vv := range v.BoundEntities {
			if vv == nil {
				continue
			}
			entityDefIdMap[fmt.Sprintf("%s:%s", vv.PackageName, vv.Name)] = vv.Id
			tmpBoundEntities = append(tmpBoundEntities, fmt.Sprintf("%s:%s", vv.PackageName, vv.Name))
		}
		if v.TaskCategory != "SUTN" {
			entityBindMap[v.NodeDefId] = tmpBoundEntities
		}
	}
	var dataCache models.RequestPreDataDto
	json.Unmarshal([]byte(request.Cache), &dataCache)
	entityOidMap := make(map[string][]string)
	entityValueMap := make(map[string]*models.RequestCacheEntityValue)
	for i, v := range dataCache.Data {
		for ii, vv := range v.Value {
			if vv.EntityName == "" {
				continue
			}
			tmpValueObj := models.RequestCacheEntityValue{Oid: vv.Id, BindFlag: "Y", EntityName: vv.EntityName, EntityDisplayName: vv.DisplayName, EntityDataOp: vv.EntityDataOp, EntityDataId: vv.DataId, FullEntityDataId: vv.FullDataId, PackageName: vv.PackageName, PreviousOids: vv.PreviousIds, SucceedingOids: vv.SucceedingIds}
			tmpEntityName := fmt.Sprintf("%s:%s", vv.PackageName, vv.EntityName)
			if _, b := entityDefIdMap[tmpEntityName]; b {
				tmpValueObj.EntityDefId = entityDefIdMap[tmpEntityName]
			}
			tmpValueObj.AttrValues = buildEntityValueAttrData(v.Title, vv.EntityData)
			if dataCache.RootEntityId == "" && i == 0 && ii == 0 {
				result.RootEntityValue = tmpValueObj
			}
			entityValueMap[vv.Id] = &tmpValueObj
			if entityOidMap[v.ItemGroup] != nil {
				entityOidMap[v.ItemGroup] = append(entityOidMap[v.ItemGroup], vv.Id)
			} else {
				entityOidMap[v.ItemGroup] = []string{vv.Id}
			}
		}
	}
	var entityNodeBind []*models.EntityNodeBindQueryObj
	dao.X.SQL("select distinct t1.node_def_id,t2.item_group from task_template t1 left join form_template t2 on t2.task_template=t1.id where t1.request_template=? and t1.node_def_id<>''", requestTemplate.Id).Find(&entityNodeBind)
	for _, v := range entityNodeBind {
		if entityBindMap[v.NodeDefId] != nil {
			entityBindMap[v.NodeDefId] = append(entityBindMap[v.NodeDefId], v.ItemGroup)
		} else {
			entityBindMap[v.NodeDefId] = []string{v.ItemGroup}
		}
	}
	for _, v := range processNodes {
		if v.TaskCategory == "" {
			continue
		}
		tmpNodeBindInfo := models.RequestCacheTaskNodeBindObj{NodeId: v.NodeId, NodeDefId: v.NodeDefId, BoundEntityValues: []*models.RequestCacheEntityValue{}}
		if entities, b := entityBindMap[v.NodeDefId]; b {
			for _, entity := range entities {
				if oids, entityExist := entityOidMap[entity]; entityExist {
					for _, oid := range oids {
						if _, oidExist := entityValueMap[oid]; oidExist {
							tmpNodeBindInfo.BoundEntityValues = append(tmpNodeBindInfo.BoundEntityValues, entityValueMap[oid])
						}
					}
				}
			}
		}
		result.TaskNodeBindInfos = append(result.TaskNodeBindInfos, &tmpNodeBindInfo)
	}
	if result.RootEntityValue.Oid == "" {
		for _, taskNode := range result.TaskNodeBindInfos {
			for _, nodeEntity := range taskNode.BoundEntityValues {
				if nodeEntity.EntityDataId == dataCache.RootEntityId {
					result.RootEntityValue = *nodeEntity
					break
				}
			}
			if result.RootEntityValue.Oid != "" {
				break
			}
		}
	}
	return
}

func buildEntityValueAttrData(titles []*models.FormItemTemplateDto, entityData map[string]interface{}) (result []*models.RequestCacheEntityAttrValue) {
	result = []*models.RequestCacheEntityAttrValue{}
	titleMap := make(map[string]*models.FormItemTemplateDto)
	for _, v := range titles {
		titleMap[v.Name] = v
	}
	for k, v := range entityData {
		if vv, b := titleMap[k]; b {
			if strings.EqualFold(vv.Multiple, models.Yes) || strings.EqualFold(vv.Multiple, models.Y) {
				var tmpV []string
				if newV, ok := v.([]interface{}); ok {
					for _, interfaceV := range newV {
						tmpV = append(tmpV, fmt.Sprintf("%s", interfaceV))
					}
					result = append(result, &models.RequestCacheEntityAttrValue{AttrDefId: vv.AttrDefId, AttrName: k, DataType: vv.AttrDefDataType, DataValue: strings.Join(tmpV, ",")})
					continue
				}
			}
			result = append(result, &models.RequestCacheEntityAttrValue{AttrDefId: vv.AttrDefId, AttrName: k, DataType: vv.AttrDefDataType, DataValue: v})
		}
	}
	return
}

func GetRequestDetailV2(requestId, taskId, userToken, language string) (result models.RequestDetail, err error) {
	// get request
	var requests []*models.RequestTable
	var actions, confirmActions []*dao.ExecAction
	var approvalList []*models.TaskTemplateDto
	dao.X.SQL("select * from request where id=?", requestId).Find(&requests)
	if len(requests) == 0 {
		return result, fmt.Errorf("can not find request with id:%s ", requestId)
	}
	if strings.Contains(requests[0].Status, "InProgress") && requests[0].ProcInstanceId != "" {
		newStatus := getInstanceStatus(requests[0].ProcInstanceId, userToken, language)
		if newStatus == "InternallyTerminated" {
			newStatus = "Termination"
		}
		if newStatus != "" && newStatus != requests[0].Status {
			if newStatus == string(models.RequestStatusCompleted) {
				// 编排的完成,并不表示 请求完成
				taskSort := GetTaskService().GenerateTaskOrderByRequestId(requestId)
				confirmActions, _ = GetRequestService().CreateRequestConfirm(*requests[0], taskSort, userToken, language)
				if len(confirmActions) > 0 {
					actions = append(actions, confirmActions...)
				}
			} else {
				actions = append(actions, &dao.ExecAction{Sql: "update request set status=?,updated_time=? where id=?",
					Param: []interface{}{newStatus, time.Now().Format(models.DateTimeFormat), requestId}})
				requests[0].Status = newStatus
			}
		}
		if len(actions) > 0 {
			updateRequestErr := dao.Transaction(actions)
			if updateRequestErr != nil {
				log.Logger.Error("Try to update request status fail", log.Error(updateRequestErr))
			}
		}
	}
	result.Request = getRequestForm(requests[0], taskId, userToken, language)
	if err != nil {
		return
	}
	if requests[0].TaskApprovalCache != "" {
		json.Unmarshal([]byte(requests[0].TaskApprovalCache), &approvalList)
		result.ApprovalList = approvalList
	}
	return
}

func getCMDBAttributes(entity, userToken string) (result []*models.EntityAttributeObj, err error) {
	req, newReqErr := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/wecmdb/api/v1/ci-types-attr/%s/attributes", models.Config.Wecube.BaseUrl, entity), nil)
	if newReqErr != nil {
		err = fmt.Errorf("try to new http request fail,%s ", newReqErr.Error())
		return
	}
	req.Header.Set("Authorization", userToken)
	//req.Header.Set("Content-Type", "application/json")
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = fmt.Errorf("try to do http request fail,%s ", respErr.Error())
		return
	}
	responseBytes, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("request cmdb attribute fail,%s ", string(responseBytes))
		return
	}
	var attrQueryResp models.EntityAttributeQueryResponse
	err = json.Unmarshal(responseBytes, &attrQueryResp)
	if err != nil {
		err = fmt.Errorf("json unmarshal attr response fail,%s ", err.Error())
		return
	}
	result = attrQueryResp.Data
	return
}

func getAttrCat(catId, userToken string) (result []*models.EntityDataObj, err error) {
	result = []*models.EntityDataObj{}
	req, newReqErr := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/wecmdb/api/v1/base-key/categories/%s", models.Config.Wecube.BaseUrl, catId), nil)
	if newReqErr != nil {
		err = fmt.Errorf("try to new http request fail,%s ", newReqErr.Error())
		return
	}
	req.Header.Set("Authorization", userToken)
	//req.Header.Set("Content-Type", "application/json")
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = fmt.Errorf("try to do http request fail,%s ", respErr.Error())
		return
	}
	responseBytes, _ := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("request cmdb categories fail,%s ", string(responseBytes))
		return
	}
	var response models.CMDBCategoriesResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		err = fmt.Errorf("json unmarshal categories response fail,%s ", err.Error())
		return
	}
	for _, v := range response.Data {
		result = append(result, &models.EntityDataObj{Id: v.Code, DisplayName: v.Value})
	}
	return
}

func BuildRequestProcessData(input models.RequestCacheData, preData *models.EntityTreeData) (result models.RequestProcessData) {
	result.ProcDefId = input.ProcDefId
	result.ProcDefKey = input.ProcDefKey
	result.RootEntityOid = input.RootEntityValue.Oid
	result.Entities = []*models.RequestCacheEntityValue{}
	result.Bindings = []*models.RequestProcessTaskNodeBindObj{}
	entityExistMap := make(map[string]int)
	if len(input.TaskNodeBindInfos) > 0 {
		for _, node := range input.TaskNodeBindInfos {
			for _, entity := range node.BoundEntityValues {
				if _, b := entityExistMap[entity.Oid]; !b {
					result.Entities = append(result.Entities, entity)
					entityExistMap[entity.Oid] = 1
				}
				if node.NodeId != "" {
					result.Bindings = append(result.Bindings, &models.RequestProcessTaskNodeBindObj{Oid: entity.Oid, NodeId: node.NodeId, NodeDefId: node.NodeDefId, EntityDataId: entity.EntityDataId, BindFlag: entity.BindFlag})
				}
			}
		}
	}
	if _, b := entityExistMap[input.RootEntityValue.Oid]; !b {
		result.Entities = append(result.Entities, &input.RootEntityValue)
	}
	if len(result.Entities) == 0 {
		tmpEntityValue := models.RequestCacheEntityValue{Oid: result.RootEntityOid, PackageName: "pseudo", EntityName: "pseudo", BindFlag: "N"}
		result.Entities = append(result.Entities, &tmpEntityValue)
	}
	if preData != nil && len(preData.EntityTreeNodes) > 0 {
		for _, preEntity := range preData.EntityTreeNodes {
			existFlag := false
			for _, entity := range result.Entities {
				if entity.Oid == preEntity.Id {
					existFlag = true
					break
				}
			}
			if !existFlag {
				tmpEntity := models.RequestCacheEntityValue{Oid: preEntity.Id, PackageName: preEntity.PackageName, EntityName: preEntity.EntityName, BindFlag: "N", EntityDataId: preEntity.DataId, EntityDisplayName: preEntity.DisplayName, FullEntityDataId: preEntity.FullDataId, PreviousOids: preEntity.PreviousIds, SucceedingOids: preEntity.SucceedingIds}
				result.Entities = append(result.Entities, &tmpEntity)
			}
		}
	}
	return result
}

func AppendUselessEntity(requestTemplateId, userToken, language string, cacheData *models.RequestCacheData, requestId string) (entityDepMap map[string][]string, preData *models.EntityTreeData, err error) {
	entityDepMap = make(map[string][]string)
	if cacheData.RootEntityValue.Oid == "" || strings.HasPrefix(cacheData.RootEntityValue.Oid, "tmp") {
		return entityDepMap, preData, nil
	}
	// get core preview data list
	rootDataId := cacheData.RootEntityValue.Oid
	if splitIndex := strings.LastIndex(rootDataId, ":"); splitIndex > 0 {
		rootDataId = rootDataId[splitIndex+1:]
	}
	if rootDataId == "" {
		err = fmt.Errorf("preview root data id can not empty")
		return
	}
	preData, err = getRequestPreviewCache(requestId)
	if err != nil {
		return
	}
	if preData == nil {
		preData, err = ProcessDataPreview(requestTemplateId, rootDataId, userToken, language)
		if err != nil {
			err = fmt.Errorf("try to get process preview data fail,%s ", err.Error())
			return
		}
	}
	// get binding entity data
	var entityList []*models.RequestCacheEntityValue
	entityExistMap := make(map[string]int)
	for _, v := range cacheData.TaskNodeBindInfos {
		for _, vv := range v.BoundEntityValues {
			if _, b := entityExistMap[vv.Oid]; !b {
				entityList = append(entityList, vv)
				entityExistMap[vv.Oid] = 1
			}
		}
	}
	// preEntityList is other entity data
	var preEntityList []*models.EntityTreeObj
	rootParent := models.RequestCacheEntityAttrValue{}
	var rootSucceeding []string
	for _, v := range preData.EntityTreeNodes {
		if _, b := entityExistMap[v.Id]; !b {
			preEntityList = append(preEntityList, v)
		}
		if v.DataId == cacheData.RootEntityValue.Oid {
			rootSucceeding = v.SucceedingIds
			rootParent.DataType = "ref"
			rootParent.AttrName = v.EntityName
			rootParent.DataValue = v.DataId
		}
	}
	for _, v := range rootSucceeding {
		entityDepMap[v] = []string{fmt.Sprintf("%s", rootParent.DataValue)}
	}
	// preEntityList -> in preData but no int boundValues
	if len(preEntityList) == 0 {
		return entityDepMap, preData, nil
	}
	dependEntityMap := make(map[string]*models.RequestCacheEntityAttrValue)
	log.Logger.Info("getDependEntity", log.StringList("rootSucceeding", rootSucceeding), log.Int("preLen", len(preEntityList)), log.Int("entityLen", len(entityList)))
	// entityList -> in boundValue entity
	getDependEntity(rootSucceeding, rootParent, preEntityList, entityList, dependEntityMap)
	for k, refAttr := range dependEntityMap {
		refDataValue := fmt.Sprintf("%s", refAttr.DataValue)
		if entityDepMap[k] != nil {
			entityDepMap[k] = append(entityDepMap[k], refDataValue)
		} else {
			entityDepMap[k] = []string{refDataValue}
		}
		log.Logger.Info("dependEntityMap", log.String("id", k), log.String("refValue", refDataValue))
	}
	if len(dependEntityMap) > 0 {
		newNode := models.RequestCacheTaskNodeBindObj{NodeId: "", NodeDefId: "", BoundEntityValues: []*models.RequestCacheEntityValue{}}
		for _, v := range preEntityList {
			if refAttr, b := dependEntityMap[v.Id]; b {
				newNode.BoundEntityValues = append(newNode.BoundEntityValues, &models.RequestCacheEntityValue{Oid: v.Id, EntityDataId: v.DataId, PackageName: v.PackageName, EntityName: v.EntityName, EntityDisplayName: v.DisplayName, FullEntityDataId: v.FullDataId, AttrValues: []*models.RequestCacheEntityAttrValue{refAttr}, PreviousOids: v.PreviousIds, SucceedingOids: v.SucceedingIds})
			}
		}
		cacheData.TaskNodeBindInfos = append(cacheData.TaskNodeBindInfos, &newNode)
	}
	return entityDepMap, preData, nil
}

func getDependEntity(succeeding []string, parent models.RequestCacheEntityAttrValue, preEntityList []*models.EntityTreeObj, entityList []*models.RequestCacheEntityValue, dependEntityMap map[string]*models.RequestCacheEntityAttrValue) {
	if len(succeeding) == 0 {
		return
	}
	// check if use by entity
	for _, v := range succeeding {
		tmpRefAttr := models.RequestCacheEntityAttrValue{}
		vDataId := v
		if strings.Contains(vDataId, ":") {
			vDataId = vDataId[strings.LastIndex(vDataId, ":")+1:]
		}
		useFlag := false
		for _, vv := range entityList {
			for _, attr := range vv.AttrValues {
				if attr.DataType == "ref" {
					tmpV := fmt.Sprintf("%s", attr.DataValue)
					if strings.Contains(tmpV, vDataId) {
						tmpRefAttr = parent
						useFlag = true
						break
					}
				}
			}
			if useFlag {
				break
			}
		}
		if !useFlag {
			// find other succeeding
			tmpParent := models.RequestCacheEntityAttrValue{}
			nextSucceeding := []string{}
			for _, vv := range preEntityList {
				if vv.DataId == v {
					tmpParent.DataType = "ref"
					tmpParent.DataValue = v
					tmpParent.AttrName = vv.EntityName
				}
				if listContains(vv.PreviousIds, v) {
					nextSucceeding = append(nextSucceeding, vv.Id)
				}
			}
			if len(nextSucceeding) > 0 {
				getDependEntity(nextSucceeding, tmpParent, preEntityList, entityList, dependEntityMap)
				for _, vv := range nextSucceeding {
					if _, b := dependEntityMap[vv]; b {
						useFlag = true
						break
					}
				}
			}
		}
		if useFlag {
			for _, vv := range preEntityList {
				if vv.Id == v {
					dependEntityMap[v] = &tmpRefAttr
				}
			}
		}
	}
}

func listContains(inputList []string, element string) bool {
	result := false
	for _, v := range inputList {
		if v == element {
			result = true
			break
		}
	}
	return result
}

func CopyRequest(requestId, createdBy string) (result models.RequestTable, err error) {
	parentRequest := &models.RequestTable{}
	var requestTable []*models.RequestTable
	var requestTemplate *models.RequestTemplateTable
	if err = dao.X.SQL("select * from request where id=?", requestId).Find(&requestTable); err != nil {
		return
	}
	if len(requestTable) == 0 {
		err = fmt.Errorf("can not find any request with id:%s ", requestId)
		return
	}
	parentRequest = requestTable[0]
	if requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(parentRequest.RequestTemplate); err != nil {
		return
	}
	if requestTemplate == nil {
		err = fmt.Errorf("requestTemplate not exist")
		return
	}
	// 重新设置请求名称
	parentRequest.Name = fmt.Sprintf("%s-%s-%s", requestTemplate.Name, requestTemplate.OperatorObjType, time.Now().Format("060102150405"))
	// 重新设置期望时间
	d, _ := time.ParseDuration(fmt.Sprintf("%dh", 24*requestTemplate.ExpireDay))
	parentRequest.ExpectTime = time.Now().Add(d).Format(models.DateTimeFormat)
	result = *parentRequest
	var actions []*dao.ExecAction
	nowTime := time.Now().Format(models.DateTimeFormat)
	result.Id = newRequestId()
	requestInsertAction := dao.ExecAction{Sql: "insert into request(id,name,request_template,reporter,emergency,report_role,status," +
		"cache,expire_time,expect_time,handler,created_by,created_time,updated_by,updated_time,parent,type,role,ref_id,ref_type,custom_form_cache,description,task_approval_cache) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
	requestInsertAction.Param = []interface{}{result.Id, parentRequest.Name, parentRequest.RequestTemplate, createdBy, parentRequest.Emergency,
		parentRequest.ReportRole, "Draft", parentRequest.Cache, "", parentRequest.ExpectTime, parentRequest.Handler, createdBy, nowTime, createdBy,
		nowTime, parentRequest.Id, parentRequest.Type, parentRequest.Role, parentRequest.RefId, parentRequest.RefType, parentRequest.CustomFormCache, parentRequest.Description, parentRequest.TaskApprovalCache}
	actions = append(actions, &requestInsertAction)
	// copy attach file
	var attachFileRows []*models.AttachFileTable
	err = dao.X.SQL("select * from attach_file where request=?", requestId).Find(&attachFileRows)
	if err != nil {
		err = fmt.Errorf("query attach file table fail:%s ", err.Error())
		return
	}
	for _, v := range attachFileRows {
		actions = append(actions, &dao.ExecAction{Sql: "insert into attach_file(id,name,s3_bucket_name,s3_key_name,request,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?,?,?)", Param: []interface{}{
			guid.CreateGuid(), v.Name, v.S3BucketName, v.S3KeyName, result.Id, createdBy, nowTime, createdBy, nowTime}})
	}
	err = dao.TransactionWithoutForeignCheck(actions)
	return
}

func GetRequestParent(requestId string) (parentRequestId string, err error) {
	var requestTable []*models.RequestTable
	err = dao.X.SQL("select `parent` from request where id=?", requestId).Find(&requestTable)
	if err != nil {
		return
	}
	if len(requestTable) == 0 {
		err = fmt.Errorf("can not find request with id:%s ", requestId)
		return
	}
	parentRequestId = requestTable[0].Parent
	return
}

func newRequestId() (requestId string) {
	dateString := time.Now().Format("2006-01-02")
	requestId = time.Now().Format("20060102")
	requestIdLock.Lock()
	defer requestIdLock.Unlock()
	result, err := dao.X.QueryString(fmt.Sprintf("select count(1) as num from request where created_time>='%s 00:00:00'", dateString))
	if err != nil {
		log.Logger.Error("try to new request id fail with count table num", log.Error(err))
		requestId = fmt.Sprintf("%s-%s", requestId, guid.CreateGuid())
		return
	}
	countNum, _ := strconv.Atoi(result[0]["num"])
	countNumString := fmt.Sprintf("%d", countNum+1)
	subId := countNumString
	for i := 0; i < 6-len(countNumString); i++ {
		subId = "0" + subId
	}
	requestId = fmt.Sprintf("%s-%s", requestId, subId)
	return
}

func GetRequestHistory(c *gin.Context, requestId string) (result *models.RequestHistory, err error) {
	result = &models.RequestHistory{}
	// 查询 request
	var requests []*models.RequestTable
	err = dao.X.SQL("select * from request where id=?", requestId).Find(&requests)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}
	if len(requests) == 0 {
		return result, fmt.Errorf("requestId: %s is invalid", requestId)
	}
	result.Request = &models.RequestForHistory{
		RequestTable: *requests[0],
	}

	// 查询 task
	var tasks []*models.TaskTable
	err = dao.X.SQL("select * from task where request=? and del_flag=0 order by sort", requestId).Find(&tasks)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}
	log.Logger.Info("history task", log.Int("taskLen", len(tasks)))

	if len(tasks) == 0 {
		return
	}

	taskIds := make([]string, 0, len(tasks))
	taskTmplIdMap := make(map[string]struct{})
	for _, task := range tasks {
		taskIds = append(taskIds, task.Id)
		taskTmplIdMap[task.TaskTemplate] = struct{}{}
	}

	taskTmplIds := make([]string, 0, len(taskTmplIdMap))
	for k := range taskTmplIdMap {
		taskTmplIds = append(taskTmplIds, k)
	}

	// 查询 task template
	var taskTemplates []*models.TaskTemplateTable
	taskTmplIdsFilterSql, taskTmplIdsFilterParams := dao.CreateListParams(taskTmplIds, "")
	err = dao.X.SQL("select * from task_template where id in ("+taskTmplIdsFilterSql+")", taskTmplIdsFilterParams...).Find(&taskTemplates)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}
	taskTmplIdMapInfo := make(map[string]*models.TaskTemplateTable)
	for _, taskTmpl := range taskTemplates {
		taskTmplIdMapInfo[taskTmpl.Id] = taskTmpl
	}

	// 查询 task handle
	var taskHandles []*models.TaskHandleTable
	taskIdsFilterSql, taskIdsFilterParams := dao.CreateListParams(taskIds, "")
	err = dao.X.SQL("select * from task_handle where task in ("+taskIdsFilterSql+") and latest_flag=1 order by updated_time asc", taskIdsFilterParams...).Find(&taskHandles)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}

	var roleDisplayMap = make(map[string]string)
	if roleDisplayMap, err = GetRoleService().GetRoleDisplayName(c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader)); err != nil {
		err = fmt.Errorf("get role display name failed: %s", err.Error())
		return
	}
	taskHandleIds := make([]string, 0, len(taskHandles))
	for _, taskHandle := range taskHandles {
		taskHandleIds = append(taskHandleIds, taskHandle.Id)
		if displayName, isExisted := roleDisplayMap[taskHandle.Role]; isExisted {
			taskHandle.Role = displayName
		}
	}

	// 查询 attach file
	attachFileTaskHandleIdMap := make(map[string][]*models.AttachFileTable)
	if len(taskHandleIds) > 0 {
		attachFiles, tmpErr := GetTaskHandleAttachFileList(taskHandleIds)
		if tmpErr != nil {
			err = tmpErr
			return
		}
		for _, attachFile := range attachFiles {
			if _, isExisted := attachFileTaskHandleIdMap[attachFile.TaskHandle]; !isExisted {
				attachFileTaskHandleIdMap[attachFile.TaskHandle] = []*models.AttachFileTable{}
			}
			attachFileTaskHandleIdMap[attachFile.TaskHandle] = append(attachFileTaskHandleIdMap[attachFile.TaskHandle], attachFile)
		}
	}

	taskIdMapHandle := make(map[string][]*models.TaskHandleForHistory)
	for _, taskHandle := range taskHandles {
		if _, isExisted := taskIdMapHandle[taskHandle.Task]; !isExisted {
			taskIdMapHandle[taskHandle.Task] = []*models.TaskHandleForHistory{}
		}

		attachFiles := make([]*models.AttachFileTable, 0)
		if _, isExisted := attachFileTaskHandleIdMap[taskHandle.Id]; isExisted {
			attachFiles = attachFileTaskHandleIdMap[taskHandle.Id]
		}
		curTaskHandleForHistory := &models.TaskHandleForHistory{
			TaskHandleTable: taskHandle,
			AttachFiles:     attachFiles,
			FilterRule:      make(map[string]interface{}),
		}
		if strings.TrimSpace(taskHandle.TaskHandleTemplate) != "" {
			taskHandleTemplate := models.TaskHandleTemplateTable{}
			if _, err = dao.X.SQL("select * from task_handle_template where id = ? ", taskHandle.TaskHandleTemplate).Get(&taskHandleTemplate); err != nil {
				return
			}
			if strings.TrimSpace(taskHandleTemplate.FilterRule) != "" {
				if err = json.Unmarshal([]byte(taskHandleTemplate.FilterRule), &curTaskHandleForHistory.FilterRule); err != nil {
					log.Logger.Error("GetRequestHistory json Unmarshal err", log.Error(err))
					return
				}
			}
		}
		taskIdMapHandle[taskHandle.Task] = append(taskIdMapHandle[taskHandle.Task], curTaskHandleForHistory)
	}

	uncompletedTasks := make([]string, 0)
	taskForHistoryList := make([]*models.TaskForHistory, 0, len(tasks))
	for _, task := range tasks {
		filterFlag := false
		if task.ConfirmResult == models.TaskConfirmResultUncompleted {
			uncompletedTasks = append(uncompletedTasks, task.Name)
		}
		if task.Type == string(models.TaskTypeImplement) {
			if task.ProcDefId != "" {
				task.Type = models.TaskTypeImplementProcess
			} else {
				task.Type = models.TaskTypeImplementCustom
			}
		}

		nextOptions := make([]string, 0)
		if task.NextOption != "" {
			nextOptions = strings.Split(task.NextOption, ",")
		}

		var handleMode string
		if templateInfo, isExisted := taskTmplIdMapInfo[task.TaskTemplate]; isExisted {
			handleMode = templateInfo.HandleMode
			var taskHandleTemplateList []*models.TaskHandleTemplateTable
			if err = dao.X.SQL("select * from task_handle_template where task_template = ?", templateInfo.Id).Find(&taskHandleTemplateList); err != nil {
				return
			}
			for _, taskHandleTemplate := range taskHandleTemplateList {
				if strings.TrimSpace(taskHandleTemplate.FilterRule) != "" {
					filterFlag = true
					break
				}
			}
		}

		editable := false
		if task.Status != string(models.TaskStatusDone) {
			editable = true
		}

		formData := make([]*models.RequestPreDataTableObj, 0)
		curTaskForHistory := &models.TaskForHistory{
			TaskTable:      *task,
			TaskHandleList: []*models.TaskHandleForHistory{},
			NextOptions:    nextOptions,
			AttachFiles:    []*models.AttachFileTable{},
			HandleMode:     handleMode,
			Editable:       editable,
			FormData:       formData,
			FilterFlag:     filterFlag,
		}
		if _, isExisted := taskIdMapHandle[task.Id]; isExisted {
			taskHandleForHistoryList := taskIdMapHandle[task.Id]
			// 审批未处理的不展示处理时间
			if task.Type == string(models.TaskTypeApprove) && len(taskHandleForHistoryList) > 0 {
				for _, history := range taskHandleForHistoryList {
					taskHandle := history.TaskHandleTable
					// 任务节点没处理,清空 创建和更新时间
					if taskHandle.HandleResult == "" && taskHandle.HandleStatus == string(models.TaskHandleResultTypeUncompleted) {
						taskHandle.UpdatedTime = ""
						taskHandle.CreatedTime = ""
					}
					// 兼容历史编排任务,判断条件选项赋值,操作重新赋值
					if taskHandle.ProcDefResult == "" && strings.TrimSpace(task.ProcDefId) != "" {
						taskHandle.ProcDefResult = taskHandle.HandleResult
						taskHandle.HandleResult = taskHandle.HandleStatus
					}
				}
			}
			curTaskForHistory.TaskHandleList = taskHandleForHistoryList
		}

		formData, err = getTaskFormData(curTaskForHistory)
		if err != nil {
			log.Logger.Error(fmt.Sprintf("get task form data for task: %s error", task.Id), log.Error(err))
			return
		}
		curTaskForHistory.FormData = formData

		taskForHistoryList = append(taskForHistoryList, curTaskForHistory)
	}
	result.Task = filterFormRowByHandleTemplate(taskForHistoryList)
	// 表单数据行 排序
	if len(result.Task) > 0 {
		for _, task := range result.Task {
			if len(task.TaskHandleList) > 0 {
				for _, taskHandle := range task.TaskHandleList {
					if len(taskHandle.FormData) > 0 {
						for _, formData := range taskHandle.FormData {
							if len(formData.Value) > 0 {
								sort.Sort(models.EntityTreeObjSort(formData.Value))
							}
						}
					}
				}
			}
			if len(task.FormData) > 0 {
				for _, formData := range task.FormData {
					if len(formData.Value) > 0 {
						sort.Sort(models.EntityTreeObjSort(formData.Value))
					}
				}
			}
		}
	}
	result.Request.UncompletedTasks = uncompletedTasks
	return
}

func getTaskFormData(taskObj *models.TaskForHistory) (result []*models.RequestPreDataTableObj, err error) {
	result = []*models.RequestPreDataTableObj{}

	// 查询 form template
	var formTemplates []*models.FormTemplateTable
	err = dao.X.SQL("select * from form_template where task_template=?", taskObj.TaskTemplate).Find(&formTemplates)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}
	if len(formTemplates) == 0 {
		log.Logger.Debug(fmt.Sprintf("can not find any form templates with taskTemplate: %s", taskObj.TaskTemplate))
		return
	}

	formTemplateIdMap := make(map[string]*models.FormTemplateTable)
	formTemplateIds := make([]string, 0, len(formTemplates))
	formTemplateRefIds := make([]string, 0, len(formTemplates))
	for _, formTmpl := range formTemplates {
		formTemplateIds = append(formTemplateIds, formTmpl.Id)
		if formTmpl.RefId != "" {
			formTemplateRefIds = append(formTemplateRefIds, formTmpl.RefId)
		}
		formTemplateIdMap[formTmpl.Id] = formTmpl
	}

	var actualFormTemplates []*models.FormTemplateTable
	actualFormTemplateIds := formTemplateIds
	if taskObj.ProcDefId == "" {
		// 非编排任务, 表单form的 form_template 用的数据表单的 form_template,
		// 需要在 form_template 里面通过id -> ref_id，这个 ref_id 才是 form 的 form_template
		actualFormTemplateIds = formTemplateRefIds

		// 查询 form 实际使用的 formTemplate
		actualFormTemplateIdsFilterSql, actualFormTemplateIdsFilterParams := dao.CreateListParams(actualFormTemplateIds, "")
		err = dao.X.SQL("select * from form_template where id in ("+actualFormTemplateIdsFilterSql+")", actualFormTemplateIdsFilterParams...).Find(&actualFormTemplates)
		if err != nil {
			err = exterror.Catch(exterror.New().DatabaseQueryError, err)
			return
		}
		if len(actualFormTemplates) == 0 {
			log.Logger.Debug(fmt.Sprintf("can not find any form templates with actualFormTemplateIds: [%s]", strings.Join(actualFormTemplateIds, ",")))
			return
		}
	} else {
		// 编排任务
		actualFormTemplateIds = make([]string, 0, len(formTemplates))
		for _, formTmpl := range formTemplates {
			if formTmpl.ItemGroupType == string(models.FormItemGroupTypeWorkflow) {
				actualFormTemplateIds = append(actualFormTemplateIds, formTmpl.Id)
			} else if formTmpl.RefId != "" {
				actualFormTemplateIds = append(actualFormTemplateIds, formTmpl.RefId)
			}
		}

		// actualFormTemplates = formTemplates
		// 查询 form 实际使用的 formTemplate
		actualFormTemplateIdsFilterSql, actualFormTemplateIdsFilterParams := dao.CreateListParams(actualFormTemplateIds, "")
		err = dao.X.SQL("select * from form_template where id in ("+actualFormTemplateIdsFilterSql+")", actualFormTemplateIdsFilterParams...).Find(&actualFormTemplates)
		if err != nil {
			err = exterror.Catch(exterror.New().DatabaseQueryError, err)
			return
		}
		if len(actualFormTemplates) == 0 {
			log.Logger.Debug(fmt.Sprintf("can not find any form templates with actualFormTemplateIds: [%s]", strings.Join(actualFormTemplateIds, ",")))
			return
		}
	}

	// 查询 form
	var taskForms []*models.FormTable
	taskFormsParamList := []interface{}{taskObj.Request}
	actualFormTemplateIdsFilterSql, actualFormTemplateIdsFilterParams := dao.CreateListParams(actualFormTemplateIds, "")
	taskFormsParamList = append(taskFormsParamList, actualFormTemplateIdsFilterParams...)
	err = dao.X.SQL("select * from form where request=? and form_template in ("+actualFormTemplateIdsFilterSql+")", taskFormsParamList...).Find(&taskForms)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}
	if len(taskForms) == 0 {
		log.Logger.Info(fmt.Sprintf("can not find any forms with request: %s and formTemplates: [%s]",
			taskObj.Request, strings.Join(actualFormTemplateIds, ",")))
		return
	}

	// 查询 request 的 form item
	taskUpdatedTime := taskObj.UpdatedTime

	var requestFormItems []*models.FormItemTable
	requestFormItemsParamList := []interface{}{taskObj.Request, taskUpdatedTime}
	err = dao.X.SQL("select * from form_item where request = ? AND updated_time <= ? order by updated_time desc", requestFormItemsParamList...).Find(&requestFormItems)

	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}
	if len(requestFormItems) == 0 {
		log.Logger.Info(fmt.Sprintf("can not find any form items with request: %s and updatedTime <= %s",
			taskObj.Request, taskUpdatedTime))
		return
	}

	// 查询 form item template, 注意：此处的 form_template 需使用根据 taskObj.TaskTemplate 过滤出来的 form_tempalte
	var itemTemplates []*models.FormItemTemplateTable
	formTemplateIdsFilterSql, formTemplateIdsFilterParams := dao.CreateListParams(formTemplateIds, "")
	err = dao.X.SQL("select * from form_item_template where form_template in ("+formTemplateIdsFilterSql+") order by item_group,sort", formTemplateIdsFilterParams...).Find(&itemTemplates)

	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}
	if len(itemTemplates) == 0 {
		log.Logger.Info(fmt.Sprintf("can not find any form item templates with formTemplates: [%s]",
			strings.Join(formTemplateIds, ",")))
		return
	}
	formResult := getItemTemplateTitle(itemTemplates)
	result = formResult

	// sort the result order by item_group_sort
	result = sortHistoryResult(result, formTemplateIdMap)

	// 通过筛选 requestFormItems 获取当前 task 的 form items
	taskFormItems := getTaskFormItems(requestFormItems, taskForms)
	if len(taskFormItems) == 0 {
		log.Logger.Info(fmt.Sprintf("can not find any form item for task: %s", taskObj.Id))
		return
	}

	var formIdMapInfo = make(map[string]*models.FormTable)
	var formTemplateIdMapInfo = make(map[string]*models.FormTemplateTable)
	if formIdMapInfo, formTemplateIdMapInfo, err = getFormAndTemplateMapInfo(taskFormItems); err != nil {
		err = fmt.Errorf("get form and form template map info failed: %s", err.Error())
		return
	}
	// sort taskFormItems by form
	sort.Slice(taskFormItems, func(i, j int) bool {
		return taskFormItems[i].Form < taskFormItems[j].Form
	})

	// itemGroup map data_ids
	itemRowMap := make(map[string][]string)
	// data_id map item
	rowItemMap := make(map[string][]*models.FormItemTable)
	for _, item := range taskFormItems {
		itemGroup := ""
		itemDataId := ""
		if tmpForm, isExisted := formIdMapInfo[item.Form]; isExisted {
			itemDataId = tmpForm.DataId
			if tmpFormTemplate, isExisted2 := formTemplateIdMapInfo[tmpForm.FormTemplate]; isExisted2 {
				itemGroup = tmpFormTemplate.ItemGroup
			} else {
				log.Logger.Debug(fmt.Sprintf("can not find itemGroup for formItem: %s", item.Id))
				continue
			}
		} else {
			log.Logger.Debug(fmt.Sprintf("can not find itemDataId for formItem: %s", item.Id))
			continue
		}

		if tmpRows, b := itemRowMap[itemGroup]; b {
			existFlag := false
			for _, v := range tmpRows {
				if itemDataId == v {
					existFlag = true
					break
				}
			}
			if !existFlag {
				itemRowMap[itemGroup] = append(itemRowMap[itemGroup], itemDataId)
			}
		} else {
			itemRowMap[itemGroup] = []string{itemDataId}
		}
		if rowItemMap[itemDataId] != nil {
			rowItemMap[itemDataId] = append(rowItemMap[itemDataId], item)
		} else {
			rowItemMap[itemDataId] = []*models.FormItemTable{item}
		}
	}

	for _, formTable := range formResult {
		if rows, b := itemRowMap[formTable.ItemGroup]; b {
			for _, row := range rows {
				tmpRowObj := models.EntityTreeObj{Id: row, DataId: row, PackageName: formTable.PackageName, EntityName: formTable.Entity}
				tmpRowObj.EntityData = make(map[string]interface{})
				for _, rowItem := range rowItemMap[row] {
					isMulti := false
					isMultiObject := false
					for _, tmpTitle := range formTable.Title {
						if tmpTitle.Name == rowItem.Name {
							if strings.EqualFold(tmpTitle.Multiple, models.Yes) || strings.EqualFold(tmpTitle.Multiple, models.Y) {
								isMulti = true
								if tmpTitle.AttrDefDataType == string(models.CmdbDataTypeMultiObject) {
									isMultiObject = true
								}
								break
							}
						}
					}
					if isMulti {
						if strings.TrimSpace(rowItem.Value) == "" {
							tmpRowObj.EntityData[rowItem.Name] = []string{}
						} else if isMultiObject {
							// 对象数组
							var jsonData interface{}
							json.Unmarshal([]byte(rowItem.Value), &jsonData)
							tmpRowObj.EntityData[rowItem.Name] = jsonData
						} else {
							tmpRowObj.EntityData[rowItem.Name] = strings.Split(rowItem.Value, ",")
						}
					} else {
						tmpRowObj.EntityData[rowItem.Name] = rowItem.Value
					}
				}
				// 可能当前表单项为当前表单模版独有表单项,需要加特殊处理
				for _, formItemTemp := range formTable.Title {
					if _, ok := tmpRowObj.EntityData[formItemTemp.Name]; !ok {
						tmpRowObj.EntityData[formItemTemp.Name] = formItemTemp.DefaultValue
					}
				}
				formTable.Value = append(formTable.Value, &tmpRowObj)
			}
		}
	}
	return
}

func sortHistoryResult(historyResult []*models.RequestPreDataTableObj, formTemplateIdMap map[string]*models.FormTemplateTable) (result []*models.RequestPreDataTableObj) {
	result = historyResult
	if len(historyResult) < 2 || formTemplateIdMap == nil {
		return
	}

	historyResultToSortList := make([]*models.HistoryResultToSort, 0, len(result))
	for _, resultData := range historyResult {
		itemGroupSort := 0
		if formTmpl, isExisted := formTemplateIdMap[resultData.FormTemplateId]; isExisted {
			itemGroupSort = formTmpl.ItemGroupSort
		}

		toSortElem := &models.HistoryResultToSort{
			ItemGroupSort:     itemGroupSort,
			HistoryResultElem: resultData,
		}
		historyResultToSortList = append(historyResultToSortList, toSortElem)
	}
	sort.Slice(historyResultToSortList, func(i int, j int) bool {
		return historyResultToSortList[i].ItemGroupSort < historyResultToSortList[j].ItemGroupSort
	})

	result = make([]*models.RequestPreDataTableObj, 0, len(historyResult))
	for _, resultElem := range historyResultToSortList {
		result = append(result, resultElem.HistoryResultElem)
	}
	return
}

func getFormAndTemplateMapInfo(taskFormItems []*models.FormItemTable) (formIdMapInfo map[string]*models.FormTable, formTemplateIdMapInfo map[string]*models.FormTemplateTable, err error) {
	formIdMapInfo = make(map[string]*models.FormTable)
	formTemplateIdMapInfo = make(map[string]*models.FormTemplateTable)

	formIdMap := make(map[string]struct{})
	for _, formItem := range taskFormItems {
		if formItem.Form != "" {
			formIdMap[formItem.Form] = struct{}{}
		}
	}

	if len(formIdMap) > 0 {
		formIds := make([]string, 0, len(formIdMap))
		for k := range formIdMap {
			formIds = append(formIds, k)
		}

		var forms []*models.FormTable
		formIdsFilterSql, formIdsFilterParams := dao.CreateListParams(formIds, "")
		err = dao.X.SQL("select * from form where id in ("+formIdsFilterSql+")", formIdsFilterParams...).Find(&forms)
		if err != nil {
			err = exterror.Catch(exterror.New().DatabaseQueryError, err)
			return
		}

		if len(forms) > 0 {
			formTemplateIdMap := make(map[string]struct{})
			for _, form := range forms {
				if form.FormTemplate != "" {
					formTemplateIdMap[form.FormTemplate] = struct{}{}
				}
				formIdMapInfo[form.Id] = form
			}

			if len(formTemplateIdMap) > 0 {
				formTemplateIds := make([]string, 0, len(formTemplateIdMap))
				for k := range formTemplateIdMap {
					formTemplateIds = append(formTemplateIds, k)
				}

				var formTemplates []*models.FormTemplateTable
				formTemplateIdsFilterSql, formTemplateIdsFilterParams := dao.CreateListParams(formTemplateIds, "")
				err = dao.X.SQL("select * from form_template where id in ("+formTemplateIdsFilterSql+")", formTemplateIdsFilterParams...).Find(&formTemplates)
				if err != nil {
					err = exterror.Catch(exterror.New().DatabaseQueryError, err)
					return
				}

				for _, formTemplate := range formTemplates {
					formTemplateIdMapInfo[formTemplate.Id] = formTemplate
				}
			}
		}
	}
	return
}

func getTaskFormItems(requestFormItems []*models.FormItemTable, taskForms []*models.FormTable) (taskFormItems []*models.FormItemTable) {
	taskFormItems = []*models.FormItemTable{}
	if len(requestFormItems) == 0 {
		return
	}
	if len(taskForms) == 0 {
		return
	}

	var distinctFormItems []*models.FormItemTable
	formNameMap := make(map[string]struct{})
	for _, item := range requestFormItems {
		tmpKey := fmt.Sprintf("%s__%s", item.Form, item.Name)
		if _, isExisted := formNameMap[tmpKey]; !isExisted {
			distinctFormItems = append(distinctFormItems, item)
			formNameMap[tmpKey] = struct{}{}
		}
	}

	taskFormIdMap := make(map[string]struct{})
	for _, taskForm := range taskForms {
		taskFormIdMap[taskForm.Id] = struct{}{}
	}

	for _, item := range distinctFormItems {
		if _, isExisted := taskFormIdMap[item.Form]; isExisted {
			taskFormItems = append(taskFormItems, item)
		}
	}
	return
}

func SaveRequestForm(requestId, operator string, param *models.RequestPreDataTableObj) (err error) {
	actions, buildActionErr := UpdateSingleRequestFormNew(requestId, operator, param)
	if buildActionErr != nil {
		err = buildActionErr
		return
	}
	if len(actions) > 0 {
		err = dao.Transaction(actions)
	}
	return
}

// 根据模版审批任务配置审批节点,filter_rule规则进行过滤
func filterFormRowByHandleTemplate(taskHistoryList []*models.TaskForHistory) []*models.TaskForHistory {
	var newTaskHistoryList []*models.TaskForHistory
	var taskHandleList []*models.TaskHandleForHistory
	var data = make(map[string]interface{})
	var itemGroup string
	var deleteRowIdMap = make(map[string]bool)
	// 根据任务处理模版过滤表单内容
	if len(taskHistoryList) > 0 {
		for _, taskHistory := range taskHistoryList {
			data = make(map[string]interface{})
			taskHandleList = []*models.TaskHandleForHistory{}
			if len(taskHistory.TaskHandleList) > 0 {
				for _, taskHandle := range taskHistory.TaskHandleList {
					itemGroup = ""
					deleteRowIdMap = make(map[string]bool)
					if len(taskHandle.FilterRule) > 0 && len(taskHistory.FormData) > 0 {
						for _, formData := range taskHistory.FormData {
							if len(formData.Title) > 0 && len(formData.Value) > 0 {
								itemGroup = formData.Title[0].ItemGroup
								for _, entity := range formData.Value {
									if len(entity.EntityData) > 0 {
										for key, value := range taskHandle.FilterRule {
											if !strings.HasPrefix(key, itemGroup+"-") {
												continue
											}
											name := strings.Replace(key, itemGroup+"-", "", 1)
											if _, ok := entity.EntityData[name]; ok {
												valueStr, ok := value.(string)
												if ok {
													if len(valueStr) == 0 {
														continue
													}
													// 单选判断,不相等直接过滤
													if valueStr != entity.EntityData[name] {
														deleteRowIdMap[entity.Id] = true
														break
													}
												} else {
													// 多选判断,都不满足才过滤
													exist := false
													filterArr, ok1 := value.([]interface{})
													if !ok1 {
														log.Logger.Error("data value  is not array", log.JsonObj("data", data))
														continue
													}
													if len(filterArr) == 0 {
														continue
													}
													// entity.EntityData里面的value 可能是[]string也有可能是 a,b形式,取决是否展示,展示前面就会处理,这个地方要兼容两种格式
													entityArr, ok2 := entity.EntityData[name].([]string)
													if !ok2 {
														str, ok3 := entity.EntityData[name].(string)
														if !ok3 {
															log.Logger.Error("entity.EntityData value is not string", log.JsonObj("data", entity.EntityData[name]))
															continue
														}
														entityArr = strings.Split(str, ",")
													}
													filterMap := convertInterfaceArray2Map(filterArr)
													for _, val := range entityArr {
														if filterMap[val] {
															exist = true
															break
														}
													}
													if !exist {
														deleteRowIdMap[entity.Id] = true
													}
												}
											}
										}
									}
								}
							}
						}
						var formDataList []*models.RequestPreDataTableObj
						for _, formData := range taskHistory.FormData {
							var valueList []*models.EntityTreeObj
							if len(formData.Value) > 0 {
								for _, entity := range formData.Value {
									if !deleteRowIdMap[entity.Id] {
										valueList = append(valueList, entity)
									}
								}
							}
							newFormData := &models.RequestPreDataTableObj{
								PackageName:    formData.PackageName,
								Entity:         formData.Entity,
								FormTemplateId: formData.FormTemplateId,
								ItemGroup:      formData.ItemGroup,
								ItemGroupName:  formData.ItemGroupName,
								ItemGroupType:  formData.ItemGroupType,
								ItemGroupRule:  formData.ItemGroupRule,
								RefEntity:      formData.RefEntity,
								SortLevel:      formData.SortLevel,
								Title:          formData.Title,
								Value:          valueList,
							}
							formDataList = append(formDataList, newFormData)
						}
						taskHandle.FormData = formDataDeepCopy(formDataList)
					} else {
						taskHandle.FormData = formDataDeepCopy(taskHistory.FormData)
					}
					// 当前任务已经执行完成,并行每个结点只能看到自己修改数据,而不是当前任务最新数据. 协同没有提交的,表单数据展示为空
					if taskHistory.Status == string(models.TaskStatusDone) {
						if taskHistory.HandleMode == string(models.TaskTemplateHandleModeAny) && strings.TrimSpace(taskHandle.HandleResult) == "" {
							taskHandle.FormData = nil
						} else if taskHistory.HandleMode == string(models.TaskTemplateHandleModeAll) && taskHistory.Request != "" {
							// 并行审批,直接读取 task_handle表的 from_data数据
							if strings.TrimSpace(taskHandle.HandleFormData) != "" {
								err := json.Unmarshal([]byte(taskHandle.HandleFormData), &taskHandle.FormData)
								if err != nil {
									log.Logger.Error("json Unmarshal err:%+v", log.Error(err))
								}
							}
						}
					}
					taskHandleList = append(taskHandleList, taskHandle)
				}
			}
			newTaskHistoryList = append(newTaskHistoryList, &models.TaskForHistory{
				TaskTable:      taskHistory.TaskTable,
				Editable:       taskHistory.Editable,
				TaskHandleList: taskHandleList,
				NextOptions:    taskHistory.NextOptions,
				AttachFiles:    taskHistory.AttachFiles,
				HandleMode:     taskHistory.HandleMode,
				FormData:       taskHistory.FormData,
				FilterFlag:     taskHistory.FilterFlag,
			})
		}
	}
	return newTaskHistoryList
}

func getRequestPreviewCache(requestId string) (result *models.EntityTreeData, err error) {
	var requestRows []*models.RequestTable
	err = dao.X.SQL("select preview_cache from request where id=?", requestId).Find(&requestRows)
	if err != nil {
		err = fmt.Errorf("query request preview cache data fail,%s ", err.Error())
		return
	}
	if len(requestRows) == 0 {
		err = fmt.Errorf("can not find request row with id:%s ", requestId)
		return
	}
	if requestRows[0].PreviewCache == "" {
		return
	}
	result = &models.EntityTreeData{}
	if err = json.Unmarshal([]byte(requestRows[0].PreviewCache), result); err != nil {
		err = fmt.Errorf("json unmarshal request preview cache data fail,%s ", err.Error())
	}
	return
}

func formDataDeepCopy(dataList []*models.RequestPreDataTableObj) []*models.RequestPreDataTableObj {
	var list []*models.RequestPreDataTableObj
	var valueList []*models.EntityTreeObj
	for _, data := range dataList {
		valueList = []*models.EntityTreeObj{}
		if len(data.Value) > 0 {
			for _, obj := range data.Value {
				entityData := make(map[string]interface{})
				if len(obj.EntityData) > 0 {
					for key, value := range obj.EntityData {
						entityData[key] = value
					}
				}
				valueList = append(valueList, &models.EntityTreeObj{
					PackageName:   obj.PackageName,
					EntityName:    obj.EntityName,
					DataId:        obj.DataId,
					DisplayName:   obj.DisplayName,
					FullDataId:    obj.FullDataId,
					Id:            obj.Id,
					EntityData:    entityData,
					PreviousIds:   obj.PreviousIds,
					SucceedingIds: obj.SucceedingIds,
					EntityDataOp:  obj.EntityDataOp,
				})
			}
		}
		list = append(list, &models.RequestPreDataTableObj{
			PackageName:    data.PackageName,
			Entity:         data.Entity,
			FormTemplateId: data.FormTemplateId,
			ItemGroup:      data.ItemGroup,
			ItemGroupName:  data.ItemGroupName,
			ItemGroupType:  data.ItemGroupType,
			ItemGroupRule:  data.ItemGroupRule,
			RefEntity:      data.RefEntity,
			SortLevel:      data.SortLevel,
			Title:          data.Title,
			Value:          valueList,
		})
	}
	return list
}

// HandleFormSensitiveData 处理表单敏感数据
func HandleFormSensitiveData(originData []map[string]interface{}, entity, token string) (targetData []map[string]interface{}, err error) {
	if len(originData) == 0 {
		originData = make([]map[string]interface{}, 0)
	}
	targetData = make([]map[string]interface{}, 0)
	var refAttributes []*models.EntityAttributeObj
	if refAttributes, err = GetCMDBCiAttrDefs(entity, token); err != nil {
		err = fmt.Errorf("query remote entity:%s attr fail:%s ", entity, err.Error())
		return
	}
	for _, dataMap := range originData {
		for _, v1 := range refAttributes {
			if strings.ToUpper(v1.Sensitive) == "YES" || strings.ToUpper(v1.Sensitive) == "Y" {
				if v, ok := dataMap[v1.PropertyName]; ok && v != "" {
					dataMap[v1.PropertyName] = models.SensitiveStyle
				}
			}
		}
		targetData = append(targetData, dataMap)
	}
	return
}
