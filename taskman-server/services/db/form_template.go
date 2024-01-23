package db

import (
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"strconv"
	"time"
)

func GetRequestFormTemplate(id string) (result models.FormTemplateDto, err error) {
	result = models.FormTemplateDto{Items: []*models.FormItemTemplateTable{}}
	requestTemplate, getErr := GetSimpleRequestTemplate(id)
	if getErr != nil {
		err = getErr
		return
	}
	var formTemplateTable []*models.FormTemplateTable
	err = x.SQL("select * from form_template where id=?", requestTemplate.FormTemplate).Find(&formTemplateTable)
	if err != nil {
		err = fmt.Errorf("Try to query form template table fail,%s ", err.Error())
		return
	}
	if len(formTemplateTable) == 0 {
		log.Logger.Warn("Can not find any data in form template", log.String("id", id))
		return
	}
	result.ExpireDay = requestTemplate.ExpireDay
	result.Id = formTemplateTable[0].Id
	result.Name = formTemplateTable[0].Name
	result.Description = formTemplateTable[0].Description
	result.UpdatedTime = formTemplateTable[0].UpdatedTime
	result.UpdatedBy = formTemplateTable[0].UpdatedBy
	var formItemTemplate []*models.FormItemTemplateTable
	x.SQL("select * from form_item_template where form_template=?", requestTemplate.FormTemplate).Find(&formItemTemplate)
	result.Items = formItemTemplate
	return
}

func CreateRequestFormTemplate(param models.FormTemplateDto, requestTemplateId string) error {
	param.NowTime = time.Now().Format(models.DateTimeFormat)
	insertFormActions, formId := getFormTemplateCreateActions(param)
	insertFormActions = append(insertFormActions, &execAction{Sql: "update request_template set form_template=?,expire_day=?,description=? where id=?", Param: []interface{}{formId, param.ExpireDay, param.Description, requestTemplateId}})
	return transactionWithoutForeignCheck(insertFormActions)
}

