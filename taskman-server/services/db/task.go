package db

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
)

func GetTaskFormStruct(procInstId, nodeDefId string) (result models.TaskMetaResult, err error) {
	result = models.TaskMetaResult{Status: "OK", Message: "Success"}
	var items []*models.FormItemTemplateTable
	err = x.SQL("select * from form_item_template where form_template in (select form_template from task_template where node_def_id=? and request_template in (select request_template from request where proc_instance_id=?))", nodeDefId, procInstId).Find(&items)
	if err != nil {
		return
	}
	if len(items) == 0 {
		err = fmt.Errorf("Can not find task item template with procInstId:%s nodeDefId:%s ", procInstId, nodeDefId)
		return
	}
	resultData := models.TaskMetaResultData{FormMetaId: items[0].FormTemplate}
	for _, item := range items {
		if item.Entity == "" {
			continue
		}
		resultData.FormItemMetas = append(resultData.FormItemMetas, &models.TaskMetaResultItem{FormItemMetaId: item.Id, PackageName: item.PackageName, EntityName: item.Entity, AttrName: item.Name})
	}
	result.Data = resultData
	return
}

func PluginTaskCreate(input *models.PluginTaskCreateRequestObj) (result *models.PluginTaskCreateOutputObj, err error) {
	result = &models.PluginTaskCreateOutputObj{CallbackParameter: input.CallbackParameter, ErrorCode: "0", ErrorMessage: ""}
	var requestTable []*models.RequestTable
	err = x.SQL("select id,form,request_template from request where proc_instance_id=?", input.ProcInstId).Find(&requestTable)
	if err != nil {
		return
	}

	return
}
