package models

import "strings"

type RequestTable struct {
	Id                  string             `json:"id" xorm:"'id' pk"`
	Name                string             `json:"name" xorm:"name"`
	Form                string             `json:"form" xorm:"form"`
	RequestTemplate     string             `json:"requestTemplate" xorm:"request_template"`
	RequestTemplateName string             `json:"requestTemplateName" xorm:"-"`
	ProcInstanceId      string             `json:"procInstanceId" xorm:"proc_instance_id"`
	ProcInstanceKey     string             `json:"procInstanceKey" xorm:"proc_instance_key"`
	Reporter            string             `json:"reporter" xorm:"reporter"`
	Handler             string             `json:"handler" xorm:"handler"`
	ReportTime          string             `json:"reportTime" xorm:"report_time"`
	Emergency           int                `json:"emergency" xorm:"emergency"`
	ReportRole          string             `json:"reportRole" xorm:"report_role"`
	Status              string             `json:"status" xorm:"status"`
	Cache               string             `json:"cache" xorm:"cache"`
	BindCache           string             `json:"bindCache" xorm:"bind_cache"`
	Result              string             `json:"result" xorm:"result"`
	ExpireTime          string             `json:"expireTime" xorm:"expire_time"`
	ExpectTime          string             `json:"expectTime" xorm:"expect_time"`
	ConfirmTime         string             `json:"confirmTime" xorm:"confirm_time"`
	CreatedBy           string             `json:"createdBy" xorm:"created_by"`
	CreatedTime         string             `json:"createdTime" xorm:"created_time"`
	UpdatedBy           string             `json:"updatedBy" xorm:"updated_by"`
	UpdatedTime         string             `json:"updatedTime" xorm:"updated_time"`
	DelFlag             int                `json:"delFlag" xorm:"del_flag"`
	HandleRoles         []string           `json:"handleRoles" xorm:"-"`
	AttachFiles         []*AttachFileTable `json:"attachFiles" xorm:"-"`
	Parent              string             `json:"parent" xorm:"parent"`
	CompletedTime       string             `json:"completedTime" xorm:"-"`
	RollbackDesc        string             `json:"rollbackDesc" xorm:"rollback_desc"`
	Type                int                `json:"type" xorm:"type"`
	OperatorObj         string             `json:"operatorObj" xorm:"operator_obj"`
	Description         string             `json:"description" xorm:"description"` // 请求描述
	Role                string             `json:"role" xorm:"role"`               // 创建请求的role
	RevokeFlag          int                `json:"revokeFlag" xorm:"revoke_flag"`  // 撤回标志 0表示没被撤回,1表示撤回
	ExpireDay           int                `json:"expireDay" xorm:"-"`             // 模板过期时间
	TemplateVersion     string             `json:"templateVersion" xorm:"-"`       // 模板版本
}

func (RequestTable) TableName() string {
	return "request"
}

// PlatformData  工作台数据
type PlatformData struct {
	Pending      string `json:"pending"`      // 待处理, eg:7;2 使用;分割开 7表示发布个数,2表示请求个数
	HasProcessed string `json:"hasProcessed"` // 已处理
	Submit       string `json:"submit"`       // 我提交的
	Draft        string `json:"draft"`        // 我暂存的
	Collect      string `json:"collect"`      // 收藏模板
}

