package request

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/services/db"
	"github.com/gin-gonic/gin"
	"time"
)

// GetRequestDetail 新版请求详情
func GetRequestDetail(c *gin.Context) {
	requestId := c.Param("requestId")
	result, err := db.GetRequestDetailV2(requestId, c.GetHeader("Authorization"))
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
		middleware.ReturnParamValidateError(c, fmt.Errorf("Param requestTemplate can not empty "))
		return
	}
	if param.Role == "" {
		middleware.ReturnParamValidateError(c, fmt.Errorf("Param role can not empty "))
		return
	}
	template, err := db.GetSimpleRequestTemplate(param.RequestTemplate)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	// 设置请求名称
	param.Name = fmt.Sprintf("%s-%s-%s", template.Name, template.OperatorObjType, time.Now().Format("060102150405"))
	// 设置请求类型
	param.Type = template.Type
	param.CreatedBy = middleware.GetRequestUser(c)
	param.ExpireDay = template.ExpireDay
	param.RequestTemplateName = template.Name
	if template.Status != "confirm" {
		param.TemplateVersion = "beta"
	} else {
		param.TemplateVersion = template.Version
	}
	err = db.CreateRequest(&param, middleware.GetRequestRoles(c), c.GetHeader("Authorization"))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	db.RecordRequestLog(param.Id, param.CreatedBy, "create")
	middleware.ReturnData(c, param)
}

func SaveRequestCache(c *gin.Context) {
	requestId := c.Param("requestId")
	cacheType := c.Param("cacheType")
	if cacheType == "data" {
		var param models.RequestProDataV2Dto
		if err := c.ShouldBindJSON(&param); err != nil {
			middleware.ReturnParamValidateError(c, err)
			return
		}
		err := db.SaveRequestCacheV2(requestId, middleware.GetRequestUser(c), c.GetHeader("Authorization"), &param)
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
		request, err := db.GetSimpleRequest(requestId)
		if err != nil {
			middleware.ReturnServerHandleError(c, err)
			return
		}
		if request.Handler != operator {
			middleware.ReturnServerHandleError(c, fmt.Errorf("hanlder %s not permission", operator))
			return
		}
		err = db.SaveRequestBindCache(requestId, operator, &param)
		if err != nil {
			middleware.ReturnServerHandleError(c, err)
		} else {
			middleware.ReturnData(c, param)
		}
	}
}

func StartRequest(c *gin.Context) {
	requestId := c.Param("requestId")
	var param models.RequestCacheData
	var operator = middleware.GetRequestUser(c)
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	request, err := db.GetSimpleRequest(requestId)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if request.Handler != operator {
		middleware.ReturnParamValidateError(c, fmt.Errorf("request handler not  permission!"))
		return
	}
	instanceId, err := db.StartRequest(requestId, operator, c.GetHeader("Authorization"), param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	db.RecordRequestLog(requestId, middleware.GetRequestUser(c), "start")
	middleware.ReturnData(c, instanceId)
}
