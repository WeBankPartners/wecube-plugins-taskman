package template

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/service"
	"github.com/gin-gonic/gin"
)

func QueryRequestTemplateGroup(c *gin.Context) {
	var param models.QueryRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	pageInfo, rowData, err := service.GetRequestTemplateGroupService().QueryRequestTemplateGroup(&param, middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnPageData(c, pageInfo, rowData)
}

func CreateRequestTemplateGroup(c *gin.Context) {
	var param models.RequestTemplateGroupTable
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := service.GetRequestTemplateGroupService().CreateRequestTemplateGroup(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func UpdateRequestTemplateGroup(c *gin.Context) {
	var param models.RequestTemplateGroupTable
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if param.Id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	err := service.GetRequestTemplateGroupService().CheckRequestTemplateGroupRoles(param.Id, middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnDataPermissionError(c, err)
		return
	}
	err = service.GetRequestTemplateGroupService().UpdateRequestTemplateGroup(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func DeleteRequestTemplateGroup(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	err := service.GetRequestTemplateGroupService().CheckRequestTemplateGroupRoles(id, middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnDataPermissionError(c, err)
		return
	}
	err = service.GetRequestTemplateGroupService().DeleteRequestTemplateGroup(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}