// PlatformDataObj 工作台返回数据
type PlatformDataObj struct {
	Id                   string `json:"id" xorm:"'id' pk"`                          // 请求ID
	Name                 string `json:"name" xorm:"name"`                           // 请求名称
	TemplateId           string `json:"templateId" xorm:"template_id"`              // 模板ID
	TemplateName         string `json:"templateName" xorm:"template_name"`          // 使用模板名称
	Type                 int    `json:"type" xorm:"type"`                           // 模板类型: 0表示请求,1表示发布
	Version              string `json:"version" xorm:"version"`                     // 模板版本
	OperatorObj          string `json:"operatorObj" xorm:"operator_obj"`            // 操作对象
	OperatorObjType      string `json:"operatorObjType" xorm:"operator_obj_type"`   // 操作对象类型
	ProcInstanceId       string `json:"procInstanceId" xorm:"proc_instance_id"`     // 编排实例id
	ProcDefId            string `json:"procDefId" xorm:"proc_def_id"`               // 编排 key
	ProcDefKey           string `json:"procDefKey" xorm:"proc_def_key"`             // 编排id
	ProcDefName          string `json:"procDefName" xorm:"proc_def_name"`           // 使用编排
	Status               string `json:"status" xorm:"status"`                       // 请求状态
	CurNode              string `json:"curNode"  xorm:"cur_node"`                   // 当前节点
	Progress             int    `json:"progress" xorm:"progress"`                   // 进展
	CreatedBy            string `json:"createdBy" xorm:"created_by"`                // 创建人
	Handler              string `json:"handler" xorm:"handler"`                     // 当前处理人
	CreatedTime          string `json:"createdTime" xorm:"created_time"`            // 创建时间
	ReportTime           string `json:"reportTime" xorm:"report_time"`              // 请求提交时间
	ExpectTime           string `json:"expectTime" xorm:"expect_time"`              // 期望完成时间
	UpdatedTime          string `json:"updatedTime" xorm:"updated_time"`            // 更新时间
	ApprovalTime         string `json:"approvalTime" xorm:"approval_time"`          // 请求处理时间
	CollectFlag          int    `json:"collectFlag" xorm:"collect_flag"`            // 收藏标记,1表示已收藏
	Role                 string `json:"role" xorm:"role"`                           // 创建请求Role
	HandleRole           string `json:"handleRole" xorm:"handle_role"`              // 处理role
	RollbackDesc         string `json:"rollbackDesc" xorm:"rollback_desc"`          // 回退原因
	RevokeFlag           int    `json:"revokeFlag" xorm:"revoke_flag"`              // 是否撤回,0表示否,1表示撤回
	StartTime            string `json:"startTime" xorm:"-"`                         // 开始时间
	EffectiveDays        int    `json:"effectiveDays" xorm:"-"`                     // 有效天数
	ParentId             string `json:"parentId" xorm:"parent_id"`                  // 模板父类id
	Cache                string `json:"-" xorm:"cache"`                             // request cache
	TaskId               string `json:"taskId" xorm:"task_id"`                      // 当前正在进行中的taskId
	TaskName             string `json:"taskName" xorm:"task_name"`                  // taskName
	TaskCreatedTime      string `json:"taskCreatedTime" xorm:"task_created_time"`   // task创建时间
	TaskApprovalTime     string `json:"taskApprovalTime" xorm:"task_approval_time"` // 任务处理时间
	TaskExpectTime       string `json:"taskExpectTime" xorm:"task_expect_time"`     // 任务期望完成时间
	TaskHandler          string `json:"taskHandler" xorm:"task_handler"`            // 任务审批人
	TaskUpdatedTime      string `json:"taskUpdatedTime" xorm:"task_updated_time"`   // 任务更新时间
	TaskStatus           string `json:"taskStatus" xorm:"task_status"`              // 当前任务状态
	RequestStayTime      int    `json:"requestStayTime" xorm:"-"`                   // 请求停留时长
	RequestStayTimeTotal int    `json:"requestStayTimeTotal" xorm:"-"`              // 请求停留时长总数
	TaskStayTime         int    `json:"taskStayTime" xorm:"-"`                      // 任务停留时长
	TaskStayTimeTotal    int    `json:"taskStayTimeTotal" xorm:"-"`                 // 任务停留时长总数
}

type RequestProgressObj struct {
	NodeId    string `json:"nodeId" xorm:"node_id"`
	NodeDefId string `json:"nodeDefId" xorm:"node_def_id"`
	Node      string `json:"node" xorm:"node"`
	Handler   string `json:"handler" xorm:"handler"`
	Status    int    `json:"status" xorm:"status"` // 状态值：1 进行中 2.未开始  3.已完成  4.报错被拒绝了
	OrderedNo string `json:"orderNo" xorm:"-"`     // 排序字段
}

