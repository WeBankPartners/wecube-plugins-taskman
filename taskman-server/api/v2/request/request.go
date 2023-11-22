package request

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/services/db"
	"github.com/gin-gonic/gin"
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
