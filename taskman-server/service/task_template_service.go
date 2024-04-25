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
	if requestTemplate.ProcDefId != "" && param.Type == string(models.TaskTypeImplement) {
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
		Assign:       string(models.TaskHandleTemplateAssignTypeTemplate),
		HandlerType:  string(models.TaskHandleTemplateHandlerTypeTemplateSuggest),
		HandleMode:   newTaskTemplate.HandleMode,
	}
	// 如果不是尾插，则需更新现有任务模板的序号
	var updateTaskTemplates []*models.TaskTemplateTable
	if param.Sort != len(taskTemplates)+1 {
		for i := param.Sort; i < len(taskTemplates)+1; i++ {
			t := taskTemplates[i-1]
			t.Sort = i + 1
			t.UpdatedBy = operator
			t.UpdatedTime = nowTime

			updateTaskTemplate := t
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

func (s *TaskTemplateService) CheckHandleTemplates(param *models.TaskTemplateDto) error {
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
			if handleTemplate.Assign == string(models.TaskHandleTemplateAssignTypeTemplate) {
				if handleTemplate.Role == "" {
					return fmt.Errorf("param assign %s not match role %s", handleTemplate.Assign, handleTemplate.Role)
				}
			} else {
				if handleTemplate.Role != "" {
					return fmt.Errorf("param assign %s not match role %s", handleTemplate.Assign, handleTemplate.Role)
				}
			}
			if handleTemplate.HandlerType == string(models.TaskHandleTemplateHandlerTypeTemplate) || handleTemplate.HandlerType == string(models.TaskHandleTemplateHandlerTypeTemplateSuggest) {
				if handleTemplate.Handler == "" {
					return fmt.Errorf("param handlerType %s not match handler %s", handleTemplate.HandlerType, handleTemplate.Handler)
				}
			} else {
				if handleTemplate.Handler != "" {
					return fmt.Errorf("param handlerType %s not match handler %s", handleTemplate.HandlerType, handleTemplate.Handler)
				}
			}
		}
	}
	return nil
}

func (s *TaskTemplateService) createProcTaskTemplates(session *xorm.Session, procDefId, requestTemplateId, userToken, language, operator string) (err error) {
	var nodeList []*models.ProcNodeObj
	// 查询编排任务节点
	nodeList, err = s.getProcTaskTemplateNodes(procDefId, userToken, language)
	if err != nil {
		return
	}
	var newTaskTemplates []*models.TaskTemplateTable
	var newTaskHandleTemplates []*models.TaskHandleTemplateTable
	nowTime := time.Now().Format(models.DateTimeFormat)
	for i, node := range nodeList {
		// 插入新任务模板
		newTaskTemplate := &models.TaskTemplateTable{
			Id:              fmt.Sprintf("ts_%s", guid.CreateGuid()),
			Type:            string(models.TaskTypeImplement),
			Sort:            i + 1,
			RequestTemplate: requestTemplateId,
			Name:            node.NodeName,
			Description:     "",
			NodeId:          node.NodeId,
			NodeDefId:       node.NodeDefId,
			NodeName:        node.NodeName,
			ExpireDay:       1,
			CreatedBy:       operator,
			CreatedTime:     nowTime,
			UpdatedBy:       operator,
			UpdatedTime:     nowTime,
			HandleMode:      string(models.TaskTemplateHandleModeCustom),
		}
		// 插入新任务处理模板
		newTaskHandleTemplate := &models.TaskHandleTemplateTable{
			Id:           guid.CreateGuid(),
			Sort:         1,
			TaskTemplate: newTaskTemplate.Id,
			Assign:       string(models.TaskHandleTemplateAssignTypeTemplate),
			HandlerType:  string(models.TaskHandleTemplateHandlerTypeTemplateSuggest),
			Role:         "",
			Handler:      "",
			HandleMode:   newTaskTemplate.HandleMode,
		}
		newTaskTemplates = append(newTaskTemplates, newTaskTemplate)
		newTaskHandleTemplates = append(newTaskHandleTemplates, newTaskHandleTemplate)
	}
	// 构造返回结果
	for _, newTaskTemplate := range newTaskTemplates {
		_, err = s.taskTemplateDao.Add(session, newTaskTemplate)
		if err != nil {
			return
		}
	}
	for _, newTaskHandleTemplate := range newTaskHandleTemplates {
		_, err = s.taskHandleTemplateDao.Add(session, newTaskHandleTemplate)
		if err != nil {
			return
		}
	}
	return
}

