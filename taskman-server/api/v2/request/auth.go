package request

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/rpc"
	"github.com/gin-gonic/gin"
	"net/http"
)

func TransAuthGetApplyRoles(c *gin.Context) {
	paramAll := c.Query("all")
	paramRoleAdmin := c.Query("roleAdmin")
	responseBytes, err := rpc.HttpGet(models.Config.Wecube.BaseUrl+fmt.Sprintf("/auth/v1/roles?all=%s&roleAdmin=%s", paramAll, paramRoleAdmin),
		c.GetHeader("Authorization"),
		c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		c.Data(http.StatusOK, "application/json; charset=utf-8", responseBytes)
	}
}

func TransAuthStartApply(c *gin.Context) {
	bodyString := c.GetString("requestBody")
	responseBytes, err := rpc.HttpPost(models.Config.Wecube.BaseUrl+"/auth/v1/roles/apply",
		c.GetHeader("Authorization"),
		c.GetHeader(middleware.AcceptLanguageHeader), []byte(bodyString))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		c.Data(http.StatusOK, "application/json; charset=utf-8", responseBytes)
	}
}

func TransAuthGetProcessableList(c *gin.Context) {
	bodyString := c.GetString("requestBody")
	responseBytes, err := rpc.HttpPost(models.Config.Wecube.BaseUrl+"/auth/v1/roles/apply/byhandler",
		c.GetHeader("Authorization"),
		c.GetHeader(middleware.AcceptLanguageHeader), []byte(bodyString))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		c.Data(http.StatusOK, "application/json; charset=utf-8", responseBytes)
	}
}

func TransAuthGetAllUser(c *gin.Context) {
	responseBytes, err := rpc.HttpGet(models.Config.Wecube.BaseUrl+"/auth/v1/users",
		c.GetHeader("Authorization"),
		c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		c.Data(http.StatusOK, "application/json; charset=utf-8", responseBytes)
	}
}

func TransAuthGetUserByRole(c *gin.Context) {
	responseBytes, err := rpc.HttpGet(models.Config.Wecube.BaseUrl+"/auth/v1/roles/"+c.Param("roleId")+"/users",
		c.GetHeader("Authorization"),
		c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		c.Data(http.StatusOK, "application/json; charset=utf-8", responseBytes)
	}
}

func TransAuthRemoveUserFromRole(c *gin.Context) {
	bodyString := c.GetString("requestBody")
	responseBytes, err := rpc.HttpPost(models.Config.Wecube.BaseUrl+"/auth/v1/roles/"+c.Param("roleId")+"/users/revoke",
		c.GetHeader("Authorization"),
		c.GetHeader(middleware.AcceptLanguageHeader), []byte(bodyString))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		c.Data(http.StatusOK, "application/json; charset=utf-8", responseBytes)
	}
}

func TransAuthAddUserForRole(c *gin.Context) {
	bodyString := c.GetString("requestBody")
	responseBytes, err := rpc.HttpPost(models.Config.Wecube.BaseUrl+"/auth/v1/roles/"+c.Param("roleId")+"/users",
		c.GetHeader("Authorization"),
		c.GetHeader(middleware.AcceptLanguageHeader), []byte(bodyString))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		c.Data(http.StatusOK, "application/json; charset=utf-8", responseBytes)
	}
}

func TransAuthHandleApplication(c *gin.Context) {
	bodyString := c.GetString("requestBody")
	responseBytes, err := rpc.HttpPut(models.Config.Wecube.BaseUrl+"/auth/v1/roles/apply",
		c.GetHeader("Authorization"),
		c.GetHeader(middleware.AcceptLanguageHeader), []byte(bodyString))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		c.Data(http.StatusOK, "application/json; charset=utf-8", responseBytes)
	}
}

func TransAuthGetApplyList(c *gin.Context) {
	bodyString := c.GetString("requestBody")
	responseBytes, err := rpc.HttpPost(models.Config.Wecube.BaseUrl+"/auth/v1/roles/apply/byapplier",
		c.GetHeader("Authorization"),
		c.GetHeader(middleware.AcceptLanguageHeader), []byte(bodyString))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		c.Data(http.StatusOK, "application/json; charset=utf-8", responseBytes)
	}
}
