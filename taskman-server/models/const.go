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
	ProcDefNodeTypeStart        ProcDefNodeType = "start"        //开始
	ProcDefNodeTypeEnd          ProcDefNodeType = "end"          //结束
	ProcDefNodeTypeAbnormal     ProcDefNodeType = "abnormal"     //异常
	ProcDefNodeTypeDecision     ProcDefNodeType = "decision"     //判断
	ProcDefNodeTypeFork         ProcDefNodeType = "fork"         //分流
	ProcDefNodeTypeMerge        ProcDefNodeType = "merge"        //汇聚
	ProcDefNodeTypeHuman        ProcDefNodeType = "human"        //人工节点
	ProcDefNodeTypeAutomatic    ProcDefNodeType = "automatic"    //自动节点
	ProcDefNodeTypeData         ProcDefNodeType = "data"         //数据节点
	ProcDefNodeTypeDate         ProcDefNodeType = "date"         //时间节点
	ProcDefNodeTypeTimeInterval ProcDefNodeType = "timeInterval" //时间间隔
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

// ApprovalTemplateRoleType 审批模板 分配方式
type ApprovalTemplateRoleType string

const (
	ApprovalTemplateRoleTypeCustom ApprovalTemplateRoleType = "custom" // 单人自定义
	ApprovalTemplateRoleTypeAny    ApprovalTemplateRoleType = "any"    // 协同
	ApprovalTemplateRoleTypeAll    ApprovalTemplateRoleType = "all"    // 并行
	ApprovalTemplateRoleTypeAdmin  ApprovalTemplateRoleType = "admin"  // 提交人角色管理员
	ApprovalTemplateRoleTypeAuto   ApprovalTemplateRoleType = "auto"   // 自动通过
)

// ApprovalTemplateRoleRoleType 审批处理模板 角色设置方式
type ApprovalTemplateRoleRoleType string

const (
	ApprovalTemplateRoleRoleTypeTemplate ApprovalTemplateRoleRoleType = "template" // 模板指定
	ApprovalTemplateRoleRoleTypeCustom   ApprovalTemplateRoleRoleType = "custom"   // 提交人指定
)

// ApprovalTemplateRoleHandlerType 审批处理模板 人员设置方式
type ApprovalTemplateRoleHandlerType string

const (
	ApprovalTemplateRoleHandlerTypeTemplate        ApprovalTemplateRoleHandlerType = "template"         // 模板指定
	ApprovalTemplateRoleHandlerTypeTemplateSuggest ApprovalTemplateRoleHandlerType = "template_suggest" // 模板建议
	ApprovalTemplateRoleHandlerTypeCustom          ApprovalTemplateRoleHandlerType = "custom"           // 提交人指定
	ApprovalTemplateRoleHandlerTypeCustomSuggest   ApprovalTemplateRoleHandlerType = "custom_suggest"   // 提交人建议
	ApprovalTemplateRoleHandlerTypeSystem          ApprovalTemplateRoleHandlerType = "system"           // 组内系统分配
	ApprovalTemplateRoleHandlerTypeClaim           ApprovalTemplateRoleHandlerType = "claim"            // 组内主动认领
)

// ApprovalRoleApprove 审批处理 是否同意
type ApprovalRoleApprove string

const (
	ApprovalRoleApproveInit    ApprovalRoleApprove = "init"    // 未处理
	ApprovalRoleApproveApprove ApprovalRoleApprove = "approve" // 同意
	ApprovalRoleApproveDeny    ApprovalRoleApprove = "deny"    // 拒绝
)

// TaskTemplateType 任务模板 类型
type TaskTemplateType string

const (
	TaskTemplateTypeProc   TaskTemplateType = "proc"   // 编排
	TaskTemplateTypeCustom TaskTemplateType = "custom" // 自定义
)

// TaskTemplateRoleType 任务模板 分配方式
type TaskTemplateRoleType string

const (
	TaskTemplateRoleTypeCustom TaskTemplateRoleType = "custom" // 单人自定义
	TaskTemplateRoleTypeAdmin  TaskTemplateRoleType = "admin"  // 提交人角色管理员
)

// TaskTemplateRoleRoleType 任务处理模板 角色设置方式
type TaskTemplateRoleRoleType string

const (
	TaskTemplateRoleRoleTypeTemplate TaskTemplateRoleRoleType = "template" // 模板指定
	TaskTemplateRoleRoleTypeCustom   TaskTemplateRoleRoleType = "custom"   // 提交人指定
)

// TaskTemplateRoleHandlerType 任务处理模板 人员设置方式
type TaskTemplateRoleHandlerType string

const (
	TaskTemplateRoleHandlerTypeTemplate        TaskTemplateRoleHandlerType = "template"         // 模板指定
	TaskTemplateRoleHandlerTypeTemplateSuggest TaskTemplateRoleHandlerType = "template_suggest" // 模板建议
	TaskTemplateRoleHandlerTypeCustom          TaskTemplateRoleHandlerType = "custom"           // 提交人指定
	TaskTemplateRoleHandlerTypeCustomSuggest   TaskTemplateRoleHandlerType = "custom_suggest"   // 提交人建议
	TaskTemplateRoleHandlerTypeSystem          TaskTemplateRoleHandlerType = "system"           // 组内系统分配
	TaskTemplateRoleHandlerTypeClaim           TaskTemplateRoleHandlerType = "claim"            // 组内主动认领
)
