package models

type TaskTemplateTable struct {
	Id              string `json:"id" xorm:"id"`
	Name            string `json:"name" xorm:"name"`
	Description     string `json:"description" xorm:"description"`
	FormTemplate    string `json:"formTemplate" xorm:"form_template"`
	RequestTemplate string `json:"requestTemplate" xorm:"request_template"`
	NodeDefId       string `json:"nodeDefId" xorm:"node_def_id"`
	NodeName        string `json:"nodeName" xorm:"node_name"`
	CreatedBy       string `json:"createdBy" xorm:"created_by"`
	CreatedTime     string `json:"createdTime" xorm:"created_time"`
	UpdatedBy       string `json:"updatedBy" xorm:"updated_by"`
	UpdatedTime     string `json:"updatedTime" xorm:"updated_time"`
	DelFlag         int    `json:"delFlag" xorm:"del_flag"`
}

type TaskTable struct {
	Id                string `json:"id" xorm:"id"`
	Name              string `json:"name" xorm:"name"`
	Description       string `json:"description" xorm:"description"`
	Form              string `json:"form" xorm:"form"`
	AttachFile        string `json:"attachFile" xorm:"attach_file"`
	Status            string `json:"status" xorm:"status"`
	Version           string `json:"version" xorm:"version"`
	Request           string `json:"request" xorm:"request"`
	Parent            string `json:"parent" xorm:"parent"`
	TaskTemplate      string `json:"taskTemplate" xorm:"task_template"`
	PackageName       string `json:"packageName" xorm:"package_name"`
	EntityName        string `json:"entityName" xorm:"entity_name"`
	ProcDefId         string `json:"procDefId" xorm:"proc_def_id"`
	ProcDefKey        string `json:"procDefKey" xorm:"proc_def_key"`
	ProcDefName       string `json:"procDefName" xorm:"proc_def_name"`
	NodeDefId         string `json:"nodeDefId" xorm:"node_def_id"`
	NodeName          string `json:"nodeName" xorm:"node_name"`
	CallbackUrl       string `json:"callbackUrl" xorm:"callback_url"`
	CallbackParameter string `json:"callbackParameter" xorm:"callback_parameter"`
	Emergency         string `json:"emergency" xorm:"emergency"`
	Result            string `json:"result" xorm:"result"`
	Reporter          string `json:"reporter" xorm:"reporter"`
	ReportTime        string `json:"reportTime" xorm:"report_time"`
	ReportRole        string `json:"reportRole" xorm:"report_role"`
	UpdatedBy         string `json:"updatedBy" xorm:"updated_by"`
	UpdatedTime       string `json:"updatedTime" xorm:"updated_time"`
	DelFlag           string `json:"delFlag" xorm:"del_flag"`
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
	UpdatedTime  string                   `json:"updatedTime"`
	UpdatedBy    string                   `json:"updatedBy"`
	MGMTRoles    []string                 `json:"mgmtRoles"`
	USERoles     []string                 `json:"useRoles"`
	MGMTRoleObjs []*RoleTable             `json:"mgmtRoleObjs"`
	USERoleObjs  []*RoleTable             `json:"useRoleObjs"`
	Items        []*FormItemTemplateTable `json:"items"`
}
