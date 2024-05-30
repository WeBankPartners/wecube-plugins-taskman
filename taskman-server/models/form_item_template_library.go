package models

type FormItemTemplateLibraryTable struct {
	Id                  string `json:"id" xorm:"'id' pk" primary-key:"id"`
	Name                string `json:"name" xorm:"name"`
	Description         string `json:"description" xorm:"description"`
	ItemGroup           string `json:"itemGroup" xorm:"item_group"`
	ItemGroupName       string `json:"itemGroupName" xorm:"item_group_name"`
	DefaultValue        string `json:"defaultValue" xorm:"default_value"`
	Sort                int    `json:"sort" xorm:"sort"`
	PackageName         string `json:"packageName" xorm:"package_name"`
	Entity              string `json:"entity" xorm:"entity"`
	AttrDefId           string `json:"attrDefId" xorm:"attr_def_id"`
	AttrDefName         string `json:"attrDefName" xorm:"attr_def_name"`
	AttrDefDataType     string `json:"attrDefDataType" xorm:"attr_def_data_type"`
	ElementType         string `json:"elementType" xorm:"element_type"`
	Title               string `json:"title" xorm:"title"`
	Width               int    `json:"width" xorm:"width"`
	RefPackageName      string `json:"refPackageName" xorm:"ref_package_name"`
	RefEntity           string `json:"refEntity" xorm:"ref_entity"`
	DataOptions         string `json:"dataOptions" xorm:"data_options"`
	Required            string `json:"required" xorm:"required"`
	Regular             string `json:"regular" xorm:"regular"`
	IsEdit              string `json:"isEdit" xorm:"is_edit"`
	IsView              string `json:"isView" xorm:"is_view"`
	IsOutput            string `json:"isOutput" xorm:"is_output"`
	InDisplayName       string `json:"inDisplayName" xorm:"in_display_name"`
	IsRefInside         string `json:"isRefInside" xorm:"is_ref_inside"`
	Multiple            string `json:"multiple" xorm:"multiple"`
	DefaultClear        string `json:"defaultClear" xorm:"default_clear"`
	RefId               string `json:"refId" xorm:"ref_id"`                              // 复制数据表单ID,数据表单删除该表单项时,需要删除审批表单,任务表单对应数据项
	RoutineExpression   string `json:"routineExpression" xorm:"routine_expression"`      // 计算表达式
	ControlSwitch       string `json:"controlSwitch" xorm:"control_switch"`              // 控制审批/任务开关
	FormTemplateLibrary string `json:"formTemplateLibrary" xorm:"form_template_library"` // 表单组件库id
	HiddenCondition     string `json:"hiddenCondition" xorm:"hidden_condition"`          // 隐藏条件
}

func (FormItemTemplateLibraryTable) TableName() string {
	return "form_item_template_library"
}
