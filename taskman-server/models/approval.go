package models

type ApprovalTable struct {
	Id               string `json:"id" xorm:"'id' pk" primary-key:"id"`
	Sort             int    `json:"sort" xorm:"sort"`
	ApprovalTemplate string `json:"approvalTemplate" xorm:"approval_template"`
	TemplateRoleType string `json:"templateRoleType" xorm:"template_role_type"`
	Request          string `json:"request" xorm:"request"`
	Name             string `json:"name" xorm:"name"`
	Approve          string `json:"approve" xorm:"approve"`
}

type ApprovalRoleTable struct {
	Id                   string `json:"id" xorm:"'id' pk" primary-key:"id"`
	Sort                 int    `json:"sort" xorm:"sort"`
	ApprovalTemplateRole string `json:"approvalTemplateRole" xorm:"approval_template_role"`
	Approval             string `json:"approval" xorm:"approval"`
	Role                 string `json:"role" xorm:"role"`
	Handler              string `json:"handler" xorm:"handler"`
	Approve              string `json:"approve" xorm:"approve"`
	HandlerType          string `json:"handlerType" xorm:"handler_type"` //人员设置方式：template.模板指定 template_suggest.模板建议 custom.提交人指定 custom_suggest.提交人建议 system.组内系统分配 claim.组内主动认领
}

type ApprovalDto struct {
	Id               string             `json:"id"`
	Sort             int                `json:"sort"`
	ApprovalTemplate string             `json:"approvalTemplate"`
	Request          string             `json:"request"`
	Name             string             `json:"name"`
	TemplateRoleType string             `json:"templateRoleType"`
	RoleObjs         []*ApprovalRoleDto `json:"roleObjs"`
}

type ApprovalRoleDto struct {
	ApprovalTemplateRole string `json:"approvalTemplateRole"`
	Role                 string `json:"role"`
	Handler              string `json:"handler"`
	Approve              string `json:"approve"`
}
