package db

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func GetEntityData(requestTemplateId string) (result models.EntityQueryResult, err error) {
	requestTemplateObj, getTemplateErr := getSimpleRequestTemplate(requestTemplateId)
	if getTemplateErr != nil {
		err = getTemplateErr
		return
	}
	if requestTemplateObj.PackageName == "" || requestTemplateObj.EntityName == "" {
		err = fmt.Errorf("RequestTemplate packageName or entityName illegal ")
		return
	}
	req, newReqErr := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s/entities/%s/query", models.Config.Wecube.BaseUrl, requestTemplateObj.PackageName, requestTemplateObj.EntityName), strings.NewReader("{\"criteria\":{}}"))
	if newReqErr != nil {
		err = fmt.Errorf("Try to new http request fail,%s ", newReqErr.Error())
		return
	}
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = fmt.Errorf("Try to do http request fail,%s ", respErr.Error())
		return
	}
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	err = json.Unmarshal(b, &result)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
	}
	return
}

func ProcessDataPreview(requestTemplateId, entityDataId string) (result models.EntityTreeResult, err error) {
	requestTemplateObj, getTemplateErr := getSimpleRequestTemplate(requestTemplateId)
	if getTemplateErr != nil {
		err = getTemplateErr
		return
	}
	if requestTemplateObj.ProcDefId == "" {
		err = fmt.Errorf("RequestTemplate proDefId illegal ")
		return
	}
	req, newReqErr := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/platform/v1/public/process/definitions/%s/preview/entities/%s", requestTemplateObj.ProcDefId, entityDataId), nil)
	if newReqErr != nil {
		err = fmt.Errorf("Try to new http request fail,%s ", newReqErr.Error())
		return
	}
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = fmt.Errorf("Try to do http request fail,%s ", respErr.Error())
		return
	}
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	err = json.Unmarshal(b, &result)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
	}
	return
}

func ListUserRequest(user string) (result []*models.RequestTable, err error) {
	result = []*models.RequestTable{}
	err = x.SQL("select id,name,form,request_template,proc_instance_id,reporter,report_time,emergency,status from request where reporter=?", user).Find(&result)
	return
}

func GetRequest(requestId string) (result models.RequestTable, err error) {
	result = models.RequestTable{}
	var requestTable []*models.RequestTable
	err = x.SQL("select * from request where id=?", requestId).Find(&requestTable)
	if err != nil {
		return
	}
	if len(requestTable) == 0 {
		err = fmt.Errorf("Can not find any request with id:%s ", requestId)
		return
	}
	result = *requestTable[0]
	return
}

func CreateRequest(param *models.RequestTable, operatorRoles []string) error {
	requestTemplateObj, err := getSimpleRequestTemplate(param.RequestTemplate)
	if err != nil {
		return err
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	formGuid := guid.CreateGuid()
	param.Id = guid.CreateGuid()
	var actions []*execAction
	formInsertAction := execAction{Sql: "insert into form(id,name,description,form_template,created_time,created_by,updated_time,updated_by) value (?,?,?,?,?,?,?,?)"}
	formInsertAction.Param = []interface{}{formGuid, param.Name + models.SysTableIdConnector + "form", "", requestTemplateObj.FormTemplate, nowTime, param.CreatedBy, nowTime, param.CreatedBy}
	actions = append(actions, &formInsertAction)
	requestInsertAction := execAction{Sql: "insert into request(id,name,form,request_template,reporter,emergency,report_role,status,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?,?,?,?,?,?)"}
	requestInsertAction.Param = []interface{}{param.Id, param.Name, formGuid, param.RequestTemplate, param.CreatedBy, param.Emergency, strings.Join(operatorRoles, ","), "created", param.CreatedBy, nowTime, param.CreatedBy, nowTime}
	actions = append(actions, &requestInsertAction)
	return transactionWithoutForeignCheck(actions)
}

func UpdateRequest(param *models.RequestTable) error {
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err := x.Exec("update request set name=?,emergency=?,updated_by=?,updated_time=? where id=?", param.Name, param.Emergency, param.UpdatedBy, nowTime)
	return err
}

func SaveRequestCache(requestId string) {

}

func GetRequestRootForm(requestId string) (result models.RequestTemplateFormStruct, err error) {
	result = models.RequestTemplateFormStruct{}
	var requestTable []*models.RequestTable
	err = x.SQL("select request_template from request where id=?", requestId).Find(&requestTable)
	if err != nil {
		return
	}
	if len(requestTable) == 0 {
		err = fmt.Errorf("Can not find request with id:%s ", requestId)
		return
	}
	requestTemplateObj, _ := getSimpleRequestTemplate(requestTable[0].RequestTemplate)
	result.Id = requestTemplateObj.Id
	result.Name = requestTemplateObj.Name
	result.PackageName = requestTemplateObj.PackageName
	result.EntityName = requestTemplateObj.EntityName
	result.ProcDefName = requestTemplateObj.ProcDefName
	result.ProcDefId = requestTemplateObj.ProcDefId
	result.ProcDefKey = requestTemplateObj.ProcDefKey
	var items []*models.FormItemTemplateTable
	x.SQL("select * from form_item_template where form_template=?", requestTemplateObj.FormTemplate).Find(&items)
	result.FormItems = items
	return
}

func GetRequestPreDataStruct(requestId string) {

}
