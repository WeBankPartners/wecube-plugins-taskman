package models

import "github.com/WeBankPartners/go-common-lib/guid"

type FormTemplateTable struct {
	Id              string `json:"id" xorm:"'id' pk" primary-key:"id"`
	RequestTemplate string `json:"requestTemplate" xorm:"request_template"`
	TaskTemplate    string `json:"taskTemplate" xorm:"task_template"`
	ItemGroup       string `json:"itemGroup" xorm:"item_group"`
	ItemGroupType   string `json:"itemGroupType" xorm:"item_group_type"` //表单组类型:workflow 编排数据,optional 自选,custom 自定义,request_form 请求表单,db判断用
	ItemGroupName   string `json:"itemGroupName" xorm:"item_group_name"`
	ItemGroupSort   int    `json:"ItemGroupSort" xorm:"item_group_sort"`     // item_group 排序
	ItemGroupRule   string `json:"itemGroupRule" xorm:"item_group_rule"`     // item_group_rule 新增一行规则,new 输入新数据,exist 选择已有数据
	RefId           string `json:"refId" xorm:"ref_id"`                      // 引用ID
	RequestFormType string `json:"requestFormType" xorm:"request_form_type"` // 请求表单类型: message 信息表单,data 数据表单
	CreatedTime     string `json:"createdTime" xorm:"created_time"`
	DelFlag         int    `json:"delFlag" xorm:"del_flag"`
}

func (FormTemplateTable) TableName() string {
	return "form_template"
}

type FormTemplateTableSort []*FormTemplateTable

func (s FormTemplateTableSort) Len() int {
	return len(s)
}

func (s FormTemplateTableSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s FormTemplateTableSort) Less(i, j int) bool {
	if s[i].ItemGroupSort < s[j].ItemGroupSort {
		return true
	}
	return false
}

