package api

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
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
	ApiCode      string
}

var (
	httpHandlerFuncList   []*handlerFuncObj
	httpHandlerFuncV2List []*handlerFuncObj
	apiCodeMap            = make(map[string]string)
)

func init() {
	httpHandlerFuncList = []*handlerFuncObj{
		{Url: "/request-template-group/query", Method: "POST", HandlerFunc: template.QueryRequestTemplateGroup, ApiCode: "request-template-group-query"},
		{Url: "/request-template-group", Method: "POST", HandlerFunc: template.CreateRequestTemplateGroup, ApiCode: "request-template-group-post"},
		{Url: "/request-template-group", Method: "PUT", HandlerFunc: template.UpdateRequestTemplateGroup, ApiCode: "request-template-group-put"},
		{Url: "/request-template-group", Method: "DELETE", HandlerFunc: template.DeleteRequestTemplateGroup, ApiCode: "request-template-group-delete"},
		{Url: "/process/list", Method: "GET", HandlerFunc: workflow.GetCoreProcessList, ApiCode: "process-list-get"},
		{Url: "/process-nodes/:id/:type", Method: "GET", HandlerFunc: workflow.GetCoreProcNodes, ApiCode: "process-nodes-get"},
		{Url: "/role/list", Method: "GET", HandlerFunc: role.GetRoleList, ApiCode: "role-list-get"},
		{Url: "/role/user/list", Method: "GET", HandlerFunc: role.GetUserByRoles, ApiCode: "role-user-list-get"},
		{Url: "/role/administrator/list", Method: "GET", HandlerFunc: role.GetRoleAdministrators, ApiCode: "role-administrator-list-get"},
		{Url: "/user/roles", Method: "GET", HandlerFunc: role.GetUserRoles, ApiCode: "user-roles-get"},
		{Url: "/request-template/query", Method: "POST", HandlerFunc: template.QueryRequestTemplate, ApiCode: "request-template-query-post"},
		{Url: "/request-template/:id", Method: "GET", HandlerFunc: template.GetRequestTemplate, ApiCode: "request-template-get"},
		{Url: "/request-template/latest/all", Method: "GET", HandlerFunc: template.GetAllLatestReleaseRequestTemplate, ApiCode: "request-template-all-get"},
		{Url: "/request-template/roles", Method: "POST", HandlerFunc: template.GetRequestTemplateRoles, ApiCode: "request-template-roles-post"},
		{Url: "/request-template", Method: "POST", HandlerFunc: template.CreateRequestTemplate, ApiCode: "request-template-post"},
		{Url: "/request-template", Method: "PUT", HandlerFunc: template.UpdateRequestTemplate, ApiCode: "request-template-put"},
		{Url: "/request-template", Method: "DELETE", HandlerFunc: template.DeleteRequestTemplate, ApiCode: "request-template-delete"},
		{Url: "/request-template/handler/update", Method: "POST", HandlerFunc: template.UpdateRequestTemplateHandler, ApiCode: "request-template-handler-update-post"},
		{Url: "/request-template/status/update", Method: "POST", HandlerFunc: template.UpdateRequestTemplateStatus, ApiCode: "request-template-status-update-post"},
		{Url: "/request-template/:id/attrs/update", Method: "PUT", HandlerFunc: template.UpdateRequestTemplateEntityAttrs, ApiCode: "request-template-attrs-update-put"},
		{Url: "/request-template/:id/attrs/get", Method: "GET", HandlerFunc: template.GetRequestTemplateEntityAttrs, ApiCode: "request-template-attrs-get"},
		{Url: "/request-template/:id/attrs/list", Method: "GET", HandlerFunc: template.ListRequestTemplateEntityAttrs, ApiCode: "request-template-attrs-list-get"},
		{Url: "/request-template/:id/entity", Method: "GET", HandlerFunc: template.QueryRequestTemplateEntity, ApiCode: "request-template-entity-get"},
		{Url: "/request-template/confirm/:id", Method: "POST", HandlerFunc: form.ConfirmRequestFormTemplate, ApiCode: "request-template-confirm-post"},
		{Url: "/request-template/confirm_count", Method: "GET", HandlerFunc: template.GetConfirmCount, ApiCode: "request-template-confirm-count-get"},
		{Url: "/request-template/fork/:id", Method: "POST", HandlerFunc: template.ForkConfirmRequestTemplate, ApiCode: "request-template-fork-post"},
		{Url: "/request-template/copy/:id", Method: "POST", HandlerFunc: template.CopyConfirmRequestTemplate, ApiCode: "request-template-copy-post"},
		{Url: "/request-template/tags/:requestTemplateGroup", Method: "GET", HandlerFunc: template.GetRequestTemplateTags, ApiCode: "request-template-tags-get"},
		{Url: "/request-template/export/:requestTemplateId", Method: "GET", HandlerFunc: template.ExportRequestTemplate, ApiCode: "request-template-export-get"},
		{Url: "/request-template/export/batch", Method: "POST", HandlerFunc: template.BatchExportRequestTemplate, ApiCode: "request-template-export-batch-post"},
		{Url: "/request-template/import", Method: "POST", HandlerFunc: template.ImportRequestTemplate, ApiCode: "request-template-import-post"},
		{Url: "/request-template/import-batch", Method: "POST", HandlerFunc: template.ImportRequestTemplateBatch, ApiCode: "request-template-import-batch-post"},
		{Url: "/request-template/import-confirm/:confirmToken", Method: "POST", HandlerFunc: template.ConfirmImportRequestTemplate, ApiCode: "request-template-import-confirm-post"},
		{Url: "/request-template/disable/:id", Method: "POST", HandlerFunc: template.DisableRequestTemplate, ApiCode: "request-template-disable-post"},
		{Url: "/request-template/enable/:id", Method: "POST", HandlerFunc: template.EnableRequestTemplate, ApiCode: "request-template-enable-post"},
		{Url: "/request-form-template/:id", Method: "GET", HandlerFunc: form.GetRequestFormTemplate, ApiCode: "request-form-template-get"},
		{Url: "/request-template/:id/data-form-clean", Method: "POST", HandlerFunc: form.CleanDataForm, ApiCode: "request-template-data-form-clean-post"},
		{Url: "/request-template/:id/filter-clean/:type", Method: "POST", HandlerFunc: form.CleanFilterCondition, ApiCode: "request-template-filter-clean-post"},
		{Url: "/request-form-template/:id", Method: "POST", HandlerFunc: form.UpdateRequestFormTemplate, ApiCode: "request-form-template-post"},
		{Url: "/request-form-template/:id/data-form", Method: "GET", HandlerFunc: form.GetDataFormTemplate, ApiCode: "request-form-template-data-form-get"},
		{Url: "/request-form-template/:id/form/:task-template-id", Method: "GET", HandlerFunc: form.GetFormTemplate, ApiCode: "request-form-template-form-get"},
		{Url: "/request-form-template/:id/global-form", Method: "GET", HandlerFunc: form.GetGlobalFormEntity, ApiCode: "request-form-template-global-form-get"},
		{Url: "/form-template/item-group-config", Method: "POST", HandlerFunc: form.UpdateFormTemplateItemGroupConfig, ApiCode: "form-template-item-group-config-post"},
		{Url: "/form-template/item-group-config", Method: "GET", HandlerFunc: form.GetFormTemplateItemGroupConfig, ApiCode: "form-template-item-group-config-get"},
		{Url: "/form-template/item-group/copy", Method: "POST", HandlerFunc: form.CopyDataFormTemplateItemGroup, ApiCode: "form-template-item-group-copy-post"},
		{Url: "/form-template/item-group", Method: "DELETE", HandlerFunc: form.DeleteFormTemplateItemGroup, ApiCode: "form-template-item-group-delete"},
		{Url: "/form-template/item-group", Method: "POST", HandlerFunc: form.UpdateFormTemplateItemGroup, ApiCode: "form-template-item-group-post"},
		{Url: "/form-template/item-group/sort", Method: "POST", HandlerFunc: form.SortFormTemplateItemGroup, ApiCode: "form-template-item-group-sort-post"},

		{Url: "/form-template-library", Method: "POST", HandlerFunc: form.AddFormTemplateLibrary, ApiCode: "form-template-library-add"},
		{Url: "/form-template-library", Method: "DELETE", HandlerFunc: form.DeleteFormTemplateLibrary, ApiCode: "form-template-library-delete"},
		{Url: "/form-template-library/query", Method: "POST", HandlerFunc: form.QueryFormTemplateLibrary, ApiCode: "form-template-library-query"},
		{Url: "/form-template-library/form-type", Method: "GET", HandlerFunc: form.QueryAllFormTemplateLibraryFormType, ApiCode: "form-template-library-form-type-get"},
		{Url: "/form-template-library/export", Method: "POST", HandlerFunc: form.ExportFormTemplateLibrary, ApiCode: "form-template-library-export"},
		{Url: "/form-template-library/export-data", Method: "GET", HandlerFunc: form.ExportFormTemplateLibraryData, ApiCode: "form-template-library-export-data"},
		{Url: "/form-template-library/import", Method: "POST", HandlerFunc: form.ImportFormTemplateLibrary, ApiCode: "form-template-library-import"},

		{Url: "/task-template/:requestTemplate", Method: "POST", HandlerFunc: task.CreateTaskTemplate, ApiCode: "task-template-create"},
		{Url: "/task-template/:requestTemplate/:id", Method: "PUT", HandlerFunc: task.UpdateTaskTemplate, ApiCode: "task-template-update"},
		{Url: "/task-template/:requestTemplate/:id", Method: "DELETE", HandlerFunc: task.DeleteTaskTemplate, ApiCode: "task-template-delete"},
		{Url: "/task-template/form-template/:requestTemplate/:id", Method: "DELETE", HandlerFunc: task.DeleteTaskTemplateFormTemplate, ApiCode: "task-template-form-template-delete"},
		{Url: "/task-template/:requestTemplate/:id", Method: "GET", HandlerFunc: task.GetTaskTemplate, ApiCode: "task-template-get"},
		{Url: "/task-template/:requestTemplate/ids", Method: "GET", HandlerFunc: task.ListTaskTemplateIds, ApiCode: "task-template-ids-list"},
		{Url: "/task-template/:requestTemplate", Method: "GET", HandlerFunc: task.ListTaskTemplates, ApiCode: "task-template-list"},
		{Url: "/task-template/workflow/options", Method: "GET", HandlerFunc: task.GetTaskTemplateWorkFlowOptions, ApiCode: "task-template-workflow-options-get"},

		{Url: "/user/template/collect", Method: "POST", HandlerFunc: collect.AddTemplateCollect, ApiCode: "user-template-collect-add"},
		{Url: "/user/template/collect/:templateId", Method: "DELETE", HandlerFunc: collect.CancelTemplateCollect, ApiCode: "user-template-collect-cancel"},
		{Url: "/user/template/collect/query", Method: "POST", HandlerFunc: collect.QueryTemplateCollect, ApiCode: "user-template-collect-query"},
		{Url: "/user/template/filter-item", Method: "POST", HandlerFunc: collect.FilterItem, ApiCode: "user-template-filter-item"},

		{Url: "/entity/data", Method: "GET", HandlerFunc: workflow.GetEntityData, ApiCode: "entity-data-get"},
		{Url: "/models/package/:packageName/entity/:entity", Method: "GET", HandlerFunc: workflow.GetEntityModel, ApiCode: "models-package-entity"},
		{Url: "/packages/:pluginPackageId/entities/:entityName/query", Method: "POST", HandlerFunc: workflow.ProcEntityDataQuery, ApiCode: "packages-entities-query"},
		{Url: "/process/preview", Method: "GET", HandlerFunc: workflow.ProcessDataPreview, ApiCode: "process-preview"},

		{Url: "/request/:requestId", Method: "GET", HandlerFunc: request.GetRequest, ApiCode: "get-request-by-id"},
		{Url: "/request", Method: "POST", HandlerFunc: request.CreateRequest, ApiCode: "create-request"},
		{Url: "/request/:requestId", Method: "PUT", HandlerFunc: request.UpdateRequest, ApiCode: "update-request-by-id"},
		{Url: "/request/:requestId", Method: "DELETE", HandlerFunc: request.DeleteRequest, ApiCode: "delete-request-by-id"},
		{Url: "/request-parent/get", Method: "GET", HandlerFunc: request.GetRequestParent, ApiCode: "get-request-parent"},
		{Url: "/request/copy/:requestId", Method: "POST", HandlerFunc: request.CopyRequest, ApiCode: "copy-request-by-id"},
		{Url: "/request-root/:requestId", Method: "GET", HandlerFunc: request.GetRequestRootForm, ApiCode: "get-request-root-form-by-id"},
		{Url: "/request-data/preview", Method: "GET", HandlerFunc: request.GetRequestPreviewData, ApiCode: "preview-request-data"},
		{Url: "/request-data/save/:requestId/:cacheType", Method: "POST", HandlerFunc: request.SaveRequestCache, ApiCode: "save-request-cache"},
		{Url: "/request-data/get/:requestId/:cacheType", Method: "GET", HandlerFunc: request.GetRequestCache, ApiCode: "get-request-cache"},
		{Url: "/request-status/:requestId/:status", Method: "POST", HandlerFunc: request.UpdateRequestStatus, ApiCode: "update-request-status"},
		{Url: "/request-data/reference/query/:formItemTemplateId/:requestId/:attrName", Method: "POST", HandlerFunc: request.GetReferenceData, ApiCode: "query-reference-data"},
		{Url: "/request-data/entity/expression/query/:formItemTemplateId/:rootDataId", Method: "GET", HandlerFunc: request.GetExpressionItemData, ApiCode: "query-expression-item-data"},
		{Url: "/request-data/form/sensitive-attr/query", Method: "POST", HandlerFunc: request.AttrSensitiveDataQuery, ApiCode: "get-request-sensitive-data"},

		{Url: "/user/platform/count", Method: "POST", HandlerFunc: request.CountPlatform, ApiCode: "count-platform"},
		{Url: "/user/platform/filter-item", Method: "POST", HandlerFunc: request.FilterItem, ApiCode: "filter-item"},
		{Url: "/user/platform/list", Method: "POST", HandlerFunc: request.DataList, ApiCode: "list-platform-data"},
		{Url: "/user/request/revoke/:requestId", Method: "POST", HandlerFunc: request.RevokeRequest, ApiCode: "revoke-request-by-id"},
		{Url: "/user/request/:permission", Method: "POST", HandlerFunc: request.ListRequest, ApiCode: "list-request-by-permission"},

		{Url: "/request/start/:requestId", Method: "POST", HandlerFunc: request.StartRequest, ApiCode: "start-request-by-id"},
		{Url: "/request/terminate/:requestId", Method: "POST", HandlerFunc: request.TerminateRequest, ApiCode: "terminate-request-by-id"},
		{Url: "/request/attach-file/upload/:requestId", Method: "POST", HandlerFunc: request.UploadRequestAttachFile, ApiCode: "upload-attach-file"},
		{Url: "/request/attach-file/download/:fileId", Method: "GET", HandlerFunc: request.DownloadAttachFile, ApiCode: "download-attach-file"},
		{Url: "/request/attach-file/remove/:fileId", Method: "DELETE", HandlerFunc: request.RemoveAttachFile, ApiCode: "remove-attach-file"},
		{Url: "/request/handler/:requestId/:latestUpdateTime", Method: "POST", HandlerFunc: request.UpdateRequestHandler, ApiCode: "update-request-handler"},
		{Url: "/request/progress", Method: "GET", HandlerFunc: request.GetRequestProgress, ApiCode: "get-request-progress"},
		{Url: "/request/process/definitions/:templateId", Method: "GET", HandlerFunc: workflow.GetProcessDefinitions, ApiCode: "get-process-definitions"},
		{Url: "/request/process/instances/:instanceId", Method: "GET", HandlerFunc: workflow.GetProcessInstance, ApiCode: "get-process-instance"},
		{Url: "/request/workflow/task_node/:procInstanceId/:nodeInstanceId", Method: "POST", HandlerFunc: workflow.GetProcDefTaskNodeContext, ApiCode: "get-proc-def-task-node-context"},
		{Url: "/request/history/list", Method: "POST", HandlerFunc: request.HistoryList, ApiCode: "list-request-history"},
		{Url: "/request/export", Method: "POST", HandlerFunc: request.Export, ApiCode: "export-request"},
		{Url: "/request/:requestId/task/list", Method: "GET", HandlerFunc: request.GetTaskList, ApiCode: "get-task-list-by-request-id"},
		{Url: "/request/confirm", Method: "POST", HandlerFunc: request.Confirm, ApiCode: "confirm-request"},
		{Url: "/request/association", Method: "POST", HandlerFunc: request.Association, ApiCode: "associate-request"},

		// For core 1:get task form template  2:create task
		{Url: "/plugin/task/create/meta", Method: "GET", HandlerFunc: task.GetTaskFormStruct, ApiCode: "plugin-task-create-meta-get"},
		{Url: "/plugin/task/create", Method: "POST", HandlerFunc: task.CreateTask, ApiCode: "plugin-task-create-post"},
		{Url: "/plugin/task/create/custom", Method: "POST", HandlerFunc: task.CreateTask, ApiCode: "plugin-task-create-custom-post"},

		{Url: "/task/save/:taskId", Method: "POST", HandlerFunc: task.SaveTaskForm, ApiCode: "task-save-post"},
		{Url: "/task/approve/:taskId", Method: "POST", HandlerFunc: task.ApproveTask, ApiCode: "task-approve-post"},
		{Url: "/task/status/:operation/:taskId/:latestUpdateTime", Method: "POST", HandlerFunc: task.ChangeTaskStatus, ApiCode: "task-status-change-post"},
		{Url: "/task-handle/update", Method: "POST", HandlerFunc: task.UpdateTaskHandle, ApiCode: "task-handle-update-post"},
		{Url: "/task/attach-file/:taskId/upload/:taskHandleId", Method: "POST", HandlerFunc: task.UploadTaskAttachFile, ApiCode: "task-attach-file-upload-post"},

		// 转发auth登录接口
		{Url: "/login/seed", Method: "GET", HandlerFunc: login.GetSeed, ApiCode: "login-seed-get"},
		{Url: "/login", Method: "POST", HandlerFunc: login.Login, ApiCode: "login-post"},
	}

	// v2 版本
	httpHandlerFuncV2List = []*handlerFuncObj{
		{Url: "/user/request-template", Method: "GET", HandlerFunc: requestNew.GetRequestTemplateByUser, ApiCode: "user-request-template"},
		{Url: "/request/detail", Method: "GET", HandlerFunc: requestNew.GetRequestDetail, ApiCode: "request-detail"},
		{Url: "/request", Method: "POST", HandlerFunc: requestNew.CreateRequest, ApiCode: "request-create"},
		{Url: "/request-data/save/:requestId/:cacheType/:event", Method: "POST", HandlerFunc: requestNew.SaveRequestCache, ApiCode: "request-data-save"},
		{Url: "/request-check/confirm/:requestId", Method: "POST", HandlerFunc: requestNew.CheckRequest, ApiCode: "request-check-confirm"},
		{Url: "/request/history/:requestId", Method: "GET", HandlerFunc: requestNew.GetRequestHistory, ApiCode: "request-history"},
		{Url: "/plugin/request/create", Method: "POST", HandlerFunc: requestNew.PluginCreateRequest, ApiCode: "plugin-request-create"},
		{Url: "/request-data/form/save/:requestId", Method: "POST", HandlerFunc: requestNew.SaveRequestFormData, ApiCode: "request-data-form-save"},
		{Url: "/request-data/form/password/decode", Method: "GET", HandlerFunc: requestNew.DecodeRequestFormDataPassword, ApiCode: "request-data-form-password-decode"},
		// 转发platform接口
		{Url: "/platform/models", Method: "GET", HandlerFunc: requestNew.GetPlatformAllModels, ApiCode: "platform-models"},
		{Url: "/platform/:package/entities/:entity/query", Method: "POST", HandlerFunc: requestNew.QueryPlatformEntityData, ApiCode: "platform-entity-query"},
		// 转发auth接口
		{Url: "/auth/roles", Method: "GET", HandlerFunc: requestNew.TransAuthGetApplyRoles, ApiCode: "auth-roles"},
		{Url: "/auth/roles/apply", Method: "POST", HandlerFunc: requestNew.TransAuthStartApply, ApiCode: "auth-roles-apply"},
		{Url: "/auth/roles/apply/byhandler", Method: "POST", HandlerFunc: requestNew.TransAuthGetProcessableList, ApiCode: "auth-roles-apply-byhandler"},
		{Url: "/auth/users", Method: "GET", HandlerFunc: requestNew.TransAuthGetAllUser, ApiCode: "auth-users"},
		{Url: "/auth/roles/:roleId/users", Method: "GET", HandlerFunc: requestNew.TransAuthGetUserByRole, ApiCode: "auth-roles-users"},
		{Url: "/auth/roles/:roleId/users/revoke", Method: "POST", HandlerFunc: requestNew.TransAuthRemoveUserFromRole, ApiCode: "auth-roles-users-revoke"},
		{Url: "/auth/roles/:roleId/users", Method: "POST", HandlerFunc: requestNew.TransAuthAddUserForRole, ApiCode: "auth-roles-users-add"},
		{Url: "/auth/roles/apply", Method: "PUT", HandlerFunc: requestNew.TransAuthHandleApplication, ApiCode: "auth-roles-apply-handle"},
		{Url: "/auth/roles/apply", Method: "DELETE", HandlerFunc: requestNew.TransApplyDelete, ApiCode: "auth-roles-apply-delete"},
		{Url: "/auth/roles/apply/byapplier", Method: "POST", HandlerFunc: requestNew.TransAuthGetApplyList, ApiCode: "auth-roles-apply-byapplier"},
		{Url: "/auth/users/register", Method: "POST", HandlerFunc: requestNew.TransAuthUserRegister, ApiCode: "auth-users-register"},
	}
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
		apiCodeMap[fmt.Sprintf("%s_%s%s", funcObj.Method, models.UrlPrefix+"/api/v1", funcObj.Url)] = funcObj.ApiCode
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
		apiCodeMap[fmt.Sprintf("%s_%s%s", funcObj.Method, models.UrlPrefix+"/api/v2", funcObj.Url)] = funcObj.ApiCode
	}

	// entity query
	r.POST(urlPrefix+"/entities/request/query", request.QueryProcDefEntity)
	middleware.InitApiMenuMap(apiCodeMap)
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
		apiCode := apiCodeMap[c.Request.Method+"_"+c.FullPath()]
		c.Writer.Header().Add("Api-Code", apiCode)
		c.Set(models.ContextApiCode, apiCode)
		c.Next()
		log.Info(c, log.LOGGER_ACCESS, zap.String("url", c.Request.RequestURI), zap.String("method", c.Request.Method), zap.Int("code", c.Writer.Status()), zap.String("operator", c.GetString("user")), zap.String("ip", middleware.GetRemoteIp(c)), zap.Float64("cost_ms", time.Since(start).Seconds()*1000), zap.String("body", string(bodyBytes)))
	}
}
