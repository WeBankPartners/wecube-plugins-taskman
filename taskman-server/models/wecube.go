package models

type EntityQueryParam struct {
	Criteria          EntityQueryObj    `json:"criteria"`
	AdditionalFilters []*EntityQueryObj `json:"additionalFilters"`
}

type EntityQueryObj struct {
	AttrName  string      `json:"attrName"`
	Op        string      `json:"op"`
	Condition interface{} `json:"condition"`
}

type EntityResponse struct {
	Status  string                   `json:"status"`
	Message string                   `json:"message"`
	Data    []map[string]interface{} `json:"data"`
}

type SyncDataModelResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    []*SyncDataModelCiType `json:"data"`
}

type SyncDataModelCiType struct {
	Name        string                 `json:"name" xorm:"'id' pk"`
	DisplayName string                 `json:"displayName" xorm:"display_name"`
	Description string                 `json:"description" xorm:"description"`
	Attributes  []*SyncDataModelCiAttr `json:"attributes" xorm:"-"`
}

type SyncDataModelCiAttr struct {
	Name             string `json:"name" xorm:"name"`
	EntityName       string `json:"entityName" xorm:"ci_type"`
	Description      string `json:"description" xorm:"description"`
	DataType         string `json:"dataType" xorm:"input_type"`
	RefPackageName   string `json:"refPackageName" xorm:"-"`
	RefEntityName    string `json:"refEntityName" xorm:"ref_ci_type"`
	RefAttributeName string `json:"refAttributeName" xorm:"-"`
	Required         string `json:"required" xorm:"nullable"`
	Multiple         string `json:"multiple"`
}

type PluginCiDataOperationRequest struct {
	RequestId string                             `json:"requestId"`
	Inputs    []*PluginCiDataOperationRequestObj `json:"inputs"`
}

type PluginCiDataOperationRequestObj struct {
	CallbackParameter string `json:"callbackParameter"`
	CiType            string `json:"ciType"`
	Operation         string `json:"operation"`
	JsonData          string `json:"jsonData"`
}

type PluginCiDataAttrValueRequest struct {
	RequestId string                             `json:"requestId"`
	Inputs    []*PluginCiDataAttrValueRequestObj `json:"inputs"`
}

type PluginCiDataAttrValueRequestObj struct {
	CallbackParameter string `json:"callbackParameter"`
	CiType            string `json:"ciType"`
	Guid              string `json:"guid"`
	CiTypeAttr        string `json:"ciTypeAttr"`
	Value             string `json:"value"`
}

type PluginCiDataOperationResp struct {
	ResultCode    string                      `json:"resultCode"`
	ResultMessage string                      `json:"resultMessage"`
	Results       PluginCiDataOperationOutput `json:"results"`
}

type PluginCiDataOperationOutput struct {
	Outputs []*PluginCiDataOperationOutputObj `json:"outputs"`
}

type PluginCiDataOperationOutputObj struct {
	CallbackParameter string `json:"callbackParameter"`
	Guid              string `json:"guid"`
	ErrorCode         string `json:"errorCode"`
	ErrorMessage      string `json:"errorMessage"`
	ErrorDetail       string `json:"errorDetail,omitempty"`
}

type CoreRoleDto struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    []CoreRoleDataObj `json:"data"`
}

type CoreRoleDataObj struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	DisplayName string `json:"displayName"`
}

type CoreUserDto struct {
	Status  string            `json:"status"`
	Message string            `json:"message"`
	Data    []CoreUserDataObj `json:"data"`
}

type CoreUserDataObj struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type EntityAttributeObj struct {
	CiTypeAttrId            string `json:"ciTypeAttrId"`
	CiTypeId                string `json:"ciTypeId"`
	PropertyName            string `json:"propertyName"`
	DisplayName             string `json:"displayName"`
	Name                    string `json:"name"`
	Description             string `json:"description"`
	Status                  string `json:"status"`
	InputType               string `json:"inputType"`
	DataType                string `json:"propertyType"`
	DataLength              int    `json:"length"`
	RegularExpressionRule   string `json:"regularExpressionRule"`
	RefCiType               string `json:"referenceId"`
	RefName                 string `json:"referenceName"`
	RefType                 string `json:"referenceType"`
	RefFilter               string `json:"referenceFilter"`
	RefUpdateStateValidate  string `json:"refUpdateStateValidate"`
	RefConfirmStateValidate string `json:"refConfirmStateValidate"`
	SelectList              string `json:"selectList"`
	UiSearchOrder           int    `json:"uiSearchOrder"`
	UiFormOrder             int    `json:"uiFormOrder"`
	UniqueConstraint        string `json:"uniqueConstraint"`
	UiNullable              string `json:"uiNullable"`
	Nullable                string `json:"nullable"`
	Editable                string `json:"editable"`
	DisplayByDefault        string `json:"displayByDefault"`
	PermissionUsage         string `json:"permissionUsage"`
	ResetOnEdit             string `json:"resetOnEdit"`
	Source                  string `json:"source"`
	Customizable            string `json:"customizable"`
	AutofillAble            string `json:"autofillable"`
	AutofillRule            string `json:"autoFillRule"`
	AutofillType            string `json:"autoFillType"`
	EditGroupControl        string `json:"editGroupControl"`
	EditGroupValues         string `json:"editGroupValues"`
	ExtRefEntity            string `json:"extRefEntity"`
	ConfirmNullable         string `json:"confirmNullable"`
	Sensitive               string `json:"sensitive"`
}

type EntityAttributeQueryResponse struct {
	StatusCode string                `json:"statusCode"`
	Data       []*EntityAttributeObj `json:"data"`
}

type CMDBCategoriesObj struct {
	CatId string `json:"catId"`
	Code  string `json:"code"`
	Value string `json:"value"`
	SeqNo int    `json:"seqNo"`
}

type CMDBCategoriesResponse struct {
	StatusCode string               `json:"statusCode"`
	Data       []*CMDBCategoriesObj `json:"data"`
}

type AttrPermissionQueryObj struct {
	HistoryId        int    `json:"historyId"`
	CiType           string `json:"ciType"`
	AttrName         string `json:"attrName"`
	Guid             string `json:"guid"`
	TmpId            string `json:"tmpId"`
	RequestId        string `json:"requestId"`
	TaskHandleId     string `json:"taskHandleId"`
	QueryPermission  bool   `json:"queryPermission"`
	UpdatePermission bool   `json:"updatePermission"`
	Value            string `json:"value"`
}

type CMDBSensitiveDataResponse struct {
	Code          int                       `json:"code"`
	StatusCode    string                    `json:"statusCode"`
	StatusMessage string                    `json:"statusMessage"`
	Data          []*AttrPermissionQueryObj `json:"data"`
}
