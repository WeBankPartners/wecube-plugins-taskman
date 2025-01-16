package template

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/exterror"
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
	pageInfo, rowData, err := service.GetRequestTemplateGroupService().QueryRequestTemplateGroup(&param, middleware.GetRequestRoles(c), c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnPageData(c, pageInfo, rowData)
}

func CreateRequestTemplateGroup(c *gin.Context) {
	var param models.RequestTemplateGroupTable
	var requestTemplateGroupList []*models.RequestTemplateGroupTable
	var err error
	if err = c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if requestTemplateGroupList, err = service.GetRequestTemplateGroupService().QueryRequestTemplateGroupByName(param.Name); err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if len(requestTemplateGroupList) > 0 {
		middleware.ReturnError(c, fmt.Errorf("the same name:%s already exists", param.Name))
		return
	}
	if err = service.GetRequestTemplateGroupService().CreateRequestTemplateGroup(&param); err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func UpdateRequestTemplateGroup(c *gin.Context) {
	var param models.RequestTemplateGroupTable
	var requestTemplateGroupList []*models.RequestTemplateGroupTable
	var err error
	if err = c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if param.Id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	if requestTemplateGroupList, err = service.GetRequestTemplateGroupService().QueryRequestTemplateGroupByName(param.Name); err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	for _, requestTemplateGroup := range requestTemplateGroupList {
		if requestTemplateGroup.Id != param.Id {
			middleware.ReturnError(c, fmt.Errorf("the same name:%s already exists", param.Name))
			return
		}
	}
	err = service.GetRequestTemplateGroupService().CheckRequestTemplateGroupRoles(param.Id, middleware.GetRequestRoles(c))
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
	var err error
	var list []*models.RequestTemplateTable
	id := c.Query("id")
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	if err = service.GetRequestTemplateGroupService().CheckRequestTemplateGroupRoles(id, middleware.GetRequestRoles(c)); err != nil {
		middleware.ReturnRequestTemplateUpdatePermissionError(c, err)
		return
	}
	if list, err = service.GetRequestTemplateService().QueryRequestTemplateListByRequestTemplateGroup(id); err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if len(list) > 0 {
		middleware.ReturnServerHandleError(c, exterror.New().TemplateGroupHasUseDeleteError)
		return
	}
	if err = service.GetRequestTemplateGroupService().DeleteRequestTemplateGroup(id); err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}
