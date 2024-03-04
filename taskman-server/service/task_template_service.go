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
	if requestTemplate.ProcDefId == "" {
		return nil, fmt.Errorf("param requestTemplate not type proc")
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
	for _, node := range nodeList {
		if node.NodeDefId == param.NodeDefId {
			sort = node.OrderedNum
			if node.NodeId != param.NodeId || node.NodeName != param.NodeDefName {
				return nil, fmt.Errorf("param nodeId %q or nodeName %q wrong", param.NodeId, param.NodeDefName)
			}
			break
		}
	}
	if param.Sort <= 0 || sort != param.Sort {
		return nil, fmt.Errorf("param sort %d or nodeDefId %s wrong", param.Sort, param.NodeDefId)
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

func (s *TaskTemplateService) getTaskTemplateHandler(requestTemplate string) (taskTemplateMap map[string]*models.TaskTemplateDto, err error) {
	taskTemplateMap = make(map[string]*models.TaskTemplateDto)
	return
}
