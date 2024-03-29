package models

type FormItemTemplateTable struct {
	Id                string           `json:"id" xorm:"'id' pk" primary-key:"id"`
	Name              string           `json:"name" xorm:"name"`
	Description       string           `json:"description" xorm:"description"`
	ItemGroup         string           `json:"itemGroup" xorm:"item_group"`
	ItemGroupName     string           `json:"itemGroupName" xorm:"item_group_name"`
	FormTemplate      string           `json:"formTemplate" xorm:"form_template"`
	DefaultValue      string           `json:"defaultValue" xorm:"default_value"`
	Sort              int              `json:"sort" xorm:"sort"`
	PackageName       string           `json:"packageName" xorm:"package_name"`
	Entity            string           `json:"entity" xorm:"entity"`
	AttrDefId         string           `json:"attrDefId" xorm:"attr_def_id"`
	AttrDefName       string           `json:"attrDefName" xorm:"attr_def_name"`
	AttrDefDataType   string           `json:"attrDefDataType" xorm:"attr_def_data_type"`
	ElementType       string           `json:"elementType" xorm:"element_type"`
	Title             string           `json:"title" xorm:"title"`
	Width             int              `json:"width" xorm:"width"`
	RefPackageName    string           `json:"refPackageName" xorm:"ref_package_name"`
	RefEntity         string           `json:"refEntity" xorm:"ref_entity"`
	DataOptions       string           `json:"dataOptions" xorm:"data_options"`
	Required          string           `json:"required" xorm:"required"`
	Regular           string           `json:"regular" xorm:"regular"`
	IsEdit            string           `json:"isEdit" xorm:"is_edit"`
	IsView            string           `json:"isView" xorm:"is_view"`
	IsOutput          string           `json:"isOutput" xorm:"is_output"`
	InDisplayName     string           `json:"inDisplayName" xorm:"in_display_name"`
	IsRefInside       string           `json:"isRefInside" xorm:"is_ref_inside"`
	Multiple          string           `json:"multiple" xorm:"multiple"`
	DefaultClear      string           `json:"defaultClear" xorm:"default_clear"`
	RefId             string           `json:"refId" xorm:"ref_id"`                         //复制数据表单ID,数据表单删除该表单项时,需要删除审批表单,任务表单对应数据项
	RoutineExpression string           `json:"routineExpression" xorm:"routine_expression"` //计算表达式
	SelectList        []*EntityDataObj `json:"selectList" xorm:"-"`
	Active            bool             `json:"active" xorm:"-"` // 是否选中状态
}

type FormItemTemplateDto struct {
	Id                string           `json:"id"`
	Name              string           `json:"name"`
	Description       string           `json:"description"`
	FormTemplate      string           `json:"itemGroupId"`
	ItemGroup         string           `json:"itemGroup"`
	ItemGroupType     string           `json:"itemGroupType"` //表单组类型:workflow 编排数据,optional 自选,custom 自定义
	ItemGroupName     string           `json:"itemGroupName"`
	ItemGroupSort     int              `json:"ItemGroupSort"` // item_group 排序
	ItemGroupRule     string           `json:"itemGroupRule"` // item_group_rule 新增一行规则,new 输入新数据,exist 选择已有数据
	DefaultValue      string           `json:"defaultValue"`
	Sort              int              `json:"sort"`
	PackageName       string           `json:"packageName"`
	Entity            string           `json:"entity"`
	AttrDefId         string           `json:"attrDefId"`
	AttrDefName       string           `json:"attrDefName"`
	AttrDefDataType   string           `json:"attrDefDataType"`
	ElementType       string           `json:"elementType"`
	Title             string           `json:"title"`
	Width             int              `json:"width"`
	RefPackageName    string           `json:"refPackageName"`
	RefEntity         string           `json:"refEntity"`
	DataOptions       string           `json:"dataOptions"`
	Required          string           `json:"required"`
	Regular           string           `json:"regular"`
	IsEdit            string           `json:"isEdit"`
	IsView            string           `json:"isView"`
	IsOutput          string           `json:"isOutput"`
	InDisplayName     string           `json:"inDisplayName"`
	IsRefInside       string           `json:"isRefInside"`
	Multiple          string           `json:"multiple"`
	DefaultClear      string           `json:"defaultClear"`
	RefId             string           `json:"copyId"`            //复制数据表单ID,数据表单删除该表单项时,需要删除审批表单,任务表单对应数据项
	RoutineExpression string           `json:"routineExpression"` //计算表达式
	SelectList        []*EntityDataObj `json:"selectList"`
	Active            bool             `json:"active"` // 是否选中状态
}

