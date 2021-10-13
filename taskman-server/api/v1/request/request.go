package request

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/services/db"
	"github.com/gin-gonic/gin"
)

func GetEntityData(c *gin.Context) {
	id := c.Query("requestId")
	if id == "" {
		middleware.ReturnParamValidateError(c, fmt.Errorf("Param requestId can not empty "))
		return
	}
	result, err := db.GetEntityData(id, c.GetHeader("Authorization"))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func ProcessDataPreview(c *gin.Context) {
	requestTemplateId := c.Query("requestTemplateId")
	entityDataId := c.Query("entityDataId")
	if requestTemplateId == "" || entityDataId == "" {
		middleware.ReturnParamValidateError(c, fmt.Errorf("Param requestTemplateId or entityDataId can not empty "))
		return
	}
	result, err := db.ProcessDataPreview(requestTemplateId, entityDataId, c.GetHeader("Authorization"))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func GetRequestPreviewData(c *gin.Context) {
	requestId := c.Query("requestId")
	entityDataId := c.Query("rootEntityId")
	result, err := db.GetRequestPreData(requestId, entityDataId, c.GetHeader("Authorization"))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, models.RequestPreDataDto{RootEntityId: entityDataId, Data: result})
	}
}

func ListRequest(c *gin.Context) {
	var param models.QueryRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	pageInfo, rowData, err := db.ListRequest(&param, middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnPageData(c, pageInfo, rowData)
	}
}

func GetRequest(c *gin.Context) {
	requestId := c.Param("requestId")
	result, err := db.GetRequest(requestId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func GetRequestRootForm(c *gin.Context) {
	requestId := c.Param("requestId")
	result, err := db.GetRequestRootForm(requestId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func CreateRequest(c *gin.Context) {
	var param models.RequestTable
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if param.Name == "" || param.RequestTemplate == "" {
		middleware.ReturnParamValidateError(c, fmt.Errorf("Param name and requestTemplate can not empty "))
		return
	}
	param.CreatedBy = middleware.GetRequestUser(c)
	err := db.CreateRequest(&param, middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, param)
	}
}

func UpdateRequest(c *gin.Context) {
	var param models.RequestTable
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if param.Id == "" || param.Name == "" {
		middleware.ReturnParamValidateError(c, fmt.Errorf("Param id and name can not empty "))
		return
	}
	param.UpdatedBy = middleware.GetRequestUser(c)
	err := db.UpdateRequest(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, param)
	}
}

func DeleteRequest(c *gin.Context) {
	requestId := c.Param("requestId")
	err := db.DeleteRequest(requestId, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnSuccess(c)
	}
}

func SaveRequestCache(c *gin.Context) {
	requestId := c.Param("requestId")
	var param models.RequestPreDataDto
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := db.SaveRequestCacheNew(requestId, middleware.GetRequestUser(c), &param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, param)
	}
}

func GetRequestCache(c *gin.Context) {
	requestId := c.Param("requestId")
	result, err := db.GetRequestCache(requestId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func StartRequest(c *gin.Context) {
	requestId := c.Param("requestId")
	var param models.RequestCacheData
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	instanceId, err := db.StartRequest(requestId, middleware.GetRequestUser(c), c.GetHeader("Authorization"), param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, instanceId)
	}
}

func TerminateRequest(c *gin.Context) {
	requestId := c.Param("requestId")
	err := db.RequestTermination(requestId, middleware.GetRequestUser(c), c.GetHeader("Authorization"))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnSuccess(c)
	}
}
