package models

// TaskHandleTable 任务处理表
type TaskHandleTable struct {
	Id                 string `json:"id" xorm:"'id' pk" primary-key:"id"`
	TaskHandleTemplate string `json:"taskHandleTemplate" xorm:"task_handle_template"`
	Task               string `json:"task" xorm:"task"`
	Role               string `json:"role" xorm:"role"`
	Handler            string `json:"handler" xorm:"handler"`
	HandlerType        string `json:"handlerType" xorm:"handler_type"`
	HandleResult       string `json:"handleResult" xorm:"handle_result"` //处理结果:approve同意,deny拒绝, redraw打回,complete完成,uncomplete未完成
	ResultDesc         string `json:"resultDesc" xorm:"result_desc"`     //处理描述
	ParentId           string `json:"parentId" xorm:"parent_id"`
	ChangeReason       string `json:"changeReason" xorm:"change_reason"` //变更原因: assign 系统分配、claim 主动领取、give 转给我
	CreatedTime        string `json:"createdTime" xorm:"created_time"`
	UpdatedTime        string `json:"updatedTime" xorm:"updated_time"`
	Sort               int    `json:"sort" xorm:"sort"`
	HandleStatus       string `json:"handleStatus" xorm:"handle_status"`    // 处理状态：complete 完成, uncomplete 未完成
	LatestFlag         int    `json:"latestFlag" xorm:"latest_flag"`        // 最新标记:1表示最新，0表示非最新
	ProcDefResult      string `json:"procDefResult" xorm:"proc_def_result"` // 编排选项结果
}

func (TaskHandleTable) TableName() string {
	return "task_handle"
}

type TaskHandleDto struct {
	Id                 string `json:"id"`
	Sort               int    `json:"sort"`
	TaskHandleTemplate string `json:"taskHandleTemplate"`
	Role               string `json:"role"`
	Handler            string `json:"handler"`
	HandleResult       string `json:"handleResult"`
	UpdatedTime        string `json:"updatedTime" xorm:"updated_time"`
}
