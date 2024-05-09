package form

import (
	"fmt"
	"strings"

	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/service"
	"github.com/gin-gonic/gin"
)

// AddFormTemplateLibrary 添加组件库
func AddFormTemplateLibrary(c *gin.Context) {
	var param models.FormTemplateLibraryParam
	var nameRepeat bool
	var err error
	if err = c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if strings.TrimSpace(param.Name) == "" {
		middleware.ReturnParamEmptyError(c, "name")
		return
	}
	// name 重名校验
	if nameRepeat, err = service.GetFormTemplateLibraryService().CheckNameExist(param.Name); err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if nameRepeat {
		middleware.ReturnServerHandleError(c, fmt.Errorf("name %s repeat", param.Name))
		return
	}
	if err = service.GetFormTemplateLibraryService().AddFormTemplateLibrary(param, middleware.GetRequestUser(c)); err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func DeleteFormTemplateLibrary(c *gin.Context) {
	var err error
	var formTemplateLibraryTable *models.FormTemplateLibraryTable
	id := c.Query("id")
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	if formTemplateLibraryTable, err = service.GetFormTemplateLibraryService().Get(id); err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	// 只有表单库创建用户才有该条记录删除权限
	if formTemplateLibraryTable.CreatedBy != middleware.GetRequestUser(c) {
		middleware.ReturnServerHandleError(c, fmt.Errorf("no deletion permission"))
		return
	}

	if err = service.GetFormTemplateLibraryService().DeleteFormTemplateLibrary(id); err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

// QueryFormTemplateLibrary 查询组件库
func QueryFormTemplateLibrary(c *gin.Context) {
	var param models.QueryFormTemplateLibraryParam
	var list []*models.FormTemplateLibraryDto
	var pageInfo models.PageInfo
	var err error
	if err = c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	if pageInfo, list, err = service.GetFormTemplateLibraryService().QueryFormTemplateLibrary(param); err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnPageData(c, pageInfo, list)
}

func QueryAllFormTemplateLibraryFormType(c *gin.Context) {
	var formTypes []string
	var err error
	if formTypes, err = service.GetFormTemplateLibraryService().QueryAllFormType(); err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, formTypes)
}
