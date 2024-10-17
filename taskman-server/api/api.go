package api

import (
	"bytes"
	"io"
	"time"

	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/middleware"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/v1/collect"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/v1/form"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/v1/login"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/v1/request"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/v1/role"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/v1/task"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/v1/template"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/v1/workflow"
	requestNew "github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/api/v2/request"
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
	httpHandlerFuncList   []*handlerFuncObj
	httpHandlerFuncV2List []*handlerFuncObj
)

func init() {
	httpHandlerFuncList = append(httpHandlerFuncList,
		&handlerFuncObj{Url: "/request-template-group/query", Method: "POST", HandlerFunc: template.QueryRequestTemplateGroup},
		&handlerFuncObj{Url: "/request-template-group", Method: "POST", HandlerFunc: template.CreateRequestTemplateGroup},
		&handlerFuncObj{Url: "/request-template-group", Method: "PUT", HandlerFunc: template.UpdateRequestTemplateGroup},
		&handlerFuncObj{Url: "/request-template-group", Method: "DELETE", HandlerFunc: template.DeleteRequestTemplateGroup},
		&handlerFuncObj{Url: "/process/list", Method: "GET", HandlerFunc: workflow.GetCoreProcessList},
		&handlerFuncObj{Url: "/process-nodes/:id/:type", Method: "GET", HandlerFunc: workflow.GetCoreProcNodes},
		&handlerFuncObj{Url: "/role/list", Method: "GET", HandlerFunc: role.GetRoleList},
		&handlerFuncObj{Url: "/role/user/list", Method: "GET", HandlerFunc: role.GetUserByRoles},
		&handlerFuncObj{Url: "/role/administrator/list", Method: "GET", HandlerFunc: role.GetRoleAdministrators},
		&handlerFuncObj{Url: "/user/roles", Method: "GET", HandlerFunc: role.GetUserRoles},
		&handlerFuncObj{Url: "/request-template/query", Method: "POST", HandlerFunc: template.QueryRequestTemplate},
		&handlerFuncObj{Url: "/request-template/all", Method: "GET", HandlerFunc: template.GetAllLatestReleaseRequestTemplate},
		&handlerFuncObj{Url: "/request-template/roles", Method: "POST", HandlerFunc: template.GetRequestTemplateRoles},
		&handlerFuncObj{Url: "/request-template", Method: "POST", HandlerFunc: template.CreateRequestTemplate},
		&handlerFuncObj{Url: "/request-template", Method: "PUT", HandlerFunc: template.UpdateRequestTemplate},
		&handlerFuncObj{Url: "/request-template", Method: "DELETE", HandlerFunc: template.DeleteRequestTemplate},
		&handlerFuncObj{Url: "/request-template/handler/update", Method: "POST", HandlerFunc: template.UpdateRequestTemplateHandler},
		&handlerFuncObj{Url: "/request-template/status/update", Method: "POST", HandlerFunc: template.UpdateRequestTemplateStatus},
		&handlerFuncObj{Url: "/request-template/:id/attrs/update", Method: "PUT", HandlerFunc: template.UpdateRequestTemplateEntityAttrs},
		&handlerFuncObj{Url: "/request-template/:id/attrs/get", Method: "GET", HandlerFunc: template.GetRequestTemplateEntityAttrs},
		&handlerFuncObj{Url: "/request-template/:id/attrs/list", Method: "GET", HandlerFunc: template.ListRequestTemplateEntityAttrs},
		&handlerFuncObj{Url: "/request-template/:id/entity", Method: "GET", HandlerFunc: template.QueryRequestTemplateEntity},
		&handlerFuncObj{Url: "/request-template/confirm/:id", Method: "POST", HandlerFunc: form.ConfirmRequestFormTemplate},
		&handlerFuncObj{Url: "/request-template/confirm_count", Method: "GET", HandlerFunc: template.GetConfirmCount},
		&handlerFuncObj{Url: "/request-template/fork/:id", Method: "POST", HandlerFunc: template.ForkConfirmRequestTemplate},
		&handlerFuncObj{Url: "/request-template/copy/:id", Method: "POST", HandlerFunc: template.CopyConfirmRequestTemplate},
		&handlerFuncObj{Url: "/request-template/tags/:requestTemplateGroup", Method: "GET", HandlerFunc: template.GetRequestTemplateTags},
		&handlerFuncObj{Url: "/request-template/export/:requestTemplateId", Method: "GET", HandlerFunc: template.ExportRequestTemplate},
		&handlerFuncObj{Url: "/request-template/export/batch", Method: "POST", HandlerFunc: template.BatchExportRequestTemplate},
		&handlerFuncObj{Url: "/request-template/import", Method: "POST", HandlerFunc: template.ImportRequestTemplate},
		&handlerFuncObj{Url: "/request-template/import-batch", Method: "POST", HandlerFunc: template.ImportRequestTemplateBatch},
		&handlerFuncObj{Url: "/request-template/import-confirm/:confirmToken", Method: "POST", HandlerFunc: template.ConfirmImportRequestTemplate},
		&handlerFuncObj{Url: "/request-template/disable/:id", Method: "POST", HandlerFunc: template.DisableRequestTemplate},
		&handlerFuncObj{Url: "/request-template/enable/:id", Method: "POST", HandlerFunc: template.EnableRequestTemplate},
		&handlerFuncObj{Url: "/request-form-template/:id", Method: "GET", HandlerFunc: form.GetRequestFormTemplate},
		&handlerFuncObj{Url: "/request-template/:id/data-form-clean", Method: "POST", HandlerFunc: form.CleanDataForm},
		&handlerFuncObj{Url: "/request-template/:id/filter-clean/:type", Method: "POST", HandlerFunc: form.CleanFilterCondition},
		&handlerFuncObj{Url: "/request-form-template/:id", Method: "POST", HandlerFunc: form.UpdateRequestFormTemplate},
		&handlerFuncObj{Url: "/request-form-template/:id/data-form", Method: "GET", HandlerFunc: form.GetDataFormTemplate},
		&handlerFuncObj{Url: "/request-form-template/:id/form/:task-template-id", Method: "GET", HandlerFunc: form.GetFormTemplate},
		&handlerFuncObj{Url: "/request-form-template/:id/global-form", Method: "GET", HandlerFunc: form.GetGlobalFormEntity},
		&handlerFuncObj{Url: "/form-template/item-group-config", Method: "POST", HandlerFunc: form.UpdateFormTemplateItemGroupConfig},
		&handlerFuncObj{Url: "/form-template/item-group-config", Method: "GET", HandlerFunc: form.GetFormTemplateItemGroupConfig},
		&handlerFuncObj{Url: "/form-template/item-group/copy", Method: "POST", HandlerFunc: form.CopyDataFormTemplateItemGroup},
		&handlerFuncObj{Url: "/form-template/item-group", Method: "DELETE", HandlerFunc: form.DeleteFormTemplateItemGroup},
		&handlerFuncObj{Url: "/form-template/item-group", Method: "POST", HandlerFunc: form.UpdateFormTemplateItemGroup},
		&handlerFuncObj{Url: "/form-template/item-group/sort", Method: "POST", HandlerFunc: form.SortFormTemplateItemGroup},

		&handlerFuncObj{Url: "/form-template-library", Method: "POST", HandlerFunc: form.AddFormTemplateLibrary},
		&handlerFuncObj{Url: "/form-template-library", Method: "DELETE", HandlerFunc: form.DeleteFormTemplateLibrary},
		&handlerFuncObj{Url: "/form-template-library/query", Method: "POST", HandlerFunc: form.QueryFormTemplateLibrary},
		&handlerFuncObj{Url: "/form-template-library/form-type", Method: "GET", HandlerFunc: form.QueryAllFormTemplateLibraryFormType},
		&handlerFuncObj{Url: "/form-template-library/export", Method: "POST", HandlerFunc: form.ExportFormTemplateLibrary},
		&handlerFuncObj{Url: "/form-template-library/export-data", Method: "GET", HandlerFunc: form.ExportFormTemplateLibraryData},
		&handlerFuncObj{Url: "/form-template-library/import", Method: "POST", HandlerFunc: form.ImportFormTemplateLibrary},

		&handlerFuncObj{Url: "/task-template/:requestTemplate", Method: "POST", HandlerFunc: task.CreateTaskTemplate},
		&handlerFuncObj{Url: "/task-template/:requestTemplate/:id", Method: "PUT", HandlerFunc: task.UpdateTaskTemplate},
		&handlerFuncObj{Url: "/task-template/:requestTemplate/:id", Method: "DELETE", HandlerFunc: task.DeleteTaskTemplate},
		&handlerFuncObj{Url: "/task-template/form-template/:requestTemplate/:id", Method: "DELETE", HandlerFunc: task.DeleteTaskTemplateFormTemplate},
		&handlerFuncObj{Url: "/task-template/:requestTemplate/:id", Method: "GET", HandlerFunc: task.GetTaskTemplate},
		&handlerFuncObj{Url: "/task-template/:requestTemplate/ids", Method: "GET", HandlerFunc: task.ListTaskTemplateIds},
		&handlerFuncObj{Url: "/task-template/:requestTemplate", Method: "GET", HandlerFunc: task.ListTaskTemplates},
		&handlerFuncObj{Url: "/task-template/workflow/options", Method: "GET", HandlerFunc: task.GetTaskTemplateWorkFlowOptions},

		&handlerFuncObj{Url: "/user/template/collect", Method: "POST", HandlerFunc: collect.AddTemplateCollect},
		&handlerFuncObj{Url: "/user/template/collect/:templateId", Method: "DELETE", HandlerFunc: collect.CancelTemplateCollect},
		&handlerFuncObj{Url: "/user/template/collect/query", Method: "POST", HandlerFunc: collect.QueryTemplateCollect},
		&handlerFuncObj{Url: "/user/template/filter-item", Method: "POST", HandlerFunc: collect.FilterItem},

		&handlerFuncObj{Url: "/entity/data", Method: "GET", HandlerFunc: workflow.GetEntityData},
		&handlerFuncObj{Url: "/models/package/:packageName/entity/:entity", Method: "GET", HandlerFunc: workflow.GetEntityModel},
		&handlerFuncObj{Url: "/packages/:pluginPackageId/entities/:entityName/query", Method: "POST", HandlerFunc: workflow.ProcEntityDataQuery},
		&handlerFuncObj{Url: "/process/preview", Method: "GET", HandlerFunc: workflow.ProcessDataPreview},

		&handlerFuncObj{Url: "/request/:requestId", Method: "GET", HandlerFunc: request.GetRequest},
		&handlerFuncObj{Url: "/request", Method: "POST", HandlerFunc: request.CreateRequest},
		&handlerFuncObj{Url: "/request/:requestId", Method: "PUT", HandlerFunc: request.UpdateRequest},
		&handlerFuncObj{Url: "/request/:requestId", Method: "DELETE", HandlerFunc: request.DeleteRequest},
		&handlerFuncObj{Url: "/request-parent/get", Method: "GET", HandlerFunc: request.GetRequestParent},
		&handlerFuncObj{Url: "/request/copy/:requestId", Method: "POST", HandlerFunc: request.CopyRequest},
		&handlerFuncObj{Url: "/request-root/:requestId", Method: "GET", HandlerFunc: request.GetRequestRootForm},
		&handlerFuncObj{Url: "/request-data/preview", Method: "GET", HandlerFunc: request.GetRequestPreviewData},
		&handlerFuncObj{Url: "/request-data/save/:requestId/:cacheType", Method: "POST", HandlerFunc: request.SaveRequestCache},
		&handlerFuncObj{Url: "/request-data/get/:requestId/:cacheType", Method: "GET", HandlerFunc: request.GetRequestCache},
		&handlerFuncObj{Url: "/request-status/:requestId/:status", Method: "POST", HandlerFunc: request.UpdateRequestStatus},
		&handlerFuncObj{Url: "/request-data/reference/query/:formItemTemplateId/:requestId/:attrName", Method: "POST", HandlerFunc: request.GetReferenceData},
		&handlerFuncObj{Url: "/request-data/entity/expression/query/:formItemTemplateId/:rootDataId", Method: "GET", HandlerFunc: request.GetExpressionItemData},

		&handlerFuncObj{Url: "/user/platform/count", Method: "POST", HandlerFunc: request.CountPlatform},
		&handlerFuncObj{Url: "/user/platform/filter-item", Method: "POST", HandlerFunc: request.FilterItem},
		&handlerFuncObj{Url: "/user/platform/list", Method: "POST", HandlerFunc: request.DataList},
		&handlerFuncObj{Url: "/user/request/revoke/:requestId", Method: "POST", HandlerFunc: request.RevokeRequest},
		&handlerFuncObj{Url: "/user/request/:permission", Method: "POST", HandlerFunc: request.ListRequest},

		&handlerFuncObj{Url: "/request/start/:requestId", Method: "POST", HandlerFunc: request.StartRequest},
		&handlerFuncObj{Url: "/request/terminate/:requestId", Method: "POST", HandlerFunc: request.TerminateRequest},
		&handlerFuncObj{Url: "/request/attach-file/upload/:requestId", Method: "POST", HandlerFunc: request.UploadRequestAttachFile},
		&handlerFuncObj{Url: "/request/attach-file/download/:fileId", Method: "GET", HandlerFunc: request.DownloadAttachFile},
		&handlerFuncObj{Url: "/request/attach-file/remove/:fileId", Method: "DELETE", HandlerFunc: request.RemoveAttachFile},
		&handlerFuncObj{Url: "/request/handler/:requestId/:latestUpdateTime", Method: "POST", HandlerFunc: request.UpdateRequestHandler},
		&handlerFuncObj{Url: "/request/progress", Method: "GET", HandlerFunc: request.GetRequestProgress},
		&handlerFuncObj{Url: "/request/process/definitions/:templateId", Method: "GET", HandlerFunc: workflow.GetProcessDefinitions},
		&handlerFuncObj{Url: "/request/process/instances/:instanceId", Method: "GET", HandlerFunc: workflow.GetProcessInstance},
		&handlerFuncObj{Url: "/request/workflow/task_node/:procInstanceId/:nodeInstanceId", Method: "POST", HandlerFunc: workflow.GetProcDefTaskNodeContext},
		&handlerFuncObj{Url: "/request/history/list", Method: "POST", HandlerFunc: request.HistoryList},
		&handlerFuncObj{Url: "/request/export", Method: "POST", HandlerFunc: request.Export},
		&handlerFuncObj{Url: "/request/:requestId/task/list", Method: "GET", HandlerFunc: request.GetTaskList},
		&handlerFuncObj{Url: "/request/confirm", Method: "POST", HandlerFunc: request.Confirm},
		&handlerFuncObj{Url: "/request/association", Method: "POST", HandlerFunc: request.Association},

		// For core 1:get task form template  2:create task
		&handlerFuncObj{Url: "/plugin/task/create/meta", Method: "GET", HandlerFunc: task.GetTaskFormStruct},
		&handlerFuncObj{Url: "/plugin/task/create", Method: "POST", HandlerFunc: task.CreateTask},
		&handlerFuncObj{Url: "/plugin/task/create/custom", Method: "POST", HandlerFunc: task.CreateTask},

		&handlerFuncObj{Url: "/task/save/:taskId", Method: "POST", HandlerFunc: task.SaveTaskForm},
		&handlerFuncObj{Url: "/task/approve/:taskId", Method: "POST", HandlerFunc: task.ApproveTask},
		&handlerFuncObj{Url: "/task/status/:operation/:taskId/:latestUpdateTime", Method: "POST", HandlerFunc: task.ChangeTaskStatus},
		&handlerFuncObj{Url: "/task-handle/update", Method: "POST", HandlerFunc: task.UpdateTaskHandle},
		&handlerFuncObj{Url: "/task/attach-file/:taskId/upload/:taskHandleId", Method: "POST", HandlerFunc: task.UploadTaskAttachFile},

		// 转发auth登录接口
		&handlerFuncObj{Url: "/login/seed", Method: "GET", HandlerFunc: login.GetSeed},
		&handlerFuncObj{Url: "/login", Method: "POST", HandlerFunc: login.Login},
	)

	// v2 版本
	httpHandlerFuncV2List = append(httpHandlerFuncV2List,
		&handlerFuncObj{Url: "/user/request-template", Method: "GET", HandlerFunc: requestNew.GetRequestTemplateByUser},
		&handlerFuncObj{Url: "/request/detail", Method: "GET", HandlerFunc: requestNew.GetRequestDetail},
		&handlerFuncObj{Url: "/request", Method: "POST", HandlerFunc: requestNew.CreateRequest},
		&handlerFuncObj{Url: "/request-data/save/:requestId/:cacheType/:event", Method: "POST", HandlerFunc: requestNew.SaveRequestCache},
		&handlerFuncObj{Url: "/request-check/confirm/:requestId", Method: "POST", HandlerFunc: requestNew.CheckRequest}, // 确认定版
		&handlerFuncObj{Url: "/request/history/:requestId", Method: "GET", HandlerFunc: requestNew.GetRequestHistory},
		&handlerFuncObj{Url: "/plugin/request/create", Method: "POST", HandlerFunc: requestNew.PluginCreateRequest},
		&handlerFuncObj{Url: "/request-data/form/save/:requestId", Method: "POST", HandlerFunc: requestNew.SaveRequestFormData},
		// 转发platform接口
		&handlerFuncObj{Url: "/platform/models", Method: "GET", HandlerFunc: requestNew.GetPlatformAllModels},
		&handlerFuncObj{Url: "/platform/:package/entities/:entity/query", Method: "POST", HandlerFunc: requestNew.QueryPlatformEntityData},
		// 转发auth接口
		&handlerFuncObj{Url: "/auth/roles", Method: "GET", HandlerFunc: requestNew.TransAuthGetApplyRoles},
		&handlerFuncObj{Url: "/auth/roles/apply", Method: "POST", HandlerFunc: requestNew.TransAuthStartApply},
		&handlerFuncObj{Url: "/auth/roles/apply/byhandler", Method: "POST", HandlerFunc: requestNew.TransAuthGetProcessableList},
		&handlerFuncObj{Url: "/auth/users", Method: "GET", HandlerFunc: requestNew.TransAuthGetAllUser},
		&handlerFuncObj{Url: "/auth/roles/:roleId/users", Method: "GET", HandlerFunc: requestNew.TransAuthGetUserByRole},
		&handlerFuncObj{Url: "/auth/roles/:roleId/users/revoke", Method: "POST", HandlerFunc: requestNew.TransAuthRemoveUserFromRole},
		&handlerFuncObj{Url: "/auth/roles/:roleId/users", Method: "POST", HandlerFunc: requestNew.TransAuthAddUserForRole},
		&handlerFuncObj{Url: "/auth/roles/apply", Method: "PUT", HandlerFunc: requestNew.TransAuthHandleApplication},
		&handlerFuncObj{Url: "/auth/roles/apply", Method: "DELETE", HandlerFunc: requestNew.TransApplyDelete},
		&handlerFuncObj{Url: "/auth/roles/apply/byapplier", Method: "POST", HandlerFunc: requestNew.TransAuthGetApplyList},
		&handlerFuncObj{Url: "/auth/users/register", Method: "POST", HandlerFunc: requestNew.TransAuthUserRegister},
	)
}

