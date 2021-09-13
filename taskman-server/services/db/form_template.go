package db

import (
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"strconv"
	"time"
)

func GetRequestFormTemplate(id string) (result models.RequestFormTemplateDto, err error) {
	var formTemplateTable []*models.FormTemplateTable
	err = x.SQL("select * from form_template where id=?", id).Find(&formTemplateTable)
	if err != nil {
		err = fmt.Errorf("Try to query form template table fail,%s ", err.Error())
		return
	}
	if len(formTemplateTable) == 0 {
		err = fmt.Errorf("Can not find any form template with id=%s ", id)
		return
	}
	result.Id = formTemplateTable[0].Id
	result.Name = formTemplateTable[0].Name
	result.Description = formTemplateTable[0].Description
	result.UpdatedTime = formTemplateTable[0].UpdatedTime
	var formItemTemplate []*models.FormItemTemplateTable
	x.SQL("select * from form_item_template where form_template=? and version is null", id).Find(&formItemTemplate)
	result.Items = formItemTemplate
	return
}

func CreateRequestFormTemplate(param models.RequestFormTemplateDto, operator string) (id string, err error) {
	param.Id = guid.CreateGuid()
	id = param.Id
	itemIds := guid.CreateGuidList(len(param.Items))
	nowTime := time.Now().Format(models.DateTimeFormat)
	var actions []*execAction
	insertAction := execAction{Sql: "insert into form_template(id,name,description,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?)"}
	insertAction.Param = []interface{}{param.Id, param.Name, param.Description, operator, nowTime, operator, nowTime}
	actions = append(actions, &insertAction)
	for i, item := range param.Items {
		tmpAction := execAction{Sql: "insert into form_item_template(id,record_id,version,form_template,name,description,default_value,sort,package_name,entity,attr_def_id,attr_def_name,attr_def_data_type,element_type,title,width,ref_package_name,ref_entity,data_options,required,regular,is_edit,is_view) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
		tmpAction.Param = []interface{}{itemIds[i], item.RecordId, item.Version, item.FormTemplate, item.Name, item.Description, item.DefaultValue, item.Sort, item.PackageName, item.Entity, item.AttrDefId, item.AttrDefName, item.AttrDefDataType, item.ElementType, item.Title, item.Width, item.RefPackageName, item.RefEntity, item.DataOptions, item.Required, item.Required, item.Regular, item.IsEdit, item.IsView}
		actions = append(actions, &tmpAction)
	}
	err = transactionWithoutForeignCheck(actions)
	return
}

func UpdateRequestFormTemplate(param models.RequestFormTemplateDto, operator string) error {
	var formTemplateTable []*models.FormTemplateTable
	err := x.SQL("select id,version,updated_time from form_template where id=?", param.Id).Find(&formTemplateTable)
	if err != nil {
		return fmt.Errorf("Try to query form template table fail,%s ", err.Error())
	}
	if len(formTemplateTable) == 0 {
		return fmt.Errorf("Can not find any form template with id=%s ", param.Id)
	}
	if param.UpdatedTime != formTemplateTable[0].UpdatedTime {
		return fmt.Errorf("Update time validate fail,please refersh data ")
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	var actions []*execAction
	updateAction := execAction{Sql: "update form_template set name=?,description=?,updated_by=?,updated_time=? where id=?"}
	updateAction.Param = []interface{}{param.Name, param.Description, operator, nowTime, param.Id}
	actions = append(actions, &updateAction)
	newItemGuidList := guid.CreateGuidList(len(param.Items))
	var formItemTemplate []*models.FormItemTemplateTable
	x.SQL("select id,version,record_id from form_item_template where form_template=? and version is null", param.Id).Find(&formItemTemplate)
	for i, inputItem := range param.Items {
		tmpAction := execAction{}
		if inputItem.Id == "" {
			inputItem.Id = newItemGuidList[i]
			tmpAction.Sql = "insert into form_item_template(id,form_template,name,description,default_value,sort,package_name,entity,attr_def_id,attr_def_name,attr_def_data_type,element_type,title,width,ref_package_name,ref_entity,data_options,required,regular,is_edit,is_view) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
			tmpAction.Param = []interface{}{inputItem.Id, inputItem.FormTemplate, inputItem.Name, inputItem.Description, inputItem.DefaultValue, inputItem.Sort, inputItem.PackageName, inputItem.Entity, inputItem.AttrDefId, inputItem.AttrDefName, inputItem.AttrDefDataType, inputItem.ElementType, inputItem.Title, inputItem.Width, inputItem.RefPackageName, inputItem.RefEntity, inputItem.DataOptions, inputItem.Required, inputItem.Required, inputItem.Regular, inputItem.IsEdit, inputItem.IsView}
		} else {
			tmpAction.Sql = "update form_item_template set name=?,description=?,default_value=?,sort=?,package_name=?,entity=?,attr_def_id=?,attr_def_name=?,attr_def_data_type=?,element_type=?,title=?,width=?,ref_package_name=?,ref_entity=?,data_options=?,required=?,regular=?,is_edit=?,is_view=? where id=?"
			tmpAction.Param = []interface{}{inputItem.Name, inputItem.Description, inputItem.DefaultValue, inputItem.Sort, inputItem.PackageName, inputItem.Entity, inputItem.AttrDefId, inputItem.AttrDefName, inputItem.AttrDefDataType, inputItem.ElementType, inputItem.Title, inputItem.Width, inputItem.RefPackageName, inputItem.RefEntity, inputItem.DataOptions, inputItem.Required, inputItem.Required, inputItem.Regular, inputItem.IsEdit, inputItem.IsView, inputItem.Id}
		}
		actions = append(actions, &tmpAction)
	}
	for _, existItem := range formItemTemplate {
		existFlag := false
		for _, inputItem := range param.Items {
			if existItem.Id == inputItem.Id {
				existFlag = true
				break
			}
		}
		if !existFlag {
			actions = append(actions, &execAction{Sql: "delete from form_item_template where id=?", Param: []interface{}{existItem.Id}})
		}
	}
	return transaction(actions)
}

func ConfirmRequestFormTemplate(id, operator string) error {
	var formTemplateTable []*models.FormTemplateTable
	err := x.SQL("select id,version,updated_time from form_template where id=?", id).Find(&formTemplateTable)
	if err != nil {
		return fmt.Errorf("Try to query form template table fail,%s ", err.Error())
	}
	if len(formTemplateTable) == 0 {
		return fmt.Errorf("Can not find any form template with id=%s ", id)
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	version := buildVersionNum(formTemplateTable[0].Version)
	var actions []*execAction
	actions = append(actions, &execAction{Sql: "update form_template set `version`=?,confirm_time=?,updated_by=?,updated_time=? where id=?", Param: []interface{}{version, nowTime, operator, nowTime, id}})
	var formItemTemplate []*models.FormItemTemplateTable
	x.SQL("select id from form_item_template where form_template=? and version is null", id).Find(&formItemTemplate)
	newGuidList := guid.CreateGuidList(len(formItemTemplate))
	for i, item := range formItemTemplate {
		actions = append(actions, &execAction{Sql: "update form_item_template set `version`=? where id=?", Param: []interface{}{version, item.Id}})
		actions = append(actions, &execAction{Sql: "insert into form_item_template(id,record_id,form_template,name,description,default_value,sort,package_name,entity,attr_def_id,attr_def_name,attr_def_data_type,element_type,title,width,ref_package_name,ref_entity,data_options,required,regular,is_edit,is_view) select '" + newGuidList[i] + "' as id,id as record_id,form_template,name,description,default_value,sort,package_name,entity,attr_def_id,attr_def_name,attr_def_data_type,element_type,title,width,ref_package_name,ref_entity,data_options,required,regular,is_edit,is_view from form_item_template where id=?", Param: []interface{}{item.Id}})
	}
	return transaction(actions)
}

func DeleteRequestFormTemplate(id string) error {
	_, err := x.Exec("update form_template set del_flag=1 where id=?", id)
	return err
}

func buildVersionNum(version string) string {
	if version == "" {
		return "v1"
	}
	tmpV, err := strconv.Atoi(version[:1])
	if err != nil {
		return fmt.Sprintf("v%d", time.Now().Unix())
	}
	return fmt.Sprintf("v%d", tmpV+1)
}
