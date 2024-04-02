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
		result.Data = models.TaskMetaResultData{}
		//err = fmt.Errorf("Can not find task item template with procInstId:%s nodeDefId:%s ", procInstId, nodeDefId)
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

func PluginTaskCreateNew(input *models.PluginTaskCreateRequestObj, callRequestId, dueDate string, nextOptions []string, userToken, language string) (result *models.PluginTaskCreateOutputObj, task models.TaskTable, err error) {
	log.Logger.Debug("task create", log.JsonObj("input", input))
	result = &models.PluginTaskCreateOutputObj{CallbackParameter: input.CallbackParameter, ErrorCode: "0", ErrorMessage: "", Comment: ""}
	var requestTable []*models.RequestTable
	err = dao.X.SQL("select * from request where proc_instance_id=?", input.ProcInstId).Find(&requestTable)
	if err != nil {
		return result, models.TaskTable{}, fmt.Errorf("Try to check proc_instance_id:%s is in request fail,%s ", input.ProcInstId, err.Error())
	}
	var actions []*dao.ExecAction
	var taskSort int
	nowTime := time.Now().Format(models.DateTimeFormat)
	input.RoleName = remakeTaskReportRole(input.RoleName)
	task = models.TaskTable{Id: "tk_" + guid.CreateGuid(), Name: input.TaskName, Status: "created", Reporter: input.Reporter, ReportRole: input.RoleName, Description: input.TaskDescription, CallbackUrl: input.CallbackUrl, CallbackParameter: input.CallbackParameter, NextOption: strings.Join(nextOptions, ","), Handler: input.Handler}
	operator := "system"
	var taskFormInput models.PluginTaskFormDto
	if input.TaskFormInput != "" {
		err = json.Unmarshal([]byte(input.TaskFormInput), &taskFormInput)
		if err != nil {
			return result, task, fmt.Errorf("Try to json unmarshal taskFormInput to json data fail,%s ", err.Error())
		}
		if task.Reporter == "" {
			task.Reporter = "taskman"
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
		taskInsertAction.Param = []interface{}{task.Id, task.Name, task.Description, task.Status, task.ProcDefId,
			task.ProcDefKey, task.NodeDefId, task.NodeName, task.CallbackUrl, task.CallbackParameter, task.Reporter,
			task.ReportRole, nowTime, customExpireTime, task.Emergency, callRequestId, task.NextOption, task.Handler, operator,
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
	task.ProcDefId = taskFormInput.ProcDefId
	task.ProcDefKey = taskFormInput.ProcDefKey
	task.Request = requestTable[0].Id
	task.NodeDefId = taskFormInput.TaskNodeDefId
	task.Emergency = requestTable[0].Emergency
	task.TemplateType = requestTable[0].Type
	var taskTemplateTable []*models.TaskTemplateTable
	dao.X.SQL("select * from task_template where request_template=? and node_def_id=?", requestTable[0].RequestTemplate, taskFormInput.TaskNodeDefId).Find(&taskTemplateTable)
	if len(taskTemplateTable) > 0 {
		task.TaskTemplate = taskTemplateTable[0].Id
		task.NodeName = taskTemplateTable[0].NodeName
		task.ExpireTime = calcExpireTime(nowTime, taskTemplateTable[0].ExpireDay)
		task.Name = taskTemplateTable[0].Name
		task.Description = taskTemplateTable[0].Description
		task.Reporter = requestTable[0].Reporter
		task.ReportTime = nowTime
		task.Handler = taskTemplateTable[0].Handler
	} else {
		log.Logger.Warn("Can not find any taskTemplate", log.String("requestTemplate", requestTable[0].RequestTemplate), log.String("nodeDefId", taskFormInput.TaskNodeDefId))
		err = fmt.Errorf("Can not find any taskTemplate in request:%s with nodeDefId:%s ", task.Request, taskFormInput.TaskNodeDefId)
		return
	}
	taskInsertAction := dao.ExecAction{Sql: "insert into task(id,name,description,form,status,request,task_template,proc_def_id,proc_def_key,node_def_id," +
		"node_name,callback_url,callback_parameter,reporter,report_role,report_time,emergency,cache,callback_request_id,next_option,expire_time," +
		"handler,created_by,created_time,updated_by,updated_time,template_type,type,sort,request_created_time) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
	taskInsertAction.Param = []interface{}{task.Id, task.Name, task.Description, task.Form, task.Status,
		task.Request, task.TaskTemplate, task.ProcDefId, task.ProcDefKey, task.NodeDefId, task.NodeName,
		task.CallbackUrl, task.CallbackParameter, task.Reporter, task.ReportRole, nowTime, task.Emergency,
		input.TaskFormInput, callRequestId, task.NextOption, task.ExpireTime, task.Handler, operator, nowTime, operator,
		nowTime, task.TemplateType, models.TaskTypeImplement, taskSort, requestTable[0].CreatedTime}
	actions = append(actions, &taskInsertAction)

	var formTemplateRows []*models.FormTemplateTable
	err = dao.X.SQL("select * from form_template where task_template=? and item_group_type='workflow'", task.TaskTemplate).Find(&formTemplateRows)
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
			log.Logger.Warn("form data entity can not find form template", log.String("task", task.Id), log.JsonObj("formDataEntity", formDataEntity))
			continue
		}
		newFormId := "form_" + guid.CreateGuid()
		actions = append(actions, &dao.ExecAction{Sql: "insert into form(id,request,task,form_template,data_id,created_by,updated_by,created_time,updated_time) values (?,?,?,?,?,?,?,?,?)", Param: []interface{}{
			newFormId, task.Request, task.Id, tmpFormTemplateId, formDataEntity.Oid, operator, operator, nowTime, nowTime,
		}})
		for _, formDataItem := range formDataEntity.FormItemValues {
			actions = append(actions, &dao.ExecAction{Sql: "insert into form_item(id,form,form_item_template,name,value,request,updated_time) values (?,?,?,?,?,?,?)", Param: []interface{}{
				"item_" + guid.CreateGuid(), newFormId, formDataItem.FormItemMetaId, formDataItem.AttrName, formDataItem.AttrValue, task.Request, nowTime,
			}})
		}
	}
	createTaskHandleAction := GetTaskHandleService().CreateTaskHandleByTemplate(task.Id, userToken, language, requestTable[0], taskTemplateTable[0])
	actions = append(actions, createTaskHandleAction...)
	err = dao.TransactionWithoutForeignCheck(actions)
	return
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
		taskForm.FormData = []*models.RequestPreDataTableObj{}
		//return taskForm, fmt.Errorf("Can not find any form item template with task:%s ", taskObj.Id)
		return
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
			return handleWorkflowTask(task, operator, userToken, param, language)
		}
		// 处理自定义任务
		return handleCustomTask(task, operator, userToken, language, param, taskSort)
	}
	return nil
}

