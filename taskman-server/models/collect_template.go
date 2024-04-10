package models

// CollectTemplateTable  收藏模板
type CollectTemplateTable struct {
	Id              string `json:"id" xorm:"'id' pk" primary-key:"id"`
	RequestTemplate string `json:"requestTemplate" xorm:"request_template"` // 收藏模板ID
	User            string `json:"user" xorm:"user"`                        // 收藏用户
	Role            string `json:"role" xorm:"role"`                        // 收藏模板时角色
	Type            int    `json:"type" xorm:"type"`                        // 类型:0表示请求 1发布
	CreatedTime     string `json:"createdTime" xorm:"created_time"`
}

func (CollectTemplateTable) TableName() string {
	return "collect_template"
}

type CollectFilterItem struct {
	TemplateGroupList   []*KeyValuePair `json:"templateList"`        // 模板组列表
	OperatorObjTypeList []string        `json:"operatorObjTypeList"` // 操作对象类型列表
	ProcDefNameList     []string        `json:"procDefNameList"`     // 使用编排
	OwnerList           []string        `json:"ownerList"`           // 属主
	TagList             []string        `json:"tagList"`             // 标签
	ManageRoleList      []string        `json:"manageRoleList"`      // 属主角色
	UseRoleList         []string        `json:"useRoleList"`         // 使用角色
}
