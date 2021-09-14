package task

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/services/db"
	"github.com/gin-gonic/gin"
)

func GetTaskTemplate(c *gin.Context) {
	requestTemplateId := c.Param("requestTemplateId")
	proNodeId := c.Param("proNodeId")
	result, err := db.GetTaskTemplate(requestTemplateId, proNodeId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func UpdateTaskTemplate(c *gin.Context) {
	id := c.Param("requestTemplateId")
	var param models.TaskTemplateDto
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	var err error
	param.UpdatedBy = middleware.GetRequestUser(c)
	if param.Id != "" {
		err = db.UpdateTaskTemplate(param)
	} else {
		err = db.CreateTaskTemplate(param, id)
	}
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnSuccess(c)
	}
}