func ApproveCustomTask(task models.TaskTable, operator, userToken, language string, param models.TaskApproveParam) (err error) {
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
	b, _ := rpc.HttpPost(models.Config.Wecube.BaseUrl+callbackUrl, userToken, language, requestBytes)
	log.Logger.Info("Custom Callback response", log.String("body", string(b)))
	err = json.Unmarshal(b, &respResult)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	if respResult.Status != "OK" {
		err = fmt.Errorf("Callback fail,%s ", respResult.Message)
		return
	}
	var actions []*dao.ExecAction
	nowTime := time.Now().Format(models.DateTimeFormat)
	actions = append(actions, &dao.ExecAction{Sql: "update task set `result`=?,chose_option=?,updated_by=?,updated_time=? where id=?", Param: []interface{}{param.Comment, param.ChoseOption, operator, nowTime, task.Id}})
	if err = dao.Transaction(actions); err != nil {
		err = fmt.Errorf("update task message fail,%s ", err.Error())
	}
	return
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
			err = dao.X.SQL("select * from task_handle where task = ? and latest_flag = 1", task.Id).Find(&taskHandleList)
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
		actions = append(actions, &dao.ExecAction{Sql: "update task set status = ?,task_result = ?,updated_by =?,updated_time =? where id = ?", Param: []interface{}{models.TaskStatusDone, GetTaskHandleService().CalcTaskResult(task.Id, param.TaskHandleId), operator, now, task.Id}})
		newApproveActions, _ = GetRequestService().CreateRequestApproval(request, task.Id, userToken, language, taskSort, false)
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
		go NotifyTaskDenyMail(request.Name, task.Name, request.CreatedBy, operator, userToken, language)
	case string(models.TaskHandleResultTypeRedraw):
		// 退回,请求变草稿,任务设置为处理完成
		actions = append(actions, &dao.ExecAction{Sql: "update task_handle set handle_result=?,handle_status=?,result_desc=?,updated_time=? where id = ?", Param: []interface{}{models.TaskHandleResultTypeRedraw, models.TaskHandleResultTypeComplete, param.Comment, now, param.TaskHandleId}})
		actions = append(actions, &dao.ExecAction{Sql: "update task set status = ?,task_result=?,description=?,updated_by=?,updated_time=? where id = ?", Param: []interface{}{models.TaskStatusDone, models.TaskHandleResultTypeRedraw, param.Comment, operator, now, task.Id}})
		actions = append(actions, &dao.ExecAction{Sql: "update request set status = ?,rollback_desc=?,updated_by=?,updated_time=? where id = ?", Param: []interface{}{models.RequestStatusDraft, param.Comment, operator, now, task.Request}})
		go NotifyTaskBackMail(request.Name, task.Name, request.CreatedBy, operator, userToken, language)
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
	actions = append(actions, &dao.ExecAction{Sql: "update task_handle set handle_result = ?,handle_status=?,result_desc = ?,updated_time =? where id= ?", Param: []interface{}{param.ChoseOption, param.HandleStatus, param.Comment, now, param.TaskHandleId}})
	if param.HandleStatus == string(models.TaskHandleResultTypeUncompleted) {
		// 任务处理未完成,直接更新任务为未完成
		actions = append(actions, &dao.ExecAction{Sql: "update task set status = ?,task_result = ?,updated_by =?,updated_time =? where id = ?", Param: []interface{}{models.TaskStatusDone, string(models.TaskHandleResultTypeUncompleted), operator, now, task.Id}})
	} else {
		actions = append(actions, &dao.ExecAction{Sql: "update task set status = ?,task_result = ?,updated_by =?,updated_time =? where id = ?", Param: []interface{}{models.TaskStatusDone, GetTaskHandleService().CalcTaskResult(task.Id, param.TaskHandleId), operator, now, task.Id}})
	}
	newApproveActions, err = GetRequestService().CreateRequestTask(request, task.Id, userToken, language, taskSort)
	if len(newApproveActions) > 0 {
		actions = append(actions, newApproveActions...)
	}
	err = dao.Transaction(actions)
	return
}

