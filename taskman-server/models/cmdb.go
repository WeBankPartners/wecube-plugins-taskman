package models

type SysCiTypeAttrTable struct {
	Id                      string `json:"ciTypeAttrId" xorm:"id"`
	CiType                  string `json:"ciTypeId" xorm:"ci_type"`
	Name                    string `json:"propertyName" xorm:"name"`
	DisplayNameTmp          string `json:"displayName" xorm:"-"`
	DisplayName             string `json:"name" xorm:"display_name"`
	Description             string `json:"description" xorm:"description"`
	Status                  string `json:"status" xorm:"status"`
	InputType               string `json:"inputType" xorm:"input_type"`
	DataType                string `json:"propertyType" xorm:"data_type"`
	DataLength              int    `json:"length" xorm:"data_length"`
	TextValidate            string `json:"regularExpressionRule" xorm:"text_validate"`
	RefCiType               string `json:"referenceId" xorm:"ref_ci_type"`
	RefName                 string `json:"referenceName" xorm:"ref_name"`
	RefType                 string `json:"referenceType" xorm:"ref_type"`
	RefFilter               string `json:"referenceFilter" xorm:"ref_filter"`
	RefUpdateStateValidate  string `json:"refUpdateStateValidate" xorm:"ref_update_state_validate"`
	RefConfirmStateValidate string `json:"refConfirmStateValidate" xorm:"ref_confirm_state_validate"`
	SelectList              string `json:"selectList" xorm:"select_list"`
	UiSearchOrder           int    `json:"uiSearchOrder" xorm:"ui_search_order"`
	UiFormOrder             int    `json:"uiFormOrder" xorm:"ui_form_order"`
	UniqueConstraint        string `json:"uniqueConstraint" xorm:"unique_constraint"`
	UiNullable              string `json:"uiNullable" xorm:"ui_nullable"`
	Nullable                string `json:"nullable" xorm:"nullable"`
	Editable                string `json:"editable" xorm:"editable"`
	DisplayByDefault        string `json:"displayByDefault" xorm:"display_by_default"`
	PermissionUsage         string `json:"permissionUsage" xorm:"permission_usage"`
	ResetOnEdit             string `json:"resetOnEdit" xorm:"reset_on_edit"`
	Source                  string `json:"source" xorm:"source"`
	Customizable            string `json:"customizable" xorm:"customizable"`
	AutofillAble            string `json:"autofillable" xorm:"autofillable"`
	AutofillRule            string `json:"autoFillRule" xorm:"autofill_rule"`
	AutofillType            string `json:"autoFillType" xorm:"autofill_type"`
	EditGroupControl        string `json:"editGroupControl" xorm:"edit_group_control"`
	EditGroupValues         string `json:"editGroupValues" xorm:"edit_group_value"`
}

type CiTypeAttrQueryResponse struct {
	StatusCode string                `json:"statusCode"`
	Data       []*SysCiTypeAttrTable `json:"data"`
}

type CiReferenceDataQueryObj struct {
	Guid    string `json:"guid"`
	KeyName string `json:"key_name"`
	IsNew   bool   `json:"is_new"`
}

type CiReferenceDataQueryResponse struct {
	StatusCode string                     `json:"statusCode"`
	Data       []*CiReferenceDataQueryObj `json:"data"`
}

type CiDataRefFilterRight struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

type CiDataRefFilterObj struct {
	Left     string               `json:"left"`
	Operator string               `json:"operator"`
	Right    CiDataRefFilterRight `json:"right"`
}

type RefSelectParam struct {
	AttrId    string             `json:"attrId"`
	RequestId string             `json:"requestId"`
	UserToken string             `json:"userToken"`
	Filter    string             `json:"filter"`
	Param     *QueryRequestParam `json:"param"`
}

type GetExpressResultParam struct {
	Filter      *CiDataRefFilterObj               `json:"filter"`
	StartCiType string                            `json:"startCiType"`
	Express     string                            `json:"express"`
	UserToken   string                            `json:"userToken"`
	FilterMap   map[string]string                 `json:"filterMap"`
	NewData     map[string]map[string]interface{} `json:"newData"`
}
