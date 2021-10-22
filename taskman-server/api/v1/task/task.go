package task

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/services/db"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetTaskFormStruct(c *gin.Context) {
	procInstId := c.Query("procInstId")
	nodeDefId := c.Query("nodeDefId")
	result, err := db.GetTaskFormStruct(procInstId, nodeDefId)
	if err != nil {
		result.Status = "ERROR"
		result.Message = err.Error()
	}
	log.Logger.Info("task form struct", log.JsonObj("response", result))
	c.JSON(http.StatusOK, result)
}

func CreateTask(c *gin.Context) {
	response := models.PluginTaskCreateResp{ResultCode: "0", ResultMessage: "success", Results: models.PluginTaskCreateOutput{}}
	var err error
	defer func() {
		if err != nil {
			log.Logger.Error("Task create handle fail", log.Error(err))
			response.ResultCode = "1"
			response.ResultMessage = err.Error()
		}
		bodyBytes, _ := json.Marshal(response)
		c.Set("responseBody", string(bodyBytes))
		c.JSON(http.StatusOK, response)
	}()
	var param models.PluginTaskCreateRequest
	c.ShouldBindJSON(&param)
	if len(param.Inputs) == 0 {
		return
	}
	for _, input := range param.Inputs {
		output, tmpErr := db.PluginTaskCreate(input, param.RequestId, param.AllowedOptions)
		if tmpErr != nil {
			output.ErrorCode = "1"
			output.ErrorMessage = tmpErr.Error()
			err = tmpErr
		}
		response.Results.Outputs = append(response.Results.Outputs, output)
	}
}

func ListTask(c *gin.Context) {
	var param models.QueryRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	pageInfo, rowData, err := db.ListTask(&param, middleware.GetRequestRoles(c), middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnPageData(c, pageInfo, rowData)
	}
}

func GetTask(c *gin.Context) {
	taskId := c.Param("taskId")
	result, err := db.GetTask(taskId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func SaveTaskForm(c *gin.Context) {
	taskId := c.Param("taskId")
	var param []*models.RequestPreDataTableObj
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := db.SaveTaskForm(taskId, middleware.GetRequestUser(c), param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnSuccess(c)
	}
}

func ApproveTask(c *gin.Context) {
	taskId := c.Param("taskId")
	var param models.TaskApproveParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := db.ApproveTask(taskId, middleware.GetRequestUser(c), c.GetHeader("Authorization"), param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnSuccess(c)
	}
}

func ChangeTaskStatus(c *gin.Context) {
	taskId := c.Param("taskId")
	operation := c.Param("operation")
	if operation != "mark" && operation != "start" && operation != "quit" {
		middleware.ReturnParamValidateError(c, fmt.Errorf("operation illegal"))
		return
	}
	taskObj, err := db.ChangeTaskStatus(taskId, middleware.GetRequestUser(c), operation)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, taskObj)
	}
}
