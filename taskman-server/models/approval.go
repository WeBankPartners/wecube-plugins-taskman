package models

type ApprovalTable struct {
	Id               string `json:"id" xorm:"'id' pk" primary-key:"id"`
	Sort             int    `json:"sort" xorm:"sort"`
	ApprovalTemplate string `json:"approvalTemplate" xorm:"approval_template"`
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
	HandlerType          string `json:"handlerType" xorm:"handler_type"`
}

type ApprovalDto struct {
	Id               string             `json:"id"`
	Sort             int                `json:"sort"`
	ApprovalTemplate string             `json:"approvalTemplate"`
	Request          string             `json:"request" xorm:"request"`
	RoleObjs         []*ApprovalRoleDto `json:"roleObjs"`
}

type ApprovalRoleDto struct {
	ApprovalTemplateRole string `json:"approvalTemplateRole"`
	Role                 string `json:"role"`
	Handler              string `json:"handler"`
	Approve              string `json:"approve"`
}
