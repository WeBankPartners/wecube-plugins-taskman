package models

type RoleTable struct {
	Id          string `json:"id" xorm:"'id' pk" primary-key:"id"`
	DisplayName string `json:"displayName" xorm:"display_name"`
	UpdatedTime string `json:"updatedTime" xorm:"updated_time"`
	CoreId      string `json:"coreId" xorm:"core_id"`
	Email       string `json:"email"`
}
