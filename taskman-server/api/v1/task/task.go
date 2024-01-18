package task

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/services/db"
	"github.com/gin-gonic/gin"
	"io/ioutil"
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
		output, taskId, tmpErr := db.PluginTaskCreate(input, param.RequestId, param.DueDate, param.AllowedOptions)
		if tmpErr != nil {
			output.ErrorCode = "1"
			output.ErrorMessage = tmpErr.Error()
			err = tmpErr
		} else {
			notifyErr := db.NotifyTaskMail(taskId)
			if notifyErr != nil {
				log.Logger.Error("Notify task mail fail", log.Error(notifyErr))
			}
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
	var param models.TaskApproveParam
	var task models.TaskTable
	var operator = middleware.GetRequestUser(c)
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	var err error
	for _, v := range param.FormData {
		tmpErr := validateFormRequire(v)
		if tmpErr != nil {
			err = tmpErr
			break
		}
	}
	if err == nil {
		err = db.ValidateRequestForm(param.FormData, c.GetHeader("Authorization"))
	}
	if err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	task, err = db.GetSimpleTask(taskId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if operator != task.Handler {
		middleware.ReturnTaskSaveNotPermissionError(c)
		return
	}
	err = db.SaveTaskForm(taskId, operator, param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		db.RecordTaskLog(taskId, task.Name, operator, "saveTask", c.Request.RequestURI, c.GetString("requestBody"))
		middleware.ReturnSuccess(c)
	}
}

func validateFormRequire(param *models.RequestPreDataTableObj) error {
	var err error
	requireMap := make(map[string]int)
	for _, v := range param.Title {
		if v.Required == "yes" {
			requireMap[v.Name] = 1
		}
	}
	for _, v := range param.Value {
		for dataKey, dataValue := range v.EntityData {
			if _, b := requireMap[dataKey]; b {
				if dataValue == nil {
					err = fmt.Errorf("Form:%s:%s data:%s can not empty ", v.PackageName, v.EntityName, dataKey)
				} else {
					if fmt.Sprintf("%s", dataValue) == "" {
						err = fmt.Errorf("Form:%s:%s data:%s can not empty ", v.PackageName, v.EntityName, dataKey)
					}
				}
			}
			if err != nil {
				break
			}
		}
		if err != nil {
			break
		}
	}
	return err
}

func ApproveTask(c *gin.Context) {
	taskId := c.Param("taskId")
	var param models.TaskApproveParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	var err error
	var operator = middleware.GetRequestUser(c)
	for _, v := range param.FormData {
		tmpErr := validateFormRequire(v)
		if tmpErr != nil {
			err = tmpErr
			break
		}
	}
	if err == nil {
		err = db.ValidateRequestForm(param.FormData, c.GetHeader("Authorization"))
	}
	if err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	taskTable, err := db.GetSimpleTask(taskId)
	if err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if taskTable.Handler != operator {
		middleware.ReturnTaskApproveNotPermissionError(c)
		return
	}
	err = db.ApproveTask(taskId, operator, c.GetHeader("Authorization"), param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		db.RecordTaskLog(taskId, taskTable.Name, operator, "approveTask", c.Request.RequestURI, c.GetString("requestBody"))
		middleware.ReturnSuccess(c)
	}
}

func ChangeTaskStatus(c *gin.Context) {
	taskId := c.Param("taskId")
	operation := c.Param("operation")
	lastedUpdateTime := c.Param("latestUpdateTime")
	if operation != "mark" && operation != "start" && operation != "quit" && operation != "give" {
		middleware.ReturnChangeTaskStatusError(c)
		return
	}
	taskObj, err := db.ChangeTaskStatus(taskId, middleware.GetRequestUser(c), operation, lastedUpdateTime)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	db.RecordTaskLog(taskId, "", middleware.GetRequestUser(c), "changeTaskStatus", c.Request.RequestURI, operation)
	middleware.ReturnData(c, taskObj)
}

func UploadTaskAttachFile(c *gin.Context) {
	taskId := c.Param("taskId")
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseErrorJson{StatusCode: "PARAM_HANDLE_ERROR", StatusMessage: "Http read upload file fail:" + err.Error(), Data: nil})
		return
	}
	if file.Size > models.UploadFileMaxSize {
		middleware.ReturnUploadFileTooLargeError(c)
		return
	}
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseErrorJson{StatusCode: "PARAM_HANDLE_ERROR", StatusMessage: "File open error:" + err.Error(), Data: nil})
		return
	}
	b, err := ioutil.ReadAll(f)
	defer f.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseErrorJson{StatusCode: "PARAM_HANDLE_ERROR", StatusMessage: "Read content fail error:" + err.Error(), Data: nil})
		return
	}
	err = db.UploadAttachFile("", taskId, file.Filename, middleware.GetRequestUser(c), b)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		db.RecordTaskLog(taskId, "", middleware.GetRequestUser(c), "uploadTaskFile", c.Request.RequestURI, file.Filename)
		middleware.ReturnData(c, db.GetTaskAttachFileList(taskId))
	}
}