// handleWorkflowTask 处理编排任务
func handleWorkflowTask(task models.TaskTable, operator, userToken string, param models.TaskApproveParam, language string) error {
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
		return fmt.Errorf("Callback fail,%s ", respResult.Message)
	}
	request, getRequestErr := GetSimpleRequest(task.Request)
	if getRequestErr != nil {
		return getRequestErr
	}
	var actions, newApproveActions []*dao.ExecAction
	actions = append(actions, &dao.ExecAction{Sql: "update task_handle set handle_result = ?,handle_status=?,result_desc = ?,updated_time =? where id= ?", Param: []interface{}{param.ChoseOption, param.HandleStatus, param.Comment, nowTime, param.TaskHandleId}})
	if param.HandleStatus == string(models.TaskHandleResultTypeUncompleted) {
		actions = append(actions, &dao.ExecAction{Sql: "update task set callback_data=?,result=?,task_result=?,chose_option=?,status=?,updated_by=?,updated_time=? where id=?", Param: []interface{}{
			string(requestBytes), param.Comment, models.TaskHandleResultTypeUncompleted, param.ChoseOption, "done", operator, nowTime, task.Id,
		}})
	} else {
		actions = append(actions, &dao.ExecAction{Sql: "update task set callback_data=?,result=?,task_result=?,chose_option=?,status=?,updated_by=?,updated_time=? where id=?", Param: []interface{}{
			string(requestBytes), param.Comment, GetTaskHandleService().CalcTaskResult(task.Id, param.TaskHandleId), param.ChoseOption, "done", operator, nowTime, task.Id,
		}})
	}
	newApproveActions, err = GetRequestService().CreateProcessTask(request, &task, userToken, language)
	if len(newApproveActions) > 0 {
		actions = append(actions, newApproveActions...)
	}
	err = dao.Transaction(actions)
	return err
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

