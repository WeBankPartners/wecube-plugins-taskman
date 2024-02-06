package models

type FormTemplateTable struct {
	Id          string `json:"id" xorm:"id"`
	Name        string `json:"name" xorm:"name"`
	Description string `json:"description" xorm:"description"`
	CreatedBy   string `json:"createdBy" xorm:"created_by"`
	CreatedTime string `json:"createdTime" xorm:"created_time"`
	UpdatedBy   string `json:"updatedBy" xorm:"updated_by"`
	UpdatedTime string `json:"updatedTime" xorm:"updated_time"`
	DelFlag     int    `json:"delFlag" xorm:"del_flag"`
}

func (FormTemplateTable) TableName() string {
	return "form_template"
}

type FormTemplateDto struct {
	Id          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	ExpireDay   int                      `json:"expireDay"`
	UpdatedBy   string                   `json:"updatedBy"`
	UpdatedTime string                   `json:"updatedTime"`
	NowTime     string                   `json:"-"`
	Items       []*FormItemTemplateTable `json:"items"`
}

type TaskFormItemQueryObj struct {
	Id               string `json:"id" xorm:"id"`
	Form             string `json:"form" xorm:"form"`
	FormItemTemplate string `json:"formItemTemplate" xorm:"form_item_template"`
	Name             string `json:"name" xorm:"name"`
	Value            string `json:"value" xorm:"value"`
	ItemGroup        string `json:"itemGroup" xorm:"item_group"`
	RowDataId        string `json:"rowDataId" xorm:"row_data_id"`
	AttrDefDataType  string `json:"attrDefDataType" xorm:"attr_def_data_type"`
	ElementType      string `json:"elementType" xorm:"element_type"`
}
