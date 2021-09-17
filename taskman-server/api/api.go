package api

import (
	"bytes"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/v1/form"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/v1/request"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/v1/task"
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
	httpHandlerFuncList = append(httpHandlerFuncList,
		&handlerFuncObj{Url: "/request-template-group/query", Method: "POST", HandlerFunc: request.QueryRequestTemplateGroup},
		&handlerFuncObj{Url: "/request-template-group", Method: "POST", HandlerFunc: request.CreateRequestTemplateGroup},
		&handlerFuncObj{Url: "/request-template-group", Method: "PUT", HandlerFunc: request.UpdateRequestTemplateGroup},
		&handlerFuncObj{Url: "/request-template-group", Method: "DELETE", HandlerFunc: request.DeleteRequestTemplateGroup},
		&handlerFuncObj{Url: "/process/list", Method: "GET", HandlerFunc: request.GetCoreProcessList},
		&handlerFuncObj{Url: "/process-nodes/:id", Method: "GET", HandlerFunc: request.GetCoreProcNodes},
		&handlerFuncObj{Url: "/role/list", Method: "GET", HandlerFunc: request.GetRoleList},
		&handlerFuncObj{Url: "/user/roles", Method: "GET", HandlerFunc: request.GetUserRoles},
		&handlerFuncObj{Url: "/request-template/query", Method: "POST", HandlerFunc: request.QueryRequestTemplate},
		&handlerFuncObj{Url: "/request-template", Method: "POST", HandlerFunc: request.CreateRequestTemplate},
		&handlerFuncObj{Url: "/request-template", Method: "PUT", HandlerFunc: request.UpdateRequestTemplate},
		&handlerFuncObj{Url: "/request-template", Method: "DELETE", HandlerFunc: request.DeleteRequestTemplate},
		&handlerFuncObj{Url: "/request-template/:id/attrs/update", Method: "PUT", HandlerFunc: request.UpdateRequestTemplateEntityAttrs},
		&handlerFuncObj{Url: "/request-template/:id/attrs/get", Method: "GET", HandlerFunc: request.GetRequestTemplateEntityAttrs},
		&handlerFuncObj{Url: "/request-template/:id/attrs/list", Method: "GET", HandlerFunc: request.ListRequestTemplateEntityAttrs},
		&handlerFuncObj{Url: "/request-template/confirm/:id", Method: "POST", HandlerFunc: form.ConfirmRequestFormTemplate},

		&handlerFuncObj{Url: "/request-form-template/:id", Method: "GET", HandlerFunc: form.GetRequestFormTemplate},
		&handlerFuncObj{Url: "/request-form-template/:id", Method: "POST", HandlerFunc: form.UpdateRequestFormTemplate},

		&handlerFuncObj{Url: "/task-template/:requestTemplateId/:proNodeId", Method: "GET", HandlerFunc: task.GetTaskTemplate},
		&handlerFuncObj{Url: "/task-template/:requestTemplateId", Method: "POST", HandlerFunc: task.UpdateTaskTemplate},

		&handlerFuncObj{Url: "/entity/data", Method: "POST", HandlerFunc: request.GetEntityData},
		&handlerFuncObj{Url: "/process/preview", Method: "GET", HandlerFunc: request.ProcessDataPreview},
	)
}

func InitHttpServer() {
	urlPrefix := models.UrlPrefix
	r := gin.New()
	// allow cross request
	if models.Config.HttpServer.Cross {
		crossHandler(r)
	}
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
