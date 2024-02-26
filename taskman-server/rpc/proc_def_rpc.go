package rpc

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"

	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
)

const (
	// pathQueryProcessDefinitions 查询编排
	pathQueryProcessDefinitions = "/platform/v1/process/definitions/list"
	// pathAllQueryProcessDefinitions 查询插件所有编排
	pathAllQueryProcessDefinitions = "/platform/v1/process/definitions/list/%s"
	// pathGetProcessDefinitions 查询编排详情
	pathGetProcessDefinitions = "/platform/v1/process/definitions/%s/outline"
	// pathQueryProcessDefinitionsTaskNodes 查询编排节点
	pathQueryProcessDefinitionsTaskNodes = "/platform/v1/public/process/definitions/%s/tasknodes"
	// pathQueryProcessDefinitionsInstance 编排实例
	pathQueryProcessDefinitionsInstance = "/platform/v1/process/instances/%s"
	// pathQueryProcessDefinitionsEntityData  查询编排实例entityData
	pathQueryProcessDefinitionsEntityData = "/platform/v1/process/definitions/%s/root-entities"
	// pathQueryPreviewProcessDefinitionsEntity 编排设计entityData预览
	pathQueryPreviewProcessDefinitionsEntity = "/platform/v1/process/definitions/%s/preview/entities/%s"
	// pathTerminationsProcessDefinitionsInstance  终止编排实例运行
	pathTerminationsProcessDefinitionsInstance = "/platform/v1/public/process/instances/%s/terminations"
	// pathStartProcDefInstance 启动编排实例
	pathStartProcDefInstance = "/platform/v1/public/process/instances"
	// pathQueryProcessDefinitionsInstanceTaskNodeContext 查询编排执行节点context
	pathQueryProcessDefinitionsInstanceTaskNodeContext = "/platform/v1/process/instances/%s/tasknodes/%s/context"

	// pathQueryModel 查询model
	pathQueryModel = "/platform/v1/models"
	// pathQueryEntities 查询entity
	pathQueryEntities = "/platform/v1/data-model/dme/all-entities"

	// @todo 后面调整回来 models.Config.Wecube.BaseUrl
	BaseUrl = "http://106.52.160.142:18080"
)

// QueryProcessDefinitionList 查询当前角色编排列表,并且根据属主角色进行过滤
func QueryProcessDefinitionList(userToken, language, manageRole string, param models.QueryProcessDefinitionParam) (processList []*models.ProcDefDto, err error) {
	var response models.QueryProcessDefinitionResponse
	var procDefMap = make(map[string]*models.ProcDefDto)
	processList = make([]*models.ProcDefDto, 0)
	postBytes, _ := json.Marshal(param)
	byteArr, err := HttpPost(BaseUrl+pathQueryProcessDefinitions, userToken, language, postBytes)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteArr, &response)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	if response.Status != "OK" {
		err = fmt.Errorf(response.Message)
		return
	}
	if len(response.Data) > 0 {
		for _, queryDto := range response.Data {
			if manageRole != "" && manageRole != queryDto.ManageRole {
				continue
			}
			if len(queryDto.ProcDefList) > 0 {
				for _, dto := range queryDto.ProcDefList {
					procDefMap[dto.Id] = dto
				}
			}
		}
	}
	if len(procDefMap) > 0 {
		if err != nil {
			return
		}
		for _, dto := range procDefMap {
			processList = append(processList, dto)
		}
	}
	return
}

// GetProcessDefineTaskNodes 获取编排节点
func GetProcessDefineTaskNodes(procDefId, userToken, language string) (list []*models.ProcNodeObj, err error) {
	var response models.ProcDefTaskNodesResponse
	list = make([]*models.ProcNodeObj, 0)
	byteArr, err := HttpGet(fmt.Sprintf(BaseUrl+pathQueryProcessDefinitionsTaskNodes, procDefId), userToken, language)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteArr, &response)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	if response.Status != "OK" {
		err = fmt.Errorf(response.Message)
		return
	}
	if len(response.Data) > 0 {
		list = response.Data
	}
	return
}

// QueryAllProcessDefinitionList 查询所有编排列表
func QueryAllProcessDefinitionList(userToken, language string) (processList []*models.ProcDef, err error) {
	var response models.QueryProcessAllDefinitionResponse
	processList = make([]*models.ProcDef, 0)
	byteArr, err := HttpGet(fmt.Sprintf(BaseUrl+pathAllQueryProcessDefinitions, "taskman"), userToken, language)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteArr, &response)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	if response.Status != "OK" {
		err = fmt.Errorf(response.Message)
		return
	}
	if len(response.Data) > 0 {
		processList = response.Data
	}
	return
}

