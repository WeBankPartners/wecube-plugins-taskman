package models

type CollectFilterItem struct {
	TemplateGroupList   []*KeyValuePair `json:"templateList"`        // 模板组列表
	OperatorObjTypeList []string        `json:"operatorObjTypeList"` // 操作对象类型列表
	ProcDefNameList     []string        `json:"procDefNameList"`     // 使用编排
	OwnerList           []string        `json:"ownerList"`           // 属主
	TagList             []string        `json:"tagList"`             // 标签
	ManageRoleList      []string        `json:"manageRoleList"`      // 属主角色
	UseRoleList         []string        `json:"useRoleList"`         // 使用角色
}
