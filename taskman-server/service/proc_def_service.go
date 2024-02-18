package service

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/rpc"
	"sort"
)

// 编排任务节点
var taskNodeTypeMap = map[string]bool{
	string(models.ProcDefNodeTypeHuman):     true,
	string(models.ProcDefNodeTypeAutomatic): true,
	string(models.ProcDefNodeTypeData):      true,
}

type ProcDefService struct {
}

// GetCoreProcessListNew 获取用户角色下的编排列表
func (s ProcDefService) GetCoreProcessListNew(userToken, language, manageRole string) (processList []*models.ProcDefObj, err error) {
	var procDefDtoList []*models.ProcDefDto
	var nodesList []*models.DataModel
	var entityMap = make(map[string]models.ProcEntity)
	processList = make([]*models.ProcDefObj, 0)
	procDefDtoList, err = rpc.QueryProcessDefinitionList(userToken, language, manageRole, models.QueryProcessDefinitionParam{Plugins: []string{"taskman"}, Status: "deployed"})
	if err != nil {
		return
	}
	nodesList, err = rpc.QueryAllModels(userToken, language)
	if err != nil {
		return
	}
	entityMap = models.ConvertModelsList2Map(nodesList)
	if len(procDefDtoList) > 0 {
		for _, dto := range procDefDtoList {
			processList = append(processList, &models.ProcDefObj{
				ProcDefId:   dto.Id,
				ProcDefKey:  dto.Key,
				ProcDefName: dto.Name,
				Status:      dto.Status,
				RootEntity:  entityMap[dto.RootEntity],
				CreatedTime: dto.CreatedTime,
			})
		}
	}
	return
}

// GetCoreProcessListAll 查询当前插件的所有编排
func (s ProcDefService) GetCoreProcessListAll(userToken, language string) (processList []*models.ProcDefObj, err error) {
	var procDefList []*models.ProcDef
	var nodesList []*models.DataModel
	var entityMap = make(map[string]models.ProcEntity)
	processList = make([]*models.ProcDefObj, 0)
	procDefList, err = rpc.QueryAllProcessDefinitionList(userToken, language)
	if err != nil {
		return
	}
	nodesList, err = rpc.QueryAllModels(userToken, language)
	if err != nil {
		return
	}
	entityMap = models.ConvertModelsList2Map(nodesList)
	if len(procDefList) > 0 {
		for _, model := range procDefList {
			processList = append(processList, &models.ProcDefObj{
				ProcDefId:   model.Id,
				ProcDefKey:  model.Key,
				ProcDefName: model.Name,
				Status:      model.Status,
				RootEntity:  entityMap[model.RootEntity],
				CreatedTime: model.CreatedTime.Format(models.DateTimeFormat),
			})
		}
	}
	return
}

// GetProcessDefine 获取编排定义
func (s ProcDefService) GetProcessDefine(procDefId, userToken, language string) (result *models.DefinitionsData, err error) {
	result, err = rpc.GetProcessDefine(procDefId, userToken, language)
	if err != nil {
		return
	}
	return
}

// GetProcessDefineTaskNodes 获取编排节点
func (s ProcDefService) GetProcessDefineTaskNodes(requestTemplate *models.RequestTemplateTable, userToken, language, filterType string) (nodeList []*models.ProcNodeObj, err error) {
	var allTaskNodeList []*models.ProcNodeObj
	nodeList = make([]*models.ProcNodeObj, 0)
	if requestTemplate == nil {
		err = fmt.Errorf("requestTemplate is empty")
		return
	}
	if requestTemplate.ProcDefId == "" {
		requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(requestTemplate.Id)
		if err != nil {
			return
		}
	}
	if requestTemplate.ProcDefId == "" {
		err = fmt.Errorf("Request template proDefId illegal ")
		return
	}
	allTaskNodeList, err = rpc.GetProcessDefineTaskNodes(requestTemplate.ProcDefId, userToken, language)
	if err != nil {
		return
	}
	if len(allTaskNodeList) > 0 {
		for _, node := range allTaskNodeList {
			if _, ok := taskNodeTypeMap[node.NodeType]; !ok {
				continue
			}
			if filterType == "template" {
				// 模板只展示人工节点
				if node.NodeType != string(models.ProcDefNodeTypeHuman) {
					continue
				}
			} else if filterType == "bind" {
				if node.DynamicBind {
					continue
				}
			}
			nodeList = append(nodeList, node)
		}
		sort.Sort(models.ProcNodeObjList(nodeList))
	}
	return
}

