package models

const (
	DateTimeFormat       = "2006-01-02 15:04:05"
	NewDateTimeFormat    = "20060102150405"
	SysTableIdConnector  = "__"
	UrlPrefix            = "/taskman"
	RowDataPermissionErr = "row data permission deny "
	AdminRole            = "SUPER_ADMIN"
	UploadFileMaxSize    = 10485760
	DefaultHttpErrorCode = "ERROR"
	WeCubeEmptySearch    = "WeCube-empty-search" //查询空数据
	Yes                  = "yes"
	Y                    = "y"
)

// RequestStatus 定义请求状态
type RequestStatus string

const (
	RequestStatusDraft               RequestStatus = "Draft"                 // 草稿
	RequestStatusPending             RequestStatus = "Pending"               // 等待定版
	RequestStatusInProgress          RequestStatus = "InProgress"            // 执行中
	RequestStatusInProgressFaulted   RequestStatus = "InProgress(Faulted)"   // 节点报错
	RequestStatusTermination         RequestStatus = "Termination"           // 手动终止
	RequestStatusCompleted           RequestStatus = "Completed"             // 成功
	RequestStatusInProgressTimeOuted RequestStatus = "InProgress(Timeouted)" // 节点超时
	RequestStatusFaulted             RequestStatus = "Faulted"               // 自动退出
	RequestStatusInApproval          RequestStatus = "InApproval"            // 审批中
	RequestStatusConfirm             RequestStatus = "Confirm"               // 请求确认
)

const ProcDefStatusTimeout = "Timeouted" //编排状态超时

// RequestTemplateStatus 请求模板状态
type RequestTemplateStatus string

const (
	RequestTemplateStatusCreated  RequestTemplateStatus = "created" // 创建
	RequestTemplateStatusDisabled RequestTemplateStatus = "disable" // 禁用,所有版本不可用
	RequestTemplateStatusCancel   RequestTemplateStatus = "cancel"  // 作废,当前版本失效
	RequestTemplateStatusPending  RequestTemplateStatus = "pending" // 待发布
	RequestTemplateStatusConfirm  RequestTemplateStatus = "confirm" // 已发布
)

// RolePermission 角色权限
type RolePermission string

const (
	RolePermissionUse  RolePermission = "USE"
	RolePermissionMGMT RolePermission = "MGMT"
)

// ProcDefNodeType 编排节点类型
type ProcDefNodeType string

const (
	ProcDefNodeTypeHuman     ProcDefNodeType = "subProcess" //人工节点
	ProcDefNodeTypeAutomatic ProcDefNodeType = "automatic"  //自动节点
	ProcDefNodeTypeData      ProcDefNodeType = "data"       //数据节点
)

// FormItemGroupType 表单组类型
type FormItemGroupType string

const (
	FormItemGroupTypeWorkflow FormItemGroupType = "workflow"
	FormItemGroupTypeOptional FormItemGroupType = "optional"
	FormItemGroupTypeCustom   FormItemGroupType = "custom"
)

// FormItemElementType 表单项类型
type FormItemElementType string

const (
	FormItemElementTypeInput     FormItemElementType = "input"
	FormItemElementTypeSelect    FormItemElementType = "select"
	FormItemElementTypeCalculate FormItemElementType = "calculate" //计算类型
)

// TaskType 任务模板、任务 类型
type TaskType string

const (
	TaskTypeNone      TaskType = ""          // 空
	TaskTypeSubmit    TaskType = "submit"    // 提交
	TaskTypeCheck     TaskType = "check"     // 定版
	TaskTypeApprove   TaskType = "approve"   // 审批
	TaskTypeImplement TaskType = "implement" // 执行类型(任务)
	TaskTypeConfirm   TaskType = "confirm"   // 请求确认
	TaskTypeRevoke    TaskType = "revoke"    // 请求撤回
)

// TaskTemplateHandleMode 任务模板 处理模式
type TaskTemplateHandleMode string

const (
	TaskTemplateHandleModeCustom TaskTemplateHandleMode = "custom" // 单人自定义
	TaskTemplateHandleModeAny    TaskTemplateHandleMode = "any"    // 协同
	TaskTemplateHandleModeAll    TaskTemplateHandleMode = "all"    // 并行
	TaskTemplateHandleModeAdmin  TaskTemplateHandleMode = "admin"  // 提交人角色管理员
	TaskTemplateHandleModeAuto   TaskTemplateHandleMode = "auto"   // 自动通过
)

