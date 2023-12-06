package request

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/services/db"
	"github.com/gin-gonic/gin"
)

// GetRequestTemplateByUser 选择模板
func GetRequestTemplateByUser(c *gin.Context) {
	result, err := db.GetRequestTemplateByUserV2(middleware.GetRequestUser(c), c.GetHeader("Authorization"), middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}
