package models

type ApprovalTemplateTable struct {
	Id              string `json:"id" xorm:"'id' pk" primary-key:"id"`
	Sort            int    `json:"sort" xorm:"sort"`
	RequestTemplate string `json:"requestTemplate" xorm:"request_template"`
	FormTemplate    string `json:"formTemplate" xorm:"form_template"`
	Name            string `json:"name" xorm:"name"`
	ExpireDay       int    `json:"expireDay" xorm:"expire_day"`
	Description     string `json:"description" xorm:"description"`
	RoleType        string `json:"roleType" xorm:"role_type"`
	CreatedBy       string `json:"createdBy" xorm:"created_by"`
	CreatedTime     string `json:"createdTime" xorm:"created_time"`
	UpdatedBy       string `json:"updatedBy" xorm:"updated_by"`
	UpdatedTime     string `json:"updatedTime" xorm:"updated_time"`
}

type ApprovalTemplateRoleTable struct {
	Id               string `json:"id" xorm:"'id' pk" primary-key:"id"`
	Sort             int    `json:"sort" xorm:"sort"`
	ApprovalTemplate string `json:"approvalTemplate" xorm:"approval_template"`
	RoleType         string `json:"roleType" xorm:"role_type"`
	HandlerType      string `json:"handlerType" xorm:"handler_type"`
	Role             string `json:"role" xorm:"role"`
	Handler          string `json:"handler" xorm:"handler"`
}

type ApprovalTemplateDto struct {
	Id              string                     `json:"id"`
	Sort            int                        `json:"sort"`
	RequestTemplate string                     `json:"requestTemplate"`
	FormTemplate    string                     `json:"formTemplate"`
	Name            string                     `json:"name"`
	ExpireDay       int                        `json:"expireDay"`
	Description     string                     `json:"description"`
	RoleType        string                     `json:"roleType"`
	RoleObjs        []*ApprovalTemplateRoleDto `json:"roleObjs"`
}

type ApprovalTemplateRoleDto struct {
	ApprovalRoleId string `json:"ApprovalRoleId"`
	RoleType       string `json:"roleType"`
	HandlerType    string `json:"handlerType"`
	Role           string `json:"role"`
	Handler        string `json:"handler"`
}

type ApprovalTemplateIdObj struct {
	Id           string `json:"id"`
	Sort         int    `json:"sort"`
	Name         string `json:"name"`
	FormTemplate string `json:"formTemplate"`
}

type ApprovalTemplateCreateParam struct {
	Sort            int    `json:"sort"`
	RequestTemplate string `json:"requestTemplate"`
	Name            string `json:"name"`
}

type ApprovalTemplateCreateResponse struct {
	Id              string                   `json:"id"`
	Sort            int                      `json:"sort"`
	RequestTemplate string                   `json:"requestTemplate"`
	Name            string                   `json:"name"`
	Ids             []*ApprovalTemplateIdObj `json:"ids"`
}

type ApprovalTemplateDeleteResponse struct {
	Ids []*ApprovalTemplateIdObj `json:"ids"`
}

type ApprovalTemplateListIdsResponse struct {
	Ids []*ApprovalTemplateIdObj `json:"ids"`
}
