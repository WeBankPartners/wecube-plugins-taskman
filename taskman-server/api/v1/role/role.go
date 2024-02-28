package role

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/service"
	"github.com/gin-gonic/gin"
	"strings"
)

func GetUserRoles(c *gin.Context) {
	service.GetRoleService().SyncCoreRole(c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	result, err := service.GetRoleService().GetRoleList(middleware.GetRequestRoles(c))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		middleware.ReturnData(c, result)
	}
}

func GetUserByRoles(c *gin.Context) {
	roleString := c.Query("roles")
	result, err := service.GetRoleService().QueryUserByRoles(strings.Split(roleString, ","),
		c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func GetRoleAdministrators(c *gin.Context) {
	role := c.Query("role")
	result, err := service.GetRoleService().GetRoleAdministrators(role,
		c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}

func GetRoleList(c *gin.Context) {
	service.GetRoleService().SyncCoreRole(c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader))
	result, err := service.GetRoleService().GetRoleList([]string{})
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	middleware.ReturnData(c, result)
}
