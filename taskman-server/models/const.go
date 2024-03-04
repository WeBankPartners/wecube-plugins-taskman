package models

const (
	DateTimeFormat       = "2006-01-02 15:04:05"
	SysTableIdConnector  = "__"
	UrlPrefix            = "/taskman"
	RowDataPermissionErr = "Row data permission deny "
	AdminRole            = "SUPER_ADMIN"
	UploadFileMaxSize    = 10485760
	DefaultHttpErrorCode = "ERROR"
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

// TemplateType 模板类型
type TemplateType string

const (
	RequestType TemplateType = "request" // 请求
	ReleaseType TemplateType = "release" // 发布
	OtherType   TemplateType = "other"   // 其他
)

// RequestTemplateStatus 请求模板状态
type RequestTemplateStatus string

const (
	RequestTemplateStatusCreated  RequestTemplateStatus = "created" // 创建
	RequestTemplateStatusDisabled RequestTemplateStatus = "disable" // 禁用
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
	TaskTypeSubmit    TaskType = "submit"    // 提交
	TaskTypeCheck     TaskType = "check"     // 定版
	TaskTypeApprove   TaskType = "approve"   // 审批
	TaskTypeImplement TaskType = "implement" // 执行类型(任务)
	TaskTypeConfirm   TaskType = "confirm"   // 请求确认
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

// TaskResultType 任务 处理结果
type TaskResultType string

const (
	TaskResultTypeComplete    TaskResultType = "complete"    // 完成
	TaskResultTypeUncompleted TaskResultType = "uncompleted" // 未完成
)

// TaskHandleResultType 任务处理 处理结果
type TaskHandleResultType string

const (
	TaskHandleResultTypeApprove    TaskHandleResultType = "approve"    // 同意
	TaskHandleResultTypeDeny       TaskHandleResultType = "deny"       // 拒绝
	TaskHandleResultTypeRedraw     TaskHandleResultType = "redraw"     // 打回
	TaskHandleResultTypeComplete   TaskHandleResultType = "complete"   // 完成
	TaskHandleResultTypeUncomplete TaskHandleResultType = "uncomplete" // 未完成
)

// TaskHandleChangeReason 任务处理 变更原因
type TaskHandleChangeReason string

const (
	TaskHandleChangeReasonAssign TaskHandleChangeReason = "assign" // 系统分配
	TaskHandleChangeReasonClaim  TaskHandleChangeReason = "claim"  // 主动领取
	TaskHandleChangeReasonGive   TaskHandleChangeReason = "give"   // 转给我
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

// ProgressStatus 请求状态
type ProgressStatus int

const (
	ProgressStatusInProgress                 ProgressStatus = 1 // 进行中
	ProgressStatusNotStart                   ProgressStatus = 2 // 未开始
	ProgressStatusCompleted                  ProgressStatus = 3 // 已完成
	ProgressStatusFail                       ProgressStatus = 4 // 报错失败,被拒绝了
	ProgressStatusAutoExitStatus             ProgressStatus = 5 // 自动退出
	ProgressStatusInternallyTerminatedStatus ProgressStatus = 6 // 自动退出
)

const (
	TaskTypeImplementProcess string = "implement_process" // 编排任务
	TaskTypeImplementCustom  string = "implement_custom"  // 自定义任务
)
