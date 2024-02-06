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

// Status 定义请求状态
type Status string

const (
	Draft               Status = "Draft"                 // 草稿
	Pending             Status = "Pending"               // 等待定版
	InProgress          Status = "InProgress"            // 执行中
	InProgressFaulted   Status = "InProgress(Faulted)"   // 节点报错
	Termination         Status = "Termination"           // 手动终止
	Completed           Status = "Completed"             // 成功
	InProgressTimeOuted Status = "InProgress(Timeouted)" // 节点超时
	Faulted             Status = "Faulted"               // 自动退出
)

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
	RequestTemplateStatusCreated  RequestTemplateStatus = "created"  // 创建
	RequestTemplateStatusDisabled RequestTemplateStatus = "disabled" // 禁用
	RequestTemplateStatusPending  RequestTemplateStatus = "pending"  // 待发布
	RequestTemplateStatusConfirm  RequestTemplateStatus = "confirm"  // 已发布
)

// RolePermission 角色权限
type RolePermission string

const (
	RolePermissionUse  RolePermission = "USE"
	RolePermissionMGMT RolePermission = "MGMT"
)
