package models

// TaskHandleTemplateTable 任务模板处理表
type TaskHandleTemplateTable struct {
	Id           string `json:"id" xorm:"'id' pk" primary-key:"id"`
	Sort         int    `json:"sort" xorm:"sort"`
	TaskTemplate string `json:"taskTemplate" xorm:"task_template"`
	Role         string `json:"role" xorm:"role"`
	Assign       string `json:"assign" xorm:"assign"`
	HandlerType  string `json:"handlerType" xorm:"handler_type"`
	Handler      string `json:"handler" xorm:"handler"`
	HandleMode   string `json:"handleMode" xorm:"handle_mode"` // 任务类型
}

func (TaskHandleTemplateTable) TableName() string {
	return "task_handle_template"
}

type TaskHandleTemplateDto struct {
	Role        string `json:"role"`
	Assign      string `json:"assign"`
	HandlerType string `json:"handlerType"`
	Handler     string `json:"handler"`
}
