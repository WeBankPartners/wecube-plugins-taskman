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

type FormTable struct {
	Id           string `json:"id" xorm:"id"`
	Name         string `json:"name" xorm:"name"`
	Description  string `json:"description" xorm:"description"`
	FormTemplate string `json:"formTemplate" xorm:"form_template"`
	CreatedBy    string `json:"createdBy" xorm:"created_by"`
	CreatedTime  string `json:"createdTime" xorm:"created_time"`
	UpdatedBy    string `json:"updatedBy" xorm:"updated_by"`
	UpdatedTime  string `json:"updatedTime" xorm:"updated_time"`
	DelFlag      int    `json:"delFlag" xorm:"del_flag"`
}

type FormItemTemplateTable struct {
	Id              string `json:"id" xorm:"id"`
	Name            string `json:"name" xorm:"name"`
	Description     string `json:"description" xorm:"description"`
	ItemGroup       string `json:"itemGroup" xorm:"item_group"`
	ItemGroupName   string `json:"itemGroupName" xorm:"item_group_name"`
	FormTemplate    string `json:"formTemplate" xorm:"form_template"`
	DefaultValue    string `json:"defaultValue" xorm:"default_value"`
	Sort            int    `json:"sort" xorm:"sort"`
	PackageName     string `json:"packageName" xorm:"package_name"`
	Entity          string `json:"entity" xorm:"entity"`
	AttrDefId       string `json:"attrDefId" xorm:"attr_def_id"`
	AttrDefName     string `json:"attrDefName" xorm:"attr_def_name"`
	AttrDefDataType string `json:"attrDefDataType" xorm:"attr_def_data_type"`
	ElementType     string `json:"elementType" xorm:"element_type"`
	Title           string `json:"title" xorm:"title"`
	Width           int    `json:"width" xorm:"width"`
	RefPackageName  string `json:"refPackageName" xorm:"ref_package_name"`
	RefEntity       string `json:"refEntity" xorm:"ref_entity"`
	DataOptions     string `json:"dataOptions" xorm:"data_options"`
	Required        string `json:"required" xorm:"required"`
	Regular         string `json:"regular" xorm:"regular"`
	IsEdit          string `json:"isEdit" xorm:"is_edit"`
	IsView          string `json:"isView" xorm:"is_view"`
	IsOutput        string `json:"isOutput" xorm:"is_output"`
}

type FormItemTable struct {
	Id               string `json:"id" xorm:"id"`
	Form             string `json:"form" xorm:"form"`
	FormItemTemplate string `json:"formItemTemplate" xorm:"form_item_template"`
	Name             string `json:"name" xorm:"name"`
	Value            string `json:"value" xorm:"value"`
	ItemGroup        string `json:"itemGroup" xorm:"item_group"`
}

type FormTemplateDto struct {
	Id          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	UpdatedBy   string                   `json:"updatedBy"`
	UpdatedTime string                   `json:"updatedTime"`
	NowTime     string                   `json:"-"`
	Items       []*FormItemTemplateTable `json:"items"`
}