func (FormItemTemplateTable) TableName() string {
	return "form_item_template"
}

type FormItemTemplateDtoSort []*FormItemTemplateDto

func (s FormItemTemplateDtoSort) Len() int {
	return len(s)
}

func (s FormItemTemplateDtoSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s FormItemTemplateDtoSort) Less(i, j int) bool {
	if s[i].Sort == 0 {
		return false
	}
	if s[j].Sort == 0 {
		return true
	}
	if s[i].Sort < s[j].Sort {
		return true
	}
	return false
}

func ConvertFormItemTemplateDto2Model(dto *FormItemTemplateDto) *FormItemTemplateTable {
	return &FormItemTemplateTable{
		Id:                dto.Id,
		Name:              dto.Name,
		Description:       dto.Description,
		FormTemplate:      dto.FormTemplate,
		ItemGroup:         dto.ItemGroup,
		ItemGroupName:     dto.ItemGroupName,
		DefaultValue:      dto.DefaultValue,
		Sort:              dto.Sort,
		PackageName:       dto.PackageName,
		Entity:            dto.Entity,
		AttrDefId:         dto.AttrDefId,
		AttrDefName:       dto.AttrDefName,
		AttrDefDataType:   dto.AttrDefDataType,
		ElementType:       dto.ElementType,
		Title:             dto.Title,
		Width:             dto.Width,
		RefPackageName:    dto.RefPackageName,
		RefEntity:         dto.RefEntity,
		DataOptions:       dto.DataOptions,
		Required:          dto.Required,
		Regular:           dto.Regular,
		IsEdit:            dto.IsEdit,
		IsView:            dto.IsView,
		IsOutput:          dto.IsOutput,
		InDisplayName:     dto.InDisplayName,
		IsRefInside:       dto.IsRefInside,
		Multiple:          dto.Multiple,
		DefaultClear:      dto.DefaultClear,
		RefId:             dto.RefId,
		SelectList:        dto.SelectList,
		Active:            dto.Active,
		RoutineExpression: dto.RoutineExpression,
	}
}

func ConvertFormItemTemplateModel2Dto(model *FormItemTemplateTable, itemGroup FormTemplateTable) *FormItemTemplateDto {
	dto := &FormItemTemplateDto{
		Id:                model.Id,
		Name:              model.Name,
		Description:       model.Description,
		ItemGroup:         model.ItemGroup,
		ItemGroupName:     model.ItemGroupName,
		FormTemplate:      model.FormTemplate,
		DefaultValue:      model.DefaultValue,
		Sort:              model.Sort,
		PackageName:       model.PackageName,
		Entity:            model.Entity,
		AttrDefId:         model.AttrDefId,
		AttrDefName:       model.AttrDefName,
		AttrDefDataType:   model.AttrDefDataType,
		ElementType:       model.ElementType,
		Title:             model.Title,
		Width:             model.Width,
		RefPackageName:    model.RefPackageName,
		RefEntity:         model.RefEntity,
		DataOptions:       model.DataOptions,
		Required:          model.Required,
		Regular:           model.Regular,
		IsEdit:            model.IsEdit,
		IsView:            model.IsView,
		IsOutput:          model.IsOutput,
		InDisplayName:     model.InDisplayName,
		IsRefInside:       model.IsRefInside,
		Multiple:          model.Multiple,
		DefaultClear:      model.DefaultClear,
		RefId:             model.RefId,
		SelectList:        model.SelectList,
		Active:            model.Active,
		RoutineExpression: model.RoutineExpression,
	}
	dto.ItemGroupType = itemGroup.ItemGroupType
	dto.ItemGroupRule = itemGroup.ItemGroupRule
	dto.ItemGroupSort = itemGroup.ItemGroupSort
	return dto
}

func ConvertFormItemTemplateModelList2Dto(tableList []*FormItemTemplateTable, itemGroup *FormTemplateTable) []*FormItemTemplateDto {
	var dtoList []*FormItemTemplateDto
	for _, model := range tableList {
		dtoList = append(dtoList, ConvertFormItemTemplateModel2Dto(model, *itemGroup))
	}
	return dtoList
}
