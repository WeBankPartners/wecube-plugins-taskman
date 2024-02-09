package models

type RequestTemplateGroupTable struct {
	Id            string    `json:"id" xorm:"'id' pk"`
	Name          string    `json:"name" xorm:"name" binding:"required"`
	Description   string    `json:"description" xorm:"description"`
	ManageRole    string    `json:"manageRole" xorm:"manage_role" binding:"required"`
	ManageRoleObj RoleTable `json:"manageRoleObj" xorm:"-"`
	CreatedBy     string    `json:"createdBy" xorm:"created_by"`
	CreatedTime   string    `json:"createdTime" xorm:"created_time"`
	UpdatedBy     string    `json:"updatedBy" xorm:"updated_by"`
	UpdatedTime   string    `json:"updatedTime" xorm:"updated_time"`
	DelFlag       int       `json:"delFlag" xorm:"del_flag"`
}

func (RequestTemplateGroupTable) TableName() string {
	return "request_template_group"
}

type TemplateGroupObj struct {
	GroupId     string                     `json:"groupId"`
	GroupName   string                     `json:"groupName"`
	CreatedTime string                     `json:"createdTime"`
	UpdatedTime string                     `json:"updatedTime"`
	Templates   []*RequestTemplateTableObj `json:"templates"`
}

type TemplateGroupSort []*TemplateGroupObj

func (s TemplateGroupSort) Len() int {
	return len(s)
}

func (s TemplateGroupSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s TemplateGroupSort) Less(i, j int) bool {
	return s[i].UpdatedTime > s[j].UpdatedTime
}