// QueryAllModels 查询所有模型
func QueryAllModels(userToken, language string) (nodesList []*models.DataModel, err error) {
	var response models.QueryAllModelsResponse
	nodesList = make([]*models.DataModel, 0)
	byteArr, err := HttpGet(BaseUrl+pathQueryModel, userToken, language)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteArr, &response)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	if response.Status != "OK" {
		err = fmt.Errorf(response.Message)
	}
	if len(response.Data) > 0 {
		nodesList = response.Data
	}
	return
}

// QueryEntityAttributes 查询所有模型
func QueryEntityAttributes(param models.QueryExpressionDataParam, userToken, language string) (entitiesList []*models.ExpressionEntities, err error) {
	var response models.QueryExpressionEntitiesResponse
	postBytes, _ := json.Marshal(param)
	entitiesList = []*models.ExpressionEntities{}
	byteArr, err := HttpPost(BaseUrl+pathQueryEntities, userToken, language, postBytes)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteArr, &response)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	if response.Status != "OK" {
		err = fmt.Errorf(response.Message)
	}
	if len(response.Data) > 0 {
		entitiesList = response.Data
	}
	return
}

// GetProcessInstance 查询编排实例
func GetProcessInstance(userToken, language, instanceId string) (processInstance *models.ProcessInstance, err error) {
	var response models.ProcessInstanceResponse
	byteArr, err := HttpGet(fmt.Sprintf(BaseUrl+pathQueryProcessDefinitionsInstance, instanceId), userToken, language)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteArr, &response)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	if response.Status != "OK" {
		err = fmt.Errorf(response.Message)
	}
	processInstance = response.Data
	return
}

// GetProcessDefine 查询编排详情
func GetProcessDefine(procDefId, userToken, language string) (result *models.DefinitionsData, err error) {
	var response models.ProcessDefinitionsResponse
	byteArr, err := HttpGet(fmt.Sprintf(BaseUrl+pathGetProcessDefinitions, procDefId), userToken, language)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteArr, &response)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	if response.Status != "OK" {
		err = fmt.Errorf(response.Message)
	}
	result = response.Data
	return
}

// GetProcDefRootEntities 查询编排实例RootEntity
func GetProcDefRootEntities(procDefId, userToken, language string) (list []*models.ProcDefEntityDataObj, err error) {
	var response models.ProcDefRootEntityResponse
	byteArr, err := HttpGet(fmt.Sprintf(BaseUrl+pathQueryProcessDefinitionsEntityData, procDefId), userToken, language)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteArr, &response)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	if response.Status != "OK" {
		err = fmt.Errorf(response.Message)
	}
	list = response.Data
	return
}

// GetProcDefDataPreview 查询编排设计 RootEntity
func GetProcDefDataPreview(procDefId, entityDataId, userToken, language string) (result *models.EntityTreeData, err error) {
	var response models.EntityTreeResponse
	byteArr, err := HttpGet(fmt.Sprintf(BaseUrl+pathQueryPreviewProcessDefinitionsEntity, procDefId, entityDataId), userToken, language)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteArr, &response)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	if response.Status != "OK" {
		err = fmt.Errorf(response.Message)
	}
	result = response.Data
	return
}

// GetProcDefTaskNodeContext 查询编排节点context
func GetProcDefTaskNodeContext(procInstanceId, taskNodeId, userToken, language string) (data interface{}, err error) {
	var byteArr []byte
	var response models.ProcDefTaskNodeContextResponse
	byteArr, err = HttpGet(fmt.Sprintf(BaseUrl+pathQueryProcessDefinitionsInstanceTaskNodeContext, procInstanceId, taskNodeId), userToken, language)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteArr, &response)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	if response.Status != "OK" {
		err = fmt.Errorf(response.Message)
	}
	data = response.Data
	return
}

// TerminationsProcDefInstance 终止编排实例运行
func TerminationsProcDefInstance(param models.TerminateInstanceParam, userToken, language string) (err error) {
	var byteArr []byte
	var response models.StartInstanceResponse
	postBytes, _ := json.Marshal(param)
	byteArr, err = HttpPost(fmt.Sprintf(BaseUrl+pathTerminationsProcessDefinitionsInstance, param.ProcInstId), userToken, language, postBytes)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteArr, &response)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	if response.Status != "OK" {
		err = fmt.Errorf(response.Message)
	}
	return
}

// StartProcDefInstances 启动编排实例
func StartProcDefInstances(param models.RequestProcessData, userToken, language string) (result *models.StartInstanceResultData, err error) {
	var byteArr []byte
	var response models.StartInstanceResponse
	postBytes, _ := json.Marshal(param)
	log.Logger.Info("Start request", log.String("param", string(postBytes)))
	byteArr, err = HttpPost(BaseUrl+pathStartProcDefInstance, userToken, language, postBytes)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteArr, &response)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	if response.Status != "OK" {
		err = fmt.Errorf(response.Message)
	}
	result = response.Data
	return
}
