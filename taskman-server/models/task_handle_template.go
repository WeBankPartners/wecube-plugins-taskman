package models

// TaskHandleTemplateTable 任务模板处理表
type TaskHandleTemplateTable struct {
	Id           string `json:"id" xorm:"'id' pk" primary-key:"id"`
	TaskTemplate string `json:"taskTemplate" xorm:"task_template"`
	Role         string `json:"role" xorm:"role"`
	Assign       string `json:"assign" xorm:"assign"`
	HandlerType  string `json:"handlerType" xorm:"handler_type"`
	Handler      string `json:"handler" xorm:"handler"`
}

func (TaskHandleTemplateTable) TableName() string {
	return "task_handle_template"
}
