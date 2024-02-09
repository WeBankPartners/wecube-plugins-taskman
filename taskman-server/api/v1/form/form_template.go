package form

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/service"
	"github.com/gin-gonic/gin"
)

func GetRequestFormTemplate(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	result, err := service.GetFormTemplateService().GetRequestFormTemplate(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func UpdateRequestFormTemplate(c *gin.Context) {
	var param models.FormTemplateDto
	var err error
	var user = middleware.GetRequestUser(c)
	id := c.Param("id")
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	param.UpdatedBy = user
	if param.Id != "" {
		err = service.GetFormTemplateService().UpdateRequestFormTemplate(param, id)
	} else {
		err = service.GetFormTemplateService().CreateRequestFormTemplate(param, id)
	}
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.SetRequestTemplateToCreated(id, user)
	result, _ := service.GetFormTemplateService().GetRequestFormTemplate(id)
	middleware.ReturnData(c, result)
}

func ConfirmRequestFormTemplate(c *gin.Context) {
	var requestTemplate *models.RequestTemplateTable
	var err error
	id := c.Param("id")
	requestTemplate, err = service.GetRequestTemplateService().GetRequestTemplate(id)
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
	service.GetOperationLogService().RecordRequestTemplateLog(id, "", middleware.GetRequestUser(c), "confirmRequestTemplate", c.Request.RequestURI, "")
	middleware.ReturnSuccess(c)
}

func DeleteRequestFormTemplate(c *gin.Context) {
	id := c.Param("id")
	err := service.GetFormTemplateService().DeleteRequestFormTemplate(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnSuccess(c)
	}
}
