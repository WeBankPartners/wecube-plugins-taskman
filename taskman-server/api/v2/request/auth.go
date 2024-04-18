package request

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/rpc"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type transResponseJson struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}

func TransAuthGetApplyRoles(c *gin.Context) {
	var allRoleResponse models.QueryRolesResponse
	var result, hasRoleList []*models.SimpleLocalRoleDto
	var roleApplyResponse models.ListRoleApplyResponse
	var err error
	var allRoleBytes, roleApplyBytes []byte
	var requestParam models.QueryRequestParam
	var exist, roleAdmin bool
	paramAll := c.Query("all")
	paramRoleAdmin := c.Query("roleAdmin")
	allRoleBytes, err = rpc.HttpGet(models.Config.Wecube.BaseUrl+fmt.Sprintf("/auth/v1/roles?all=%s&roleAdmin=%s", paramAll, paramRoleAdmin),
		c.GetHeader("Authorization"),
		c.GetHeader(middleware.AcceptLanguageHeader))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if paramRoleAdmin != "" {
		roleAdmin, err = strconv.ParseBool(paramRoleAdmin)
		if err != nil {
			middleware.ReturnServerHandleError(c, err)
			return
		}
		if roleAdmin {
			c.Data(http.StatusOK, "application/json; charset=utf-8", allRoleBytes)
			return
		}
	}
	// 已经拥有角色&已经申请并且没有过期角色都需要过滤掉
	if err = json.Unmarshal(allRoleBytes, &allRoleResponse); err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if allRoleResponse.Status != "OK" {
		err = fmt.Errorf(allRoleResponse.Message)
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if len(allRoleResponse.Data) == 0 {
		middleware.ReturnSuccess(c)
		return
	}
	if hasRoleList, err = rpc.QueryUserRoles(middleware.GetRequestUser(c), c.GetHeader("Authorization"), c.GetHeader(middleware.AcceptLanguageHeader)); err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if len(hasRoleList) == 0 {
		hasRoleList = make([]*models.SimpleLocalRoleDto, 0)
	}
	requestParam.Filters = []*models.QueryRequestFilterObj{{Name: "status", Operator: "in", Value: "init"}}
	requestParam.Pageable = &models.PageInfo{
		StartIndex: 0,
		PageSize:   10000,
	}
	requestParam.Paging = true
	requestParamBytes, _ := json.Marshal(requestParam)
	roleApplyBytes, err = rpc.HttpPost(models.Config.Wecube.BaseUrl+"/auth/v1/roles/apply/byhandler",
		c.GetHeader("Authorization"),
		c.GetHeader(middleware.AcceptLanguageHeader), requestParamBytes)
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if err = json.Unmarshal(roleApplyBytes, &roleApplyResponse); err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if roleApplyResponse.Status != "OK" {
		err = fmt.Errorf(roleApplyResponse.Message)
		middleware.ReturnServerHandleError(c, err)
		return
	}
	if roleApplyResponse.Data == nil {
		roleApplyResponse.Data = &models.RoleApply{Contents: make([]*models.RoleApplyDto, 0)}
	} else {
		var tempApply = &models.RoleApply{Contents: make([]*models.RoleApplyDto, 0)}
		// 过滤掉申请过期数据
		for _, roleApply := range roleApplyResponse.Data.Contents {
			if roleApply.ExpireTime != "" {
				t, _ := time.Parse(models.DateTimeFormat, roleApply.ExpireTime)
				if time.Now().Before(t) {
					tempApply.Contents = append(tempApply.Contents, roleApply)
				}
			}
		}
		roleApplyResponse.Data = tempApply
	}
	for _, roleDto := range allRoleResponse.Data {
		exist = false
		for _, hasRoleDto := range hasRoleList {
			if hasRoleDto.ID == roleDto.ID {
				exist = true
				continue
			}
		}
		if exist {
			continue
		}
		for _, roleApply := range roleApplyResponse.Data.Contents {
			if roleApply.Role != nil && roleApply.Role.ID == roleDto.ID {
				exist = true
				continue
			}
		}
		if !exist {
			result = append(result, roleDto)
		}
	}
	response := transResponseJson{Status: "ok", Data: result}
	bodyBytes, _ := json.Marshal(response)
	c.Set("responseBody", string(bodyBytes))
	c.JSON(http.StatusOK, response)
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

func TransAuthUserRegister(c *gin.Context) {
	bodyString := c.GetString("requestBody")
	responseBytes, err := rpc.HttpPost(models.Config.Wecube.BaseUrl+"/auth/v1/users/register",
		c.GetHeader("Authorization"),
		c.GetHeader(middleware.AcceptLanguageHeader), []byte(bodyString))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
	} else {
		c.Data(http.StatusOK, "application/json; charset=utf-8", responseBytes)
	}
}
