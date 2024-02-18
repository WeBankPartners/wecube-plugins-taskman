package form

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/service"
	"github.com/gin-gonic/gin"
)

func GetRequestFormTemplate(c *gin.Context) {
	requestTemplateId := c.Param("id")
	if requestTemplateId == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	result, err := service.GetFormTemplateService().GetRequestFormTemplate(requestTemplateId)
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
	requestTemplateId := c.Param("id")
	if requestTemplateId == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	param.UpdatedBy = user
	if param.Id != "" {
		err = service.GetFormTemplateService().UpdateRequestFormTemplate(param, requestTemplateId)
	} else {
		err = service.GetFormTemplateService().CreateRequestFormTemplate(param, requestTemplateId)
	}
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	err = service.GetRequestTemplateService().UpdateRequestTemplateStatusToCreated(requestTemplateId, user)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	result, _ := service.GetFormTemplateService().GetRequestFormTemplate(requestTemplateId)
	middleware.ReturnData(c, result)
}

func ConfirmRequestFormTemplate(c *gin.Context) {
	var requestTemplate *models.RequestTemplateTable
	var err error
	requestTemplateId := c.Param("id")
	if requestTemplateId == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	requestTemplate, err = service.GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
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
	err = service.ConfirmRequestTemplate(requestTemplateId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.GetOperationLogService().RecordRequestTemplateLog(requestTemplateId, "", middleware.GetRequestUser(c), "confirmRequestTemplate", c.Request.RequestURI, "")
	middleware.ReturnSuccess(c)
}

// GetDataFormTemplate 获取数据表单模板
func GetDataFormTemplate(c *gin.Context) {
	var result *models.DataFormTemplateDto
	var err error
	requestTemplateId := c.Param("id")
	if requestTemplateId == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	result, err = service.GetFormTemplateService().GetDataFormTemplate(requestTemplateId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

// UpdateDataFormTemplate 新增更新全局表单模板
func UpdateDataFormTemplate(c *gin.Context) {
	var param models.DataFormTemplateDto
	var err error
	var user = middleware.GetRequestUser(c)
	requestTemplateId := c.Param("id")
	if requestTemplateId == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	param.UpdatedBy = user
	if param.FormTemplateId != "" {
		err = service.GetFormTemplateService().UpdateDataFormTemplate(param, requestTemplateId)
	} else {
		err = service.GetFormTemplateService().CreateDataFormTemplate(param, requestTemplateId)
	}
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
}

// GetConfigureForm 获取配置表单
func GetConfigureForm(c *gin.Context) {
	var configureDto *models.FormTemplateGroupConfigureDto
	var err error
	formTemplateId := c.Query("form-template-id")
	itemGroupName := c.Query("item-group-name")
	if formTemplateId == "" || itemGroupName == "" {
		middleware.ReturnParamEmptyError(c, "id or form-name")
		return
	}
	configureDto, err = service.GetFormTemplateService().GetConfigureForm(formTemplateId, itemGroupName)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, configureDto)
}
