package models

type TaskTemplateTable struct {
	Id              string `json:"id" xorm:"'id' pk" primary-key:"id"`
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
	HandleMode      string `json:"handleMode" xorm:"handle_mode"`
	Type            string `json:"type" xorm:"type"`
}

func (TaskTemplateTable) TableName() string {
	return "task_template"
}

// type TaskTemplateVo struct {
// 	Id      string `json:"id" xorm:"'id' pk"`
// 	Handler string `json:"handler" xorm:"handler"`
// 	Role    string `json:"role" xorm:"role"`
// }

type TaskTemplateDto struct {
	Id              string                   `json:"id"`
	Type            string                   `json:"type"`
	NodeId          string                   `json:"nodeId"`
	NodeDefId       string                   `json:"nodeDefId"`
	NodeDefName     string                   `json:"nodeDefName"`
	Name            string                   `json:"name"`
	Description     string                   `json:"description"`
	ExpireDay       int                      `json:"expireDay"`
	UpdatedTime     string                   `json:"updatedTime"`
	UpdatedBy       string                   `json:"updatedBy"`
	MGMTRoles       []string                 `json:"mgmtRoles"`    // 还有用吗？
	USERoles        []string                 `json:"useRoles"`     // 还有用吗？
	MGMTRoleObjs    []*RoleTable             `json:"mgmtRoleObjs"` // 还有用吗？
	USERoleObjs     []*RoleTable             `json:"useRoleObjs"`  // 还有用吗？
	Items           []*FormItemTemplateDto   `json:"items"`
	RequestTemplate string                   `json:"requestTemplate"`
	Sort            int                      `json:"sort"`
	HandleMode      string                   `json:"handleMode"`
	HandleTemplates []*TaskHandleTemplateDto `json:"handleTemplates"`
}

type TaskTemplateIdObj struct {
	Id        string `json:"id"`
	Sort      int    `json:"sort"`
	Name      string `json:"name"`
	NodeDefId string `json:"nodeDefId"`
}

type TaskTemplateCreateResponse struct {
	TaskTemplate *TaskTemplateDto     `json:"taskTemplate"`
	Ids          []*TaskTemplateIdObj `json:"ids"`
}

type TaskTemplateDeleteResponse struct {
	Type string               `json:"type"`
	Ids  []*TaskTemplateIdObj `json:"ids"`
}

type TaskTemplateListIdsResponse struct {
	Type string               `json:"type"`
	Ids  []*TaskTemplateIdObj `json:"ids"`
}
