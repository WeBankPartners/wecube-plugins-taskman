package models

import (
	"time"
)

type RequestTemplateTable struct {
	Id              string `json:"id"  xorm:"'id' pk" primary-key:"id"`
	Group           string `json:"group" xorm:"group"`
	Name            string `json:"name" xorm:"name"`
	Description     string `json:"description" xorm:"description"`
	Tags            string `json:"tags" xorm:"tags"`
	Status          string `json:"status" xorm:"status"`
	RecordId        string `json:"recordId" xorm:"record_id"`
	Version         string `json:"version" xorm:"version"`
	ConfirmTime     string `json:"confirmTime" xorm:"confirm_time"`
	PackageName     string `json:"packageName" xorm:"package_name"`
	EntityName      string `json:"entityName" xorm:"entity_name"`
	ProcDefKey      string `json:"procDefKey" xorm:"proc_def_key"`
	ProcDefId       string `json:"procDefId" xorm:"proc_def_id"`
	ProcDefName     string `json:"procDefName" xorm:"proc_def_name"`
	CreatedBy       string `json:"createdBy" xorm:"created_by"`
	CreatedTime     string `json:"createdTime" xorm:"created_time"`
	UpdatedBy       string `json:"updatedBy" xorm:"updated_by"`
	UpdatedTime     string `json:"updatedTime" xorm:"updated_time"`
	EntityAttrs     string `json:"entityAttrs" xorm:"entity_attrs"`
	ExpireDay       int    `json:"expireDay" xorm:"expire_day"`
	Handler         string `json:"handler" xorm:"handler"`
	DelFlag         int    `json:"delFlag" xorm:"del_flag"`
	Type            int    `json:"type" xorm:"type"`                         // 请求类型, 0表示请求,1表示发布,2为变更,3为事件,4为问题
	OperatorObjType string `json:"operatorObjType" xorm:"operator_obj_type"` // 操作对象类型
	ParentId        string `json:"parentId" xorm:"parent_id"`                // 父类ID
	ApproveBy       string `json:"approveBy" xorm:"approve_by"`              // 模板发布审批人
	CheckSwitch     bool   `json:"checkSwitch" xorm:"check_switch"`          // 是否加入确认定版流程
	ConfirmSwitch   bool   `json:"confirmSwitch" xorm:"confirm_switch"`      // 是否加入确认流程
	BackDesc        string `json:"backDesc" xorm:"back_desc"`                // 退回理由
}

func (RequestTemplateTable) TableName() string {
	return "request_template"
}

