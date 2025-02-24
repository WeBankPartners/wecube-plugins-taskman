package form

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/exterror"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/try"
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
		middleware.ReturnError(c, exterror.New().FormTemplateLibraryAddNameRepeatError)
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
	if formTemplateLibraryTable == nil {
		middleware.ReturnParamValidateError(c, fmt.Errorf("id is invalid"))
		return
	}
	// 只有表单库创建用户才有该条记录删除权限
	if formTemplateLibraryTable.CreatedBy != middleware.GetRequestUser(c) {
		middleware.ReturnError(c, exterror.New().FormTemplateLibraryDeletePermissionError)
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

func ExportFormTemplateLibraryData(c *gin.Context) {
	var result []*models.FormTemplateLibraryTableData
	var err error
	defer try.ExceptionStack(func(e interface{}, err interface{}) {
		retErr := fmt.Errorf("%v", err)
		middleware.ReturnError(c, exterror.Catch(exterror.New().ServerHandleError, retErr))
		log.Error(nil, log.LOGGER_APP, e.(string))
	})
	if result, err = service.ExportFormTemplateLibrary(c); err != nil {
		middleware.ReturnError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

// 导出组件库
func ExportFormTemplateLibrary(c *gin.Context) {
	defer try.ExceptionStack(func(e interface{}, err interface{}) {
		retErr := fmt.Errorf("%v", err)
		middleware.ReturnError(c, exterror.Catch(exterror.New().ServerHandleError, retErr))
		log.Error(nil, log.LOGGER_APP, e.(string))
	})

	var err error
	retData, err := service.ExportFormTemplateLibrary(c)
	if err != nil {
		middleware.ReturnError(c, err)
	} else {
		fileName := "empty"
		if len(retData) > 0 {
			fileName = fmt.Sprintf("%s et al.%d", retData[0].Id, len(retData))
		}
		fileName = fmt.Sprintf("%s-%s.json", fileName, time.Now().Format("20060102150405"))

		retDataBytes, tmpErr := json.Marshal(retData)
		if tmpErr != nil {
			err = fmt.Errorf("marshal exportFormTemplateLibrary failed: %s", tmpErr.Error())
			middleware.ReturnError(c, err)
			return
		}
		c.Header("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fileName))
		c.Data(http.StatusOK, "application/octet-stream", retDataBytes)
	}
}

// 导入组件库
func ImportFormTemplateLibrary(c *gin.Context) {
	defer try.ExceptionStack(func(e interface{}, err interface{}) {
		retErr := fmt.Errorf("%v", err)
		middleware.ReturnError(c, exterror.Catch(exterror.New().ServerHandleError, retErr))
		log.Error(nil, log.LOGGER_APP, e.(string))
	})

	_, fileBytes, err := middleware.ReadFormFile(c, "file")
	if err != nil {
		middleware.ReturnError(c, err)
		return
	}

	var formTemplateLibraryData []*models.FormTemplateLibraryTableData
	if err = json.Unmarshal(fileBytes, &formTemplateLibraryData); err != nil {
		middleware.ReturnError(c, fmt.Errorf("json unmarshal form template library data failed: %s", err.Error()))
		return
	}

	err = service.ImportFormTemplateLibrary(c, formTemplateLibraryData)
	if err != nil {
		middleware.ReturnError(c, err)
	} else {
		middleware.ReturnSuccess(c)
	}
}
