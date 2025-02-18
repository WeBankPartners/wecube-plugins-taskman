package request

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/try"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/exterror"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/rpc"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/service"
	"github.com/gin-gonic/gin"
)

func GetRequestPreviewData(c *gin.Context) {
	requestId := c.Query("requestId")
	entityDataId := c.Query("rootEntityId")
	result, _, err := service.GetRequestPreData(requestId, entityDataId, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, models.RequestPreDataDto{RootEntityId: entityDataId, Data: result})
}

// CountPlatform 个人工作台数量统计-new
func CountPlatform(c *gin.Context) {
	var param models.CountPlatformParam
	var err error
	if err = c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	// 默认三个月前时间
	if param.QueryTimeStart == "" {
		param.QueryTimeStart = time.Now().AddDate(0, -3, 0).Format(models.DateTimeFormat)
	}
	if param.QueryTimeEnd == "" {
		param.QueryTimeEnd = time.Now().Format(models.DateTimeFormat)
	}
	platformData, err := service.GetPlatformCount(param, middleware.GetRequestUser(c), middleware.GetRequestRoles(c))
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
	// 默认三个月前时间
	if param.QueryTimeStart == "" {
		param.QueryTimeStart = time.Now().AddDate(0, -3, 0).Format(models.DateTimeFormat)
	}
	if param.QueryTimeEnd == "" {
		param.QueryTimeEnd = time.Now().Format(models.DateTimeFormat)
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
	err = service.RequestConfirm(param, user, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
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
		middleware.ReturnParamValidateError(c, fmt.Errorf("param name and requestTemplate can not empty "))
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
		middleware.ReturnParamValidateError(c, fmt.Errorf("param id and name can not empty "))
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
	var request models.RequestTable
	var err error
	if err = c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if request, err = service.GetSimpleRequest(requestId); err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	instanceId, err := service.StartRequest(request, middleware.GetRequestUser(c), c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader), param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.GetOperationLogService().RecordRequestLog(requestId, "", middleware.GetRequestUser(c), "startRequest", c.Request.RequestURI, c.GetString("requestBody"))
	middleware.ReturnData(c, instanceId)
}

func UpdateRequestStatus(c *gin.Context) {
	defer try.ExceptionStack(func(e interface{}, err interface{}) {
		retErr := fmt.Errorf("%v", err)
		middleware.ReturnError(c, exterror.Catch(exterror.New().ServerHandleError, retErr))
		log.Logger.Error(e.(string))
	})
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
	if param.Dialect != nil {
		if param.Dialect.AssociatedData != nil {
			for k, v := range param.Dialect.AssociatedData {
				if strings.HasPrefix(v, "tmp"+models.SysTableIdConnector) {
					delete(param.Dialect.AssociatedData, k)
				}
			}
		}
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
	b, err := io.ReadAll(f)
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
	var err error
	var attachFile models.AttachFileTable
	var fileContent []byte
	var fileName string
	var checkPermission bool
	fileId := c.Param("fileId")
	if attachFile, err = service.GetAttachFileInfo(fileId); err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if checkPermission, err = service.CheckDownloadPermission(attachFile, middleware.GetRequestRoles(c)); err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if !checkPermission {
		middleware.ReturnDataPermissionDenyError(c)
		return
	}
	if fileContent, fileName, err = service.DownloadAttachFile(attachFile); err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename*=UTF-8''%s", fileName))
	c.Data(http.StatusOK, "application/octet-stream", fileContent)
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
	if err := service.CheckAttachFilePermission(fileId, middleware.GetRequestUser(c), middleware.GetRequestRoles(c)); err != nil {
		middleware.ReturnDataPermissionDenyError(c)
		return
	}
	fileObj, err := service.RemoveAttachFile(fileId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		if fileObj.Request != "" {
			middleware.ReturnData(c, service.GetRequestAttachFileList(fileObj.Request))
		} else if fileObj.TaskHandle != "" {
			middleware.ReturnData(c, service.GetAttachFileListByTaskHandleId(fileObj.TaskHandle))
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
	requestId := c.Query("requestId")
	rowData, err := service.GetRequestProgress(requestId, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, rowData)
}

func GetExpressionItemData(c *gin.Context) {
	formItemTemplateId := c.Param("formItemTemplateId")
	rootDataId := c.Param("rootDataId")
	formItemTemplateRow, err := service.GetFormItemTemplateService().GetFormItemTemplate(formItemTemplateId)
	if err != nil {
		middleware.ReturnError(c, err)
		return
	}
	if formItemTemplateRow.RoutineExpression == "" {
		middleware.ReturnError(c, fmt.Errorf("expression is empty"))
		return
	}
	result, queryErr := rpc.QueryEntityExpressionData(formItemTemplateRow.RoutineExpression, rootDataId, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if queryErr != nil {
		middleware.ReturnError(c, queryErr)
	} else {
		middleware.ReturnData(c, result)
	}
}

// Association 请求关联单
func Association(c *gin.Context) {
	var param models.RequestAssociationParam
	var pageInfo models.PageInfo
	var rowsData []*models.SimpleRequestDto
	var err error
	if err = c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if param.PageSize == 0 {
		param.PageSize = 50
	}
	if pageInfo, rowsData, err = service.GetRequestService().Association(param); err != nil {
		middleware.ReturnError(c, err)
		return
	}
	middleware.ReturnPageData(c, pageInfo, rowsData)
}

// AttrSensitiveDataQuery 1.整个itsm流程中 ，被提交人、审批人、定版人、任务人修改or新增过，这个表单的查看按钮不做权限控制——原因：这个改动value属于itsm表单，不需要权限控制
// 2.整个itsm流程中 ，没有被提交人、审批人、定版人、任务人修改or新增过，这个表单的查看按钮使用对应人的cmdb查看权限做权限控制——原因：这个原value属于cmdb，不能在itsm里绕过cmdb的查看权限去查看数据库里的密码字段/**
func AttrSensitiveDataQuery(c *gin.Context) {
	var paramList, guidNotEmptyParamList []*models.RequestFormSensitiveDataParam
	var result, subResult []*models.AttrPermissionQueryObj
	var idValueMap = make(map[string]string)
	var finalIdMap = make(map[string]string)
	var err error
	var formItemList []*models.FormItemTable
	if err = c.ShouldBindJSON(&paramList); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if len(paramList) == 0 {
		middleware.ReturnError(c, fmt.Errorf("param empty"))
		return
	}
	if strings.TrimSpace(paramList[0].RequestId) == "" {
		middleware.ReturnError(c, fmt.Errorf("param requestId empty"))
		return
	}
	if formItemList, err = service.GetFormService().QueryFormItemListByRequest(paramList[0].RequestId); err != nil {
		middleware.ReturnError(c, err)
		return
	}
	// 处理 paramList 过滤掉 guid为空数据,以及guid为临时Id,查询CMDB返回行guid不为空
	for _, param := range paramList {
		// 直接解密
		originVal, _ := service.HandleSensitiveValDecode(param.AttrVal)
		if strings.TrimSpace(param.Guid) == "" || strings.HasPrefix(param.Guid, "tmp"+models.SysTableIdConnector) {
			result = append(result, &models.AttrPermissionQueryObj{
				Guid:             param.Guid,
				CiType:           param.CiType,
				AttrName:         param.AttrName,
				TmpId:            param.TmpId,
				QueryPermission:  true,
				UpdatePermission: true,
				Value:            originVal,
				TaskHandleId:     param.TaskHandleId,
			})
		} else {
			// 审批时候, guid可能会拼接 wecmdb:business_product:business_product_6241b387c3d04818b45ab 前缀,需要截取取最后一个:
			if len(param.Guid) > 7 && strings.HasPrefix(param.Guid, "wecmdb:") {
				tmp := param.Guid
				param.Guid = param.Guid[strings.LastIndex(param.Guid, ":")+1:]
				finalIdMap[param.Guid] = tmp
			}
			guidNotEmptyParamList = append(guidNotEmptyParamList, param)
			idValueMap[param.Guid+param.AttrName] = originVal
		}
	}
	// 从CMDB查询拿到初始化数据,跟taskMan数据比对
	if len(guidNotEmptyParamList) > 0 {
		if subResult, err = service.GetCMDBCiAttrSensitiveData(guidNotEmptyParamList, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader)); err != nil {
			middleware.ReturnError(c, err)
			return
		}
		for _, item := range subResult {
			// 没有数据查询权限,数据返回为空,需要对比表单数据
			if item.Value == "" && !item.QueryPermission {
				for _, formItem := range formItemList {
					if strings.HasSuffix(formItem.RowDataId, item.Guid) && formItem.Name == item.AttrName && item.TaskHandleId == formItem.TaskHandle {
						if formItem.ModifyFlag == 1 {
							// 修改过
							item.QueryPermission = true
							item.UpdatePermission = true
							item.Value = formItem.Value
						}
					}
				}
			} else {
				// 敏感数据被修改了,直接展示
				if v, ok := idValueMap[item.Guid+item.AttrName]; ok {
					if v == item.Value && !item.QueryPermission {
						// 没有查询权限,数据相等,也不能表示没有修改,需要查询表单数据数据确定是否修改过
						for _, formItem := range formItemList {
							if strings.HasSuffix(formItem.RowDataId, item.Guid) && formItem.Name == item.AttrName && item.TaskHandleId == formItem.TaskHandle {
								if formItem.ModifyFlag == 1 {
									// 修改过
									item.QueryPermission = true
									item.UpdatePermission = true
								}
							}
						}
					} else {
						// 不相等 一定是修改过
						item.QueryPermission = true
						item.UpdatePermission = true
						item.Value = v
					}
				}
			}
			// 返回给web的guid需要返回原值
			if v, ok := finalIdMap[item.Guid]; ok {
				item.Guid = v
			}
			result = append(result, item)
		}
	}

	middleware.ReturnData(c, result)
}
