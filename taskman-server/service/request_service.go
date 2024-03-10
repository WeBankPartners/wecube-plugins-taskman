package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	requestDao *dao.RequestDao
}

var (
	requestIdLock   = new(sync.RWMutex)
	templateTypeArr = []int{int(models.SceneTypeRelease), int(models.SceneTypeRequest)} // 模版类型: 1表示请求,2表示发布
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
	return
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
	nowTime := time.Now().Format(models.DateTimeFormat)
	request, err = GetSimpleRequest(requestId)
	if err != nil {
		return
	}
	if request.Id == "" {
		err = fmt.Errorf("param requestId:%s not exist", requestId)
		return
	}
	if request.Status != "Pending" {
		err = exterror.New().RevokeRequestError
		return
	}
	if request.CreatedBy != user {
		err = fmt.Errorf("request not yours")
		return
	}
	_, err = dao.X.Exec("update request set status ='Draft',revoke_flag=1,updated_time=? where id=?", nowTime, request.Id)
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
		err = fmt.Errorf("Can not find any request with id:%s ", requestId)
		return
	}
	result = *requestTable[0]
	if result.Cache != "" {
		var cacheObj models.RequestPreDataDto
		err = json.Unmarshal([]byte(result.Cache), &cacheObj)
		if err != nil {
			err = fmt.Errorf("Try to json unmarshal cache data fail,%s ", err.Error())
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
	err = dao.X.SQL("select id,name,form,request_template,proc_instance_id,proc_instance_key,reporter,report_time,emergency,status from request where id=?", requestId).Find(&requestTable)
	if err != nil {
		return
	}
	if len(requestTable) == 0 {
		err = fmt.Errorf("Can not find any request with id:%s ", requestId)
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
			return fmt.Errorf("Try to sync proDefId fail,%s ", err.Error())
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
	err := ValidateRequestForm(param.Data, userToken)
	if err != nil {
		return err
	}
	var formItemNameQuery []*models.FormItemTemplateTable
	err = dao.X.SQL("select item_group,group_concat(name,',') as name from form_item_template where in_display_name='yes' and form_template in (select form_template from request_template where id in (select request_template from request where id=?)) group by item_group", requestId).Find(&formItemNameQuery)
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
		return fmt.Errorf("Try to json marshal param fail,%s ", err.Error())
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	actions := UpdateRequestFormItem(requestId, operator, nowTime, param)
	actions = append(actions, &dao.ExecAction{Sql: "update request set cache=?,updated_by=?,updated_time=?,operator_obj=? where id=?", Param: []interface{}{string(paramBytes), operator, nowTime, param.EntityName, requestId}})
	return dao.Transaction(actions)
}

func SaveRequestCacheV2(requestId, operator, userToken string, param *models.RequestProDataV2Dto) error {
	var customFormCache []byte
	var taskApprovalCache string
	err := ValidateRequestForm(param.Data, userToken)
	if err != nil {
		return err
	}
	var formItemNameQuery []*models.FormItemTemplateTable
	err = dao.X.SQL("select item_group,group_concat(name,',') as name from form_item_template where in_display_name='yes' and form_template in (select request_template from form_template where id in (select request_template from request where id=?)) group by item_group", requestId).Find(&formItemNameQuery)
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
	newParam := &models.RequestPreDataDto{
		RootEntityId: param.RootEntityId,
		Data:         param.Data,
	}
	paramBytes, err := json.Marshal(newParam)
	if err != nil {
		return fmt.Errorf("Try to json marshal param fail,%s ", err.Error())
	}
	if len(param.CustomForm.Value) > 0 {
		customFormCache, err = json.Marshal(param.CustomForm)
		if err != nil {
			return fmt.Errorf("Try to json marshal param fail,%s ", err.Error())
		}
	}
	if len(param.ApprovalList) > 0 {
		approvalBytes, _ := json.Marshal(param.ApprovalList)
		if approvalBytes != nil {
			taskApprovalCache = string(approvalBytes)
		}
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	actions := UpdateRequestFormItem(requestId, operator, nowTime, newParam)
	actions = append(actions, &dao.ExecAction{Sql: "update request set cache=?,updated_by=?,updated_time=?,name=?,description=?,expect_time=?,operator_obj=?,custom_form_cache=?,task_approval_cache=?" +
		" where id=?", Param: []interface{}{string(paramBytes), operator, nowTime, param.Name, param.Description, param.ExpectTime, param.EntityName, string(customFormCache), taskApprovalCache, requestId}})
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
		for _, valueObj := range v.Value {
			// 每行记录新增一个表单
			newFormId = guid.CreateGuid()
			actions = append(actions, &dao.ExecAction{Sql: "insert into form(id,request,form_template,data_id,created_by,created_time,updated_by," +
				"updated_time) values (?,?,?,?,?,?,?,?)", Param: []interface{}{newFormId, requestId, v.FormTemplateId, valueObj.Id, operator, now, operator, now}})
			tmpGuidList := guid.CreateGuidList(len(v.Title))
			for i, title := range v.Title {
				if title.Multiple == "Y" {
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

func GetRequestCache(requestId, cacheType string) (result interface{}, err error) {
	var requestTable []*models.RequestTable
	if cacheType == "data" {
		err = dao.X.SQL("select cache from request where id=?", requestId).Find(&requestTable)
		if err != nil {
			return
		}
		if len(requestTable) == 0 {
			err = fmt.Errorf("Can not find any request with id:%s ", requestId)
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
			err = fmt.Errorf("Can not find any request with id:%s ", requestId)
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
		err = fmt.Errorf("Can not find request with id:%s ", requestId)
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

func GetRequestPreData(requestId, entityDataId, userToken, language string) (result []*models.RequestPreDataTableObj, err error) {
	var requestTables []*models.RequestTable
	err = dao.X.SQL("select cache from request where id=?", requestId).Find(&requestTables)
	if err != nil {
		return
	}
	if len(requestTables) == 0 {
		return result, fmt.Errorf("Can not find requestId:%s ", requestId)
	}
	if requestTables[0].Cache != "" {
		var cacheObj models.RequestPreDataDto
		err = json.Unmarshal([]byte(requestTables[0].Cache), &cacheObj)
		if err != nil {
			return result, fmt.Errorf("Try to json unmarshal cache data fail,%s ", err.Error())
		}
		if cacheObj.RootEntityId == entityDataId {
			result = cacheObj.Data
			return
		}
	}
	result = []*models.RequestPreDataTableObj{}
	requestTemplateId, tmpErr := getRequestTemplateByRequest(requestId)
	if tmpErr != nil {
		return result, tmpErr
	}
	result = getRequestPreDataByTemplateId(requestTemplateId)
	if entityDataId == "" {
		return
	}
	previewData, previewErr := ProcessDataPreview(requestTemplateId, entityDataId, userToken, language)
	if previewErr != nil {
		return result, previewErr
	}
	if len(previewData.EntityTreeNodes) == 0 {
		return
	}
	for _, entity := range result {
		for _, tmpData := range previewData.EntityTreeNodes {
			if tmpData.EntityName == entity.Entity {
				tmpValueData := make(map[string]interface{})
				for _, title := range entity.Title {
					tmpValueData[title.Name] = tmpData.EntityData[title.Name]
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
	var itemGroupType, itemGroupRule string
	var formTemplateId string
	tmpPackageName := items[0].PackageName
	tmpEntity := items[0].Entity
	tmpItemGroup := items[0].ItemGroup
	tmpItemGroupName := items[0].ItemGroupName
	var tmpRefEntity []string
	var tmpItems []*models.FormItemTemplateDto
	existItemMap := make(map[string]int)
	for _, v := range items {
		tmpKey := fmt.Sprintf("%s__%s", v.ItemGroup, v.Name)
		itemGroup := models.FormTemplateTable{}
		if v.FormTemplate != "" {
			dao.X.SQL("select * from form_template where id=?", v.FormTemplate).Get(&itemGroup)
			itemGroupType = itemGroup.ItemGroupType
			itemGroupRule = itemGroup.ItemGroupRule
			formTemplateId = itemGroup.Id
		}
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
					FormTemplateId: formTemplateId,
					ItemGroup:      tmpItemGroup,
					ItemGroupName:  tmpItemGroupName,
					ItemGroupType:  itemGroupType,
					ItemGroupRule:  itemGroupRule,
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
			tmpRefEntity = []string{}
		} else {
			if tmpEntity == "" && v.Entity != "" {
				tmpEntity = v.Entity
				tmpPackageName = v.PackageName
			}
		}
		tmpItems = append(tmpItems, models.ConvertFormItemTemplateModel2Dto(v, itemGroup))
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
			FormTemplateId: formTemplateId,
			ItemGroup:      tmpItemGroup,
			ItemGroupName:  tmpItemGroupName,
			ItemGroupType:  itemGroupType,
			ItemGroupRule:  itemGroupRule,
			RefEntity:      tmpRefEntity,
			SortLevel:      0,
			Title:          tmpItems,
			Value:          []*models.EntityTreeObj{},
		})
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
	entityMap := make(map[string][]string)
	entityNumMap := make(map[string]int)
	for _, v := range param {
		for _, vv := range v.RefEntity {
			if vv == v.Entity {
				continue
			}
			if _, b := entityMap[vv]; b {
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
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(request.RequestTemplate)
	if err != nil {
		return
	}
	if requestTemplate == nil {
		return
	}
	entityDepMap, err = AppendUselessEntity(requestTemplate.Id, userToken, language, &cacheData)
	if err != nil {
		err = fmt.Errorf("Try to append useless entity fail,%s ", err.Error())
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
	actions = append(actions, &dao.ExecAction{Sql: "update task_handle set handle_result=?,updated_time=? where id=?",
		Param: []interface{}{models.TaskHandleResultTypeApprove, nowTime, checkTaskHandle.Id}})
	// 更新任务为完成
	actions = append(actions, &dao.ExecAction{Sql: "update task set status=?,updated_by=?,updated_time=? where id=?",
		Param: []interface{}{models.TaskStatusDone, operator, nowTime, task.Id}})
	approvalActions, err = GetRequestService().CreateRequestApproval(request, "", userToken, language)
	if err != nil {
		return
	}
	if len(approvalActions) > 0 {
		actions = append(actions, approvalActions...)
	}
	err = dao.Transaction(actions)
	return
}

func StartRequest(requestId, operator, userToken, language string, cacheData models.RequestCacheData) (result *models.StartInstanceResultData, err error) {
	var requestTemplateTable []*models.RequestTemplateTable
	dao.X.SQL("select * from request_template where id in (select request_template from request where id=?)", requestId).Find(&requestTemplateTable)
	if len(requestTemplateTable) == 0 {
		return result, fmt.Errorf("Can not find requestTemplate with request:%s ", requestId)
	}
	cacheData.ProcDefId = requestTemplateTable[0].ProcDefId
	cacheData.ProcDefKey = requestTemplateTable[0].ProcDefKey
	entityDepMap, tmpErr := AppendUselessEntity(requestTemplateTable[0].Id, userToken, language, &cacheData)
	if tmpErr != nil {
		return result, fmt.Errorf("Try to append useless entity fail,%s ", tmpErr.Error())
	}
	fillBindingWithRequestData(requestId, userToken, language, &cacheData, entityDepMap)
	cacheBytes, _ := json.Marshal(cacheData)
	log.Logger.Info("cacheByte", log.String("cacheBytes", string(cacheBytes)))
	startParam := BuildRequestProcessData(cacheData)
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
	_, err = dao.X.Exec("update request set handler=?,proc_instance_id=?,proc_instance_key=?,confirm_time=?,expire_time=?,status=?,bind_cache=?,updated_by=?,updated_time=? where id=?", operator, procInstId, result.ProcInstKey, nowTime, expireTime, result.Status, string(cacheBytes), operator, nowTime, requestId)
	return
}

func StartRequestNew(request models.RequestTable, userToken, language string, cacheData models.RequestCacheData) (actions []*dao.ExecAction, err error) {
	var requestTemplateTable []*models.RequestTemplateTable
	var result *models.StartInstanceResultData
	actions = []*dao.ExecAction{}
	dao.X.SQL("select * from request_template where  id=?)", request.RequestTemplate).Find(&requestTemplateTable)
	if len(requestTemplateTable) == 0 {
		err = fmt.Errorf("Can not find requestTemplate with request:%s ", request.Id)
		return
	}
	cacheData.ProcDefId = requestTemplateTable[0].ProcDefId
	cacheData.ProcDefKey = requestTemplateTable[0].ProcDefKey
	entityDepMap, tmpErr := AppendUselessEntity(requestTemplateTable[0].Id, userToken, language, &cacheData)
	if tmpErr != nil {
		err = fmt.Errorf("Try to append useless entity fail,%s ", tmpErr.Error())
		return
	}
	fillBindingWithRequestData(request.Id, userToken, language, &cacheData, entityDepMap)
	startParam := BuildRequestProcessData(cacheData)
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
	actions = append(actions, &dao.ExecAction{Sql: "update request set proc_instance_id=?,proc_instance_key=?,status=?,updated_time=? " +
		"where id=?", Param: []interface{}{procInstId, result.ProcInstKey, result.Status, nowTime, request.Id}})
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
				return fmt.Errorf("Try to build bind data fail,%s ", bindErr.Error())
			}
			bindCacheBytes, _ := json.Marshal(bindData)
			bindCache = string(bindCacheBytes)
		}
		// 请求定版, 根据模板配置开启是否确认定版
		err = GetRequestService().CreateRequestCheck(request, operator, bindCache, userToken, language)
	} else if status == "Draft" {
		if request.Handler != operator {
			err = exterror.New().UpdateRequestHandlerStatusError
			return err
		}
		_, err = dao.X.Exec("update request set status=?,rollback_desc=?,updated_by=?,handler=?,updated_time=?,confirm_time=? where id=?", status, description, operator, operator, nowTime, nowTime, requestId)
	} else {
		_, err = dao.X.Exec("update request set status=?,updated_by=?,updated_time=? where id=?", status, operator, nowTime, requestId)
	}
	return err
}

func fillBindingWithRequestData(requestId, userToken, language string, cacheData *models.RequestCacheData, existDepMap map[string][]string) {
	var items []*models.FormItemTemplateTable
	dao.X.SQL("select * from form_item_template where form_template in (select form_template from request_template where id in (select request_template from request where id=?)) order by entity,sort", requestId).Find(&items)
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
		err = fmt.Errorf("Json marshal param data fail,%s ", tmpErr.Error())
		return
	}
	req, newReqErr := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/wecmdb/api/v1/ci-data/reference-data/query/%s", models.Config.Wecube.BaseUrl, attrId), bytes.NewReader(paramBytes))
	if newReqErr != nil {
		err = fmt.Errorf("Try to new http request fail,%s ", newReqErr.Error())
		return
	}
	req.Header.Set("Authorization", userToken)
	req.Header.Set("Content-Type", "application/json")
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = fmt.Errorf("Try to do http request fail,%s ", respErr.Error())
		return
	}
	statusCode = resp.StatusCode
	result, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return
}

func GetRequestPreBindData(request models.RequestTable, requestTemplate *models.RequestTemplateTable, userToken, language string) (result models.RequestCacheData, err error) {
	if request.Cache == "" {
		return result, fmt.Errorf("Can not find request cache data with id:%s ", request.Id)
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
			if _, b := entityOidMap[v.ItemGroup]; b {
				entityOidMap[v.ItemGroup] = append(entityOidMap[v.ItemGroup], vv.Id)
			} else {
				entityOidMap[v.ItemGroup] = []string{vv.Id}
			}
		}
	}
	var entityNodeBind []*models.EntityNodeBindQueryObj
	dao.X.SQL("select distinct t1.node_def_id,t2.item_group from task_template t1 left join form_template t2 on t2.task_template=t1.id where t1.request_template=? and t1.node_def_id<>''", requestTemplate.Id).Find(&entityNodeBind)
	for _, v := range entityNodeBind {
		if _, b := entityBindMap[v.NodeDefId]; b {
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
			if vv.Multiple == "Y" {
				var tmpV []string
				for _, interfaceV := range v.([]interface{}) {
					tmpV = append(tmpV, fmt.Sprintf("%s", interfaceV))
				}
				result = append(result, &models.RequestCacheEntityAttrValue{AttrDefId: vv.AttrDefId, AttrName: k, DataType: vv.AttrDefDataType, DataValue: strings.Join(tmpV, ",")})
			} else {
				result = append(result, &models.RequestCacheEntityAttrValue{AttrDefId: vv.AttrDefId, AttrName: k, DataType: vv.AttrDefDataType, DataValue: v})
			}
		}
	}
	return
}

func GetRequestTaskList(requestId string) (result models.TaskQueryResult, err error) {
	var taskTable []*models.TaskTable
	err = dao.X.SQL("select id from task where request=? order by created_time desc", requestId).Find(&taskTable)
	if err != nil {
		return
	}
	if len(taskTable) > 0 {
		result, err = GetTask(taskTable[0].Id)
		return
	}
	// get request
	var requests []*models.RequestTable
	dao.X.SQL("select * from request where id=?", requestId).Find(&requests)
	if len(requests) == 0 {
		return result, fmt.Errorf("Can not find request with id:%s ", requestId)
	}
	var requestCache models.RequestPreDataDto
	err = json.Unmarshal([]byte(requests[0].Cache), &requestCache)
	if err != nil {
		return result, fmt.Errorf("Try to json unmarshal request cache fail,%s ", err.Error())
	}
	var requestTemplateTable []*models.RequestTemplateTable
	dao.X.SQL("select * from request_template where id in (select request_template from request where id=?)", requestId).Find(&requestTemplateTable)
	requestQuery := models.TaskQueryObj{RequestId: requestId, RequestName: requests[0].Name, Reporter: requests[0].Reporter, ReportTime: requests[0].ReportTime, Comment: requests[0].Result, Editable: false}
	requestQuery.FormData = requestCache.Data
	requestQuery.AttachFiles = GetRequestAttachFileList(requestId)
	requestQuery.ExpireTime = requests[0].ExpireTime
	requestQuery.ExpectTime = requests[0].ExpectTime
	requestQuery.ProcInstanceId = requests[0].ProcInstanceId
	if len(requestTemplateTable) > 0 {
		requestQuery.RequestTemplate = requestTemplateTable[0].Name
	}
	result.Data = []*models.TaskQueryObj{&requestQuery}
	result.TimeStep, err = getRequestTimeStep(requests[0].RequestTemplate)
	if len(result.TimeStep) > 0 {
		result.TimeStep[0].Active = true
	}
	return
}

func GetRequestTaskListV2(requestId string) (taskQueryList []*models.TaskQueryObj, err error) {
	var taskTable []*models.TaskTable
	err = dao.X.SQL("select id from task where request=? order by created_time desc", requestId).Find(&taskTable)
	if err != nil {
		return
	}
	if len(taskTable) > 0 {
		taskQueryList, err = GetTaskV2(taskTable[0].Id)
		return
	}
	// get request
	var requests []*models.RequestTable
	dao.X.SQL("select * from request where id=?", requestId).Find(&requests)
	if len(requests) == 0 {
		err = fmt.Errorf("Can not find request with id:%s ", requestId)
		return
	}
	requestQuery := models.TaskQueryObj{RequestId: requestId, RequestName: requests[0].Name, Reporter: requests[0].Reporter, ReportTime: requests[0].ReportTime, Comment: requests[0].Result, Editable: false, RollbackDesc: requests[0].RollbackDesc}
	if requests[0].Cache != "" {
		var requestCache models.RequestPreDataDto
		err = json.Unmarshal([]byte(requests[0].Cache), &requestCache)
		if err != nil {
			return
		}
		requestQuery.FormData = requestCache.Data
	}
	requestQuery.AttachFiles = GetRequestAttachFileList(requestId)
	requestQuery.ExpireTime = requests[0].ExpireTime
	requestQuery.ExpectTime = requests[0].ExpectTime
	requestQuery.ProcInstanceId = requests[0].ProcInstanceId
	requestQuery.CreatedTime = requests[0].CreatedTime
	requestQuery.HandleTime = requests[0].ReportTime
	requestQuery.HandleRoleName = requests[0].Role
	requestQuery.Handler = requests[0].CreatedBy
	taskQueryList = append(taskQueryList, []*models.TaskQueryObj{&requestQuery, getPendingRequestData(requests[0])}...)
	return
}

func getPendingRequestData(request *models.RequestTable) *models.TaskQueryObj {
	var role string
	// 请求在定版状态,从模板角色表中读取
	rtRoleMap := getRequestTemplateMGMTRole()
	roles := rtRoleMap[request.RequestTemplate]
	if len(roles) > 0 {
		role = roles[0]
	}
	taskQueryObj := &models.TaskQueryObj{
		RequestId:      request.Id,
		RequestName:    request.Name,
		Editable:       false,
		Status:         "",
		ExpireTime:     request.ExpireTime,
		ExpectTime:     request.ExpectTime,
		Handler:        request.Handler,
		HandleTime:     request.UpdatedTime,
		FormData:       nil,
		IsHistory:      false,
		HandleRoleName: role,
		CreatedTime:    request.ReportTime,
		RollbackDesc:   request.RollbackDesc,
	}
	if request.Status != "Draft" && request.Status != "Pending" {
		taskQueryObj.HandleTime = request.UpdatedTime
	}
	return taskQueryObj
}

func GetRequestDetailV2(requestId, userToken, language string) (result models.RequestDetail, err error) {
	// get request
	var requests []*models.RequestTable
	var taskQueryList []*models.TaskQueryObj
	var actions []*dao.ExecAction
	var approvalList []*models.TaskTemplateDto
	dao.X.SQL("select * from request where id=?", requestId).Find(&requests)
	if len(requests) == 0 {
		return result, fmt.Errorf("Can not find request with id:%s ", requestId)
	}
	if strings.Contains(requests[0].Status, "InProgress") && requests[0].ProcInstanceId != "" {
		newStatus := getInstanceStatus(requests[0].ProcInstanceId, userToken, language)
		if newStatus == "InternallyTerminated" {
			newStatus = "Termination"
		}
		if newStatus != "" && newStatus != requests[0].Status {
			actions = append(actions, &dao.ExecAction{Sql: "update request set status=?,updated_time=? where id=?",
				Param: []interface{}{newStatus, time.Now().Format(models.DateTimeFormat), requests[0].Id}})
			requests[0].Status = newStatus
		}
		if len(actions) > 0 {
			updateRequestErr := dao.Transaction(actions)
			if updateRequestErr != nil {
				log.Logger.Error("Try to update request status fail", log.Error(updateRequestErr))
			}
		}
	}
	result.Request = getRequestForm(requests[0], userToken, language)
	taskQueryList, err = GetRequestTaskListV2(requestId)
	if err != nil {
		return
	}
	if requests[0].TaskApprovalCache != "" {
		json.Unmarshal([]byte(requests[0].TaskApprovalCache), &approvalList)
		result.ApprovalList = approvalList
	}
	result.Data = taskQueryList
	return
}

func getCMDBSelectList(input []*models.RequestPreDataTableObj, userToken string) (output []*models.RequestPreDataTableObj, err error) {
	ciAttrMap := make(map[string][]string)
	ciAttrSelectMap := make(map[string][]*models.EntityDataObj)
	for _, v := range input {
		if v.Entity == "" {
			continue
		}
		for _, vv := range v.Title {
			if vv.ElementType == "select" && vv.RefEntity == "" {
				if _, b := ciAttrMap[v.Entity]; b {
					ciAttrMap[v.Entity] = append(ciAttrMap[v.Entity], vv.Name)
				} else {
					ciAttrMap[v.Entity] = []string{vv.Name}
				}
			}
			vv.SelectList = []*models.EntityDataObj{}
		}
	}
	output = input
	if len(ciAttrMap) <= 0 {
		return
	}
	for k, v := range ciAttrMap {
		tmpV, tmpErr := getCMDBAttributeSelectList(k, userToken, v)
		if tmpErr != nil {
			err = tmpErr
			break
		}
		for kk, vv := range tmpV {
			ciAttrSelectMap[kk] = vv
		}
	}
	if err != nil {
		return
	}
	for _, v := range output {
		if v.Entity == "" {
			continue
		}
		for _, vv := range v.Title {
			tmpKey := v.Entity + "_" + vv.Name
			if tmpSelectList, b := ciAttrSelectMap[tmpKey]; b {
				vv.SelectList = tmpSelectList
			}
		}
	}
	return
}

func getCMDBAttributes(entity, userToken string) (result []*models.EntityAttributeObj, err error) {
	req, newReqErr := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/wecmdb/api/v1/ci-types-attr/%s/attributes", models.Config.Wecube.BaseUrl, entity), nil)
	if newReqErr != nil {
		err = fmt.Errorf("Try to new http request fail,%s ", newReqErr.Error())
		return
	}
	req.Header.Set("Authorization", userToken)
	//req.Header.Set("Content-Type", "application/json")
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = fmt.Errorf("Try to do http request fail,%s ", respErr.Error())
		return
	}
	responseBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Request cmdb attribute fail,%s ", string(responseBytes))
		return
	}
	var attrQueryResp models.EntityAttributeQueryResponse
	err = json.Unmarshal(responseBytes, &attrQueryResp)
	if err != nil {
		err = fmt.Errorf("Json unmarshal attr response fail,%s ", err.Error())
		return
	}
	result = attrQueryResp.Data
	return
}

func getCMDBAttributeSelectList(entity, userToken string, attributes []string) (result map[string][]*models.EntityDataObj, err error) {
	result = make(map[string][]*models.EntityDataObj)
	attrList, queryErr := getCMDBAttributes(entity, userToken)
	if queryErr != nil {
		err = queryErr
		return
	}
	for _, v := range attrList {
		existFlag := false
		for _, vv := range attributes {
			if v.PropertyName == vv {
				existFlag = true
				break
			}
		}
		if existFlag && v.SelectList != "" {
			tmpSelectList, tmpErr := getAttrCat(v.SelectList, userToken)
			if tmpErr != nil {
				err = tmpErr
				break
			}
			result[entity+"_"+v.PropertyName] = tmpSelectList
		}
	}
	return
}

func getAttrCat(catId, userToken string) (result []*models.EntityDataObj, err error) {
	result = []*models.EntityDataObj{}
	req, newReqErr := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/wecmdb/api/v1/base-key/categories/%s", models.Config.Wecube.BaseUrl, catId), nil)
	if newReqErr != nil {
		err = fmt.Errorf("Try to new http request fail,%s ", newReqErr.Error())
		return
	}
	req.Header.Set("Authorization", userToken)
	//req.Header.Set("Content-Type", "application/json")
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = fmt.Errorf("Try to do http request fail,%s ", respErr.Error())
		return
	}
	responseBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Request cmdb categories fail,%s ", string(responseBytes))
		return
	}
	var response models.CMDBCategoriesResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		err = fmt.Errorf("Json unmarshal categories response fail,%s ", err.Error())
		return
	}
	for _, v := range response.Data {
		result = append(result, &models.EntityDataObj{Id: v.Code, DisplayName: v.Value})
	}
	return
}

func BuildRequestProcessData(input models.RequestCacheData) (result models.RequestProcessData) {
	result.ProcDefId = input.ProcDefId
	result.ProcDefKey = input.ProcDefKey
	result.RootEntityOid = input.RootEntityValue.Oid
	result.Entities = []*models.RequestCacheEntityValue{}
	result.Bindings = []*models.RequestProcessTaskNodeBindObj{}
	entityExistMap := make(map[string]int)
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
	if _, b := entityExistMap[input.RootEntityValue.Oid]; !b {
		result.Entities = append(result.Entities, &input.RootEntityValue)
	}
	if len(result.Entities) == 0 {
		tmpEntityValue := models.RequestCacheEntityValue{Oid: result.RootEntityOid, PackageName: "pseudo", EntityName: "pseudo", BindFlag: "N"}
		result.Entities = append(result.Entities, &tmpEntityValue)
	}
	return result
}

func AppendUselessEntity(requestTemplateId, userToken, language string, cacheData *models.RequestCacheData) (entityDepMap map[string][]string, err error) {
	entityDepMap = make(map[string][]string)
	if cacheData.RootEntityValue.Oid == "" || strings.HasPrefix(cacheData.RootEntityValue.Oid, "tmp") {
		return entityDepMap, nil
	}
	// get core preview data list
	preData, preErr := ProcessDataPreview(requestTemplateId, cacheData.RootEntityValue.Oid, userToken, language)
	if preErr != nil {
		return entityDepMap, fmt.Errorf("Try to get process preview data fail,%s ", preErr.Error())
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
		return entityDepMap, nil
	}
	dependEntityMap := make(map[string]*models.RequestCacheEntityAttrValue)
	log.Logger.Info("getDependEntity", log.StringList("rootSucceeding", rootSucceeding), log.Int("preLen", len(preEntityList)), log.Int("entityLen", len(entityList)))
	// entityList -> in boundValue entity
	getDependEntity(rootSucceeding, rootParent, preEntityList, entityList, dependEntityMap)
	for k, refAttr := range dependEntityMap {
		refDataValue := fmt.Sprintf("%s", refAttr.DataValue)
		if _, b := entityDepMap[k]; b {
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
	return entityDepMap, nil
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

func notifyRoleMail(requestId, userToken, language string) error {
	if !models.MailEnable {
		return nil
	}
	log.Logger.Info("Start notify request mail", log.String("requestId", requestId))
	var roleTable []*models.RoleTable
	err := dao.X.SQL("select id,email from `role` where id in (select `role` from request_template_role where role_type='MGMT' and request_template in (select request_template from request where id=?))", requestId).Find(&roleTable)
	if err != nil {
		return fmt.Errorf("Notify role mail query roles fail,%s ", err.Error())
	}
	if len(roleTable) == 0 {
		return nil
	}
	mailList := GetRoleService().GetRoleMail(roleTable, userToken, language)
	if len(mailList) == 0 {
		log.Logger.Warn("Notify role mail break,email is empty", log.String("role", roleTable[0].Id))
		return nil
	}
	var requestTable []*models.RequestTable
	dao.X.SQL("select t1.id,t1.name,t2.name as request_template,t1.reporter,t1.report_time,t1.emergency from request t1 left join request_template t2 on t1.request_template=t2.id where t1.id=?", requestId).Find(&requestTable)
	if len(requestTable) == 0 {
		return nil
	}
	var subject, content string
	subject = fmt.Sprintf("Taskman Request [%s] %s[%s]", models.PriorityLevelMap[requestTable[0].Emergency], requestTable[0].Name, requestTable[0].RequestTemplate)
	content = fmt.Sprintf("Taskman Request \nID:%s \nPriority:%s \nName:%s \nTemplate:%s \nReporter:%s \nReportTime:%s\n", requestTable[0].Id, models.PriorityLevelMap[requestTable[0].Emergency], requestTable[0].Name, requestTable[0].RequestTemplate, requestTable[0].Reporter, requestTable[0].ReportTime)
	err = models.MailSender.Send(subject, content, mailList)
	if err != nil {
		log.Logger.Error("Notify role mail fail", log.Error(err))
		return fmt.Errorf("Notify role mail fail,%s ", err.Error())
	}
	return nil
}

func CopyRequest(requestId, createdBy string) (result models.RequestTable, err error) {
	parentRequest := &models.RequestTable{}
	var requestTable []*models.RequestTable
	var requestTemplate *models.RequestTemplateTable
	err = dao.X.SQL("select * from request where id=?", requestId).Find(&requestTable)
	if err != nil {
		return
	}
	if len(requestTable) == 0 {
		err = fmt.Errorf("Can not find any request with id:%s ", requestId)
		return
	}
	parentRequest = requestTable[0]
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(parentRequest.RequestTemplate)
	if err != nil {
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
	var formTable []*models.FormTable
	err = dao.X.SQL("select * from form where id=?", parentRequest.Form).Find(&formTable)
	if err != nil {
		return
	}
	if len(formTable) == 0 {
		err = fmt.Errorf("Can not find form %s from parent request:%s ", parentRequest.Form, parentRequest.Id)
		return
	}
	parentForm := formTable[0]
	var formItemTable []*models.FormItemTable
	err = dao.X.SQL("select * from form_item where form=?", parentRequest.Form).Find(&formItemTable)
	if err != nil {
		return
	}
	var actions []*dao.ExecAction
	nowTime := time.Now().Format(models.DateTimeFormat)
	formGuid := guid.CreateGuid()
	newRequestId := newRequestId()
	result.Id = newRequestId
	formInsertAction := dao.ExecAction{Sql: "insert into form(id,name,description,form_template,created_time,created_by,updated_time,updated_by) value (?,?,?,?,?,?,?,?)"}
	formInsertAction.Param = []interface{}{formGuid, parentRequest.Name + models.SysTableIdConnector + "form", "", parentForm.FormTemplate, nowTime, createdBy, nowTime, createdBy}
	actions = append(actions, &formInsertAction)
	requestInsertAction := dao.ExecAction{Sql: "insert into request(id,name,form,request_template,reporter,emergency,report_role,status," +
		"cache,expire_time,expect_time,handler,created_by,created_time,updated_by,updated_time,parent,type,role) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
	requestInsertAction.Param = []interface{}{newRequestId, parentRequest.Name, formGuid, parentRequest.RequestTemplate, createdBy, parentRequest.Emergency, parentRequest.ReportRole, "Draft", parentRequest.Cache,
		"", parentRequest.ExpectTime, parentRequest.Handler, createdBy, nowTime, createdBy, nowTime, parentRequest.Id, parentRequest.Type, parentRequest.Role}
	actions = append(actions, &requestInsertAction)
	for _, formItemRow := range formItemTable {
		actions = append(actions, &dao.ExecAction{Sql: "INSERT INTO form_item (id,form,form_item_template,name,value,item_group,row_data_id) VALUES (?,?,?,?,?,?,?)", Param: []interface{}{
			guid.CreateGuid(), formGuid, formItemRow.FormItemTemplate, formItemRow.Name, formItemRow.Value, formItemRow.ItemGroup, formItemRow.RowDataId,
		}})
	}
	// copy attach file
	var attachFileRows []*models.AttachFileTable
	err = dao.X.SQL("select * from attach_file where request=?", requestId).Find(&attachFileRows)
	if err != nil {
		err = fmt.Errorf("query attach file table fail:%s ", err.Error())
		return
	}
	for _, v := range attachFileRows {
		actions = append(actions, &dao.ExecAction{Sql: "insert into attach_file(id,name,s3_bucket_name,s3_key_name,request,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?,?,?)", Param: []interface{}{
			guid.CreateGuid(), v.Name, v.S3BucketName, v.S3KeyName, newRequestId, createdBy, nowTime, createdBy, nowTime}})
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
	requestId = fmt.Sprintf("%s", time.Now().Format("20060102"))
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
	err = dao.X.Context(c).Table(models.RequestTable{}.TableName()).
		Where("id = ?", requestId).
		Find(&requests)
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
	err = dao.X.SQL("select * from task where request=? order by created_time", requestId).Find(&tasks)
	//err = dao.X.Context(c).Table(models.TaskTable{}.TableName()).
	//	Where("request = ?", requestId).
	//	Asc("created_time").
	//	Find(&tasks)
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
	err = dao.X.Context(c).Table(models.TaskTemplateTable{}.TableName()).
		In("id", taskTmplIds).
		Find(&taskTemplates)
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
	err = dao.X.Context(c).Table(models.TaskHandleTable{}.TableName()).
		In("task", taskIds).
		Asc("created_time").
		Find(&taskHandles)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}

	taskHandleIds := make([]string, 0, len(taskHandles))
	for _, taskHandle := range taskHandles {
		taskHandleIds = append(taskHandleIds, taskHandle.Id)
	}

	// 查询 attach file
	attachFileTaskHandleIdMap := make(map[string][]*models.AttachFileTable)
	if len(taskHandleIds) > 0 {
		attachFiles, tmpErr := GetTaskHandleAttachFileList(c, taskHandleIds)
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
			TaskHandleTable: *taskHandle,
			AttachFiles:     attachFiles,
		}
		taskIdMapHandle[taskHandle.Task] = append(taskIdMapHandle[taskHandle.Task], curTaskHandleForHistory)
	}

	uncompletedTasks := make([]string, 0)
	taskForHistoryList := make([]*models.TaskForHistory, 0, len(tasks))
	log.Logger.Debug(fmt.Sprintf("start handle tasks: [%s]", strings.Join(taskIds, ",")))
	for _, task := range tasks {
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
		}
		if _, isExisted := taskIdMapHandle[task.Id]; isExisted {
			curTaskForHistory.TaskHandleList = taskIdMapHandle[task.Id]
		}

		log.Logger.Debug(fmt.Sprintf("get task: %s form data", task.Id))
		formData, err = getTaskFormData(c, curTaskForHistory)
		if err != nil {
			log.Logger.Error(fmt.Sprintf("get task form data for task: %s error", task.Id), log.Error(err))
			return
		}
		curTaskForHistory.FormData = formData

		taskForHistoryList = append(taskForHistoryList, curTaskForHistory)
	}
	log.Logger.Debug(fmt.Sprintf("finish handle tasks: [%s]", strings.Join(taskIds, ",")))
	result.Task = taskForHistoryList
	result.Request.UncompletedTasks = uncompletedTasks
	return
}

func getTaskFormData(c *gin.Context, taskObj *models.TaskForHistory) (result []*models.RequestPreDataTableObj, err error) {
	result = []*models.RequestPreDataTableObj{}

	// 查询 form template
	var formTemplates []*models.FormTemplateTable
	err = dao.X.Context(c).Table(models.FormTemplateTable{}.TableName()).
		Where("task_template = ?", taskObj.TaskTemplate).
		Find(&formTemplates)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}
	if len(formTemplates) == 0 {
		log.Logger.Error(fmt.Sprintf("can not find any form templates with taskTemplate: %s", taskObj.TaskTemplate))
		return
	}

	formTemplateIds := make([]string, 0, len(formTemplates))
	formTemplateRefIds := make([]string, 0, len(formTemplates))
	for _, formTmpl := range formTemplates {
		formTemplateIds = append(formTemplateIds, formTmpl.Id)
		if formTmpl.RefId != "" {
			formTemplateRefIds = append(formTemplateRefIds, formTmpl.RefId)
		}
	}

	var actualFormTemplates []*models.FormTemplateTable
	actualFormTemplateIds := formTemplateIds
	if taskObj.ProcDefId == "" {
		// 非编排任务, 表单form的 form_template 用的数据表单的 form_template,
		// 需要在 form_template 里面通过id -> ref_id，这个 ref_id 才是 form 的 form_template
		actualFormTemplateIds = formTemplateRefIds

		// 查询 form 实际使用的 formTemplate
		err = dao.X.Context(c).Table(models.FormTemplateTable{}.TableName()).
			In("id", actualFormTemplateIds).
			Find(&actualFormTemplates)
		if err != nil {
			err = exterror.Catch(exterror.New().DatabaseQueryError, err)
			return
		}
		if len(actualFormTemplates) == 0 {
			log.Logger.Error(fmt.Sprintf("can not find any form templates with actualFormTemplateIds: [%s]", strings.Join(actualFormTemplateIds, ",")))
			return
		}
	} else {
		actualFormTemplates = formTemplates
	}
	/*
		actualFormTemplateIdMapInfo := make(map[string]*models.FormTemplateTable)
		for _, formTemplate := range actualFormTemplates {
			actualFormTemplateIdMapInfo[formTemplate.Id] = formTemplate
		}
	*/

	// 查询 form
	var taskForms []*models.FormTable
	err = dao.X.Context(c).Table(models.FormTable{}.TableName()).
		Where("request = ?", taskObj.Request).
		In("form_template", actualFormTemplateIds).
		Find(&taskForms)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}
	if len(taskForms) == 0 {
		log.Logger.Error(fmt.Sprintf("can not find any forms with request: %s and formTemplates: [%s]",
			taskObj.Request, strings.Join(actualFormTemplateIds, ",")))
		return
	}
	/*
		taskFormIdMapInfo := make(map[string]*models.FormTable)
		for _, form := range taskForms {
			taskFormIdMapInfo[form.Id] = form
		}
	*/

	// 查询 request 的 form item
	taskUpdatedTime := taskObj.UpdatedTime
	taskHandleCnt := len(taskObj.TaskHandleList)
	if taskHandleCnt > 0 {
		taskUpdatedTime = taskObj.TaskHandleList[taskHandleCnt-1].UpdatedTime
	}
	var requestFormItems []*models.FormItemTable
	err = dao.X.Context(c).Table(models.FormItemTable{}.TableName()).
		Where("request = ? AND updated_time <= ?", taskObj.Request, taskUpdatedTime).
		Desc("updated_time").
		Find(&requestFormItems)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}
	if len(requestFormItems) == 0 {
		log.Logger.Error(fmt.Sprintf("can not find any form items with request: %s and updatedTime <= %s",
			taskObj.Request, taskUpdatedTime))
		return
	}

	// 查询 form item template
	var itemTemplates []*models.FormItemTemplateTable
	// err = dao.X.Context(c).SQL("select * from form_item_template where form_template in (select form_template from task_template where id=?) order by item_group,sort", taskObj.TaskTemplate).Find(&itemTemplates)
	err = dao.X.Context(c).Table(models.FormItemTemplateTable{}.TableName()).
		In("form_template", actualFormTemplateIds).
		OrderBy("item_group,sort").
		Find(&itemTemplates)
	if err != nil {
		err = exterror.Catch(exterror.New().DatabaseQueryError, err)
		return
	}
	if len(itemTemplates) == 0 {
		// err = fmt.Errorf("can not find any form item template with task: %s", taskObj.Id)
		log.Logger.Error(fmt.Sprintf("can not find any form item templates with formTemplates: [%s]",
			strings.Join(actualFormTemplateIds, ",")))
		return
	}
	formResult := getItemTemplateTitle(itemTemplates)
	result = formResult

	// 通过筛选 requestFormItems 获取当前 task 的 form items
	/*
		dao.X.Context(c).SQL("select * from form_item where form=? order by item_group,row_data_id", taskObj.Form).Find(&items)
	*/
	taskFormItems := getTaskFormItems(requestFormItems, taskForms)
	if len(taskFormItems) == 0 {
		log.Logger.Error(fmt.Sprintf("can not find any form item for task: %s", taskObj.Id))
		return
	}

	formIdMapInfo := make(map[string]*models.FormTable)
	formTemplateIdMapInfo := make(map[string]*models.FormTemplateTable)
	formIdMapInfo, formTemplateIdMapInfo, err = getFormAndTemplateMapInfo(c, taskFormItems)
	if err != nil {
		err = fmt.Errorf("get form and form template map info failed: %s", err.Error())
		return
	}

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
				log.Logger.Error(fmt.Sprintf("can not find itemGroup for formItem: %s", item.Id))
				continue
			}
		} else {
			log.Logger.Error(fmt.Sprintf("can not find itemDataId for formItem: %s", item.Id))
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
		if _, b := rowItemMap[itemDataId]; b {
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
					for _, tmpTitle := range formTable.Title {
						if tmpTitle.Name == rowItem.Name {
							if tmpTitle.Multiple == "Y" {
								isMulti = true
								break
							}
						}
					}
					if isMulti {
						tmpRowObj.EntityData[rowItem.Name] = strings.Split(rowItem.Value, ",")
					} else {
						tmpRowObj.EntityData[rowItem.Name] = rowItem.Value
					}
				}
				formTable.Value = append(formTable.Value, &tmpRowObj)
			}
		}
	}
	return
}

func getFormAndTemplateMapInfo(c *gin.Context, taskFormItems []*models.FormItemTable) (formIdMapInfo map[string]*models.FormTable, formTemplateIdMapInfo map[string]*models.FormTemplateTable, err error) {
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
		err = dao.X.Context(c).Table(models.FormTable{}.TableName()).
			In("id", formIds).
			Find(&forms)
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
				err = dao.X.Context(c).Table(models.FormTemplateTable{}.TableName()).
					In("id", formTemplateIds).
					Find(&formTemplates)
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
