package models

// TaskApproveRoleTable 任务审批表
type TaskApproveRoleTable struct {
	Id          string `json:"id" xorm:"'id' pk" primary-key:"id"`
	Task        string `json:"task" xorm:"task"`       // 任务审批ID
	Role        string `json:"role" xorm:"role"`       // 角色
	handler     string `json:"handler" xorm:"handler"` // 处理人
	CreatedTime string `json:"createdTime" xorm:"created_time"`
}

func (TaskApproveRoleTable) TableName() string {
	return "task_approve_role"
}
