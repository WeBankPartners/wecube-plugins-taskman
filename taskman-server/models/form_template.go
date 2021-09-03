package models

import "time"

type FormTemplateTable struct {
	Id          string    `json:"id" xorm:"id"`
	Name        string    `json:"name" xorm:"name"`
	Description string    `json:"description" xorm:"description"`
	Role        string    `json:"role" xorm:"role"`
	CreatedBy   string    `json:"created_by" xorm:"created_by"`
	CreatedTime time.Time `json:"created_time" xorm:"created_time"`
	UpdatedBy   string    `json:"updated_by" xorm:"updated_by"`
	UpdatedTime time.Time `json:"updated_time" xorm:"updated_time"`
	DelFlag     int       `json:"del_flags" xorm:"del_flags"`
}

type FormItemTemplateTable struct {
	Id              string `json:"id" xorm:"id"`
	Name            string `json:"name" xorm:"name"`
	Description     string `json:"description" xorm:"description"`
	FormTemplateId  string `json:"form_template_id" json:"form_template_id"`
	DefaultValue    string `json:"default_value" json:"default_value"`
	Sort            int    `json:"sort" json:"sort"`
	PackageName     string `json:"package_name" json:"package_name"`
	Entity          string `json:"entity" json:"entity"`
	AttrDefId       string `json:"attr_def_id" json:"attr_def_id"`
	AttrDefDataType string `json:"attr_def_data_type" json:"attr_def_data_type"`
	ElementType     string `json:"element_type" json:"element_type"`
	Title           string `json:"title" json:"title"`
	Width           int    `json:"width" json:"width"`
	RefPackageName  string `json:"ref_package_name" json:"ref_package_name"`
	RefEntity       string `json:"ref_entity" json:"ref_entity"`
	RefFilters      string `json:"ref_filters" json:"ref_filters"`
	DataOptions     string `json:"data_options" json:"data_options"`
	Required        int    `json:"required" json:"required"`
	Regular         string `json:"regular" json:"regular"`
	IsEdit          int    `json:"is_edit" json:"is_edit"`
	IsView          int    `json:"is_view" json:"is_view"`
}
