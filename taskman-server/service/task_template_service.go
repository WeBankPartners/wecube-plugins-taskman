package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
)

type TaskTemplateService struct {
	taskTemplateDao     *dao.TaskTemplateDao
	taskTemplateRoleDao *dao.TaskTemplateRoleDao
}

func GetTaskTemplate(requestTemplateId, proNodeId, nodeId string) (result models.TaskTemplateDto, err error) {
	result = models.TaskTemplateDto{Items: []*models.FormItemTemplateDto{}}
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
	var formItemTemplate []*models.FormItemTemplateTable
	var formItemTemplateGroup models.FormTemplateTable
	dao.X.SQL("select * from form_item_template where form_template=?", taskTemplate.FormTemplate).Find(&formItemTemplate)
	if len(formItemTemplate) > 0 && formItemTemplate[0].FormTemplate != "" {
		dao.X.SQL("select * from form_item_template_group where id=?", formItemTemplate[0].FormTemplate).Get(&formItemTemplateGroup)
		for _, formItem := range formItemTemplate {
			result.Items = append(result.Items, models.ConvertFormItemTemplateModel2Dto(formItem, formItemTemplateGroup))
		}
	}
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

func (s TaskTemplateService) CreateCustomTaskTemplate(param *models.CustomTaskTemplateCreateParam) (*models.CustomTaskTemplateCreateResponse, error) {
	actions := []*dao.ExecAction{}
	// 查询现有任务模板列表
	var taskTemplates []*models.TaskTemplateTable
	err := s.taskTemplateDao.DB.SQL("SELECT * FROM task_template WHERE request_template = ? ORDER BY sort", param.RequestTemplate).Find(&taskTemplates)
	if err != nil {
		return nil, err
	}
	// 校验参数
	if param.Sort <= 0 || param.Sort > len(taskTemplates)+1 {
		return nil, errors.New("param sort out of range")
	}
	// 插入新任务模板
	nowTime := time.Now().Format(models.DateTimeFormat)
	newTaskTemplate := &models.TaskTemplateTable{
		Id:              guid.CreateGuid(),
		Type:            string(models.TaskTemplateTypeCustom),
		Sort:            param.Sort,
		RequestTemplate: param.RequestTemplate,
		Name:            param.Name,
		RoleType:        string(models.TaskTemplateRoleTypeCustom),
		CreatedTime:     nowTime,
		UpdatedTime:     nowTime,
	}
	action := &dao.ExecAction{Sql: "INSERT INTO task_template (id,type,sort,request_template,name,expire_day,role_type,created_time,updated_time) VALUES (?,?,?,?,?,?,?,?,?)"}
	action.Param = []interface{}{newTaskTemplate.Id, newTaskTemplate.Type, newTaskTemplate.Sort, newTaskTemplate.RequestTemplate, newTaskTemplate.Name, newTaskTemplate.ExpireDay, newTaskTemplate.RoleType, newTaskTemplate.CreatedTime, newTaskTemplate.UpdatedTime}
	actions = append(actions, action)
	// 插入新任务处理模板
	newTaskTemplateRole := &models.TaskTemplateRoleTable{
		Id:             guid.CreateGuid(),
		TaskTemplate:   newTaskTemplate.Id,
		CustomRoleType: string(models.TaskTemplateRoleRoleTypeTemplate),
		HandlerType:    string(models.TaskTemplateRoleHandlerTypeTemplate),
	}
	action = &dao.ExecAction{Sql: "INSERT INTO task_template_role (id,task_template,custom_role_type,handler_type) VALUES (?,?,?,?)"}
	action.Param = []interface{}{newTaskTemplateRole.Id, newTaskTemplateRole.TaskTemplate, newTaskTemplateRole.CustomRoleType, newTaskTemplateRole.HandlerType}
	actions = append(actions, action)
	// 如果不是尾插，则需更新现有数据的序号
	if param.Sort != len(taskTemplates)+1 {
		for i := param.Sort; i < len(taskTemplates)+1; i++ {
			t := taskTemplates[i-1]
			t.Sort += 1
			t.UpdatedTime = nowTime
			action = &dao.ExecAction{Sql: "UPDATE task_template SET sort = ?, updated_time = ? WHERE id = ?"}
			action.Param = []interface{}{t.Sort, t.UpdatedTime, t.Id}
			actions = append(actions, action)
		}
		tmp := make([]*models.TaskTemplateTable, len(taskTemplates)+1)
		copy(tmp, taskTemplates[:param.Sort-1])
		tmp[param.Sort-1] = newTaskTemplate
		copy(tmp[param.Sort:], taskTemplates[param.Sort-1:])
		taskTemplates = tmp
	} else {
		taskTemplates = append(taskTemplates, newTaskTemplate)
	}
	// 执行事务
	err = dao.Transaction(actions)
	if err != nil {
		return nil, err
	}
	// 构造返回结果
	result := &models.CustomTaskTemplateCreateResponse{
		Id:              newTaskTemplate.Id,
		Sort:            newTaskTemplate.Sort,
		RequestTemplate: newTaskTemplate.RequestTemplate,
		Name:            newTaskTemplate.Name,
		Ids:             make([]*models.CustomTaskTemplateIdObj, len(taskTemplates)),
	}
	for i, taskTemplate := range taskTemplates {
		result.Ids[i] = &models.CustomTaskTemplateIdObj{
			Id:   taskTemplate.Id,
			Sort: taskTemplate.Sort,
			Name: taskTemplate.Name,
		}
	}
	return result, nil
}

func (s TaskTemplateService) UpdateTaskTemplate(param *models.TaskTemplateDto) error {
	actions := []*dao.ExecAction{}
	// 查询现有任务模板
	var taskTemplates []*models.TaskTemplateTable
	err := s.taskTemplateDao.DB.SQL("SELECT * FROM task_template WHERE id = ?", param.Id).Find(&taskTemplates)
	if err != nil {
		return err
	}
	if len(taskTemplates) == 0 {
		return errors.New("no task_template record found")
	}
	taskTemplate := taskTemplates[0]
	// 校验参数
	if taskTemplate.Type != param.Type || taskTemplate.Sort != param.Sort || taskTemplate.RequestTemplate != param.RequestTemplateId {
		return errors.New("param type or sort or requestTemplate wrong")
	}
	for _, roleObj := range param.RoleObjs {
		if roleObj.RoleType != string(models.TaskTemplateRoleRoleTypeTemplate) &&
			(roleObj.HandlerType == string(models.TaskTemplateRoleHandlerTypeTemplate) || roleObj.HandlerType == string(models.TaskTemplateRoleHandlerTypeTemplateSuggest)) {
			return fmt.Errorf("roleType %s not match handlerType %s", roleObj.RoleType, roleObj.HandlerType)
		}
		if roleObj.RoleType != string(models.TaskTemplateRoleRoleTypeTemplate) && roleObj.Role != "" {
			return fmt.Errorf("roleType %s not match role %s", roleObj.RoleType, roleObj.Role)
		}
		if roleObj.HandlerType != string(models.TaskTemplateRoleHandlerTypeTemplate) && roleObj.HandlerType != string(models.TaskTemplateRoleHandlerTypeTemplateSuggest) && roleObj.Handler != "" {
			return fmt.Errorf("handlerType %s not match handler %s", roleObj.HandlerType, roleObj.Handler)
		}
	}
	// 更新现有任务模板
	nowTime := time.Now().Format(models.DateTimeFormat)
	action := &dao.ExecAction{Sql: "UPDATE task_template SET name = ?, expire_day = ?, description = ?, role_type = ?, updated_time = ? WHERE id = ?"}
	action.Param = []interface{}{param.Name, param.ExpireDay, param.Description, param.RoleType, nowTime, param.Id}
	actions = append(actions, action)
	// 增删改现有任务处理模板
	if param.RoleType == string(models.TaskTemplateRoleTypeAdmin) {
		action = &dao.ExecAction{Sql: "DELETE FROM task_template_role WHERE task_template = ?"}
		action.Param = []interface{}{param.Id}
		actions = append(actions, action)
	} else {
		// 查询现有任务处理模板
		var taskTemplateRoles []*models.TaskTemplateRoleTable
		err = s.taskTemplateRoleDao.DB.SQL("SELECT * FROM task_template_role WHERE task_template = ?", param.Id).Find(&taskTemplateRoles)
		if err != nil {
			return err
		}
		// 对比增删改
		for i, taskTemplateRole := range taskTemplateRoles {
			if i < len(param.RoleObjs) {
				roleObj := param.RoleObjs[i]
				action = &dao.ExecAction{Sql: "UPDATE task_template_role SET custom_role_type = ?, handler_type = ?, custom_role = ?, handler = ? WHERE id = ?"}
				action.Param = []interface{}{roleObj.RoleType, roleObj.HandlerType, roleObj.Role, roleObj.Handler, taskTemplateRole.Id}
				actions = append(actions, action)
			} else {
				action = &dao.ExecAction{Sql: "DELETE FROM task_template_role WHERE id = ?"}
				action.Param = []interface{}{taskTemplateRole.Id}
				actions = append(actions, action)
			}
		}
		for sort := len(taskTemplateRoles) + 1; sort <= len(param.RoleObjs); sort++ {
			roleObj := param.RoleObjs[sort-1]
			action = &dao.ExecAction{Sql: "INSERT INTO task_template_role (id,task_template,custom_role_type,handler_type,custom_role,handler) VALUES (?,?,?,?,?,?)"}
			action.Param = []interface{}{guid.CreateGuid(), param.Id, roleObj.RoleType, roleObj.HandlerType, roleObj.Role, roleObj.Handler}
			actions = append(actions, action)
		}
	}
	// 执行事务
	err = dao.Transaction(actions)
	if err != nil {
		return err
	}
	return nil
}

func (s TaskTemplateService) DeleteCustomTaskTemplate(id string) (*models.CustomTaskTemplateDeleteResponse, error) {
	actions := []*dao.ExecAction{}
	// 查询现有任务模板
	var taskTemplates []*models.TaskTemplateTable
	err := s.taskTemplateDao.DB.SQL("SELECT * FROM task_template WHERE id = ?", id).Find(&taskTemplates)
	if err != nil {
		return nil, err
	}
	if len(taskTemplates) == 0 {
		return nil, errors.New("no task_template record found")
	}
	taskTemplate := taskTemplates[0]
	// 校验参数
	if taskTemplate.Type != string(models.TaskTemplateTypeCustom) {
		return nil, errors.New("type wrong")
	}
	// 删除现有任务处理模板
	action := &dao.ExecAction{Sql: "DELETE FROM task_template_role WHERE task_template = ?"}
	action.Param = []interface{}{id}
	actions = append(actions, action)
	// 删除现有任务模板
	action = &dao.ExecAction{Sql: "DELETE FROM task_template WHERE id = ?"}
	action.Param = []interface{}{id}
	actions = append(actions, action)
	// 查询剩余任务模板列表
	taskTemplates = nil
	err = s.taskTemplateDao.DB.SQL("SELECT * FROM task_template WHERE request_template = ? AND id != ? ORDER BY sort", taskTemplate.RequestTemplate, id).Find(&taskTemplates)
	if err != nil {
		return nil, err
	}
	// 如果不是尾删，则需更新现有数据的序号
	if taskTemplate.Sort != len(taskTemplates)+1 {
		nowTime := time.Now().Format(models.DateTimeFormat)
		for i := taskTemplate.Sort; i < len(taskTemplates)+1; i++ {
			t := taskTemplates[i-1]
			t.Sort = i
			t.UpdatedTime = nowTime
			action = &dao.ExecAction{Sql: "UPDATE task_template SET sort = ?, updated_time = ? WHERE id = ?"}
			action.Param = []interface{}{t.Sort, t.UpdatedTime, t.Id}
			actions = append(actions, action)
		}
	}
	// 执行事务
	err = dao.Transaction(actions)
	if err != nil {
		return nil, err
	}
	// 构造返回结果
	result := &models.CustomTaskTemplateDeleteResponse{
		Ids: make([]*models.CustomTaskTemplateIdObj, len(taskTemplates)),
	}
	for i, taskTemplate := range taskTemplates {
		result.Ids[i] = &models.CustomTaskTemplateIdObj{
			Id:   taskTemplate.Id,
			Sort: taskTemplate.Sort,
			Name: taskTemplate.Name,
		}
	}
	return result, nil
}

func (s TaskTemplateService) GetTaskTemplate(id string) (*models.TaskTemplateDto, error) {
	// 查询现有任务模板
	var taskTemplates []*models.TaskTemplateTable
	err := s.taskTemplateDao.DB.SQL("SELECT * FROM task_template WHERE id = ?", id).Find(&taskTemplates)
	if err != nil {
		return nil, err
	}
	if len(taskTemplates) == 0 {
		return nil, errors.New("no task_template record found")
	}
	taskTemplate := taskTemplates[0]
	// 查询现有任务处理模板
	var taskTemplateRoles []*models.TaskTemplateRoleTable
	err = s.taskTemplateRoleDao.DB.SQL("SELECT * FROM task_template_role WHERE task_template = ?", id).Find(&taskTemplateRoles)
	if err != nil {
		return nil, err
	}
	// 构造返回结果
	result := &models.TaskTemplateDto{
		Id:                taskTemplate.Id,
		Type:              taskTemplate.Type,
		Sort:              taskTemplate.Sort,
		RequestTemplateId: taskTemplate.RequestTemplate,
		Name:              taskTemplate.Name,
		ExpireDay:         taskTemplate.ExpireDay,
		Description:       taskTemplate.Description,
		RoleType:          taskTemplate.RoleType,
		RoleObjs:          make([]*models.TaskTemplateRoleDto, len(taskTemplateRoles)),
	}
	for i, taskTemplateRole := range taskTemplateRoles {
		result.RoleObjs[i] = &models.TaskTemplateRoleDto{
			RoleType:    taskTemplateRole.CustomRoleType,
			HandlerType: taskTemplateRole.HandlerType,
			Role:        taskTemplateRole.CustomRole,
			Handler:     taskTemplateRole.Handler,
		}
	}
	return result, nil
}

func (s TaskTemplateService) ListCustomTaskTemplateIds(requestTemplateId string) (*models.CustomTaskTemplateListIdsResponse, error) {
	// 查询现有任务模板列表
	var taskTemplates []*models.TaskTemplateTable
	err := s.taskTemplateDao.DB.SQL("SELECT * FROM task_template WHERE request_template = ? AND type = ? ORDER BY sort", requestTemplateId, string(models.TaskTemplateTypeCustom)).Find(&taskTemplates)
	if err != nil {
		return nil, err
	}
	// 构造返回结果
	result := &models.CustomTaskTemplateListIdsResponse{
		Ids: make([]*models.CustomTaskTemplateIdObj, len(taskTemplates)),
	}
	for i, taskTemplate := range taskTemplates {
		result.Ids[i] = &models.CustomTaskTemplateIdObj{
			Id:   taskTemplate.Id,
			Sort: taskTemplate.Sort,
			Name: taskTemplate.Name,
		}
	}
	return result, nil
}

func (s TaskTemplateService) ListTaskTemplates(requestTemplateId string) ([]*models.TaskTemplateDto, error) {
	result := []*models.TaskTemplateDto{}
	// 查询现有任务模板列表
	var taskTemplates []*models.TaskTemplateTable
	err := s.taskTemplateDao.DB.SQL("SELECT * FROM task_template WHERE request_template = ? ORDER BY sort", requestTemplateId).Find(&taskTemplates)
	if err != nil {
		return nil, err
	}
	if len(taskTemplates) == 0 {
		return result, nil
	}
	// 汇总任务模板列表
	taskTemplateIds := make([]string, len(taskTemplates))
	for i, taskTemplate := range taskTemplates {
		taskTemplateIds[i] = taskTemplate.Id
	}
	// 查询现有任务处理模板列表
	var taskTemplateRoles []*models.TaskTemplateRoleTable
	err = s.taskTemplateRoleDao.DB.SQL("SELECT * FROM task_template_role WHERE task_template IN ('" + strings.Join(taskTemplateIds, "','") + "') ORDER BY task_template").Find(&taskTemplateRoles)
	if err != nil {
		return nil, err
	}
	// 汇总任务处理模板列表
	taskTemplateRoleMap := make(map[string][]*models.TaskTemplateRoleTable)
	taskTemplateId := ""
	for _, taskTemplateRole := range taskTemplateRoles {
		if taskTemplateId != taskTemplateRole.TaskTemplate {
			taskTemplateId = taskTemplateRole.TaskTemplate
			taskTemplateRoleMap[taskTemplateId] = make([]*models.TaskTemplateRoleTable, 0)
		}
		taskTemplateRoleMap[taskTemplateId] = append(taskTemplateRoleMap[taskTemplateId], taskTemplateRole)
	}
	// 构造返回结果
	result = make([]*models.TaskTemplateDto, len(taskTemplates))
	for i, taskTemplate := range taskTemplates {
		result[i] = &models.TaskTemplateDto{
			Id:                taskTemplate.Id,
			Type:              taskTemplate.Type,
			Sort:              taskTemplate.Sort,
			RequestTemplateId: taskTemplate.RequestTemplate,
			Name:              taskTemplate.Name,
			ExpireDay:         taskTemplate.ExpireDay,
			Description:       taskTemplate.Description,
			RoleType:          taskTemplate.RoleType,
		}
		if roleObjs, ok := taskTemplateRoleMap[taskTemplate.Id]; ok {
			result[i].RoleObjs = make([]*models.TaskTemplateRoleDto, len(roleObjs))
			for j, roleObj := range roleObjs {
				result[i].RoleObjs[j] = &models.TaskTemplateRoleDto{
					RoleType:    roleObj.CustomRoleType,
					HandlerType: roleObj.HandlerType,
					Role:        roleObj.CustomRole,
					Handler:     roleObj.Handler,
				}
			}
		}
	}
	return result, nil
}

func (s *TaskTemplateService) DeleteTaskTemplates(requestTemplateId string) ([]*dao.ExecAction, error) {
	// 查询现有任务模板列表
	var taskTemplates []*models.TaskTemplateTable
	err := s.taskTemplateDao.DB.SQL("SELECT * FROM task_template WHERE request_template = ? ORDER BY sort", requestTemplateId).Find(&taskTemplates)
	if err != nil {
		return nil, err
	}
	if len(taskTemplates) == 0 {
		return nil, nil
	}
	// 汇总任务模板列表
	taskTemplateIds := make([]string, len(taskTemplates))
	for i, taskTemplate := range taskTemplates {
		taskTemplateIds[i] = taskTemplate.Id
	}
	actions := []*dao.ExecAction{}
	// 删除现有任务处理模板
	action := &dao.ExecAction{Sql: "DELETE FROM task_template_role WHERE task_template IN ('" + strings.Join(taskTemplateIds, "','") + "')"}
	actions = append(actions, action)
	// 删除现有任务模板
	action = &dao.ExecAction{Sql: "DELETE FROM task_template WHERE id IN ('" + strings.Join(taskTemplateIds, "','") + "')"}
	actions = append(actions, action)
	return actions, nil
}
