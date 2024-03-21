package template

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/exterror"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/service"
	"github.com/gin-gonic/gin"
)

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
	var err error
	var list []*models.RequestTemplateTable
	if err = c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if err = validateRequestTemplateParam(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	// 校验名称是否重复
	list, err = service.GetRequestTemplateService().QueryListByName(param.Name)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if len(list) > 0 {
		middleware.ReturnServerHandleError(c, fmt.Errorf("name repeat"))
		return
	}
	param.CreatedBy = middleware.GetRequestUser(c)
	result, err := service.GetRequestTemplateService().CreateRequestTemplate(param, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
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
		middleware.ReturnError(c, fmt.Errorf("requestTemplate not exist"))
		return
	}
	// 处理时间校验
	if common.GetLowVersionUnixMillis(requestTemplate.UpdatedTime) != param.LatestUpdateTime {
		err = exterror.New().DealWithAtTheSameTimeError
		middleware.ReturnError(c, err)
		return
	}
	if requestTemplate.Status == string(models.RequestTemplateStatusConfirm) || requestTemplate.Status == string(models.RequestTemplateStatusPending) {
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
		middleware.ReturnError(c, fmt.Errorf("requestTemplate not exist"))
		return
	}
	if requestTemplate.Status != param.Status {
		middleware.ReturnParamValidateError(c, fmt.Errorf("param status invalid"))
		return
	}
	if param.TargetStatus != string(models.RequestTemplateStatusCreated) && param.TargetStatus != string(models.RequestTemplateStatusDisabled) &&
		param.TargetStatus != string(models.RequestTemplateStatusPending) && param.TargetStatus != string(models.RequestTemplateStatusConfirm) {
		middleware.ReturnParamValidateError(c, fmt.Errorf("param targetStatus invalid"))
		return
	}
	// 模板发布
	if param.TargetStatus == string(models.RequestTemplateStatusConfirm) {
		err = service.GetRequestTemplateService().ConfirmRequestTemplate(param.RequestTemplateId, middleware.GetRequestUser(c), c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
		if err != nil {
			middleware.ReturnServerHandleError(c, err)
			return
		}
		middleware.ReturnSuccess(c)
		return
	}
	if requestTemplate.Status == string(models.RequestTemplateStatusCreated) {
		// 校验是否有修改权限
		err = service.GetRequestTemplateService().CheckPermission(param.RequestTemplateId, middleware.GetRequestUser(c))
		if err != nil {
			middleware.ReturnServerHandleError(c, err)
			return
		}
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
	var result models.RequestTemplateQueryObj
	var list []*models.RequestTemplateTable
	var err error
	var repeatFlag = true
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
	// 校验是否有修改权限
	err = service.GetRequestTemplateService().CheckPermission(param.Id, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	// 校验名称重复
	list, err = service.GetRequestTemplateService().QueryListByName(param.Name)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if len(list) > 0 {
		// 只有一条数据,判断是否为当前数据
		if len(list) == 1 && param.Id == list[0].Id {
			repeatFlag = false
		} else {
			for _, template := range list {
				// 说明是 相同模版的不同版本
				if template.Id == param.RecordId {
					repeatFlag = false
				}
			}
		}
	} else {
		// 没有查询到数据repeat为false
		repeatFlag = false
	}
	if repeatFlag {
		middleware.ReturnServerHandleError(c, fmt.Errorf("name repeat"))
		return
	}
	param.UpdatedBy = middleware.GetRequestUser(c)
	result, err = service.GetRequestTemplateService().UpdateRequestTemplate(&param, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
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
	_, err := service.GetRequestTemplateService().DeleteRequestTemplate(id, false)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func ListRequestTemplateEntityAttrs(c *gin.Context) {
	id := c.Param("id")
	result, err := service.GetRequestTemplateService().ListRequestTemplateEntityAttrs(id, "", c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
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
	result, err := service.GetRequestTemplateService().GetRequestTemplateEntityAttrs(id)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func QueryRequestTemplateEntity(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	result, err := service.GetRequestTemplateService().QueryRequestTemplateEntity(id, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func UpdateRequestTemplateEntityAttrs(c *gin.Context) {
	id := c.Param("id")
	var user = middleware.GetRequestUser(c)
	if id == "" {
		middleware.ReturnParamEmptyError(c, "id")
		return
	}
	var param []*models.ProcEntityAttributeObj
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	err := service.GetRequestTemplateService().UpdateRequestTemplateEntityAttrs(id, param, user)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.GetOperationLogService().RecordRequestTemplateLog(id, "", middleware.GetRequestUser(c), "updateRequestTemplateAttr", c.Request.RequestURI, c.GetString("requestBody"))
	middleware.ReturnSuccess(c)
}

func ForkConfirmRequestTemplate(c *gin.Context) {
	requestTemplateId := c.Param("id")
	err := service.GetRequestTemplateService().ForkConfirmRequestTemplate(requestTemplateId, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.GetOperationLogService().RecordRequestTemplateLog(requestTemplateId, "", middleware.GetRequestUser(c), "forkRequestTemplate", c.Request.RequestURI, "")
	middleware.ReturnSuccess(c)
}

func GetConfirmCount(c *gin.Context) {
	var count int
	count, err := service.GetRequestTemplateService().GetConfirmCount(middleware.GetRequestUser(c), c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, count)
}

func GetRequestTemplateTags(c *gin.Context) {
	group := c.Param("requestTemplateGroup")
	result, err := service.GetRequestTemplateService().GetRequestTemplateTags(group)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func ExportRequestTemplate(c *gin.Context) {
	requestTemplateId := c.Param("requestTemplateId")
	result, err := service.GetRequestTemplateService().RequestTemplateExport(requestTemplateId)
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
	templateName, backToken, importErr := service.GetRequestTemplateService().RequestTemplateImport(paramObj, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader), "", middleware.GetRequestUser(c))
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
	_, _, err := service.GetRequestTemplateService().RequestTemplateImport(models.RequestTemplateExport{}, c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader), confirmToken, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnSuccess(c)
}

func DisableRequestTemplate(c *gin.Context) {
	requestTemplateId := c.Param("id")
	err := service.GetRequestTemplateService().DisableRequestTemplate(requestTemplateId, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.GetOperationLogService().RecordRequestTemplateLog(requestTemplateId, "", middleware.GetRequestUser(c), "disableRequestTemplate", c.Request.RequestURI, "")
	middleware.ReturnSuccess(c)
}

func EnableRequestTemplate(c *gin.Context) {
	requestTemplateId := c.Param("id")
	err := service.GetRequestTemplateService().EnableRequestTemplate(requestTemplateId, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.GetOperationLogService().RecordRequestTemplateLog(requestTemplateId, "", middleware.GetRequestUser(c), "enableRequestTemplate", c.Request.RequestURI, "")
	middleware.ReturnSuccess(c)
}