func (s *TaskTemplateService) createProcTaskTemplatesSql(procDefId, requestTemplateId, userToken, language, operator string) ([]*dao.ExecAction, error) {
	var actions []*dao.ExecAction
	// 查询编排任务节点
	nodeList, err := s.getProcTaskTemplateNodes(procDefId, userToken, language)
	if err != nil {
		return nil, err
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	for i, node := range nodeList {
		// 插入新任务模板
		taskId := fmt.Sprintf("ts_%s", guid.CreateGuid())
		handleMode := string(models.TaskTemplateHandleModeCustom)
		action := &dao.ExecAction{Sql: "INSERT INTO task_template (id,type,sort,request_template,name,node_id,node_def_id,node_name,expire_day,created_by,created_time,updated_by,updated_time,handle_mode) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
		action.Param = []interface{}{taskId, string(models.TaskTypeImplement), i + 1, requestTemplateId, node.NodeName, node.NodeId, node.NodeDefId, node.NodeName, 1, operator, nowTime, operator, nowTime, handleMode}
		actions = append(actions, action)
		// 插入新任务处理模板
		action = &dao.ExecAction{Sql: "INSERT INTO task_handle_template (id,sort,task_template,assign,handler_type,handle_mode) VALUES (?,?,?,?,?,?)"}
		action.Param = []interface{}{guid.CreateGuid(), 1, taskId, string(models.TaskHandleTemplateAssignTypeTemplate), string(models.TaskHandleTemplateHandlerTypeTemplateSuggest), handleMode}
		actions = append(actions, action)
	}
	return actions, nil
}

// deleteTaskTemplateSql 删除任务模板,模板没有发布不能创建请求,所以只删除关联模板表
func (s *TaskTemplateService) deleteTaskTemplateSql(requestTemplateId, taskTemplateId string) ([]*dao.ExecAction, error) {
	var actions []*dao.ExecAction
	var formTemplateList []*models.FormTemplateTable
	// 查询任务模版
	taskTemplate, err := s.taskTemplateDao.Get(taskTemplateId)
	if err != nil {
		return nil, err
	}
	if taskTemplate == nil {
		return nil, errors.New("no task_template record found")
	}
	// 校验参数
	if taskTemplate.Type != string(models.TaskTypeImplement) {
		return nil, fmt.Errorf("type wrong: %s", taskTemplate.Type)
	}
	if taskTemplate.RequestTemplate != requestTemplateId {
		return nil, fmt.Errorf("param requestTemplate wrong: %s", requestTemplateId)
	}
	dao.X.SQL("select * from form_template where request_template = ? and task_template = ?", requestTemplateId, taskTemplateId).Find(&formTemplateList)
	if len(formTemplateList) > 0 {
		for _, formTemplate := range formTemplateList {
			// 删除表单项模板表
			actions = append(actions, &dao.ExecAction{Sql: "delete from form_item_template WHERE form_template = ? ", Param: []interface{}{formTemplate.Id}})
		}
	}
	// 删除任务处理模板表
	actions = append(actions, &dao.ExecAction{Sql: "delete from task_handle_template WHERE task_template = ? ", Param: []interface{}{taskTemplateId}})
	// 删除表单模板表
	actions = append(actions, &dao.ExecAction{Sql: "delete from form_template WHERE task_template = ?", Param: []interface{}{taskTemplateId}})
	// 删除任务模版表
	actions = append(actions, &dao.ExecAction{Sql: "delete from task_template WHERE id = ?", Param: []interface{}{taskTemplateId}})
	return actions, nil
}

func (s *TaskTemplateService) updateProcTaskTemplatesSql(taskTemplateId, nodeId, nodeName string) ([]*dao.ExecAction, error) {
	var actions []*dao.ExecAction
	nowTime := time.Now().Format(models.DateTimeFormat)
	// 更新任务模板
	action := &dao.ExecAction{Sql: "UPDATE task_template SET node_id=?, node_name=?, updated_time=? WHERE id=?"}
	action.Param = []interface{}{nodeId, nodeName, nowTime, taskTemplateId}
	actions = append(actions, action)
	return actions, nil
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
	if err := s.CheckHandleTemplates(param); err != nil {
		return nil, err
	}
	// 更新任务模板
	nowTime := time.Now().Format(models.DateTimeFormat)
	taskTemplate.Name = param.Name
	taskTemplate.Description = param.Description
	taskTemplate.ExpireDay = param.ExpireDay
	taskTemplate.UpdatedBy = operator
	taskTemplate.UpdatedTime = nowTime
	taskTemplate.HandleMode = param.HandleMode
	updateTaskTemplate := taskTemplate
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
				taskHandleTemplate.Assign = newHandleTemplate.Assign
				taskHandleTemplate.HandlerType = newHandleTemplate.HandlerType
				taskHandleTemplate.Role = newHandleTemplate.Role
				taskHandleTemplate.Handler = newHandleTemplate.Handler
				taskHandleTemplate.HandleMode = param.HandleMode
				taskHandleTemplate.FilterRule = param.HandleTemplates[i].FilterRule
				taskHandleTemplate.AssignRule = param.HandleTemplates[i].AssignRule
				updateTaskTemplate := taskHandleTemplate
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
				FilterRule:   newHandleTemplate.FilterRule,
				AssignRule:   newHandleTemplate.AssignRule,
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
			t := taskTemplates[i]
			t.Sort = i
			updateTaskTemplate := t
			updateTaskTemplates = append(updateTaskTemplates, updateTaskTemplate)
		}
	}
	// 查询表单模板列表
	var formTemplateList []*models.FormTemplateTable
	dao.X.SQL("select * from form_template where request_template = ? and task_template = ?", requestTemplateId, deleteTaskTemplateId).Find(&formTemplateList)
	// 执行事务
	err = transaction(func(session *xorm.Session) error {
		// 先删表单模板
		if len(formTemplateList) > 0 {
			for _, formTemplate := range formTemplateList {
				err = GetFormTemplateService().DeleteFormTemplateItemGroupTransaction(session, formTemplate.Id)
				if err != nil {
					return err
				}
			}
		}
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

func (s *TaskTemplateService) GetTaskTemplate(requestTemplateId, id, typ string) (*models.TaskTemplateDto, error) {
	// 查询请求模板
	requestTemplate, err := GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if err != nil {
		return nil, err
	}
	if requestTemplate == nil {
		return nil, errors.New("no request_template record found")
	}
	// 查询任务模板
	taskTemplate, err := s.taskTemplateDao.Get(id)
	if err != nil {
		return nil, err
	}
	if taskTemplate == nil {
		return nil, errors.New("no task_template record found")
	}
	if taskTemplate.RequestTemplate != requestTemplateId || taskTemplate.Type != typ {
		return nil, errors.New("param requestTemplate or type wrong")
	}
	result, err := s.genTaskTemplateDto(taskTemplate.Id)
	if err != nil {
		return nil, err
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
	result := &models.TaskTemplateListIdsResponse{Type: typ, ProcDefId: requestTemplate.ProcDefId}
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
	return result, nil
}

func (s *TaskTemplateService) ListTaskTemplates(requestTemplateId, typ string) ([]*models.TaskTemplateDto, error) {
	// 查询请求模板
	requestTemplate, err := GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if err != nil {
		return nil, err
	}
	if requestTemplate == nil {
		return nil, errors.New("no request_template record found")
	}
	var result []*models.TaskTemplateDto
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
			Id:          taskTemplateHandle.Id,
			Role:        taskTemplateHandle.Role,
			Assign:      taskTemplateHandle.Assign,
			HandlerType: taskTemplateHandle.HandlerType,
			Handler:     taskTemplateHandle.Handler,
		}
	}
	return result, nil
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
	case string(models.TaskTypeSubmit):
		return "su", nil
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

func (s *TaskTemplateService) QueryTaskTemplateListByRequestTemplateAndType(requestTemplateId, templateType string) (list []*models.TaskTemplateTable, err error) {
	return s.taskTemplateDao.QueryByRequestTemplateAndType(requestTemplateId, templateType)
}

func (s *TaskTemplateService) Get(taskTemplateId string) (result *models.TaskTemplateTable, err error) {
	if taskTemplateId == "" {
		return
	}
	return s.taskTemplateDao.Get(taskTemplateId)
}

func (s *TaskTemplateService) GetTaskTemplateListByRequestId(requestId string) (list []*models.TaskTemplateTable, err error) {
	var requestList []*models.RequestTable
	var requestTemplate string
	err = dao.X.SQL("select * from request where id = ?", requestId).Find(&requestList)
	if err != nil {
		return
	}
	if len(requestList) > 0 {
		requestTemplate = requestList[0].RequestTemplate
	}
	err = dao.X.SQL("select * from task_template where request_template = ?", requestTemplate).Find(&list)
	return
}

func (s *TaskTemplateService) QueryTaskTemplateListByRequestTemplate(requestTemplateId string) (list []*models.TaskTemplateTable, err error) {
	return s.taskTemplateDao.QueryByRequestTemplate(requestTemplateId)
}

func (s *TaskTemplateService) QueryTaskTemplateDtoListByRequestTemplate(requestTemplateId string) (result []*models.TaskTemplateDto, err error) {
	var taskTemplateList []*models.TaskTemplateTable
	if taskTemplateList, err = s.taskTemplateDao.QueryByRequestTemplate(requestTemplateId); err != nil {
		return
	}
	if len(taskTemplateList) > 0 {
		// 构造返回结果
		result = make([]*models.TaskTemplateDto, len(taskTemplateList))
		for i, taskTemplate := range taskTemplateList {
			if result[i], err = s.genTaskTemplateDto(taskTemplate.Id); err != nil {
				return
			}
		}
	}
	return
}

func (s *TaskTemplateService) GetCheckRoleAndHandler(requestTemplateId string) (role, handler string) {
	// 先查询 定版配置
	var taskTemplateList []*models.TaskTemplateTable
	var taskHandleTemplateList []*models.TaskHandleTemplateTable
	var requestTemplate *models.RequestTemplateTable
	var requestTemplateRoleList []*models.RequestTemplateRoleTable
	dao.X.SQL("select * from task_template where request_template = ? and type = ?", requestTemplateId, models.TaskTypeCheck).Find(&taskTemplateList)
	if len(taskTemplateList) > 0 {
		dao.X.SQL("select * from task_handle_template where task_template = ?", taskTemplateList[0].Id).Find(&taskHandleTemplateList)
		if len(taskHandleTemplateList) > 0 {
			role = taskHandleTemplateList[0].Role
			handler = taskHandleTemplateList[0].Handler
		}
	} else {
		// 没有找到定版模版,说明不需要定版,直接返回空
		return
	}
	// 定版没配置定版角色,查找模版属主
	if role == "" {
		requestTemplate, _ = GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
		if requestTemplate != nil {
			handler = requestTemplate.Handler
		}
		requestTemplateRoleList, _ = GetRequestTemplateService().GetRequestTemplateRole(requestTemplateId)
		if len(requestTemplateRoleList) > 0 {
			for _, requestTemplateRole := range requestTemplateRoleList {
				if requestTemplateRole.RoleType == string(models.RolePermissionMGMT) {
					role = requestTemplateRole.Role
					break
				}
			}
		}
	}
	return
}
