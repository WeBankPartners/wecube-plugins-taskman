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
	entityDataId := c.Query("entityDataId")
	result, err := db.GetRequestPreData(requestId, entityDataId, c.GetHeader("Authorization"))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func ListUserRequest(c *gin.Context) {
	result, err := db.ListUserRequest(middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
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

func SaveRequest(c *gin.Context) {
	//requestId := c.Param("requestId")
	//var param models.RequestCacheData
	//if err := c.ShouldBindJSON(&param);err != nil {
	//	middleware.ReturnParamValidateError(c, err)
	//	return
	//}

}