type FormTemplateDto struct {
	Id              string                 `json:"id"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	ExpireDay       int                    `json:"expireDay"`
	UpdatedBy       string                 `json:"updatedBy"`
	UpdatedTime     string                 `json:"updatedTime"`
	NowTime         string                 `json:"-"`
	RequestTemplate string                 `json:"-"`
	RequestFormType string                 `json:"-"`
	Items           []*FormItemTemplateDto `json:"items"`
}

// DataFormTemplateDto 全局表单模板 dto
type DataFormTemplateDto struct {
	AssociationWorkflow bool                    `json:"associationWorkflow"` // 是否关联编排
	UpdatedBy           string                  `json:"updatedBy"`           // 更新人
	Groups              []*FormTemplateGroupDto `json:"groups"`
}

type SimpleFormTemplateDto struct {
	TaskTemplateId string                  `json:"taskTemplateId"` // 数据表单模板ID
	UpdatedBy      string                  `json:"updatedBy"`      // 更新人
	Groups         []*FormTemplateGroupDto `json:"groups"`
}

// FormTemplateGroupDto 表单模板组dto
type FormTemplateGroupDto struct {
	ItemGroupId   string                 `json:"itemGroupId"` //表单组ID
	ItemGroup     string                 `json:"itemGroup"`
	ItemGroupType string                 `json:"itemGroupType"` // 表单组类型:workflow 编排数据,optional 自选,custom 自定义
	ItemGroupName string                 `json:"itemGroupName"`
	ItemGroupSort int                    `json:"itemGroupSort"` // 组排序
	Items         []*FormItemTemplateDto `json:"items"`         // 表单项
}

// FormTemplateGroupConfigureDto 表单组配置在dto
type FormTemplateGroupConfigureDto struct {
	RequestTemplateId string                    `json:"requestTemplateId"` // 模板Id
	TaskTemplateId    string                    `json:"taskTemplateId"`    // 任务模板ID
	FormTemplateId    string                    `json:"itemGroupId"`
	ItemGroup         string                    `json:"itemGroup"`
	ItemGroupType     string                    `json:"itemGroupType"` // 表单组类型:workflow 编排数据,optional 自选,custom 自定义
	ItemGroupName     string                    `json:"itemGroupName"`
	ItemGroupRule     string                    `json:"itemGroupRule"` // item_group_rule 新增一行规则,new 输入新数据,exist 选择已有数据
	ItemGroupSort     int                       `json:"itemGroupSort"` // 表单组排序
	SystemItems       []*ProcEntityAttributeObj `json:"systemItems"`   // 系统表单项
	CustomItems       []*FormItemTemplateDto    `json:"customItems"`   // 自定义表单项
}

// FormTemplateGroupCustomDataDto 表单组自定义数据dto
type FormTemplateGroupCustomDataDto struct {
	RequestTemplateId string                 `json:"requestTemplateId"` // 模板Id
	FormTemplateId    string                 `json:"itemGroupId"`
	Items             []*FormItemTemplateDto `json:"items"` // 表单项
}

// FormTemplateGroupSortDto 表单组排序dto
type FormTemplateGroupSortDto struct {
	RequestTemplateId string   `json:"requestTemplateId"` // 模板Id
	TaskTemplateId    string   `json:"taskTemplateId"`    // 任务模板Id
	ItemGroupIdSort   []string `json:"itemGroupIdSort"`   // 排序
}

type TaskFormItemQueryObj struct {
	Id               string `json:"id" xorm:"'id' pk" primary-key:"id"`
	Form             string `json:"form" xorm:"form"`
	FormItemTemplate string `json:"formItemTemplate" xorm:"form_item_template"`
	Name             string `json:"name" xorm:"name"`
	Value            string `json:"value" xorm:"value"`
	ItemGroup        string `json:"itemGroup" xorm:"item_group"`
	RowDataId        string `json:"rowDataId" xorm:"row_data_id"`
	AttrDefDataType  string `json:"attrDefDataType" xorm:"attr_def_data_type"`
	ElementType      string `json:"elementType" xorm:"element_type"`
}

func ConvertProcEntityAttributeObj2FormItemTemplate(param FormTemplateGroupConfigureDto, workflowEntityAttribute *ProcEntityAttributeObj, newItemGroupId string) *FormItemTemplateTable {
	var elementType = string(FormItemElementTypeInput)
	if workflowEntityAttribute.DataType == "ref" {
		elementType = string(FormItemElementTypeSelect)
	}
	return &FormItemTemplateTable{
		Id:              guid.CreateGuid(),
		Name:            workflowEntityAttribute.Name,
		Description:     workflowEntityAttribute.Description,
		FormTemplate:    newItemGroupId,
		ItemGroup:       param.ItemGroup,
		ItemGroupName:   param.ItemGroupName,
		Sort:            0,
		PackageName:     workflowEntityAttribute.EntityPackage,
		Entity:          workflowEntityAttribute.EntityName,
		AttrDefId:       workflowEntityAttribute.Id,
		AttrDefName:     workflowEntityAttribute.Name,
		AttrDefDataType: workflowEntityAttribute.DataType,
		ElementType:     elementType,
		Title:           workflowEntityAttribute.Description,
		Width:           24,
		RefPackageName:  workflowEntityAttribute.EntityPackage,
		RefEntity:       workflowEntityAttribute.Name,
		DataOptions:     "",
		Required:        "no",
		Regular:         "",
		IsEdit:          "yes",
		IsView:          "yes",
		IsOutput:        "no",
		InDisplayName:   "no",
		IsRefInside:     "no",
		Multiple:        workflowEntityAttribute.Multiple,
		DefaultClear:    "no",
		RefId:           "",
		SelectList:      nil,
		Active:          false,
	}
}

type FormTemplateGroupDtoSort []*FormTemplateGroupDto

func (s FormTemplateGroupDtoSort) Len() int {
	return len(s)
}

func (s FormTemplateGroupDtoSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s FormTemplateGroupDtoSort) Less(i, j int) bool {
	if s[i].ItemGroupSort < s[j].ItemGroupSort {
		return true
	}
	return false
}
