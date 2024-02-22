package models

type FormItemTemplateGroupTable struct {
	Id            string `json:"id" xorm:"'id' pk" primary-key:"id"`
	Name          string `json:"name" xorm:"name"`
	Description   string `json:"description" xorm:"description"`
	ItemGroup     string `json:"itemGroup" xorm:"item_group"`
	ItemGroupType string `json:"itemGroupType" xorm:"item_group_type"` //表单组类型:workflow 编排数据,optional 自选,custom 自定义
	ItemGroupName string `json:"itemGroupName" xorm:"item_group_name"`
	ItemGroupSort int    `json:"ItemGroupSort" xorm:"item_group_sort"` // item_group 排序
	ItemGroupRule string `json:"itemGroupRule" xorm:"item_group_rule"` // item_group_rule 新增一行规则,new 输入新数据,exist 选择已有数据
	FormTemplate  string `json:"formTemplate" xorm:"form_template"`
}
