package models

// RequestApproveRoleTable 请求审批角色表
type RequestApproveRoleTable struct {
	Id             string `json:"id" xorm:"'id' pk"`
	RequestApprove string `json:"requestApprove" xorm:"request_approve"` // 请求审批ID
	Role           string `json:"role" xorm:"role"`                      // 角色
	handler        string `json:"handler" xorm:"handler"`                // 处理人
	CreatedTime    string `json:"createdTime" xorm:"created_time"`
}

func (RequestApproveRoleTable) TableName() string {
	return "request_approve_role"
}
