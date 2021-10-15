package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"io/ioutil"
	"net/http"
	"strings"
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

func PluginTaskCreate(input *models.PluginTaskCreateRequestObj, callRequestId string) (result *models.PluginTaskCreateOutputObj, err error) {
	result = &models.PluginTaskCreateOutputObj{CallbackParameter: input.CallbackParameter, ErrorCode: "0", ErrorMessage: "", Comment: ""}
	var requestTable []*models.RequestTable
	err = x.SQL("select id,form,request_template,emergency from request where proc_instance_id=?", input.ProcInstId).Find(&requestTable)
	if err != nil {
		return result, fmt.Errorf("Try to check proc_instance_id:%s is in request fail,%s ", input.ProcInstId, err.Error())
	}
	var actions []*execAction
	nowTime := time.Now().Format(models.DateTimeFormat)
	newTaskFormObj := models.FormTable{Id: guid.CreateGuid(), Name: "form_" + input.TaskName}
	input.RoleName = remakeTaskReportRole(input.RoleName)
	newTaskObj := models.TaskTable{Id: guid.CreateGuid(), Name: input.TaskName, Status: "created", Form: newTaskFormObj.Id, Reporter: input.Reporter, ReportRole: input.RoleName, Description: input.TaskDescription, CallbackUrl: input.CallbackUrl}
	var taskFormInput models.PluginTaskFormDto
	if input.TaskFormInput != "" {
		err = json.Unmarshal([]byte(input.TaskFormInput), &taskFormInput)
		if err != nil {
			return result, fmt.Errorf("Try to json unmarshal taskFormInput to json data fail,%s ", err.Error())
		}
	} else {
		taskInsertAction := execAction{Sql: "insert into task(id,name,description,status,proc_def_id,proc_def_key,node_def_id,node_name,callback_url,reporter,report_role,report_time,emergency,callback_request_id,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
		taskInsertAction.Param = []interface{}{newTaskObj.Id, newTaskObj.Name, newTaskObj.Description, newTaskObj.Status, newTaskObj.ProcDefId, newTaskObj.ProcDefKey, newTaskObj.NodeDefId, newTaskObj.NodeName, newTaskObj.CallbackUrl, newTaskObj.Reporter, newTaskObj.ReportRole, nowTime, newTaskObj.Emergency, callRequestId, "system", nowTime, "system", nowTime}
		actions = append(actions, &taskInsertAction)
		err = transaction(actions)
		return
	}
	itemTemplateMap := getFormItemTemplateMap(taskFormInput.FormMetaId)
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
	taskInsertAction := execAction{Sql: "insert into task(id,name,description,form,status,request,task_template,proc_def_id,proc_def_key,node_def_id,node_name,callback_url,reporter,report_role,report_time,emergency,cache,callback_request_id,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
	taskInsertAction.Param = []interface{}{newTaskObj.Id, newTaskObj.Name, newTaskObj.Description, newTaskObj.Form, newTaskObj.Status, newTaskObj.Request, newTaskObj.TaskTemplate, newTaskObj.ProcDefId, newTaskObj.ProcDefKey, newTaskObj.NodeDefId, newTaskObj.NodeName, newTaskObj.CallbackUrl, newTaskObj.Reporter, newTaskObj.ReportRole, nowTime, newTaskObj.Emergency, input.TaskFormInput, callRequestId, "system", nowTime, "system", nowTime}
	actions = append(actions, &taskInsertAction)
	actions = append(actions, &execAction{Sql: "insert into form(id,name,form_template) value (?,?,?)", Param: []interface{}{newTaskFormObj.Id, newTaskFormObj.Name, newTaskFormObj.FormTemplate}})
	for _, formDataEntity := range taskFormInput.FormDataEntities {
		tmpGuidList := guid.CreateGuidList(len(formDataEntity.FormItemValues))
		for i, formDataItem := range formDataEntity.FormItemValues {
			tmpItemGroup := itemTemplateMap[formDataItem.FormItemMetaId].ItemGroup
			actions = append(actions, &execAction{Sql: "insert into form_item(id,form,form_item_template,name,value,item_group,row_data_id) value (?,?,?,?,?,?,?)", Param: []interface{}{tmpGuidList[i], newTaskFormObj.Id, formDataItem.FormItemMetaId, formDataItem.AttrName, formDataItem.AttrValue, tmpItemGroup, formDataItem.Oid}})
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

func ListTask(param *models.QueryRequestParam, userRoles []string) (pageInfo models.PageInfo, rowData []*models.TaskTable, err error) {
	rowData = []*models.TaskTable{}
	roleFilterSql := "1=1"
	if len(userRoles) > 0 {
		roleFilterList := []string{}
		for _, v := range userRoles {
			roleFilterList = append(roleFilterList, "report_role like '%,"+v+",%'")
		}
		roleFilterSql = strings.Join(roleFilterList, " or ")
	}
	filterSql, _, queryParam := transFiltersToSQL(param, &models.TransFiltersParam{IsStruct: true, StructObj: models.TaskTable{}, PrimaryKey: "id", Prefix: "t1"})
	baseSql := fmt.Sprintf("select t1.id,t1.name,t1.description,t1.form,t1.status,t1.`version`,t1.request,t1.task_template,t1.node_name,t1.reporter,t1.report_time,t1.emergency from (select * from task where task_template in (select task_template from task_template_role where role_type='USE' and `role` in ('"+strings.Join(userRoles, "','")+"')) union select * from task where task_template is null and (%s)) t1 where t1.del_flag=0 %s ", roleFilterSql, filterSql)
	if param.Paging {
		pageInfo.StartIndex = param.Pageable.StartIndex
		pageInfo.PageSize = param.Pageable.PageSize
		pageInfo.TotalRows = queryCount(baseSql, queryParam...)
		pageSql, pageParam := transPageInfoToSQL(*param.Pageable)
		baseSql += pageSql
		queryParam = append(queryParam, pageParam...)
	}
	err = x.SQL(baseSql, queryParam...).Find(&rowData)
	return
}

func GetTask() {

}

func getSimpleTask(taskId string) (result models.TaskTable, err error) {
	var taskTable []*models.TaskTable
	err = x.SQL("select * from task where id=?", taskId).Find(&taskTable)
	if err != nil {
		return
	}
	if len(taskTable) == 0 {
		return result, fmt.Errorf("Can not find any task with id:%s ", taskId)
	}
	result = *taskTable[0]
	return
}

func ApproveTask(taskId, operator, userToken string) error {
	requestParam, callbackUrl, err := getApproveCallbackParam(taskId)
	if err != nil {
		return err
	}
	requestBytes, _ := json.Marshal(requestParam)
	req, newReqErr := http.NewRequest(http.MethodPost, models.Config.Wecube.BaseUrl+callbackUrl, bytes.NewReader(requestBytes))
	if newReqErr != nil {
		return fmt.Errorf("Try to new http request fail,%s ", newReqErr.Error())
	}
	req.Header.Set("Authorization", userToken)
	req.Header.Set("Content-Type", "application/json")
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		return fmt.Errorf("Try to do http request fail,%s ", respErr.Error())
	}
	var respResult models.CallbackResult
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	log.Logger.Info("Callback response", log.String("body", string(b)))
	err = json.Unmarshal(b, &respResult)
	if err != nil {
		return fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
	}
	if respResult.Status != "OK" {
		return fmt.Errorf("Callback fail,%s ", respResult.Message)
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err = x.Exec("update task set callback_parameter=?,status=?,updated_by=?,updated_time=? where id=?", string(requestBytes), "done", operator, nowTime, taskId)
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
		resultObj.Outputs = []*models.PluginTaskCreateOutputObj{{Comment: taskObj.Result, ErrorCode: "0"}}
		result.Results = resultObj
		return
	}
	var taskFormOutput models.PluginTaskFormDto
	err = json.Unmarshal([]byte(taskObj.Cache), &taskFormOutput)
	if err != nil {
		return result, callbackUrl, fmt.Errorf("Try to json unmarshal cache data fail:%s ", err.Error())
	}
	var items []*models.TaskFormItemQueryObj
	err = x.SQL("select t1.*,t2.attr_def_data_type,t2.element_type from form_item t1 left join form_item_template t2 on t1.form_item_template=t2.id where form=?", taskObj.Form).Find(&items)
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
	resultObj.Outputs = []*models.PluginTaskCreateOutputObj{{Comment: taskObj.Result, ErrorCode: "0", TaskFormOutput: string(formBytes)}}
	result.Results = resultObj
	return
}
