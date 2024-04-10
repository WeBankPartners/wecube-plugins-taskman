package task

import (
	"errors"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/rpc"

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
	id := c.Param("id")
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
	if param.Id != id {
		middleware.ReturnParamValidateError(c, errors.New("param id wrong"))
		return
	}
	// 校验权限
	user := middleware.GetRequestUser(c)
	err := service.GetRequestTemplateService().CheckPermission(requestTemplateId, user)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	// 更新任务模板
	result, err := service.GetTaskTemplateService().UpdateTaskTemplate(&param, user)
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

func DeleteTaskTemplateFormTemplate(c *gin.Context) {
	var result *models.SimpleFormTemplateDto
	var err error
	requestTemplateId := c.Param("requestTemplate")
	id := c.Param("id")
	result, err = service.GetFormTemplateService().GetFormTemplate(requestTemplateId, id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if result != nil && len(result.Groups) > 0 {
		for _, group := range result.Groups {
			err = service.GetFormTemplateService().DeleteFormTemplateItemGroup(group.ItemGroupId)
			if err != nil {
				middleware.ReturnServerHandleError(c, err)
				return
			}
		}
	}
	middleware.ReturnSuccess(c)
}

// 读取任务模板
func GetTaskTemplate(c *gin.Context) {
	requestTemplateId := c.Param("requestTemplate")
	id := c.Param("id")
	typ := c.Query("type")
	// 校验参数
	if requestTemplateId == "" || id == "" || typ == "" {
		middleware.ReturnParamValidateError(c, errors.New("param empty"))
		return
	}
	result, err := service.GetTaskTemplateService().GetTaskTemplate(requestTemplateId, id, typ)
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
	result, err := service.GetTaskTemplateService().ListTaskTemplates(requestTemplateId, typ)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

// GetTaskTemplateWorkFlowOptions 获取任务模版的编排节点选项配置
func GetTaskTemplateWorkFlowOptions(c *gin.Context) {
	var taskTemplate *models.TaskTemplateTable
	var requestTemplate *models.RequestTemplateTable
	var options []string
	var err error
	taskTemplateId := c.Query("taskTemplateId")
	// 校验参数
	if taskTemplateId == "" {
		middleware.ReturnParamValidateError(c, errors.New("taskTemplateId param empty"))
		return
	}
	taskTemplate, err = service.GetTaskTemplateService().Get(taskTemplateId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if taskTemplate == nil {
		middleware.ReturnServerHandleError(c, fmt.Errorf("taskTemplateId invalid"))
		return
	}
	if taskTemplate.NodeId == "" {
		middleware.ReturnSuccess(c)
	}
	requestTemplate, err = service.GetRequestTemplateService().GetRequestTemplate(taskTemplate.RequestTemplate)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if requestTemplate == nil {
		middleware.ReturnServerHandleError(c, fmt.Errorf("taskTemplateId invalid"))
		return
	}
	options, err = rpc.GetProcessNodeAllowOptions(requestTemplate.ProcDefId, taskTemplate.NodeId, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, options)
}
