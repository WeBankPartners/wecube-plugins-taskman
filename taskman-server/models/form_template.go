package models

type FormTemplateTable struct {
	Id          string `json:"id" xorm:"'id' pk"`
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
	Id          string                   `json:"id"`
	Name        string                   `json:"name"`
	Description string                   `json:"description"`
	ExpireDay   int                      `json:"expireDay"`
	UpdatedBy   string                   `json:"updatedBy"`
	UpdatedTime string                   `json:"updatedTime"`
	NowTime     string                   `json:"-"`
	Items       []*FormItemTemplateTable `json:"items"`
}

// DataFormTemplateDto 全局表单模板 dto
type DataFormTemplateDto struct {
	FormTemplateId      string                      `json:"formTemplateId"`      // 全局表单模板ID
	AssociationWorkflow bool                        `json:"associationWorkflow"` // 是否关联编排
	UpdatedBy           string                      `json:"updatedBy"`           // 更新人
	UpdatedTime         string                      `json:"updatedTime"`         // 更新时间
	Groups              []*DataFormTemplateGroupDto `json:"groups"`
}

// DataFormTemplateGroupDto 全局表单模板组dto
type DataFormTemplateGroupDto struct {
	ItemGroup     string                   `json:"itemGroup"`
	ItemGroupType string                   `json:"itemGroupType"` // 表单组类型:workflow 编排数据,optional 自选,custom 自定义
	ItemGroupName string                   `json:"itemGroupName"`
	ItemGroupSort int                      `json:"itemGroupSort"` // 组排序
	Items         []*FormItemTemplateTable `json:"items"`         // 表单项
}

// FormTemplateGroupConfigureDto 表单组配置在dto
type FormTemplateGroupConfigureDto struct {
	ItemGroup     string                   `json:"itemGroup"`
	ItemGroupType string                   `json:"itemGroupType"` // 表单组类型:workflow 编排数据,optional 自选,custom 自定义
	ItemGroupName string                   `json:"itemGroupName"`
	ItemGroupRule string                   `json:"item_group_rule"` // item_group_rule 新增一行规则,new 输入新数据,exist 选择已有数据
	SystemItems   []*FormItemTemplateTable `json:"systemItems"`     // 系统表单项
	CustomItems   []*FormItemTemplateTable `json:"customItems"`     // 自定义表单项
}

type TaskFormItemQueryObj struct {
	Id               string `json:"id" xorm:"'id' pk"`
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
	var items = make([]*FormItemTemplateTable, 0)
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
		Id:          dto.FormTemplateId,
		Name:        "globalFormTemplate",
		Description: "",
		ExpireDay:   0,
		UpdatedBy:   dto.UpdatedBy,
		UpdatedTime: dto.UpdatedTime,
		Items:       items,
	}
}

type DataFormTemplateGroupDtoSort []*DataFormTemplateGroupDto

func (s DataFormTemplateGroupDtoSort) Len() int {
	return len(s)
}

func (s DataFormTemplateGroupDtoSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s DataFormTemplateGroupDtoSort) Less(i, j int) bool {
	if s[i].ItemGroupSort < s[j].ItemGroupSort {
		return true
	}
	return false
}