// TaskHandleTemplateAssignType 任务处理模板 分派方式
type TaskHandleTemplateAssignType string

const (
	TaskHandleTemplateAssignTypeTemplate TaskHandleTemplateAssignType = "template" // 模板指定
	TaskHandleTemplateAssignTypeCustom   TaskHandleTemplateAssignType = "custom"   // 提交人指定
)

// TaskHandleTemplateHandlerType 任务处理模板 人员设置方式
type TaskHandleTemplateHandlerType string

const (
	TaskHandleTemplateHandlerTypeTemplate        TaskHandleTemplateHandlerType = "template"         // 模板指定
	TaskHandleTemplateHandlerTypeTemplateSuggest TaskHandleTemplateHandlerType = "template_suggest" // 模板建议
	TaskHandleTemplateHandlerTypeCustom          TaskHandleTemplateHandlerType = "custom"           // 提交人指定
	TaskHandleTemplateHandlerTypeCustomSuggest   TaskHandleTemplateHandlerType = "custom_suggest"   // 提交人建议
	TaskHandleTemplateHandlerTypeSystem          TaskHandleTemplateHandlerType = "system"           // 组内系统分配
	TaskHandleTemplateHandlerTypeClaim           TaskHandleTemplateHandlerType = "claim"            // 组内主动认领
)

// TaskHandleResultType 任务处理 处理结果
type TaskHandleResultType string

const (
	TaskHandleResultTypeApprove     TaskHandleResultType = "approve"     // 同意
	TaskHandleResultTypeDeny        TaskHandleResultType = "deny"        // 拒绝
	TaskHandleResultTypeRedraw      TaskHandleResultType = "redraw"      // 打回
	TaskHandleResultTypeComplete    TaskHandleResultType = "complete"    // 完成
	TaskHandleResultTypeUncompleted TaskHandleResultType = "uncompleted" // 未完成
)

// RequestFormType 请求表单类型
type RequestFormType string

const (
	RequestFormTypeMessage RequestFormType = "message" // 请求表单-信息表单
	RequestFormTypeData    RequestFormType = "data"    // 请求表单-数据表单
)

// TaskStatus 任务状态
type TaskStatus string

const (
	TaskStatusCreated TaskStatus = "created" // 任务创建
	TaskStatusMarked  TaskStatus = "marked"  // 任务认领
	TaskStatusDoing   TaskStatus = "doing"   // 任务进行中
	TaskStatusDone    TaskStatus = "done"    // 任务完成
)

const (
	TaskTypeImplementProcess string = "implement_process" // 编排任务
	TaskTypeImplementCustom  string = "implement_custom"  // 自定义任务
)

const (
	TaskConfirmResultComplete    string = "complete"    // 任务已完成
	TaskConfirmResultUncompleted string = "uncompleted" // 任务未完成
)

// SceneType 场景类型
type SceneType int

const (
	SceneTypeRelease SceneType = 1 // 发布
	SceneTypeRequest SceneType = 2 // 请求
	SceneTypeProblem SceneType = 3 // 问题
	SceneTypeEvent   SceneType = 4 // 事件
	SceneTypeChange  SceneType = 5 // 变更
)

type PlatTab string

const (
	PlatTabRequest        PlatTab = "Request" // 请求
	PlatTabRelease        PlatTab = "Release" // 发布
	PlatTabProblem        PlatTab = "Problem" // 问题
	PlatTabProblemEvent   PlatTab = "Event"   // 事件
	PlatTabProblemChange  PlatTab = "Change"  // 变更
	PlatTabProblemApprove PlatTab = "Approve" // 审批
	PlatTabProblemTask    PlatTab = "Task"    // 任务处理
	PlatTabProblemCheck   PlatTab = "Check"   // 请求定版
	PlatTabProblemConfirm PlatTab = "Confirm" // 请求确认
)

// TaskExecStatus 任务执行状态       // 状态值：1 进行中 2.未开始  3.已完成  4.报错被拒绝了
type TaskExecStatus int

const (
	TaskExecStatusDoing                TaskExecStatus = 1 // 进行中
	TaskExecStatusNotStart             TaskExecStatus = 2 // 未开始
	TaskExecStatusCompleted            TaskExecStatus = 3 // 已完成
	TaskExecStatusFail                 TaskExecStatus = 4 // 报错失败,被拒绝了
	TaskExecStatusAutoExitStatus       TaskExecStatus = 5 // 自动退出
	TaskExecStatusInternallyTerminated TaskExecStatus = 6 // 手动终止
)
