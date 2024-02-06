package models

type ProcDefQueryDto struct {
	ManageRole  string        `json:"manageRole"` //管理角色
	ProcDefList []*ProcDefDto `json:"dataList"`   // 编排列表
}

type QueryProcessDefinitionParam struct {
	ProcDefId        string   `json:"procDefId"`        // 编排Id
	ProcDefName      string   `json:"procDefName"`      // 编排名称
	Plugins          []string `json:"plugins"`          // 授权插件
	UpdatedTimeStart string   `json:"updatedTimeStart"` // 更新时间开始
	UpdatedTimeEnd   string   `json:"updatedTimeEnd"`   // 更新时间结束
	Status           string   `json:"status"`           // disabled 禁用 draft草稿 deployed 发布状态
	CreatedBy        string   `json:"createdBy"`        // 创建人
	UpdatedBy        string   `json:"updatedBy"`        // 更新人
	Scene            string   `json:"scene"`            // 使用场景
	UserRoles        []string // 用户角色
}

type ProcDefDto struct {
	Id               string   `json:"id"`               // 唯一标识
	Key              string   `json:"key"`              // 编排key
	Name             string   `json:"name"`             // 编排名称
	Version          string   `json:"version"`          // 版本
	RootEntity       string   `json:"rootEntity"`       // 根节点
	Status           string   `json:"status"`           // 状态
	Tags             string   `json:"tags"`             // 标签
	AuthPlugins      []string `json:"authPlugins"`      // 授权插件
	Scene            string   `json:"scene"`            // 使用场景
	ConflictCheck    bool     `json:"conflictCheck"`    // 冲突检测
	CreatedBy        string   `json:"createdBy"`        // 创建人
	CreatedTime      string   `json:"createdTime"`      // 创建时间
	UpdatedBy        string   `json:"updatedBy"`        // 更新人
	UpdatedTime      string   `json:"updatedTime"`      // 更新时间
	EnableCreated    bool     `json:"enableCreated"`    // 能否创建新版本
	EnableModifyName bool     `json:"enableModifyName"` // 能否修改名称
	UseRoles         []string `json:"userRoles"`        // 使用角色
}

type EntityProDefDto struct {
	Id          string `json:"id"`          // 唯一标识
	PackageName string `json:"packageName"` // 包名称
	Name        string `json:"name"`        // 名称
	DisplayName string `json:"displayName"` // 显示名称
	Description string `json:"description"` // 描述
}

type ProcDefObj struct {
	ProcDefId   string     `json:"procDefId"`
	ProcDefKey  string     `json:"procDefKey"`
	ProcDefName string     `json:"procDefName"`
	Status      string     `json:"status"`
	RootEntity  ProcEntity `json:"rootEntity"`
	CreatedTime string     `json:"createdTime"`
}

type ProcAllDefObj struct {
	ProcDefId   string `json:"procDefId"`
	ProcDefKey  string `json:"procDefKey"`
	ProcDefName string `json:"procDefName"`
	Status      string `json:"status"`
	RootEntity  string `json:"rootEntity"`
	CreatedTime string `json:"createdTime"`
}

type ProcAllQueryResponse struct {
	Status  string           `json:"status"`
	Message string           `json:"message"`
	Data    []*ProcAllDefObj `json:"data"`
}

type ProcNodeObj struct {
	NodeId        string        `json:"nodeId"`
	NodeName      string        `json:"nodeName"`
	NodeType      string        `json:"nodeType"`
	NodeDefId     string        `json:"nodeDefId"`
	TaskCategory  string        `json:"taskCategory"`
	RoutineExp    string        `json:"routineExp"`
	ServiceId     string        `json:"serviceId"`
	ServiceName   string        `json:"serviceName"`
	OrderedNo     string        `json:"orderedNo"`
	OrderedNum    int           `json:"-"`
	DynamicBind   string        `json:"dynamicbind"`
	BoundEntities []*ProcEntity `json:"boundEntities"`
}

type CodeProcessQueryObj struct {
	ExcludeMode     string `json:"excludeMode"`
	ProcDefId       string `json:"procDefId"`
	ProcDefKey      string `json:"procDefKey"`
	ProcDefName     string `json:"procDefName"`
	ProcDefVersion  string `json:"procDefVersion"`
	RootEntity      string `json:"rootEntity"`
	Status          string `json:"status"`
	CreatedTime     string `json:"createdTime"`
	CreatedUnixTime int64  `json:"-"`
	Tags            string `json:"tags"`
}

type WorkflowRsp struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Data    Workflow `json:"data"`
}

