package collect

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/services/db"
	"github.com/gin-gonic/gin"
)

// AddTemplateCollect 添加模板收藏
func AddTemplateCollect(c *gin.Context) {
	templateId := c.Param("templateId")
	param := &models.CollectTemplateTable{
		RequestTemplate: templateId,
		User:            middleware.GetRequestUser(c),
	}
	err := db.AddTemplateCollect(param)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnSuccess(c)
	}
}

// CancelTemplateCollect  取消收藏模板
func CancelTemplateCollect(c *gin.Context) {
	templateId := c.Param("templateId")
	err := db.DeleteTemplateCollect(templateId, middleware.GetRequestUser(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnSuccess(c)
	}
}

// QueryTemplateCollect  获取收藏模板列表
func QueryTemplateCollect(c *gin.Context) {
	var param models.QueryCollectTemplateObj
	if err := c.ShouldBindJSON(&param); err != nil {
		middleware.ReturnParamValidateError(c, err)
		return
	}
	if param.PageSize == 0 {
		param.PageSize = 10
	}
	pageInfo, rowData, err := db.QueryTemplateCollect(&param, middleware.GetRequestUser(c), c.GetHeader("Authorization"))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnPageData(c, pageInfo, rowData)
}
