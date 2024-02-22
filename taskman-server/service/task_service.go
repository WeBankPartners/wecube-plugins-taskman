package service

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/exterror"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/rpc"
	"strconv"
	"strings"
	"time"
)

type TaskService struct {
	taskDao dao.TaskDao
	taskRoleDao dao.TaskRoleDao
}

func GetTaskFormStruct(procInstId, nodeDefId string) (result models.TaskMetaResult, err error) {
	result = models.TaskMetaResult{Status: "OK", Message: "Success"}
	var items []*models.FormItemTemplateTable
	err = dao.X.SQL("select * from form_item_template where form_template in (select form_template from task_template where node_def_id=? and request_template in (select request_template from request where proc_instance_id=?))", nodeDefId, procInstId).Find(&items)
	if err != nil {
		return
	}
	if len(items) == 0 {
		err = fmt.Errorf("Can not find task item template with procInstId:%s nodeDefId:%s ", procInstId, nodeDefId)
		return
	}
	resultData := models.TaskMetaResultData{FormMetaId: items[0].FormTemplate}
	for _, item := range items {
		if item.Entity == "" {
			continue
		}
		resultData.FormItemMetas = append(resultData.FormItemMetas, &models.TaskMetaResultItem{FormItemMetaId: item.Id, PackageName: item.PackageName, EntityName: item.Entity, AttrName: item.Name})
	}
	result.Data = resultData
	return
}

