package models

type ProcessDefinitionsResponse struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Data    *DefinitionsData `json:"data"`
}
type FlowNodes struct {
	NodeID            string      `json:"nodeId"`
	NodeName          string      `json:"nodeName"`
	NodeType          string      `json:"nodeType"`
	NodeDefID         string      `json:"nodeDefId"`
	Status            string      `json:"status"`
	OrderedNo         string      `json:"orderedNo"`
	ProcDefID         string      `json:"procDefId"`
	ProcDefKey        string      `json:"procDefKey"`
	RoutineExpression string      `json:"routineExpression"`
	TaskCategory      interface{} `json:"taskCategory"`
	ServiceID         string      `json:"serviceId"`
	DynamicBind       string      `json:"dynamicBind"`
	Description       string      `json:"description"`
	PreviousNodeIds   []string    `json:"previousNodeIds"`
	SucceedingNodeIds []string    `json:"succeedingNodeIds"`
}
type DefinitionsData struct {
	ProcDefID        string      `json:"procDefId"`
	ProcDefKey       string      `json:"procDefKey"`
	ProcDefName      string      `json:"procDefName"`
	ProcDefVersion   string      `json:"procDefVersion"`
	Status           string      `json:"status"`
	ProcDefData      interface{} `json:"procDefData"`
	RootEntity       string      `json:"rootEntity"`
	CreatedTime      interface{} `json:"createdTime"`
	ExcludeMode      string      `json:"excludeMode"`
	Tags             interface{} `json:"tags"`
	PermissionToRole interface{} `json:"permissionToRole"`
	FlowNodes        []FlowNodes `json:"flowNodes"`
}

type ProcessInstanceResponse struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Data    *ProcessInstance `json:"data"`
}

type TaskNodeInstances struct {
	NodeID            string        `json:"nodeId"`
	NodeName          string        `json:"nodeName"`
	NodeType          string        `json:"nodeType"`
	NodeDefID         string        `json:"nodeDefId"`
	Status            string        `json:"status"`
	OrderedNo         interface{}   `json:"orderedNo"`
	ProcDefID         string        `json:"procDefId"`
	ProcDefKey        string        `json:"procDefKey"`
	RoutineExpression interface{}   `json:"routineExpression"`
	TaskCategory      interface{}   `json:"taskCategory"`
	ServiceID         interface{}   `json:"serviceId"`
	DynamicBind       interface{}   `json:"dynamicBind"`
	Description       interface{}   `json:"description"`
	PreviousNodeIds   []interface{} `json:"previousNodeIds"`
	SucceedingNodeIds []string      `json:"succeedingNodeIds"`
	ProcInstID        int           `json:"procInstId"`
	ProcInstKey       string        `json:"procInstKey"`
	ID                int           `json:"id"`
}
type ProcessInstance struct {
	ID                int                 `json:"id"`
	ProcInstKey       string              `json:"procInstKey"`
	ProcInstName      string              `json:"procInstName"`
	CreatedTime       string              `json:"createdTime"`
	Operator          string              `json:"operator"`
	Status            string              `json:"status"`
	ProcDefID         string              `json:"procDefId"`
	EntityTypeID      string              `json:"entityTypeId"`
	EntityDataID      string              `json:"entityDataId"`
	TaskNodeInstances []TaskNodeInstances `json:"taskNodeInstances"`
}

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
