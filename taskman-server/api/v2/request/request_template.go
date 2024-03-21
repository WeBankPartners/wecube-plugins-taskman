package request

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/rpc"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetRequestTemplateByUser 选择模板
func GetRequestTemplateByUser(c *gin.Context) {
	service.GetRoleService().SyncCoreRole(c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	result, err := service.GetRequestTemplateService().GetRequestTemplateByUserV2(middleware.GetRequestUser(c), c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader), middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

// GetPlatformAllModels 转发platform models查询
func GetPlatformAllModels(c *gin.Context) {
	modelList, err := rpc.QueryAllModels(c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, modelList)
	}
}

// QueryPlatformEntityData 转发platform entity数据查询
func QueryPlatformEntityData(c *gin.Context) {
	bodyString := c.GetString("requestBody")
	responseBytes, err := rpc.HttpPost(fmt.Sprintf("%s/%s/entities/%s/query", models.Config.Wecube.BaseUrl, c.Param("package"), c.Param("entity")),
		c.GetHeader("Authorization"),
		c.GetHeader(middleware.AcceptLanguageHeader), []byte(bodyString))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		c.JSON(http.StatusOK, responseBytes)
	}
}
