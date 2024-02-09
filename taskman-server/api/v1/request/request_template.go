package request

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/exterror"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/service"
	"github.com/gin-gonic/gin"
)

func QueryRequestTemplateGroup(c *gin.Context) {
	var param models.QueryRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	pageInfo, rowData, err := service.QueryRequestTemplateGroup(&param, middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnPageData(c, pageInfo, rowData)
}

func CreateRequestTemplateGroup(c *gin.Context) {
	var param models.RequestTemplateGroupTable
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := service.CreateRequestTemplateGroup(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func UpdateRequestTemplateGroup(c *gin.Context) {
	var param models.RequestTemplateGroupTable
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if param.Id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	err := service.CheckRequestTemplateGroupRoles(param.Id, middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnDataPermissionError(c, err)
		return
	}
	err = service.UpdateRequestTemplateGroup(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func DeleteRequestTemplateGroup(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	err := service.CheckRequestTemplateGroupRoles(id, middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnDataPermissionError(c, err)
		return
	}
	err = service.DeleteRequestTemplateGroup(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func GetCoreProcessList(c *gin.Context) {
	mangeRole := c.Query("role")
	result, err := service.GetProcDefService().GetCoreProcessListNew(c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader), mangeRole)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func GetCoreProcNodes(c *gin.Context) {
	var nodeList []*models.ProcNodeObj
	var err error
	requestTemplateId := c.Param("id")
	getType := c.Param("type")
	nodeList, err = service.GetProcDefService().GetProcessDefineTaskNodes(models.RequestTemplateTable{Id: requestTemplateId}, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader), getType)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, nodeList)
}

func GetRoleList(c *gin.Context) {
	service.SyncCoreRole(c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	result, err := service.GetRoleList([]string{})
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func GetUserByRoles(c *gin.Context) {
	roleString := c.Query("roles")
	result, err := service.QueryUserByRoles(strings.Split(roleString, ","), c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func GetUserRoles(c *gin.Context) {
	service.SyncCoreRole(c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	result, err := service.GetRoleList(middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func QueryRequestTemplate(c *gin.Context) {
	var param models.QueryRequestParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	pageInfo, rowData, err := service.GetRequestTemplateService().QueryRequestTemplate(&param, models.CommonParam{Token: c.GetHeader("Authorization"),
		Language: c.GetHeader(middleware.AcceptLanguageHeader), Roles: middleware.GetRequestRoles(c)})
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnPageData(c, pageInfo, rowData)
}

func CreateRequestTemplate(c *gin.Context) {
	var param models.RequestTemplateUpdateParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if err := validateRequestTemplateParam(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	param.CreatedBy = middleware.GetRequestUser(c)
	result, err := service.GetRequestTemplateService().CreateRequestTemplate(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.GetOperationLogService().RecordRequestTemplateLog(result.Id, result.Name, param.CreatedBy, "createRequestTemplate", c.Request.RequestURI, c.GetString("requestBody"))
	middleware.ReturnData(c, result)
}

// UpdateRequestTemplateHandler 请求模板处理:转给我
func UpdateRequestTemplateHandler(c *gin.Context) {
	var param models.RequestTemplateHandlerDto
	var requestTemplate *models.RequestTemplateTable
	var err error
	if err = c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if param.RequestTemplateId == "" {
		middleware.ReturnParamEmptyError(c, "requestTemplateId")
		return
	}
	requestTemplate, err = service.GetRequestTemplateService().GetRequestTemplate(param.RequestTemplateId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if requestTemplate == nil {
		middleware.ReturnError(c, fmt.Errorf("requestTemplateId not exist"))
		return
	}
	// 处理时间校验
	if common.GetLowVersionUnixMillis(requestTemplate.UpdatedTime) != param.LatestUpdateTime {
		err = exterror.New().DealWithAtTheSameTimeError
		middleware.ReturnError(c, err)
		return
	}
	if requestTemplate.Status == string(models.RequestTemplateStatusConfirm) {
		middleware.ReturnError(c, fmt.Errorf("request template has deployed"))
		return
	}
	if err := service.GetRequestTemplateService().CheckRequestTemplateRoles(param.RequestTemplateId, middleware.GetRequestRoles(c)); err != nil {
		middleware.ReturnDataPermissionError(c, err)
		return
	}
	err = service.GetRequestTemplateService().UpdateRequestTemplateHandler(requestTemplate.Id, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

// UpdateRequestTemplateStatus 更新请求模板状态
func UpdateRequestTemplateStatus(c *gin.Context) {
	var param models.RequestTemplateStatusUpdateParam
	var requestTemplate *models.RequestTemplateTable
	var err error
	if err = c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	requestTemplate, err = service.GetRequestTemplateService().GetRequestTemplate(param.RequestTemplateId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if requestTemplate == nil {
		middleware.ReturnError(c, fmt.Errorf("requestTemplateId not exist"))
		return
	}
	if requestTemplate.Status != param.Status {
		middleware.ReturnParamValidateError(c, fmt.Errorf("param status invalid"))
		return
	}
	err = service.GetRequestTemplateService().UpdateRequestTemplateStatus(param.RequestTemplateId, middleware.GetRequestUser(c), param.TargetStatus, param.Reason)
	if err != nil {
		middleware.ReturnError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func UpdateRequestTemplate(c *gin.Context) {
	var param models.RequestTemplateUpdateParam
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if err := validateRequestTemplateParam(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if param.Id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	if err := service.GetRequestTemplateService().CheckRequestTemplateRoles(param.Id, middleware.GetRequestRoles(c)); err != nil {
		middleware.ReturnDataPermissionError(c, err)
		return
	}
	param.UpdatedBy = middleware.GetRequestUser(c)
	result, err := service.UpdateRequestTemplate(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.GetOperationLogService().RecordRequestTemplateLog(result.Id, result.Name, param.CreatedBy, "updateRequestTemplate", c.Request.RequestURI, c.GetString("requestBody"))
	middleware.ReturnData(c, result)
}

func validateRequestTemplateParam(param *models.RequestTemplateUpdateParam) error {
	if param.Name == "" {
		return fmt.Errorf("Param name can not empty ")
	}
	if param.Group == "" {
		return fmt.Errorf("Param group can not empty ")
	}
	if param.ProcDefKey == "" || param.ProcDefId == "" || param.ProcDefName == "" {
		return fmt.Errorf("Param procDef can not empty ")
	}
	if len(param.MGMTRoles) == 0 {
		return fmt.Errorf("Param mgmt can not empty ")
	}
	return nil
}

func DeleteRequestTemplate(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	if err := service.GetRequestTemplateService().CheckRequestTemplateRoles(id, middleware.GetRequestRoles(c)); err != nil {
		middleware.ReturnDataPermissionError(c, err)
		return
	}
	_, err := service.DeleteRequestTemplate(id, false)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func ListRequestTemplateEntityAttrs(c *gin.Context) {
	id := c.Param("id")
	result, err := service.ListRequestTemplateEntityAttrs(id, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func GetRequestTemplateEntityAttrs(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	result, err := service.GetRequestTemplateEntityAttrs(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func UpdateRequestTemplateEntityAttrs(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	var param []*models.ProcEntityAttributeObj
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := service.UpdateRequestTemplateEntityAttrs(id, param, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.SetRequestTemplateToCreated(id, middleware.GetRequestUser(c))
	service.GetOperationLogService().RecordRequestTemplateLog(id, "", middleware.GetRequestUser(c), "updateRequestTemplateAttr", c.Request.RequestURI, c.GetString("requestBody"))
	middleware.ReturnSuccess(c)
}

func GetRequestTemplateByUser(c *gin.Context) {
	result, err := service.GetRequestTemplateByUser(middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func ForkConfirmRequestTemplate(c *gin.Context) {
	requestTemplateId := c.Param("id")
	err := service.ForkConfirmRequestTemplate(requestTemplateId, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.GetOperationLogService().RecordRequestTemplateLog(requestTemplateId, "", middleware.GetRequestUser(c), "forkRequestTemplate", c.Request.RequestURI, "")
	middleware.ReturnSuccess(c)
}

func GetRequestTemplateTags(c *gin.Context) {
	group := c.Param("requestTemplateGroup")
	result, err := service.GetRequestTemplateTags(group)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func ExportRequestTemplate(c *gin.Context) {
	requestTemplateId := c.Param("requestTemplateId")
	result, err := service.RequestTemplateExport(requestTemplateId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	b, jsonErr := json.Marshal(result)
	if jsonErr != nil {
		middleware.ReturnServerHandleError(c, fmt.Errorf("Export requestTemplate config fail, json marshal object error:%s ", jsonErr.Error()))
		return
	}
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s_%s.json", "rt_"+requestTemplateId, time.Now().Format("20060102150405")))
	c.Data(http.StatusOK, "application/octet-stream", b)
}

func ImportRequestTemplate(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseErrorJson{StatusCode: "PARAM_HANDLE_ERROR", StatusMessage: "Http read upload file fail:" + err.Error(), Data: nil})
		return
	}
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseErrorJson{StatusCode: "PARAM_HANDLE_ERROR", StatusMessage: "File open error:" + err.Error(), Data: nil})
		return
	}
	var paramObj models.RequestTemplateExport
	b, err := ioutil.ReadAll(f)
	defer f.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseErrorJson{StatusCode: "PARAM_HANDLE_ERROR", StatusMessage: "Read content fail error:" + err.Error(), Data: nil})
		return
	}
	err = json.Unmarshal(b, &paramObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ResponseErrorJson{StatusCode: "PARAM_HANDLE_ERROR", StatusMessage: "Json unmarshal fail error:" + err.Error(), Data: nil})
		return
	}
	templateName, backToken, importErr := service.RequestTemplateImport(paramObj, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader), "", middleware.GetRequestUser(c))
	if importErr != nil {
		middleware.ReturnServerHandleError(c, importErr)
		return
	}
	if backToken != "" {
		c.JSON(http.StatusOK, models.ResponseJson{StatusCode: "CONFIRM", Data: models.ImportData{Token: backToken, TemplateName: templateName}})
		return
	}
	middleware.ReturnSuccess(c)
}

func ConfirmImportRequestTemplate(c *gin.Context) {
	confirmToken := c.Param("confirmToken")
	_, _, err := service.RequestTemplateImport(models.RequestTemplateExport{}, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader), confirmToken, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func DisableRequestTemplate(c *gin.Context) {
	requestTemplateId := c.Param("id")
	err := service.DisableRequestTemplate(requestTemplateId, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.GetOperationLogService().RecordRequestTemplateLog(requestTemplateId, "", middleware.GetRequestUser(c), "disableRequestTemplate", c.Request.RequestURI, "")
	middleware.ReturnSuccess(c)
}

func EnableRequestTemplate(c *gin.Context) {
	requestTemplateId := c.Param("id")
	err := service.EnableRequestTemplate(requestTemplateId, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.GetOperationLogService().RecordRequestTemplateLog(requestTemplateId, "", middleware.GetRequestUser(c), "enableRequestTemplate", c.Request.RequestURI, "")
	middleware.ReturnSuccess(c)
}
