package models

type FormItemTemplateTable struct {
	Id              string           `json:"id" xorm:"'id' pk" primary-key:"id"`
	Name            string           `json:"name" xorm:"name"`
	Description     string           `json:"description" xorm:"description"`
	ItemGroup       string           `json:"itemGroup" xorm:"item_group"`
	ItemGroupType   string           `json:"itemGroupType" xorm:"item_group_type"` //表单组类型:workflow 编排数据,optional 自选,custom 自定义
	ItemGroupName   string           `json:"itemGroupName" xorm:"item_group_name"`
	ItemGroupSort   int              `json:"ItemGroupSort" xorm:"item_group_sort"` // item_group 排序
	ItemGroupRule   string           `json:"itemGroupRule" xorm:"item_group_rule"` // item_group_rule 新增一行规则,new 输入新数据,exist 选择已有数据
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
	CopyId          string           `json:"copyId" xorm:"copy_id"` //复制数据表单ID,数据表单删除该表单项时,需要删除审批表单,任务表单对应数据项
	SelectList      []*EntityDataObj `json:"selectList" xorm:"-"`
	Active          bool             `json:"active" xorm:"-"` // 是否选中状态
}

type FormItemTemplateDto struct {
	Id              string           `json:"id"`
	Name            string           `json:"name"`
	Description     string           `json:"description"`
	ItemGroup       string           `json:"itemGroup"`
	ItemGroupType   string           `json:"itemGroupType"` //表单组类型:workflow 编排数据,optional 自选,custom 自定义
	ItemGroupName   string           `json:"itemGroupName"`
	ItemGroupSort   int              `json:"ItemGroupSort"` // item_group 排序
	ItemGroupRule   string           `json:"itemGroupRule"` // item_group_rule 新增一行规则,new 输入新数据,exist 选择已有数据
	FormTemplate    string           `json:"formTemplate"`
	DefaultValue    string           `json:"defaultValue"`
	Sort            int              `json:"sort"`
	PackageName     string           `json:"packageName"`
	Entity          string           `json:"entity"`
	AttrDefId       string           `json:"attrDefId"`
	AttrDefName     string           `json:"attrDefName"`
	AttrDefDataType string           `json:"attrDefDataType"`
	ElementType     string           `json:"elementType"`
	Title           string           `json:"title"`
	Width           int              `json:"width"`
	RefPackageName  string           `json:"refPackageName"`
	RefEntity       string           `json:"refEntity"`
	DataOptions     string           `json:"dataOptions"`
	Required        string           `json:"required"`
	Regular         string           `json:"regular"`
	IsEdit          string           `json:"isEdit"`
	IsView          string           `json:"isView"`
	IsOutput        string           `json:"isOutput"`
	InDisplayName   string           `json:"inDisplayName"`
	IsRefInside     string           `json:"isRefInside"`
	Multiple        string           `json:"multiple"`
	DefaultClear    string           `json:"defaultClear"`
	CopyId          string           `json:"copyId"` //复制数据表单ID,数据表单删除该表单项时,需要删除审批表单,任务表单对应数据项
	SelectList      []*EntityDataObj `json:"selectList"`
	Active          bool             `json:"active"` // 是否选中状态
}

func (FormItemTemplateTable) TableName() string {
	return "form_item_template"
}

type FormItemTemplateTableSort []*FormItemTemplateTable

func (s FormItemTemplateTableSort) Len() int {
	return len(s)
}

func (s FormItemTemplateTableSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s FormItemTemplateTableSort) Less(i, j int) bool {
	if s[i].Sort < s[j].Sort {
		return true
	}
	return false
}

func ConvertFormItemTemplateList2Map(formItemTemplateList []*FormItemTemplateTable) map[string]*FormItemTemplateTable {
	hashMap := make(map[string]*FormItemTemplateTable)
	for _, table := range formItemTemplateList {
		hashMap[table.Id] = table
	}
	return hashMap
}
