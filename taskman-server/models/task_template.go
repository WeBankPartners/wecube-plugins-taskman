package models

type TaskTemplateTable struct {
	Id              string `json:"id" xorm:"id"`
	Name            string `json:"name" xorm:"name"`
	Description     string `json:"description" xorm:"description"`
	FormTemplate    string `json:"formTemplate" xorm:"form_template"`
	RequestTemplate string `json:"requestTemplate" xorm:"request_template"`
	NodeDefId       string `json:"nodeDefId" xorm:"node_def_id"`
	NodeName        string `json:"nodeName" xorm:"node_name"`
	ExpireDay       int    `json:"expireDay" xorm:"expire_day"`
	Handler         string `json:"handler" xorm:"handler"`
	CreatedBy       string `json:"createdBy" xorm:"created_by"`
	CreatedTime     string `json:"createdTime" xorm:"created_time"`
	UpdatedBy       string `json:"updatedBy" xorm:"updated_by"`
	UpdatedTime     string `json:"updatedTime" xorm:"updated_time"`
	DelFlag         int    `json:"delFlag" xorm:"del_flag"`
}

type TaskTemplateRoleTable struct {
	Id           string `json:"id" xorm:"id"`
	TaskTemplate string `json:"taskTemplate" xorm:"task_template"`
	Role         string `json:"role" xorm:"role"`
	RoleType     string `json:"roleType" xorm:"role_type"`
}

type TaskTemplateDto struct {
	Id           string                   `json:"id"`
	NodeDefId    string                   `json:"nodeDefId"`
	NodeDefName  string                   `json:"nodeDefName"`
	Name         string                   `json:"name"`
	Description  string                   `json:"description"`
	ExpireDay    int                      `json:"expireDay"`
	Handler      string                   `json:"handler"`
	UpdatedTime  string                   `json:"updatedTime"`
	UpdatedBy    string                   `json:"updatedBy"`
	MGMTRoles    []string                 `json:"mgmtRoles"`
	USERoles     []string                 `json:"useRoles"`
	MGMTRoleObjs []*RoleTable             `json:"mgmtRoleObjs"`
	USERoleObjs  []*RoleTable             `json:"useRoleObjs"`
	Items        []*FormItemTemplateTable `json:"items"`
}
