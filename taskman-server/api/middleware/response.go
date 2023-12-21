package middleware

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/exterror"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

const AcceptLanguageHeader = "Accept-Language"

func ReturnPageData(c *gin.Context, pageInfo models.PageInfo, contents interface{}) {
	if contents == nil {
		contents = []string{}
	}
	obj := models.ResponseJson{StatusCode: "OK", Data: models.ResponsePageData{PageInfo: pageInfo, Contents: contents}}
	bodyBytes, _ := json.Marshal(obj)
	c.Set("responseBody", string(bodyBytes))
	c.JSON(http.StatusOK, obj)
}

func ReturnEmptyPageData(c *gin.Context) {
	c.JSON(http.StatusOK, models.ResponseJson{StatusCode: "OK", Data: models.ResponsePageData{PageInfo: models.PageInfo{StartIndex: 0, PageSize: 0, TotalRows: 0}, Contents: []string{}}})
}

func ReturnData(c *gin.Context, data interface{}) {
	if data == nil {
		data = []string{}
	}
	obj := models.ResponseJson{StatusCode: "OK", Data: data}
	bodyBytes, _ := json.Marshal(obj)
	c.Set("responseBody", string(bodyBytes))
	c.JSON(http.StatusOK, obj)
}

func ReturnSuccess(c *gin.Context) {
	c.Set("responseBody", "{\"statusCode\":\"OK\",\"data\":[]}")
	c.JSON(http.StatusOK, models.ResponseJson{StatusCode: "OK", Data: []string{}})
}

func ReturnError(c *gin.Context, err error) {
	if _, ok := err.(exterror.CustomError); !ok {
		err = exterror.Catch(exterror.New().ServerHandleError, err)
	}
	_, errorKey, errorMessage := exterror.GetErrorResult(c.GetHeader(AcceptLanguageHeader), err)
	log.Logger.Error("Handle error", log.String("statusCode", errorKey), log.String("message", errorMessage))
	obj := models.ResponseErrorJson{StatusCode: errorKey, StatusMessage: errorMessage}
	bodyBytes, _ := json.Marshal(obj)
	c.Set("responseBody", string(bodyBytes))
	c.JSON(http.StatusOK, obj)
}

func ReturnParamValidateError(c *gin.Context, err error) {
	if _, ok := err.(exterror.CustomError); !ok {
		err = exterror.Catch(exterror.New().RequestParamValidateError, err)
	}
	ReturnError(c, err)
}

func ReturnParamEmptyError(c *gin.Context, paramName string) {
	ReturnParamValidateError(c, fmt.Errorf("param %s empty", paramName))
	//ReturnError(c, "PARAM_EMPTY_ERROR", paramName, nil)
}

func ReturnServerHandleError(c *gin.Context, err error) {
	if _, ok := err.(exterror.CustomError); !ok {
		err = exterror.Catch(exterror.New().ServerHandleError, err)
	}
	ReturnError(c, err)
	//log.Logger.Error("Request server handle error", log.Error(err))
	//ReturnError(c, "SERVER_HANDLE_ERROR", err.Error(), nil)
}

func ReturnTokenValidateError(c *gin.Context, err error) {
	if _, ok := err.(exterror.CustomError); !ok {
		err = exterror.Catch(exterror.New().RequestTokenValidateError, err)
	}
	ReturnError(c, err)
	//c.JSON(http.StatusUnauthorized, models.ResponseErrorJson{StatusCode: "TOKEN_VALIDATE_ERROR", StatusMessage: err.Error(), Data: nil})
}

func ReturnDataPermissionError(c *gin.Context, err error) {
	if _, ok := err.(exterror.CustomError); !ok {
		err = exterror.Catch(exterror.New().DataPermissionDeny, err)
	}
	ReturnError(c, err)
	//ReturnError(c, "DATA_PERMISSION_ERROR", err.Error(), nil)
}

func ReturnDataPermissionDenyError(c *gin.Context) {
	ReturnError(c, exterror.New().DataPermissionDeny)
	//ReturnError(c, "DATA_PERMISSION_DENY", "permission deny", nil)
}

func ReturnApiPermissionError(c *gin.Context) {
	ReturnError(c, exterror.New().ApiPermissionDeny)
	//ReturnError(c, "API_PERMISSION_ERROR", "api permission deny", nil)
}

func ReturnTemplateAlreadyCollectError(c *gin.Context) {
	ReturnError(c, exterror.New().TemplateAlreadyCollect)
}

func ReturnUploadFileTooLargeError(c *gin.Context) {
	ReturnError(c, exterror.New().UploadFileTooLarge)
}

func ReturnChangeTaskStatusError(c *gin.Context) {
	ReturnError(c, exterror.New().ChangeTaskStatusError)
}

func ReturnSubmitRequestNotPermissionError(c *gin.Context) {
	ReturnError(c, exterror.New().SubmitRequestNotPermission)
}

func ReturnSaveRequestNotPermissionError(c *gin.Context) {
	ReturnError(c, exterror.New().SaveRequestNotPermission)
}

func ReturnTaskApproveNotPermissionError(c *gin.Context) {
	ReturnError(c, exterror.New().TaskApproveNotPermission)
}

func ReturnTaskSaveNotPermissionError(c *gin.Context) {
	ReturnError(c, exterror.New().TaskSaveNotPermission)
}

func ReturnUpdateRequestHandlerStatusError(c *gin.Context) {
	ReturnError(c, exterror.New().UpdateRequestHandlerStatusError)
}

func ReturnRevokeRequestErrorError(c *gin.Context) {
	ReturnError(c, exterror.New().RevokeRequestError)
}

func ReturnGetRequestPreviewDataError(c *gin.Context) {
	ReturnError(c, exterror.New().GetRequestPreviewDataError)
}

func InitHttpError() {
	err := exterror.InitErrorTemplateList(models.Config.HttpServer.ErrorTemplateDir, models.Config.HttpServer.ErrorDetailReturn)
	if err != nil {
		log.Logger.Error("Init error template list fail", log.Error(err))
	}
}