type RequestTemplateDto struct {
	Id               string `json:"id"`
	Group            string `json:"group"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	FormTemplate     string `json:"formTemplate"`
	Tags             string `json:"tags"`
	Status           string `json:"status"`
	RecordId         string `json:"recordId"`
	Version          string `json:"version"`
	ConfirmTime      string `json:"confirmTime"`
	PackageName      string `json:"packageName"`
	EntityName       string `json:"entityName"`
	ProcDefKey       string `json:"procDefKey"`
	ProcDefId        string `json:"procDefId"`
	ProcDefName      string `json:"procDefName"`
	CreatedBy        string `json:"createdBy"`
	CreatedTime      string `json:"createdTime"`
	UpdatedBy        string `json:"updatedBy"`
	UpdatedTime      string `json:"updatedTime"`
	EntityAttrs      string `json:"entityAttrs"`
	ExpireDay        int    `json:"expireDay"`
	Handler          string `json:"handler"`
	DelFlag          int    `json:"delFlag"`
	Type             int    `json:"type"`             // 请求类型,0表示请求,1表示发布
	OperatorObjType  string `json:"operatorObjType"`  // 操作对象类型
	ParentId         string `json:"parentId"`         // 父类ID
	ApproveBy        string `json:"approveBy"`        // 模板发布审批人
	CheckSwitch      bool   `json:"pendingSwitch"`    // 是否加入确认定版流程
	CheckRole        string `json:"pendingRole"`      // 定版角色
	CheckExpireDay   int    `json:"pendingExpireDay"` // 定版时效
	CheckHandler     string `json:"pendingHandler"`   // 定版处理人
	ConfirmSwitch    bool   `json:"confirmSwitch"`    // 是否加入确认流程
	ConfirmExpireDay int    `json:"confirmExpireDay"` // 确认过期时间
	BackDesc         string `json:"rollbackDesc"`     // 退回理由
}

// CollectDataObj 收藏数据项
type CollectDataObj struct {
	ParentId          string   `json:"parentId" xorm:"parent_id"`                    // 父类ID
	Id                string   `json:"id" xorm:"'id' pk"`                            // 模版ID
	Name              string   `json:"name" xorm:"name"`                             // 模版名称
	Version           string   `json:"version" xorm:"version"`                       // 模版名称
	Status            int      `json:"status" xorm:"status"`                         // 模版状态: 1可使用 2已禁用 3权限被移除
	TemplateGroupId   string   `json:"templateGourpId" xorm:"template_group_id"`     // 模版组ID
	TemplateGroup     string   `json:"templateGroup" xorm:"template_group"`          // 模版组
	TemplateGroupRole string   `json:"templateGroupRole" xorm:"template_group_role"` // 模版组角色
	OperatorObjType   string   `json:"operatorObjType" xorm:"operator_obj_type"`     // 操作对象类型
	ProcDefName       string   `json:"procDefName" xorm:"proc_def_name"`             // 使用编排
	ManageRole        string   `json:"manageRole" xorm:"manage_role"`                // 属主角色
	Owner             string   `json:"owner" xorm:"owner"`                           // 属主
	UseRole           string   `json:"useRole" xorm:"use_role"`                      // 使用角色
	Tags              string   `json:"tags" xorm:"tags"`                             // 标签
	WorkNode          []string `json:"workNode" xorm:"work_node"`                    // 人工任务
	Approves          []string `json:"approves" xorm:"-"`                            // 审批列表
	Tasks             []string `json:"tasks" xorm:"-"`                               // 任务节点
	CreatedTime       string   `json:"createdTime" xorm:"created_time"`              // 创建时间
	UpdatedTime       string   `json:"updatedTime" xorm:"updated_time"`              // 更新时间
}

type RequestTemplateHandlerDto struct {
	RequestTemplateId string `json:"request_template_id"` //模板id
	LatestUpdateTime  string `json:"latestUpdateTime"`    //最后更新时间
}

type RequestTemplateQueryObj struct {
	RequestTemplateDto
	MGMTRoles      []*RoleTable `json:"mgmtRoles"`
	USERoles       []*RoleTable `json:"useRoles"`
	OperateOptions []string     `json:"operateOptions"`
	ModifyType     bool         `json:"modifyType"`    // 是否能够修改模板类型
	Administrator  string       `json:"administrator"` // 角色管理员
}

type RequestTemplateStatusUpdateParam struct {
	RequestTemplateId string `json:"requestTemplateId"` // 请求模板ID
	Status            string `json:"status"`            // 当前状态
	TargetStatus      string `json:"targetStatus"`      // 目标状态
	Reason            string `json:"reason"`            // 原因
}

type RequestTemplateUpdateParam struct {
	RequestTemplateDto
	MGMTRoles []string `json:"mgmtRoles"`
	USERoles  []string `json:"useRoles"`
}

type RequestTemplateTableObj struct {
	Id              string `json:"id" xorm:"'id' pk"`
	Name            string `json:"name" xorm:"name"`
	Version         string `json:"version" xorm:"version"`
	Tags            string `json:"tags" xorm:"tags"`
	Status          string `json:"status" xorm:"status"`
	UpdatedBy       string `json:"updatedBy" xorm:"updated_by"`
	Handler         string `json:"handler" xorm:"handler"`
	Role            string `json:"role" xorm:"role"`
	RoleDisplay     string `json:"roleDisplay" xorm:"-"`
	UpdatedTime     string `json:"updatedTime" xorm:"updated_time"`
	CollectFlag     int    `json:"collectFlag" xorm:"collect_flag"`          // 是否收藏 1表示已收藏
	Type            int    `json:"type" xorm:"type"`                         // 请求类型, 0表示请求,1表示发布
	OperatorObjType string `json:"operatorObjType" xorm:"operator_obj_type"` // 操作对象类型
}

type UserRequestTemplateQueryObjNew struct {
	ManageRole        string              `json:"manageRole"`        //管理角色
	ManageRoleDisplay string              `json:"manageRoleDisplay"` //管理角色
	Groups            []*TemplateGroupObj `json:"groups"`
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

type RequestTemplateExport struct {
	RequestTemplate      RequestTemplateDto          `json:"requestTemplate"`
	FormTemplate         []*FormTemplateTable        `json:"formTemplate"`
	FormItemTemplate     []*FormItemTemplateTable    `json:"formItemTemplate"`
	RequestTemplateRole  []*RequestTemplateRoleTable `json:"requestTemplateRole"`
	TaskTemplate         []*TaskTemplateTable        `json:"taskTemplate"`
	TaskTemplateRole     []*TaskHandleTemplateTable  `json:"taskTemplateRole"`
	RequestTemplateGroup RequestTemplateGroupTable   `json:"requestTemplateGroup"`
}

type RequestTemplateTmp struct {
	ProcDefId         string `json:"procDefId" xorm:"proc_def_id`
	TemplateName      string `json:"templateName" xorm:"template_name`
	TemplateGroupName string `json:"templateGroupName" xorm:"template_group_name`
	Version           string `json:"version" xorm:"version"`
	Status            string `json:"status" xorm:"status"`
	ExpireDay         int    `json:"expireDay" xorm:"expire_day"`
}

type ImportData struct {
	Token        string `json:"token"`
	TemplateName string `json:"templateName"`
}

