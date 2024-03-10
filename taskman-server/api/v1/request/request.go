package request

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/exterror"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/service"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func GetRequestPreviewData(c *gin.Context) {
	requestId := c.Query("requestId")
	entityDataId := c.Query("rootEntityId")
	result, err := service.GetRequestPreData(requestId, entityDataId, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, models.RequestPreDataDto{RootEntityId: entityDataId, Data: result})
}

// CountRequest 个人工作台统计
func CountRequest(c *gin.Context) {
	platformData, err := service.GetRequestCount(middleware.GetRequestUser(c), middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, platformData)
}

// FilterItem 过滤数据
func FilterItem(c *gin.Context) {
	var param models.FilterRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	data, err := service.GetFilterItem(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, data)
}

// DataList  工作台数据列表
func DataList(c *gin.Context) {
	var param models.PlatformRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if param.Tab == "" {
		param.Tab = "pending"
	}
	if param.Action == 0 {
		param.Action = 1
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	pageInfo, rowData, err := service.DataList(&param, middleware.GetRequestRoles(c), c.GetHeader("Authorization"), middleware.GetRequestUser(c), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnPageData(c, pageInfo, rowData)
}

// RevokeRequest 撤回请求
func RevokeRequest(c *gin.Context) {
	requestId := c.Param("requestId")
	err := service.RevokeRequest(requestId, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

// HistoryList 发布/请求历史
func HistoryList(c *gin.Context) {
	var param models.RequestHistoryParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	pageInfo, rowData, err := service.HistoryList(&param, middleware.GetRequestRoles(c), c.GetHeader("Authorization"), middleware.GetRequestUser(c), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnPageData(c, pageInfo, rowData)
}

// Export 请求导出
func Export(c *gin.Context) {
	var param models.RequestHistoryParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := service.Export(c.Writer, &param, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader), middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

// GetTaskList 获取任务列表
func GetTaskList(c *gin.Context) {
	var taskList []*models.TaskTable
	var err error
	requestId := c.Param("requestId")
	taskList, err = service.GetTaskService().ListImplementTasks(requestId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, taskList)
}

// Confirm 请求确认
func Confirm(c *gin.Context) {
	var param models.RequestConfirmParam
	var request models.RequestTable
	var err error
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if param.Id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	// 权限校验,创建人才能确认
	user := middleware.GetRequestUser(c)
	request, err = service.GetSimpleRequest(param.Id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if request.CreatedBy != user {
		middleware.ReturnServerHandleError(c, exterror.New().ServerHandleError)
		return
	}
	if request.Status != string(models.RequestStatusConfirm) {
		middleware.ReturnServerHandleError(c, fmt.Errorf("request status not confirm"))
		return
	}
	err = service.RequestConfirm(param, user)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func ListRequest(c *gin.Context) {
	permission := c.Param("permission")
	var param models.QueryRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	pageInfo, rowData, err := service.ListRequest(&param, middleware.GetRequestRoles(c), c.GetHeader("Authorization"), permission, middleware.GetRequestUser(c), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnPageData(c, pageInfo, rowData)
	}
}

func GetRequest(c *gin.Context) {
	requestId := c.Param("requestId")
	result, err := service.GetRequestWithRoot(requestId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func GetRequestDetail(c *gin.Context) {
	requestId := c.Param("requestId")
	result, err := service.GetRequestTaskList(requestId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func GetRequestRootForm(c *gin.Context) {
	requestId := c.Param("requestId")
	result, err := service.GetRequestRootForm(requestId)
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
	err := service.CreateRequest(&param, middleware.GetRequestRoles(c), c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.GetOperationLogService().RecordRequestLog(param.Id, param.Name, param.CreatedBy, "createRequest", c.Request.RequestURI, c.GetString("requestBody"))
	middleware.ReturnData(c, param)
}

func UpdateRequest(c *gin.Context) {
	var param models.RequestTable
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	param.Id = c.Param("requestId")
	if param.Id == "" || param.Name == "" {
		middleware.ReturnParamValidateError(c, fmt.Errorf("Param id and name can not empty "))
		return
	}
	param.UpdatedBy = middleware.GetRequestUser(c)
	err := service.UpdateRequest(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.GetOperationLogService().RecordRequestLog(param.Id, param.Name, param.UpdatedBy, "updateRequest", c.Request.RequestURI, c.GetString("requestBody"))
	middleware.ReturnData(c, param)
}

func DeleteRequest(c *gin.Context) {
	requestId := c.Param("requestId")
	err := service.DeleteRequest(requestId, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.GetOperationLogService().RecordRequestLog(requestId, "", middleware.GetRequestUser(c), "deleteRequest", c.Request.RequestURI, "")
	middleware.ReturnSuccess(c)
}

func SaveRequestCache(c *gin.Context) {
	requestId := c.Param("requestId")
	cacheType := c.Param("cacheType")
	if cacheType == "data" {
		var param models.RequestPreDataDto
		if err := c.ShouldBindJSON(&param); err != nil {
			middleware.ReturnParamValidateError(c, err)
			return
		}
		err := service.SaveRequestCacheNew(requestId, middleware.GetRequestUser(c), c.GetHeader("Authorization"), &param)
		if err != nil {
			middleware.ReturnServerHandleError(c, err)
		} else {
			middleware.ReturnData(c, param)
		}
	} else {
		var param models.RequestCacheData
		if err := c.ShouldBindJSON(&param); err != nil {
			middleware.ReturnParamValidateError(c, err)
			return
		}
		err := service.SaveRequestBindCache(requestId, middleware.GetRequestUser(c), &param)
		if err != nil {
			middleware.ReturnServerHandleError(c, err)
		} else {
			middleware.ReturnData(c, param)
		}
	}
}

func GetRequestCache(c *gin.Context) {
	requestId := c.Param("requestId")
	cacheType := c.Param("cacheType")
	result, err := service.GetRequestCache(requestId, cacheType)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func StartRequest(c *gin.Context) {
	requestId := c.Param("requestId")
	var param models.RequestCacheData
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	instanceId, err := service.StartRequest(requestId, middleware.GetRequestUser(c), c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader), param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.GetOperationLogService().RecordRequestLog(requestId, "", middleware.GetRequestUser(c), "startRequest", c.Request.RequestURI, c.GetString("requestBody"))
	middleware.ReturnData(c, instanceId)
}

func UpdateRequestStatus(c *gin.Context) {
	requestId := c.Param("requestId")
	status := c.Param("status")
	if requestId == "" || status == "" {
		middleware.ReturnParamValidateError(c, fmt.Errorf("url param can not empty"))
		return
	}
	var description string
	var param models.UpdateRequestStatusParam
	if bindErr := c.ShouldBindJSON(&param); bindErr == nil {
		description = param.Description
	}
	err := service.UpdateRequestStatus(requestId, status, middleware.GetRequestUser(c), c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader), description)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.GetOperationLogService().RecordRequestLog(requestId, "", middleware.GetRequestUser(c), "setRequestStatus", c.Request.RequestURI, status)
	middleware.ReturnSuccess(c)
}

func TerminateRequest(c *gin.Context) {
	requestId := c.Param("requestId")
	err := service.RequestTermination(requestId, middleware.GetRequestUser(c), c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnSuccess(c)
	}
}

func GetCmdbReferenceData(c *gin.Context) {
	attrId := c.Param("attrId")
	var param models.QueryRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	resultBytes, statusCode, err := service.GetCmdbReferenceData(attrId, c.GetHeader("Authorization"), param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		c.Data(statusCode, "application/json", resultBytes)
	}
}

func GetReferenceData(c *gin.Context) {
	formItemTemplateId := c.Param("formItemTemplateId")
	requestId := c.Param("requestId")
	var param models.QueryRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	input := models.RefSelectParam{FormItemTemplateId: formItemTemplateId, RequestId: requestId, Param: &param, UserToken: c.GetHeader("Authorization")}
	result, err := service.GetCMDBRefSelectResult(&input)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		result = service.FilterInSideData(result, input.FormItemTemplate, requestId)
		middleware.ReturnData(c, result)
	}
}

func UploadRequestAttachFile(c *gin.Context) {
	requestId := c.Param("requestId")
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
	err = service.UploadAttachFile(requestId, "", "", file.Filename, middleware.GetRequestUser(c), b)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, service.GetRequestAttachFileList(requestId))
	}
}

func DownloadAttachFile(c *gin.Context) {
	fileId := c.Param("fileId")
	if err := service.CheckAttachFilePermission(fileId, middleware.GetRequestUser(c), "download", middleware.GetRequestRoles(c)); err != nil {
		middleware.ReturnDataPermissionDenyError(c)
		return
	}
	fileContent, fileName, err := service.DownloadAttachFile(fileId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", fileName))
		c.Data(http.StatusOK, "application/octet-stream", fileContent)
	}
}

// UpdateRequestHandler 更新请求处理人,包括认领&转给我逻辑
func UpdateRequestHandler(c *gin.Context) {
	requestId := c.Param("requestId")
	lastedUpdateTime := c.Param("latestUpdateTime")
	request, err := service.GetSimpleRequest(requestId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if common.GetLowVersionUnixMillis(request.UpdatedTime) != lastedUpdateTime {
		middleware.ReturnDealWithAtTheSameTimeError(c)
		return
	}
	// 请求在Pending状态才有转给我
	if request.Status != "Pending" {
		middleware.ReturnUpdateRequestHandlerStatusError(c)
		return
	}
	err = service.UpdateRequestHandler(requestId, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func RemoveAttachFile(c *gin.Context) {
	fileId := c.Param("fileId")
	if err := service.CheckAttachFilePermission(fileId, middleware.GetRequestUser(c), "delete", middleware.GetRequestRoles(c)); err != nil {
		middleware.ReturnDataPermissionDenyError(c)
		return
	}
	fileObj, err := service.RemoveAttachFile(fileId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		if fileObj.Request != "" {
			middleware.ReturnData(c, service.GetRequestAttachFileList(fileObj.Request))
		} else {
			middleware.ReturnData(c, service.GetTaskAttachFileList(fileObj.Task))
		}
	}
}

func QueryProcDefEntity(c *gin.Context) {
	result := models.ProcDefRootEntityResponse{
		HttpResponseMeta: models.HttpResponseMeta{Status: "OK", Message: "Success"},
		Data:             []*models.ProcDefEntityDataObj{},
	}
	result.Data = append(result.Data, &models.ProcDefEntityDataObj{Id: "taskman_request_id", DisplayName: "request"})
	c.JSON(http.StatusOK, result)
}

func CopyRequest(c *gin.Context) {
	requestId := c.Param("requestId")
	createdBy := middleware.GetRequestUser(c)
	result, err := service.CopyRequest(requestId, createdBy)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		service.GetOperationLogService().RecordRequestLog(requestId, "", middleware.GetRequestUser(c), "copyRequest", c.Request.RequestURI, "")
		middleware.ReturnData(c, result)
	}
}

func GetRequestParent(c *gin.Context) {
	requestId := c.Query("requestId")
	result, err := service.GetRequestParent(requestId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

// GetRequestProgress  获取请求进度
func GetRequestProgress(c *gin.Context) {
	var param models.RequestQueryParam
	var rowData *models.RequestProgressObj
	var err error
	if err = c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	rowData, err = service.GetRequestProgress(param.RequestId, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, rowData)
}