type ExpireObj struct {
	Percent    float64 `json:"percent"`
	ReportTime string  `json:"reportTime"`
	ExpireTime string  `json:"expireTime"`
	NowTime    string  `json:"nowTime"`
	TotalDay   float64 `json:"totalDay"`
	LeftDay    float64 `json:"leftDay"`
	UseDay     float64 `json:"useDay"`
}

type RequestCacheData struct {
	ProcDefId         string                         `json:"procDefId"`
	ProcDefKey        string                         `json:"procDefKey"`
	RootEntityValue   RequestCacheEntityValue        `json:"rootEntityValue"`
	TaskNodeBindInfos []*RequestCacheTaskNodeBindObj `json:"taskNodeBindInfos"`
}

type RequestCacheTaskNodeBindObj struct {
	BoundEntityValues []*RequestCacheEntityValue `json:"boundEntityValues"`
	NodeDefId         string                     `json:"nodeDefId"`
	NodeId            string                     `json:"nodeId"`
}

type RequestCacheEntityValue struct {
	AttrValues        []*RequestCacheEntityAttrValue `json:"attrValues"`
	BindFlag          string                         `json:"bindFlag"`
	EntityDataId      string                         `json:"entityDataId"`
	EntityDataOp      string                         `json:"entityDataOp"`
	EntityDataState   string                         `json:"entityDataState"`
	EntityDefId       string                         `json:"entityDefId"`
	EntityName        string                         `json:"entityName"`
	EntityDisplayName string                         `json:"entityDisplayName"`
	FullEntityDataId  interface{}                    `json:"fullEntityDataId"`
	Oid               string                         `json:"oid"`
	PackageName       string                         `json:"packageName"`
	PreviousOids      []string                       `json:"previousOids"`
	Processed         bool                           `json:"processed"`
	SucceedingOids    []string                       `json:"succeedingOids"`
}

type RequestCacheEntityAttrValue struct {
	DataOid   string      `json:"-"`
	AttrDefId string      `json:"attrDefId"`
	AttrName  string      `json:"attrName"`
	DataType  string      `json:"dataType"`
	DataValue interface{} `json:"dataValue"`
}

type RequestProcessData struct {
	ProcDefId     string                           `json:"procDefId"`
	ProcDefKey    string                           `json:"procDefKey"`
	RootEntityOid string                           `json:"rootEntityOid"`
	Entities      []*RequestCacheEntityValue       `json:"entities"`
	Bindings      []*RequestProcessTaskNodeBindObj `json:"bindings"`
}

type RequestProcessTaskNodeBindObj struct {
	NodeId       string `json:"nodeId"`
	NodeDefId    string `json:"nodeDefId"`
	Oid          string `json:"oid"`
	EntityDataId string `json:"entityDataId"`
	BindFlag     string `json:"bindFlag"`
}

type RequestPreDataTableObj struct {
	PackageName   string                   `json:"packageName"`
	Entity        string                   `json:"entity"`
	ItemGroup     string                   `json:"itemGroup"`
	ItemGroupName string                   `json:"itemGroupName"`
	RefEntity     []string                 `json:"-"`
	SortLevel     int                      `json:"-"`
	Title         []*FormItemTemplateTable `json:"title"`
	Value         []*EntityTreeObj         `json:"value"`
}

type RequestPreDataSort []*RequestPreDataTableObj

func (s RequestPreDataSort) Len() int {
	return len(s)
}

func (s RequestPreDataSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s RequestPreDataSort) Less(i, j int) bool {
	if s[i].SortLevel == s[j].SortLevel {
		return s[i].Entity < s[j].Entity
	}
	return s[i].SortLevel > s[j].SortLevel
}

type RequestPreDataDto struct {
	RootEntityId string                    `json:"rootEntityId"`
	EntityName   string                    `json:"entityName"` // 操作单元
	Data         []*RequestPreDataTableObj `json:"data"`
}

