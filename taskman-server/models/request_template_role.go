package models

type RequestTemplateRoleTable struct {
	Id              string `json:"id" xorm:"id"`
	RequestTemplate string `json:"requestTemplate" xorm:"request_template"`
	Role            string `json:"role" xorm:"role"`
	RoleType        string `json:"roleType" xorm:"role_type"`
}

func (RequestTemplateRoleTable) TableName() string {
	return "request_template_role"
}

func CreateRequestTemplateRoleTable(requestTemplateId, role string, permission RolePermission) *RequestTemplateRoleTable {
	return &RequestTemplateRoleTable{
		Id:              requestTemplateId + SysTableIdConnector + role + SysTableIdConnector + string(permission),
		RequestTemplate: requestTemplateId,
		Role:            role,
		RoleType:        string(permission),
	}
}
