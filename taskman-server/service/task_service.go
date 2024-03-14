package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/exterror"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/rpc"
	"xorm.io/xorm"
)

type TaskService struct {
	taskDao       *dao.TaskDao
	taskHandleDao *dao.TaskHandleDao
}

func GetTaskFormStruct(procInstId, nodeDefId string) (result models.TaskMetaResult, err error) {
	result = models.TaskMetaResult{Status: "OK", Message: "Success"}
	var items []*models.FormItemTemplateTable
	//err = dao.X.SQL("select * from form_item_template where form_template in (select form_template from task_template where node_def_id=? and request_template in (select request_template from request where proc_instance_id=?))", nodeDefId, procInstId).Find(&items)
	err = dao.X.SQL("select * from form_item_template where form_template in (select id from form_template where task_template in (select id from task_template where node_def_id=?) and request_template in (select request_template from request where proc_instance_id=?))", nodeDefId, procInstId).Find(&items)
	if err != nil {
		return
	}
	if len(items) == 0 {
		err = fmt.Errorf("Can not find task item template with procInstId:%s nodeDefId:%s ", procInstId, nodeDefId)
		return
	}
	queryRows, queryErr := dao.X.QueryString("select id,task_template from form_template where id=?", items[0].FormTemplate)
	if queryErr != nil {
		err = fmt.Errorf("query task tempalte with form template id:%s fail,%s ", items[0].FormTemplate, queryErr.Error())
		return
	}
	resultData := models.TaskMetaResultData{FormMetaId: items[0].FormTemplate}
	if len(queryRows) > 0 {
		resultData.FormMetaId = queryRows[0]["task_template"]
	}
	for _, item := range items {
		if item.Entity == "" {
			continue
		}
		resultData.FormItemMetas = append(resultData.FormItemMetas, &models.TaskMetaResultItem{FormItemMetaId: item.Id, PackageName: item.PackageName, EntityName: item.Entity, AttrName: item.Name})
	}
	result.Data = resultData
	return
}

