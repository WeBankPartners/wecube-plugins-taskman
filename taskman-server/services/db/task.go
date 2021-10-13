package db

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"time"
)

func GetTaskFormStruct(procInstId, nodeDefId string) (result models.TaskMetaResult, err error) {
	result = models.TaskMetaResult{Status: "OK", Message: "Success"}
	var items []*models.FormItemTemplateTable
	err = x.SQL("select * from form_item_template where form_template in (select form_template from task_template where node_def_id=? and request_template in (select request_template from request where proc_instance_id=?))", nodeDefId, procInstId).Find(&items)
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

func PluginTaskCreate(input *models.PluginTaskCreateRequestObj) (result *models.PluginTaskCreateOutputObj, err error) {
	result = &models.PluginTaskCreateOutputObj{CallbackParameter: input.CallbackParameter, ErrorCode: "0", ErrorMessage: "", Comment: "OK"}
	var requestTable []*models.RequestTable
	err = x.SQL("select id,form,request_template,emergency from request where proc_instance_id=?", input.ProcInstId).Find(&requestTable)
	if err != nil {
		return result, fmt.Errorf("Try to check proc_instance_id:%s is in request fail,%s ", input.ProcInstId, err.Error())
	}
	var taskFormInput models.PluginTaskFormDto
	err = json.Unmarshal([]byte(input.TaskFormInput), &taskFormInput)
	if err != nil {
		return result, fmt.Errorf("Try to json unmarshal taskFormInput to json data fail,%s ", err.Error())
	}
	var actions []*execAction
	itemTemplateMap := getFormItemTemplateMap(taskFormInput.FormMateId)
	newTaskFormObj := models.FormTable{Id: guid.CreateGuid(), Name: "form_" + input.TaskName}
	newTaskObj := models.TaskTable{Id: guid.CreateGuid(), Name: input.TaskName, Status: "created", Form: newTaskFormObj.Id, Reporter: input.Reporter, ReportRole: input.RoleName, Description: input.TaskDescription, CallbackUrl: input.CallbackUrl}
	newTaskObj.ProcDefId = taskFormInput.ProcDefId
	newTaskObj.ProcDefKey = taskFormInput.ProcDefKey
	if len(requestTable) > 0 {
		newTaskObj.Request = requestTable[0].Id
		newTaskObj.NodeDefId = taskFormInput.TaskNodeDefId
		newTaskObj.Emergency = requestTable[0].Emergency
		var taskTemplateTable []*models.TaskTemplateTable
		x.SQL("select * from task_template where request_template=? and node_def_id=?", requestTable[0].RequestTemplate, taskFormInput.TaskNodeDefId).Find(&taskTemplateTable)
		if len(taskTemplateTable) > 0 {
			newTaskObj.TaskTemplate = taskTemplateTable[0].Id
			newTaskObj.NodeName = taskTemplateTable[0].NodeName
			newTaskFormObj.FormTemplate = taskTemplateTable[0].FormTemplate
		} else {
			log.Logger.Warn("Can not find any taskTemplate", log.String("requestTemplate", requestTable[0].RequestTemplate), log.String("nodeDefId", taskFormInput.TaskNodeDefId))
		}
	}
	tmpTaskRequest, tmpTaskTemplate := "NULL", "NULL"
	if newTaskObj.Request == "" {
		tmpTaskRequest = fmt.Sprintf("'%s'", newTaskObj.Request)
	}
	if newTaskObj.TaskTemplate == "" {
		tmpTaskTemplate = fmt.Sprintf("'%s'", newTaskObj.TaskTemplate)
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	taskInsertAction := execAction{Sql: "insert into task(id,name,description,form,status,request,task_template,proc_def_id,proc_def_key,node_def_id,node_name,callback_url,reporter,report_role,report_time,emergency) value (?,?,?,?,?," + tmpTaskRequest + "," + tmpTaskTemplate + ",?,?,?,?,?,?,?,?,?)"}
	taskInsertAction.Param = []interface{}{newTaskObj.Id, newTaskObj.Name, newTaskObj.Description, newTaskObj.Form, newTaskObj.Status, newTaskObj.ProcDefId, newTaskObj.ProcDefKey, newTaskObj.NodeDefId, newTaskObj.NodeName, newTaskObj.CallbackUrl, newTaskObj.Reporter, newTaskObj.ReportRole, nowTime, newTaskObj.Emergency}
	actions = append(actions, &taskInsertAction)
	if newTaskFormObj.FormTemplate == "" {
		actions = append(actions, &execAction{Sql: "insert into form(id,name) value (?,?)", Param: []interface{}{newTaskFormObj.Id, newTaskFormObj.Name}})
	} else {
		actions = append(actions, &execAction{Sql: "insert into form(id,name,form_template) value (?,?,?)", Param: []interface{}{newTaskFormObj.Id, newTaskFormObj.Name, newTaskFormObj.FormTemplate}})
	}
	for _, formDataEntity := range taskFormInput.FormDataEntities {
		tmpGuidList := guid.CreateGuidList(len(formDataEntity.FormItemValues))
		for i, formDataItem := range formDataEntity.FormItemValues {
			if formDataItem.FormItemMateId != "" {
				tmpItemGroup := itemTemplateMap[formDataItem.FormItemMateId].ItemGroup
				actions = append(actions, &execAction{Sql: "insert into form_item(id,form,form_item_template,name,value,item_group,row_data_id) value (?,?,?,?,?,?,?)", Param: []interface{}{tmpGuidList[i], newTaskFormObj.Id, formDataItem.FormItemMateId, formDataItem.AttrName, formDataItem.AttrValue, tmpItemGroup, formDataItem.EntityDataId}})
			} else {
				actions = append(actions, &execAction{Sql: "insert into form_item(id,form,name,value,row_data_id) value (?,?,?,?,?)", Param: []interface{}{tmpGuidList[i], newTaskFormObj.Id, formDataItem.AttrName, formDataItem.AttrValue, formDataItem.EntityDataId}})
			}
		}
	}
	err = transactionWithoutForeignCheck(actions)
	return
}

func getFormItemTemplateMap(formTemplateId string) map[string]*models.FormItemTemplateTable {
	resultMap := make(map[string]*models.FormItemTemplateTable)
	var itemTemplateTable []*models.FormItemTemplateTable
	x.SQL("select * from form_item_template where form_template=?", formTemplateId).Find(&itemTemplateTable)
	for _, v := range itemTemplateTable {
		resultMap[v.Id] = v
	}
	return resultMap
}
