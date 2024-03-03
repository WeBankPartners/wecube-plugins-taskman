package service

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/rpc"
	"xorm.io/xorm"
)

type TaskTemplateService struct {
	taskTemplateDao       *dao.TaskTemplateDao
	taskHandleTemplateDao *dao.TaskHandleTemplateDao
}

// func GetTaskTemplate(requestTemplateId, proNodeId, nodeId string) (result models.TaskTemplateDto, err error) {
// 	result = models.TaskTemplateDto{Items: []*models.FormItemTemplateDto{}}
// 	taskTemplate, getTaskErr := getSimpleTaskTemplate("", requestTemplateId, proNodeId, nodeId)
// 	if getTaskErr != nil {
// 		log.Logger.Warn("GetTaskTemplate warning", log.Error(getTaskErr))
// 		return
// 	}
// 	result.Id = taskTemplate.Id
// 	result.NodeDefId = taskTemplate.NodeDefId
// 	result.NodeDefName = taskTemplate.NodeName
// 	result.ExpireDay = taskTemplate.ExpireDay
// 	result.Handler = taskTemplate.Handler
// 	result.MGMTRoleObjs = []*models.RoleTable{}
// 	result.MGMTRoles = []string{}
// 	result.USERoleObjs = []*models.RoleTable{}
// 	result.USERoles = []string{}
// 	result.RequestTemplateId = requestTemplateId
// 	var formTemplateTable []*models.FormTemplateTable
// 	err = dao.X.SQL("select * from form_template where id=?", taskTemplate.FormTemplate).Find(&formTemplateTable)
// 	if err != nil {
// 		err = fmt.Errorf("Try to query form template table fail,%s ", err.Error())
// 		return
// 	}
// 	if len(formTemplateTable) == 0 {
// 		err = fmt.Errorf("Can not find any form template with id=%s ", taskTemplate.FormTemplate)
// 		return
// 	}
// 	var formItemTemplate []*models.FormItemTemplateTable
// 	var formItemTemplateGroup models.FormTemplateTable
// 	dao.X.SQL("select * from form_item_template where form_template=?", taskTemplate.FormTemplate).Find(&formItemTemplate)
// 	if len(formItemTemplate) > 0 && formItemTemplate[0].FormTemplate != "" {
// 		dao.X.SQL("select * from form_item_template_group where id=?", formItemTemplate[0].FormTemplate).Get(&formItemTemplateGroup)
// 		for _, formItem := range formItemTemplate {
// 			result.Items = append(result.Items, models.ConvertFormItemTemplateModel2Dto(formItem, formItemTemplateGroup))
// 		}
// 	}
// 	roleMap, _ := getRoleMap()
// 	var taskRoleTable []*models.TaskTemplateRoleTable
// 	dao.X.SQL("select `role`,role_type from task_template_role where task_template=?", taskTemplate.Id).Find(&taskRoleTable)
// 	for _, role := range taskRoleTable {
// 		if role.RoleType == "MGMT" {
// 			result.MGMTRoleObjs = append(result.MGMTRoleObjs, roleMap[role.Role])
// 			result.MGMTRoles = append(result.MGMTRoles, role.Role)
// 		} else {
// 			result.USERoleObjs = append(result.USERoleObjs, roleMap[role.Role])
// 			result.USERoles = append(result.USERoles, role.Role)
// 		}
// 	}
// 	return
// }

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