func SaveTaskFormNew(task *models.TaskTable, operator string, param *models.TaskApproveParam) (err error) {
	var actions []*dao.ExecAction
	nowTime := time.Now().Format(models.DateTimeFormat)
	if param.ChoseOption != "" {
		actions = append(actions, &dao.ExecAction{Sql: "update task set `result`=?,chose_option=?,updated_by=?,updated_time=? where id=?", Param: []interface{}{param.Comment, param.ChoseOption, operator, nowTime, task.Id}})
		if task.Request == "" {
			// 编排单独触发任务型
			return dao.Transaction(actions)
		}
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
						var tmpV []string
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
						} else if item.UpdatedTime.UnixMilli() == poolItem.UpdatedTime.UnixMilli() {
							if item.FormItemId > poolItem.FormItemId {
								poolItem = item
							}
						}
					}
				}
			}
		}
	}
	return
}

func UpdateTaskHandle(param models.TaskHandleUpdateParam, operator, userToken, language string) (err error) {
	var task models.TaskTable
	var request models.RequestTable
	var taskHandleList []*models.TaskHandleTable
	var requestName string
	if task, err = getSimpleTask(param.TaskId); err != nil {
		return
	}
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
	if taskHandleList[0].TaskHandleTemplate != "" {
		actions = append(actions, &dao.ExecAction{Sql: "insert into task_handle(id,task_handle_template,task,role,handler,handler_type,parent_id,created_time," +
			"updated_time,change_reason) values(?,?,?,?,?,?,?,?,?,?)", Param: []interface{}{guid.CreateGuid(), taskHandleList[0].TaskHandleTemplate, taskHandleList[0].Task,
			taskHandleList[0].Role, operator, taskHandleList[0].HandlerType, param.TaskHandleId, nowTime, nowTime, param.ChangeReason}})
	} else {
		actions = append(actions, &dao.ExecAction{Sql: "insert into task_handle(id,task,role,handler,handler_type,parent_id,created_time," +
			"updated_time,change_reason) values(?,?,?,?,?,?,?,?,?)", Param: []interface{}{guid.CreateGuid(), taskHandleList[0].Task,
			taskHandleList[0].Role, operator, taskHandleList[0].HandlerType, param.TaskHandleId, nowTime, nowTime, param.ChangeReason}})
	}
	actions = append(actions, &dao.ExecAction{Sql: "update task_handle set latest_flag = 0 where id = ?", Param: []interface{}{param.TaskHandleId}})
	// 给原处理人发送邮件
	if task.Request != "" {
		if request, err = GetSimpleRequest(task.Request); err != nil {
			return
		}
		requestName = request.Name
	}
	go NotifyTaskHandlerUpdateMail(requestName, task.Name, taskHandleList[0].Handler, operator, userToken, language)
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
				return nil, errors.New("no task_handle_template record found")
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
	err = dao.X.SQL("select * from task where request = ? and type = ? and del_flag = 0 order by sort desc limit 0,1", requestId, models.TaskTypeCheck).Find(&taskList)
	if err != nil {
		return
	}
	if len(taskList) > 0 {
		task = taskList[0]
	}
	return
}

