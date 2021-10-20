package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

func GetEntityData(requestId, userToken string) (result models.EntityQueryResult, err error) {
	requestTemplateId, tmpErr := getRequestTemplateByRequest(requestId)
	if tmpErr != nil {
		return result, tmpErr
	}
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
	req.Header.Set("Authorization", userToken)
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

func ProcessDataPreview(requestTemplateId, entityDataId, userToken string) (result models.EntityTreeResult, err error) {
	requestTemplateObj, getTemplateErr := getSimpleRequestTemplate(requestTemplateId)
	if getTemplateErr != nil {
		err = getTemplateErr
		return
	}
	if requestTemplateObj.ProcDefId == "" {
		err = fmt.Errorf("RequestTemplate proDefId illegal ")
		return
	}
	req, newReqErr := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/platform/v1/public/process/definitions/%s/preview/entities/%s", models.Config.Wecube.BaseUrl, requestTemplateObj.ProcDefId, entityDataId), nil)
	if newReqErr != nil {
		err = fmt.Errorf("Try to new http request fail,%s ", newReqErr.Error())
		return
	}
	req.Header.Set("Authorization", userToken)
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

func ListRequest(param *models.QueryRequestParam, userRoles []string) (pageInfo models.PageInfo, rowData []*models.RequestTable, err error) {
	rowData = []*models.RequestTable{}
	filterSql, _, queryParam := transFiltersToSQL(param, &models.TransFiltersParam{IsStruct: true, StructObj: models.RequestTable{}, PrimaryKey: "id"})
	baseSql := fmt.Sprintf("select id,name,form,request_template,proc_instance_id,proc_instance_key,reporter,report_time,emergency,status,created_by,created_time,updated_by,updated_time from request where del_flag=0 and request_template in (select id from request_template where id in (select request_template from request_template_role where role_type='USE' and `role` in ('"+strings.Join(userRoles, "','")+"'))) %s ", filterSql)
	if param.Paging {
		pageInfo.StartIndex = param.Pageable.StartIndex
		pageInfo.PageSize = param.Pageable.PageSize
		pageInfo.TotalRows = queryCount(baseSql, queryParam...)
		pageSql, pageParam := transPageInfoToSQL(*param.Pageable)
		baseSql += pageSql
		queryParam = append(queryParam, pageParam...)
	}
	err = x.SQL(baseSql, queryParam...).Find(&rowData)
	if len(rowData) > 0 {
		var requestTemplateTable []*models.RequestTemplateTable
		x.SQL("select id,name from request_template").Find(&requestTemplateTable)
		rtMap := make(map[string]string)
		for _, v := range requestTemplateTable {
			rtMap[v.Id] = v.Name
		}
		for _, v := range rowData {
			v.RequestTemplate = rtMap[v.RequestTemplate]
		}
	}
	return
}

func GetRequestWithRoot(requestId string) (result models.RequestTable, err error) {
	result = models.RequestTable{}
	var requestTable []*models.RequestTable
	err = x.SQL("select id,name,form,request_template,proc_instance_id,proc_instance_key,reporter,report_time,emergency,status,cache from request where id=?", requestId).Find(&requestTable)
	if err != nil {
		return
	}
	if len(requestTable) == 0 {
		err = fmt.Errorf("Can not find any request with id:%s ", requestId)
		return
	}
	result = *requestTable[0]
	if result.Cache != "" {
		var cacheObj models.RequestPreDataDto
		err = json.Unmarshal([]byte(result.Cache), &cacheObj)
		if err != nil {
			err = fmt.Errorf("Try to json unmarshal cache data fail,%s ", err.Error())
			return
		}
		result.Cache = cacheObj.RootEntityId
	}
	return
}

func GetRequest(requestId string) (result models.RequestTable, err error) {
	result = models.RequestTable{}
	var requestTable []*models.RequestTable
	err = x.SQL("select id,name,form,request_template,proc_instance_id,proc_instance_key,reporter,report_time,emergency,status from request where id=?", requestId).Find(&requestTable)
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
	_, err := x.Exec("update request set name=?,emergency=?,updated_by=?,updated_time=? where id=?", param.Name, param.Emergency, param.UpdatedBy, nowTime, param.Id)
	return err
}

func DeleteRequest(requestId, operator string) error {
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err := x.Exec("update request set del_flag=1,updated_by=?,updated_time=? where id=?", operator, nowTime, requestId)
	return err
}

func SaveRequestCacheNew(requestId, operator string, param *models.RequestPreDataDto) error {
	for _, v := range param.Data {
		for _, value := range v.Value {
			if value.Id == "" {
				value.PackageName = v.PackageName
				value.EntityName = v.Entity
				value.EntityDataOp = "create"
				value.Id = fmt.Sprintf("tmp%s%s", models.SysTableIdConnector, guid.CreateGuid())
				value.DataId = value.Id
			}
		}
	}
	paramBytes, err := json.Marshal(param)
	if err != nil {
		return fmt.Errorf("Try to json marshal param fail,%s ", err.Error())
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	actions := UpdateRequestFormItem(requestId, param)
	actions = append(actions, &execAction{Sql: "update request set cache=?,updated_by=?,updated_time=? where id=?", Param: []interface{}{string(paramBytes), operator, nowTime, requestId}})
	return transaction(actions)
}

func UpdateRequestFormItem(requestId string, param *models.RequestPreDataDto) []*execAction {
	var actions []*execAction
	requestObj, _ := GetRequest(requestId)
	actions = append(actions, &execAction{Sql: "delete from form_item where form in (select form from request where id=?)", Param: []interface{}{requestId}})
	for _, v := range param.Data {
		for _, valueObj := range v.Value {
			tmpGuidList := guid.CreateGuidList(len(v.Title))
			for i, title := range v.Title {
				actions = append(actions, &execAction{Sql: "insert into form_item(id,form,form_item_template,name,value,item_group,row_data_id) value (?,?,?,?,?,?,?)", Param: []interface{}{tmpGuidList[i], requestObj.Form, title.Id, title.Name, valueObj.EntityData[title.Name], title.ItemGroup, valueObj.DataId}})
			}
		}
	}
	return actions
}

func GetRequestCache(requestId string) (result models.RequestPreDataDto, err error) {
	var requestTable []*models.RequestTable
	err = x.SQL("select cache from request where id=?", requestId).Find(&requestTable)
	if err != nil {
		return
	}
	if len(requestTable) == 0 {
		err = fmt.Errorf("Can not find any request with id:%s ", requestId)
		return
	}
	err = json.Unmarshal([]byte(requestTable[0].Cache), &result)
	return
}

func getRequestTemplateByRequest(requestId string) (templateId string, err error) {
	var requestTable []*models.RequestTable
	err = x.SQL("select request_template from request where id=?", requestId).Find(&requestTable)
	if err != nil {
		return
	}
	if len(requestTable) == 0 {
		err = fmt.Errorf("Can not find request with id:%s ", requestId)
		return
	}
	templateId = requestTable[0].RequestTemplate
	return
}

func GetRequestRootForm(requestId string) (result models.RequestTemplateFormStruct, err error) {
	result = models.RequestTemplateFormStruct{}
	requestTemplateId, tmpErr := getRequestTemplateByRequest(requestId)
	if tmpErr != nil {
		return result, tmpErr
	}
	requestTemplateObj, _ := getSimpleRequestTemplate(requestTemplateId)
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

func GetRequestPreData(requestId, entityDataId, userToken string) (result []*models.RequestPreDataTableObj, err error) {
	var requestTables []*models.RequestTable
	err = x.SQL("select cache from request where id=?", requestId).Find(&requestTables)
	if err != nil {
		return
	}
	if len(requestTables) == 0 {
		return result, fmt.Errorf("Can not find requestId:%s ", requestId)
	}
	if requestTables[0].Cache != "" {
		var cacheObj models.RequestPreDataDto
		err = json.Unmarshal([]byte(requestTables[0].Cache), &cacheObj)
		if err != nil {
			return result, fmt.Errorf("Try to json unmarshal cache data fail,%s ", err.Error())
		}
		if cacheObj.RootEntityId == entityDataId {
			result = cacheObj.Data
			return
		}
	}
	result = []*models.RequestPreDataTableObj{}
	requestTemplateId, tmpErr := getRequestTemplateByRequest(requestId)
	if tmpErr != nil {
		return result, tmpErr
	}
	var items []*models.FormItemTemplateTable
	err = x.SQL("select * from form_item_template where form_template in (select form_template from request_template where id=?) order by item_group,sort", requestTemplateId).Find(&items)
	if err != nil {
		return
	}
	if len(items) == 0 {
		return result, fmt.Errorf("RequestTemplate:%s have no task form items ", requestTemplateId)
	}
	result = getItemTemplateTitle(items)
	if entityDataId == "" {
		return
	}
	previewData, previewErr := ProcessDataPreview(requestTemplateId, entityDataId, userToken)
	if previewErr != nil {
		return result, previewErr
	}
	if len(previewData.Data.EntityTreeNodes) == 0 {
		return
	}
	for _, entity := range result {
		for _, tmpData := range previewData.Data.EntityTreeNodes {
			if tmpData.EntityName == entity.Entity {
				tmpValueData := make(map[string]interface{})
				for _, title := range entity.Title {
					tmpValueData[title.Name] = tmpData.EntityData[title.Name]
				}
				entity.Value = append(entity.Value, &models.EntityTreeObj{Id: tmpData.Id, PackageName: tmpData.PackageName, EntityName: tmpData.EntityName, DataId: tmpData.DataId, PreviousIds: tmpData.PreviousIds, SucceedingIds: tmpData.SucceedingIds, DisplayName: tmpData.DisplayName, FullDataId: tmpData.FullDataId, EntityData: tmpValueData})
			}
		}
	}
	return
}

func getItemTemplateTitle(items []*models.FormItemTemplateTable) []*models.RequestPreDataTableObj {
	result := []*models.RequestPreDataTableObj{}
	tmpPackageName := items[0].PackageName
	tmpEntity := items[0].Entity
	tmpItemGroup := items[0].ItemGroup
	tmpItemGroupName := items[0].ItemGroupName
	tmpRefEntity := []string{}
	tmpItems := []*models.FormItemTemplateTable{}
	existItemMap := make(map[string]int)
	for _, v := range items {
		tmpKey := fmt.Sprintf("%s__%s", v.ItemGroup, v.Name)
		if _, b := existItemMap[tmpKey]; b {
			continue
		} else {
			existItemMap[tmpKey] = 1
		}
		if v.ItemGroup != tmpItemGroup {
			if tmpItemGroup != "" {
				result = append(result, &models.RequestPreDataTableObj{Entity: tmpEntity, ItemGroup: tmpItemGroup, ItemGroupName: tmpItemGroupName, PackageName: tmpPackageName, Title: tmpItems, RefEntity: tmpRefEntity, Value: []*models.EntityTreeObj{}})
			}
			tmpItems = []*models.FormItemTemplateTable{}
			tmpEntity = v.Entity
			tmpPackageName = v.PackageName
			tmpItemGroup = v.ItemGroup
			tmpItemGroupName = v.ItemGroupName
			tmpRefEntity = []string{}
		} else {
			if tmpEntity == "" && v.Entity != "" {
				tmpEntity = v.Entity
				tmpPackageName = v.PackageName
			}
		}
		tmpItems = append(tmpItems, v)
		if v.RefEntity != "" {
			existFlag := false
			for _, vv := range tmpRefEntity {
				if vv == v.RefEntity {
					existFlag = true
					break
				}
			}
			if !existFlag {
				tmpRefEntity = append(tmpRefEntity, v.RefEntity)
			}
		}
	}
	if len(tmpItems) > 0 {
		tmpEntity = items[len(items)-1].Entity
		tmpPackageName = items[len(items)-1].PackageName
		tmpItemGroup = items[len(items)-1].ItemGroup
		tmpItemGroupName = items[len(items)-1].ItemGroupName
		result = append(result, &models.RequestPreDataTableObj{Entity: tmpEntity, ItemGroup: tmpItemGroup, ItemGroupName: tmpItemGroupName, PackageName: tmpPackageName, Title: tmpItems, RefEntity: tmpRefEntity, Value: []*models.EntityTreeObj{}})
	}
	// sort result by dependence
	result = sortRequestEntity(result)
	return result
}

func sortRequestEntity(param []*models.RequestPreDataTableObj) models.RequestPreDataSort {
	var result models.RequestPreDataSort
	entityMap := make(map[string][]string)
	entityNumMap := make(map[string]int)
	for _, v := range param {
		for _, vv := range v.RefEntity {
			if vv == v.Entity {
				continue
			}
			if _, b := entityMap[vv]; b {
				entityMap[vv] = append(entityMap[vv], v.Entity)
			} else {
				entityMap[vv] = []string{v.Entity}
			}
		}
	}
	countNum := 0
	levelNum := 1
	for _, v := range param {
		if v.Entity == "" {
			v.SortLevel = levelNum
			entityNumMap[v.Entity] = v.SortLevel
			countNum = countNum + 1
			continue
		}
		if _, b := entityMap[v.Entity]; !b {
			v.SortLevel = levelNum
			entityNumMap[v.Entity] = v.SortLevel
			countNum = countNum + 1
		}
	}
	for countNum < len(param) {
		levelNum = levelNum + 1
		for k, v := range entityMap {
			ready := true
			for _, vv := range v {
				if _, b := entityNumMap[vv]; !b {
					ready = false
					break
				} else {
					if entityNumMap[vv] == levelNum {
						ready = false
						break
					}
				}
			}
			if ready {
				for _, vv := range param {
					if vv.Entity == k {
						vv.SortLevel = levelNum
						entityNumMap[vv.Entity] = vv.SortLevel
						countNum = countNum + 1
					}
				}
				delete(entityMap, k)
			}
		}
	}
	for _, v := range param {
		result = append(result, v)
	}
	sort.Sort(result)
	return result
}

func StartRequest(requestId, operator, userToken string, cacheData models.RequestCacheData) (result models.StartInstanceResultData, err error) {
	fillBindingWithRequestData(requestId, &cacheData)
	cacheBytes, tmpErr := json.Marshal(cacheData)
	if tmpErr != nil {
		err = fmt.Errorf("Json marshal cache data fail,%s ", tmpErr.Error())
		return
	}
	req, newReqErr := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/platform/v1/public/process/instances", models.Config.Wecube.BaseUrl), bytes.NewReader(cacheBytes))
	if newReqErr != nil {
		err = fmt.Errorf("Try to new http request fail,%s ", newReqErr.Error())
		return
	}
	req.Header.Set("Authorization", userToken)
	req.Header.Set("Content-Type", "application/json")
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = fmt.Errorf("Try to do http request fail,%s ", respErr.Error())
		return
	}
	var respResult models.StartInstanceResult
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	err = json.Unmarshal(b, &respResult)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	if respResult.Status != "OK" {
		err = fmt.Errorf("Start instance fail,%s ", respResult.Message)
		return
	}
	result = respResult.Data
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err = x.Exec("update request set proc_instance_id=?,proc_instance_key=?,report_time=?,status=?,bind_cache=?,updated_by=?,updated_time=? where id=?", strconv.Itoa(result.Id), result.ProcInstKey, nowTime, respResult.Data.Status, string(cacheBytes), operator, nowTime, requestId)
	return
}

func fillBindingWithRequestData(requestId string, cacheData *models.RequestCacheData) {
	var items []*models.FormItemTemplateTable
	x.SQL("select * from form_item_template where form_template in (select form_template from request_template where id=?) order by entity,sort", requestId).Find(&items)
	itemMap := make(map[string][]string)
	for _, item := range items {
		if item.Entity == "" || item.RefEntity == "" {
			continue
		}
		if _, b := itemMap[item.Entity]; !b {
			itemMap[item.Entity] = []string{item.Name}
		} else {
			itemMap[item.Entity] = append(itemMap[item.Entity], item.Name)
		}
	}
	entityNewMap := make(map[string][]string)
	for _, taskNode := range cacheData.TaskNodeBindInfos {
		for _, entityValue := range taskNode.BoundEntityValues {
			if entityValue.EntityDataId == cacheData.RootEntityValue.Oid {
				cacheData.RootEntityValue.Oid = entityValue.Oid
				cacheData.RootEntityValue.EntityName = entityValue.EntityName
				cacheData.RootEntityValue.EntityDataOp = entityValue.EntityDataOp
				cacheData.RootEntityValue.AttrValues = entityValue.AttrValues
				cacheData.RootEntityValue.PackageName = entityValue.PackageName
				cacheData.RootEntityValue.BindFlag = entityValue.BindFlag
				cacheData.RootEntityValue.EntityDataId = entityValue.EntityDataId
				cacheData.RootEntityValue.EntityDataState = entityValue.EntityDataState
				cacheData.RootEntityValue.EntityDefId = entityValue.EntityDefId
				cacheData.RootEntityValue.FullEntityDataId = entityValue.FullEntityDataId
				cacheData.RootEntityValue.PreviousOids = entityValue.PreviousOids
				cacheData.RootEntityValue.SucceedingOids = entityValue.SucceedingOids
			}
			if _, b := entityNewMap[entityValue.Oid]; b {
				continue
			}
			if entityRefs, b := itemMap[entityValue.EntityName]; b {
				if entityValue.EntityDataOp == "create" {
					tmpRefOidList := []string{}
					for _, attrValueObj := range entityValue.AttrValues {
						for _, entityRef := range entityRefs {
							if attrValueObj.AttrName == entityRef {
								tmpRefOidList = append(tmpRefOidList, fmt.Sprintf("%s", attrValueObj.DataValue))
							}
						}
					}
					entityNewMap[entityValue.Oid] = tmpRefOidList
				} else {
					tmpRefOidList := []string{}
					for _, attrValueObj := range entityValue.AttrValues {
						for _, entityRef := range entityRefs {
							if attrValueObj.AttrName == entityRef {
								valueString := fmt.Sprintf("%s", attrValueObj.DataValue)
								if strings.HasPrefix(valueString, "tmp") {
									tmpRefOidList = append(tmpRefOidList, valueString)
								}
							}
						}
					}
					entityNewMap[entityValue.Oid] = tmpRefOidList
				}
			}
		}
	}
	if len(entityNewMap) > 0 {
		for _, taskNode := range cacheData.TaskNodeBindInfos {
			for _, entityValue := range taskNode.BoundEntityValues {
				if refOids, b := entityNewMap[entityValue.Oid]; b {
					entityValue.PreviousOids = append(entityValue.PreviousOids, refOids...)
				}
				for tmpOid, refOids := range entityNewMap {
					inRefFlag := false
					for _, refOid := range refOids {
						if entityValue.Oid == refOid {
							inRefFlag = true
						}
					}
					if inRefFlag {
						entityValue.SucceedingOids = append(entityValue.SucceedingOids, tmpOid)
					}
				}
				entityValue.PreviousOids = listToSet(entityValue.PreviousOids)
				entityValue.SucceedingOids = listToSet(entityValue.SucceedingOids)
			}
		}
	}
}

func listToSet(input []string) []string {
	result := []string{}
	tmpMap := make(map[string]int)
	for _, v := range input {
		if _, b := tmpMap[v]; !b {
			result = append(result, v)
			tmpMap[v] = 1
		}
	}
	return result
}

func RequestTermination(requestId, operator, userToken string) error {
	requestObj, err := GetRequest(requestId)
	if err != nil {
		return err
	}
	param := models.TerminateInstanceParam{ProcInstId: requestObj.ProcInstanceId, ProcInstKey: requestObj.ProcInstanceKey}
	paramBytes, tmpErr := json.Marshal(param)
	if tmpErr != nil {
		return fmt.Errorf("Json marshal param data fail,%s ", tmpErr.Error())
	}
	req, newReqErr := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/platform/v1/public/process/instances/%s/terminations", models.Config.Wecube.BaseUrl, requestObj.ProcInstanceId), bytes.NewReader(paramBytes))
	if newReqErr != nil {
		return fmt.Errorf("Try to new http request fail,%s ", newReqErr.Error())
	}
	req.Header.Set("Authorization", userToken)
	req.Header.Set("Content-Type", "application/json")
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		return fmt.Errorf("Try to do http request fail,%s ", respErr.Error())
	}
	var respResult models.StartInstanceResult
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	err = json.Unmarshal(b, &respResult)
	if err != nil {
		return fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
	}
	if respResult.Status != "OK" {
		return fmt.Errorf("Terminate instance fail,%s ", respResult.Message)
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err = x.Exec("update request set status='Termination',updated_by=?,updated_time=? where id=?", operator, nowTime, requestId)
	return err
}
