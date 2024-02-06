package form

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/service"
	"github.com/gin-gonic/gin"
)

var (
	requestTemplateService = service.GetRequestTemplateService()
	operationLogService    = service.GetOperationLogService()
)

func GetRequestFormTemplate(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	result, err := service.GetRequestFormTemplate(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func UpdateRequestFormTemplate(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	var param models.FormTemplateDto
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	var err error
	param.UpdatedBy = middleware.GetRequestUser(c)
	if param.Id != "" {
		err = service.UpdateRequestFormTemplate(param)
	} else {
		err = service.CreateRequestFormTemplate(param, id)
	}
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.SetRequestTemplateToCreated(id, middleware.GetRequestUser(c))
	result, _ := service.GetRequestFormTemplate(id)
	middleware.ReturnData(c, result)
}

func ConfirmRequestFormTemplate(c *gin.Context) {
	var requestTemplate *models.RequestTemplateTable
	var err error
	id := c.Param("id")
	requestTemplate, err = requestTemplateService.GetRequestTemplate(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if requestTemplate == nil {
		middleware.ReturnError(c, fmt.Errorf("param id is vailid"))
		return
	}
	// 只要待确认状态下才能去发布
	if requestTemplate.Status != string(models.RequestTemplateStatusPending) {
		middleware.ReturnError(c, fmt.Errorf("illegal operation"))
		return
	}
	err = service.ConfirmRequestTemplate(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	operationLogService.RecordRequestTemplateLog(id, "", middleware.GetRequestUser(c), "confirmRequestTemplate", c.Request.RequestURI, "")
	middleware.ReturnSuccess(c)
}

func DeleteRequestFormTemplate(c *gin.Context) {
	id := c.Param("id")
	err := service.DeleteRequestFormTemplate(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnSuccess(c)
	}
}