func PluginTaskCreate(input *models.PluginTaskCreateRequestObj, callRequestId, dueDate string, nextOptions []string) (result *models.PluginTaskCreateOutputObj, taskId string, err error) {
	log.Logger.Debug("task create", log.JsonObj("input", input))
	result = &models.PluginTaskCreateOutputObj{CallbackParameter: input.CallbackParameter, ErrorCode: "0", ErrorMessage: "", Comment: ""}
	var requestTable []*models.RequestTable
	err = dao.X.SQL("select id,form,request_template,emergency,type from request where proc_instance_id=?", input.ProcInstId).Find(&requestTable)
	if err != nil {
		return result, taskId, fmt.Errorf("Try to check proc_instance_id:%s is in request fail,%s ", input.ProcInstId, err.Error())
	}
	var actions []*dao.ExecAction
	nowTime := time.Now().Format(models.DateTimeFormat)
	newTaskFormObj := models.FormTable{Id: guid.CreateGuid(), Name: "form_" + input.TaskName}
	input.RoleName = remakeTaskReportRole(input.RoleName)
	newTaskObj := models.TaskTable{Id: guid.CreateGuid(), Name: input.TaskName, Status: "created", Form: newTaskFormObj.Id, Reporter: input.Reporter, ReportRole: input.RoleName, Description: input.TaskDescription, CallbackUrl: input.CallbackUrl, CallbackParameter: input.CallbackParameter, NextOption: strings.Join(nextOptions, ","), Handler: input.Handler}
	taskId = newTaskObj.Id
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
		// Custom task create
		customExpireTime := ""
		dueMin, _ := strconv.Atoi(dueDate)
		if dueMin > 0 {
			customExpireTime = time.Now().Add(time.Duration(dueMin) * time.Minute).Format(models.DateTimeFormat)
		}
		taskInsertAction := dao.ExecAction{Sql: "insert into task(id,name,description,status,proc_def_id,proc_def_key,node_def_id,node_name,callback_url,callback_parameter,reporter,report_role,report_time,expire_time,emergency,callback_request_id,next_option,handler,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
		taskInsertAction.Param = []interface{}{newTaskObj.Id, newTaskObj.Name, newTaskObj.Description, newTaskObj.Status, newTaskObj.ProcDefId, newTaskObj.ProcDefKey, newTaskObj.NodeDefId, newTaskObj.NodeName, newTaskObj.CallbackUrl, newTaskObj.CallbackParameter, newTaskObj.Reporter, newTaskObj.ReportRole, nowTime, customExpireTime, newTaskObj.Emergency, callRequestId, newTaskObj.NextOption, newTaskObj.Handler, "system", nowTime, "system", nowTime}
		actions = append(actions, &taskInsertAction)
		err = dao.Transaction(actions)
		return
	}
	itemTemplateMap := getFormItemTemplateMap(taskFormInput.FormMetaId)
	newTaskObj.ProcDefId = taskFormInput.ProcDefId
	newTaskObj.ProcDefKey = taskFormInput.ProcDefKey
	if len(requestTable) > 0 {
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
			newTaskFormObj.FormTemplate = taskTemplateTable[0].FormTemplate
		} else {
			log.Logger.Warn("Can not find any taskTemplate", log.String("requestTemplate", requestTable[0].RequestTemplate), log.String("nodeDefId", taskFormInput.TaskNodeDefId))
			err = fmt.Errorf("Can not find any taskTemplate in request:%s with nodeDefId:%s ", newTaskObj.Request, taskFormInput.TaskNodeDefId)
			return
		}
	}
	log.Logger.Debug("debug1", log.JsonObj("newTaskFormObj", newTaskFormObj))
	taskInsertAction := dao.ExecAction{Sql: "insert into task(id,name,description,form,status,request,task_template,proc_def_id,proc_def_key,node_def_id," +
		"node_name,callback_url,callback_parameter,reporter,report_role,report_time,emergency,cache,callback_request_id,next_option,expire_time," +
		"handler,created_by,created_time,updated_by,updated_time,template_type) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
	taskInsertAction.Param = []interface{}{newTaskObj.Id, newTaskObj.Name, newTaskObj.Description, newTaskObj.Form, newTaskObj.Status,
		newTaskObj.Request, newTaskObj.TaskTemplate, newTaskObj.ProcDefId, newTaskObj.ProcDefKey, newTaskObj.NodeDefId, newTaskObj.NodeName,
		newTaskObj.CallbackUrl, newTaskObj.CallbackParameter, newTaskObj.Reporter, newTaskObj.ReportRole, nowTime, newTaskObj.Emergency,
		input.TaskFormInput, callRequestId, newTaskObj.NextOption, newTaskObj.ExpireTime, newTaskObj.Handler, "system", nowTime, "system",
		nowTime, newTaskObj.TemplateType}
	actions = append(actions, &taskInsertAction)
	actions = append(actions, &dao.ExecAction{Sql: "insert into form(id,name,form_template) value (?,?,?)", Param: []interface{}{newTaskFormObj.Id, newTaskFormObj.Name, newTaskFormObj.FormTemplate}})
	for _, formDataEntity := range taskFormInput.FormDataEntities {
		tmpGuidList := guid.CreateGuidList(len(formDataEntity.FormItemValues))
		for i, formDataItem := range formDataEntity.FormItemValues {
			tmpItemGroup := itemTemplateMap[formDataItem.FormItemMetaId].ItemGroup
			actions = append(actions, &dao.ExecAction{Sql: "insert into form_item(id,form,form_item_template,name,value,item_group,row_data_id) value (?,?,?,?,?,?,?)", Param: []interface{}{tmpGuidList[i], newTaskFormObj.Id, formDataItem.FormItemMetaId, formDataItem.AttrName, formDataItem.AttrValue, tmpItemGroup, formDataItem.Oid}})
		}
	}
	customItems, tmpErr := getLastCustomFormItem(newTaskObj.Request, newTaskFormObj.FormTemplate, newTaskFormObj.Id)
	if tmpErr != nil {
		err = fmt.Errorf("Try to get custom item fail,%s ", tmpErr.Error())
		return
	}
	if len(customItems) > 0 {
		tmpGuidList := guid.CreateGuidList(len(customItems))
		for i, customItem := range customItems {
			actions = append(actions, &dao.ExecAction{Sql: "insert into form_item(id,form,form_item_template,name,value,item_group,row_data_id) value (?,?,?,?,?,?,?)", Param: []interface{}{tmpGuidList[i], customItem.Form, customItem.FormItemTemplate, customItem.Name, customItem.Value, customItem.ItemGroup, customItem.RowDataId}})
		}
	}
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
	err = dao.X.SQL("select * from form_item_template where form_template in (select form_template from task_template where id=?) order by item_group,sort", taskObj.TaskTemplate).Find(&itemTemplates)
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

func getTaskMapByRequestId(requestId string) (taskMap map[string]*models.TaskTable, err error) {
	taskMap = make(map[string]*models.TaskTable)
	var taskTable []*models.TaskTable
	err = dao.X.SQL("select * from task where request = ?", requestId).Find(&taskTable)
	if err != nil {
		return
	}
	if len(taskTable) == 0 {
		err = fmt.Errorf("Can not find any task with request:%s ", requestId)
		return
	}
	for _, task := range taskTable {
		taskMap[task.NodeDefId] = task
	}
	return
}

func ApproveTask(taskId, operator, userToken string, param models.TaskApproveParam) error {
	err := SaveTaskForm(taskId, operator, param)
	if err != nil {
		return err
	}
	requestParam, callbackUrl, getDataErr := getApproveCallbackParam(taskId)
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
		if strings.Contains(respResult.Message, "None process instance found") {
			dao.X.Exec("update task set status='done',updated_by=?,updated_time=? where id=?", operator, nowTime, taskId)
		}
		return fmt.Errorf("Callback fail,%s ", respResult.Message)
	}
	_, err = dao.X.Exec("update task set callback_data=?,result=?,chose_option=?,status=?,updated_by=?,updated_time=? where id=?", string(requestBytes), param.Comment, param.ChoseOption, "done", operator, nowTime, taskId)
	if err != nil {
		return fmt.Errorf("Callback succeed,but update database task row fail,%s ", err.Error())
	}
	return nil
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

func NotifyTaskMail(taskId, userToken, language string) error {
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

func (s TaskService) CreateTasks(param []*models.TaskDto) ([]*dao.ExecAction, error) {
	actions := []*dao.ExecAction{}
	return actions, nil
}

func (s TaskService) UpdateTasks(param []*models.TaskDto) ([]*dao.ExecAction, error) {
	actions := []*dao.ExecAction{}
	return actions, nil
}

func (s TaskService) DeleteTasks(param []*models.TaskDto) ([]*dao.ExecAction, error) {
	actions := []*dao.ExecAction{}
	return actions, nil
}

func (s TaskService) ListTask(requestId string) ([]*models.TaskDto, error) {
	result := []*models.TaskDto{}
	return result, nil
}
