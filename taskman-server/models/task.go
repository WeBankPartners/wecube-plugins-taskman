package models

type TaskTable struct {
	Id                string `json:"id" xorm:"id"`
	Name              string `json:"name" xorm:"name"`
	Description       string `json:"description" xorm:"description"`
	Form              string `json:"form" xorm:"form"`
	AttachFile        string `json:"attachFile" xorm:"attach_file"`
	Status            string `json:"status" xorm:"status"`
	Version           string `json:"version" xorm:"version"`
	Request           string `json:"request" xorm:"request"`
	Parent            string `json:"parent" xorm:"parent"`
	TaskTemplate      string `json:"taskTemplate" xorm:"task_template"`
	PackageName       string `json:"packageName" xorm:"package_name"`
	EntityName        string `json:"entityName" xorm:"entity_name"`
	ProcDefId         string `json:"procDefId" xorm:"proc_def_id"`
	ProcDefKey        string `json:"procDefKey" xorm:"proc_def_key"`
	ProcDefName       string `json:"procDefName" xorm:"proc_def_name"`
	NodeDefId         string `json:"nodeDefId" xorm:"node_def_id"`
	NodeName          string `json:"nodeName" xorm:"node_name"`
	CallbackUrl       string `json:"callbackUrl" xorm:"callback_url"`
	CallbackParameter string `json:"callbackParameter" xorm:"callback_parameter"`
	Emergency         string `json:"emergency" xorm:"emergency"`
	Result            string `json:"result" xorm:"result"`
	Cache             string `json:"cache" xorm:"cache"`
	CallbackRequestId string `json:"callbackRequestId" xorm:"callback_request_id"`
	Reporter          string `json:"reporter" xorm:"reporter"`
	ReportTime        string `json:"reportTime" xorm:"report_time"`
	ReportRole        string `json:"reportRole" xorm:"report_role"`
	Owner             string `json:"owner" xorm:"owner"`
	CreatedBy         string `json:"createdBy" xorm:"created_by"`
	CreatedTime       string `json:"createdTime" xorm:"created_time"`
	UpdatedBy         string `json:"updatedBy" xorm:"updated_by"`
	UpdatedTime       string `json:"updatedTime" xorm:"updated_time"`
	DelFlag           string `json:"delFlag" xorm:"del_flag"`
}

type TaskMetaResult struct {
	Status  string             `json:"status"`
	Message string             `json:"message"`
	Data    TaskMetaResultData `json:"data"`
}

type TaskMetaResultData struct {
	FormMetaId    string                `json:"formMetaId"`
	FormItemMetas []*TaskMetaResultItem `json:"formItemMetas"`
}

type TaskMetaResultItem struct {
	FormItemMetaId string `json:"formItemMetaId"`
	PackageName    string `json:"packageName"`
	EntityName     string `json:"entityName"`
	AttrName       string `json:"attrName"`
}

type PluginTaskCreateRequest struct {
	RequestId string                        `json:"requestId"`
	Inputs    []*PluginTaskCreateRequestObj `json:"inputs"`
}

type PluginTaskCreateRequestObj struct {
	CallbackParameter string `json:"callbackParameter"`
	ProcInstId        string `json:"procInstId"`
	CallbackUrl       string `json:"callbackUrl"`
	Reporter          string `json:"reporter"`
	RoleName          string `json:"roleName"`
	TaskName          string `json:"taskName"`
	TaskDescription   string `json:"taskDescription"`
	TaskFormInput     string `json:"taskFormInput"`
}

type PluginTaskCreateResp struct {
	ResultCode    string                 `json:"resultCode"`
	ResultMessage string                 `json:"resultMessage"`
	Results       PluginTaskCreateOutput `json:"results"`
}

type PluginTaskCreateOutput struct {
	RequestId string                       `json:"requestId"`
	Outputs   []*PluginTaskCreateOutputObj `json:"outputs"`
}

type PluginTaskCreateOutputObj struct {
	CallbackParameter string `json:"callbackParameter"`
	Comment           string `json:"comment"`
	TaskFormOutput    string `json:"taskFormOutput"`
	ErrorCode         string `json:"errorCode"`
	ErrorMessage      string `json:"errorMessage"`
	ErrorDetail       string `json:"errorDetail,omitempty"`
}

type PluginTaskFormDto struct {
	FormMetaId       string                  `json:"formMetaId"`
	ProcDefId        string                  `json:"procDefId"`
	ProcDefKey       string                  `json:"procDefKey"`
	ProcInstId       int                     `json:"procInstId"`
	ProcInstKey      string                  `json:"procInstKey"`
	TaskNodeDefId    string                  `json:"taskNodeDefId"`
	TaskNodeInstId   int                     `json:"taskNodeInstId"`
	FormDataEntities []*PluginTaskFormEntity `json:"formDataEntities"`
}

type PluginTaskFormEntity struct {
	FormMetaId       string                 `json:"formMetaId"`
	PackageName      string                 `json:"packageName"`
	EntityName       string                 `json:"entityName"`
	Oid              string                 `json:"oid"`
	EntityDataId     string                 `json:"entityDataId"`
	FullEntityDataId string                 `json:"fullEntityDataId"`
	EntityDataState  string                 `json:"entityDataState"`
	EntityDataOp     string                 `json:"entityDataOp"`
	BindFlag         string                 `json:"bindFlag"`
	FormItemValues   []*PluginTaskFormValue `json:"formItemValues"`
}

type PluginTaskFormValue struct {
	FormItemMetaId   string      `json:"formItemMetaId"`
	PackageName      string      `json:"packageName"`
	EntityName       string      `json:"entityName"`
	AttrName         string      `json:"attrName"`
	Oid              string      `json:"oid"`
	EntityDataId     string      `json:"entityDataId"`
	FullEntityDataId string      `json:"fullEntityDataId"`
	AttrValue        interface{} `json:"attrValue"`
}

type CallbackResult struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type TaskQueryResult struct {
	TimeStep []*TaskQueryTimeStep `json:"timeStep"`
	Data     []*TaskQueryObj      `json:"data"`
}

type TaskQueryTimeStep struct {
	RequestTemplateId string `json:"requestTemplateId"`
	TaskTemplateId    string `json:"taskTemplateId"`
	Name              string `json:"name"`
	Active            bool   `json:"active"`
}

type TaskQueryObj struct {
	RequestId   string                    `json:"requestId"`
	TaskId      string                    `json:"taskId"`
	Reporter    string                    `json:"reporter"`
	ReportTime  string                    `json:"reportTime"`
	Comment     string                    `json:"comment"`
	AttachFiles []string                  `json:"attachFiles"`
	Editable    bool                      `json:"editable"`
	Status      string                    `json:"status"`
	FormData    []*RequestPreDataTableObj `json:"formData"`
}
