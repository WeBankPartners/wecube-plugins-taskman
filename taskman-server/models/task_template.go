package models

type TaskTemplateTable struct {
	Id              string `json:"id" xorm:"'id' pk" primary-key:"id"`
	Type            string `json:"type" xorm:"type"`
	Name            string `json:"name" xorm:"name"`
	Description     string `json:"description" xorm:"description"`
	FormTemplate    string `json:"formTemplate" xorm:"form_template"`
	RequestTemplate string `json:"requestTemplate" xorm:"request_template"`
	NodeId          string `json:"nodeId" xorm:"node_id"`
	NodeDefId       string `json:"nodeDefId" xorm:"node_def_id"`
	NodeName        string `json:"nodeName" xorm:"node_name"`
	ExpireDay       int    `json:"expireDay" xorm:"expire_day"`
	Handler         string `json:"handler" xorm:"handler"`
	CreatedBy       string `json:"createdBy" xorm:"created_by"`
	CreatedTime     string `json:"createdTime" xorm:"created_time"`
	UpdatedBy       string `json:"updatedBy" xorm:"updated_by"`
	UpdatedTime     string `json:"updatedTime" xorm:"updated_time"`
	DelFlag         int    `json:"delFlag" xorm:"del_flag"`
	Sort            int    `json:"sort" xorm:"sort"`
	RoleType        string `json:"roleType" xorm:"role_type"`
}

type TaskTemplateVo struct {
	Id      string `json:"id" xorm:"'id' pk"`
	Handler string `json:"handler" xorm:"handler"`
	Role    string `json:"role" xorm:"role"`
}

type TaskTemplateRoleTable struct {
	Id             string `json:"id" xorm:"'id' pk" primary-key:"id"`
	TaskTemplate   string `json:"taskTemplate" xorm:"task_template"`
	Role           string `json:"role" xorm:"role"`
	RoleType       string `json:"roleType" xorm:"role_type"`
	CustomRoleType string `json:"customRoleType" xorm:"custom_role_type"`
	HandlerType    string `json:"handlerType" xorm:"handler_type"`
	CustomRole     string `json:"customRole" xorm:"custom_role"`
	Handler        string `json:"handler" xorm:"handler"`
}

type TaskTemplateDto struct {
	Id                string                 `json:"id"`
	Type              string                 `json:"type"`
	NodeId            string                 `json:"nodeId"`
	NodeDefId         string                 `json:"nodeDefId"`
	NodeDefName       string                 `json:"nodeDefName"`
	Name              string                 `json:"name"`
	Description       string                 `json:"description"`
	ExpireDay         int                    `json:"expireDay"`
	Handler           string                 `json:"handler"`
	UpdatedTime       string                 `json:"updatedTime"`
	UpdatedBy         string                 `json:"updatedBy"`
	MGMTRoles         []string               `json:"mgmtRoles"`
	USERoles          []string               `json:"useRoles"`
	MGMTRoleObjs      []*RoleTable           `json:"mgmtRoleObjs"`
	USERoleObjs       []*RoleTable           `json:"useRoleObjs"`
	Items             []*FormItemTemplateDto `json:"items"`
	RequestTemplateId string                 `json:"requestTemplateId"`
	Sort              int                    `json:"sort"`
	RoleType          string                 `json:"roleType"`
	RoleObjs          []*TaskTemplateRoleDto `json:"roleObjs"`
}

type TaskTemplateRoleDto struct {
	RoleType    string `json:"roleType"`
	HandlerType string `json:"handlerType"`
	Role        string `json:"role"`
	Handler     string `json:"handler"`
}

type CustomTaskTemplateIdObj struct {
	Id   string `json:"id"`
	Sort int    `json:"sort"`
	Name string `json:"name"`
}

type CustomTaskTemplateCreateParam struct {
	Sort            int    `json:"sort"`
	RequestTemplate string `json:"requestTemplate"`
	Name            string `json:"name"`
}

type CustomTaskTemplateCreateResponse struct {
	Id              string                     `json:"id"`
	Sort            int                        `json:"sort"`
	RequestTemplate string                     `json:"requestTemplate"`
	Name            string                     `json:"name"`
	Ids             []*CustomTaskTemplateIdObj `json:"ids"`
}

type CustomTaskTemplateDeleteResponse struct {
	Ids []*CustomTaskTemplateIdObj `json:"ids"`
}

type CustomTaskTemplateListIdsResponse struct {
	Ids []*CustomTaskTemplateIdObj `json:"ids"`
}