type RequestTemplateEntityDto struct {
	FormType string   `json:"formType"` //表单类型
	Entities []string `json:"entities"` //实例
}

type RequestTemplateSort []*RequestTemplateTableObj

func (s RequestTemplateSort) Len() int {
	return len(s)
}

func (s RequestTemplateSort) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s RequestTemplateSort) Less(i, j int) bool {
	return s[i].UpdatedTime > s[j].UpdatedTime
}

func ConvertRequestTemplateUpdateParam2RequestTemplate(param RequestTemplateUpdateParam) *RequestTemplateTable {
	nowTime := time.Now().Format(DateTimeFormat)
	return &RequestTemplateTable{
		Id:              param.Id,
		Group:           param.Group,
		Name:            param.Name,
		Description:     param.Description,
		Tags:            param.Tags,
		Status:          param.Status,
		PackageName:     param.PackageName,
		EntityName:      param.EntityName,
		ProcDefKey:      param.ProcDefKey,
		ProcDefId:       param.ProcDefId,
		ProcDefName:     param.ProcDefName,
		CreatedBy:       param.CreatedBy,
		CreatedTime:     nowTime,
		UpdatedBy:       param.CreatedBy,
		UpdatedTime:     nowTime,
		EntityAttrs:     "",
		ExpireDay:       param.ExpireDay,
		Handler:         param.Handler,
		DelFlag:         0,
		Type:            param.Type,
		OperatorObjType: param.OperatorObjType,
		ParentId:        param.Id,
		ApproveBy:       param.ApproveBy,
		CheckSwitch:     param.CheckSwitch,
		ConfirmSwitch:   param.ConfirmSwitch,
		BackDesc:        "",
	}
}

func ConvertRequestTemplateDto2Model(param RequestTemplateDto) *RequestTemplateTable {
	return &RequestTemplateTable{
		Id:              param.Id,
		Group:           param.Group,
		Name:            param.Name,
		Description:     param.Description,
		Tags:            param.Tags,
		Status:          param.Status,
		RecordId:        param.RecordId,
		Version:         param.Version,
		ConfirmTime:     param.ConfirmTime,
		PackageName:     param.PackageName,
		EntityName:      param.EntityName,
		ProcDefKey:      param.ProcDefKey,
		ProcDefId:       param.ProcDefId,
		ProcDefName:     param.ProcDefName,
		CreatedBy:       param.CreatedBy,
		CreatedTime:     param.CreatedTime,
		UpdatedBy:       param.UpdatedBy,
		UpdatedTime:     param.UpdatedTime,
		EntityAttrs:     param.EntityAttrs,
		ExpireDay:       param.ExpireDay,
		Handler:         param.Handler,
		DelFlag:         param.DelFlag,
		Type:            param.Type,
		OperatorObjType: param.OperatorObjType,
		ParentId:        param.ParentId,
		ApproveBy:       param.ApproveBy,
		CheckSwitch:     param.CheckSwitch,
		ConfirmSwitch:   param.ConfirmSwitch,
		BackDesc:        param.BackDesc,
	}
}

func ConvertRequestTemplateModel2Dto(requestTemplate *RequestTemplateTable) *RequestTemplateDto {
	return &RequestTemplateDto{
		Id:              requestTemplate.Id,
		Group:           requestTemplate.Group,
		Name:            requestTemplate.Name,
		Description:     requestTemplate.Description,
		Tags:            requestTemplate.Tags,
		Status:          requestTemplate.Status,
		RecordId:        requestTemplate.RecordId,
		Version:         requestTemplate.Version,
		ConfirmTime:     requestTemplate.ConfirmTime,
		PackageName:     requestTemplate.PackageName,
		EntityName:      requestTemplate.EntityName,
		ProcDefKey:      requestTemplate.ProcDefKey,
		ProcDefId:       requestTemplate.ProcDefId,
		ProcDefName:     requestTemplate.ProcDefName,
		CreatedBy:       requestTemplate.CreatedBy,
		CreatedTime:     requestTemplate.CreatedTime,
		UpdatedBy:       requestTemplate.UpdatedBy,
		UpdatedTime:     requestTemplate.UpdatedTime,
		EntityAttrs:     requestTemplate.EntityAttrs,
		ExpireDay:       requestTemplate.ExpireDay,
		Handler:         requestTemplate.Handler,
		DelFlag:         requestTemplate.DelFlag,
		Type:            requestTemplate.Type,
		OperatorObjType: requestTemplate.OperatorObjType,
		ParentId:        requestTemplate.ParentId,
		ApproveBy:       requestTemplate.ApproveBy,
		CheckSwitch:     requestTemplate.CheckSwitch,
		CheckRole:       "",
		CheckHandler:    "",
		ConfirmSwitch:   requestTemplate.ConfirmSwitch,
		BackDesc:        requestTemplate.BackDesc,
	}
}
