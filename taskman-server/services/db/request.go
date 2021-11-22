package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
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
	req, newReqErr := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/platform/v1/process/definitions/%s/root-entities", models.Config.Wecube.BaseUrl, requestTemplateObj.ProcDefId), strings.NewReader(""))
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

func ListRequest(param *models.QueryRequestParam, userRoles []string, userToken, permission, operator string) (pageInfo models.PageInfo, rowData []*models.RequestTable, err error) {
	rowData = []*models.RequestTable{}
	if strings.ToLower(permission) == "mgmt" {
		permission = "MGMT"
	} else {
		permission = "USE"
	}
	filterSql, _, queryParam := transFiltersToSQL(param, &models.TransFiltersParam{IsStruct: true, StructObj: models.RequestTable{}, PrimaryKey: "id"})
	baseSql := fmt.Sprintf("select id,name,form,request_template,proc_instance_id,proc_instance_key,reporter,handler,report_time,emergency,status,expire_time,expect_time,confirm_time,created_by,created_time,updated_by,updated_time from request where del_flag=0 and (created_by='"+operator+"' or request_template in (select id from request_template where id in (select request_template from request_template_role where role_type='"+permission+"' and `role` in ('"+strings.Join(userRoles, "','")+"')))) %s ", filterSql)
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
		x.SQL("select id,name,status,version from request_template").Find(&requestTemplateTable)
		rtMap := make(map[string]string)
		for _, v := range requestTemplateTable {
			if v.Status != "confirm" {
				rtMap[v.Id] = fmt.Sprintf("%s(beta)", v.Name)
			} else {
				rtMap[v.Id] = fmt.Sprintf("%s(%s)", v.Name, v.Version)
			}
		}
		rtRoleMap := getRequestTemplateMGMTRole()
		var actions []*execAction
		for _, v := range rowData {
			v.RequestTemplateName = rtMap[v.RequestTemplate]
			if tmpRoles, b := rtRoleMap[v.RequestTemplate]; b {
				v.HandleRoles = tmpRoles
			} else {
				v.HandleRoles = []string{}
			}
			if strings.Contains(v.Status, "InProgress") && v.ProcInstanceId != "" {
				newStatus := getInstanceStatus(v.ProcInstanceId, userToken)
				if newStatus == "InternallyTerminated" {
					newStatus = "Termination"
				}
				if newStatus != "" && newStatus != v.Status {
					actions = append(actions, &execAction{Sql: "update request set status=? where id=?", Param: []interface{}{newStatus, v.Id}})
					v.Status = newStatus
				}
			}
		}
		if len(actions) > 0 {
			updateStatusErr := transaction(actions)
			if updateStatusErr != nil {
				log.Logger.Error("Try to update request status fail", log.Error(updateStatusErr))
			}
		}
	}
	return
}

func getRequestTemplateMGMTRole() (result map[string][]string) {
	result = make(map[string][]string)
	var requestTemplateRole []*models.RequestTemplateRoleTable
	x.SQL("select * from request_template_role where role_type='MGMT' order by request_template").Find(&requestTemplateRole)
	if len(requestTemplateRole) == 0 {
		return result
	}
	tmpTemplate := requestTemplateRole[0].RequestTemplate
	tmpRoles := []string{}
	for _, v := range requestTemplateRole {
		if v.RequestTemplate != tmpTemplate {
			result[tmpTemplate] = tmpRoles
			tmpTemplate = v.RequestTemplate
			tmpRoles = []string{}
		}
		tmpRoles = append(tmpRoles, v.Role)
	}
	if len(tmpRoles) > 0 {
		tmpTemplate = requestTemplateRole[len(requestTemplateRole)-1].RequestTemplate
		result[tmpTemplate] = tmpRoles
	}
	return result
}

func calcExpireObj(param *models.ExpireObj) {
	if param.ReportTime == "" || param.ExpireTime == "" {
		return
	}
	reportT, _ := time.Parse(models.DateTimeFormat, param.ReportTime)
	expireT, _ := time.Parse(models.DateTimeFormat, param.ExpireTime)
	nowT, _ := time.Parse(models.DateTimeFormat, param.NowTime)
	max := expireT.Sub(reportT).Seconds()
	use := nowT.Sub(reportT).Seconds()
	param.Percent = (use / max) * 100
	param.TotalDay = max / 86400
	param.UseDay = use / 86400
	param.LeftDay = (max - use) / 86400
	param.Percent, _ = strconv.ParseFloat(fmt.Sprintf("%.0f", param.Percent), 64)
	param.TotalDay, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", param.TotalDay), 64)
	param.LeftDay, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", param.LeftDay), 64)
	param.UseDay, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", param.UseDay), 64)
	return
}

