package models

type FormTemplateLibraryTable struct {
	Id          string `json:"id" xorm:"'id' pk" primary-key:"id"`
	Name        string `json:"name" xorm:"name"`
	FormType    string `json:"formType" xorm:"form_type"`
	CreatedTime string `json:"createdTime" xorm:"created_time"`
	UpdatedTime string `json:"updatedTime" xorm:"updated_time"`
	CreatedBy   string `json:"createdBy" xorm:"created_by"` // 创建人
	DelFlag     int    `json:"delFlag" xorm:"del_flag"`
}

func (FormTemplateLibraryTable) TableName() string {
	return "form_template_library"
}

type FormTemplateLibraryDto struct {
	Id          string                        `json:"id"`
	Name        string                        `json:"name"`
	FormType    string                        `json:"formType"`
	CreatedTime string                        `json:"createdTime"`
	CreatedBy   string                        `json:"createdBy"` // 创建人
	FormItems   string                        `json:"formItems"` //表单项,逗号隔开
	Items       []*FormItemTemplateLibraryDto `json:"items"`
}

type FormTemplateLibraryTableData struct {
	Id          string                              `json:"id" xorm:"id"`
	Name        string                              `json:"name" xorm:"name"`
	FormType    string                              `json:"formType" xorm:"form_type"`
	CreatedTime string                              `json:"createdTime" xorm:"created_time"`
	UpdatedTime string                              `json:"updatedTime" xorm:"updated_time"`
	CreatedBy   string                              `json:"createdBy" xorm:"created_by"` // 创建人
	DelFlag     int                                 `json:"delFlag" xorm:"del_flag"`
	Items       []*FormItemTemplateLibraryTableData `json:"items" xorm:"-"`
}
