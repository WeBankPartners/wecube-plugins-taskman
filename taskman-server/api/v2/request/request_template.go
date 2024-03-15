package request

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/service"
	"github.com/gin-gonic/gin"
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
