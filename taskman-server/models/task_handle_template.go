package models

// TaskHandleTemplateTable 任务模板处理表
type TaskHandleTemplateTable struct {
	Id           string `json:"id" xorm:"'id' pk" primary-key:"id"`
	Sort         int    `json:"sort" xorm:"sort"`
	TaskTemplate string `json:"taskTemplate" xorm:"task_template"`
	Role         string `json:"role" xorm:"role"`
	Assign       string `json:"assign" xorm:"assign"`
	HandlerType  string `json:"handlerType" xorm:"handler_type"` //template.模板指定 custom.提交人指定
	Handler      string `json:"handler" xorm:"handler"`
	HandleMode   string `json:"handleMode" xorm:"handle_mode"` // 处理模式：custom.单人自定义 any.协同 all.并行 admin.提交人角色管理员 auto.自动通过
	AssignRule   string `json:"assignRule" xorm:"assign_rule"` // 分配规则
	FilterRule   string `json:"filterRule" xorm:"filter_rule"` // 下拉框过滤规则
}

func (TaskHandleTemplateTable) TableName() string {
	return "task_handle_template"
}

type TaskHandleTemplateDto struct {
	Id          string                 `json:"id"`
	Role        string                 `json:"role"`
	Assign      string                 `json:"assign"`
	HandlerType string                 `json:"handlerType"`
	Handler     string                 `json:"handler"`
	AssignRule  map[string]interface{} `json:"assignRule"` // 分配规则
	FilterRule  map[string]interface{} `json:"filterRule"` // 下拉框过滤规则
}
