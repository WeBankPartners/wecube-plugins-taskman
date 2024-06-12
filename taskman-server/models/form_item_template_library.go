package models

import (
	"encoding/json"
	"strings"
)

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

type FormItemTemplateLibraryDto struct {
	Id                  string                   `json:"id"`
	Name                string                   `json:"name"`
	Description         string                   `json:"description"`
	ItemGroup           string                   `json:"itemGroup"`
	ItemGroupName       string                   `json:"itemGroupName"`
	DefaultValue        string                   `json:"defaultValue"`
	Sort                int                      `json:"sort"`
	PackageName         string                   `json:"packageName"`
	Entity              string                   `json:"entity"`
	AttrDefId           string                   `json:"attrDefId"`
	AttrDefName         string                   `json:"attrDefName"`
	AttrDefDataType     string                   `json:"attrDefDataType"`
	ElementType         string                   `json:"elementType"`
	Title               string                   `json:"title"`
	Width               int                      `json:"width"`
	RefPackageName      string                   `json:"refPackageName"`
	RefEntity           string                   `json:"refEntity"`
	DataOptions         string                   `json:"dataOptions"`
	Required            string                   `json:"required"`
	Regular             string                   `json:"regular"`
	IsEdit              string                   `json:"isEdit"`
	IsView              string                   `json:"isView"`
	IsOutput            string                   `json:"isOutput"`
	InDisplayName       string                   `json:"inDisplayName"`
	IsRefInside         string                   `json:"isRefInside"`
	Multiple            string                   `json:"multiple"`
	DefaultClear        string                   `json:"defaultClear"`
	RefId               string                   `json:"refId"`               // 复制数据表单ID,数据表单删除该表单项时,需要删除审批表单,任务表单对应数据项
	RoutineExpression   string                   `json:"routineExpression"`   // 计算表达式
	ControlSwitch       string                   `json:"controlSwitch"`       // 控制审批/任务开关
	FormTemplateLibrary string                   `json:"formTemplateLibrary"` // 表单组件库id
	HiddenCondition     []*QueryRequestFilterObj `json:"hiddenCondition"`     // 隐藏条件
}

func (FormItemTemplateLibraryTable) TableName() string {
	return "form_item_template_library"
}

func ConvertFormItemTemplateLibraryDto2Model(dto *FormItemTemplateLibraryDto) *FormItemTemplateLibraryTable {
	var hiddenCondition string
	if len(dto.HiddenCondition) > 0 {
		byteArr, _ := json.Marshal(dto.HiddenCondition)
		if len(byteArr) > 0 {
			hiddenCondition = string(byteArr)
		}
	}
	return &FormItemTemplateLibraryTable{
		Id:                  dto.Id,
		Name:                dto.Name,
		Description:         dto.Description,
		ItemGroup:           dto.ItemGroup,
		ItemGroupName:       dto.ItemGroupName,
		DefaultValue:        dto.DefaultValue,
		Sort:                dto.Sort,
		PackageName:         dto.PackageName,
		Entity:              dto.Entity,
		AttrDefId:           dto.AttrDefId,
		AttrDefName:         dto.AttrDefName,
		AttrDefDataType:     dto.AttrDefDataType,
		ElementType:         dto.ElementType,
		Title:               dto.Title,
		Width:               dto.Width,
		RefPackageName:      dto.RefPackageName,
		RefEntity:           dto.RefEntity,
		DataOptions:         dto.DataOptions,
		Required:            dto.Required,
		Regular:             dto.Regular,
		IsEdit:              dto.IsEdit,
		IsView:              dto.IsView,
		IsOutput:            dto.IsOutput,
		InDisplayName:       dto.InDisplayName,
		IsRefInside:         dto.IsRefInside,
		Multiple:            dto.Multiple,
		DefaultClear:        dto.DefaultClear,
		RefId:               dto.RefId,
		RoutineExpression:   dto.RoutineExpression,
		ControlSwitch:       dto.ControlSwitch,
		FormTemplateLibrary: dto.FormTemplateLibrary,
		HiddenCondition:     hiddenCondition,
	}
}

func ConvertFormItemTemplateLibraryModel2Dto(list []*FormItemTemplateLibraryTable) []*FormItemTemplateLibraryDto {
	var dtoList []*FormItemTemplateLibraryDto
	if len(list) > 0 {
		for _, library := range list {
			var newHiddenCondition []*QueryRequestFilterObj
			if strings.TrimSpace(library.HiddenCondition) != "" {
				json.Unmarshal([]byte(library.HiddenCondition), &newHiddenCondition)
			}
			dto := &FormItemTemplateLibraryDto{
				Id:                  library.Id,
				Name:                library.Name,
				Description:         library.Description,
				ItemGroup:           library.ItemGroup,
				ItemGroupName:       library.ItemGroupName,
				DefaultValue:        library.DefaultValue,
				Sort:                library.Sort,
				PackageName:         library.PackageName,
				Entity:              library.Entity,
				AttrDefId:           library.AttrDefId,
				AttrDefName:         library.AttrDefName,
				AttrDefDataType:     library.AttrDefDataType,
				ElementType:         library.ElementType,
				Title:               library.Title,
				Width:               library.Width,
				RefPackageName:      library.RefPackageName,
				RefEntity:           library.RefEntity,
				DataOptions:         library.DataOptions,
				Required:            library.Required,
				Regular:             library.Regular,
				IsEdit:              library.IsEdit,
				IsView:              library.IsView,
				IsOutput:            library.IsOutput,
				InDisplayName:       library.InDisplayName,
				IsRefInside:         library.IsRefInside,
				Multiple:            library.Multiple,
				DefaultClear:        library.DefaultClear,
				RefId:               library.RefId,
				RoutineExpression:   library.RoutineExpression,
				ControlSwitch:       library.ControlSwitch,
				FormTemplateLibrary: library.FormTemplateLibrary,
				HiddenCondition:     newHiddenCondition,
			}
			dtoList = append(dtoList, dto)
		}
	}
	return dtoList
}

type FormItemTemplateLibraryTableSort []*FormItemTemplateLibraryTable

func (s FormItemTemplateLibraryTableSort) Len() int {
	return len(s)
}

func (s FormItemTemplateLibraryTableSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s FormItemTemplateLibraryTableSort) Less(i, j int) bool {
	return strings.Compare(s[i].Id, s[j].Id) < 0
}
