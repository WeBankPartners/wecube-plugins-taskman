package request

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/services/db"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func GetEntityData(c *gin.Context) {
	id := c.Query("requestId")
	if id == "" {
		middleware.ReturnParamValidateError(c, fmt.Errorf("Param requestId can not empty "))
		return
	}
	result, err := db.GetEntityData(id, c.GetHeader("Authorization"))
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
	result, err := db.ProcessDataPreview(requestTemplateId, entityDataId, c.GetHeader("Authorization"))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func GetRequestPreviewData(c *gin.Context) {
	requestId := c.Query("requestId")
	entityDataId := c.Query("rootEntityId")
	result, err := db.GetRequestPreData(requestId, entityDataId, c.GetHeader("Authorization"))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, models.RequestPreDataDto{RootEntityId: entityDataId, Data: result})
	}
}

// CountRequest 个人工作台统计
func CountRequest(c *gin.Context) {
	platformData, err := db.GetRequestCount(middleware.GetRequestUser(c), middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, platformData)
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
	pageInfo, rowData, err := db.DataList(&param, middleware.GetRequestRoles(c), c.GetHeader("Authorization"), middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnPageData(c, pageInfo, rowData)
}

func ListRequest(c *gin.Context) {
	permission := c.Param("permission")
	var param models.QueryRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	pageInfo, rowData, err := db.ListRequest(&param, middleware.GetRequestRoles(c), c.GetHeader("Authorization"), permission, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnPageData(c, pageInfo, rowData)
	}
}

func GetRequest(c *gin.Context) {
	requestId := c.Param("requestId")
	result, err := db.GetRequestWithRoot(requestId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func GetRequestDetail(c *gin.Context) {
	requestId := c.Param("requestId")
	result, err := db.GetRequestTaskList(requestId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func GetRequestRootForm(c *gin.Context) {
	requestId := c.Param("requestId")
	result, err := db.GetRequestRootForm(requestId)
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
	err := db.CreateRequest(&param, middleware.GetRequestRoles(c), c.GetHeader("Authorization"))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		db.RecordRequestLog(param.Id, param.CreatedBy, "create")
		middleware.ReturnData(c, param)
	}
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
	err := db.UpdateRequest(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		db.RecordRequestLog(param.Id, param.UpdatedBy, "update")
		middleware.ReturnData(c, param)
	}
}

func DeleteRequest(c *gin.Context) {
	requestId := c.Param("requestId")
	err := db.DeleteRequest(requestId, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		db.RecordRequestLog(requestId, middleware.GetRequestUser(c), "delete")
		middleware.ReturnSuccess(c)
	}
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
		err := db.SaveRequestCacheNew(requestId, middleware.GetRequestUser(c), c.GetHeader("Authorization"), &param)
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
		err := db.SaveRequestBindCache(requestId, middleware.GetRequestUser(c), &param)
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
	result, err := db.GetRequestCache(requestId, cacheType)
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
	instanceId, err := db.StartRequest(requestId, middleware.GetRequestUser(c), c.GetHeader("Authorization"), param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		db.RecordRequestLog(requestId, middleware.GetRequestUser(c), "start")
		middleware.ReturnData(c, instanceId)
	}
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
	err := db.UpdateRequestStatus(requestId, status, middleware.GetRequestUser(c), c.GetHeader("Authorization"), description)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		db.RecordRequestLog(requestId, middleware.GetRequestUser(c), status)
		middleware.ReturnSuccess(c)
	}
}

func TerminateRequest(c *gin.Context) {
	requestId := c.Param("requestId")
	err := db.RequestTermination(requestId, middleware.GetRequestUser(c), c.GetHeader("Authorization"))
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
	resultBytes, statusCode, err := db.GetCmdbReferenceData(attrId, c.GetHeader("Authorization"), param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		c.Data(statusCode, "application/json", resultBytes)
	}
}

func GetReferenceData(c *gin.Context) {
	attrId := c.Param("attrId")
	requestId := c.Param("requestId")
	var param models.QueryRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	input := models.RefSelectParam{AttrId: attrId, RequestId: requestId, Param: &param, UserToken: c.GetHeader("Authorization")}
	result, err := db.GetCMDBRefSelectResult(&input)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		result = db.FilterInSideData(result, attrId, requestId)
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
		c.JSON(http.StatusInternalServerError, models.ResponseErrorJson{StatusCode: "PARAM_HANDLE_ERROR", StatusMessage: "File too large ", Data: nil})
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
	err = db.UploadAttachFile(requestId, "", file.Filename, middleware.GetRequestUser(c), b)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, db.GetRequestAttachFileList(requestId))
	}
}

func DownloadAttachFile(c *gin.Context) {
	fileId := c.Param("fileId")
	if err := db.CheckAttachFilePermission(fileId, middleware.GetRequestUser(c), "download", middleware.GetRequestRoles(c)); err != nil {
		middleware.ReturnDataPermissionDenyError(c)
		return
	}
	fileContent, fileName, err := db.DownloadAttachFile(fileId)
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
	request, err := db.GetSimpleRequest(requestId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	// @todo 需要添加转给人权限校验
	if request.Status != "Pending" {
		middleware.ReturnServerHandleError(c, fmt.Errorf("request status invalid,requestId:%s ", requestId))
		return
	}
	err = db.UpdateRequestHandler(requestId, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func RemoveAttachFile(c *gin.Context) {
	fileId := c.Param("fileId")
	if err := db.CheckAttachFilePermission(fileId, middleware.GetRequestUser(c), "delete", middleware.GetRequestRoles(c)); err != nil {
		middleware.ReturnDataPermissionDenyError(c)
		return
	}
	fileObj, err := db.RemoveAttachFile(fileId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		if fileObj.Request != "" {
			middleware.ReturnData(c, db.GetRequestAttachFileList(fileObj.Request))
		} else {
			middleware.ReturnData(c, db.GetTaskAttachFileList(fileObj.Task))
		}
	}
}

func QueryWorkflowEntity(c *gin.Context) {
	result := models.WorkflowEntityQuery{Status: "OK", Message: "Success", Data: []*models.WorkflowEntityDataObj{}}
	result.Data = append(result.Data, &models.WorkflowEntityDataObj{Id: "taskman_request_id", DisplayName: "request"})
	c.JSON(http.StatusOK, result)
}

func CopyRequest(c *gin.Context) {
	requestId := c.Param("requestId")
	createdBy := middleware.GetRequestUser(c)
	result, err := db.CopyRequest(requestId, createdBy)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		db.RecordRequestLog(requestId, createdBy, "copy")
		middleware.ReturnData(c, result)
	}
}

func GetRequestParent(c *gin.Context) {
	requestId := c.Query("requestId")
	result, err := db.GetRequestParent(requestId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}
