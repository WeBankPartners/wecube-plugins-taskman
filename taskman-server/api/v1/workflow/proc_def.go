package workflow

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/service"
	"github.com/gin-gonic/gin"
)

// GetProcessDefinitions 流程定义
func GetProcessDefinitions(c *gin.Context) {
	var rowData *models.DefinitionsData
	var err error
	templateId := c.Param("templateId")
	if templateId == "" {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	rowData, err = service.GetProcessDefinitions(templateId, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, rowData)
}

// GetProcessInstance 获取请求工作流
func GetProcessInstance(c *gin.Context) {
	var rowData *models.ProcessInstance
	var err error
	instanceId := c.Param("instanceId")
	if instanceId == "" {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	rowData, err = service.GetProcDefService().GetProcessDefineInstance(instanceId, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, rowData)
}

// GetProcDefTaskNodeContext  获取工作流执行节点
func GetProcDefTaskNodeContext(c *gin.Context) {
	procInstanceId := c.Param("procInstanceId")
	nodeInstanceId := c.Param("nodeInstanceId")
	rowData, err := service.GetProcDefService().GetProcDefTaskNodeContext(procInstanceId, nodeInstanceId, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, rowData)
}

func ProcessDataPreview(c *gin.Context) {
	requestTemplateId := c.Query("requestTemplateId")
	entityDataId := c.Query("entityDataId")
	if requestTemplateId == "" || entityDataId == "" {
		middleware.ReturnParamValidateError(c, fmt.Errorf("Param requestTemplateId or entityDataId can not empty "))
		return
	}
	result, err := service.ProcessDataPreview(requestTemplateId, entityDataId, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func GetEntityData(c *gin.Context) {
	id := c.Query("requestId")
	if id == "" {
		middleware.ReturnParamValidateError(c, fmt.Errorf("Param requestId can not empty "))
		return
	}
	result, err := service.GetEntityData(id, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func GetCoreProcNodes(c *gin.Context) {
	var nodeList []*models.ProcNodeObj
	var err error
	requestTemplateId := c.Param("id")
	getType := c.Param("type")
	nodeList, err = service.GetProcDefService().GetProcessDefineTaskNodes(&models.RequestTemplateTable{Id: requestTemplateId}, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader), getType)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, nodeList)
}

func GetCoreProcessList(c *gin.Context) {
	mangeRole := c.Query("role")
	result, err := service.GetProcDefService().GetCoreProcessListNew(c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader), mangeRole)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}
