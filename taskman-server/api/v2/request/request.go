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
	result, err := db.GetRequestTaskListV2(requestId, c.GetHeader("Authorization"))
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
	param.Name = fmt.Sprintf("%s-%s-%s", template.Name, template.OperatorObjType, time.Now().Format("20060102150405"))
	// 设置请求类型
	param.Type = template.Type
	param.CreatedBy = middleware.GetRequestUser(c)
	err = db.CreateRequest(&param, middleware.GetRequestRoles(c), c.GetHeader("Authorization"))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		db.RecordRequestLog(param.Id, param.CreatedBy, "create")
		middleware.ReturnData(c, param)
	}
}

func SaveRequestCache(c *gin.Context) {
	requestId := c.Param("requestId")
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
}
