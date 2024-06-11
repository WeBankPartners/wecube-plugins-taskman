package middleware

import (
	"fmt"
	"github.com/WeBankPartners/go-common-lib/token"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func GetRequestUser(c *gin.Context) string {
	return c.GetString("user")
}

func GetRequestRoles(c *gin.Context) []string {
	return c.GetStringSlice("roles")
}

var (
	whiteListUrl = map[string]struct{}{
		models.UrlPrefix + "/api/v1/login/seed": {},
		models.UrlPrefix + "/api/v1/login":      {},
		//models.UrlPrefix + "/api/v2/auth/roles":          {},
		models.UrlPrefix + "/api/v2/auth/roles/apply":    {},
		models.UrlPrefix + "/api/v2/auth/users/register": {},
	}
)

func isWhiteListUrl(url string) (result bool) {
	if paramIndex := strings.Index(url, "?"); paramIndex > 0 {
		url = url[:paramIndex]
	}
	_, result = whiteListUrl[url]
	return
}

func AuthCoreRequestToken() gin.HandlerFunc {
	var index int
	return func(c *gin.Context) {
		if isWhiteListUrl(c.Request.RequestURI) {
			c.Next()
		} else {
			uri := c.Request.RequestURI
			if index = strings.Index(uri, "?"); index > 0 {
				// 非登录情况下放行,登录情况下需要鉴权初始化用户和角色
				if uri[:index] == "/api/v2/auth/roles" && strings.TrimSpace(c.GetHeader("Authorization")) == "" {
					c.Next()
					return
				}
			}
			err := authCoreRequest(c)
			if err != nil {
				log.Logger.Error("Validate core token fail", log.Error(err))
				c.JSON(http.StatusUnauthorized, models.EntityResponse{Status: "ERROR", Message: "Core token validate fail "})
				c.Abort()
			} else {
				c.Next()
			}
		}
	}
}

func authCoreRequest(c *gin.Context) error {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return fmt.Errorf("can not find Request Header Authorization ")
	}
	authToken, err := token.DecodeJwtToken(authHeader, models.Config.Wecube.JwtSigningKey)
	if err != nil {
		return err
	}
	if authToken.User == "" {
		return fmt.Errorf("token content is illegal,main message is empty ")
	}
	c.Set("user", strings.ReplaceAll(authToken.User, " ", ""))
	var roles []string
	for _, v := range authToken.Roles {
		roles = append(roles, strings.ReplaceAll(v, " ", ""))
	}
	c.Set("roles", roles)
	return nil
}
