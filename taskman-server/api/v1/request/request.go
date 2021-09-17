package request

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
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
