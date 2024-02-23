package models

import "github.com/WeBankPartners/go-common-lib/guid"

type FormTemplateTable struct {
	Id          string `json:"id" xorm:"'id' pk" primary-key:"id"`
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
	Id          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	ExpireDay   int                    `json:"expireDay"`
	UpdatedBy   string                 `json:"updatedBy"`
	UpdatedTime string                 `json:"updatedTime"`
	NowTime     string                 `json:"-"`
	Items       []*FormItemTemplateDto `json:"items"`
}

// DataFormTemplateDto 全局表单模板 dto
type DataFormTemplateDto struct {
	FormTemplateId      string                  `json:"formTemplateId"`      // 数据表单模板ID
	AssociationWorkflow bool                    `json:"associationWorkflow"` // 是否关联编排
	UpdatedBy           string                  `json:"updatedBy"`           // 更新人
	Groups              []*FormTemplateGroupDto `json:"groups"`
}

type SimpleFormTemplateDto struct {
	FormTemplateId string                  `json:"formTemplateId"` // 数据表单模板ID
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
	FormTemplateId    string                    `json:"formTemplateId"`    // 表单模板ID
	ItemGroupId       string                    `json:"itemGroupId"`
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
	FormTemplateId    string                 `json:"formTemplateId"`    // 表单模板ID
	ItemGroupId       string                 `json:"itemGroupId"`
	Items             []*FormItemTemplateDto `json:"items"` // 表单项
}

// FormTemplateGroupSortDto 表单组排序dto
type FormTemplateGroupSortDto struct {
	RequestTemplateId string   `json:"requestTemplateId"` // 模板Id
	FormTemplateId    string   `json:"formTemplateId"`    // 表单模板ID
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

func CovertFormTemplateDto2Model(dto FormTemplateDto) *FormTemplateTable {
	return &FormTemplateTable{
		Id:          dto.Id,
		Name:        dto.Name,
		Description: dto.Description,
		CreatedBy:   dto.UpdatedBy,
		CreatedTime: dto.NowTime,
		UpdatedBy:   dto.UpdatedBy,
		UpdatedTime: dto.NowTime,
	}
}

func ConvertDataFormTemplate2FormTemplateDto(dto DataFormTemplateDto) FormTemplateDto {
	var items = make([]*FormItemTemplateDto, 0)
	if len(dto.Groups) > 0 {
		for _, group := range dto.Groups {
			if group != nil && len(group.Items) > 0 {
				for _, item := range group.Items {
					if item != nil {
						items = append(items, item)
					}
				}
			}
		}
	}
	return FormTemplateDto{
		Id:        dto.FormTemplateId,
		Name:      "global_form_template",
		ExpireDay: 0,
		UpdatedBy: dto.UpdatedBy,
		Items:     items,
	}
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
		ItemGroupId:     newItemGroupId,
		ItemGroup:       param.ItemGroup,
		ItemGroupName:   param.ItemGroupName,
		FormTemplate:    param.FormTemplateId,
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
		CopyId:          "",
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