type RequestProDataV2Dto struct {
	Name         string                    `json:"name"`
	Description  string                    `json:"description"`
	ExpectTime   string                    `json:"expectTime"` // 期望完成时间
	EntityName   string                    `json:"entityName"` // 操作单元
	RootEntityId string                    `json:"rootEntityId"`
	Data         []*RequestPreDataTableObj `json:"data"`
}

type RequestForm struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	RequestType       int    `json:"requestType"`       // 请求类型,0表示请求,1表示发布
	Progress          int    `json:"progress"`          // 请求进度
	Status            string `json:"status"`            // 请求状态
	CurNode           string `json:"curNode"`           // 当前节点
	Handler           string `json:"handler"`           // 当前处理人
	CreatedBy         string `json:"createdBy"`         // 创建人
	Role              string `json:"role"`              // 创建人角色
	TemplateName      string `json:"templateName"`      // 使用模板
	Version           string `json:"version"`           // 模板版本
	TemplateGroupName string `json:"templateGroupName"` // 使用模板组
	Description       string `json:"description"`       // 请求描述
	CreatedTime       string `json:"createdTime"`       // 创建时间
	ExpectTime        string `json:"expectTime" `       // 期望时间
	OperatorObj       string `json:"operatorObj"`       // 发布操作对象
	ProcInstanceId    string `json:"procInstanceId"`    // 编排实例ID
	ExpireDay         int    `json:"expireDay"`         // 模板过期时间
}

type FilterItem struct {
	TemplateList        []*KeyValuePair `json:"templateList"`        // 模板列表
	RequestTemplateList []*KeyValuePair `json:"requestTemplateList"` // 请求模板列表
	ReleaseTemplateList []*KeyValuePair `json:"releaseTemplateList"` // 发布模板列表
	OperatorObjTypeList []string        `json:"operatorObjTypeList"` // 操作对象类型列表
	ProcDefNameList     []string        `json:"procDefNameList"`     // 使用编排
	CreatedByList       []string        `json:"createdByList"`       // 创建人列表
	HandlerList         []string        `json:"handlerList"`         // 处理人列表
}

type KeyValuePair struct {
	TemplateId   string `json:"templateId"`   // 使用模板
	TemplateName string `json:"templateName"` // 使用模板
	Version      string `json:"version"`      // 模板版本
}

type KeyValueSort []*KeyValuePair

func (q KeyValueSort) Len() int {
	return len(q)
}

func (q KeyValueSort) Less(i, j int) bool {
	t := strings.Compare(q[i].TemplateName, q[j].TemplateName)
	if t < 0 {
		return true
	}
	return false
}

func (q KeyValueSort) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

type FilterObj struct {
	Id              string `json:"id" xorm:"'id' pk"`                        // requestID
	TemplateId      string `json:"templateId" xorm:"template_id"`            // 模板IDa
	TemplateName    string `json:"templateName" xorm:"template_name"`        // 模板名称
	Version         string `json:"version" xorm:"version"`                   // 模板版本
	OperatorObjType string `json:"operatorObjType" xorm:"operator_obj_type"` // 操作对象类型
	ProcDefName     string `json:"procDefName" xorm:"proc_def_name"`         // 使用编排
	CreatedBy       string `json:"createdBy" xorm:"created_by"`              // 创建人
	Handler         string `json:"handler" xorm:"handler"`                   // 处理人
	TemplateType    int    `json:"-" xorm:"template_type"`                   // 模板类型
}

type RequestDetail struct {
	Request RequestForm     `json:"request"` // 请求信息
	Data    []*TaskQueryObj `json:"data"`
}

type UpdateRequestStatusParam struct {
	Description string `json:"description"`
}

type RequestProgressObjSort []*RequestProgressObj

func (q RequestProgressObjSort) Len() int {
	return len(q)
}

func (q RequestProgressObjSort) Less(i, j int) bool {
	t := strings.Compare(q[i].OrderedNo, q[j].OrderedNo)
	if t < 0 {
		return true
	}
	return false
}

func (q RequestProgressObjSort) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}
