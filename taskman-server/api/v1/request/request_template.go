package request

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/services/db"
	"github.com/gin-gonic/gin"
	"strings"
)

func QueryRequestTemplateGroup(c *gin.Context) {
	var param models.QueryRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	pageInfo, rowData, err := db.QueryRequestTemplateGroup(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnPageData(c, pageInfo, rowData)
	}
}

func CreateRequestTemplateGroup(c *gin.Context) {
	var param models.RequestTemplateGroupTable
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := db.CreateRequestTemplateGroup(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnSuccess(c)
	}
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
	err := db.UpdateRequestTemplateGroup(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnSuccess(c)
	}
}

func DeleteRequestTemplateGroup(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	err := db.DeleteRequestTemplateGroup(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnSuccess(c)
	}
}

func GetCoreProcessList(c *gin.Context) {
	result, err := db.GetCoreProcessListNew(c.GetHeader("Authorization"))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func GetCoreProcNodes(c *gin.Context) {
	requestTemplateId := c.Param("id")
	getType := c.Param("type")
	result, err := db.GetProcessNodesByProc(requestTemplateId, c.GetHeader("Authorization"), getType)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func GetRoleList(c *gin.Context) {
	db.SyncCoreRole()
	result, err := db.GetRoleList([]string{})
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func GetUserByRoles(c *gin.Context) {
	roleString := c.Query("roles")
	result, err := db.QueryUserByRoles(strings.Split(roleString, ","), c.GetHeader("Authorization"))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func GetUserRoles(c *gin.Context) {
	db.SyncCoreRole()
	result, err := db.GetRoleList(middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func QueryRequestTemplate(c *gin.Context) {
	var param models.QueryRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	pageInfo, rowData, err := db.QueryRequestTemplate(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnPageData(c, pageInfo, rowData)
	}
}

func CreateRequestTemplate(c *gin.Context) {
	var param models.RequestTemplateUpdateParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if err := validateRequestTemplateParam(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	param.CreatedBy = middleware.GetRequestUser(c)
	result, err := db.CreateRequestTemplate(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func UpdateRequestTemplate(c *gin.Context) {
	var param models.RequestTemplateUpdateParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if err := validateRequestTemplateParam(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if param.Id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	param.UpdatedBy = middleware.GetRequestUser(c)
	result, err := db.UpdateRequestTemplate(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func validateRequestTemplateParam(param *models.RequestTemplateUpdateParam) error {
	if param.Name == "" {
		return fmt.Errorf("Param name can not empty ")
	}
	if param.Group == "" {
		return fmt.Errorf("Param group can not empty ")
	}
	if param.ProcDefKey == "" || param.ProcDefId == "" || param.ProcDefName == "" {
		return fmt.Errorf("Param procDef can not empty ")
	}
	if len(param.MGMTRoles) == 0 {
		return fmt.Errorf("Param mgmt can not empty ")
	}
	return nil
}

func DeleteRequestTemplate(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	err := db.DeleteRequestTemplate(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnSuccess(c)
	}
}

func ListRequestTemplateEntityAttrs(c *gin.Context) {
	id := c.Param("id")
	result, err := db.ListRequestTemplateEntityAttrs(id, c.GetHeader("Authorization"))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func GetRequestTemplateEntityAttrs(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	result, err := db.GetRequestTemplateEntityAttrs(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func UpdateRequestTemplateEntityAttrs(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	var param []*models.ProcEntityAttributeObj
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := db.UpdateRequestTemplateEntityAttrs(id, param, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		db.SetRequestTemplateToCreated(id, middleware.GetRequestUser(c))
		middleware.ReturnSuccess(c)
	}
}

func GetRequestTemplateByUser(c *gin.Context) {
	result, err := db.GetRequestTemplateByUser(middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}
