package models

type RequestTemplateTable struct {
	Id           string `json:"id" xorm:"id"`
	Group        string `json:"group" xorm:"group"`
	Name         string `json:"name" xorm:"name"`
	Description  string `json:"description" xorm:"description"`
	FormTemplate string `json:"formTemplate" xorm:"form_template"`
	Tags         string `json:"tags" xorm:"tags"`
	Status       string `json:"status" xorm:"status"`
	RecordId     string `json:"recordId" xorm:"record_id"`
	Version      string `json:"version" xorm:"version"`
	ConfirmTime  string `json:"confirmTime" xorm:"confirm_time"`
	PackageName  string `json:"packageName" xorm:"package_name"`
	EntityName   string `json:"entityName" xorm:"entity_name"`
	ProcDefKey   string `json:"procDefKey" xorm:"proc_def_key"`
	ProcDefId    string `json:"procDefId" xorm:"proc_def_id"`
	ProcDefName  string `json:"procDefName" xorm:"proc_def_name"`
	CreatedBy    string `json:"createdBy" xorm:"created_by"`
	CreatedTime  string `json:"createdTime" xorm:"created_time"`
	UpdatedBy    string `json:"updatedBy" xorm:"updated_by"`
	UpdatedTime  string `json:"updatedTime" xorm:"updated_time"`
	EntityAttrs  string `json:"entityAttrs" xorm:"entity_attrs"`
	ExpireDay    int    `json:"expireDay" xorm:"expire_day"`
	Handler      string `json:"handler" xorm:"handler"`
	DelFlag      int    `json:"delFlag" xorm:"del_flag"`
}

type RequestTemplateGroupTable struct {
	Id            string    `json:"id" xorm:"id"`
	Name          string    `json:"name" xorm:"name" binding:"required"`
	Description   string    `json:"description" xorm:"description"`
	ManageRole    string    `json:"manageRole" xorm:"manage_role" binding:"required"`
	ManageRoleObj RoleTable `json:"manageRoleObj" xorm:"-"`
	CreatedBy     string    `json:"createdBy" xorm:"created_by"`
	CreatedTime   string    `json:"createdTime" xorm:"created_time"`
	UpdatedBy     string    `json:"updatedBy" xorm:"updated_by"`
	UpdatedTime   string    `json:"updatedTime" xorm:"updated_time"`
	DelFlag       int       `json:"delFlag" xorm:"del_flag"`
}

type RoleTable struct {
	Id          string `json:"id" xorm:"id"`
	DisplayName string `json:"displayName" xorm:"display_name"`
	UpdatedTime string `json:"updatedTime" xorm:"updated_time"`
	CoreId      string `json:"coreId" xorm:"core_id"`
}

type RequestTemplateRoleTable struct {
	Id              string `json:"id" xorm:"id"`
	RequestTemplate string `json:"requestTemplate" xorm:"request_template"`
	Role            string `json:"role" xorm:"role"`
	RoleType        string `json:"roleType" xorm:"role_type"`
}

type CoreProcessQueryResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Data    []*CodeProcessQueryObj `json:"data"`
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

type RequestTemplateQueryObj struct {
	RequestTemplateTable
	MGMTRoles      []*RoleTable `json:"mgmtRoles"`
	USERoles       []*RoleTable `json:"useRoles"`
	OperateOptions []string     `json:"operateOptions"`
}

type RequestTemplateUpdateParam struct {
	RequestTemplateTable
	MGMTRoles []string `json:"mgmtRoles"`
	USERoles  []string `json:"useRoles"`
}

type ProcEntityAttributeObj struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	DataType          string `json:"dataType"`
	Mandatory         bool   `json:"mandatory"`
	RefPackageName    string `json:"refPackageName"`
	RefEntityName     string `json:"refEntityName"`
	RefAttrName       string `json:"refAttrName"`
	ReferenceId       string `json:"referenceId"`
	Active            bool   `json:"active"`
	EntityId          string `json:"entityId"`
	EntityName        string `json:"entityName"`
	EntityDisplayName string `json:"entityDisplayName"`
	EntityPackage     string `json:"entityPackage"`
	Multiple          string `json:"multiple"`
}

type ProcEntity struct {
	Id          string                    `json:"id"`
	PackageName string                    `json:"packageName"`
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	DisplayName string                    `json:"displayName"`
	Attributes  []*ProcEntityAttributeObj `json:"attributes"`
}

type ProcDefObj struct {
	ProcDefId   string     `json:"procDefId"`
	ProcDefKey  string     `json:"procDefKey"`
	ProcDefName string     `json:"procDefName"`
	Status      string     `json:"status"`
	RootEntity  ProcEntity `json:"rootEntity"`
	CreatedTime string     `json:"createdTime"`
}

type ProcQueryResponse struct {
	Status  string        `json:"status"`
	Message string        `json:"message"`
	Data    []*ProcDefObj `json:"data"`
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

type ProcNodeQueryResponse struct {
	Status  string         `json:"status"`
	Message string         `json:"message"`
	Data    []*ProcNodeObj `json:"data"`
}

type UserRequestTemplateQueryObj struct {
	GroupId          string                       `json:"groupId"`
	GroupName        string                       `json:"groupName"`
	GroupDescription string                       `json:"groupDescription"`
	Templates        []*RequestTemplateTable      `json:"-"`
	Tags             []*UserRequestTemplateTagObj `json:"tags"`
}

type UserRequestTemplateTagObj struct {
	Tag       string                  `json:"tag"`
	Templates []*RequestTemplateTable `json:"templates"`
}

type RequestTemplateFormStruct struct {
	Id            string                    `json:"id"`
	Name          string                    `json:"name"`
	PackageName   string                    `json:"packageName"`
	EntityName    string                    `json:"entityName"`
	ProcDefKey    string                    `json:"procDefKey"`
	ProcDefId     string                    `json:"procDefId"`
	ProcDefName   string                    `json:"procDefName"`
	FormItems     []*FormItemTemplateTable  `json:"formItems"`
	TaskTemplates []*TaskTemplateFormStruct `json:"taskTemplates"`
}

type TaskTemplateFormStruct struct {
	Id          string                   `json:"id"`
	Name        string                   `json:"name"`
	NodeDefId   string                   `json:"nodeDefId"`
	NodeDefName string                   `json:"nodeDefName"`
	FormItems   []*FormItemTemplateTable `json:"formItems"`
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
