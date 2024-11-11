package template

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"io"
	"net/http"
	"time"

	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/exterror"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/service"
	"github.com/gin-gonic/gin"
)

var requestTemplateStatusMap = map[models.RequestTemplateStatus]bool{
	models.RequestTemplateStatusCreated:  true,
	models.RequestTemplateStatusDisabled: true,
	models.RequestTemplateStatusPending:  true,
	models.RequestTemplateStatusConfirm:  true,
	models.RequestTemplateStatusCancel:   true,
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

func GetAllLatestReleaseRequestTemplate(c *gin.Context) {
	var result []*models.RequestTemplateSimpleQueryObj
	var err error
	if result, err = service.GetRequestTemplateService().GetAllLatestReleaseRequestTemplate(models.CommonParam{Token: c.GetHeader("Authorization"),
		Language: c.GetHeader(middleware.AcceptLanguageHeader)}); err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func GetRequestTemplateRoles(c *gin.Context) {
	var param models.GetRequestTemplateRolesParam
	var err error
	var result, subRoles []string
	if err = c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}

	// 查询模版角色
	if result, err = service.GetRequestTemplateService().GetRequestTemplateRoles(param.RequestTemplateIds); err != nil {
		middleware.ReturnError(c, err)
		return
	}
	// 查询任务模版配置的角色
	if subRoles, err = service.GetRequestTemplateService().GetTaskHandleTemplateRolesByRequestTemplateIds(param.RequestTemplateIds); err != nil {
		middleware.ReturnError(c, err)
		return
	}
	result = append(result, subRoles...)
	middleware.ReturnData(c, result)
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
	list, err = service.GetRequestTemplateService().QueryListByNameNotContainsCancel(param.Name)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if len(list) > 0 {
		middleware.ReturnError(c, exterror.New().RequestTemplateNameRepeatError)
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
		middleware.ReturnRequestTemplateUpdatePermissionError(c, err)
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
	var taskTemplateList []*models.TaskTemplateDto
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
	if _, ok := requestTemplateStatusMap[models.RequestTemplateStatus(param.TargetStatus)]; !ok {
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
	// 版本弃用,需要删除 recordId以及引用该版本的recordId需要使用上一个版本的recordId
	if param.TargetStatus == string(models.RequestTemplateStatusCancel) {
		var requestTemplateTempList []*models.RequestTemplateTable
		var newRecordId string
		dao.X.SQL("select * from request_template where name = ? order by id asc", requestTemplate.Name).Find(&requestTemplateTempList)
		if len(requestTemplateTempList) > 0 {
			for i, templateTemp := range requestTemplateTempList {
				if templateTemp.Id == requestTemplate.Id && i >= 1 {
					newRecordId = requestTemplateTempList[i-1].Id
				}
			}
			for _, templateTemp := range requestTemplateTempList {
				if templateTemp.RecordId == requestTemplate.Id {
					templateTemp.RecordId = newRecordId
					if err = service.GetRequestTemplateService().UpdateRequestTemplateStatusAndRecordId(templateTemp.Id, middleware.GetRequestUser(c), templateTemp.Status, templateTemp.RecordId); err != nil {
						middleware.ReturnError(c, err)
						return
					}
				}
			}
		}
		if err = service.GetRequestTemplateService().UpdateRequestTemplateStatusAndRecordId(param.RequestTemplateId, middleware.GetRequestUser(c), param.TargetStatus, ""); err != nil {
			middleware.ReturnError(c, err)
			return
		}
		middleware.ReturnSuccess(c)
		return
	}
	if requestTemplate.Status == string(models.RequestTemplateStatusCreated) {
		// 校验是否有修改权限
		if err = service.GetRequestTemplateService().CheckPermission(param.RequestTemplateId, middleware.GetRequestUser(c)); err != nil {
			middleware.ReturnServerHandleError(c, err)
			return
		}
		// 提交审核,需要校验任务角色是否补充完整
		if param.TargetStatus == string(models.RequestTemplateStatusPending) {
			if taskTemplateList, err = service.GetTaskTemplateService().QueryTaskTemplateDtoListByRequestTemplate(requestTemplate.Id); err != nil {
				middleware.ReturnServerHandleError(c, err)
				return
			}
			for _, taskTemplate := range taskTemplateList {
				if taskTemplate.Type == string(models.TaskTypeApprove) || taskTemplate.Type == string(models.TaskTypeImplement) {
					if err = service.GetTaskTemplateService().CheckHandleTemplates(taskTemplate); err != nil {
						if taskTemplate.Type == string(models.TaskTypeApprove) {
							err = exterror.New().TemplateSubmitApproveHandlerEmptyError.WithParam(taskTemplate.Name)
						} else if taskTemplate.Type == string(models.TaskTypeImplement) {
							err = exterror.New().TemplateSubmitTaskHandlerEmptyError.WithParam(taskTemplate.Name)
						}
						middleware.ReturnError(c, err)
						return
					}
				}
			}
		}
	}
	if err = service.GetRequestTemplateService().UpdateRequestTemplateStatus(param.RequestTemplateId, middleware.GetRequestUser(c), param.TargetStatus, param.Reason); err != nil {
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
	list, err = service.GetRequestTemplateService().QueryListByNameNotContainsCancel(param.Name)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	switch len(list) {
	case 0:
		// 没有查询到数据repeat为false
		repeatFlag = false
	case 1:
		// 只有一条数据,判断是否为当前数据
		if param.Id == list[0].Id {
			repeatFlag = false
		}
	default:
		for _, template := range list {
			// 说明是 相同模版的不同版本
			if template.Id == param.RecordId {
				repeatFlag = false
			}
		}
	}
	if repeatFlag {
		middleware.ReturnError(c, exterror.New().RequestTemplateNameRepeatError)
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
		return fmt.Errorf("param name can not empty ")
	}
	if param.Group == "" {
		return fmt.Errorf("param group can not empty ")
	}
	if len(param.MGMTRoles) == 0 {
		return fmt.Errorf("param mgmt can not empty ")
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
		middleware.ReturnRequestTemplateUpdatePermissionError(c, err)
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

func CopyConfirmRequestTemplate(c *gin.Context) {
	requestTemplateId := c.Param("id")
	err := service.GetRequestTemplateService().CopyConfirmRequestTemplate(requestTemplateId, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.GetOperationLogService().RecordRequestTemplateLog(requestTemplateId, "", middleware.GetRequestUser(c), "CopyConfirmRequestTemplate", c.Request.RequestURI, "")
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
		middleware.ReturnServerHandleError(c, fmt.Errorf("export requestTemplate config fail, json marshal object error:%s ", jsonErr.Error()))
		return
	}
	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s_%s.json", "rt_"+requestTemplateId, time.Now().Format("20060102150405")))
	c.Data(http.StatusOK, "application/octet-stream", b)
}

func BatchExportRequestTemplate(c *gin.Context) {
	var err error
	var param models.BatchExportRequestTemplateParam
	var requestTemplateExport models.RequestTemplateExport
	var result []models.RequestTemplateExport
	if err = c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	for _, requestTemplateId := range param.RequestTemplateIds {
		if requestTemplateExport, err = service.GetRequestTemplateService().RequestTemplateExport(requestTemplateId); err != nil {
			middleware.ReturnError(c, err)
			return
		}
		if requestTemplateExport.RequestTemplate.Id != "" {
			result = append(result, requestTemplateExport)
		}
	}
	middleware.ReturnData(c, result)
}

// ImportRequestTemplateBatch 批量导入完然后发布
func ImportRequestTemplateBatch(c *gin.Context) {
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
	var paramObj []models.RequestTemplateExport
	b, err := io.ReadAll(f)
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
	for _, obj := range paramObj {
		exportParam := models.RequestTemplateExportParam{
			Input:        obj,
			UserToken:    c.GetHeader("Authorization"),
			Language:     c.GetHeader(middleware.AcceptLanguageHeader),
			ConfirmToken: "",
			Operator:     middleware.GetRequestUser(c),
			CoverRole:    false,
			UserRoles:    middleware.GetRequestRoles(c),
		}
		_, templateName, backToken, importErr := service.GetRequestTemplateService().RequestTemplateImport(exportParam)
		if importErr != nil {
			middleware.ReturnServerHandleError(c, importErr)
			return
		}
		if backToken != "" {
			c.JSON(http.StatusOK, models.ResponseJson{StatusCode: "CONFIRM", Data: models.ImportData{Token: backToken, TemplateName: templateName}})
			return
		}
	}
	middleware.ReturnSuccess(c)
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
	b, err := io.ReadAll(f)
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
	exportParam := models.RequestTemplateExportParam{
		Input:        paramObj,
		UserToken:    c.GetHeader("Authorization"),
		Language:     c.GetHeader(middleware.AcceptLanguageHeader),
		ConfirmToken: "",
		Operator:     middleware.GetRequestUser(c),
		CoverRole:    true,
		UserRoles:    middleware.GetRequestRoles(c),
	}
	_, templateName, backToken, importErr := service.GetRequestTemplateService().RequestTemplateImport(exportParam)
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
	exportParam := models.RequestTemplateExportParam{
		Input:        models.RequestTemplateExport{},
		UserToken:    c.GetHeader("Authorization"),
		Language:     c.GetHeader(middleware.AcceptLanguageHeader),
		ConfirmToken: confirmToken,
		Operator:     middleware.GetRequestUser(c),
		CoverRole:    true,
		UserRoles:    middleware.GetRequestRoles(c),
	}
	_, _, _, err := service.GetRequestTemplateService().RequestTemplateImport(exportParam)
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
