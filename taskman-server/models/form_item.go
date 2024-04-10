package models

type FormItemTable struct {
	Id               string `json:"id" xorm:"'id' pk" primary-key:"id"`
	Form             string `json:"form" xorm:"form"`
	FormItemTemplate string `json:"formItemTemplate" xorm:"form_item_template"`
	Name             string `json:"name" xorm:"name"`
	Value            string `json:"value" xorm:"value"`
	ItemGroup        string `json:"itemGroup" xorm:"item_group"`
	RowDataId        string `json:"rowDataId" xorm:"row_data_id"`
	Request          string `json:"request" xorm:"request"`
	UpdatedTime      string `json:"updatedTime" xorm:"updated_time"`
	OriginalId       string `json:"originalId" xorm:"original_id"`
	TaskHandle       string `json:"taskHandle" xorm:"task_handle"`
	DelFlag          bool   `json:"delFlag" xorm:"del_flag"`
}

func (FormItemTable) TableName() string {
	return "form_item"
}
