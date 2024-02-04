package request

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/service"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
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
	result, err := service.GetCoreProcessListNew(c.GetHeader("Authorization"))
	//procList, err := service.GetCoreProcessListAll(c.GetHeader("Authorization"), "MGMT", models.ProcessFetchTabs)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	//result := []*models.ProcDefObj{}
	//for _, v := range procList {
	//	tmpData := models.ProcDefObj{ProcDefId: v.ProcDefId, ProcDefKey: v.ProcDefKey, ProcDefName: v.ProcDefName, Status: v.Status, CreatedTime: v.CreatedTime, RootEntity: models.ProcEntity{}}
	//	tmpEntity := v.RootEntity
	//	if filterIndex := strings.Index(tmpEntity, "{"); filterIndex > 0 {
	//		tmpEntity = tmpEntity[:filterIndex]
	//	}
	//	if strings.Contains(tmpEntity, ":") {
	//		tmpEntitySplit := strings.Split(tmpEntity, ":")
	//		tmpData.RootEntity.PackageName = tmpEntitySplit[0]
	//		tmpData.RootEntity.Name = tmpEntitySplit[1]
	//	}
	//	result = append(result, &tmpData)
	//}
	middleware.ReturnData(c, result)
}

func GetCoreProcNodes(c *gin.Context) {
	requestTemplateId := c.Param("id")
	getType := c.Param("type")
	result, err := service.GetProcessNodesByProc(models.RequestTemplateTable{Id: requestTemplateId}, c.GetHeader("Authorization"), getType)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func GetRoleList(c *gin.Context) {
	service.SyncCoreRole()
	result, err := service.GetRoleList([]string{})
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func GetUserByRoles(c *gin.Context) {
	roleString := c.Query("roles")
	result, err := service.QueryUserByRoles(strings.Split(roleString, ","), c.GetHeader("Authorization"))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func GetUserRoles(c *gin.Context) {
	service.SyncCoreRole()
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
	pageInfo, rowData, err := service.QueryRequestTemplate(&param, c.GetHeader("Authorization"), middleware.GetRequestRoles(c))
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
	result, err := service.CreateRequestTemplate(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.RecordRequestTemplateLog(result.Id, result.Name, param.CreatedBy, "createRequestTemplate", c.Request.RequestURI, c.GetString("requestBody"))
	middleware.ReturnData(c, result)
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
	if err := service.CheckRequestTemplateRoles(param.Id, middleware.GetRequestRoles(c)); err != nil {
		middleware.ReturnDataPermissionError(c, err)
		return
	}
	param.UpdatedBy = middleware.GetRequestUser(c)
	result, err := service.UpdateRequestTemplate(&param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.RecordRequestTemplateLog(result.Id, result.Name, param.CreatedBy, "updateRequestTemplate", c.Request.RequestURI, c.GetString("requestBody"))
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
	if err := service.CheckRequestTemplateRoles(id, middleware.GetRequestRoles(c)); err != nil {
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
	result, err := service.ListRequestTemplateEntityAttrs(id, c.GetHeader("Authorization"))
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
	service.RecordRequestTemplateLog(id, "", middleware.GetRequestUser(c), "updateRequestTemplateAttr", c.Request.RequestURI, c.GetString("requestBody"))
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
	service.RecordRequestTemplateLog(requestTemplateId, "", middleware.GetRequestUser(c), "forkRequestTemplate", c.Request.RequestURI, "")
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
	templateName, backToken, importErr := service.RequestTemplateImport(paramObj, c.GetHeader("Authorization"), "", middleware.GetRequestUser(c))
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
	_, _, err := service.RequestTemplateImport(models.RequestTemplateExport{}, c.GetHeader("Authorization"), confirmToken, middleware.GetRequestUser(c))
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
	service.RecordRequestTemplateLog(requestTemplateId, "", middleware.GetRequestUser(c), "disableRequestTemplate", c.Request.RequestURI, "")
	middleware.ReturnSuccess(c)
}

func EnableRequestTemplate(c *gin.Context) {
	requestTemplateId := c.Param("id")
	err := service.EnableRequestTemplate(requestTemplateId, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	service.RecordRequestTemplateLog(requestTemplateId, "", middleware.GetRequestUser(c), "enableRequestTemplate", c.Request.RequestURI, "")
	middleware.ReturnSuccess(c)
}
