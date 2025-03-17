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
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

// GetRequestDetail 新版请求详情
func GetRequestDetail(c *gin.Context) {
	requestId := c.Query("requestId")
	taskId := c.Query("taskId")
	result, err := service.GetRequestDetailV2(requestId, taskId, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
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
		middleware.ReturnParamValidateError(c, fmt.Errorf("param requestTemplate can not empty "))
		return
	}
	if param.Role == "" {
		middleware.ReturnParamValidateError(c, fmt.Errorf("param role can not empty "))
		return
	}
	param.CreatedBy = middleware.GetRequestUser(c)
	_, err := handleCreateRequest(&param, middleware.GetRequestRoles(c), c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.GetOperationLogService().RecordRequestLog(param.Id, param.Name, param.CreatedBy, "createRequest", c.Request.RequestURI, c.GetString("requestBody"))
	middleware.ReturnData(c, param)
}

func handleCreateRequest(param *models.RequestTable, roles []string, userToken, language string) (newRequest *models.RequestTable, err error) {
	template, getTemplateErr := service.GetRequestTemplateService().GetRequestTemplate(param.RequestTemplate)
	if getTemplateErr != nil {
		err = getTemplateErr
		return
	}
	// 设置请求名称
	param.Name = fmt.Sprintf("%s-%s-%s", template.Name, template.OperatorObjType, time.Now().Format("060102150405"))
	// 设置请求类型
	param.Type = template.Type
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
	param.CustomForm.Title = models.ConvertFormItemTemplateModelList2Dto(items, &models.FormTemplateTable{})
	if template.ProcDefId != "" {
		param.AssociationWorkflow = true
	}
	newRequest = param
	err = service.CreateRequest(param, roles, userToken, language)
	return
}

func SaveRequestCache(c *gin.Context) {
	requestId := c.Param("requestId")
	cacheType := c.Param("cacheType")
	event := c.Param("event")
	user := middleware.GetRequestUser(c)
	var request models.RequestTable
	var err error
	request, err = service.GetSimpleRequest(requestId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if cacheType == "data" {
		var param models.RequestProDataV2Dto
		if err := c.ShouldBindJSON(&param); err != nil {
			middleware.ReturnParamValidateError(c, err)
			return
		}
		if request.CreatedBy != user {
			middleware.ReturnReportRequestNotPermissionError(c)
			return
		}
		for _, entityData := range param.Data {
			if err = service.HandleSensitiveDataDecode(entityData); err != nil {
				middleware.ReturnServerHandleError(c, err)
				return
			}
			passwordAttrMap := make(map[string]bool)
			for _, title := range entityData.Title {
				cmdbAttrModel := models.EntityAttributeObj{}
				if strings.TrimSpace(title.CmdbAttr) == "" {
					continue
				}
				if err = json.Unmarshal([]byte(title.CmdbAttr), &cmdbAttrModel); err != nil {
					return
				}
				if cmdbAttrModel.InputType == string(models.FormItemElementTypePassword) {
					passwordAttrMap[title.Name] = true
				}
			}
			// 密码处理,web传递原密码,需要加密处理
			for _, entityItem := range entityData.Value {
				for key, value := range entityItem.EntityData {
					// 空数据不用加密
					if value == nil || fmt.Sprintf("%+v", value) == "" {
						continue
					}
					inputValue := fmt.Sprintf("%+v", value)
					if passwordAttrMap[key] && !strings.HasPrefix(strings.ToLower(inputValue), models.EncryptPasswordPrefix) && !strings.HasPrefix(strings.ToLower(inputValue), models.EncryptPasswordPrefixC) {
						if inputValue, err = service.AesEnPasswordByGuid("", models.Config.EncryptSeed, inputValue, service.DEFALT_CIPHER_C); err != nil {
							err = fmt.Errorf("try to encrypt password type column:%s value:%s fail,%s  ", key, inputValue, err.Error())
							return
						}
						entityItem.EntityData[models.ModifyPrefixConstant+key] = 1
						entityItem.EntityData[key] = inputValue
					}
				}
			}
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
		if request.Status == string(models.RequestStatusDraft) {
			middleware.ReturnError(c, exterror.New().RequestHandleError)
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
	if request.Status == string(models.RequestStatusDraft) {
		middleware.ReturnError(c, exterror.New().RequestHandleError)
		return
	}
	task, err = service.GetTaskService().GetLatestCheckTask(requestId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if task == nil {
		err = fmt.Errorf("requestId:%s not has check task", requestId)
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
		middleware.ReturnParamValidateError(c, fmt.Errorf("request handler not permission"))
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
		log.Error(nil, log.LOGGER_APP, e.(string))
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
	var users []string
	var exist bool
	defer func() {
		if err != nil {
			log.Error(nil, log.LOGGER_APP, "Plugin request create handle fail", zap.Error(err))
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
	requestToken := c.GetHeader("Authorization")
	requestLanguage := "en"
	for _, input := range param.Inputs {
		users, err = service.GetRoleService().QueryUserByRoles([]string{input.ReportRole},
			c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
		exist = false
		for _, user := range users {
			if user == input.ReportUser {
				exist = true
			}
		}
		if !exist {
			// 用户和角色填写不匹配,返回错误
			err = fmt.Errorf("role:%s and user:%s do not match", input.ReportRole, input.ReportUser)
			return
		}
		output, tmpErr := handlePluginRequestCreate(input, param.RequestId, requestToken, requestLanguage)
		if tmpErr != nil {
			output.ErrorCode = "1"
			output.ErrorMessage = tmpErr.Error()
			err = tmpErr
		}
		response.Results.Outputs = append(response.Results.Outputs, output)
	}
}

func handlePluginRequestCreate(input *models.PluginRequestCreateParamObj, callRequestId, token, language string) (result *models.PluginRequestCreateOutputObj, err error) {
	var newRequest *models.RequestTable
	result = &models.PluginRequestCreateOutputObj{CallbackParameter: input.CallbackParameter}
	requestObj := models.RequestTable{RequestTemplate: input.RequestTemplate, Role: input.ReportRole, CreatedBy: input.ReportUser}
	// 创建请求
	if newRequest, err = handleCreateRequest(&requestObj, []string{input.ReportRole}, token, language); err != nil {
		return
	}
	if requestObj.Id == "" {
		err = fmt.Errorf("create request fail,requestId empty")
		return
	}
	result.RequestId = requestObj.Id
	if newRequest != nil {
		result.RequestName = newRequest.Name
		result.RequestTemplate = newRequest.RequestTemplate
		result.RequestTemplateType = newRequest.Type
	}
	// 预览根数据表单
	previewData, platformPreviewData, previewErr := service.GetRequestPreData(requestObj.Id, input.RootDataId, token, language)
	if previewErr != nil {
		err = previewErr
		return
	}
	// 保存请求表单数据
	saveParam := models.RequestProDataV2Dto{Data: previewData, RootEntityId: input.RootDataId, Name: requestObj.Name, ExpectTime: requestObj.ExpectTime, Description: requestObj.Description, CustomForm: models.CustomForm{}, ApprovalList: []*models.TaskTemplateDto{}}
	for _, dataRow := range platformPreviewData.EntityTreeNodes {
		if dataRow.DataId == input.RootDataId {
			saveParam.EntityName = dataRow.DisplayName
			break
		}
	}
	if approveList, getApproveErr := service.GetTaskTemplateService().ListTaskTemplates(requestObj.RequestTemplate, string(models.TaskTypeApprove)); getApproveErr != nil {
		err = getApproveErr
		return
	} else {
		saveParam.ApprovalList = append(saveParam.ApprovalList, approveList...)
	}
	if approveList, getApproveErr := service.GetTaskTemplateService().ListTaskTemplates(requestObj.RequestTemplate, string(models.TaskTypeImplement)); getApproveErr != nil {
		err = getApproveErr
		return
	} else {
		saveParam.ApprovalList = append(saveParam.ApprovalList, approveList...)
	}

	if err = service.SaveRequestCacheV2(requestObj.Id, "system", token, &saveParam); err != nil {
		return
	}
	// 请求默认创建为草稿态,根据判断是否提交请求
	if input.IsDraftStatus != "true" {
		// 更新请求状态
		if err = service.UpdateRequestStatus(requestObj.Id, "Pending", "system", token, language, requestObj.Description); err != nil {
			return
		}
	}
	return
}

func SaveRequestFormData(c *gin.Context) {
	requestId := c.Param("requestId")
	user := middleware.GetRequestUser(c)
	var param models.RequestPreDataTableObj
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
	if err = service.HandleSensitiveDataDecode(&param); err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	err = service.SaveRequestForm(requestId, user, &param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, param)
	}
}

func DecodeRequestFormDataPassword(c *gin.Context) {
	encryptPwd := c.Query("encryptPwd")
	var err error
	var result string
	if strings.TrimSpace(encryptPwd) == "" {
		middleware.ReturnParamValidateError(c, fmt.Errorf("encryptPwd is empty"))
		return
	}
	// {cipher_a} 表示cmdb密码加密,{cipher_b} 表示taskman密码加密
	if !strings.HasPrefix(strings.ToLower(encryptPwd), models.EncryptSensitivePrefix) {
		middleware.ReturnParamValidateError(c, fmt.Errorf("encryptPwd format is invalid"))
		return
	}
	if result, err = service.AesDePasswordByGuid("", models.Config.EncryptSeed, encryptPwd); err != nil {
		middleware.ReturnError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}
