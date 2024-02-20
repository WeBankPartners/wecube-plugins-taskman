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
	err = service.GetRequestTemplateService().ConfirmRequestTemplate(requestTemplateId)
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
	if param.DataFormTemplateId != "" {
		err = service.GetFormTemplateService().UpdateDataFormTemplate(param, requestTemplateId)
	} else {
		err = service.GetFormTemplateService().CreateDataFormTemplate(param, requestTemplateId)
	}
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

// GetFormTemplateItemGroup 获取配置表单组
func GetFormTemplateItemGroup(c *gin.Context) {
	var configureDto *models.FormTemplateGroupConfigureDto
	var err error
	formTemplateId := c.Query("form-template-id")
	itemGroupName := c.Query("item-group-name")
	requestTemplateId := c.Query("request-template-id")
	if itemGroupName == "" || requestTemplateId == "" {
		middleware.ReturnParamEmptyError(c, "request-template-id or item-group-name")
		return
	}
	// formTemplateId 为空,查询数据表单模板是否为空,为空则新建(只有数据表单的formTemplateId会传递"",任务和审批表单不会)
	if formTemplateId == "" {
		formTemplateId, err = createDataFormTemplate(requestTemplateId)
		if err != nil {
			middleware.ReturnServerHandleError(c, err)
			return
		}
		if formTemplateId == "" {
			middleware.ReturnParamEmptyError(c, "form-template-id")
			return
		}
	}
	configureDto, err = service.GetFormTemplateService().GetConfigureForm(formTemplateId, itemGroupName, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, configureDto)
}

// UpdateFormTemplateItemGroup 新增更新表单组
func UpdateFormTemplateItemGroup(c *gin.Context) {
	var param models.FormTemplateGroupConfigureDto
	var err error
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err = service.GetFormItemTemplateService().UpdateFormTemplateItemGroup(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

// UpdateFormTemplateItemGroupCustomData 更新表单组自定义数据
func UpdateFormTemplateItemGroupCustomData(c *gin.Context) {
	var param models.FormTemplateGroupCustomDataDto
	var err error
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	//err = service.GetFormItemTemplateService().UpdateFormTemplateItemGroup(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

// SortFormTemplateItemGroup 表单组排序
func SortFormTemplateItemGroup(c *gin.Context) {
	var param models.FormTemplateGroupSortDto
	var err error
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err = service.GetFormItemTemplateService().SortFormTemplateItemGroup(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

// DeleteFormTemplateItemGroup 删除表单组
func DeleteFormTemplateItemGroup(c *gin.Context) {
	formTemplateId := c.Query("form-template-id")
	itemGroupName := c.Query("item-group-name")
	if formTemplateId == "" || itemGroupName == "" {
		middleware.ReturnParamEmptyError(c, "form-template-id or item-group-name")
		return
	}
	err := service.GetFormItemTemplateService().DeleteFormTemplateItemGroup(formTemplateId, itemGroupName)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

// createDataFormTemplate 创建数据表单
func createDataFormTemplate(requestTemplateId string) (formTemplateId string, err error) {
	var requestTemplate *models.RequestTemplateTable
	requestTemplate, err = service.GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if err != nil {
		return
	}
	if requestTemplate == nil {
		err = fmt.Errorf("param request-template-id is vailid")
		return
	}
	// 新建数据表单
	if requestTemplate.DataFormTemplate == "" {
		err = service.GetFormTemplateService().CreateDataFormTemplate(models.DataFormTemplateDto{}, requestTemplateId)
		if err != nil {
			return
		}
		requestTemplate, _ = service.GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
		formTemplateId = requestTemplate.DataFormTemplate
	} else {
		err = fmt.Errorf("form-template-id is empty")
	}
	return
}
