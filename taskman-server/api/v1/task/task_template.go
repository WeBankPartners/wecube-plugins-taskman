package task

import (
	"errors"

	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/service"
	"github.com/gin-gonic/gin"
)

// CreateTaskTemplate 新建任务模板
func CreateTaskTemplate(c *gin.Context) {
	requestTemplateId := c.Param("requestTemplate")
	var param models.TaskTemplateDto
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	// 校验参数
	if param.Type == "" || param.RequestTemplate == "" || param.Name == "" || param.ExpireDay <= 0 || param.Sort <= 0 {
		middleware.ReturnParamValidateError(c, errors.New("param empty"))
		return
	}
	if param.Id != "" || param.NodeDefId != "" {
		middleware.ReturnParamValidateError(c, errors.New("param not empty"))
		return
	}
	if param.RequestTemplate != requestTemplateId {
		middleware.ReturnParamValidateError(c, errors.New("param requestTemplate wrong"))
		return
	}
	if param.Type != string(models.TaskTypeApprove) && param.Type != string(models.TaskTypeImplement) {
		middleware.ReturnParamValidateError(c, errors.New("param type wrong"))
		return
	}
	// 校验权限
	user := middleware.GetRequestUser(c)
	err := service.GetRequestTemplateService().CheckPermission(requestTemplateId, user)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	result, err := service.GetTaskTemplateService().CreateTaskTemplate(&param, user)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

// UpdateTaskTemplate 更新任务模板/创建编排任务模板
func UpdateTaskTemplate(c *gin.Context) {
	requestTemplateId := c.Param("requestTemplate")
	id := c.Param("id") // id是taskTemplateId，或nodeDefId，前者为空时传后者
	var param models.TaskTemplateDto
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	// 校验参数
	if id == "" || param.Type == "" || param.RequestTemplate == "" || param.Name == "" || param.ExpireDay <= 0 || param.Sort <= 0 {
		middleware.ReturnParamValidateError(c, errors.New("param empty"))
		return
	}
	if param.RequestTemplate != requestTemplateId {
		middleware.ReturnParamValidateError(c, errors.New("param requestTemplate wrong"))
		return
	}
	// 校验权限
	user := middleware.GetRequestUser(c)
	err := service.GetRequestTemplateService().CheckPermission(requestTemplateId, user)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	var result *models.TaskTemplateDto
	if param.Id == "" {
		// 创建编排任务模板
		if param.NodeDefId != id {
			middleware.ReturnParamValidateError(c, errors.New("param id wrong"))
			return
		}
		result, err = service.GetTaskTemplateService().CreateProcTaskTemplate(&param, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader), user)
	} else {
		// 更新任务模板
		if param.Id != id && param.NodeDefId != id {
			middleware.ReturnParamValidateError(c, errors.New("param id wrong"))
			return
		}
		result, err = service.GetTaskTemplateService().UpdateTaskTemplate(&param, user)
	}
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

// 删除任务模板
func DeleteTaskTemplate(c *gin.Context) {
	requestTemplateId := c.Param("requestTemplate")
	id := c.Param("id")
	// 校验参数
	if requestTemplateId == "" || id == "" {
		middleware.ReturnParamValidateError(c, errors.New("param empty"))
		return
	}
	// 校验权限
	user := middleware.GetRequestUser(c)
	err := service.GetRequestTemplateService().CheckPermission(requestTemplateId, user)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	result, err := service.GetTaskTemplateService().DeleteTaskTemplate(requestTemplateId, id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

// 读取任务模板
func GetTaskTemplate(c *gin.Context) {
	requestTemplateId := c.Param("requestTemplate")
	id := c.Param("id") // id是taskTemplateId，或nodeDefId，后者非空时传后者
	typ := c.Query("type")
	// 校验参数
	if requestTemplateId == "" || id == "" || typ == "" {
		middleware.ReturnParamValidateError(c, errors.New("param empty"))
		return
	}
	result, err := service.GetTaskTemplateService().GetTaskTemplate(requestTemplateId, id, typ, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

// 任务模板id列表
func ListTaskTemplateIds(c *gin.Context) {
	requestTemplateId := c.Param("requestTemplate")
	typ := c.Query("type")
	// 校验参数
	if requestTemplateId == "" || typ == "" {
		middleware.ReturnParamValidateError(c, errors.New("param empty"))
		return
	}
	result, err := service.GetTaskTemplateService().ListTaskTemplateIds(requestTemplateId, typ, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

// 任务模板列表
func ListTaskTemplates(c *gin.Context) {
	requestTemplateId := c.Param("requestTemplate")
	typ := c.Query("type")
	// 校验参数
	if requestTemplateId == "" || typ == "" {
		middleware.ReturnParamValidateError(c, errors.New("param empty"))
		return
	}
	result, err := service.GetTaskTemplateService().ListTaskTemplates(requestTemplateId, typ, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}
