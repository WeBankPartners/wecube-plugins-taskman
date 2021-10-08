package request

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/services/db"
	"github.com/gin-gonic/gin"
)

func GetEntityData(c *gin.Context) {
	id := c.Query("requestTemplateId")
	if id == "" {
		middleware.ReturnParamValidateError(c, fmt.Errorf("Param requestTemplateId can not empty "))
		return
	}
	result, err := db.GetEntityData(id)
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
	result, err := db.ProcessDataPreview(requestTemplateId, entityDataId)
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
	err := db.CreateRequest(&param, middleware.GetRequestUser(c), middleware.GetRequestRoles(c))
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