// GetDoingTask 任务存在并行情况,并行情况传递了taskId就取taskId的任务,没传递按照任务模版排序
func (s *TaskService) GetDoingTask(requestId, templateId, taskId string) (task *models.TaskTable, err error) {
	var taskList []*models.TaskTable
	var taskTemplateList []*models.TaskTemplateTable
	var doingTaskMap = make(map[string]*models.TaskTable)
	if taskId != "" {
		if err = dao.X.SQL("select * from task where id = ? and status <> 'done' and del_flag = 0", taskId).Find(&taskList); err != nil {
			return
		}
		if len(taskList) > 0 {
			task = taskList[0]
			return
		}
	}
	if err = dao.X.SQL("select * from task where request = ? and status <> 'done' and del_flag = 0 order by sort asc", requestId).Find(&taskList); err != nil {
		return
	}
	if len(taskList) > 0 {
		if len(taskList) == 1 {
			task = taskList[0]
			return
		}
		for _, taskTable := range taskList {
			if taskTable.TaskTemplate != "" {
				doingTaskMap[taskTable.TaskTemplate] = taskTable
			}
		}
		err = dao.X.SQL("select * from task_template where request_template = ? order by sort asc", templateId).Find(&taskTemplateList)
		if err != nil {
			return
		}
		if len(taskTemplateList) > 0 {
			for _, taskTemplate := range taskTemplateList {
				if v, ok := doingTaskMap[taskTemplate.Id]; ok {
					task = v
					return
				}
			}
		}
	}
	return
}

// GetTaskMapByRequestId 取最后一次任务提交之前的数据
func (s *TaskService) GetTaskMapByRequestId(request models.RequestTable) (taskMap map[string]*models.TaskTable, err error) {
	var taskList []*models.TaskTable
	taskMap = make(map[string]*models.TaskTable)
	if request.Status == string(models.RequestStatusDraft) {
		return
	}
	err = dao.X.SQL("select * from task where request = ? and del_flag = 0 order by sort desc", request.Id).Find(&taskList)
	if err != nil {
		return
	}
	if len(taskList) > 0 {
		for _, task := range taskList {
			if task.Type == string(models.TaskTypeSubmit) {
				if _, ok := taskMap[task.TaskTemplate]; !ok {
					taskMap[task.TaskTemplate] = task
				}
				break
			}
			if task.TaskTemplate != "" {
				if _, ok := taskMap[task.TaskTemplate]; !ok {
					taskMap[task.TaskTemplate] = task
				}
			}
		}
	}
	return
}

// GetDoneTaskByRequestId 已完成对任务,取最后一次提交请求后的完成任务
func (s *TaskService) GetDoneTaskByRequestId(request models.RequestTable) (taskList []*models.TaskTable, err error) {
	if request.Status == string(models.RequestStatusDraft) {
		return
	}
	var allTaskList []*models.TaskTable
	taskList = []*models.TaskTable{}
	err = dao.X.SQL("select * from task where request = ? and del_flag = 0 order by sort desc", request.Id).Find(&allTaskList)
	if err != nil {
		return
	}
	if len(allTaskList) > 0 {
		for _, task := range allTaskList {
			if task.Type == string(models.TaskTypeSubmit) {
				taskList = append(taskList, task)
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

// GetLatestTaskListByRequestIdAndTaskTemplateId 获取任务模版最新的执行任务,取最后一次提交请求后
func (s *TaskService) GetLatestTaskListByRequestIdAndTaskTemplateId(requestId, taskTemplateId string) (taskList []*models.TaskTable, err error) {
	var allTaskList []*models.TaskTable
	taskList = []*models.TaskTable{}
	err = dao.X.SQL("select * from task where request = ? and del_flag = 0 order by sort desc", requestId).Find(&allTaskList)
	if err != nil {
		return
	}
	if len(allTaskList) > 0 {
		for _, task := range allTaskList {
			if task.Type == string(models.TaskTypeSubmit) {
				break
			}
			if task.TaskTemplate == taskTemplateId {
				taskList = append(taskList, task)
			}
		}
	}
	return
}
