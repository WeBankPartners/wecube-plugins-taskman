package models

type FormTable struct {
	Id           string `json:"id" xorm:"'id' pk" primary-key:"id"`
	Name         string `json:"name" xorm:"name"`
	Description  string `json:"description" xorm:"description"`
	FormTemplate string `json:"formTemplate" xorm:"form_template"`
	CreatedBy    string `json:"createdBy" xorm:"created_by"`
	CreatedTime  string `json:"createdTime" xorm:"created_time"`
	UpdatedBy    string `json:"updatedBy" xorm:"updated_by"`
	UpdatedTime  string `json:"updatedTime" xorm:"updated_time"`
	DelFlag      int    `json:"delFlag" xorm:"del_flag"`
}

func (FormTable) TableName() string {
	return "form"
}

// FormNewTable 新表单
type FormNewTable struct {
	Id           string `json:"id" xorm:"'id' pk" primary-key:"id"`
	Request      string `json:"request" xorm:"request"`
	Task         string `json:"task" xorm:"task"`
	FormTemplate string `json:"formTemplate" xorm:"form_template"`
	DataId       string `json:"dataId" xorm:"data_id"`
	CreatedBy    string `json:"createdBy" xorm:"created_by"`
	CreatedTime  string `json:"createdTime" xorm:"created_time"`
	UpdatedBy    string `json:"updatedBy" xorm:"updated_by"`
	UpdatedTime  string `json:"updatedTime" xorm:"updated_time"`
}

func (FormNewTable) TableName() string {
	return "form_new"
}
