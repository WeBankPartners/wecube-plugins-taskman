package task

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/service"
	"github.com/gin-gonic/gin"
)

var (
	operationLogService = service.GetOperationLogService()
)

func GetTaskTemplate(c *gin.Context) {
	requestTemplateId := c.Param("requestTemplateId")
	proNodeId := c.Param("proNodeId")
	result, err := service.GetTaskTemplate(requestTemplateId, proNodeId, "")
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func UpdateTaskTemplate(c *gin.Context) {
	id := c.Param("requestTemplateId")
	var param models.TaskTemplateDto
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	var err error
	param.UpdatedBy = middleware.GetRequestUser(c)
	err = validateTaskTemplateParam(param)
	if err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if param.Id != "" {
		err = service.UpdateTaskTemplate(param)
	} else {
		err = service.CreateTaskTemplate(param, id)
	}
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.SetRequestTemplateToCreated(id, middleware.GetRequestUser(c))
	operationLogService.RecordRequestTemplateLog(id, "", middleware.GetRequestUser(c), "updateTaskTemplate", c.Request.RequestURI, c.GetString("requestBody"))
	result, _ := service.GetTaskTemplate(id, param.NodeDefId, "")
	middleware.ReturnData(c, result)
}

func validateTaskTemplateParam(param models.TaskTemplateDto) error {
	if param.Name == "" {
		return fmt.Errorf("Param name can not empty ")
	}
	if len(param.USERoles) == 0 {
		return fmt.Errorf("Param user roles can not empty ")
	}
	if param.ExpireDay <= 0 {
		return fmt.Errorf("Param expire day can not empty ")
	}
	if param.NodeDefId == "" {
		return fmt.Errorf("Param nodeDefId can not empty ")
	}
	return nil
}