func PluginTaskCreateNew(input *models.PluginTaskCreateRequestObj, callRequestId, dueDate string, nextOptions []string, userToken, language string) (result *models.PluginTaskCreateOutputObj, taskId string, err error) {
	log.Logger.Debug("task create", log.JsonObj("input", input))
	result = &models.PluginTaskCreateOutputObj{CallbackParameter: input.CallbackParameter, ErrorCode: "0", ErrorMessage: "", Comment: ""}
	var requestTable []*models.RequestTable
	err = dao.X.SQL("select * from request where proc_instance_id=?", input.ProcInstId).Find(&requestTable)
	if err != nil {
		return result, taskId, fmt.Errorf("Try to check proc_instance_id:%s is in request fail,%s ", input.ProcInstId, err.Error())
	}
	var actions []*dao.ExecAction
	var taskSort int
	nowTime := time.Now().Format(models.DateTimeFormat)
	input.RoleName = remakeTaskReportRole(input.RoleName)
	newTaskObj := models.TaskTable{Id: "tk_" + guid.CreateGuid(), Name: input.TaskName, Status: "created", Reporter: input.Reporter, ReportRole: input.RoleName, Description: input.TaskDescription, CallbackUrl: input.CallbackUrl, CallbackParameter: input.CallbackParameter, NextOption: strings.Join(nextOptions, ","), Handler: input.Handler}
	taskId = newTaskObj.Id
	operator := "system"
	var taskFormInput models.PluginTaskFormDto
	if input.TaskFormInput != "" {
		err = json.Unmarshal([]byte(input.TaskFormInput), &taskFormInput)
		if err != nil {
			return result, taskId, fmt.Errorf("Try to json unmarshal taskFormInput to json data fail,%s ", err.Error())
		}
		if newTaskObj.Reporter == "" {
			newTaskObj.Reporter = "taskman"
		}
	} else {
		// 自定义任务的发起，不需要表单，只创建任务
		customExpireTime := ""
		dueMin, _ := strconv.Atoi(dueDate)
		if dueMin > 0 {
			customExpireTime = time.Now().Add(time.Duration(dueMin) * time.Minute).Format(models.DateTimeFormat)
		}

		taskInsertAction := dao.ExecAction{Sql: "insert into task(id,name,description,status,proc_def_id,proc_def_key,node_def_id,node_name,callback_url," +
			"callback_parameter,reporter,report_role,report_time,expire_time,emergency,callback_request_id,next_option,handler,created_by,created_time," +
			"updated_by,updated_time,type) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
		taskInsertAction.Param = []interface{}{newTaskObj.Id, newTaskObj.Name, newTaskObj.Description, newTaskObj.Status, newTaskObj.ProcDefId,
			newTaskObj.ProcDefKey, newTaskObj.NodeDefId, newTaskObj.NodeName, newTaskObj.CallbackUrl, newTaskObj.CallbackParameter, newTaskObj.Reporter,
			newTaskObj.ReportRole, nowTime, customExpireTime, newTaskObj.Emergency, callRequestId, newTaskObj.NextOption, newTaskObj.Handler, operator,
			nowTime, operator, nowTime, models.TaskTypeImplement}
		actions = append(actions, &taskInsertAction)

		err = dao.Transaction(actions)
		return
	}
	if len(requestTable) == 0 {
		err = fmt.Errorf("can not find request with proc_instance_id:%s", input.ProcInstId)
		return
	}
	taskSort = GetTaskService().GenerateTaskOrderByRequestId(requestTable[0].Id)
	newTaskObj.ProcDefId = taskFormInput.ProcDefId
	newTaskObj.ProcDefKey = taskFormInput.ProcDefKey
	newTaskObj.Request = requestTable[0].Id
	newTaskObj.NodeDefId = taskFormInput.TaskNodeDefId
	newTaskObj.Emergency = requestTable[0].Emergency
	newTaskObj.TemplateType = requestTable[0].Type
	var taskTemplateTable []*models.TaskTemplateTable
	dao.X.SQL("select * from task_template where request_template=? and node_def_id=?", requestTable[0].RequestTemplate, taskFormInput.TaskNodeDefId).Find(&taskTemplateTable)
	if len(taskTemplateTable) > 0 {
		newTaskObj.TaskTemplate = taskTemplateTable[0].Id
		newTaskObj.NodeName = taskTemplateTable[0].NodeName
		newTaskObj.ExpireTime = calcExpireTime(nowTime, taskTemplateTable[0].ExpireDay)
		newTaskObj.Name = taskTemplateTable[0].Name
		newTaskObj.Description = taskTemplateTable[0].Description
		newTaskObj.Reporter = requestTable[0].Reporter
		newTaskObj.ReportTime = nowTime
		newTaskObj.Handler = taskTemplateTable[0].Handler
	} else {
		log.Logger.Warn("Can not find any taskTemplate", log.String("requestTemplate", requestTable[0].RequestTemplate), log.String("nodeDefId", taskFormInput.TaskNodeDefId))
		err = fmt.Errorf("Can not find any taskTemplate in request:%s with nodeDefId:%s ", newTaskObj.Request, taskFormInput.TaskNodeDefId)
		return
	}
	taskInsertAction := dao.ExecAction{Sql: "insert into task(id,name,description,form,status,request,task_template,proc_def_id,proc_def_key,node_def_id," +
		"node_name,callback_url,callback_parameter,reporter,report_role,report_time,emergency,cache,callback_request_id,next_option,expire_time," +
		"handler,created_by,created_time,updated_by,updated_time,template_type,type,sort) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
	taskInsertAction.Param = []interface{}{newTaskObj.Id, newTaskObj.Name, newTaskObj.Description, newTaskObj.Form, newTaskObj.Status,
		newTaskObj.Request, newTaskObj.TaskTemplate, newTaskObj.ProcDefId, newTaskObj.ProcDefKey, newTaskObj.NodeDefId, newTaskObj.NodeName,
		newTaskObj.CallbackUrl, newTaskObj.CallbackParameter, newTaskObj.Reporter, newTaskObj.ReportRole, nowTime, newTaskObj.Emergency,
		input.TaskFormInput, callRequestId, newTaskObj.NextOption, newTaskObj.ExpireTime, newTaskObj.Handler, operator, nowTime, operator,
		nowTime, newTaskObj.TemplateType, models.TaskTypeImplement, taskSort}
	actions = append(actions, &taskInsertAction)
	// 新增form
	var formTemplateRows []*models.FormTemplateTable
	err = dao.X.SQL("select * from form_template where task_template=? and item_group_type='workflow'", newTaskObj.TaskTemplate).Find(&formTemplateRows)
	if err != nil {
		err = fmt.Errorf("query form template table fail,%s ", err.Error())
		return
	}
	for _, formDataEntity := range taskFormInput.FormDataEntities {
		entityItemGroup := fmt.Sprintf("%s:%s", formDataEntity.PackageName, formDataEntity.EntityName)
		tmpFormTemplateId := ""
		for _, formTemplateRow := range formTemplateRows {
			if formTemplateRow.ItemGroup == entityItemGroup {
				tmpFormTemplateId = formTemplateRow.Id
				break
			}
		}
		if tmpFormTemplateId == "" {
			log.Logger.Warn("form data entity can not find form template", log.String("task", taskId), log.JsonObj("formDataEntity", formDataEntity))
			continue
		}
		newFormId := "form_" + guid.CreateGuid()
		actions = append(actions, &dao.ExecAction{Sql: "insert into form(id,request,task,form_template,data_id,created_by,updated_by,created_time,updated_time) values (?,?,?,?,?,?,?,?,?)", Param: []interface{}{
			newFormId, newTaskObj.Request, taskId, tmpFormTemplateId, formDataEntity.Oid, operator, operator, nowTime, nowTime,
		}})
		for _, formDataItem := range formDataEntity.FormItemValues {
			actions = append(actions, &dao.ExecAction{Sql: "insert into form_item(id,form,form_item_template,name,value,request,updated_time) values (?,?,?,?,?,?,?)", Param: []interface{}{
				"item_" + guid.CreateGuid(), newFormId, formDataItem.FormItemMetaId, formDataItem.AttrName, formDataItem.AttrValue, newTaskObj.Request, nowTime,
			}})
		}
	}
	createTaskHandleAction := GetTaskHandleService().CreateTaskHandleByTemplate(newTaskObj.Id, userToken, language, requestTable[0], taskTemplateTable[0])
	actions = append(actions, createTaskHandleAction...)
	err = dao.TransactionWithoutForeignCheck(actions)
	return
}

func getLastCustomFormItem(requestId, taskFormTemplateId, newTaskFormId string) (result []*models.FormItemTable, err error) {
	result = []*models.FormItemTable{}
	if requestId == "" || taskFormTemplateId == "" {
		return
	}
	var formItemTemplates []*models.FormItemTemplateTable
	err = dao.X.SQL("select * from form_item_template where entity='' and form_template=?", taskFormTemplateId).Find(&formItemTemplates)
	if len(formItemTemplates) == 0 || err != nil {
		return
	}
	groupNameTemplateIdMap := make(map[string]string)
	filterList := []string{}
	for _, v := range formItemTemplates {
		groupNameTemplateIdMap[fmt.Sprintf("%s_%s", v.ItemGroup, v.Name)] = v.Id
		filterList = append(filterList, fmt.Sprintf("(item_group='%s' and name='%s')", v.ItemGroup, v.Name))
	}
	var formItems []*models.FormItemTable
	err = dao.X.SQL("select * from form_item where ("+strings.Join(filterList, " or ")+") and form in (select form from request where id=? union select form from task where request=?) order by item_group,name,id desc", requestId, requestId).Find(&formItems)
	//if len(formItems) == 0 {
	//	for _, v := range formItemTemplates {
	//		tmpKey := fmt.Sprintf("%s_%s", v.ItemGroup, v.Name)
	//		result = append(result, &models.FormItemTable{Form: newTaskFormId, FormItemTemplate: groupNameTemplateIdMap[tmpKey], Name: v.Name, ItemGroup: v.ItemGroup, Value: "", RowDataId: ""})
	//	}
	//	return
	//}
	groupNameExistMap := make(map[string]int)
	for _, v := range formItems {
		tmpKey := fmt.Sprintf("%s_%s", v.ItemGroup, v.Name)
		if _, b := groupNameExistMap[tmpKey]; b {
			continue
		}
		result = append(result, &models.FormItemTable{Form: newTaskFormId, FormItemTemplate: groupNameTemplateIdMap[tmpKey], Name: v.Name, ItemGroup: v.ItemGroup, Value: v.Value, RowDataId: v.RowDataId})
		groupNameExistMap[tmpKey] = 1
	}
	return
}

func getFormItemTemplateMap(formTemplateId string) map[string]*models.FormItemTemplateTable {
	resultMap := make(map[string]*models.FormItemTemplateTable)
	var itemTemplateTable []*models.FormItemTemplateTable
	dao.X.SQL("select * from form_item_template where form_template=?", formTemplateId).Find(&itemTemplateTable)
	for _, v := range itemTemplateTable {
		resultMap[v.Id] = v
	}
	return resultMap
}

func remakeTaskReportRole(reportRoles string) string {
	resultList := []string{}
	for _, v := range strings.Split(reportRoles, ",") {
		if v != "" {
			resultList = append(resultList, v)
		}
	}
	if len(resultList) > 0 {
		return fmt.Sprintf(",%s,", strings.Join(resultList, ","))
	}
	return ""
}

func ListTask(param *models.QueryRequestParam, userRoles []string, operator string) (pageInfo models.PageInfo, rowData []*models.TaskListObj, err error) {
	rowData = []*models.TaskListObj{}
	roleFilterSql := "1=1"
	if len(userRoles) > 0 {
		roleFilterList := []string{}
		for _, v := range userRoles {
			roleFilterList = append(roleFilterList, "report_role like '%,"+v+",%'")
		}
		roleFilterSql = strings.Join(roleFilterList, " or ")
	}
	filterSql, _, queryParam := dao.TransFiltersToSQL(param, &models.TransFiltersParam{IsStruct: true, StructObj: models.TaskTable{}, PrimaryKey: "id", Prefix: "t1"})
	baseSql := fmt.Sprintf("select t1.id,t1.name,t1.description,t1.form,t1.status,t1.`version`,t1.request,t1.task_template,t1.node_name,t1.reporter,t1.report_role,t1.report_time,t1.emergency,t1.handler,t1.created_by,t1.created_time,t1.updated_by,t1.updated_time,t1.expire_time from (select * from task where task_template in (select task_template from task_template_role where role_type='USE' and `role` in ('"+strings.Join(userRoles, "','")+"')) union select * from task where task_template is null and (%s)) t1 where t1.del_flag=0 %s ", roleFilterSql, filterSql)
	if param.Paging {
		pageInfo.StartIndex = param.Pageable.StartIndex
		pageInfo.PageSize = param.Pageable.PageSize
		pageInfo.TotalRows = dao.QueryCount(baseSql, queryParam...)
		pageSql, pageParam := dao.TransPageInfoToSQL(*param.Pageable)
		baseSql += pageSql
		queryParam = append(queryParam, pageParam...)
	}
	err = dao.X.SQL(baseSql, queryParam...).Find(&rowData)
	if err != nil {
		return models.PageInfo{}, nil, err
	}
	var requestIdList []string
	roleMap := getTaskTemplateRoles()
	for _, v := range rowData {
		buildTaskOperation(v, operator)
		requestIdList = append(requestIdList, v.Request)
		if tmpRoles, b := roleMap[v.TaskTemplate]; b {
			v.HandleRoles = tmpRoles
		} else {
			if v.ReportRole != "" {
				v.HandleRoles = strings.Split(v.ReportRole[1:len(v.ReportRole)-1], ",")
			} else {
				v.HandleRoles = []string{}
			}
		}
	}
	var requestTables []*models.RequestTable
	dao.X.SQL("select t1.id,t1.name,t1.reporter,t1.report_time,t2.name as request_template from request t1 left join request_template t2 on t1.request_template=t2.id where t1.id in ('" + strings.Join(requestIdList, "','") + "')").Find(&requestTables)
	requestMap := make(map[string]*models.RequestTable)
	for _, v := range requestTables {
		requestMap[v.Id] = v
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	for _, v := range rowData {
		if _, b := requestMap[v.Request]; b {
			v.RequestObj = *requestMap[v.Request]
			v.Reporter = "taskman"
		}
		if v.ExpireTime != "" {
			timeObj := models.ExpireObj{ReportTime: v.ReportTime, ExpireTime: v.ExpireTime, NowTime: nowTime}
			if v.Status == "done" {
				timeObj = models.ExpireObj{ReportTime: v.ReportTime, ExpireTime: v.ExpireTime, NowTime: v.UpdatedTime}
			}
			calcExpireObj(&timeObj)
			v.ExpirePercentObj = timeObj
		}
	}
	return
}

func getTaskTemplateRoles() map[string][]string {
	result := make(map[string][]string)
	var taskRoles []*models.TaskTemplateRoleTable
	dao.X.SQL("select * from task_template_role where role_type='USE' order by task_template").Find(&taskRoles)
	if len(taskRoles) == 0 {
		return result
	}
	tmpTemplate := taskRoles[0].TaskTemplate
	tmpList := []string{}
	for _, v := range taskRoles {
		if v.TaskTemplate != tmpTemplate {
			result[tmpTemplate] = tmpList
			tmpTemplate = v.TaskTemplate
			tmpList = []string{}
		}
		tmpList = append(tmpList, v.Role)
	}
	if len(tmpList) > 0 {
		tmpTemplate = taskRoles[len(taskRoles)-1].TaskTemplate
		result[tmpTemplate] = tmpList
	}
	return result
}

func GetTask(taskId string) (result models.TaskQueryResult, err error) {
	result = models.TaskQueryResult{}
	taskObj, tmpErr := getSimpleTask(taskId)
	if tmpErr != nil {
		return result, tmpErr
	}
	if taskObj.Request == "" {
		taskForm, tmpErr := queryTaskForm(&taskObj)
		if tmpErr != nil {
			return result, tmpErr
		}
		result.Data = []*models.TaskQueryObj{&taskForm}
		result.TimeStep = []*models.TaskQueryTimeStep{}
		return
	}
	// get request
	var requests []*models.RequestTable
	dao.X.SQL("select * from request where id=?", taskObj.Request).Find(&requests)
	if len(requests) == 0 {
		return result, fmt.Errorf("Can not find request with id:%s ", taskObj.Request)
	}
	if requests[0].Parent != "" {
		if parentRequest, getParentErr := GetRequestTaskList(requests[0].Parent); getParentErr != nil {
			err = getParentErr
			return
		} else {
			for _, v := range parentRequest.Data {
				v.IsHistory = true
			}
			result.Data = parentRequest.Data
		}
	}
	var requestCache models.RequestPreDataDto
	err = json.Unmarshal([]byte(requests[0].Cache), &requestCache)
	if err != nil {
		return result, fmt.Errorf("Try to json unmarshal request cache fail,%s ", err.Error())
	}
	var requestTemplateTable []*models.RequestTemplateTable
	dao.X.SQL("select * from request_template where id in (select request_template from request where id=?)", taskObj.Request).Find(&requestTemplateTable)
	requestQuery := models.TaskQueryObj{RequestId: taskObj.Request, RequestName: requests[0].Name, Reporter: requests[0].Reporter, ReportTime: requests[0].ReportTime, Comment: requests[0].Result, Editable: false}
	requestQuery.AttachFiles = GetRequestAttachFileList(taskObj.Request)
	requestQuery.ExpireTime = requests[0].ExpireTime
	requestQuery.ExpectTime = requests[0].ExpectTime
	requestQuery.ProcInstanceId = requests[0].ProcInstanceId
	requestQuery.FormData = requestCache.Data
	if len(requestTemplateTable) > 0 {
		requestQuery.RequestTemplate = requestTemplateTable[0].Name
	}
	result.Data = append(result.Data, &requestQuery)
	result.TimeStep, err = getRequestTimeStep(requests[0].RequestTemplate)
	if err != nil {
		return
	}
	// get task list
	var taskHandlerRows []*models.TaskHandlerQueryData
	dao.X.SQL("select t1.id,t2.handler,t4.display_name from task t1 left join task_template t2 on t1.task_template=t2.id left join task_template_role t3 on t1.task_template=t3.task_template left join `role` t4 on t3.`role`=t4.id where t1.request=?", taskObj.Request).Find(&taskHandlerRows)
	taskHandlerMap := make(map[string]*models.TaskHandlerQueryData)
	for _, row := range taskHandlerRows {
		taskHandlerMap[row.Id] = row
	}
	var taskList []*models.TaskTable
	dao.X.SQL("select * from task where request=? and report_time<='"+taskObj.ReportTime+"' order by created_time", taskObj.Request).Find(&taskList)
	for _, v := range taskList {
		tmpTaskForm, tmpErr := queryTaskForm(v)
		if tmpErr != nil {
			err = tmpErr
			break
		}
		if handlerObj, b := taskHandlerMap[tmpTaskForm.TaskId]; b {
			tmpTaskForm.Handler = handlerObj.Handler
			tmpTaskForm.HandleRoleName = handlerObj.DisplayName
		}
		if v.Status == "done" {
			tmpTaskForm.Handler = v.UpdatedBy
		}
		result.Data = append(result.Data, &tmpTaskForm)
	}
	if err != nil {
		return
	}
	for _, v := range result.TimeStep {
		if v.TaskTemplateId == taskObj.TaskTemplate {
			v.Active = true
			break
		}
	}
	return
}

func GetTaskV2(taskId string) (taskQueryList []*models.TaskQueryObj, err error) {
	var taskObj models.TaskTable
	var taskForm models.TaskQueryObj
	taskObj, err = getSimpleTask(taskId)
	if err != nil {
		return
	}
	if taskObj.Request == "" {
		taskForm, err = queryTaskForm(&taskObj)
		if err != nil {
			return
		}
		taskQueryList = []*models.TaskQueryObj{&taskForm}
		return
	}
	// get request
	var requests []*models.RequestTable
	dao.X.SQL("select * from request where id=?", taskObj.Request).Find(&requests)
	if len(requests) == 0 {
		err = fmt.Errorf("Can not find request with id:%s ", taskObj.Request)
		return
	}
	/*// 当前请求可能是历史请求重新发起生成,则需要展示历史请求ID的请求和任务数据
	if requests[0].Parent != "" {
		if parentRequestTaskQueryList, getParentErr := GetRequestTaskListV2(requests[0].Parent); getParentErr != nil {
			err = getParentErr
			return
		} else {
			if len(parentRequestTaskQueryList) > 0 {
				for _, taskQuery := range parentRequestTaskQueryList {
					taskQuery.IsHistory = true
				}
				taskQueryList = parentRequestTaskQueryList
			}
		}
	}*/
	var requestCache models.RequestPreDataDto
	err = json.Unmarshal([]byte(requests[0].Cache), &requestCache)
	if err != nil {
		return
	}
	var requestTemplateTable []*models.RequestTemplateTable
	dao.X.SQL("select * from request_template where id in (select request_template from request where id=?)", taskObj.Request).Find(&requestTemplateTable)
	requestQuery := models.TaskQueryObj{RequestId: taskObj.Request, RequestName: requests[0].Name, Reporter: requests[0].Reporter, ReportTime: requests[0].ReportTime, Comment: requests[0].Result, Editable: false, RollbackDesc: requests[0].RollbackDesc}
	requestQuery.AttachFiles = GetRequestAttachFileList(taskObj.Request)
	requestQuery.ExpireTime = requests[0].ExpireTime
	requestQuery.ExpectTime = requests[0].ExpectTime
	requestQuery.ProcInstanceId = requests[0].ProcInstanceId
	requestQuery.FormData = requestCache.Data
	requestQuery.CreatedTime = requests[0].CreatedTime
	requestQuery.HandleTime = requests[0].ReportTime
	requestQuery.Handler = requests[0].CreatedBy
	requestQuery.HandleRoleName = requests[0].Role
	if len(requestTemplateTable) > 0 {
		requestQuery.RequestTemplate = requestTemplateTable[0].Name
	}
	taskQueryList = append(taskQueryList, []*models.TaskQueryObj{&requestQuery, getPendingRequestData(requests[0])}...)
	// get task list
	var taskHandlerRows []*models.TaskHandlerQueryData
	dao.X.SQL("select t1.id,t2.handler,t4.display_name from task t1 left join task_template t2 on t1.task_template=t2.id left join task_template_role t3 on t1.task_template=t3.task_template left join `role` t4 on t3.`role`=t4.id where t1.request=?", taskObj.Request).Find(&taskHandlerRows)
	taskHandlerMap := make(map[string]*models.TaskHandlerQueryData)
	for _, row := range taskHandlerRows {
		taskHandlerMap[row.Id] = row
	}
	var taskList []*models.TaskTable
	dao.X.SQL("select * from task where request=? and report_time<='"+taskObj.ReportTime+"' order by created_time", taskObj.Request).Find(&taskList)
	for _, v := range taskList {
		tmpTaskForm, tmpErr := queryTaskForm(v)
		if tmpErr != nil {
			err = tmpErr
			break
		}
		if handlerObj, b := taskHandlerMap[tmpTaskForm.TaskId]; b {
			tmpTaskForm.Handler = handlerObj.Handler
			tmpTaskForm.HandleRoleName = handlerObj.DisplayName
		}
		if v.Status == "done" {
			tmpTaskForm.Handler = v.UpdatedBy
			tmpTaskForm.HandleTime = v.UpdatedTime
		}
		if len(strings.TrimSpace(v.NodeName)) > 0 {
			tmpTaskForm.TaskName = v.NodeName
		}
		tmpTaskForm.CreatedTime = v.CreatedTime
		taskQueryList = append(taskQueryList, &tmpTaskForm)
	}
	return
}

func queryTaskForm(taskObj *models.TaskTable) (taskForm models.TaskQueryObj, err error) {
	taskForm = models.TaskQueryObj{TaskId: taskObj.Id, TaskName: taskObj.Name, Description: taskObj.Description, RequestId: taskObj.Request, Reporter: taskObj.Reporter, ReportTime: taskObj.ReportTime, Comment: taskObj.Result, Status: taskObj.Status, NextOption: []string{}, ExpireTime: taskObj.ExpireTime, FormData: []*models.RequestPreDataTableObj{}}
	taskForm.AttachFiles = GetTaskAttachFileList(taskObj.Id)
	if taskObj.Status != "done" {
		taskForm.Editable = true
	} else {
		taskForm.Handler = taskObj.UpdatedBy
		taskForm.HandleTime = taskObj.UpdatedTime
	}
	if taskObj.NextOption != "" {
		taskForm.NextOption = strings.Split(taskObj.NextOption, ",")
		taskForm.ChoseOption = taskObj.ChoseOption
	}
	if taskObj.Request == "" {
		return
	}
	var itemTemplates []*models.FormItemTemplateTable
	//err = dao.X.SQL("select * from form_item_template where form_template in (select form_template from task_template where id=?) order by item_group,sort", taskObj.TaskTemplate).Find(&itemTemplates)
	err = dao.X.SQL("select * from form_item_template where form_template in (select id from form_template where task_template=?) order by item_group,sort", taskObj.TaskTemplate).Find(&itemTemplates)
	if err != nil {
		return
	}
	if len(itemTemplates) == 0 {
		return taskForm, fmt.Errorf("Can not find any form item template with task:%s ", taskObj.Id)
	}
	formResult := getItemTemplateTitle(itemTemplates)
	taskForm.FormData = formResult
	var items []*models.FormItemTable
	dao.X.SQL("select * from form_item where form=? order by item_group,row_data_id", taskObj.Form).Find(&items)
	if len(items) == 0 {
		return
	}
	itemRowMap := make(map[string][]string)
	rowItemMap := make(map[string][]*models.FormItemTable)
	for _, item := range items {
		if tmpRows, b := itemRowMap[item.ItemGroup]; b {
			existFlag := false
			for _, v := range tmpRows {
				if item.RowDataId == v {
					existFlag = true
					break
				}
			}
			if !existFlag {
				itemRowMap[item.ItemGroup] = append(itemRowMap[item.ItemGroup], item.RowDataId)
			}
		} else {
			itemRowMap[item.ItemGroup] = []string{item.RowDataId}
		}
		if _, b := rowItemMap[item.RowDataId]; b {
			rowItemMap[item.RowDataId] = append(rowItemMap[item.RowDataId], item)
		} else {
			rowItemMap[item.RowDataId] = []*models.FormItemTable{item}
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
	taskForm.FormData = formResult
	return
}

func getRequestTimeStep(requestTemplateId string) (result []*models.TaskQueryTimeStep, err error) {
	var requestTemplateTable []*models.RequestTemplateTable
	err = dao.X.SQL("select id,name from request_template where id=?", requestTemplateId).Find(&requestTemplateTable)
	if err != nil {
		return
	}
	if len(requestTemplateTable) == 0 {
		return result, fmt.Errorf("Can not find requestTemplate with id:%s ", requestTemplateId)
	}
	result = append(result, &models.TaskQueryTimeStep{RequestTemplateId: requestTemplateTable[0].Id, Name: "Start", Active: false})
	var taskTemplateTable []*models.TaskTemplateTable
	dao.X.SQL("select id,name from task_template where request_template=?", requestTemplateId).Find(&taskTemplateTable)
	for _, v := range taskTemplateTable {
		result = append(result, &models.TaskQueryTimeStep{RequestTemplateId: requestTemplateId, TaskTemplateId: v.Id, Name: v.Name, Active: false})
	}
	return
}

func getSimpleTask(taskId string) (result models.TaskTable, err error) {
	var taskTable []*models.TaskTable
	err = dao.X.SQL("select * from task where id=?", taskId).Find(&taskTable)
	if err != nil {
		return
	}
	if len(taskTable) == 0 {
		return result, fmt.Errorf("Can not find any task with id:%s ", taskId)
	}
	result = *taskTable[0]
	return
}

func GetTaskListByRequestId(requestId string) (taskList []*models.TaskTable, err error) {
	err = dao.X.SQL("select * from task where request = ? order by created_time desc ", requestId).Find(&taskList)
	if err != nil {
		return
	}
	return
}

func ApproveTask(task models.TaskTable, operator, userToken, language string, param models.TaskApproveParam) error {
	var err error
	var taskSort int
	err = SaveTaskFormNew(&task, operator, &param)
	if err != nil {
		return err
	}
	taskSort = GetTaskService().GenerateTaskOrderByRequestId(task.Request)
	switch models.TaskType(task.Type) {
	case models.TaskTypeApprove:
		return handleApprove(task, operator, userToken, language, param, taskSort)
	case models.TaskTypeImplement:
		// 编排任务,走编排逻辑.
		if task.ProcDefKey != "" && task.ProcDefId != "" {
			return handleWorkflowTask(task, operator, userToken, param, language, taskSort)
		}
		// 处理自定义任务
		return handleCustomTask(task, operator, userToken, language, param, taskSort)
	}
	return nil
}

// handleApprove 处理审批
func handleApprove(task models.TaskTable, operator, userToken, language string, param models.TaskApproveParam, taskSort int) (err error) {
	var actions, newApproveActions []*dao.ExecAction
	var request models.RequestTable
	now := time.Now().Format(models.DateTimeFormat)
	request, err = GetSimpleRequest(task.Request)
	if err != nil {
		return
	}
	if requestTemplate, getTemplateErr := GetRequestTemplateService().GetRequestTemplate(request.RequestTemplate); getTemplateErr != nil {
		err = getTemplateErr
		return
	} else {
		if requestTemplate.ProcDefId != "" {
			request.AssociationWorkflow = true
		}
	}
	switch param.ChoseOption {
	case string(models.TaskHandleResultTypeApprove):
		// 当前审批通过,需要通过查看 task_template里面handle_mode 判断协同,并行
		var taskHandleList []*models.TaskHandleTable
		var taskTemplateList []*models.TaskTemplateTable
		err = dao.X.SQL("select * from task_template where id = ?", task.TaskTemplate).Find(&taskTemplateList)
		if err != nil {
			return
		}
		if len(taskTemplateList) == 0 {
			err = fmt.Errorf("task:%s taskTemplate is empty", task.Id)
			return
		}
		if taskTemplateList[0].HandleMode == string(models.TaskTemplateHandleModeAll) {
			// 并行模式,都要审批完成才能到下一步
			err = dao.X.SQL("select * from task_handle where task = ?", task.Id).Find(&taskHandleList)
			if err != nil {
				return
			}
			if len(taskHandleList) == 0 {
				err = fmt.Errorf("task:%s taskHandleList is empty", task.Id)
				return
			}
			for _, taskHandle := range taskHandleList {
				// 存在任务节点 没有审批通过,并且不是当前节点,更新当前处理节点为完成后,return 等待其他审批人处理
				if taskHandle.HandleResult != string(models.TaskHandleResultTypeApprove) && taskHandle.Id != param.TaskHandleId {
					_, err = dao.X.Exec("update task_handle set handle_result = ?,handle_status = ?,result_desc = ?,updated_time =? where id = ?", models.TaskHandleResultTypeApprove, models.TaskHandleResultTypeComplete, param.Comment, now, param.TaskHandleId)
					if err != nil {
						return
					}
					return
				}
			}
		}
		actions = append(actions, &dao.ExecAction{Sql: "update task_handle set handle_result = ?,handle_status = ?,result_desc = ?,updated_time =? where id= ?", Param: []interface{}{models.TaskHandleResultTypeApprove, models.TaskHandleResultTypeComplete, param.Comment, now, param.TaskHandleId}})
		actions = append(actions, &dao.ExecAction{Sql: "update task set status = ?,task_result = ?,updated_by =?,updated_time =? where id = ?", Param: []interface{}{models.TaskStatusDone, models.TaskHandleResultTypeApprove, operator, now, task.Id}})
		newApproveActions, _ = GetRequestService().CreateRequestApproval(request, task.Id, userToken, language, taskSort)
		if len(newApproveActions) > 0 {
			actions = append(actions, newApproveActions...)
		}
		err = dao.Transaction(actions)
		return
	case string(models.TaskHandleResultTypeDeny):
		// 拒绝, 任务处理结果设置为拒绝,请求状态设置自动退回
		actions = append(actions, &dao.ExecAction{Sql: "update task_handle set handle_result=?,handle_status=?,result_desc=?,updated_time=? where id = ?", Param: []interface{}{models.TaskHandleResultTypeDeny, models.TaskHandleResultTypeComplete, param.Comment, now, param.TaskHandleId}})
		actions = append(actions, &dao.ExecAction{Sql: "update task set status = ?,task_result=?,updated_by=?,updated_time=? where id = ?", Param: []interface{}{models.TaskStatusDone, models.TaskHandleResultTypeDeny, operator, now, task.Id}})
		actions = append(actions, &dao.ExecAction{Sql: "update request set status = ?,updated_by=?,updated_time=? where id = ?", Param: []interface{}{models.RequestStatusFaulted, operator, now, task.Request}})

	case string(models.TaskHandleResultTypeRedraw):
		// 退回,请求变草稿,任务设置为处理完成
		actions = append(actions, &dao.ExecAction{Sql: "update task_handle set handle_result=?,handle_status=?,result_desc=?,updated_time=? where id = ?", Param: []interface{}{models.TaskHandleResultTypeRedraw, models.TaskHandleResultTypeComplete, param.Comment, now, param.TaskHandleId}})
		actions = append(actions, &dao.ExecAction{Sql: "update task set status = ?,task_result=?,updated_by=?,updated_time=? where id = ?", Param: []interface{}{models.TaskStatusDone, models.TaskHandleResultTypeRedraw, operator, now, task.Id}})
		actions = append(actions, &dao.ExecAction{Sql: "update request set status = ?,updated_by=?,updated_time=? where id = ?", Param: []interface{}{models.RequestStatusDraft, operator, now, task.Request}})
	}
	if len(actions) > 0 {
		err = dao.Transaction(actions)
	}
	return
}

// handleCustomTask 处理自定义任务
func handleCustomTask(task models.TaskTable, operator, userToken, language string, param models.TaskApproveParam, taskSort int) (err error) {
	var actions, newApproveActions []*dao.ExecAction
	var request models.RequestTable
	now := time.Now().Format(models.DateTimeFormat)
	request, err = GetSimpleRequest(task.Request)
	if err != nil {
		return
	}
	actions = append(actions, &dao.ExecAction{Sql: "update task_handle set handle_result = ?,handle_status=?,result_desc = ?,updated_time =? where id= ?", Param: []interface{}{param.HandleStatus, models.TaskHandleResultTypeComplete, param.Comment, now, param.TaskHandleId}})
	actions = append(actions, &dao.ExecAction{Sql: "update task set status = ?,task_result = ?,updated_by =?,updated_time =? where id = ?", Param: []interface{}{models.TaskStatusDone, param.HandleStatus, operator, now, task.Id}})
	newApproveActions, err = GetRequestService().CreateRequestTask(request, task.Id, userToken, language, taskSort)
	if len(newApproveActions) > 0 {
		actions = append(actions, newApproveActions...)
	}
	err = dao.Transaction(actions)
	return
}

// handleWorkflowTask 处理编排任务
func handleWorkflowTask(task models.TaskTable, operator, userToken string, param models.TaskApproveParam, language string, taskSort int) error {
	var err error
	requestParam, callbackUrl, getDataErr := getApproveCallbackParamNew(task.Id)
	if getDataErr != nil {
		return getDataErr
	}
	if param.ChoseOption != "" {
		requestParam.ResultCode = param.ChoseOption
	}
	for _, v := range requestParam.Results.Outputs {
		v.Comment = param.Comment
	}
	var respResult models.CallbackResult
	requestBytes, _ := json.Marshal(requestParam)
	b, _ := rpc.HttpPost(models.Config.Wecube.BaseUrl+callbackUrl, userToken, "", requestBytes)
	log.Logger.Info("Callback response", log.String("body", string(b)))
	err = json.Unmarshal(b, &respResult)
	if err != nil {
		return fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	if respResult.Status != "OK" {
		//if strings.Contains(respResult.Message, "None process instance found") {
		//	dao.X.Exec("update task set status='done',updated_by=?,updated_time=? where id=?", operator, nowTime, task.Id)
		//}
		return fmt.Errorf("Callback fail,%s ", respResult.Message)
	}
	request, getRequestErr := GetSimpleRequest(task.Request)
	if getRequestErr != nil {
		return getRequestErr
	}
	var actions, newApproveActions []*dao.ExecAction
	actions = append(actions, &dao.ExecAction{Sql: "update task set callback_data=?,result=?,task_result=?,chose_option=?,status=?,updated_by=?,updated_time=? where id=?", Param: []interface{}{
		string(requestBytes), param.Comment, models.TaskResultTypeComplete, param.ChoseOption, "done", operator, nowTime, task.Id,
	}})
	actions = append(actions, &dao.ExecAction{Sql: "update task_handle set handle_result = ?,handle_status=?,result_desc = ?,updated_time =? where id= ?", Param: []interface{}{param.ChoseOption, param.HandleStatus, param.Comment, nowTime, param.TaskHandleId}})
	newApproveActions, err = GetRequestService().CreateProcessTask(request, &task, userToken, language, taskSort)
	if len(newApproveActions) > 0 {
		actions = append(actions, newApproveActions...)
	}
	err = dao.Transaction(actions)
	return err
}

func getApproveCallbackParam(taskId string) (result models.PluginTaskCreateResp, callbackUrl string, err error) {
	result = models.PluginTaskCreateResp{ResultCode: "0"}
	taskObj, tmpErr := getSimpleTask(taskId)
	if tmpErr != nil {
		return result, callbackUrl, tmpErr
	}
	callbackUrl = taskObj.CallbackUrl
	resultObj := models.PluginTaskCreateOutput{RequestId: taskObj.CallbackRequestId}
	if taskObj.Cache == "" {
		resultObj.Outputs = []*models.PluginTaskCreateOutputObj{{CallbackParameter: taskObj.CallbackParameter, Comment: taskObj.Result, ErrorCode: "0"}}
		result.Results = resultObj
		return
	}
	var taskFormOutput models.PluginTaskFormDto
	err = json.Unmarshal([]byte(taskObj.Cache), &taskFormOutput)
	if err != nil {
		return result, callbackUrl, fmt.Errorf("Try to json unmarshal cache data fail:%s ", err.Error())
	}
	var items []*models.TaskFormItemQueryObj
	err = dao.X.SQL("select t1.*,t2.attr_def_data_type,t2.element_type from form_item t1 left join form_item_template t2 on t1.form_item_template=t2.id where form=?", taskObj.Form).Find(&items)
	if err != nil {
		return result, callbackUrl, fmt.Errorf("Try to query form item fail:%s ", err.Error())
	}
	itemValueMap := make(map[string]interface{})
	for _, v := range items {
		itemValueMap[fmt.Sprintf("%s_%s", v.FormItemTemplate, v.RowDataId)] = v.Value
	}
	for _, formEntity := range taskFormOutput.FormDataEntities {
		for _, itemValueObj := range formEntity.FormItemValues {
			tmpKey := fmt.Sprintf("%s_%s", itemValueObj.FormItemMetaId, itemValueObj.Oid)
			itemValueObj.AttrValue = itemValueMap[tmpKey]
		}
	}
	formBytes, _ := json.Marshal(taskFormOutput)
	resultObj.Outputs = []*models.PluginTaskCreateOutputObj{{CallbackParameter: taskObj.CallbackParameter, Comment: taskObj.Result, ErrorCode: "0", TaskFormOutput: string(formBytes)}}
	result.Results = resultObj
	return
}

func getApproveCallbackParamNew(taskId string) (result models.PluginTaskCreateResp, callbackUrl string, err error) {
	result = models.PluginTaskCreateResp{ResultCode: "0"}
	taskObj, tmpErr := getSimpleTask(taskId)
	if tmpErr != nil {
		return result, callbackUrl, tmpErr
	}
	callbackUrl = taskObj.CallbackUrl
	resultObj := models.PluginTaskCreateOutput{RequestId: taskObj.CallbackRequestId}
	if taskObj.Cache == "" {
		resultObj.Outputs = []*models.PluginTaskCreateOutputObj{{CallbackParameter: taskObj.CallbackParameter, Comment: taskObj.Result, ErrorCode: "0"}}
		result.Results = resultObj
		return
	}
	var taskFormOutput models.PluginTaskFormDto
	err = json.Unmarshal([]byte(taskObj.Cache), &taskFormOutput)
	if err != nil {
		return result, callbackUrl, fmt.Errorf("Try to json unmarshal cache data fail:%s ", err.Error())
	}
	taskDataRows := models.RequestPoolDataQueryRows{}
	err = dao.X.SQL("select t1.id as form_id,t1.form_template,t3.item_group,t3.item_group_type ,t1.data_id,t2.id as form_item_id,t2.form_item_template,t2.name,t2.value,t2.updated_time from form t1 left join form_item t2 on t1.id=t2.form left join form_template t3 on t1.form_template=t3.id where t1.task=?", taskId).Find(&taskDataRows)
	if err != nil {
		return result, callbackUrl, fmt.Errorf("Try to query form item fail:%s ", err.Error())
	}
	poolForms := taskDataRows.DataParse()
	for _, formEntity := range taskFormOutput.FormDataEntities {
		for _, itemValueObj := range formEntity.FormItemValues {
			poolItem := getRequestPoolLatestItem(poolForms, itemValueObj.Oid, itemValueObj.AttrName)
			itemValueObj.AttrValue = poolItem.Value
		}
	}
	formBytes, _ := json.Marshal(taskFormOutput)
	resultObj.Outputs = []*models.PluginTaskCreateOutputObj{{CallbackParameter: taskObj.CallbackParameter, Comment: taskObj.Result, ErrorCode: "0", TaskFormOutput: string(formBytes)}}
	result.Results = resultObj
	return
}

func SaveTaskForm(taskId, operator string, param models.TaskApproveParam) error {
	var actions []*dao.ExecAction
	taskObj, err := getSimpleTask(taskId)
	if err != nil {
		return err
	}
	var existFormItemTable []*models.FormItemTable
	dao.X.SQL("select * from form_item where form in (select form from task where id=?)", taskId).Find(&existFormItemTable)
	existFormItemMap := make(map[string]int)
	inputItemMap := make(map[string]int)
	for _, v := range existFormItemTable {
		existFormItemMap[fmt.Sprintf("%s^%s^%s", v.ItemGroup, v.Name, v.RowDataId)] = 1
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	actions = append(actions, &dao.ExecAction{Sql: "update task set `result`=?,chose_option=?,updated_by=?,updated_time=? where id=?", Param: []interface{}{param.Comment, param.ChoseOption, operator, nowTime, taskId}})
	if taskObj.Request != "" {
		for _, tableForm := range param.FormData {
			tmpColumnMap := make(map[string]string)
			tmpMultiMap := make(map[string]int)
			for _, title := range tableForm.Title {
				tmpColumnMap[title.Name] = title.Id
				if title.Multiple == "Y" {
					tmpMultiMap[title.Name] = 1
				}
			}
			for _, valueObj := range tableForm.Value {
				if valueObj.Id == "" {
					valueObj.Id = fmt.Sprintf("tmp%s%s", models.SysTableIdConnector, guid.CreateGuid())
				}
				for k, v := range valueObj.EntityData {
					if titleId, b := tmpColumnMap[k]; b {
						tmpExistKey := fmt.Sprintf("%s^%s^%s", tableForm.ItemGroup, k, valueObj.Id)
						inputItemMap[tmpExistKey] = 1
						if _, bb := tmpMultiMap[k]; bb {
							tmpV := []string{}
							for _, interfaceV := range v.([]interface{}) {
								tmpV = append(tmpV, fmt.Sprintf("%s", interfaceV))
							}
							if _, bbb := existFormItemMap[tmpExistKey]; !bbb {
								actions = append(actions, &dao.ExecAction{Sql: "insert into form_item(id,form,form_item_template,name,value,item_group,row_data_id) value (?,?,?,?,?,?,?)", Param: []interface{}{guid.CreateGuid(), taskObj.Form, titleId, k, strings.Join(tmpV, ","), tableForm.ItemGroup, valueObj.Id}})
							} else {
								actions = append(actions, &dao.ExecAction{Sql: "update form_item set value=? where form=? and row_data_id=? and name=?", Param: []interface{}{strings.Join(tmpV, ","), taskObj.Form, valueObj.Id, k}})
							}
						} else {
							if _, bbb := existFormItemMap[tmpExistKey]; !bbb {
								actions = append(actions, &dao.ExecAction{Sql: "insert into form_item(id,form,form_item_template,name,value,item_group,row_data_id) value (?,?,?,?,?,?,?)", Param: []interface{}{guid.CreateGuid(), taskObj.Form, titleId, k, v, tableForm.ItemGroup, valueObj.Id}})
							} else {
								actions = append(actions, &dao.ExecAction{Sql: "update form_item set value=? where form=? and row_data_id=? and name=?", Param: []interface{}{v, taskObj.Form, valueObj.Id, k}})
							}
						}
					}
				}
			}
		}
		for _, row := range existFormItemTable {
			if _, ok := inputItemMap[fmt.Sprintf("%s^%s^%s", row.ItemGroup, row.Name, row.RowDataId)]; !ok {
				actions = append(actions, &dao.ExecAction{Sql: "delete from form_item where id=?", Param: []interface{}{row.Id}})
			}
		}
	}
	return dao.Transaction(actions)
}

func SaveTaskFormNew(task *models.TaskTable, operator string, param *models.TaskApproveParam) (err error) {
	var actions []*dao.ExecAction
	nowTime := time.Now().Format(models.DateTimeFormat)
	actions = append(actions, &dao.ExecAction{Sql: "update task set `result`=?,chose_option=?,updated_by=?,updated_time=? where id=?", Param: []interface{}{param.Comment, param.ChoseOption, operator, nowTime, task.Id}})
	if task.Request == "" {
		// 编排单独触发任务型
		return dao.Transaction(actions)
	}
	// 查请求数据池里的数据(里面的数据可能包含当前任务之前保存的数据)
	requestPoolRows := models.RequestPoolDataQueryRows{}
	if err = dao.X.SQL("select t1.id as form_id,t1.task,t1.form_template,t3.item_group,t3.item_group_type ,t1.data_id,t2.id as form_item_id,t2.form_item_template,t2.name,t2.value,t2.updated_time from form t1 left join form_item t2 on t1.id=t2.form left join form_template t3 on t1.form_template=t3.id where t1.request=?", task.Request).Find(&requestPoolRows); err != nil {
		return fmt.Errorf("query request item pool data fail,%s ", err.Error())
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
	for _, tableForm := range param.FormData {
		columnNameIdMap := make(map[string]string)
		isColumnMultiMap := make(map[string]int)
		for _, title := range tableForm.Title {
			columnNameIdMap[title.Name] = title.Id
			if title.Multiple == "Y" {
				isColumnMultiMap[title.Name] = 1
			}
		}
		poolForms := itemGroupFormMap[tableForm.ItemGroup]
		for _, valueObj := range tableForm.Value {
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
				actions = append(actions, &dao.ExecAction{Sql: "insert into form(id,request,task,form_template,data_id,created_by,updated_by,created_time,updated_time) values (?,?,?,?,?,?,?,?,?)", Param: []interface{}{
					formId, task.Request, task.Id, tableForm.FormTemplateId, valueObj.Id, operator, operator, nowTime, nowTime,
				}})
			}
			// 判断数据行属性的变化
			for k, v := range valueObj.EntityData {
				// 判断属性合不合法，是不是属性该表单的属性
				formItemTemplateId, nameLegalCheck := columnNameIdMap[k]
				if !nameLegalCheck {
					continue
				}
				// 整理属性值，特殊处理数组
				valueString := fmt.Sprintf("%s", v)
				if _, multipleFlag := isColumnMultiMap[k]; multipleFlag {
					if vInterfaceList, assertOk := v.([]interface{}); assertOk {
						tmpV := []string{}
						for _, interfaceV := range vInterfaceList {
							tmpV = append(tmpV, fmt.Sprintf("%s", interfaceV))
						}
						valueString = strings.Join(tmpV, ",")
					} else {
						err = fmt.Errorf("row:%s key:%s value:%v is not array,format to []interface{} fail", valueObj.Id, k, v)
						return
					}
				}
				// 从数据池里尝试查找有没有已存在的数据(同一个itemGroup，同一个数据行下的同一属性)
				latestPoolItem := getRequestPoolLatestItem(poolForms, valueObj.Id, k)
				if latestPoolItem.FormItemId == "" {
					// 没有在数据池里找到相关数据行的该属性
					actions = append(actions, &dao.ExecAction{Sql: "insert into form_item(id,form,form_item_template,name,value,request,updated_time,task_handle) values (?,?,?,?,?,?,?,?)", Param: []interface{}{
						"item_" + guid.CreateGuid(), formId, formItemTemplateId, k, valueString, task.Request, nowTime, param.TaskHandleId,
					}})
				} else {
					if latestPoolItem.Value != valueString {
						// 数据有更新
						if latestPoolItem.Task == task.Id {
							// 属于该任务，更新数据值
							actions = append(actions, &dao.ExecAction{Sql: "update form_item set value=?,updated_time=?,task_handle=? where id=?", Param: []interface{}{
								valueString, nowTime, param.TaskHandleId, latestPoolItem.FormItemId,
							}})
						} else {
							// 不属于该任务，新增数据纪录
							actions = append(actions, &dao.ExecAction{Sql: "insert into form_item(id,form,form_item_template,name,value,request,updated_time,original_id,task_handle) values (?,?,?,?,?,?,?,?,?)", Param: []interface{}{
								"item_" + guid.CreateGuid(), formId, formItemTemplateId, k, valueString, task.Request, nowTime, latestPoolItem.FormItemId, param.TaskHandleId,
							}})
						}
					}
				}
			}
		}
		// 如果之前是该任务保存的数据行但又没传过来了，说明已经删除行
		for _, poolForm := range poolForms {
			if poolForm.Task == task.Id {
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
	err = dao.Transaction(actions)
	if err != nil {
		err = fmt.Errorf("save task:%s form data fail,%s ", task.Id, err.Error())
	}
	return
}

func getRequestPoolLatestItem(poolForms []*models.RequestPoolForm, dataId, name string) (poolItem *models.RequestPoolDataQueryRow) {
	poolItem = &models.RequestPoolDataQueryRow{}
	for _, formObj := range poolForms {
		if formObj.DataId == dataId {
			for _, item := range formObj.Items {
				if item.Name == name {
					if poolItem.UpdatedTime.IsZero() {
						poolItem = item
					} else {
						if item.UpdatedTime.UnixMilli() > poolItem.UpdatedTime.UnixMilli() {
							poolItem = item
						}
					}
				}
			}
		}
	}
	return
}

func UpdateTaskHandle(param models.TaskHandleUpdateParam, operator string) (err error) {
	var task models.TaskTable
	var taskHandleList []*models.TaskHandleTable
	task, err = getSimpleTask(param.TaskId)
	if common.GetLowVersionUnixMillis(task.UpdatedTime) != param.LatestUpdateTime {
		err = exterror.New().DealWithAtTheSameTimeError
		return
	}
	if task.Status == string(models.TaskStatusDone) {
		err = fmt.Errorf("Task already done with %s %s ", task.UpdatedBy, task.UpdatedTime)
		return
	}
	dao.X.SQL("select * from task_handle where id = ?", param.TaskHandleId).Find(&taskHandleList)
	if len(taskHandleList) == 0 {
		err = fmt.Errorf("taskHandle is empty")
		return
	}
	var actions []*dao.ExecAction
	nowTime := time.Now().Format(models.DateTimeFormat)
	actions = append(actions, &dao.ExecAction{Sql: "update task set status=?,handler=?,updated_by=?,updated_time=? where id=?", Param: []interface{}{"marked",
		operator, operator, nowTime, param.TaskId}})
	//添加认领记录
	actions = append(actions, &dao.ExecAction{Sql: "insert into task_handle(id,task_handle_template,task,role,handler,handler_type,parent_id,created_time," +
		"updated_time,change_reason) values(?,?,?,?,?,?,?,?,?,?)", Param: []interface{}{guid.CreateGuid(), taskHandleList[0].TaskHandleTemplate, taskHandleList[0].Task,
		taskHandleList[0].Role, operator, taskHandleList[0].HandlerType, param.TaskHandleId, nowTime, nowTime, param.ChangeReason}})
	actions = append(actions, &dao.ExecAction{Sql: "update task_handle set latest_flag = 0 where id = ?", Param: []interface{}{param.TaskHandleId}})
	return dao.Transaction(actions)
}

func ChangeTaskStatus(taskId, operator, operation, lastedUpdateTime string) (taskObj models.TaskTable, err error) {
	taskObj, err = getSimpleTask(taskId)
	if err != nil {
		return
	}
	if common.GetLowVersionUnixMillis(taskObj.UpdatedTime) != lastedUpdateTime {
		err = exterror.New().DealWithAtTheSameTimeError
		return
	}
	if taskObj.Status == "done" {
		return taskObj, fmt.Errorf("Task aleary done with %s %s ", taskObj.UpdatedBy, taskObj.UpdatedTime)
	}
	var actions []*dao.ExecAction
	nowTime := time.Now().Format(models.DateTimeFormat)
	if operation == "mark" {
		actions = append(actions, &dao.ExecAction{Sql: "update task set status=?,handler=?,updated_by=?,updated_time=? where id=?", Param: []interface{}{"marked", operator, operator, nowTime, taskId}})
	} else if operation == "start" {
		if operator != taskObj.Handler {
			return taskObj, fmt.Errorf("Task handler is %s ", taskObj.Handler)
		}
		actions = append(actions, &dao.ExecAction{Sql: "update task set status=?,updated_by=?,updated_time=? where id=?", Param: []interface{}{"doing", operator, nowTime, taskId}})
	} else if operation == "quit" {
		if operator != taskObj.Handler {
			return taskObj, fmt.Errorf("Task handler is %s ", taskObj.Handler)
		}
		actions = append(actions, &dao.ExecAction{Sql: "update task set status=?,updated_by=?,updated_time=? where id=?", Param: []interface{}{"marked", operator, nowTime, taskId}})
	} else if operation == "give" {
		// 转给我
		if taskObj.Status == "done" {
			return taskObj, fmt.Errorf("Task status:%s is not marked ", taskObj.Status)
		}
		actions = append(actions, &dao.ExecAction{Sql: "update task set status=?,handler=?,updated_by=?,updated_time=? where id=?", Param: []interface{}{"marked", operator, operator, nowTime, taskId}})
	}
	err = dao.Transaction(actions)
	if err != nil {
		return taskObj, err
	}
	taskObj, _ = getSimpleTask(taskId)
	if taskObj.Status == "created" {
		taskObj.OperationOptions = []string{"mark"}
	} else if taskObj.Status == "marked" || taskObj.Status == "doing" {
		if taskObj.Handler == operator {
			taskObj.OperationOptions = []string{"start"}
		} else {
			taskObj.OperationOptions = []string{"mark"}
		}
	} else {
		taskObj.OperationOptions = []string{}
	}
	return taskObj, nil
}

func buildTaskOperation(taskObj *models.TaskListObj, operator string) {
	if taskObj.Status == "created" {
		taskObj.OperationOptions = []string{"mark"}
	} else if taskObj.Status == "marked" || taskObj.Status == "doing" {
		if taskObj.Handler == operator {
			taskObj.OperationOptions = []string{"start"}
		} else {
			taskObj.OperationOptions = []string{"mark"}
		}
	} else {
		taskObj.OperationOptions = []string{}
	}
}

func NotifyTaskMail(taskId, userToken, language, mailSubject, mailContent string) error {
	if !models.MailEnable {
		return nil
	}
	taskObj, _ := getSimpleTask(taskId)
	log.Logger.Info("Start notify task mail", log.String("taskId", taskId))
	var roleTable []*models.RoleTable
	dao.X.SQL("select id,email from `role` where id in (select `role` from task_template_role where role_type='USE' and task_template in (select task_template from task where id=?))", taskId).Find(&roleTable)
	reportRoleString := taskObj.ReportRole
	reportRoleString = strings.ReplaceAll(reportRoleString, "[", "")
	reportRoleString = strings.ReplaceAll(reportRoleString, "]", "")
	for _, v := range strings.Split(reportRoleString, ",") {
		if v != "" {
			roleTable = append(roleTable, &models.RoleTable{Id: v})
		}
	}
	if len(roleTable) == 0 {
		return fmt.Errorf("can not find handle role with task:%s ", taskId)
	}
	mailList := GetRoleService().GetRoleMail(roleTable, userToken, language)
	if len(mailList) == 0 {
		log.Logger.Warn("Notify task mail break,email is empty", log.String("role", roleTable[0].Id))
		return fmt.Errorf("handle role email is empty ")
	}
	var taskTable []*models.TaskTable
	dao.X.SQL("select t1.id,t1.name,t1.description,t2.name as request,t1.node_name,t1.emergency,t1.reporter,t1.created_time from task t1 left join request t2 on t1.request=t2.id where t1.id=?", taskId).Find(&taskTable)
	if len(taskTable) == 0 {
		return fmt.Errorf("can not find task with id:%s ", taskId)
	}
	var subject, content string
	subject = fmt.Sprintf("Taskman task [%s] %s[%s]", models.PriorityLevelMap[taskTable[0].Emergency], taskTable[0].Name, taskTable[0].Request)
	content = fmt.Sprintf("Taskman task \nID:%s \nPriority:%s \nName:%s \nRequest:%s \nDescription:%s \nReporter:%s \nCreateTime:%s \n", taskTable[0].Id, models.PriorityLevelMap[taskTable[0].Emergency], taskTable[0].Name, taskTable[0].Request, taskTable[0].Description, taskTable[0].Reporter, taskTable[0].CreatedTime)
	if mailSubject != "" {
		subject = mailSubject
	}
	if mailContent != "" {
		content = mailContent
	}
	err := models.MailSender.Send(subject, content, mailList)
	if err != nil {
		return fmt.Errorf("send notify email fail:%s ", err.Error())
	}
	return nil
}

func GetSimpleTask(taskId string) (task models.TaskTable, err error) {
	var taskTable []*models.TaskTable
	err = dao.X.SQL("select * from task where id=?", taskId).Find(&taskTable)
	if err != nil {
		return
	}
	if len(taskTable) == 0 {
		return task, fmt.Errorf("Can not find any task with id:%s ", taskId)
	}
	task = *taskTable[0]
	return
}

func (s *TaskService) CreateTasks(param []*models.TaskDto) (func(*xorm.Session) error, error) {
	if len(param) == 0 {
		return nil, nil
	}
	// 校验参数
	taskTemplateId := param[0].TaskTemplate
	requestId := param[0].Request
	taskType := param[0].Type
	guidPrefix, _ := GetTaskTemplateService().genTaskIdPrefix(taskType)
	taskHandleHandlerTypes := make([][]string, len(param))
	for i, task := range param {
		if task.TaskTemplate != taskTemplateId {
			return nil, errors.New("param taskTemplate not the same")
		}
		if task.Request != requestId {
			return nil, errors.New("param request not the same")
		}
		if task.Type != taskType {
			return nil, errors.New("param type not the same")
		}
		if len(task.Handles) == 0 {
			continue
		}
		taskHandleHandlerTypes[i] = make([]string, len(task.Handles))
		for j, handle := range task.Handles {
			taskHandleTemplateId := handle.TaskHandleTemplate
			// 查询任务处理模板
			taskHandleTemplate, err := GetTaskTemplateService().taskHandleTemplateDao.Get(taskHandleTemplateId)
			if err != nil {
				return nil, err
			}
			if taskHandleTemplate == nil {
				return nil, errors.New("no task_handle_template record found")
			}
			if taskHandleTemplate.TaskTemplate != taskTemplateId {
				return nil, errors.New("task_handle_template not match task_template")
			}
			taskHandleHandlerTypes[i][j] = taskHandleTemplate.HandlerType
		}
	}
	// 查询任务模板
	taskTemplate, err := GetTaskTemplateService().taskTemplateDao.Get(taskTemplateId)
	if err != nil {
		return nil, err
	}
	if taskTemplate == nil {
		return nil, errors.New("no task_template record found")
	}
	// 查询任务列表
	tasks, err := s.taskDao.QueryByRequest(requestId)
	if err != nil {
		return nil, err
	}
	if len(tasks) > 0 {
		return nil, errors.New("task record already exist")
	}
	// 构造返回结果
	result := func(session *xorm.Session) error {
		nowTime := time.Now().Format(models.DateTimeFormat)
		// 插入数据
		for i, task := range param {
			task.Id = guid.CreateGuid()
			task.Sort = i + 1
			// 插入任务
			newTask := &models.TaskTable{
				Id:           fmt.Sprintf("%s_%s", guidPrefix, guid.CreateGuid()),
				Name:         task.Name,
				Type:         task.Type,
				Sort:         i + 1,
				TaskTemplate: taskTemplateId,
				Request:      requestId,
			}
			_, err := s.taskDao.Add(session, newTask)
			if err != nil {
				return err
			}
			// 插入任务处理
			for j, handle := range task.Handles {
				newTaskHandle := &models.TaskHandleTable{
					Id:                 guid.CreateGuid(),
					TaskHandleTemplate: handle.TaskHandleTemplate,
					Task:               newTask.Id,
					Role:               handle.Role,
					Handler:            handle.Handler,
					HandlerType:        taskHandleHandlerTypes[i][j],
					CreatedTime:        nowTime,
					UpdatedTime:        nowTime,
				}
				_, err = s.taskHandleDao.Add(session, newTaskHandle)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
	return result, nil
}

func (s *TaskService) UpdateTasks(param []*models.TaskDto) (func(*xorm.Session) error, error) {
	if len(param) == 0 {
		return nil, nil
	}
	// 校验参数
	taskTemplateId := param[0].TaskTemplate
	requestId := param[0].Request
	taskType := param[0].Type
	guidPrefix, _ := GetTaskTemplateService().genTaskIdPrefix(taskType)
	taskHandleHandlerTypes := make([][]string, len(param))
	taskHandleHandlerTypeMap := make([]map[string]string, len(param))
	for i, task := range param {
		if task.TaskTemplate != taskTemplateId {
			return nil, errors.New("param taskTemplate not the same")
		}
		if task.Request != requestId {
			return nil, errors.New("param request not the same")
		}
		if task.Type != taskType {
			return nil, errors.New("param type not the same")
		}
		if task.Sort != i+1 {
			return nil, errors.New("param sort wrong")
		}
		if len(task.Handles) == 0 {
			continue
		}
		taskHandleHandlerTypes[i] = make([]string, len(task.Handles))
		taskHandleHandlerTypeMap[i] = make(map[string]string)
		for j, handle := range task.Handles {
			taskHandleTemplateId := handle.TaskHandleTemplate
			// 查询任务处理模板
			taskHandleTemplate, err := GetTaskTemplateService().taskHandleTemplateDao.Get(taskHandleTemplateId)
			if err != nil {
				return nil, err
			}
			if taskHandleTemplate == nil {
				return nil, errors.New("no task_template_role record found")
			}
			if taskHandleTemplate.TaskTemplate != taskTemplateId {
				return nil, errors.New("task_handle_template not match task_template")
			}
			taskHandleHandlerTypes[i][j] = taskHandleTemplate.HandlerType
			taskHandleHandlerTypeMap[i][taskHandleTemplateId] = taskHandleTemplate.HandlerType
		}
	}
	// 查询任务模板
	taskTemplate, err := GetTaskTemplateService().taskTemplateDao.Get(taskTemplateId)
	if err != nil {
		return nil, err
	}
	if taskTemplate == nil {
		return nil, errors.New("no task_template record found")
	}
	// 查询任务列表
	tasks, err := s.taskDao.QueryByRequestAndType(requestId, taskType)
	if err != nil {
		return nil, err
	}
	// 汇总任务列表
	taskIds := make([]string, len(tasks))
	for i, task := range tasks {
		taskIds[i] = task.Id
	}
	// 查询任务处理列表
	taskHandles, err := s.taskHandleDao.QueryByTasks(taskIds)
	if err != nil {
		return nil, err
	}
	// 汇总任务处理列表
	taskHandleMap := make(map[string][]*models.TaskHandleTable)
	taskId := ""
	for _, taskHandle := range taskHandles {
		if taskId != taskHandle.Task {
			taskId = taskHandle.Task
			taskHandleMap[taskId] = make([]*models.TaskHandleTable, 0)
		}
		taskHandleMap[taskId] = append(taskHandleMap[taskId], taskHandle)
	}
	// 对比请求和现有数据
	var newTasks, updateTasks []*models.TaskTable
	var newTaskHandles, updateTaskHandles []*models.TaskHandleTable
	var deleteTaskIds, deleteTaskHandleIds []string
	nowTime := time.Now().Format(models.DateTimeFormat)
	for i, task := range param {
		if task.Id == "" {
			task.Id = fmt.Sprintf("%s_%s", guidPrefix, guid.CreateGuid())
		}
		taskHandles, ok := taskHandleMap[task.Id]
		if !ok {
			// 插入任务
			newTask := &models.TaskTable{
				Id:           task.Id,
				Name:         task.Name,
				Request:      task.Request,
				TaskTemplate: task.TaskTemplate,
				Type:         task.Type,
				Sort:         task.Sort,
			}
			newTasks = append(newTasks, newTask)
			// 插入任务处理
			for j, handle := range task.Handles {
				newTaskHandle := &models.TaskHandleTable{
					Id:                 guid.CreateGuid(),
					TaskHandleTemplate: handle.TaskHandleTemplate,
					Task:               newTask.Id,
					Role:               handle.Role,
					Handler:            handle.Handler,
					HandlerType:        taskHandleHandlerTypes[i][j],
					CreatedTime:        nowTime,
					UpdatedTime:        nowTime,
					Sort:               j + 1,
				}
				newTaskHandles = append(newTaskHandles, newTaskHandle)
			}
		} else {
			// 更新任务
			updateTask := &models.TaskTable{
				Id:           task.Id,
				Name:         task.Name,
				TaskTemplate: task.TaskTemplate,
				Sort:         task.Sort,
			}
			updateTasks = append(updateTasks, updateTask)
			// 增删改任务处理
			for j, taskHandle := range taskHandles {
				if j < len(task.Handles) {
					handle := task.Handles[j]
					updateTaskHandle := &models.TaskHandleTable{
						Id:                 taskHandle.Id,
						TaskHandleTemplate: handle.TaskHandleTemplate,
						Role:               handle.Role,
						Handler:            handle.Handler,
						HandlerType:        taskHandleHandlerTypeMap[i][handle.TaskHandleTemplate],
						UpdatedTime:        nowTime,
						Sort:               j + 1,
					}
					updateTaskHandles = append(updateTaskHandles, updateTaskHandle)
				} else {
					deleteTaskHandleIds = append(deleteTaskHandleIds, taskHandle.Id)
				}
			}
			for sort := len(taskHandles) + 1; sort <= len(task.Handles); sort++ {
				handle := task.Handles[sort-1]
				newTaskHandle := &models.TaskHandleTable{
					Id:                 guid.CreateGuid(),
					TaskHandleTemplate: handle.TaskHandleTemplate,
					Task:               task.Id,
					Role:               handle.Role,
					Handler:            handle.Handler,
					HandlerType:        taskHandleHandlerTypeMap[i][handle.TaskHandleTemplate],
					CreatedTime:        nowTime,
					UpdatedTime:        nowTime,
					Sort:               sort,
				}
				newTaskHandles = append(newTaskHandles, newTaskHandle)
			}
			// 处理完，从任务处理列表删除
			delete(taskHandleMap, task.Id)
		}
	}
	// 删除剩余的任务和任务处理
	for taskId, taskHandles := range taskHandleMap {
		for _, taskHandle := range taskHandles {
			deleteTaskHandleIds = append(deleteTaskHandleIds, taskHandle.Id)
		}
		deleteTaskIds = append(deleteTaskIds, taskId)
	}
	// 构造返回结果
	result := func(session *xorm.Session) error {
		err := s.taskHandleDao.Deletes(session, deleteTaskHandleIds)
		if err != nil {
			return err
		}
		err = s.taskDao.Deletes(session, deleteTaskIds)
		if err != nil {
			return err
		}
		for _, newTask := range newTasks {
			_, err = s.taskDao.Add(session, newTask)
			if err != nil {
				return err
			}
		}
		for _, newTaskHandle := range newTaskHandles {
			_, err = s.taskHandleDao.Add(session, newTaskHandle)
			if err != nil {
				return err
			}
		}
		for _, updateTask := range updateTasks {
			err = s.taskDao.Update(session, updateTask)
			if err != nil {
				return err
			}
		}
		for _, updateTaskHandle := range updateTaskHandles {
			err = s.taskHandleDao.Update(session, updateTaskHandle)
			if err != nil {
				return err
			}
		}
		return nil
	}
	return result, nil
}

func (s *TaskService) DeleteTasks(requestId, typ string) (func(*xorm.Session) error, error) {
	// 查询任务列表
	var tasks []*models.TaskTable
	tasks, err := s.taskDao.QueryByRequestAndType(requestId, typ)
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, nil
	}
	// 汇总任务列表
	taskIds := make([]string, len(tasks))
	for i, task := range tasks {
		taskIds[i] = task.Id
	}
	// 构造返回结果
	result := func(session *xorm.Session) error {
		// 删除任务处理列表
		err := s.taskHandleDao.DeleteByTasks(session, taskIds)
		if err != nil {
			return err
		}
		// 删除任务列表
		err = s.taskDao.Deletes(session, taskIds)
		if err != nil {
			return err
		}
		return nil
	}
	return result, nil
}

func (s *TaskService) ListTasks(requestId, typ string) ([]*models.TaskDto, error) {
	// 查询任务列表
	tasks, err := s.taskDao.QueryByRequestAndType(requestId, typ)
	if err != nil {
		return nil, err
	}
	if len(tasks) == 0 {
		return nil, nil
	}
	// 汇总任务列表
	taskIds := make([]string, len(tasks))
	for i, task := range tasks {
		taskIds[i] = task.Id
	}
	// 查询任务处理列表
	taskHandles, err := s.taskHandleDao.QueryByTasks(taskIds)
	if err != nil {
		return nil, err
	}
	// 汇总任务处理列表
	taskHandleMap := make(map[string][]*models.TaskHandleTable)
	taskId := ""
	for _, taskHandle := range taskHandles {
		if taskId != taskHandle.Task {
			taskId = taskHandle.Task
			taskHandleMap[taskId] = make([]*models.TaskHandleTable, 0)
		}
		taskHandleMap[taskId] = append(taskHandleMap[taskId], taskHandle)
	}
	// 构造返回结果
	result := make([]*models.TaskDto, len(tasks))
	for i, task := range tasks {
		result[i] = &models.TaskDto{
			Id:           task.Id,
			Name:         task.Name,
			Type:         task.Type,
			Sort:         task.Sort,
			TaskTemplate: task.TaskTemplate,
			Request:      task.Request,
		}
		if taskHandles, ok := taskHandleMap[task.Id]; ok {
			result[i].Handles = make([]*models.TaskHandleDto, len(taskHandles))
			for j, taskHandle := range taskHandles {
				result[i].Handles[j] = &models.TaskHandleDto{
					Id:                 taskHandle.Id,
					Sort:               taskHandle.Sort,
					TaskHandleTemplate: taskHandle.TaskHandleTemplate,
					Role:               taskHandle.Role,
					Handler:            taskHandle.Handler,
					HandleResult:       taskHandle.HandleResult,
					UpdatedTime:        taskHandle.UpdatedTime,
				}
			}
		}
	}
	return result, nil
}

func (s *TaskService) ListImplementTasks(requestId string) (list []*models.TaskTable, err error) {
	list = []*models.TaskTable{}
	err = dao.X.SQL("select * from task where request = ? and type =?", requestId, models.TaskTypeImplement).Find(&list)
	return
}

func (s *TaskService) GetLatestCheckTask(requestId string) (task *models.TaskTable, err error) {
	var taskList []*models.TaskTable
	err = dao.X.SQL("select * from task where request = ? and type = ? order by sort desc limit 0,1", requestId, models.TaskTypeCheck).Find(&taskList)
	if err != nil {
		return
	}
	if len(taskList) > 0 {
		task = taskList[0]
	}
	return
}

func (s *TaskService) GetLatestTask(requestId string) (task *models.TaskTable, err error) {
	var taskList []*models.TaskTable
	err = dao.X.SQL("select * from task where request = ? order by sort desc limit 0,1", requestId).Find(&taskList)
	if err != nil {
		return
	}
	if len(taskList) > 0 {
		task = taskList[0]
	}
	return
}

func (s *TaskService) GetTaskMapByRequestId(requestId string) (taskMap map[string]*models.TaskTable, err error) {
	var taskList []*models.TaskTable
	taskMap = make(map[string]*models.TaskTable)
	err = dao.X.SQL("select * from task where request = ? order by sort desc", requestId).Find(&taskList)
	if err != nil {
		return
	}
	if len(taskList) > 0 {
		for _, task := range taskList {
			if task.TaskTemplate != "" {
				if _, ok := taskMap[task.TaskTemplate]; !ok {
					taskMap[task.TaskTemplate] = task
				}
			}
		}
	}
	return
}

func (s *TaskService) GetDoingTaskByRequestIdAndType(requestId string, taskType models.TaskType) (task *models.TaskTable, err error) {
	var taskList []*models.TaskTable
	err = dao.X.SQL("select * from task where request = ? and type = ? and status != ?", requestId, taskType, models.TaskStatusDone).Find(&taskList)
	if err != nil {
		return
	}
	if len(taskList) > 0 {
		task = taskList[0]
	}
	return
}

func (s *TaskService) GetDoingTaskByRequestId(requestId string) (task *models.TaskTable, err error) {
	var taskList []*models.TaskTable
	err = dao.X.SQL("select * from task where request = ?  and status != ?", requestId, models.TaskStatusDone).Find(&taskList)
	if err != nil {
		return
	}
	if len(taskList) > 0 {
		task = taskList[0]
	}
	return
}

// GetDoneTaskByRequestId 已完成对任务,取最后一次提交请求后的完成任务
func (s *TaskService) GetDoneTaskByRequestId(requestId string) (taskList []*models.TaskTable, err error) {
	var allTaskList []*models.TaskTable
	taskList = []*models.TaskTable{}
	err = dao.X.SQL("select * from task where request = ? order by sort desc", requestId).Find(&allTaskList)
	if err != nil {
		return
	}
	if len(allTaskList) > 0 {
		for _, task := range allTaskList {
			if task.Type == string(models.TaskTypeSubmit) {
				break
			}
			if task.Status == string(models.TaskStatusDone) {
				taskList = append(taskList, task)
			}
		}
	}
	return
}

// GenerateTaskOrderByRequestId 获取任务最新sort
func (s *TaskService) GenerateTaskOrderByRequestId(requestId string) (order int) {
	order = 1
	var latestTaskList []*models.TaskTable
	dao.X.SQL("select * from task where request = ?  order by sort desc limit 0,1", requestId).Find(&latestTaskList)
	if len(latestTaskList) > 0 {
		order = latestTaskList[0].Sort + 1
	}
	return
}
