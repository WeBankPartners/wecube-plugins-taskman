package models

type TaskFormInput struct {
	FormMetaID       string        `json:"formMetaId"`
	ProcDefID        string        `json:"procDefId"`
	ProcDefKey       string        `json:"procDefKey"`
	ProcInstID       int           `json:"procInstId"`
	ProcInstKey      string        `json:"procInstKey"`
	TaskNodeDefID    string        `json:"taskNodeDefId"`
	TaskNodeInstID   int           `json:"taskNodeInstId"`
	FormDataEntities []interface{} `json:"formDataEntities"`
}
type Inputs struct {
	TaskFormInput   TaskFormInput `json:"taskFormInput"`
	TaskDescription string        `json:"taskDescription"`
	RoleName        string        `json:"roleName"`
	CallbackURL     string        `json:"callbackUrl"`
	TaskName        string        `json:"taskName"`
	Reporter        string        `json:"reporter"`
	ProcInstID      int           `json:"procInstId"`
}
type Outputs struct {
	TaskFormOutput string `json:"taskFormOutput"`
	ErrorMessage   string `json:"errorMessage"`
	ErrorCode      string `json:"errorCode"`
	Comment        string `json:"comment"`
}
type RequestObjects struct {
	CallbackParameter string    `json:"callbackParameter"`
	Inputs            []Inputs  `json:"inputs"`
	Outputs           []Outputs `json:"outputs"`
}
type ExecutionNode struct {
	NodeName       string           `json:"nodeName"`
	NodeID         string           `json:"nodeId"`
	NodeDefID      string           `json:"nodeDefId"`
	NodeInstID     int              `json:"nodeInstId"`
	NodeType       string           `json:"nodeType"`
	NodeExpression string           `json:"nodeExpression"`
	PluginInfo     string           `json:"pluginInfo"`
	RequestID      string           `json:"requestId"`
	BeginTime      string           `json:"beginTime"`
	EndTime        string           `json:"endTime"`
	RequestObjects []RequestObjects `json:"requestObjects"`
}

type ExecutionResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    *ExecutionNode `json:"data"`
}
