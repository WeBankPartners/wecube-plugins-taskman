package models

type FormItemTemplateTable struct {
	Id              string           `json:"id" xorm:"'id' pk"`
	Name            string           `json:"name" xorm:"name"`
	Description     string           `json:"description" xorm:"description"`
	ItemGroup       string           `json:"itemGroup" xorm:"item_group"`
	ItemGroupType   string           `json:"itemGroupType" xorm:"item_group_type"` //表单组类型:workflow 编排数据,optional 自选,custom 自定义
	ItemGroupName   string           `json:"itemGroupName" xorm:"item_group_name"`
	FormTemplate    string           `json:"formTemplate" xorm:"form_template"`
	DefaultValue    string           `json:"defaultValue" xorm:"default_value"`
	Sort            int              `json:"sort" xorm:"sort"`
	PackageName     string           `json:"packageName" xorm:"package_name"`
	Entity          string           `json:"entity" xorm:"entity"`
	AttrDefId       string           `json:"attrDefId" xorm:"attr_def_id"`
	AttrDefName     string           `json:"attrDefName" xorm:"attr_def_name"`
	AttrDefDataType string           `json:"attrDefDataType" xorm:"attr_def_data_type"`
	ElementType     string           `json:"elementType" xorm:"element_type"`
	Title           string           `json:"title" xorm:"title"`
	Width           int              `json:"width" xorm:"width"`
	RefPackageName  string           `json:"refPackageName" xorm:"ref_package_name"`
	RefEntity       string           `json:"refEntity" xorm:"ref_entity"`
	DataOptions     string           `json:"dataOptions" xorm:"data_options"`
	Required        string           `json:"required" xorm:"required"`
	Regular         string           `json:"regular" xorm:"regular"`
	IsEdit          string           `json:"isEdit" xorm:"is_edit"`
	IsView          string           `json:"isView" xorm:"is_view"`
	IsOutput        string           `json:"isOutput" xorm:"is_output"`
	InDisplayName   string           `json:"inDisplayName" xorm:"in_display_name"`
	IsRefInside     string           `json:"isRefInside" xorm:"is_ref_inside"`
	Multiple        string           `json:"multiple" xorm:"multiple"`
	DefaultClear    string           `json:"defaultClear" xorm:"default_clear"`
	SelectList      []*EntityDataObj `json:"selectList" xorm:"-"`
}

func (FormItemTemplateTable) TableName() string {
	return "form_item_template"
}
