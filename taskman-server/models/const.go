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