type Workflow struct {
	Status          string           `json:"status"`
	DefinitionsInfo *DefinitionsInfo `json:"define_data"`
	InstancesInfo   *InstancesInfo   `json:"instance_data"`
}

type DefinitionsInfo struct {
	ProcDefId        string          `json:"procDefId"`
	ProcDefKey       string          `json:"procDefKey"`
	ProcDefName      string          `json:"procDefName"`
	ProcDefVersion   string          `json:"procDefVersion"`
	Status           string          `json:"status"`
	ProcDefData      string          `json:"procDefData"`
	CreatedTime      string          `json:"createdTime"`
	ExcludeMode      string          `json:"excludeMode"`
	Tags             string          `json:"tags"`
	PermissionToRole string          `json:"permissionToRole"`
	FlowNodes        []*WorkflowNode `json:"flowNodes"`
}

type InstancesInfo struct {
	Id           int             `json:"id"`
	ProcInstKey  string          `json:"procInstKey"`
	ProcInstName string          `json:"procInstName"`
	CreatedTime  string          `json:"createdTime"`
	Operator     string          `json:"operator"`
	Status       string          `json:"status"`
	ProcDefId    string          `json:"procDefId"`
	EntityTypeId string          `json:"entityTypeId"`
	EntityDataId string          `json:"entityDataId"`
	TaskNodes    []*WorkflowNode `json:"taskNodeInstances"`
}

// WorkflowNode 任务编排节点
type WorkflowNode struct {
	Id                int      `json:"id"`
	ProInstId         int      `json:"proInstId"`
	ProInstKey        string   `json:"proInstKey"`
	NodeId            string   `json:"nodeId"`
	NodeName          string   `json:"nodeName"`
	NodeType          string   `json:"nodeType"`
	NodeDefId         string   `json:"nodeDefId"`
	Status            string   `json:"status"`
	OrderedNo         string   `json:"orderedNo"`
	ProcDefId         string   `json:"procDefId"`
	ProcDefKey        string   `json:"procDefKey"`
	RoutineExpression string   `json:"routineExpression"`
	TaskCategory      string   `json:"taskCategory"`
	ServiceId         string   `json:"serviceId"`
	DynamicBind       string   `json:"dynamicBind"`
	Description       string   `json:"description"`
	PreviousNodeIds   []string `json:"previousNodeIds"`
	SucceedingNodeIds []string `json:"succeedingNodeIds"`
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

type EntityNodeBindQueryObj struct {
	NodeDefId string `xorm:"node_def_id"`
	ItemGroup string `xorm:"item_group"`
}

type InstanceStatusQuery struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    InstanceStatusQueryObj `json:"data"`
}

type InstanceStatusQueryObj struct {
	Id                int                        `json:"id"`
	ProcDefId         string                     `json:"procDefId"`
	ProcInstKey       string                     `json:"procInstKey"`
	ProcInstName      string                     `json:"procInstName"`
	Status            string                     `json:"status"`
	TaskNodeInstances []*InstanceStatusQueryNode `json:"taskNodeInstances"`
}

type InstanceStatusQueryNode struct {
	Id        int    `json:"id"`
	NodeId    string `json:"nodeId"`
	NodeDefId string `json:"nodeDefId"`
	NodeName  string `json:"nodeName"`
	NodeType  string `json:"nodeType"`
	Status    string `json:"status"`
	OrderedNo string `json:"orderedNo"`
}

type StartInstanceResultData struct {
	Id          int    `json:"id"`
	ProcInstKey string `json:"procInstKey"`
	ProcDefId   string `json:"procDefId"`
	ProcDefKey  string `json:"procDefKey"`
	Status      string `json:"status"`
}

type ProcNodeObjList []*ProcNodeObj

func (s ProcNodeObjList) Len() int {
	return len(s)
}

func (s ProcNodeObjList) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s ProcNodeObjList) Less(i, j int) bool {
	return s[i].OrderedNum < s[j].OrderedNum
}

func ConvertModelsList2Map(nodesList []*DataModel) map[string]ProcEntity {
	var entityMap = make(map[string]ProcEntity)
	if len(nodesList) > 0 {
		for _, model := range nodesList {
			for _, entity := range model.Entities {
				if entity.PackageName != "" && entity.Name != "" {
					entityMap[entity.PackageName+":"+entity.Name] = ProcEntity{
						Id:          entity.DataModelId,
						PackageName: entity.PackageName,
						Name:        entity.Name,
						Description: entity.Description,
						DisplayName: entity.DisplayName,
						Attributes:  nil,
					}
				}
			}
		}
	}
	return entityMap
}
