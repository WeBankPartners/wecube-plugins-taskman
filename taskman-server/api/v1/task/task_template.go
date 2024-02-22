package task

import (
	"errors"
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
	err = service.GetRequestTemplateService().UpdateRequestTemplateStatusToCreated(id, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}

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

// 新建自定义任务模板
func CreateCustomTaskTemplate(c *gin.Context) {
	var param models.CustomTaskTemplateCreateParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	result, err := service.GetTaskTemplateService().CreateCustomTaskTemplate(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

// 更新任务模板
func UpdateCustomTaskTemplate(c *gin.Context) {
	var param models.TaskTemplateDto
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := service.GetTaskTemplateService().UpdateCustomTaskTemplate(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

// 删除自定义任务模板
func DeleteCustomTaskTemplate(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	result, err := service.GetTaskTemplateService().DeleteCustomTaskTemplate(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

// 读取自定义任务模板
func GetCustomTaskTemplate(c *gin.Context) {
	requestTemplateId := c.Param("requestTemplateId")
	id := c.Param("id")
	result, err := service.GetTaskTemplateService().GetCustomTaskTemplate(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if result.RequestTemplateId != requestTemplateId {
		err = errors.New("param id and requestTemplateId not match")
		middleware.ReturnParamValidateError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

// 自定义任务模板id列表
func ListCustomTaskTemplateIds(c *gin.Context) {
	requestTemplateId := c.Param("requestTemplateId")
	result, err := service.GetTaskTemplateService().ListCustomTaskTemplateIds(requestTemplateId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

// 任务模板列表
func ListCustomTaskTemplate(c *gin.Context) {
	requestTemplateId := c.Param("requestTemplateId")
	result, err := service.GetTaskTemplateService().ListCustomTaskTemplate(requestTemplateId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}
