package middleware

import (
	"errors"
	"fmt"
	"strconv"

	"encoding/json"
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

func ReturnApiPermissionError(c *gin.Context) {
	ReturnError(c, exterror.New().ApiPermissionDeny)
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

func Return(c *gin.Context, response interface{}) {
	bodyBytes, _ := json.Marshal(response)
	c.Set("responseBody", string(bodyBytes))
	c.JSON(http.StatusOK, response)
}
func ReturnSuccess(c *gin.Context) {
	c.Set("responseBody", "{\"statusCode\":\"OK\",\"data\":[]}")
	c.JSON(http.StatusOK, models.ResponseJson{StatusCode: "OK", Data: []string{}})
}

func ReturnError(c *gin.Context, err error) {
	var customError exterror.CustomError
	if !errors.As(err, &customError) {
		err = exterror.Catch(exterror.New().ServerHandleError, err)
	}
	errorCode, errorKey, errorMessage := exterror.GetErrorResult(c.GetHeader(AcceptLanguageHeader), err)
	log.Logger.Error("Handle error", log.String("statusCode", errorKey), log.String("message", errorMessage))
	obj := models.ResponseErrorJson{StatusCode: strconv.Itoa(errorCode), StatusMessage: errorMessage}
	bodyBytes, _ := json.Marshal(obj)
	c.Writer.Header().Add("Error-Code", strconv.Itoa(errorCode))
	c.Set("responseBody", string(bodyBytes))
	c.JSON(http.StatusOK, obj)
}

func ReturnParamValidateError(c *gin.Context, err error) {
	var customError exterror.CustomError
	if !errors.As(err, &customError) {
		err = exterror.Catch(exterror.New().RequestParamValidateError, err)
	}
	c.Writer.Header().Add("Error-Code", strconv.Itoa(exterror.New().RequestParamValidateError.Code))
	ReturnError(c, err)
}

func ReturnParamEmptyError(c *gin.Context, paramName string) {
	ReturnParamValidateError(c, fmt.Errorf("param %s empty", paramName))
}

func ReturnServerHandleError(c *gin.Context, err error) {
	var customError exterror.CustomError
	if !errors.As(err, &customError) {
		err = exterror.Catch(exterror.New().ServerHandleError, err)
	}
	ReturnError(c, err)
}

func ReturnTokenValidateError(c *gin.Context, err error) {
	var customError exterror.CustomError
	if !errors.As(err, &customError) {
		err = exterror.Catch(exterror.New().RequestTokenValidateError, err)
	}
	ReturnError(c, err)
}

func ReturnDataPermissionError(c *gin.Context, err error) {
	var customError exterror.CustomError
	if !errors.As(err, &customError) {
		err = exterror.Catch(exterror.New().DataPermissionDeny, err)
	}
	ReturnError(c, err)
}

func ReturnRequestTemplateUpdatePermissionError(c *gin.Context, err error) {
	var customError exterror.CustomError
	if !errors.As(err, &customError) {
		err = exterror.Catch(exterror.New().TemplateUpdatePermissionDeny, err)
	}
	ReturnError(c, err)
}

func ReturnDataPermissionDenyError(c *gin.Context) {
	ReturnError(c, exterror.New().DataPermissionDeny)
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

func ReturnUpdateRequestHandlerStatusError(c *gin.Context) {
	ReturnError(c, exterror.New().UpdateRequestHandlerStatusError)
}

func ReturnReportRequestNotPermissionError(c *gin.Context) {
	ReturnError(c, exterror.New().ReportRequestNotPermission)
}

func ReturnDealWithAtTheSameTimeError(c *gin.Context) {
	ReturnError(c, exterror.New().DealWithAtTheSameTimeError)
}

func InitHttpError() {
	err := exterror.InitErrorTemplateList(models.Config.HttpServer.ErrorTemplateDir, models.Config.HttpServer.ErrorDetailReturn)
	if err != nil {
		log.Logger.Error("Init error template list fail", log.Error(err))
	}
}
