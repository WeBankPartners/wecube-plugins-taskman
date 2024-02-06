package models

// RequestApproveTable 请求审批表
type RequestApproveTable struct {
	Id           string `json:"id" xorm:"id"`
	Request      string `json:"request" xorm:"request"`            // 请求ID
	Name         string `json:"name" xorm:"name"`                  // 审批名称
	ExpireDay    string `json:"expireDay" xorm:"expire_day"`       // 审批时效
	Type         string `json:"type" xorm:"type"`                  // 审批类型:审批类型: administrator 管理员,pass 自动通过,custom 自定义
	Sort         int    `json:"sort" xorm:"sort"`                  // 执行顺序
	FormTemplate string `json:"formTemplate" xorm:"form_template"` // 表单ID
	Method       string `json:"method" xorm:"method"`              // 分配类型:单人:single, 协同: anyone 并行:all
	CreatedBy    string `json:"createdBy" xorm:"created_by"`
	CreatedTime  string `json:"createdTime" xorm:"created_time"`
	UpdatedBy    string `json:"updatedBy" xorm:"updated_by"`
	UpdatedTime  string `json:"updatedTime" xorm:"updated_time"`
}

func (RequestApproveTable) TableName() string {
	return "request_approve"
}
