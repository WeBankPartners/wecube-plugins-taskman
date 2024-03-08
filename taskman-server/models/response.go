package models

type PageInfo struct {
	StartIndex int `json:"startIndex"`
	PageSize   int `json:"pageSize"`
	TotalRows  int `json:"totalRows"`
}

type ResponsePageData struct {
	PageInfo PageInfo    `json:"pageInfo"`
	Contents interface{} `json:"contents"`
}

type ResponseJson struct {
	StatusCode string      `json:"statusCode"`
	Data       interface{} `json:"data"`
}

type ResponseErrorObj struct {
	ErrorMessage string `json:"errorMessage"`
}

type ResponseErrorJson struct {
	StatusCode    string      `json:"statusCode"`
	StatusMessage string      `json:"statusMessage"`
	Data          interface{} `json:"data"`
}

type SysLogTable struct {
	LogCat      string `json:"logCat" xorm:"log_cat"`
	Operator    string `json:"operator" xorm:"operator"`
	Operation   string `json:"operation" xorm:"operation"`
	Content     string `json:"content" xorm:"content"`
	RequestUrl  string `json:"requestUrl" xorm:"request_url"`
	ClientHost  string `json:"clientHost" xorm:"client_host"`
	CreatedDate string `json:"createdDate" xorm:"created_date"`
	DataCiType  string `json:"dataCiType" xorm:"data_ci_type"`
	DataGuid    string `json:"dataGuid" xorm:"data_guid"`
	DataKeyName string `json:"dataKeyName" xorm:"data_key_name"`
	Response    string `json:"response" xorm:"response"`
}

type HttpResponseMeta struct {
	Code    int    `json:"-"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type QueryRolesResponse struct {
	HttpResponseMeta
	Data []*SimpleLocalRoleDto `json:"data"`
}

type QueryUserResponse struct {
	HttpResponseMeta
	Data []*UserDto `json:"data"`
}

type QueryProcessDefinitionResponse struct {
	HttpResponseMeta
	Data []*ProcDefQueryDto `json:"data"`
}

type QueryProcessAllDefinitionResponse struct {
	HttpResponseMeta
	Data []*ProcDef `json:"data"`
}

type QueryAllModelsResponse struct {
	HttpResponseMeta
	Data []*DataModel `json:"data"`
}

type QueryExpressionEntitiesResponse struct {
	HttpResponseMeta
	Data []*ExpressionEntities `json:"data"`
}

type ProcNodeQueryResponse struct {
	HttpResponseMeta
	Data []*ProcNodeObj `json:"data"`
}

type ProcDefTaskNodesResponse struct {
	HttpResponseMeta
	Data []*ProcNodeObj `json:"data"`
}

type CoreProcessQueryResponse struct {
	HttpResponseMeta
	Data []*CodeProcessQueryObj `json:"data"`
}

type ProcessDefinitionsResponse struct {
	HttpResponseMeta
	Data *DefinitionsData `json:"data"`
}

type ProcessInstanceResponse struct {
	HttpResponseMeta
	Data *ProcessInstance `json:"data"`
}

type ProcDefTaskNodeContextResponse struct {
	HttpResponseMeta
	Data interface{} `json:"data"`
}

type ProcQueryResponse struct {
	HttpResponseMeta
	Data []*ProcDefObj `json:"data"`
}

type StartInstanceResponse struct {
	HttpResponseMeta
	Data *StartInstanceResultData `json:"data"`
}

type ProcDefRootEntityResponse struct {
	HttpResponseMeta
	Data []*ProcDefEntityDataObj `json:"data"`
}

type EntityTreeResponse struct {
	HttpResponseMeta
	Data *EntityTreeData `json:"data"`
}

type DataModelEntityResponse struct {
	HttpResponseMeta
	Data *DataModelEntity `json:"data"`
}

type RemoteLoginResp struct {
	HttpResponseMeta
	Data interface{} `json:"data"`
}
