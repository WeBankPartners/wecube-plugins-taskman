package form

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/services/db"
	"github.com/gin-gonic/gin"
)

func GetRequestFormTemplate(c *gin.Context) {
	id := c.Param("id")
	result, err := db.GetRequestFormTemplate(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func CreateRequestFormTemplate(c *gin.Context) {
	var param models.RequestFormTemplateDto
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	id, err := db.CreateRequestFormTemplate(param, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, id)
	}
}

func UpdateRequestFormTemplate(c *gin.Context) {
	var param models.RequestFormTemplateDto
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := db.UpdateRequestFormTemplate(param, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnSuccess(c)
	}
}

func ConfirmRequestFormTemplate(c *gin.Context) {
	id := c.Param("id")
	err := db.ConfirmRequestFormTemplate(id, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnSuccess(c)
	}
}

func DeleteRequestFormTemplate(c *gin.Context) {
	id := c.Param("id")
	err := db.DeleteRequestFormTemplate(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnSuccess(c)
	}
}