func getInstanceStatus(instanceId, userToken string) string {
	req, newReqErr := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/platform/v1/process/instances/%s", models.Config.Wecube.BaseUrl, instanceId), nil)
	if newReqErr != nil {
		log.Logger.Error("GetInstanceStatus fail", log.String("msg", "new http request fail"), log.Error(newReqErr))
		return ""
	}
	req.Header.Set("Authorization", userToken)
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		log.Logger.Error("GetInstanceStatus fail", log.String("msg", "Try to do http request fail"), log.Error(respErr))
		return ""
	}
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	var response models.InstanceStatusQuery
	err := json.Unmarshal(b, &response)
	if err != nil {
		log.Logger.Error("GetInstanceStatus fail", log.String("msg", "Try to json unmarshal body fail"), log.Error(err))
		return ""
	}
	if response.Status != "OK" {
		log.Logger.Error("GetInstanceStatus fail", log.String("msg", response.Message))
		return ""
	}
	if response.Data.Status != "InProgress" {
		return response.Data.Status
	}
	status := "InProgress"
	for _, v := range response.Data.TaskNodeInstances {
		if v.Status == "Faulted" {
			status = "InProgress(Faulted)"
			break
		}
		if v.Status == "Timeouted" {
			status = "InProgress(Timeouted)"
			break
		}
	}
	return status
}

func calcExpireTime(reportTime string, expireDay int) (expire string) {
	t, err := time.Parse(models.DateTimeFormat, reportTime)
	if err != nil {
		return
	}
	expire = t.Add(time.Duration(expireDay*24) * time.Hour).Format(models.DateTimeFormat)
	return
}

