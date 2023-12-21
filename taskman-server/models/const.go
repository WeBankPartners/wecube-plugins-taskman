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
