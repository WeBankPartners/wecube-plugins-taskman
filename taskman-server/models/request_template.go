package models

import "time"

type RequestTemplateTable struct {
	Id           string    `json:"id" xorm:"id"`
	Group        string    `json:"group" xorm:"group"`
	Name         string    `json:"name" xorm:"name"`
	Description  string    `json:"description" xorm:"description"`
	FormTemplate string    `json:"formTemplate" xorm:"form_template"`
	Tags         string    `json:"tags" xorm:"tags"`
	Status       string    `json:"status" xorm:"status"`
	PackageName  string    `json:"packageName" xorm:"package_name"`
	EntityName   string    `json:"entityName" xorm:"entity_name"`
	ProcDefKey   string    `json:"procDefKey" xorm:"proc_def_key"`
	ProcDefId    string    `json:"procDefId" xorm:"proc_def_id"`
	ProcDefName  string    `json:"procDefName" xorm:"proc_def_name"`
	CreatedBy    string    `json:"createdBy" xorm:"created_by"`
	CreatedTime  time.Time `json:"createdTime" xorm:"created_time"`
	UpdatedBy    string    `json:"updatedBy" xorm:"updated_by"`
	UpdatedTime  time.Time `json:"updatedTime" xorm:"updated_time"`
	EntityAttrs  string    `json:"entityAttrs" xorm:"entity_attrs"`
	DelFlag      int       `json:"delFlags" xorm:"del_flags"`
}

type RequestTemplateGroupTable struct {
	Id          string    `json:"id" xorm:"id"`
	Name        string    `json:"name" xorm:"name" binding:"required"`
	Description string    `json:"description" xorm:"description"`
	ManageRole  string    `json:"manageRole" xorm:"manage_role" binding:"required"`
	CreatedBy   string    `json:"createdBy" xorm:"created_by"`
	CreatedTime time.Time `json:"createdTime" xorm:"created_time"`
	UpdatedBy   string    `json:"updatedBy" xorm:"updated_by"`
	UpdatedTime time.Time `json:"updatedTime" xorm:"updated_time"`
	DelFlag     int       `json:"delFlags" xorm:"del_flags"`
}

type RoleTable struct {
	Id          string    `json:"id" xorm:"id"`
	DisplayName string    `json:"displayName" xorm:"display_name"`
	UpdatedTime time.Time `json:"updated_time" xorm:"updated_time"`
}

type RequestTemplateRoleTable struct {
	Id              string `json:"id" xorm:"id"`
	RequestTemplate string `json:"request_template" xorm:"request_template"`
	Role            string `json:"role" xorm:"role"`
	RoleType        string `json:"role_type" xorm:"role_type"`
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
	MGMTRoles []*RoleTable `json:"mgmtRoles"`
	USERoles  []*RoleTable `json:"useRoles"`
}

type RequestTemplateUpdateParam struct {
	RequestTemplateTable
	MGMTRoles []string `json:"mgmtRoles"`
	USERoles  []string `json:"useRoles"`
}