func getFormTemplateCreateActions(param models.FormTemplateDto) (actions []*execAction, id string) {
	param.Id = guid.CreateGuid()
	id = param.Id
	itemIds := guid.CreateGuidList(len(param.Items))
	insertAction := execAction{Sql: "insert into form_template(id,name,description,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?)"}
	insertAction.Param = []interface{}{param.Id, param.Name, param.Description, param.UpdatedBy, param.NowTime, param.UpdatedBy, param.NowTime}
	actions = append(actions, &insertAction)
	for i, item := range param.Items {
		tmpAction := execAction{Sql: "insert into form_item_template(id,form_template,name,description,item_group,item_group_name,default_value,sort,package_name,entity,attr_def_id,attr_def_name,attr_def_data_type,element_type,title,width,ref_package_name,ref_entity,data_options,required,regular,is_edit,is_view,is_output,in_display_name,is_ref_inside,multiple,default_clear) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
		tmpAction.Param = []interface{}{itemIds[i], id, item.Name, item.Description, item.ItemGroup, item.ItemGroupName, item.DefaultValue, item.Sort, item.PackageName, item.Entity, item.AttrDefId, item.AttrDefName, item.AttrDefDataType, item.ElementType, item.Title, item.Width, item.RefPackageName, item.RefEntity, item.DataOptions, item.Required, item.Regular, item.IsEdit, item.IsView, item.IsOutput, item.InDisplayName, item.IsRefInside, item.Multiple, item.DefaultClear}
		actions = append(actions, &tmpAction)
	}
	return
}

func UpdateRequestFormTemplate(param models.FormTemplateDto) error {
	param.NowTime = time.Now().Format(models.DateTimeFormat)
	updateActions, err := getFormTemplateUpdateActions(param)
	if err != nil {
		return err
	}
	updateActions = append(updateActions, &execAction{Sql: "update request_template set expire_day=?,description=? where form_template=?", Param: []interface{}{param.ExpireDay, param.Description, param.Id}})
	return transaction(updateActions)
}

func getFormTemplateUpdateActions(param models.FormTemplateDto) (actions []*execAction, err error) {
	var formTemplateTable []*models.FormTemplateTable
	err = x.SQL("select id,updated_time from form_template where id=?", param.Id).Find(&formTemplateTable)
	if err != nil {
		err = fmt.Errorf("Try to query form template table fail,%s ", err.Error())
		return
	}
	if len(formTemplateTable) == 0 {
		err = fmt.Errorf("Can not find any form template with id=%s ", param.Id)
		return
	}
	if param.UpdatedTime != formTemplateTable[0].UpdatedTime {
		err = fmt.Errorf("Update time validate fail,please refersh data ")
		return
	}
	updateAction := execAction{Sql: "update form_template set name=?,description=?,updated_by=?,updated_time=? where id=?"}
	updateAction.Param = []interface{}{param.Name, param.Description, param.UpdatedBy, param.NowTime, param.Id}
	actions = append(actions, &updateAction)
	newItemGuidList := guid.CreateGuidList(len(param.Items))
	var formItemTemplate []*models.FormItemTemplateTable
	x.SQL("select id from form_item_template where form_template=?", param.Id).Find(&formItemTemplate)
	for i, inputItem := range param.Items {
		tmpAction := execAction{}
		if inputItem.Id == "" {
			inputItem.Id = newItemGuidList[i]
			tmpAction.Sql = "insert into form_item_template(id,form_template,name,description,item_group,item_group_name,default_value,sort,package_name,entity,attr_def_id,attr_def_name,attr_def_data_type,element_type,title,width,ref_package_name,ref_entity,data_options,required,regular,is_edit,is_view,is_output,in_display_name,is_ref_inside,multiple,default_clear) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
			tmpAction.Param = []interface{}{inputItem.Id, param.Id, inputItem.Name, inputItem.Description, inputItem.ItemGroup, inputItem.ItemGroupName, inputItem.DefaultValue, inputItem.Sort, inputItem.PackageName, inputItem.Entity, inputItem.AttrDefId, inputItem.AttrDefName, inputItem.AttrDefDataType, inputItem.ElementType, inputItem.Title, inputItem.Width, inputItem.RefPackageName, inputItem.RefEntity, inputItem.DataOptions, inputItem.Required, inputItem.Regular, inputItem.IsEdit, inputItem.IsView, inputItem.IsOutput, inputItem.InDisplayName, inputItem.IsRefInside, inputItem.Multiple, inputItem.DefaultClear}
		} else {
			tmpAction.Sql = "update form_item_template set name=?,description=?,item_group=?,item_group_name=?,default_value=?,sort=?,package_name=?,entity=?,attr_def_id=?,attr_def_name=?,attr_def_data_type=?,element_type=?,title=?,width=?,ref_package_name=?,ref_entity=?,data_options=?,required=?,regular=?,is_edit=?,is_view=?,is_output=?,in_display_name=?,is_ref_inside=?,multiple=?,default_clear=? where id=?"
			tmpAction.Param = []interface{}{inputItem.Name, inputItem.Description, inputItem.ItemGroup, inputItem.ItemGroupName, inputItem.DefaultValue, inputItem.Sort, inputItem.PackageName, inputItem.Entity, inputItem.AttrDefId, inputItem.AttrDefName, inputItem.AttrDefDataType, inputItem.ElementType, inputItem.Title, inputItem.Width, inputItem.RefPackageName, inputItem.RefEntity, inputItem.DataOptions, inputItem.Required, inputItem.Regular, inputItem.IsEdit, inputItem.IsView, inputItem.IsOutput, inputItem.InDisplayName, inputItem.IsRefInside, inputItem.Multiple, inputItem.DefaultClear, inputItem.Id}
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
			actions = append(actions, &execAction{Sql: "delete from form_item where form_item_template=?", Param: []interface{}{existItem.Id}})
			actions = append(actions, &execAction{Sql: "delete from form_item_template where id=?", Param: []interface{}{existItem.Id}})
		}
	}
	return
}

func DeleteRequestFormTemplate(id string) error {
	_, err := x.Exec("update form_template set del_flag=1 where id=?", id)
	return err
}

func buildVersionNum(version string) string {
	if version == "" {
		return "v1"
	}
	tmpV, err := strconv.Atoi(version[1:])
	if err != nil {
		return fmt.Sprintf("v%d", time.Now().Unix())
	}
	return fmt.Sprintf("v%d", tmpV+1)
}
