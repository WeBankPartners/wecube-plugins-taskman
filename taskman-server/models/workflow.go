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

type ExecutionResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
