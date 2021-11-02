package db

import (
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"time"
)

func GetTaskTemplate(requestTemplateId, proNodeId string) (result models.TaskTemplateDto, err error) {
	result = models.TaskTemplateDto{}
	taskTemplate, getTaskErr := getSimpleTaskTemplate("", requestTemplateId, proNodeId)
	if getTaskErr != nil {
		log.Logger.Warn("GetTaskTemplate warning", log.Error(getTaskErr))
		return
	}
	result.Id = taskTemplate.Id
	result.NodeDefId = taskTemplate.NodeDefId
	result.NodeDefName = taskTemplate.NodeName
	result.ExpireDay = taskTemplate.ExpireDay
	result.Handler = taskTemplate.Handler
	result.MGMTRoleObjs = []*models.RoleTable{}
	result.MGMTRoles = []string{}
	result.USERoleObjs = []*models.RoleTable{}
	result.USERoles = []string{}
	var formTemplateTable []*models.FormTemplateTable
	err = x.SQL("select * from form_template where id=?", taskTemplate.FormTemplate).Find(&formTemplateTable)
	if err != nil {
		err = fmt.Errorf("Try to query form template table fail,%s ", err.Error())
		return
	}
	if len(formTemplateTable) == 0 {
		err = fmt.Errorf("Can not find any form template with id=%s ", taskTemplate.FormTemplate)
		return
	}
	result.Name = formTemplateTable[0].Name
	result.Description = formTemplateTable[0].Description
	result.UpdatedTime = formTemplateTable[0].UpdatedTime
	result.UpdatedBy = formTemplateTable[0].UpdatedBy
	var formItemTemplate []*models.FormItemTemplateTable
	x.SQL("select * from form_item_template where form_template=?", taskTemplate.FormTemplate).Find(&formItemTemplate)
	result.Items = formItemTemplate
	roleMap, _ := getRoleMap()
	var taskRoleTable []*models.TaskTemplateRoleTable
	x.SQL("select `role`,role_type from task_template_role where task_template=?", taskTemplate.Id).Find(&taskRoleTable)
	for _, role := range taskRoleTable {
		if role.RoleType == "MGMT" {
			result.MGMTRoleObjs = append(result.MGMTRoleObjs, roleMap[role.Role])
			result.MGMTRoles = append(result.MGMTRoles, role.Role)
		} else {
			result.USERoleObjs = append(result.USERoleObjs, roleMap[role.Role])
			result.USERoles = append(result.USERoles, role.Role)
		}
	}
	return
}

func CreateTaskTemplate(param models.TaskTemplateDto, requestTemplateId string) error {
	_, checkExistErr := getSimpleTaskTemplate("", requestTemplateId, param.NodeDefId)
	if checkExistErr == nil {
		return fmt.Errorf("RequestTemplate:%s nodeDefId:%s already have task form,please reload ", requestTemplateId, param.NodeDefId)
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	formCreateParam := models.FormTemplateDto{Name: param.Name, Description: param.Description, UpdatedBy: param.UpdatedBy, Items: param.Items, NowTime: nowTime}
	actions, formId := getFormTemplateCreateActions(formCreateParam)
	taskTemplateId := guid.CreateGuid()
	actions = append(actions, &execAction{Sql: "insert into task_template(id,name,description,form_template,request_template,node_id,node_def_id,node_name,expire_day,handler) value (?,?,?,?,?,?,?,?,?,?)", Param: []interface{}{taskTemplateId, param.Name, param.Description, formId, requestTemplateId, param.NodeId, param.NodeDefId, param.NodeDefName, param.ExpireDay, param.Handler}})
	for _, v := range param.MGMTRoles {
		actions = append(actions, &execAction{Sql: "insert into task_template_role(id,task_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{taskTemplateId + models.SysTableIdConnector + v + models.SysTableIdConnector + "MGMT", taskTemplateId, v, "MGMT"}})
	}
	for _, v := range param.USERoles {
		actions = append(actions, &execAction{Sql: "insert into task_template_role(id,task_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{taskTemplateId + models.SysTableIdConnector + v + models.SysTableIdConnector + "USE", taskTemplateId, v, "USE"}})
	}
	return transactionWithoutForeignCheck(actions)
}

func UpdateTaskTemplate(param models.TaskTemplateDto) error {
	taskTemplate, err := getSimpleTaskTemplate(param.Id, "", "")
	if err != nil {
		return err
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	formUpdateParam := models.FormTemplateDto{Id: taskTemplate.FormTemplate, Name: param.Name, Description: param.Description, UpdatedBy: param.UpdatedBy, UpdatedTime: param.UpdatedTime, Items: param.Items, NowTime: nowTime}
	actions, getActionErr := getFormTemplateUpdateActions(formUpdateParam)
	if getActionErr != nil {
		return getActionErr
	}
	actions = append(actions, &execAction{Sql: "delete from task_template_role where task_template=?", Param: []interface{}{param.Id}})
	for _, v := range param.MGMTRoles {
		actions = append(actions, &execAction{Sql: "insert into task_template_role(id,task_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{param.Id + models.SysTableIdConnector + v + models.SysTableIdConnector + "MGMT", param.Id, v, "MGMT"}})
	}
	for _, v := range param.USERoles {
		actions = append(actions, &execAction{Sql: "insert into task_template_role(id,task_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{param.Id + models.SysTableIdConnector + v + models.SysTableIdConnector + "USE", param.Id, v, "USE"}})
	}
	actions = append(actions, &execAction{Sql: "update task_template set name=?,description=?,expire_day=?,handler=? where id=?", Param: []interface{}{param.Name, param.Description, param.ExpireDay, param.Handler, param.Id}})
	return transaction(actions)
}

func getSimpleTaskTemplate(id, requestTemplate, proNodeId string) (result models.TaskTemplateTable, err error) {
	var taskTemplateTable []*models.TaskTemplateTable
	baseSql := "select * from task_template where 1=1 "
	var params []interface{}
	if id != "" {
		baseSql += " and id=?"
		params = append(params, id)
	}
	if requestTemplate != "" {
		baseSql += " and request_template=?"
		params = append(params, requestTemplate)
	}
	if proNodeId != "" {
		baseSql += " and node_def_id=?"
		params = append(params, proNodeId)
	}
	err = x.SQL(baseSql, params...).Find(&taskTemplateTable)
	if err != nil {
		err = fmt.Errorf("Try to query database fail,%s ", err.Error())
		return
	}
	if len(taskTemplateTable) == 0 {
		err = fmt.Errorf("Can not find task template with id:%s ", id)
		result = models.TaskTemplateTable{}
		return
	}
	result = *taskTemplateTable[0]
	return
}

func getRoleMap() (result map[string]*models.RoleTable, err error) {
	result = make(map[string]*models.RoleTable)
	var roleTable []*models.RoleTable
	err = x.SQL("select id,display_name from `role`").Find(&roleTable)
	if err != nil {
		return
	}
	for _, v := range roleTable {
		result[v.Id] = v
	}
	return
}
