package models

import "time"

type TaskTemplateTable struct {
	Id                string    `json:"id" xorm:"id"`
	Name              string    `json:"name" xorm:"name"`
	Description       string    `json:"description" xorm:"description"`
	FormTemplate      string    `json:"form_template" json:"form_template"`
	RequestTemplateId string    `json:"request_template_id" json:"request_template_id"`
	ProcDefId         string    `json:"proc_def_id" json:"proc_def_id"`
	ProcDefKey        string    `json:"proc_def_key"`
	ProcDefName       string    `json:"proc_def_name"`
	NodeDefId         string    `json:"node_def_id"`
	NodeName          string    `json:"node_name"`
	CreatedBy         string    `json:"created_by" xorm:"created_by"`
	CreatedTime       time.Time `json:"created_time" xorm:"created_time"`
	UpdatedBy         string    `json:"updated_by" xorm:"updated_by"`
	UpdatedTime       time.Time `json:"updated_time" xorm:"updated_time"`
	DelFlag           int       `json:"del_flags" xorm:"del_flags"`
}
