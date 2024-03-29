package form

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/services/db"
	"github.com/gin-gonic/gin"
)

func GetRequestFormTemplate(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	result, err := db.GetRequestFormTemplate(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func UpdateRequestFormTemplate(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	var param models.FormTemplateDto
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	var err error
	param.UpdatedBy = middleware.GetRequestUser(c)
	if param.Id != "" {
		err = db.UpdateRequestFormTemplate(param)
	} else {
		err = db.CreateRequestFormTemplate(param, id)
	}
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	db.SetRequestTemplateToCreated(id, middleware.GetRequestUser(c))
	result, _ := db.GetRequestFormTemplate(id)
	middleware.ReturnData(c, result)
}

func ConfirmRequestFormTemplate(c *gin.Context) {
	id := c.Param("id")
	err := db.ConfirmRequestTemplate(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	db.RecordRequestTemplateLog(id, "", middleware.GetRequestUser(c), "confirmRequestTemplate", c.Request.RequestURI, "")
	middleware.ReturnSuccess(c)
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