func GetRequestWithRoot(requestId string) (result models.RequestTable, err error) {
	result = models.RequestTable{}
	var requestTable []*models.RequestTable
	err = x.SQL("select id,name,form,request_template,proc_instance_id,proc_instance_key,reporter,report_time,emergency,status,cache,expire_time,expect_time,confirm_time,handler from request where id=?", requestId).Find(&requestTable)
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
	result.AttachFiles = GetRequestAttachFileList(requestId)
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

func CreateRequest(param *models.RequestTable, operatorRoles []string, userToken string) error {
	requestTemplateObj, err := getSimpleRequestTemplate(param.RequestTemplate)
	if err != nil {
		return err
	}
	var actions []*execAction
	err = SyncProcDefId(requestTemplateObj.Id, requestTemplateObj.ProcDefId, requestTemplateObj.ProcDefName, "", userToken)
	if err != nil {
		return fmt.Errorf("Try to sync proDefId fail,%s ", err.Error())
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	formGuid := guid.CreateGuid()
	param.Id = guid.CreateGuid()
	formInsertAction := execAction{Sql: "insert into form(id,name,description,form_template,created_time,created_by,updated_time,updated_by) value (?,?,?,?,?,?,?,?)"}
	formInsertAction.Param = []interface{}{formGuid, param.Name + models.SysTableIdConnector + "form", "", requestTemplateObj.FormTemplate, nowTime, param.CreatedBy, nowTime, param.CreatedBy}
	actions = append(actions, &formInsertAction)
	requestInsertAction := execAction{Sql: "insert into request(id,name,form,request_template,reporter,emergency,report_role,status,expire_time,expect_time,handler,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
	requestInsertAction.Param = []interface{}{param.Id, param.Name, formGuid, param.RequestTemplate, param.CreatedBy, param.Emergency, strings.Join(operatorRoles, ","), "Draft", "", param.ExpectTime, requestTemplateObj.Handler, param.CreatedBy, nowTime, param.CreatedBy, nowTime}
	actions = append(actions, &requestInsertAction)
	return transactionWithoutForeignCheck(actions)
}

func UpdateRequest(param *models.RequestTable) error {
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err := x.Exec("update request set name=?,expect_time=?,emergency=?,handler=?,updated_by=?,updated_time=? where id=?", param.Name, param.ExpectTime, param.Emergency, param.Handler, param.UpdatedBy, nowTime, param.Id)
	return err
}

func DeleteRequest(requestId, operator string) error {
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err := x.Exec("update request set del_flag=1,updated_by=?,updated_time=? where id=?", operator, nowTime, requestId)
	return err
}

func SaveRequestCacheNew(requestId, operator string, param *models.RequestPreDataDto) error {
	var formItemNameQuery []*models.FormItemTemplateTable
	err := x.SQL("select item_group,group_concat(name,',') as name from form_item_template where in_display_name='yes' and form_template in (select form_template from request_template where id in (select request_template from request where id=?)) group by item_group", requestId).Find(&formItemNameQuery)
	itemGroupNameMap := make(map[string][]string)
	for _, v := range formItemNameQuery {
		itemGroupNameMap[v.ItemGroup] = strings.Split(v.Name, ",")
	}
	for _, v := range param.Data {
		nameList := itemGroupNameMap[v.ItemGroup]
		for _, value := range v.Value {
			if value.Id == "" {
				value.PackageName = v.PackageName
				value.EntityName = v.Entity
				value.EntityDataOp = "create"
				value.Id = fmt.Sprintf("tmp%s%s", models.SysTableIdConnector, guid.CreateGuid())
				value.DisplayName = concatItemDisplayName(value.EntityData, nameList)
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

func SaveRequestBindCache(requestId, operator string, param *models.RequestCacheData) error {
	cacheBytes, _ := json.Marshal(param)
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err := x.Exec("update request set status='Draft',bind_cache=?,updated_by=?,updated_time=? where id=?", string(cacheBytes), operator, nowTime, requestId)
	return err
}

func concatItemDisplayName(rowData map[string]interface{}, nameList []string) string {
	displayNameList := []string{}
	for _, v := range nameList {
		if v == "" {
			continue
		}
		if _, b := rowData[v]; b {
			tmpValue := fmt.Sprintf("%s", rowData[v])
			if rowData[v] == nil {
				tmpValue = ""
			}
			displayNameList = append(displayNameList, tmpValue)
		}
	}
	return strings.Join(displayNameList, "__")
}

func UpdateRequestFormItem(requestId string, param *models.RequestPreDataDto) []*execAction {
	var actions []*execAction
	requestObj, _ := GetRequest(requestId)
	actions = append(actions, &execAction{Sql: "delete from form_item where form in (select form from request where id=?)", Param: []interface{}{requestId}})
	for _, v := range param.Data {
		for _, valueObj := range v.Value {
			tmpGuidList := guid.CreateGuidList(len(v.Title))
			for i, title := range v.Title {
				if title.Multiple == "Y" {
					if tmpV, b := valueObj.EntityData[title.Name]; b {
						tmpStringV := []string{}
						for _, interfaceV := range tmpV.([]interface{}) {
							tmpStringV = append(tmpStringV, fmt.Sprintf("%s", interfaceV))
						}
						actions = append(actions, &execAction{Sql: "insert into form_item(id,form,form_item_template,name,value,item_group,row_data_id) value (?,?,?,?,?,?,?)", Param: []interface{}{tmpGuidList[i], requestObj.Form, title.Id, title.Name, strings.Join(tmpStringV, ","), title.ItemGroup, valueObj.Id}})
					}
				} else {
					actions = append(actions, &execAction{Sql: "insert into form_item(id,form,form_item_template,name,value,item_group,row_data_id) value (?,?,?,?,?,?,?)", Param: []interface{}{tmpGuidList[i], requestObj.Form, title.Id, title.Name, valueObj.EntityData[title.Name], title.ItemGroup, valueObj.Id}})
				}
			}
		}
	}
	return actions
}

func GetRequestCache(requestId, cacheType string) (result interface{}, err error) {
	var requestTable []*models.RequestTable
	if cacheType == "data" {
		err = x.SQL("select cache from request where id=?", requestId).Find(&requestTable)
		if err != nil {
			return
		}
		if len(requestTable) == 0 {
			err = fmt.Errorf("Can not find any request with id:%s ", requestId)
			return
		}
		if requestTable[0].Cache == "" {
			return models.RequestPreDataDto{Data: []*models.RequestPreDataTableObj{}}, nil
		}
		var dataCache models.RequestPreDataDto
		err = json.Unmarshal([]byte(requestTable[0].Cache), &dataCache)
		return dataCache, err
	} else {
		err = x.SQL("select bind_cache from request where id=?", requestId).Find(&requestTable)
		if err != nil {
			return
		}
		if len(requestTable) == 0 {
			err = fmt.Errorf("Can not find any request with id:%s ", requestId)
			return
		}
		if requestTable[0].BindCache == "" {
			return models.RequestCacheData{TaskNodeBindInfos: []*models.RequestCacheTaskNodeBindObj{}}, nil
		}
		var bindCache models.RequestCacheData
		err = json.Unmarshal([]byte(requestTable[0].BindCache), &bindCache)
		return bindCache, err
	}
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
		if items[len(items)-1].Entity != "" {
			tmpEntity = items[len(items)-1].Entity
		}
		if items[len(items)-1].PackageName != "" {
			tmpPackageName = items[len(items)-1].PackageName
		}
		tmpItemGroup = items[len(items)-1].ItemGroup
		tmpItemGroupName = items[len(items)-1].ItemGroupName
		result = append(result, &models.RequestPreDataTableObj{Entity: tmpEntity, ItemGroup: tmpItemGroup, ItemGroupName: tmpItemGroupName, PackageName: tmpPackageName, Title: tmpItems, RefEntity: tmpRefEntity, Value: []*models.EntityTreeObj{}})
	}
	// sort result by dependence
	result = sortRequestEntity(result)
	var err error
	result, err = getCMDBSelectList(result, models.CoreToken.GetCoreToken())
	if err != nil {
		log.Logger.Error("Try to get selectList fail", log.Error(err))
	}
	for _, v := range result {
		for _, vv := range v.Title {
			if vv.SelectList == nil {
				vv.SelectList = []*models.EntityDataObj{}
			}
		}
	}
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
	var requestTemplateTable []*models.RequestTemplateTable
	x.SQL("select * from request_template where id in (select request_template from request where id=?)", requestId).Find(&requestTemplateTable)
	if len(requestTemplateTable) == 0 {
		return result, fmt.Errorf("Can not find requestTemplate with request:%s ", requestId)
	}
	cacheData.ProcDefId = requestTemplateTable[0].ProcDefId
	cacheData.ProcDefKey = requestTemplateTable[0].ProcDefKey
	err = AppendUselessEntity(requestTemplateTable[0].Id, userToken, &cacheData)
	if err != nil {
		return result, fmt.Errorf("Try to append useless entity fail,%s ", err.Error())
	}
	fillBindingWithRequestData(requestId, &cacheData)
	cacheBytes, _ := json.Marshal(cacheData)
	startParam := BuildRequestProcessData(cacheData)
	startParamBytes, tmpErr := json.Marshal(startParam)
	if tmpErr != nil {
		err = fmt.Errorf("Json marshal cache data fail,%s ", tmpErr.Error())
		return
	}
	log.Logger.Info("Start request", log.String("param", string(startParamBytes)))
	req, newReqErr := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/platform/v1/public/process/instances", models.Config.Wecube.BaseUrl), bytes.NewReader(startParamBytes))
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
	expireTime := calcExpireTime(nowTime, requestTemplateTable[0].ExpireDay)
	_, err = x.Exec("update request set handler=?,proc_instance_id=?,proc_instance_key=?,confirm_time=?,expire_time=?,status=?,bind_cache=?,updated_by=?,updated_time=? where id=?", operator, strconv.Itoa(result.Id), result.ProcInstKey, nowTime, expireTime, respResult.Data.Status, string(cacheBytes), operator, nowTime, requestId)
	return
}

func UpdateRequestStatus(requestId, status, operator, userToken string) error {
	var err error
	nowTime := time.Now().Format(models.DateTimeFormat)
	if status == "Pending" {
		bindData, bindErr := GetRequestPreBindData(requestId, userToken)
		if bindErr != nil {
			return fmt.Errorf("Try to build bind data fail,%s ", bindErr.Error())
		}
		bindCacheBytes, _ := json.Marshal(bindData)
		bindCache := string(bindCacheBytes)
		_, err = x.Exec("update request set status=?,reporter=?,report_time=?,bind_cache=?,updated_by=?,updated_time=? where id=?", status, operator, nowTime, bindCache, operator, nowTime, requestId)
		if err == nil {
			err = notifyRoleMail(requestId)
		}
	} else {
		_, err = x.Exec("update request set status=?,updated_by=?,updated_time=? where id=?", status, operator, nowTime, requestId)
	}
	return err
}

func fillBindingWithRequestData(requestId string, cacheData *models.RequestCacheData) {
	var items []*models.FormItemTemplateTable
	x.SQL("select * from form_item_template where form_template in (select form_template from request_template where id in (select request_template from request where id=?)) order by entity,sort", requestId).Find(&items)
	itemMap := make(map[string][]string)
	for _, item := range items {
		if item.Entity == "" {
			continue
		}
		if _, b := itemMap[item.Entity]; !b {
			itemMap[item.Entity] = []string{item.Name}
		} else {
			itemMap[item.Entity] = append(itemMap[item.Entity], item.Name)
		}
	}
	for k, v := range itemMap {
		log.Logger.Info("itemMap", log.String("key", k), log.StringList("value", v))
	}
	entityNewMap := make(map[string][]string)
	entityOidMap := make(map[string]int)
	dataIdOidMap := make(map[string]string)
	for _, taskNode := range cacheData.TaskNodeBindInfos {
		for _, entityValue := range taskNode.BoundEntityValues {
			entityOidMap[entityValue.Oid] = 1
			if _, b := entityNewMap[entityValue.Oid]; b {
				continue
			}
			if entityRefs, b := itemMap[entityValue.EntityName]; b {
				if entityValue.EntityDataOp == "create" {
					tmpRefOidList := []string{}
					for _, attrValueObj := range entityValue.AttrValues {
						for _, entityRef := range entityRefs {
							if attrValueObj.AttrName == entityRef {
								tmpV := fmt.Sprintf("%s", attrValueObj.DataValue)
								if strings.Contains(tmpV, ",") {
									tmpRefOidList = append(tmpRefOidList, strings.Split(tmpV, ",")...)
								} else {
									tmpRefOidList = append(tmpRefOidList, tmpV)
								}
							}
						}
					}
					entityNewMap[entityValue.Oid] = tmpRefOidList
				} else {
					dataIdOidMap[entityValue.EntityDataId] = entityValue.Oid
					tmpRefOidList := []string{}
					for _, attrValueObj := range entityValue.AttrValues {
						for _, entityRef := range entityRefs {
							if attrValueObj.AttrName == entityRef {
								valueString := fmt.Sprintf("%s", attrValueObj.DataValue)
								if strings.Contains(valueString, ",") {
									for _, tmpV := range strings.Split(valueString, ",") {
										if strings.HasPrefix(tmpV, "tmp") {
											tmpRefOidList = append(tmpRefOidList, tmpV)
										}
									}
								} else {
									if strings.HasPrefix(valueString, "tmp") {
										tmpRefOidList = append(tmpRefOidList, valueString)
									}
								}
							}
						}
					}
					entityNewMap[entityValue.Oid] = tmpRefOidList
				}
			}
		}
	}
	for k, v := range entityNewMap {
		tmpRefOidList := []string{}
		for _, vv := range v {
			if oid, b := dataIdOidMap[vv]; b {
				tmpRefOidList = append(tmpRefOidList, oid)
			} else {
				tmpRefOidList = append(tmpRefOidList, vv)
			}
		}
		entityNewMap[k] = tmpRefOidList
		log.Logger.Info("entityNewMap", log.String("key", k), log.StringList("value", tmpRefOidList))
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
				entityValue.PreviousOids = listToSet(entityValue.PreviousOids, entityOidMap)
				entityValue.SucceedingOids = listToSet(entityValue.SucceedingOids, entityOidMap)
			}
		}
	}
	for _, taskNode := range cacheData.TaskNodeBindInfos {
		existFlag := false
		for _, entityValue := range taskNode.BoundEntityValues {
			if entityValue.EntityDataId == cacheData.RootEntityValue.Oid {
				existFlag = true
				cacheData.RootEntityValue.Oid = entityValue.Oid
				cacheData.RootEntityValue.EntityName = entityValue.EntityName
				cacheData.RootEntityValue.EntityDataOp = entityValue.EntityDataOp
				cacheData.RootEntityValue.AttrValues = entityValue.AttrValues
				cacheData.RootEntityValue.PackageName = entityValue.PackageName
				cacheData.RootEntityValue.EntityDisplayName = entityValue.EntityDisplayName
				cacheData.RootEntityValue.BindFlag = entityValue.BindFlag
				cacheData.RootEntityValue.EntityDataId = entityValue.EntityDataId
				cacheData.RootEntityValue.EntityDataState = entityValue.EntityDataState
				cacheData.RootEntityValue.EntityDefId = entityValue.EntityDefId
				cacheData.RootEntityValue.FullEntityDataId = entityValue.FullEntityDataId
				cacheData.RootEntityValue.PreviousOids = entityValue.PreviousOids
				cacheData.RootEntityValue.SucceedingOids = entityValue.SucceedingOids
				break
			}
		}
		if existFlag {
			break
		}
	}
}

func listToSet(input []string, itemMap map[string]int) []string {
	result := []string{}
	tmpMap := make(map[string]int)
	for _, v := range input {
		if v == "" {
			continue
		}
		if _, b := tmpMap[v]; !b {
			if _, bb := itemMap[v]; bb {
				result = append(result, v)
				tmpMap[v] = 1
			}
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

func GetCmdbReferenceData(attrId, userToken string, param models.QueryRequestParam) (result []byte, statusCode int, err error) {
	paramBytes, tmpErr := json.Marshal(param)
	if tmpErr != nil {
		err = fmt.Errorf("Json marshal param data fail,%s ", tmpErr.Error())
		return
	}
	req, newReqErr := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/wecmdb/api/v1/ci-data/reference-data/query/%s", models.Config.Wecube.BaseUrl, attrId), bytes.NewReader(paramBytes))
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
	statusCode = resp.StatusCode
	result, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return
}

func GetRequestPreBindData(requestId, userToken string) (result models.RequestCacheData, err error) {
	var requestTable []*models.RequestTable
	err = x.SQL("select * from request where id=?", requestId).Find(&requestTable)
	if err != nil {
		return result, fmt.Errorf("Try to query request fail,%s ", err.Error())
	}
	if len(requestTable) == 0 {
		return result, fmt.Errorf("Can not find request with id:%s ", requestId)
	}
	if requestTable[0].Cache == "" {
		return result, fmt.Errorf("Can not find request cache data with id:%s ", requestId)
	}
	processNodes, processErr := GetProcessNodesByProc(models.RequestTemplateTable{Id: requestTable[0].RequestTemplate}, userToken, "bind")
	if processErr != nil {
		return result, processErr
	}
	entityDefIdMap := make(map[string]string)
	entityBindMap := make(map[string][]string)
	for _, v := range processNodes {
		tmpBoundEntities := []string{}
		for _, vv := range v.BoundEntities {
			entityDefIdMap[fmt.Sprintf("%s:%s", vv.PackageName, vv.Name)] = vv.Id
			tmpBoundEntities = append(tmpBoundEntities, fmt.Sprintf("%s:%s", vv.PackageName, vv.Name))
		}
		if v.TaskCategory != "SUTN" {
			entityBindMap[v.NodeDefId] = tmpBoundEntities
		}
	}
	var dataCache models.RequestPreDataDto
	json.Unmarshal([]byte(requestTable[0].Cache), &dataCache)
	entityOidMap := make(map[string][]string)
	entityValueMap := make(map[string]*models.RequestCacheEntityValue)
	for i, v := range dataCache.Data {
		for ii, vv := range v.Value {
			if vv.EntityName == "" {
				continue
			}
			tmpValueObj := models.RequestCacheEntityValue{Oid: vv.Id, BindFlag: "Y", EntityName: vv.EntityName, EntityDisplayName: vv.DisplayName, EntityDataOp: vv.EntityDataOp, EntityDataId: vv.DataId, FullEntityDataId: vv.FullDataId, PackageName: vv.PackageName, PreviousOids: vv.PreviousIds, SucceedingOids: vv.SucceedingIds}
			tmpEntityName := fmt.Sprintf("%s:%s", vv.PackageName, vv.EntityName)
			if _, b := entityDefIdMap[tmpEntityName]; b {
				tmpValueObj.EntityDefId = entityDefIdMap[tmpEntityName]
			}
			tmpValueObj.AttrValues = buildEntityValueAttrData(v.Title, vv.EntityData)
			if dataCache.RootEntityId == "" && i == 0 && ii == 0 {
				result.RootEntityValue = tmpValueObj
			}
			entityValueMap[vv.Id] = &tmpValueObj
			if _, b := entityOidMap[v.ItemGroup]; b {
				entityOidMap[v.ItemGroup] = append(entityOidMap[v.ItemGroup], vv.Id)
			} else {
				entityOidMap[v.ItemGroup] = []string{vv.Id}
			}
		}
	}
	var entityNodeBind []*models.EntityNodeBindQueryObj
	x.SQL("select distinct t1.node_def_id,t2.item_group from task_template t1 left join form_item_template t2 on t1.form_template=t2.form_template where t1.request_template=?", requestTable[0].RequestTemplate).Find(&entityNodeBind)
	for _, v := range entityNodeBind {
		if _, b := entityBindMap[v.NodeDefId]; b {
			entityBindMap[v.NodeDefId] = append(entityBindMap[v.NodeDefId], v.ItemGroup)
		} else {
			entityBindMap[v.NodeDefId] = []string{v.ItemGroup}
		}
	}
	for _, v := range processNodes {
		if v.TaskCategory == "" {
			continue
		}
		tmpNodeBindInfo := models.RequestCacheTaskNodeBindObj{NodeId: v.NodeId, NodeDefId: v.NodeDefId, BoundEntityValues: []*models.RequestCacheEntityValue{}}
		if entities, b := entityBindMap[v.NodeDefId]; b {
			for _, entity := range entities {
				if oids, entityExist := entityOidMap[entity]; entityExist {
					for _, oid := range oids {
						if _, oidExist := entityValueMap[oid]; oidExist {
							tmpNodeBindInfo.BoundEntityValues = append(tmpNodeBindInfo.BoundEntityValues, entityValueMap[oid])
						}
					}
				}
			}
		}
		result.TaskNodeBindInfos = append(result.TaskNodeBindInfos, &tmpNodeBindInfo)
	}
	return
}

func buildEntityValueAttrData(titles []*models.FormItemTemplateTable, entityData map[string]interface{}) (result []*models.RequestCacheEntityAttrValue) {
	result = []*models.RequestCacheEntityAttrValue{}
	titleMap := make(map[string]*models.FormItemTemplateTable)
	for _, v := range titles {
		titleMap[v.Name] = v
	}
	for k, v := range entityData {
		if vv, b := titleMap[k]; b {
			if vv.Multiple == "Y" {
				tmpV := []string{}
				for _, interfaceV := range v.([]interface{}) {
					tmpV = append(tmpV, fmt.Sprintf("%s", interfaceV))
				}
				result = append(result, &models.RequestCacheEntityAttrValue{AttrDefId: vv.AttrDefId, AttrName: k, DataType: vv.AttrDefDataType, DataValue: strings.Join(tmpV, ",")})
			} else {
				result = append(result, &models.RequestCacheEntityAttrValue{AttrDefId: vv.AttrDefId, AttrName: k, DataType: vv.AttrDefDataType, DataValue: v})
			}
		}
	}
	return
}

func RecordRequestLog(requestId, operator, operation string) {
	_, err := x.Exec("insert into operation_log(id,request,operation,operator,op_time) value (?,?,?,?,?)", guid.CreateGuid(), requestId, operation, operator, time.Now().Format(models.DateTimeFormat))
	if err != nil {
		log.Logger.Error("Record request operation log fail", log.Error(err))
	}
}

func GetRequestTaskList(requestId string) (result models.TaskQueryResult, err error) {
	var taskTable []*models.TaskTable
	err = x.SQL("select id from task where request=? order by created_time desc", requestId).Find(&taskTable)
	if err != nil {
		return
	}
	if len(taskTable) > 0 {
		result, err = GetTask(taskTable[0].Id)
		return
	}
	// get request
	var requests []*models.RequestTable
	x.SQL("select * from request where id=?", requestId).Find(&requests)
	if len(requests) == 0 {
		return result, fmt.Errorf("Can not find request with id:%s ", requestId)
	}
	var requestCache models.RequestPreDataDto
	err = json.Unmarshal([]byte(requests[0].Cache), &requestCache)
	if err != nil {
		return result, fmt.Errorf("Try to json unmarshal request cache fail,%s ", err.Error())
	}
	requestQuery := models.TaskQueryObj{RequestId: requestId, RequestName: requests[0].Name, Reporter: requests[0].Reporter, ReportTime: requests[0].ReportTime, Comment: requests[0].Result, Editable: false}
	requestQuery.FormData = requestCache.Data
	requestQuery.AttachFiles = GetRequestAttachFileList(requestId)
	result.Data = []*models.TaskQueryObj{&requestQuery}
	result.TimeStep, err = getRequestTimeStep(requests[0].RequestTemplate)
	if len(result.TimeStep) > 0 {
		result.TimeStep[0].Active = true
	}
	return
}

func getCMDBSelectList(input []*models.RequestPreDataTableObj, userToken string) (output []*models.RequestPreDataTableObj, err error) {
	ciAttrMap := make(map[string][]string)
	ciAttrSelectMap := make(map[string][]*models.EntityDataObj)
	for _, v := range input {
		if v.Entity == "" {
			continue
		}
		for _, vv := range v.Title {
			if vv.ElementType == "select" && vv.RefEntity == "" {
				if _, b := ciAttrMap[v.Entity]; b {
					ciAttrMap[v.Entity] = append(ciAttrMap[v.Entity], vv.Name)
				} else {
					ciAttrMap[v.Entity] = []string{vv.Name}
				}
			}
			vv.SelectList = []*models.EntityDataObj{}
		}
	}
	output = input
	if len(ciAttrMap) <= 0 {
		return
	}
	for k, v := range ciAttrMap {
		tmpV, tmpErr := getCMDBAttributes(k, userToken, v)
		if tmpErr != nil {
			err = tmpErr
			break
		}
		for kk, vv := range tmpV {
			ciAttrSelectMap[kk] = vv
		}
	}
	if err != nil {
		return
	}
	for _, v := range output {
		if v.Entity == "" {
			continue
		}
		for _, vv := range v.Title {
			tmpKey := v.Entity + "_" + vv.Name
			if tmpSelectList, b := ciAttrSelectMap[tmpKey]; b {
				vv.SelectList = tmpSelectList
			}
		}
	}
	return
}

func getCMDBAttributes(entity, userToken string, attributes []string) (result map[string][]*models.EntityDataObj, err error) {
	result = make(map[string][]*models.EntityDataObj)
	req, newReqErr := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/wecmdb/api/v1/ci-types-attr/%s/attributes", models.Config.Wecube.BaseUrl, entity), nil)
	if newReqErr != nil {
		err = fmt.Errorf("Try to new http request fail,%s ", newReqErr.Error())
		return
	}
	req.Header.Set("Authorization", userToken)
	//req.Header.Set("Content-Type", "application/json")
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = fmt.Errorf("Try to do http request fail,%s ", respErr.Error())
		return
	}
	responseBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Request cmdb attribute fail,%s ", string(responseBytes))
		return
	}
	var attrQueryResp models.EntityAttributeQueryResponse
	err = json.Unmarshal(responseBytes, &attrQueryResp)
	if err != nil {
		err = fmt.Errorf("Json unmarshal attr response fail,%s ", err.Error())
		return
	}
	for _, v := range attrQueryResp.Data {
		existFlag := false
		for _, vv := range attributes {
			if v.PropertyName == vv {
				existFlag = true
				break
			}
		}
		if existFlag && v.SelectList != "" {
			tmpSelectList, tmpErr := getAttrCat(v.SelectList, userToken)
			if tmpErr != nil {
				err = tmpErr
				break
			}
			result[entity+"_"+v.PropertyName] = tmpSelectList
		}
	}
	return
}

func getAttrCat(catId, userToken string) (result []*models.EntityDataObj, err error) {
	result = []*models.EntityDataObj{}
	req, newReqErr := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/wecmdb/api/v1/base-key/categories/%s", models.Config.Wecube.BaseUrl, catId), nil)
	if newReqErr != nil {
		err = fmt.Errorf("Try to new http request fail,%s ", newReqErr.Error())
		return
	}
	req.Header.Set("Authorization", userToken)
	//req.Header.Set("Content-Type", "application/json")
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = fmt.Errorf("Try to do http request fail,%s ", respErr.Error())
		return
	}
	responseBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Request cmdb categories fail,%s ", string(responseBytes))
		return
	}
	var response models.CMDBCategoriesResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		err = fmt.Errorf("Json unmarshal categories response fail,%s ", err.Error())
		return
	}
	for _, v := range response.Data {
		result = append(result, &models.EntityDataObj{Id: v.Code, DisplayName: v.Value})
	}
	return
}

func BuildRequestProcessData(input models.RequestCacheData) (result models.RequestProcessData) {
	result.ProcDefId = input.ProcDefId
	result.ProcDefKey = input.ProcDefKey
	result.RootEntityOid = input.RootEntityValue.Oid
	result.Entities = []*models.RequestCacheEntityValue{}
	result.Bindings = []*models.RequestProcessTaskNodeBindObj{}
	entityExistMap := make(map[string]int)
	for _, node := range input.TaskNodeBindInfos {
		for _, entity := range node.BoundEntityValues {
			if _, b := entityExistMap[entity.Oid]; !b {
				result.Entities = append(result.Entities, entity)
				entityExistMap[entity.Oid] = 1
			}
			if node.NodeId != "" {
				result.Bindings = append(result.Bindings, &models.RequestProcessTaskNodeBindObj{Oid: entity.Oid, NodeId: node.NodeId, NodeDefId: node.NodeDefId, EntityDataId: entity.EntityDataId, BindFlag: entity.BindFlag})
			}
		}
	}
	return result
}

func AppendUselessEntity(requestTemplateId, userToken string, cacheData *models.RequestCacheData) error {
	if cacheData.RootEntityValue.Oid == "" || strings.HasPrefix(cacheData.RootEntityValue.Oid, "tmp") {
		return nil
	}
	preData, preErr := ProcessDataPreview(requestTemplateId, cacheData.RootEntityValue.Oid, userToken)
	if preErr != nil {
		return fmt.Errorf("Try to get process preview data fail,%s ", preErr.Error())
	}
	entityList := []*models.RequestCacheEntityValue{}
	entityExistMap := make(map[string]int)
	for _, v := range cacheData.TaskNodeBindInfos {
		for _, vv := range v.BoundEntityValues {
			if _, b := entityExistMap[vv.Oid]; !b {
				entityList = append(entityList, vv)
				entityExistMap[vv.Oid] = 1
			}
		}
	}
	preEntityList := []*models.EntityTreeObj{}
	rootSucceeding := []string{}
	for _, v := range preData.Data.EntityTreeNodes {
		if _, b := entityExistMap[v.Id]; !b {
			preEntityList = append(preEntityList, v)
		}
		if v.DataId == cacheData.RootEntityValue.Oid {
			rootSucceeding = v.SucceedingIds
		}
	}
	if len(preEntityList) == 0 {
		return nil
	}
	dependEntityMap := make(map[string]int)
	log.Logger.Debug("getDependEntity", log.StringList("rootSucceeding", rootSucceeding), log.Int("preLen", len(preEntityList)), log.Int("entityLen", len(entityList)))
	getDependEntity(rootSucceeding, preEntityList, entityList, dependEntityMap)
	for k, _ := range dependEntityMap {
		log.Logger.Debug("dependEntityMap", log.String("id", k))
	}
	if len(dependEntityMap) > 0 {
		newNode := models.RequestCacheTaskNodeBindObj{NodeId: "", NodeDefId: "", BoundEntityValues: []*models.RequestCacheEntityValue{}}
		for _, v := range preEntityList {
			if _, b := dependEntityMap[v.Id]; b {
				newNode.BoundEntityValues = append(newNode.BoundEntityValues, &models.RequestCacheEntityValue{Oid: v.Id, EntityDataId: v.DataId, PackageName: v.PackageName, EntityName: v.EntityName, EntityDisplayName: v.DisplayName, FullEntityDataId: v.FullDataId, AttrValues: []*models.RequestCacheEntityAttrValue{}, PreviousOids: v.PreviousIds, SucceedingOids: v.SucceedingIds})
			}
		}
		cacheData.TaskNodeBindInfos = append(cacheData.TaskNodeBindInfos, &newNode)
	}
	return nil
}

func getDependEntity(succeeding []string, preEntityList []*models.EntityTreeObj, entityList []*models.RequestCacheEntityValue, dependEntityMap map[string]int) {
	if len(succeeding) == 0 {
		return
	}
	// check if use by entity
	for _, v := range succeeding {
		vDataId := v
		if strings.Contains(vDataId, ":") {
			vDataId = vDataId[strings.LastIndex(vDataId, ":")+1:]
		}
		useFlag := false
		for _, vv := range entityList {
			for _, attr := range vv.AttrValues {
				if attr.DataType == "ref" {
					tmpV := fmt.Sprintf("%s", attr.DataValue)
					if strings.Contains(tmpV, vDataId) {
						useFlag = true
						break
					}
				}
			}
			if useFlag {
				break
			}
		}
		if !useFlag {
			// find other succeeding
			nextSucceeding := []string{}
			for _, vv := range preEntityList {
				if listContains(vv.PreviousIds, v) {
					nextSucceeding = append(nextSucceeding, vv.Id)
				}
			}
			if len(nextSucceeding) > 0 {
				getDependEntity(nextSucceeding, preEntityList, entityList, dependEntityMap)
				for _, vv := range nextSucceeding {
					if _, b := dependEntityMap[vv]; b {
						useFlag = true
						break
					}
				}
			}
		}
		if useFlag {
			for _, vv := range preEntityList {
				if vv.Id == v {
					dependEntityMap[v] = 1
				}
			}
		}
	}
}

func listContains(inputList []string, element string) bool {
	result := false
	for _, v := range inputList {
		if v == element {
			result = true
			break
		}
	}
	return result
}

func notifyRoleMail(requestId string) error {
	if !models.MailEnable {
		return nil
	}
	var roleTable []*models.RoleTable
	err := x.SQL("select id,email from `role` where id in (select `role` from request_template_role where role_type='MGMT' and request_template in (select request_template from request where id=?))", requestId).Find(&roleTable)
	if err != nil {
		return fmt.Errorf("Notify role mail query roles fail,%s ", err.Error())
	}
	if len(roleTable) == 0 {
		return nil
	}
	if roleTable[0].Email == "" {
		return nil
	}
	var requestTable []*models.RequestTable
	x.SQL("select t1.id,t1.name,t2.name as request_template,t1.reporter,t1.report_time,t1.emergency from request t1 left join request_template t2 on t1.request_template=t2.id where t1.id=?", requestId).Find(&requestTable)
	if len(requestTable) == 0 {
		return nil
	}
	var subject, content string
	subject = fmt.Sprintf("Taskman Request [%s] %s[%s]", models.PriorityLevelMap[requestTable[0].Emergency], requestTable[0].Name, requestTable[0].RequestTemplate)
	content = fmt.Sprintf("Taskman Request \nID:%s \nPriority:%s \nName:%s \nTemplate:%s \nReporter:%s \nReportTime:%s\n", requestTable[0].Id, models.PriorityLevelMap[requestTable[0].Emergency], requestTable[0].Name, requestTable[0].RequestTemplate, requestTable[0].Reporter, requestTable[0].ReportTime)
	err = models.MailSender.Send(subject, content, []string{roleTable[0].Email})
	if err != nil {
		return fmt.Errorf("Notify role mail fail,%s ", err.Error())
	}
	return nil
}
