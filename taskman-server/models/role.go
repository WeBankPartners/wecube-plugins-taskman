package models

import "strings"

type RoleTable struct {
	Id          string `json:"id" xorm:"'id' pk" primary-key:"id"`
	DisplayName string `json:"displayName" xorm:"display_name"`
	UpdatedTime string `json:"updatedTime" xorm:"updated_time"`
	CoreId      string `json:"coreId" xorm:"core_id"`
	Email       string `json:"email"`
}

type RoleApply struct {
	PageInfo *PageInfo       `json:"pageInfo"` // 分页信息
	Contents []*RoleApplyDto `json:"contents"` // 列表内容
}

type RoleApplyDto struct {
	ID          string              `json:"id"`
	CreatedBy   string              `json:"createdBy"`
	UpdatedBy   string              `json:"updatedBy"`
	CreatedTime string              `json:"createdTime"`
	UpdatedTime string              `json:"updatedTime"`
	EmailAddr   string              `json:"emailAddr"`
	Role        *SimpleLocalRoleDto `json:"role"`
	Status      string              `json:"status"`     // init,approve,deny,expire,preExpried
	ExpireTime  string              `json:"expireTime"` //角色过期时间,""表示永久生效
}

type RoleTableSort []*RoleTable

func (s RoleTableSort) Len() int {
	return len(s)
}

func (s RoleTableSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s RoleTableSort) Less(i, j int) bool {
	return strings.Compare(s[i].CoreId, s[j].CoreId) < 0
}
