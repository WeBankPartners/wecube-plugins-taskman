package api

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/gin-gonic/gin"
)

type handlerFuncObj struct {
	HandlerFunc  func(c *gin.Context)
	Method       string
	Url          string
	LogOperation bool
}

var (
	httpHandlerFuncList []*handlerFuncObj
)

func init() {
	// baseKey

}

func InitHttpServer() {
	urlPrefix := models.UrlPrefix
	r := gin.New()
	// allow cross request
	if models.Config.HttpServer.Cross {
		crossHandler(r)
	}
	r.Static(fmt.Sprintf("%s/fonts", urlPrefix), fmt.Sprintf("public%s/fonts", urlPrefix))
	// access log
	if models.Config.Log.AccessLogEnable {
		r.Use(httpLogHandle())
	}
	// register handler func with auth
	authRouter := r.Group(urlPrefix+"/api/v1", middleware.AuthCoreRequestToken())
	for _, funcObj := range httpHandlerFuncList {
		switch funcObj.Method {
		case "GET":
			authRouter.GET(funcObj.Url, funcObj.HandlerFunc)
			break
		case "POST":
			authRouter.POST(funcObj.Url, funcObj.HandlerFunc)
			break
		case "PUT":
			authRouter.PUT(funcObj.Url, funcObj.HandlerFunc)
			break
		case "DELETE":
			authRouter.DELETE(funcObj.Url, funcObj.HandlerFunc)
			break
		}
	}
	r.Run(":" + models.Config.HttpServer.Port)
}

func crossHandler(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		if c.GetHeader("Origin") != "" {
			c.Header("Access-Control-Allow-Origin", c.GetHeader("Origin"))
		}
	})
}

func httpLogHandle() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		bodyBytes, _ := ioutil.ReadAll(c.Request.Body)
		c.Request.Body.Close()
		c.Request.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
		c.Set("requestBody", string(bodyBytes))
		c.Next()
		log.AccessLogger.Info("request", log.String("url", c.Request.RequestURI), log.String("method", c.Request.Method), log.Int("code", c.Writer.Status()), log.String("operator", c.GetString("user")), log.String("ip", middleware.GetRemoteIp(c)), log.Float64("cost_ms", time.Now().Sub(start).Seconds()*1000), log.String("body", string(bodyBytes)))
	}
}