// func CreateTaskTemplate(param models.TaskTemplateDto, requestTemplateId string) error {
// 	_, checkExistErr := getSimpleTaskTemplate("", requestTemplateId, param.NodeDefId, "")
// 	if checkExistErr == nil {
// 		return fmt.Errorf("RequestTemplate:%s nodeDefId:%s already have task form,please reload ", requestTemplateId, param.NodeDefId)
// 	}
// 	nowTime := time.Now().Format(models.DateTimeFormat)
// 	formCreateParam := models.FormTemplateDto{Name: param.Name, Description: param.Description, UpdatedBy: param.UpdatedBy, Items: param.Items, NowTime: nowTime}
// 	actions, formId := getFormTemplateCreateActions(formCreateParam)
// 	taskTemplateId := guid.CreateGuid()
// 	actions = append(actions, &dao.ExecAction{Sql: "insert into task_template(id,name,description,form_template,request_template,node_id,node_def_id,node_name,expire_day,handler) value (?,?,?,?,?,?,?,?,?,?)", Param: []interface{}{taskTemplateId, param.Name, param.Description, formId, requestTemplateId, param.NodeId, param.NodeDefId, param.NodeDefName, param.ExpireDay, param.Handler}})
// 	for _, v := range param.MGMTRoles {
// 		actions = append(actions, &dao.ExecAction{Sql: "insert into task_template_role(id,task_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{taskTemplateId + models.SysTableIdConnector + v + models.SysTableIdConnector + "MGMT", taskTemplateId, v, "MGMT"}})
// 	}
// 	for _, v := range param.USERoles {
// 		actions = append(actions, &dao.ExecAction{Sql: "insert into task_template_role(id,task_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{taskTemplateId + models.SysTableIdConnector + v + models.SysTableIdConnector + "USE", taskTemplateId, v, "USE"}})
// 	}
// 	return dao.TransactionWithoutForeignCheck(actions)
// }

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

// func getTaskTemplateHandler(requestTemplate string) (taskTemplateMap map[string]*models.TaskTemplateVo, err error) {
// 	taskTemplateMap = make(map[string]*models.TaskTemplateVo)
// 	var taskTemplateList []*models.TaskTemplateVo
// 	err = dao.X.SQL("select  t.id,t.handler,tr.role from task_template t join task_template_role tr on  "+
// 		"t.id=tr.task_template where t.request_template = ?", requestTemplate).Find(&taskTemplateList)
// 	if len(taskTemplateList) > 0 {
// 		for _, taskTemplate := range taskTemplateList {
// 			taskTemplateMap[taskTemplate.Id] = taskTemplate
// 		}
// 	}
// 	return
// }

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

