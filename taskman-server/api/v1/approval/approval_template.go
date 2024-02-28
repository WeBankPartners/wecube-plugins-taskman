package approval

import (
	"errors"

	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/service"
	"github.com/gin-gonic/gin"
)

// 新建审批模板
func CreateApprovalTemplate(c *gin.Context) {
	var param models.ApprovalTemplateCreateParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	result, err := service.GetApprovalTemplateService().CreateApprovalTemplate(&param, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

// 更新审批模板
func UpdateApprovalTemplate(c *gin.Context) {
	var param models.ApprovalTemplateDto
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := service.GetApprovalTemplateService().UpdateApprovalTemplate(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

// 删除审批模板
func DeleteApprovalTemplate(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	result, err := service.GetApprovalTemplateService().DeleteApprovalTemplate(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

// 读取审批模板
func GetApprovalTemplate(c *gin.Context) {
	requestTemplateId := c.Param("requestTemplateId")
	id := c.Param("id")
	result, err := service.GetApprovalTemplateService().GetApprovalTemplate(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if result.RequestTemplate != requestTemplateId {
		err = errors.New("param id and requestTemplateId not match")
		middleware.ReturnParamValidateError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

// 审批模板id列表
func ListApprovalTemplateIds(c *gin.Context) {
	requestTemplateId := c.Param("requestTemplateId")
	result, err := service.GetApprovalTemplateService().ListApprovalTemplateIds(requestTemplateId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

// 审批模板列表
func ListApprovalTemplates(c *gin.Context) {
	requestTemplateId := c.Param("requestTemplateId")
	result, err := service.GetApprovalTemplateService().ListApprovalTemplates(requestTemplateId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}
