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

// GetFormTemplate 获取表单模板
func GetFormTemplate(c *gin.Context) {
	var result *models.SimpleFormTemplateDto
	var err error
	formTemplateId := c.Param("from-template-id")
	if formTemplateId == "" {
		middleware.ReturnParamEmptyError(c, "from-template-id")
		return
	}
	result, err = service.GetFormTemplateService().GetFormTemplate(formTemplateId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

// GetGlobalFormEntity 获取全局表单 entity
func GetGlobalFormEntity(c *gin.Context) {
	var entityList []string
	var err error
	requestTemplateId := c.Param("id")
	if requestTemplateId == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	entityList, err = service.GetFormTemplateService().GetDataFormTemplateItemGroups(requestTemplateId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, entityList)
}

// GetDataFormTemplateItemGroupConfig 获取配置表单组配置
func GetDataFormTemplateItemGroupConfig(c *gin.Context) {
	var configureDto *models.FormTemplateGroupConfigureDto
	var err error
	formTemplateId := c.Query("form-template-id")
	entity := c.Query("entity")
	formType := c.Query("form-type")
	itemGroupId := c.Query("item-group-id")
	requestTemplateId := c.Query("request-template-id")
	if requestTemplateId == "" || formTemplateId == "" {
		middleware.ReturnParamEmptyError(c, "request-template-id or item-group-id or form-template-id")
		return
	}
	configureDto, err = service.GetFormTemplateService().GetDataFormConfig(formTemplateId, itemGroupId, formType, entity, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, configureDto)
}

// GetFormTemplateItemGroupConfig 获取配置表单组配置
func GetFormTemplateItemGroupConfig(c *gin.Context) {
	var configureDto *models.FormTemplateGroupConfigureDto
	var err error
	formTemplateId := c.Query("form-template-id")
	itemGroupId := c.Query("item-group-id")
	requestTemplateId := c.Query("request-template-id")
	if itemGroupId == "" || formTemplateId == "" || requestTemplateId == "" {
		middleware.ReturnParamEmptyError(c, "request-template-id or item-group-name or form-template-id")
		return
	}
	configureDto, err = service.GetFormTemplateService().GetFormConfig(requestTemplateId, formTemplateId, itemGroupId, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, configureDto)
}

// UpdateFormTemplateItemGroupConfig 新增更新表单组
func UpdateFormTemplateItemGroupConfig(c *gin.Context) {
	var param models.FormTemplateGroupConfigureDto
	var err error
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if err = validateFormTemplateItemGroupConfigParam(param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	// 校验是否有修改权限
	err = service.GetRequestTemplateService().CheckPermission(param.RequestTemplateId, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	err = service.GetFormItemTemplateService().UpdateFormTemplateItemGroupConfig(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

// UpdateFormTemplateItemGroup 更新表单组
func UpdateFormTemplateItemGroup(c *gin.Context) {
	var param models.FormTemplateGroupCustomDataDto
	var err error
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	// 校验是否有修改权限
	err = service.GetRequestTemplateService().CheckPermission(param.RequestTemplateId, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	err = service.GetFormItemTemplateService().UpdateFormTemplateItemGroup(param)
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
	// 校验是否有修改权限
	err = service.GetRequestTemplateService().CheckPermission(param.RequestTemplateId, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
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
	itemGroupId := c.Query("item-group-id")
	if formTemplateId == "" || itemGroupId == "" {
		middleware.ReturnParamEmptyError(c, "form-template-id or item-group-id")
		return
	}
	err := service.GetFormItemTemplateService().DeleteFormTemplateItemGroup(formTemplateId, itemGroupId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

// CopyDataFormTemplateItemGroup 数据表单模板组copy
func CopyDataFormTemplateItemGroup(c *gin.Context) {
	formTemplateId := c.Query("form-template-id")
	itemGroupId := c.Query("item-group-id")
	if formTemplateId == "" || itemGroupId == "" {
		middleware.ReturnParamEmptyError(c, "form-template-id or item-group-name")
		return
	}
	err := service.GetFormItemTemplateService().CopyDataFormTemplateItemGroup(formTemplateId, itemGroupId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func validateFormTemplateItemGroupConfigParam(param models.FormTemplateGroupConfigureDto) error {
	if param.RequestTemplateId == "" {
		return fmt.Errorf("param RequestTemplateId is empty")
	}
	if param.FormTemplateId == "" {
		return fmt.Errorf("param FormTemplateId is empty")
	}
	if param.ItemGroupType == "" || param.ItemGroupName == "" || param.ItemGroupRule == "" {
		return fmt.Errorf("param ItemGroup is empty")
	}
	return nil
}
