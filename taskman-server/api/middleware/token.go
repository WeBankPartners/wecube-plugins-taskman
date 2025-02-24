package middleware

import (
	"fmt"
	"github.com/WeBankPartners/go-common-lib/token"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
	"regexp"
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
	whitePathMap = map[string]bool{
		models.UrlPrefix + "/entities/${model}/query": true,
	}
	ApiMenuMap         = make(map[string][]string) // key -> apiCode  value -> menuList
	HomePageApiCodeMap = make(map[string]bool)
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
				log.Error(nil, log.LOGGER_APP, "Validate core token fail", zap.Error(err))
				c.JSON(http.StatusUnauthorized, models.EntityResponse{Status: "ERROR", Message: "Core token validate fail "})
				c.Abort()
			} else {
				// 白名单URL直接放行
				for path, _ := range whitePathMap {
					re := regexp.MustCompile(buildRegexPattern(path))
					if re.MatchString(c.Request.URL.Path) {
						c.Next()
						return
					}
				}
				contextApiCode := c.GetString(models.ContextApiCode)
				// 首页菜单直接放行
				if HomePageApiCodeMap[contextApiCode] {
					c.Next()
					return
				}
				if strings.HasPrefix(contextApiCode, "plugin-") {
					c.Next()
					return
				}
				if models.Config.MenuApiMap.Enable == "true" || strings.TrimSpace(models.Config.MenuApiMap.Enable) == "" || strings.ToUpper(models.Config.MenuApiMap.Enable) == "Y" {
					legal := false
					if allowMenuList, ok := ApiMenuMap[contextApiCode]; ok {
						legal = compareStringList(GetRequestRoles(c), allowMenuList)
					} else {
						legal = validateMenuApi(GetRequestRoles(c), c.Request.URL.Path, c.Request.Method)
					}
					if legal {
						c.Next()
					} else {
						ReturnApiPermissionError(c)
						c.Abort()
					}
				} else {
					c.Next()
				}
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

func ReadFormFile(c *gin.Context, fileKey string) (fileName string, fileBytes []byte, err error) {
	file, getFileErr := c.FormFile(fileKey)
	if getFileErr != nil {
		err = getFileErr
		return
	}
	fileName = file.Filename
	fileHandler, openFileErr := file.Open()
	if openFileErr != nil {
		err = openFileErr
		return
	}
	fileBytes, err = io.ReadAll(fileHandler)
	fileHandler.Close()
	return
}

func validateMenuApi(roles []string, path, method string) (legal bool) {
	// 防止ip 之类数据配置不上
	path = strings.ReplaceAll(path, ".", "")
	for _, menuApi := range models.MenuApiGlobalList {
		for _, role := range roles {
			if strings.ToLower(menuApi.Menu) == strings.ToLower(role) {
				for _, item := range menuApi.Urls {
					if strings.TrimSpace(item.Url) == "" {
						continue
					}
					if strings.ToLower(item.Method) == strings.ToLower(method) {
						re := regexp.MustCompile(buildRegexPattern(item.Url))
						if re.MatchString(path) {
							legal = true
							return
						}
					}
				}
			}
		}
	}
	return
}

func buildRegexPattern(template string) string {
	// 将 ${variable} 替换为 (\w+) ,并且严格匹配
	return "^" + regexp.MustCompile(`\$\{[\w.-]+\}`).ReplaceAllString(template, `(\w+)`) + "$"
}

func InitApiMenuMap(apiMenuCodeMap map[string]string) {
	var exist bool
	matchUrlMap := make(map[string]int)
	for k, code := range apiMenuCodeMap {
		exist = false
		re := regexp.MustCompile("^" + regexp.MustCompile(":[\\w\\-]+").ReplaceAllString(strings.ToLower(k), "([\\w\\.\\-\\$\\{\\}:\\[\\]]+)") + "$")
		for _, menuApi := range models.MenuApiGlobalList {
			for _, item := range menuApi.Urls {
				key := strings.ToLower(item.Method + "_" + item.Url)
				if re.MatchString(key) {
					exist = true
					if existList, existFlag := ApiMenuMap[code]; existFlag {
						ApiMenuMap[code] = append(existList, menuApi.Menu)
					} else {
						ApiMenuMap[code] = []string{menuApi.Menu}
					}
					if menuApi.Menu == models.HomePage {
						HomePageApiCodeMap[code] = true
					}
					matchUrlMap[item.Method+"_"+item.Url] = 1
				}
			}
		}
		if !exist {
			log.Info(nil, log.LOGGER_APP, "", zap.String("path", k), zap.String("code", code))
		}
	}
	for _, menuApi := range models.MenuApiGlobalList {
		for _, item := range menuApi.Urls {
			if _, ok := matchUrlMap[item.Method+"_"+item.Url]; !ok {
				//log.Info(nil, log.LOGGER_APP, "InitApiMenuMap can not match menuUrl", zap.String("menu", menuApi.Menu), zap.String("method", item.Method), zap.String("url", item.Url))
			}
		}
	}
	for k, v := range ApiMenuMap {
		if len(v) > 1 {
			ApiMenuMap[k] = DistinctStringList(v, []string{})
		}
	}
	log.Debug(nil, log.LOGGER_APP, "InitApiMenuMap done", log.JsonObj("ApiMenuMap", ApiMenuMap))
}

func DistinctStringList(input, excludeList []string) (output []string) {
	if len(input) == 0 {
		return
	}
	existMap := make(map[string]int)
	for _, v := range excludeList {
		existMap[v] = 1
	}
	for _, v := range input {
		if _, ok := existMap[v]; !ok {
			output = append(output, v)
			existMap[v] = 1
		}
	}
	return
}

func compareStringList(from, target []string) bool {
	match := false
	for _, f := range from {
		for _, t := range target {
			if f == t {
				match = true
				break
			}
		}
		if match {
			break
		}
	}
	return match
}