func InitHttpServer() {
	middleware.InitHttpError()
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
		case "POST":
			authRouter.POST(funcObj.Url, funcObj.HandlerFunc)
		case "PUT":
			authRouter.PUT(funcObj.Url, funcObj.HandlerFunc)
		case "DELETE":
			authRouter.DELETE(funcObj.Url, funcObj.HandlerFunc)
		}
	}

	authRouterV2 := r.Group(urlPrefix+"/api/v2", middleware.AuthCoreRequestToken())
	for _, funcObj := range httpHandlerFuncV2List {
		switch funcObj.Method {
		case "GET":
			authRouterV2.GET(funcObj.Url, funcObj.HandlerFunc)
		case "POST":
			authRouterV2.POST(funcObj.Url, funcObj.HandlerFunc)
		case "PUT":
			authRouterV2.PUT(funcObj.Url, funcObj.HandlerFunc)
		case "DELETE":
			authRouterV2.DELETE(funcObj.Url, funcObj.HandlerFunc)
		}
	}

	// entity query
	r.POST(urlPrefix+"/entities/request/query", request.QueryProcDefEntity)
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
		bodyBytes, _ := io.ReadAll(c.Request.Body)
		c.Request.Body.Close()
		c.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		c.Set("requestBody", string(bodyBytes))
		c.Next()
		log.AccessLogger.Info("request", log.String("url", c.Request.RequestURI), log.String("method", c.Request.Method), log.Int("code", c.Writer.Status()), log.String("operator", c.GetString("user")), log.String("ip", middleware.GetRemoteIp(c)), log.Float64("cost_ms", time.Since(start).Seconds()*1000), log.String("body", string(bodyBytes)))
	}
}