// GetProcessDefineInstance 获取编排实例
func (s ProcDefService) GetProcessDefineInstance(instanceId, userToken, language string) (processInstance *models.ProcessInstance, err error) {
	return rpc.GetProcessInstance(userToken, language, instanceId)
}

// GetProcDefRootEntities 获取编排RootEntity
func (s ProcDefService) GetProcDefRootEntities(procDefId, userToken, language string) (list []*models.ProcDefEntityDataObj, err error) {
	return rpc.GetProcDefRootEntities(procDefId, userToken, language)
}

// GetProcDefDataPreview 查询编排设计 RootEntity
func (s ProcDefService) GetProcDefDataPreview(procDefId, entityDataId, userToken, language string) (result *models.EntityTreeData, err error) {
	return rpc.GetProcDefDataPreview(procDefId, entityDataId, userToken, language)
}

// GetProcDefTaskNodeContext 查询编排设计 RootEntity
func (s ProcDefService) GetProcDefTaskNodeContext(procInstanceId, taskNodeId, userToken, language string) (data interface{}, err error) {
	return rpc.GetProcDefTaskNodeContext(procInstanceId, taskNodeId, userToken, language)
}

// TerminationsProcDefInstance 终止编排实例运行
func (s ProcDefService) TerminationsProcDefInstance(param models.TerminateInstanceParam, userToken, language string) (err error) {
	return rpc.TerminationsProcDefInstance(param, userToken, language)
}

// StartProcDefInstances 启动编排实例
func (s ProcDefService) StartProcDefInstances(param models.RequestProcessData, userToken, language string) (result *models.StartInstanceResultData, err error) {
	return rpc.StartProcDefInstances(param, userToken, language)
}

func (s ProcDefService) CheckProDefId(proDefId, proDefName, proDefKey, userToken, language string) (exist bool, newProDefId string, err error) {
	exist = false
	var processList []*models.ProcDefObj
	if proDefKey != "" {
		tmpProcessList, tmpErr := GetProcDefService().GetCoreProcessListNew(userToken, language, "")
		if tmpErr != nil {
			err = tmpErr
		} else {
			for _, v := range tmpProcessList {
				processList = append(processList, &models.ProcDefObj{ProcDefId: v.ProcDefId, ProcDefName: v.ProcDefName, ProcDefKey: v.ProcDefKey})
			}
		}
	} else {
		allProcessList, tmpErr := GetProcDefService().GetCoreProcessListAll(userToken, language)
		if tmpErr != nil {
			err = tmpErr
		} else {
			for _, v := range allProcessList {
				processList = append(processList, &models.ProcDefObj{ProcDefId: v.ProcDefId, ProcDefName: v.ProcDefName, ProcDefKey: v.ProcDefKey})
			}
		}
	}
	if err != nil {
		return
	}
	for _, v := range processList {
		if v.ProcDefId == proDefId {
			exist = true
			break
		}
	}
	if exist {
		return
	}
	count := 0
	for _, v := range processList {
		if proDefKey != "" {
			if proDefKey == v.ProcDefKey {
				count = count + 1
				newProDefId = v.ProcDefId
			}
			continue
		}
		if v.ProcDefName == proDefName {
			count = count + 1
			newProDefId = v.ProcDefId
		}
	}
	if count != 1 {
		if proDefKey != "" {
			err = fmt.Errorf("Find %d record from process list by query proDefKey:%s ", count, proDefKey)
		} else {
			err = fmt.Errorf("Find %d record from process list by query proDefName:%s ", count, proDefName)
		}
	}
	return
}
