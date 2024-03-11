package request

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/exterror"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/try"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// GetRequestDetail 新版请求详情
func GetRequestDetail(c *gin.Context) {
	requestId := c.Param("requestId")
	result, err := service.GetRequestDetailV2(requestId, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func CreateRequest(c *gin.Context) {
	var param models.RequestTable
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if param.RequestTemplate == "" {
		middleware.ReturnParamValidateError(c, fmt.Errorf("Param requestTemplate can not empty "))
		return
	}
	if param.Role == "" {
		middleware.ReturnParamValidateError(c, fmt.Errorf("Param role can not empty "))
		return
	}
	template, err := service.GetRequestTemplateService().GetRequestTemplate(param.RequestTemplate)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	// 设置请求名称
	param.Name = fmt.Sprintf("%s-%s-%s", template.Name, template.OperatorObjType, time.Now().Format("060102150405"))
	// 设置请求类型
	param.Type = template.Type
	param.CreatedBy = middleware.GetRequestUser(c)
	param.ExpireDay = template.ExpireDay
	param.RequestTemplateName = template.Name
	d, _ := time.ParseDuration(fmt.Sprintf("%dh", 24*param.ExpireDay))
	param.ExpectTime = time.Now().Add(d).Format(models.DateTimeFormat)
	if template.Status != "confirm" {
		param.TemplateVersion = "beta"
	} else {
		param.TemplateVersion = template.Version
	}
	var items []*models.FormItemTemplateTable
	dao.X.SQL("select * from form_item_template where form_template in (select id from form_template where request_template=? and request_form_type = ?) order by sort", template.Id, models.RequestFormTypeMessage).Find(&items)
	param.CustomForm.Title = items
	if template.ProcDefId != "" {
		param.AssociationWorkflow = true
	}
	err = service.CreateRequest(&param, middleware.GetRequestRoles(c), c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.GetOperationLogService().RecordRequestLog(param.Id, param.Name, param.CreatedBy, "createRequest", c.Request.RequestURI, c.GetString("requestBody"))
	middleware.ReturnData(c, param)
}

func SaveRequestCache(c *gin.Context) {
	requestId := c.Param("requestId")
	cacheType := c.Param("cacheType")
	event := c.Param("event")
	user := middleware.GetRequestUser(c)
	if cacheType == "data" {
		var param models.RequestProDataV2Dto
		if err := c.ShouldBindJSON(&param); err != nil {
			middleware.ReturnParamValidateError(c, err)
			return
		}
		request, err := service.GetSimpleRequest(requestId)
		if err != nil {
			middleware.ReturnServerHandleError(c, err)
			return
		}
		if request.CreatedBy != user {
			middleware.ReturnReportRequestNotPermissionError(c)
			return
		}
		err = service.SaveRequestCacheV2(requestId, user, c.GetHeader("Authorization"), &param)
		if err != nil {
			middleware.ReturnServerHandleError(c, err)
		} else {
			middleware.ReturnData(c, param)
		}
	} else {
		var param models.RequestCacheData
		var operator = middleware.GetRequestUser(c)
		if err := c.ShouldBindJSON(&param); err != nil {
			middleware.ReturnParamValidateError(c, err)
			return
		}
		taskHandle, err := service.GetTaskHandleService().GetLatestRequestCheckTaskHandleByRequestId(requestId)
		if err != nil {
			middleware.ReturnError(c, err)
			return
		}
		if taskHandle == nil {
			middleware.ReturnParamValidateError(c, fmt.Errorf("requestId:%s not has check taskHandle", requestId))
			return
		}
		if taskHandle.Handler != operator {
			switch event {
			// 暂存
			case "save":
				middleware.ReturnSaveRequestNotPermissionError(c)
				return
				//退回 or 确认定版
			case "submit":
				middleware.ReturnSubmitRequestNotPermissionError(c)
				return
			}
		}
		err = service.SaveRequestBindCache(requestId, operator, &param)
		if err != nil {
			middleware.ReturnServerHandleError(c, err)
		} else {
			middleware.ReturnData(c, param)
		}
	}
}

// CheckRequest 确认定版
func CheckRequest(c *gin.Context) {
	requestId := c.Param("requestId")
	var param models.RequestCacheData
	var request models.RequestTable
	var operator = middleware.GetRequestUser(c)
	var err error
	var task *models.TaskTable
	if err = c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	request, err = service.GetSimpleRequest(requestId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	task, err = service.GetTaskService().GetLatestCheckTask(requestId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if task == nil {
		err = fmt.Errorf("requestId:%s not has check task")
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if task.Status == string(models.TaskStatusDone) {
		err = fmt.Errorf("taskId:%s status  is done", task.Id)
		middleware.ReturnServerHandleError(c, err)
		return
	}
	taskHandle, err := service.GetTaskHandleService().GetLatestRequestCheckTaskHandleByRequestId(requestId)
	if err != nil {
		middleware.ReturnError(c, err)
		return
	}
	if taskHandle == nil {
		middleware.ReturnParamValidateError(c, fmt.Errorf("requestId:%s not has check taskHandle", requestId))
		return
	}
	if taskHandle.Handler != operator {
		middleware.ReturnParamValidateError(c, fmt.Errorf("request handler not  permission!"))
		return
	}
	err = service.CheckRequest(request, task, operator, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader), param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.GetOperationLogService().RecordRequestLog(requestId, "", middleware.GetRequestUser(c), "checkRequest", c.Request.RequestURI, c.GetString("requestBody"))
	middleware.ReturnSuccess(c)
}

func GetRequestHistory(c *gin.Context) {
	defer try.ExceptionStack(func(e interface{}, err interface{}) {
		retErr := fmt.Errorf("%v", err)
		middleware.ReturnError(c, exterror.Catch(exterror.New().ServerHandleError, retErr))
		log.Logger.Error(e.(string))
	})

	requestId := c.Param("requestId")
	result, err := service.GetRequestHistory(c, requestId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func PluginCreateRequest(c *gin.Context) {
	response := models.PluginRequestCreateResp{ResultCode: "0", ResultMessage: "success", Results: models.PluginRequestCreateOutput{}}
	var err error
	defer func() {
		if err != nil {
			log.Logger.Error("Plugin request create handle fail", log.Error(err))
			response.ResultCode = "1"
			response.ResultMessage = err.Error()
		}
		bodyBytes, _ := json.Marshal(response)
		c.Set("responseBody", string(bodyBytes))
		c.JSON(http.StatusOK, response)
	}()
	var param models.PluginRequestCreateParam
	if err = c.ShouldBindJSON(&param); err != nil {
		return
	}
	if len(param.Inputs) == 0 {
		return
	}
	for _, input := range param.Inputs {
		output, tmpErr := handlePluginRequestCreate(input, param.RequestId)
		if tmpErr != nil {
			output.ErrorCode = "1"
			output.ErrorMessage = tmpErr.Error()
			err = tmpErr
		}
		response.Results.Outputs = append(response.Results.Outputs, output)
	}
}

func handlePluginRequestCreate(input *models.PluginRequestCreateParamObj, callRequestId string) (result *models.PluginRequestCreateOutputObj, err error) {

	return
}
