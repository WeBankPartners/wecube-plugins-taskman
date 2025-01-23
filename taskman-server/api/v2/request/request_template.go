package request

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/rpc"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// GetRequestTemplateByUser 选择模板
func GetRequestTemplateByUser(c *gin.Context) {
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
	plugin := c.Param("package")
	entity := c.Param("entity")
	var responseBytes []byte
	var err error
	var refAttributes []*models.EntityAttributeObj
	responseBytes, err = rpc.HttpPost(fmt.Sprintf("%s/%s/entities/%s/query", models.Config.Wecube.BaseUrl, plugin, entity),
		c.GetHeader("Authorization"),
		c.GetHeader(middleware.AcceptLanguageHeader), []byte(bodyString))
	if err != nil {
		middleware.ReturnServerHandleError(c, err)
		return
	}
	// cmdb 涉及到敏感数据需要隐藏显示
	if plugin == "wecmdb" {
		var response models.EntityResponse
		if refAttributes, err = service.GetCMDBCiAttrDefs(entity, c.GetHeader("Authorization")); err != nil {
			err = fmt.Errorf("query remote entity:%s attr fail:%s ", entity, err.Error())
			middleware.ReturnServerHandleError(c, err)
			return
		}
		if err = json.Unmarshal(responseBytes, &response); err != nil {
			middleware.ReturnServerHandleError(c, fmt.Errorf("json Unmarshal err,%s", err.Error()))
			return
		}
		if response.Status == "OK" {
			for _, dataMap := range response.Data {
				for _, v1 := range refAttributes {
					if strings.ToUpper(v1.Sensitive) == "Y" {
						if _, ok := dataMap[v1.PropertyName]; ok {
							dataMap[v1.PropertyName] = "******"
						}
					}
				}
			}
		}
	}
	c.Data(http.StatusOK, "application/json; charset=utf-8", responseBytes)
}