func (s *TaskTemplateService) CreateTaskTemplate(param *models.TaskTemplateDto, operator string) (*models.TaskTemplateCreateResponse, error) {
	// 校验参数
	guidPrefix, err := s.genTaskIdPrefix(param.Type)
	if err != nil {
		return nil, err
	}
	// 查询请求模板
	requestTemplate, err := GetRequestTemplateService().GetRequestTemplate(param.RequestTemplate)
	if err != nil {
		return nil, err
	}
	if requestTemplate == nil {
		return nil, fmt.Errorf("no request_template record found: %s", param.RequestTemplate)
	}
	if requestTemplate.ProcDefId != "" {
		return nil, fmt.Errorf("param requestTemplate type proc: %s", param.RequestTemplate)
	}
	// 查询任务模板列表
	taskTemplates, err := s.taskTemplateDao.QueryByRequestTemplateAndType(param.RequestTemplate, param.Type)
	if err != nil {
		return nil, err
	}
	if param.Sort <= 0 || param.Sort > len(taskTemplates)+1 {
		return nil, fmt.Errorf("param sort out of range: %d", param.Sort)
	}
	// 插入新任务模板
	nowTime := time.Now().Format(models.DateTimeFormat)
	newTaskTemplate := &models.TaskTemplateTable{
		Id:              fmt.Sprintf("%s_%s", guidPrefix, guid.CreateGuid()),
		Name:            param.Name,
		RequestTemplate: param.RequestTemplate,
		ExpireDay:       param.ExpireDay,
		CreatedBy:       operator,
		CreatedTime:     nowTime,
		UpdatedBy:       operator,
		UpdatedTime:     nowTime,
		Type:            param.Type,
		Sort:            param.Sort,
		HandleMode:      string(models.TaskTemplateHandleModeCustom),
	}
	// 插入新任务处理模板
	newTaskHandleTemplate := &models.TaskHandleTemplateTable{
		Id:           guid.CreateGuid(),
		TaskTemplate: newTaskTemplate.Id,
		Assign:       string(models.TaskHandleTemplateAssignTypeCustom),
		HandlerType:  string(models.TaskHandleTemplateHandlerTypeCustom),
		HandleMode:   newTaskTemplate.HandleMode,
	}
	// 如果不是尾插，则需更新现有任务模板的序号
	var updateTaskTemplates []*models.TaskTemplateTable
	if param.Sort != len(taskTemplates)+1 {
		for i := param.Sort; i < len(taskTemplates)+1; i++ {
			t := taskTemplates[i-1]
			t.Sort += 1

			updateTaskTemplate := &models.TaskTemplateTable{
				Id:          t.Id,
				Sort:        t.Sort,
				UpdatedBy:   operator,
				UpdatedTime: nowTime,
			}
			updateTaskTemplates = append(updateTaskTemplates, updateTaskTemplate)
		}
	}
	// 执行事务
	err = transaction(func(session *xorm.Session) error {
		_, err := s.taskTemplateDao.Add(session, newTaskTemplate)
		if err != nil {
			return err
		}
		_, err = s.taskHandleTemplateDao.Add(session, newTaskHandleTemplate)
		if err != nil {
			return err
		}
		for _, updateTaskTemplate := range updateTaskTemplates {
			err = s.taskTemplateDao.Update(session, updateTaskTemplate)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// 构造返回结果
	resultDto, err := s.genTaskTemplateDto(newTaskTemplate.Id)
	if err != nil {
		return nil, err
	}
	resultIds, err := s.genTaskTemplateIds(param.RequestTemplate, param.Type)
	if err != nil {
		return nil, err
	}
	result := &models.TaskTemplateCreateResponse{
		TaskTemplate: resultDto,
		Ids:          resultIds,
	}
	return result, nil
}

func (s *TaskTemplateService) checkHandleTemplates(param *models.TaskTemplateDto) error {
	if param.HandleMode == string(models.TaskTemplateHandleModeAdmin) ||
		param.HandleMode == string(models.TaskTemplateHandleModeAuto) {
		if len(param.HandleTemplates) > 0 {
			return fmt.Errorf("param handleMode %s should not has handleTemplates", param.HandleMode)
		}
	} else {
		if len(param.HandleTemplates) == 0 {
			return fmt.Errorf("param handleMode %s need handleTemplates", param.HandleMode)
		}
		for _, handleTemplate := range param.HandleTemplates {
			if handleTemplate.Assign != string(models.TaskHandleTemplateAssignTypeTemplate) &&
				(handleTemplate.HandlerType == string(models.TaskHandleTemplateHandlerTypeTemplate) || handleTemplate.HandlerType == string(models.TaskHandleTemplateHandlerTypeTemplateSuggest)) {
				return fmt.Errorf("param assign %s not match handlerType %s", handleTemplate.Assign, handleTemplate.HandlerType)
			}
			if handleTemplate.Assign != string(models.TaskHandleTemplateAssignTypeTemplate) && handleTemplate.Role != "" {
				return fmt.Errorf("param assign %s not match role %s", handleTemplate.Assign, handleTemplate.Role)
			}
			if handleTemplate.HandlerType != string(models.TaskHandleTemplateHandlerTypeTemplate) && handleTemplate.HandlerType != string(models.TaskHandleTemplateHandlerTypeTemplateSuggest) && handleTemplate.Handler != "" {
				return fmt.Errorf("param handlerType %s not match handler %s", handleTemplate.HandlerType, handleTemplate.Handler)
			}
		}
	}
	return nil
}

func (s *TaskTemplateService) CreateProcTaskTemplate(param *models.TaskTemplateDto, userToken, language, operator string) (*models.TaskTemplateDto, error) {
	// 校验参数
	if param.Type != string(models.TaskTypeImplement) {
		return nil, fmt.Errorf("param type wrong: %s", param.Type)
	}
	if param.NodeDefId == "" {
		return nil, errors.New("param empty")
	}
	if err := s.checkHandleTemplates(param); err != nil {
		return nil, err
	}
	// 查询请求模板
	requestTemplate, err := GetRequestTemplateService().GetRequestTemplate(param.RequestTemplate)
	if err != nil {
		return nil, err
	}
	if requestTemplate == nil {
		return nil, errors.New("no request_template record found")
	}
	if requestTemplate.ProcDefId != param.NodeDefId {
		return nil, fmt.Errorf("param nodeDefId %s wrong", param.NodeDefId)
	}
	// 查询现有任务模板
	taskTemplate, err := s.taskTemplateDao.GetProc(param.RequestTemplate, param.NodeDefId)
	if err != nil {
		return nil, err
	}
	if taskTemplate != nil {
		return nil, errors.New("task_template record already exist")
	}
	// 查询编排任务节点
	nodeList, err := s.getProcTaskTemplateNodes(requestTemplate.ProcDefId, userToken, language)
	if err != nil {
		return nil, err
	}
	sort := 0
	for i, node := range nodeList {
		if node.NodeDefId == param.NodeDefId {
			sort = i + 1
			if node.NodeId != param.NodeId || node.NodeName != param.NodeDefName {
				return nil, fmt.Errorf("param nodeId %q or nodeName %q wrong", param.NodeId, param.NodeDefName)
			}
			break
		}
	}
	if param.Sort <= 0 || sort != param.Sort {
		return nil, fmt.Errorf("param sort %d wrong", param.Sort)
	}
	// 插入新任务模板
	nowTime := time.Now().Format(models.DateTimeFormat)
	newTaskTemplate := &models.TaskTemplateTable{
		Id:              fmt.Sprintf("ts_%s", guid.CreateGuid()),
		Type:            param.Type,
		Sort:            param.Sort,
		RequestTemplate: param.RequestTemplate,
		Name:            param.Name,
		Description:     param.Description,
		NodeId:          param.NodeId,
		NodeDefId:       param.NodeDefId,
		NodeName:        param.NodeDefName,
		ExpireDay:       param.ExpireDay,
		CreatedBy:       operator,
		CreatedTime:     nowTime,
		UpdatedBy:       operator,
		UpdatedTime:     nowTime,
		HandleMode:      param.HandleMode,
	}
	// 插入新任务处理模板
	newTaskHandleTemplates := make([]*models.TaskHandleTemplateTable, len(param.HandleTemplates))
	for i, handleTemplate := range param.HandleTemplates {
		newTaskHandleTemplates[i] = &models.TaskHandleTemplateTable{
			Id:           guid.CreateGuid(),
			Sort:         i + 1,
			TaskTemplate: newTaskTemplate.Id,
			Assign:       handleTemplate.Assign,
			HandlerType:  handleTemplate.HandlerType,
			Role:         handleTemplate.Role,
			Handler:      handleTemplate.Handler,
			HandleMode:   param.HandleMode,
		}
	}
	// 执行事务
	err = transaction(func(session *xorm.Session) error {
		_, err := s.taskTemplateDao.Add(session, newTaskTemplate)
		if err != nil {
			return err
		}
		for _, newTaskHandleTemplate := range newTaskHandleTemplates {
			_, err = s.taskHandleTemplateDao.Add(session, newTaskHandleTemplate)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// 构造返回结果
	result, err := s.genTaskTemplateDto(newTaskTemplate.Id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *TaskTemplateService) UpdateTaskTemplate(param *models.TaskTemplateDto, operator string) (*models.TaskTemplateDto, error) {
	// 查询任务模板
	taskTemplate, err := s.taskTemplateDao.Get(param.Id)
	if err != nil {
		return nil, err
	}
	if taskTemplate == nil {
		return nil, errors.New("no task_template record found")
	}
	// 校验参数
	if taskTemplate.Type != param.Type || taskTemplate.Sort != param.Sort || taskTemplate.RequestTemplate != param.RequestTemplate || taskTemplate.NodeId != param.NodeId || taskTemplate.NodeDefId != param.NodeDefId || taskTemplate.NodeName != param.NodeDefName {
		return nil, errors.New("param wrong")
	}
	if err := s.checkHandleTemplates(param); err != nil {
		return nil, err
	}
	// 更新任务模板
	nowTime := time.Now().Format(models.DateTimeFormat)
	updateTaskTemplate := &models.TaskTemplateTable{
		Id:          param.Id,
		Name:        param.Name,
		Description: param.Description,
		ExpireDay:   param.ExpireDay,
		UpdatedBy:   operator,
		UpdatedTime: nowTime,
		HandleMode:  param.HandleMode,
	}
	// 更新任务处理模板
	deleteTaskHandleTemplateAll := false
	var deleteTaskHandleTemplateIds []string
	var updateTaskHandleTemplates []*models.TaskHandleTemplateTable
	var newTaskHandleTemplates []*models.TaskHandleTemplateTable
	if param.HandleMode == string(models.TaskTemplateHandleModeAdmin) || param.HandleMode == string(models.TaskTemplateHandleModeAuto) {
		deleteTaskHandleTemplateAll = true
	} else {
		// 查询任务处理模板
		taskHandleTemplates, err := s.taskHandleTemplateDao.QueryByTaskTemplate(updateTaskTemplate.Id)
		if err != nil {
			return nil, err
		}
		// 对比增删改
		for i, taskHandleTemplate := range taskHandleTemplates {
			if i < len(param.HandleTemplates) {
				newHandleTemplate := param.HandleTemplates[i]
				updateTaskTemplate := &models.TaskHandleTemplateTable{
					Id:          taskHandleTemplate.Id,
					Assign:      newHandleTemplate.Assign,
					HandlerType: newHandleTemplate.HandlerType,
					Role:        newHandleTemplate.Role,
					Handler:     newHandleTemplate.Handler,
					HandleMode:  param.HandleMode,
				}
				updateTaskHandleTemplates = append(updateTaskHandleTemplates, updateTaskTemplate)
			} else {
				deleteTaskHandleTemplateIds = append(deleteTaskHandleTemplateIds, taskHandleTemplate.Id)
			}
		}
		for sort := len(taskHandleTemplates) + 1; sort <= len(param.HandleTemplates); sort++ {
			newHandleTemplate := param.HandleTemplates[sort-1]
			newTaskHandleTemplate := &models.TaskHandleTemplateTable{
				Id:           guid.CreateGuid(),
				Sort:         sort,
				TaskTemplate: updateTaskTemplate.Id,
				Assign:       newHandleTemplate.Assign,
				HandlerType:  newHandleTemplate.HandlerType,
				Role:         newHandleTemplate.Role,
				Handler:      newHandleTemplate.Handler,
				HandleMode:   param.HandleMode,
			}
			newTaskHandleTemplates = append(newTaskHandleTemplates, newTaskHandleTemplate)
		}
	}
	// 执行事务
	err = transaction(func(session *xorm.Session) error {
		err := s.taskTemplateDao.Update(session, updateTaskTemplate)
		if err != nil {
			return err
		}
		if deleteTaskHandleTemplateAll {
			err = s.taskHandleTemplateDao.DeleteByTaskTemplate(session, updateTaskTemplate.Id)
			if err != nil {
				return err
			}
		} else {
			for _, updateTaskHandleTemplate := range updateTaskHandleTemplates {
				err = s.taskHandleTemplateDao.Update(session, updateTaskHandleTemplate)
				if err != nil {
					return err
				}
			}
			err = s.taskHandleTemplateDao.Deletes(session, deleteTaskHandleTemplateIds)
			if err != nil {
				return err
			}
			for _, newTaskHandleTemplate := range newTaskHandleTemplates {
				_, err = s.taskHandleTemplateDao.Add(session, newTaskHandleTemplate)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// 构造返回结果
	result, err := s.genTaskTemplateDto(updateTaskTemplate.Id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *TaskTemplateService) DeleteTaskTemplate(requestTemplateId, id string) (*models.TaskTemplateDeleteResponse, error) {
	// 查询任务模版
	taskTemplate, err := s.taskTemplateDao.Get(id)
	if err != nil {
		return nil, err
	}
	if taskTemplate == nil {
		return nil, errors.New("no task_template record found")
	}
	// 校验参数
	if taskTemplate.Type != string(models.TaskTypeApprove) && taskTemplate.Type != string(models.TaskTypeImplement) {
		return nil, fmt.Errorf("type wrong: %s", taskTemplate.Type)
	}
	if taskTemplate.NodeDefId != "" {
		return nil, fmt.Errorf("nodeDefId not empty: %s", taskTemplate.NodeDefId)
	}
	if taskTemplate.RequestTemplate != requestTemplateId {
		return nil, fmt.Errorf("param requestTemplate wrong: %s", requestTemplateId)
	}
	// 删除任务处理模板
	deleteTaskHandleTemplateAll := true
	// 删除任务模板
	deleteTaskTemplateId := id
	// 查询任务模板列表
	taskTemplates, err := s.taskTemplateDao.QueryByRequestTemplateAndType(taskTemplate.RequestTemplate, taskTemplate.Type)
	if err != nil {
		return nil, err
	}
	// 如果不是尾删，则需更新现有数据的序号
	var updateTaskTemplates []*models.TaskTemplateTable
	if taskTemplate.Sort != len(taskTemplates) {
		for i := taskTemplate.Sort; i < len(taskTemplates); i++ {
			updateTaskTemplate := &models.TaskTemplateTable{
				Id:   taskTemplates[i].Id,
				Sort: i,
			}
			updateTaskTemplates = append(updateTaskTemplates, updateTaskTemplate)
		}
	}
	// 执行事务
	err = transaction(func(session *xorm.Session) error {
		if deleteTaskHandleTemplateAll {
			err := s.taskHandleTemplateDao.DeleteByTaskTemplate(session, deleteTaskTemplateId)
			if err != nil {
				return err
			}
		}
		err := s.taskTemplateDao.Delete(session, deleteTaskTemplateId)
		if err != nil {
			return err
		}
		for _, updateTaskTemplate := range updateTaskTemplates {
			err = s.taskTemplateDao.Update(session, updateTaskTemplate)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	// 构造返回结果
	result := &models.TaskTemplateDeleteResponse{
		Type: taskTemplate.Type,
	}
	result.Ids, err = s.genTaskTemplateIds(taskTemplate.RequestTemplate, taskTemplate.Type)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *TaskTemplateService) GetTaskTemplate(requestTemplateId, id, userToken, language string) (*models.TaskTemplateDto, error) {
	// 查询请求模板
	requestTemplate, err := GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if err != nil {
		return nil, err
	}
	if requestTemplate == nil {
		return nil, errors.New("no request_template record found")
	}
	var result *models.TaskTemplateDto
	if requestTemplate.ProcDefId != "" {
		// 编排任务
		// 查询任务模板
		taskTemplate, err := s.taskTemplateDao.GetProc(requestTemplateId, id)
		if err != nil {
			return nil, err
		}
		if taskTemplate != nil {
			result, err = s.genTaskTemplateDto(taskTemplate.Id)
			if err != nil {
				return nil, err
			}
		} else {
			// 查询任务节点
			nodeList, err := s.getProcTaskTemplateNodes(requestTemplate.ProcDefId, userToken, language)
			if err != nil {
				return nil, err
			}
			for _, node := range nodeList {
				if node.NodeDefId == id {
					result = s.genProcTaskTemplateDto(node, requestTemplateId)
					break
				}
			}
			if result == nil {
				return nil, errors.New("proc task template not found")
			}
		}
	} else {
		// 其他类型
		// 查询任务模板
		taskTemplate, err := s.taskTemplateDao.Get(id)
		if err != nil {
			return nil, err
		}
		if taskTemplate == nil {
			return nil, errors.New("no task_template record found")
		}
		result, err = s.genTaskTemplateDto(taskTemplate.Id)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (s *TaskTemplateService) getProcTaskTemplateNodes(procDefId string, userToken string, language string) ([]*models.ProcNodeObj, error) {
	allTaskNodeList, err := rpc.GetProcessDefineTaskNodes(procDefId, userToken, language)
	if err != nil {
		return nil, err
	}
	var nodeList []*models.ProcNodeObj
	for _, node := range allTaskNodeList {
		if node.NodeType != string(models.ProcDefNodeTypeHuman) {
			continue
		}
		if node.TaskCategory != "SUTN" {
			continue
		}
		if node.OrderedNo == "" {
			node.OrderedNum = 0
		} else {
			node.OrderedNum, _ = strconv.Atoi(node.OrderedNo)
		}
		nodeList = append(nodeList, node)
	}
	sort.Sort(models.ProcNodeObjList(nodeList))
	return nodeList, nil
}

func (s *TaskTemplateService) ListTaskTemplateIds(requestTemplateId, typ, userToken, language string) (*models.TaskTemplateListIdsResponse, error) {
	// 查询请求模板
	requestTemplate, err := GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if err != nil {
		return nil, err
	}
	if requestTemplate == nil {
		return nil, errors.New("no request_template record found")
	}
	result := &models.TaskTemplateListIdsResponse{}
	if requestTemplate.ProcDefId != "" {
		// 编排任务
		result.Type = string(models.TaskTypeImplement)
		// 查询任务节点
		nodeList, err := s.getProcTaskTemplateNodes(requestTemplate.ProcDefId, userToken, language)
		if err != nil {
			return nil, err
		}
		// 查询任务模板列表
		taskTemplates, err := s.taskTemplateDao.QueryByRequestTemplateAndType(requestTemplateId, typ)
		if err != nil {
			return nil, err
		}
		// 构造返回结果
		result.Ids = make([]*models.TaskTemplateIdObj, len(nodeList))
		taskTemplateIdx := 0
		for i, node := range nodeList {
			if taskTemplateIdx < len(taskTemplates) {
				taskTemplate := taskTemplates[taskTemplateIdx]
				if node.NodeDefId == taskTemplate.NodeDefId {
					result.Ids[i] = &models.TaskTemplateIdObj{
						Id:        taskTemplate.Id,
						Sort:      taskTemplate.Sort,
						Name:      taskTemplate.Name,
						NodeDefId: taskTemplate.NodeDefId,
					}
					taskTemplateIdx++
				} else {
					result.Ids[i] = &models.TaskTemplateIdObj{
						Sort:      node.OrderedNum,
						Name:      node.NodeName,
						NodeDefId: node.NodeDefId,
					}
				}
			} else {
				result.Ids[i] = &models.TaskTemplateIdObj{
					Sort:      node.OrderedNum,
					Name:      node.NodeName,
					NodeDefId: node.NodeDefId,
				}
			}
		}
	} else {
		// 其他类型
		// 查询任务模板列表
		taskTemplates, err := s.taskTemplateDao.QueryByRequestTemplateAndType(requestTemplateId, typ)
		if err != nil {
			return nil, err
		}
		// 构造返回结果
		result.Ids = make([]*models.TaskTemplateIdObj, len(taskTemplates))
		for i, taskTemplate := range taskTemplates {
			result.Ids[i] = &models.TaskTemplateIdObj{
				Id:        taskTemplate.Id,
				Sort:      taskTemplate.Sort,
				Name:      taskTemplate.Name,
				NodeDefId: taskTemplate.NodeDefId,
			}
		}
		if len(taskTemplates) > 0 {
			result.Type = taskTemplates[0].Type
		}
	}
	return result, nil
}

func (s *TaskTemplateService) ListTaskTemplates(requestTemplateId, typ, userToken, language string) ([]*models.TaskTemplateDto, error) {
	// 查询请求模板
	requestTemplate, err := GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if err != nil {
		return nil, err
	}
	if requestTemplate == nil {
		return nil, errors.New("no request_template record found")
	}
	var result []*models.TaskTemplateDto
	if requestTemplate.ProcDefId != "" {
		// 编排任务
		// 查询任务节点
		nodeList, err := s.getProcTaskTemplateNodes(requestTemplate.ProcDefId, userToken, language)
		if err != nil {
			return nil, err
		}
		// 查询任务模板列表
		taskTemplates, err := s.taskTemplateDao.QueryByRequestTemplateAndType(requestTemplateId, typ)
		if err != nil {
			return nil, err
		}
		// 构造返回结果
		result = make([]*models.TaskTemplateDto, len(nodeList))
		taskTemplateIdx := 0
		for i, node := range nodeList {
			if taskTemplateIdx < len(taskTemplates) {
				taskTemplate := taskTemplates[taskTemplateIdx]
				if node.NodeDefId == taskTemplate.NodeDefId {
					dto, err := s.genTaskTemplateDto(taskTemplate.Id)
					if err != nil {
						return nil, err
					}
					result[i] = dto
					taskTemplateIdx++
				} else {
					result[i] = s.genProcTaskTemplateDto(node, requestTemplateId)
				}
			} else {
				result[i] = s.genProcTaskTemplateDto(node, requestTemplateId)
			}
		}
	} else {
		// 其他类型
		// 查询任务模板列表
		taskTemplates, err := s.taskTemplateDao.QueryByRequestTemplateAndType(requestTemplateId, typ)
		if err != nil {
			return nil, err
		}
		// 构造返回结果
		result = make([]*models.TaskTemplateDto, len(taskTemplates))
		for i, taskTemplate := range taskTemplates {
			result[i], err = s.genTaskTemplateDto(taskTemplate.Id)
			if err != nil {
				return nil, err
			}
		}
	}
	return result, nil
}

func (s *TaskTemplateService) DeleteTaskTemplates(requestTemplateId string) (func(*xorm.Session) error, error) {
	// 查询任务模板列表
	taskTemplates, err := s.taskTemplateDao.QueryByRequestTemplate(requestTemplateId)
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
	result := func(session *xorm.Session) error {
		// 删除任务处理模板
		err := s.taskHandleTemplateDao.DeleteByTaskTemplates(session, taskTemplateIds)
		if err != nil {
			return err
		}
		// 删除任务模板
		err = s.taskTemplateDao.Deletes(session, taskTemplateIds)
		if err != nil {
			return err
		}
		return nil
	}
	return result, nil
}

func (s *TaskTemplateService) genTaskTemplateDto(taskTemplateId string) (*models.TaskTemplateDto, error) {
	taskTemplate, err := s.taskTemplateDao.Get(taskTemplateId)
	if err != nil {
		return nil, err
	}
	if taskTemplate == nil {
		return nil, nil
	}
	taskTemplateHandles, err := s.taskHandleTemplateDao.QueryByTaskTemplate(taskTemplateId)
	if err != nil {
		return nil, err
	}
	result := &models.TaskTemplateDto{
		Id:              taskTemplate.Id,
		Type:            taskTemplate.Type,
		NodeId:          taskTemplate.NodeId,
		NodeDefId:       taskTemplate.NodeDefId,
		NodeDefName:     taskTemplate.NodeName,
		Name:            taskTemplate.Name,
		Description:     taskTemplate.Description,
		ExpireDay:       taskTemplate.ExpireDay,
		UpdatedTime:     taskTemplate.UpdatedTime,
		UpdatedBy:       taskTemplate.UpdatedBy,
		RequestTemplate: taskTemplate.RequestTemplate,
		Sort:            taskTemplate.Sort,
		HandleMode:      taskTemplate.HandleMode,
		HandleTemplates: make([]*models.TaskHandleTemplateDto, len(taskTemplateHandles)),
	}
	for i, taskTemplateHandle := range taskTemplateHandles {
		result.HandleTemplates[i] = &models.TaskHandleTemplateDto{
			Role:        taskTemplateHandle.Role,
			Assign:      taskTemplateHandle.Assign,
			HandlerType: taskTemplateHandle.HandlerType,
			Handler:     taskTemplateHandle.Handler,
		}
	}
	return result, nil
}

func (s *TaskTemplateService) genProcTaskTemplateDto(node *models.ProcNodeObj, requestTemplateId string) *models.TaskTemplateDto {
	return &models.TaskTemplateDto{
		Type:            string(models.TaskTypeImplement),
		NodeId:          node.NodeId,
		NodeDefId:       node.NodeDefId,
		NodeDefName:     node.NodeName,
		Name:            node.NodeName,
		ExpireDay:       1,
		RequestTemplate: requestTemplateId,
		Sort:            node.OrderedNum,
		HandleMode:      string(models.TaskTemplateHandleModeCustom),
		HandleTemplates: []*models.TaskHandleTemplateDto{
			{
				Assign:      string(models.TaskHandleTemplateAssignTypeTemplate),
				HandlerType: string(models.TaskHandleTemplateHandlerTypeTemplate),
			},
		},
	}
}

func (s *TaskTemplateService) genTaskTemplateIds(requestTemplateId, typ string) ([]*models.TaskTemplateIdObj, error) {
	taskTemplates, err := s.taskTemplateDao.QueryByRequestTemplateAndType(requestTemplateId, typ)
	if err != nil {
		return nil, err
	}
	result := make([]*models.TaskTemplateIdObj, len(taskTemplates))
	for i, taskTemplate := range taskTemplates {
		result[i] = &models.TaskTemplateIdObj{
			Id:        taskTemplate.Id,
			Sort:      taskTemplate.Sort,
			Name:      taskTemplate.Name,
			NodeDefId: taskTemplate.NodeDefId,
		}
	}
	return result, nil
}

func (s *TaskTemplateService) genTaskIdPrefix(typ string) (string, error) {
	switch typ {
	case string(models.TaskTypeCheck):
		return "ch", nil
	case string(models.TaskTypeApprove):
		return "ap", nil
	case string(models.TaskTypeImplement):
		return "ts", nil
	case string(models.TaskTypeConfirm):
		return "co", nil
	default:
		return "", fmt.Errorf("type invalid: %s", typ)
	}
}
