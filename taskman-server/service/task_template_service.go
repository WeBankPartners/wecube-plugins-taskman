package service

import (
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"time"
)

type TaskTemplateService struct {
	taskTemplateDao     dao.TaskTemplateDao
	taskTemplateRoleDao dao.TaskTemplateRoleDao
}

func GetTaskTemplate(requestTemplateId, proNodeId, nodeId string) (result models.TaskTemplateDto, err error) {
	result = models.TaskTemplateDto{}
	taskTemplate, getTaskErr := getSimpleTaskTemplate("", requestTemplateId, proNodeId, nodeId)
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
	result.RequestTemplateId = requestTemplateId
	var formTemplateTable []*models.FormTemplateTable
	err = dao.X.SQL("select * from form_template where id=?", taskTemplate.FormTemplate).Find(&formTemplateTable)
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
	dao.X.SQL("select * from form_item_template where form_template=?", taskTemplate.FormTemplate).Find(&formItemTemplate)
	result.Items = formItemTemplate
	roleMap, _ := getRoleMap()
	var taskRoleTable []*models.TaskTemplateRoleTable
	dao.X.SQL("select `role`,role_type from task_template_role where task_template=?", taskTemplate.Id).Find(&taskRoleTable)
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

func getFormTemplateCreateActions(param models.FormTemplateDto) (actions []*dao.ExecAction, id string) {
	param.Id = guid.CreateGuid()
	id = param.Id
	itemIds := guid.CreateGuidList(len(param.Items))
	insertAction := dao.ExecAction{Sql: "insert into form_template(id,name,description,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?)"}
	insertAction.Param = []interface{}{param.Id, param.Name, param.Description, param.UpdatedBy, param.NowTime, param.UpdatedBy, param.NowTime}
	actions = append(actions, &insertAction)
	for i, item := range param.Items {
		tmpAction := dao.ExecAction{Sql: "insert into form_item_template(id,form_template,name,description,item_group,item_group_name,default_value,sort,package_name,entity,attr_def_id,attr_def_name,attr_def_data_type,element_type,title,width,ref_package_name,ref_entity,data_options,required,regular,is_edit,is_view,is_output,in_display_name,is_ref_inside,multiple,default_clear) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
		tmpAction.Param = []interface{}{itemIds[i], id, item.Name, item.Description, item.ItemGroup, item.ItemGroupName, item.DefaultValue, item.Sort, item.PackageName, item.Entity, item.AttrDefId, item.AttrDefName, item.AttrDefDataType, item.ElementType, item.Title, item.Width, item.RefPackageName, item.RefEntity, item.DataOptions, item.Required, item.Regular, item.IsEdit, item.IsView, item.IsOutput, item.InDisplayName, item.IsRefInside, item.Multiple, item.DefaultClear}
		actions = append(actions, &tmpAction)
	}
	return
}

func CreateTaskTemplate(param models.TaskTemplateDto, requestTemplateId string) error {
	_, checkExistErr := getSimpleTaskTemplate("", requestTemplateId, param.NodeDefId, "")
	if checkExistErr == nil {
		return fmt.Errorf("RequestTemplate:%s nodeDefId:%s already have task form,please reload ", requestTemplateId, param.NodeDefId)
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	formCreateParam := models.FormTemplateDto{Name: param.Name, Description: param.Description, UpdatedBy: param.UpdatedBy, Items: param.Items, NowTime: nowTime}
	actions, formId := getFormTemplateCreateActions(formCreateParam)
	taskTemplateId := guid.CreateGuid()
	actions = append(actions, &dao.ExecAction{Sql: "insert into task_template(id,name,description,form_template,request_template,node_id,node_def_id,node_name,expire_day,handler) value (?,?,?,?,?,?,?,?,?,?)", Param: []interface{}{taskTemplateId, param.Name, param.Description, formId, requestTemplateId, param.NodeId, param.NodeDefId, param.NodeDefName, param.ExpireDay, param.Handler}})
	for _, v := range param.MGMTRoles {
		actions = append(actions, &dao.ExecAction{Sql: "insert into task_template_role(id,task_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{taskTemplateId + models.SysTableIdConnector + v + models.SysTableIdConnector + "MGMT", taskTemplateId, v, "MGMT"}})
	}
	for _, v := range param.USERoles {
		actions = append(actions, &dao.ExecAction{Sql: "insert into task_template_role(id,task_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{taskTemplateId + models.SysTableIdConnector + v + models.SysTableIdConnector + "USE", taskTemplateId, v, "USE"}})
	}
	return dao.TransactionWithoutForeignCheck(actions)
}

func UpdateTaskTemplate(param models.TaskTemplateDto) error {
	taskTemplate, err := getSimpleTaskTemplate(param.Id, "", "", "")
	if err != nil {
		return err
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	formUpdateParam := models.FormTemplateDto{Id: taskTemplate.FormTemplate, Name: param.Name, Description: param.Description, UpdatedBy: param.UpdatedBy, UpdatedTime: param.UpdatedTime, Items: param.Items, NowTime: nowTime}
	actions, getActionErr := getFormTemplateUpdateActions(formUpdateParam)
	if getActionErr != nil {
		return getActionErr
	}
	actions = append(actions, &dao.ExecAction{Sql: "delete from task_template_role where task_template=?", Param: []interface{}{param.Id}})
	for _, v := range param.MGMTRoles {
		actions = append(actions, &dao.ExecAction{Sql: "insert into task_template_role(id,task_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{param.Id + models.SysTableIdConnector + v + models.SysTableIdConnector + "MGMT", param.Id, v, "MGMT"}})
	}
	for _, v := range param.USERoles {
		actions = append(actions, &dao.ExecAction{Sql: "insert into task_template_role(id,task_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{param.Id + models.SysTableIdConnector + v + models.SysTableIdConnector + "USE", param.Id, v, "USE"}})
	}
	actions = append(actions, &dao.ExecAction{Sql: "update task_template set name=?,description=?,expire_day=?,handler=? where id=?", Param: []interface{}{param.Name, param.Description, param.ExpireDay, param.Handler, param.Id}})
	return dao.Transaction(actions)
}

func getFormTemplateUpdateActions(param models.FormTemplateDto) (actions []*dao.ExecAction, err error) {
	var formTemplateTable []*models.FormTemplateTable
	err = dao.X.SQL("select id,updated_time from form_template where id=?", param.Id).Find(&formTemplateTable)
	if err != nil {
		err = fmt.Errorf("Try to query form template table fail,%s ", err.Error())
		return
	}
	if len(formTemplateTable) == 0 {
		err = fmt.Errorf("Can not find any form template with id=%s ", param.Id)
		return
	}
	if param.UpdatedTime != formTemplateTable[0].UpdatedTime {
		err = fmt.Errorf("Update time validate fail,please refersh data ")
		return
	}
	updateAction := dao.ExecAction{Sql: "update form_template set name=?,description=?,updated_by=?,updated_time=? where id=?"}
	updateAction.Param = []interface{}{param.Name, param.Description, param.UpdatedBy, param.NowTime, param.Id}
	actions = append(actions, &updateAction)
	newItemGuidList := guid.CreateGuidList(len(param.Items))
	var formItemTemplate []*models.FormItemTemplateTable
	dao.X.SQL("select id from form_item_template where form_template=?", param.Id).Find(&formItemTemplate)
	for i, inputItem := range param.Items {
		tmpAction := dao.ExecAction{}
		if inputItem.Id == "" {
			inputItem.Id = newItemGuidList[i]
			tmpAction.Sql = "insert into form_item_template(id,form_template,name,description,item_group,item_group_name,default_value,sort,package_name,entity,attr_def_id,attr_def_name,attr_def_data_type,element_type,title,width,ref_package_name,ref_entity,data_options,required,regular,is_edit,is_view,is_output,in_display_name,is_ref_inside,multiple,default_clear) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
			tmpAction.Param = []interface{}{inputItem.Id, param.Id, inputItem.Name, inputItem.Description, inputItem.ItemGroup, inputItem.ItemGroupName, inputItem.DefaultValue, inputItem.Sort, inputItem.PackageName, inputItem.Entity, inputItem.AttrDefId, inputItem.AttrDefName, inputItem.AttrDefDataType, inputItem.ElementType, inputItem.Title, inputItem.Width, inputItem.RefPackageName, inputItem.RefEntity, inputItem.DataOptions, inputItem.Required, inputItem.Regular, inputItem.IsEdit, inputItem.IsView, inputItem.IsOutput, inputItem.InDisplayName, inputItem.IsRefInside, inputItem.Multiple, inputItem.DefaultClear}
		} else {
			tmpAction.Sql = "update form_item_template set name=?,description=?,item_group=?,item_group_name=?,default_value=?,sort=?,package_name=?,entity=?,attr_def_id=?,attr_def_name=?,attr_def_data_type=?,element_type=?,title=?,width=?,ref_package_name=?,ref_entity=?,data_options=?,required=?,regular=?,is_edit=?,is_view=?,is_output=?,in_display_name=?,is_ref_inside=?,multiple=?,default_clear=? where id=?"
			tmpAction.Param = []interface{}{inputItem.Name, inputItem.Description, inputItem.ItemGroup, inputItem.ItemGroupName, inputItem.DefaultValue, inputItem.Sort, inputItem.PackageName, inputItem.Entity, inputItem.AttrDefId, inputItem.AttrDefName, inputItem.AttrDefDataType, inputItem.ElementType, inputItem.Title, inputItem.Width, inputItem.RefPackageName, inputItem.RefEntity, inputItem.DataOptions, inputItem.Required, inputItem.Regular, inputItem.IsEdit, inputItem.IsView, inputItem.IsOutput, inputItem.InDisplayName, inputItem.IsRefInside, inputItem.Multiple, inputItem.DefaultClear, inputItem.Id}
		}
		actions = append(actions, &tmpAction)
	}
	for _, existItem := range formItemTemplate {
		existFlag := false
		for _, inputItem := range param.Items {
			if existItem.Id == inputItem.Id {
				existFlag = true
				break
			}
		}
		if !existFlag {
			actions = append(actions, &dao.ExecAction{Sql: "delete from form_item where form_item_template=?", Param: []interface{}{existItem.Id}})
			actions = append(actions, &dao.ExecAction{Sql: "delete from form_item_template where id=?", Param: []interface{}{existItem.Id}})
		}
	}
	return
}

func getSimpleTaskTemplate(id, requestTemplate, proNodeId, nodeId string) (result models.TaskTemplateTable, err error) {
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
	if nodeId != "" {
		baseSql += " and node_id=?"
		params = append(params, nodeId)
	}
	err = dao.X.SQL(baseSql, params...).Find(&taskTemplateTable)
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
	err = dao.X.SQL("select id,display_name from `role`").Find(&roleTable)
	if err != nil {
		return
	}
	for _, v := range roleTable {
		result[v.Id] = v
	}
	return
}

func getTaskTemplateHandler(requestTemplate string) (taskTemplateMap map[string]*models.TaskTemplateVo, err error) {
	taskTemplateMap = make(map[string]*models.TaskTemplateVo)
	var taskTemplateList []*models.TaskTemplateVo
	err = dao.X.SQL("select  t.id,t.handler,tr.role from task_template t join task_template_role tr on  "+
		"t.id=tr.task_template where t.request_template = ?", requestTemplate).Find(&taskTemplateList)
	if len(taskTemplateList) > 0 {
		for _, taskTemplate := range taskTemplateList {
			taskTemplateMap[taskTemplate.Id] = taskTemplate
		}
	}
	return
}

func GetTaskTemplateMapByRequestTemplate(requestTemplate string) (taskTemplateMap map[string]int, err error) {
	taskTemplateMap = make(map[string]int)
	var rowsData []*models.TaskTemplateTable
	sql := "select * from task_template where request_template = ?"
	err = dao.X.SQL(sql, requestTemplate).Find(&rowsData)
	if len(rowsData) > 0 {
		for _, row := range rowsData {
			taskTemplateMap[row.Name] = row.ExpireDay
		}
	}
	return
}