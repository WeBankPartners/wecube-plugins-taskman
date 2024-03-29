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
	pathQueryPreviewProcessDefinitionsEntity = "/platform/v1/public/process/definitions/%s/preview/entities/%s"
	// pathTerminationsProcessDefinitionsInstance  终止编排实例运行
	pathTerminationsProcessDefinitionsInstance = "/platform/v1/public/process/instances/%s/terminations"
	// pathStartProcDefInstance 启动编排实例
	pathStartProcDefInstance = "/platform/v1/public/process/instances"
	// pathQueryProcessDefinitionsInstanceTaskNodeContext 查询编排执行节点context
	pathQueryProcessDefinitionsInstanceTaskNodeContext = "/platform/v1/process/instances/%s/tasknodes/%s/context"

	// pathQueryModel 查询model
	pathQueryModel = "/platform/v1/models?withAttr=no"
	// pathQueryEntities 查询entity
	pathQueryEntities = "/platform/v1/data-model/dme/all-entities"
	// pathQueryEntityModel
	pathQueryEntityModel = "/platform/v1/models/package/%s/entity/%s"
	// pathSyncWorkflowUseRole
	pathSyncWorkflowUseRole = "/platform/v1/public/process/definitions/syncUseRole"
	// pathQueryEntityExpressionData 查询表达式数据
	pathQueryEntityExpressionData = "/platform/v1/public/data-model/dme/integrated-query"
	// pathQueryProcessDefinitionsNodeOptions 查询编排节点判断选项
	pathQueryProcessDefinitionsNodeOptions = "/platform/v1/public/process/definitions/%s/options/%s"
)

// QueryProcessDefinitionList 查询当前角色编排列表,并且根据属主角色进行过滤
func QueryProcessDefinitionList(userToken, language, manageRole string, param models.QueryProcessDefinitionParam) (processList []*models.ProcDefDto, err error) {
	var response models.QueryProcessDefinitionResponse
	procDefDistinctMap := make(map[string]int)
	processList = make([]*models.ProcDefDto, 0)
	postBytes, _ := json.Marshal(param)
	byteArr, err := HttpPost(models.Config.Wecube.BaseUrl+pathQueryProcessDefinitions, userToken, language, postBytes)
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
					if _, existFlag := procDefDistinctMap[dto.Id]; !existFlag {
						processList = append(processList, dto)
						procDefDistinctMap[dto.Id] = 1
					}
				}
			}
		}
	}
	return
}

// GetProcessDefineTaskNodes 获取编排节点
func GetProcessDefineTaskNodes(procDefId, userToken, language string) (list []*models.ProcNodeObj, err error) {
	var response models.ProcDefTaskNodesResponse
	list = make([]*models.ProcNodeObj, 0)
	byteArr, err := HttpGet(fmt.Sprintf(models.Config.Wecube.BaseUrl+pathQueryProcessDefinitionsTaskNodes, procDefId), userToken, language)
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
	byteArr, err := HttpGet(fmt.Sprintf(models.Config.Wecube.BaseUrl+pathAllQueryProcessDefinitions, "taskman"), userToken, language)
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
func QueryAllModels(userToken, language string) (modelList []*models.DataModel, err error) {
	var response models.QueryAllModelsResponse
	modelList = make([]*models.DataModel, 0)
	byteArr, err := HttpGet(models.Config.Wecube.BaseUrl+pathQueryModel, userToken, language)
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
		modelList = response.Data
	}
	return
}

// QueryEntityAttributes 查询所有模型
func QueryEntityAttributes(param models.QueryExpressionDataParam, userToken, language string) (entitiesList []*models.ExpressionEntities, err error) {
	var response models.QueryExpressionEntitiesResponse
	postBytes, _ := json.Marshal(param)
	entitiesList = []*models.ExpressionEntities{}
	byteArr, err := HttpPost(models.Config.Wecube.BaseUrl+pathQueryEntities, userToken, language, postBytes)
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
	byteArr, err := HttpGet(fmt.Sprintf(models.Config.Wecube.BaseUrl+pathQueryProcessDefinitionsInstance, instanceId), userToken, language)
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
	byteArr, err := HttpGet(fmt.Sprintf(models.Config.Wecube.BaseUrl+pathGetProcessDefinitions, procDefId), userToken, language)
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
	byteArr, err := HttpGet(fmt.Sprintf(models.Config.Wecube.BaseUrl+pathQueryProcessDefinitionsEntityData, procDefId), userToken, language)
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
	byteArr, err := HttpGet(fmt.Sprintf(models.Config.Wecube.BaseUrl+pathQueryPreviewProcessDefinitionsEntity, procDefId, entityDataId), userToken, language)
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
	byteArr, err = HttpGet(fmt.Sprintf(models.Config.Wecube.BaseUrl+pathQueryProcessDefinitionsInstanceTaskNodeContext, procInstanceId, taskNodeId), userToken, language)
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
	byteArr, err = HttpPost(fmt.Sprintf(models.Config.Wecube.BaseUrl+pathTerminationsProcessDefinitionsInstance, param.ProcInstId), userToken, language, postBytes)
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
	byteArr, err = HttpPost(models.Config.Wecube.BaseUrl+pathStartProcDefInstance, userToken, language, postBytes)
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

// GetEntityModel 查询entity
func GetEntityModel(packageName, entityName, userToken, language string) (data interface{}, err error) {
	var byteArr []byte
	var response models.DataModelEntityResponse
	byteArr, err = HttpGet(fmt.Sprintf(models.Config.Wecube.BaseUrl+pathQueryEntityModel, packageName, entityName), userToken, language)
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

// SyncWorkflowUseRole 同步编排使用角色
func SyncWorkflowUseRole(param models.SyncUseRoleParam, userToken, language string) (response models.HttpResponseMeta, err error) {
	var byteArr []byte
	postBytes, _ := json.Marshal(param)
	byteArr, err = HttpPost(models.Config.Wecube.BaseUrl+pathSyncWorkflowUseRole, userToken, language, postBytes)
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

func QueryEntityExpressionData(expression, rootDataId, userToken, language string) (result []map[string]interface{}, err error) {
	param := models.PluginQueryExpressionDataParam{
		DataModelExpression: expression,
		RootDataId:          rootDataId,
		Token:               userToken,
	}
	var response models.PluginQueryExpressionDataResponse
	var byteArr []byte
	postBytes, _ := json.Marshal(param)
	byteArr, err = HttpPost(models.Config.Wecube.BaseUrl+pathQueryEntityExpressionData, userToken, language, postBytes)
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
	} else {
		result = response.Data
	}
	return
}

func GetProcessNodeAllowOptions(procDefId, procNodeId, userToken, language string) (options []string, err error) {
	var response models.GetProcessNodeAllowOptionsResponse
	var byteArr []byte
	byteArr, err = HttpGet(models.Config.Wecube.BaseUrl+fmt.Sprintf(pathQueryProcessDefinitionsNodeOptions, procDefId, procNodeId), userToken, language)
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
	} else {
		options = response.Data
	}
	return
}
