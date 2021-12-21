package db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"io/ioutil"
	"net/http"
	"strings"
)

func GetCMDBRefSelectResult(input *models.RefSelectParam) (result []*models.EntityDataObj, err error) {
	result = []*models.EntityDataObj{}
	// if param map data have no new data -> get remote data + same entity new data
	refFlag, options, tmpErr := checkIfNeedAnalyze(input)
	log.Logger.Info("isContainNewMap", log.String("isContainNewMap", fmt.Sprintf("%d", refFlag)))
	if refFlag == -1 {
		return result, fmt.Errorf("AttrId:%s illegal ", input.AttrId)
	}
	if refFlag == 2 {
		return options, tmpErr
	}
	if refFlag == 0 {
		result, err = getRefDataWithoutFilter(input)
		return
	}
	// get attr filter
	filterString, getFilterErr := getCMDBRefFilter(input.AttrId, input.UserToken)
	if getFilterErr != nil {
		return result, getFilterErr
	}
	// if filter empty -> get remote data + same entity new data
	if filterString == "" {
		result, err = getRefDataWithoutFilter(input)
		return
	}
	input.Filter = filterString
	// analyze filter with map data
	result, err = analyzeFilterData(input)
	return
}

func checkIfNeedAnalyze(input *models.RefSelectParam) (refFlag int, options []*models.EntityDataObj, err error) {
	// refFlag 0 -> ref without new  1 -> ref with new  2 -> no ref use options  -1 -> attrId illegal
	refFlag = 0
	var entity, attrName, attrRef, dataOptions string
	if strings.Contains(input.AttrId, models.SysTableIdConnector) {
		entity = strings.Split(input.AttrId, models.SysTableIdConnector)[0]
		attrName = strings.Split(input.AttrId, models.SysTableIdConnector)[1]
	} else {
		refFlag = -1
		return refFlag, options, nil
	}
	var formItemTemplates []*models.FormItemTemplateTable
	//x.SQL("select id,name,ref_package_name,ref_entity,data_options from form_item_template where entity=? and form_template in (select form_template from request_template where id in (select request_template from request where id=?))", entity, input.RequestId).Find(&formItemTemplates)
	x.SQL("select distinct name,ref_package_name,ref_entity,data_options from form_item_template where entity=? and form_template in (select form_template from request_template where id in (select request_template from request where id=?) union select form_template from task_template where id in (select task_template from task where request=?))", entity, input.RequestId, input.RequestId).Find(&formItemTemplates)
	refColumnMap := make(map[string]int)
	for _, v := range formItemTemplates {
		if v.Name == attrName {
			attrRef = v.RefEntity
			dataOptions = v.DataOptions
		}
		if v.RefEntity != "" {
			refColumnMap[v.Name] = 1
		}
	}
	if attrRef == "" {
		refFlag = 2
		if strings.Contains(dataOptions, "http") {
			options, err = getRemoteEntityOptions(dataOptions, input.UserToken, input.Param.Dialect.AssociatedData)
		} else {
			options, err = getCMDBAttributeOptions(entity, attrName, input.UserToken)
		}
		return refFlag, options, err
	}
	if input.Param.Dialect != nil {
		for k, v := range input.Param.Dialect.AssociatedData {
			if _, b := refColumnMap[k]; b {
				if strings.HasPrefix(v, "tmp") {
					refFlag = 1
					break
				}
			}
		}
	}
	return refFlag, options, nil
}

func getCMDBRefFilter(attrId, userToken string) (filterString string, err error) {
	req, newReqErr := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/wecmdb/api/v1/ci-types-attr/%s/attributes", models.Config.Wecube.BaseUrl, strings.Split(attrId, models.SysTableIdConnector)[0]), nil)
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
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Try to request cmdb ci attr fail,response code:%d ", resp.StatusCode)
		return
	}
	var responseData models.CiTypeAttrQueryResponse
	respBodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal(respBodyBytes, &responseData)
	for _, v := range responseData.Data {
		if v.Id == attrId {
			filterString = v.RefFilter
			break
		}
	}
	return
}

func getRefDataWithoutFilter(input *models.RefSelectParam) (result []*models.EntityDataObj, err error) {
	log.Logger.Info("getRefDataWithoutFilter", log.JsonObj("input", input))
	result = []*models.EntityDataObj{}
	remoteRefData, refErr := getCMDBRefData(input)
	if refErr != nil {
		err = refErr
		return
	}
	for _, v := range remoteRefData {
		result = append(result, &models.EntityDataObj{Id: v.Guid, DisplayName: v.KeyName, IsNew: v.IsNew})
	}
	cacheData, cacheErr := getRequestCacheNewData(input.RequestId, input.AttrId)
	if cacheErr != nil {
		err = cacheErr
		return
	}
	for _, v := range cacheData {
		result = append(result, &models.EntityDataObj{Id: v.Guid, DisplayName: v.KeyName, IsNew: v.IsNew})
	}
	return
}

func getCMDBRefData(input *models.RefSelectParam) (result []*models.CiReferenceDataQueryObj, err error) {
	result = []*models.CiReferenceDataQueryObj{}
	paramBytes, tmpErr := json.Marshal(input.Param)
	if tmpErr != nil {
		err = fmt.Errorf("Json marshal param data fail,%s ", tmpErr.Error())
		return
	}
	req, newReqErr := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/wecmdb/api/v1/ci-data/reference-data/query/%s", models.Config.Wecube.BaseUrl, input.AttrId), bytes.NewReader(paramBytes))
	if newReqErr != nil {
		err = fmt.Errorf("Try to new http request fail,%s ", newReqErr.Error())
		return
	}
	req.Header.Set("Authorization", input.UserToken)
	req.Header.Set("Content-Type", "application/json")
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = fmt.Errorf("Try to do http request fail,%s ", respErr.Error())
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Try to request cmdb reference data fail,statusCode:%d ", resp.StatusCode)
		return
	}
	var response models.CiReferenceDataQueryResponse
	responseBody, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal(responseBody, &response)
	result = response.Data
	return
}

func getRefSelectEntity(requestId, attrId string) (refEntity string, err error) {
	var formItemTemplates []*models.FormItemTemplateTable
	attrSplit := strings.Split(attrId, models.SysTableIdConnector)
	x.SQL("select id,ref_package_name,ref_entity from form_item_template where entity=? and name=? and form_template in (select form_template from request_template where id in (select request_template from request where id=?)  union select form_template from task_template where request_template in (select request_template from request where id=?))", attrSplit[0], attrSplit[1], requestId, requestId).Find(&formItemTemplates)
	if len(formItemTemplates) == 0 {
		return refEntity, fmt.Errorf("Can not find form item template with entity:%s name:%s ", attrSplit[0], attrSplit[1])
	}
	return formItemTemplates[0].RefEntity, nil
}

func getRequestCacheNewData(requestId, attrId string) (result []*models.CiReferenceDataQueryObj, err error) {
	refEntity, tmpErr := getRefSelectEntity(requestId, attrId)
	if tmpErr != nil {
		return result, tmpErr
	}
	cacheDataObj, cacheErr := GetRequestCache(requestId, "data")
	if cacheErr != nil {
		return result, cacheErr
	}
	cacheData := cacheDataObj.(models.RequestPreDataDto)
	for _, v := range cacheData.Data {
		if v.Entity == refEntity {
			for _, vv := range v.Value {
				if strings.HasPrefix(vv.Id, "tmp") {
					result = append(result, &models.CiReferenceDataQueryObj{Guid: vv.Id, KeyName: vv.DisplayName, IsNew: true})
				}
			}
		}
	}
	return
}

func analyzeFilterData(input *models.RefSelectParam) (result []*models.EntityDataObj, err error) {
	log.Logger.Debug("analyzeFilterData", log.JsonObj("input", input))
	var filters []map[string]models.CiDataRefFilterObj
	err = json.Unmarshal([]byte(input.Filter), &filters)
	if err != nil {
		err = fmt.Errorf("Json unmarshal filters string fail,%s ", err.Error())
		return
	}
	if len(filters) == 0 {
		err = fmt.Errorf("Get ci reference data fail,filters string illgeal ")
		return
	}
	newData, tmpErr := getRequestNewData(input.RequestId)
	if tmpErr != nil {
		return result, fmt.Errorf("Try to get request new data fail,%s ", tmpErr.Error())
	}
	refEntity, tmpErr := getRefSelectEntity(input.RequestId, input.AttrId)
	if tmpErr != nil {
		return result, tmpErr
	}
	attrName := strings.Split(input.AttrId, models.SysTableIdConnector)[1]
	filterMap := input.Param.Dialect.AssociatedData
	if _, b := filterMap[attrName]; b {
		delete(filterMap, attrName)
	}
	var filterQueryParam models.EntityQueryParam
	for _, filter := range filters[0] {
		tmpParam := models.GetExpressResultParam{
			Filter:    &filter,
			UserToken: input.UserToken,
			FilterMap: filterMap,
			NewData:   newData,
		}
		tmpFilterQuery, tmpErr := getRefFilterParam(&tmpParam)
		if tmpErr != nil {
			err = tmpErr
			break
		}
		filterQueryParam.AdditionalFilters = append(filterQueryParam.AdditionalFilters, &tmpFilterQuery)
	}
	if err != nil {
		return
	}
	// query ci data
	queryResult, queryErr := getCiData(filterQueryParam, refEntity, input.UserToken, newData)
	if queryErr != nil {
		return result, queryErr
	}
	for _, v := range queryResult {
		result = append(result, &models.EntityDataObj{Id: fmt.Sprintf("%s", v["id"]), DisplayName: fmt.Sprintf("%s", v["displayName"])})
	}
	return
}

func getRefFilterParam(param *models.GetExpressResultParam) (queryParam models.EntityQueryObj, err error) {
	filter := param.Filter
	startCiType := filter.Left[:strings.LastIndex(filter.Left, "[")]
	if _, b := param.FilterMap[startCiType]; b {
		delete(param.FilterMap, startCiType)
	}
	column := filter.Left[strings.LastIndex(filter.Left, "[")+1 : strings.LastIndex(filter.Left, "]")]
	var valueList []string
	if filter.Right.Type == "expression" {
		//Example: {"type":"expression","value":"app_instance.unit>unit.resource_set>resource_set~(resource_set)host_resource:[guid]"}
		param.Express = fmt.Sprintf("%s", filter.Right.Value)
		param.StartCiType = startCiType
		valueList, err = getExpressResultList(param)
		if err != nil {
			return
		}
		if len(valueList) == 0 {
			if _, b := param.FilterMap[column]; b {
				if param.FilterMap[column] != "" {
					valueList = append(valueList, param.FilterMap[column])
				}
			}
		}
	} else if filter.Right.Type == "array" {
		//Example: ["JAVA","NGINX","MYSQL"]
		for _, rv := range filter.Right.Value.([]interface{}) {
			valueList = append(valueList, rv.(string))
		}
	}
	queryParam.AttrName = column
	queryParam.Op = filter.Operator
	queryParam.Condition = valueList
	return
}

type expressionSqlObj struct {
	Table           string                   `json:"table"`
	IndexTableName  string                   `json:"index_table_name"`
	LeftJoinColumn  string                   `json:"left_join_column"`
	RightJoinColumn string                   `json:"right_join_column"`
	WhereSql        string                   `json:"where_sql"`
	Filters         []*models.EntityQueryObj `json:"filters"`
	ResultColumn    string                   `json:"result_column"`
	RefColumn       string                   `json:"ref_column"`
	MultiRefTable   string                   `json:"multi_ref_table"`
}

func getExpressResultList(param *models.GetExpressResultParam) (result []string, err error) {
	result = []string{}
	express := param.Express
	// Example expression -> "host_resource_instance.resource_set>resource_set~(resource_set)unit[{key_name eq 'hhh'},{code in ['u','v']}]:[guid]"
	var ciList, filterParams, tmpSplitList []string
	// replace content 'xxx' to '$1' in case of content have '>~.:()[]'
	if strings.Contains(express, "'") {
		tmpSplitList = strings.Split(express, "'")
		express = ""
		for i, v := range tmpSplitList {
			if i%2 == 0 {
				if i == len(tmpSplitList)-1 {
					express += v
				} else {
					express += fmt.Sprintf("%s'$%d'", v, i/2)
				}
			} else {
				filterParams = append(filterParams, strings.ReplaceAll(v, "'", ""))
			}
		}
	}
	// split with > or ~
	var cursor int
	for i, v := range express {
		if v == 62 || v == 126 {
			ciList = append(ciList, express[cursor:i])
			cursor = i
		}
	}
	ciList = append(ciList, express[cursor:])
	// analyze each ci segment
	var expressionSqlList []*expressionSqlObj
	for i, ci := range ciList {
		eso := expressionSqlObj{IndexTableName: fmt.Sprintf("t%d", i)}
		if strings.HasPrefix(ci, ">") {
			eso.LeftJoinColumn = ciList[i-1][strings.LastIndex(ciList[i-1], ".")+1:]
			ci = ci[1:]
		} else if strings.HasPrefix(ci, "~") {
			eso.RightJoinColumn = ci[2:strings.Index(ci, ")")]
			eso.RefColumn = eso.RightJoinColumn
			ci = ci[strings.Index(ci, ")")+1:]
		}
		// ASCII . -> 46 , [ -> 91 , ] -> 93 , : -> 58 , { -> 123 , } -> 125
		for j, v := range ci {
			if v == 46 || v == 58 || v == 91 {
				eso.Table = ci[:j]
				ci = ci[j:]
				break
			}
		}
		if eso.Table == "" {
			eso.Table = ci
		}
		if ci[0] == 91 {
			tmpFilterStr := ci[2 : len(ci)-2]
			for j, v := range ci {
				if v == 46 || v == 58 {
					tmpFilterStr = ci[2 : j-2]
					ci = ci[j:]
					break
				}
			}
			for _, v := range strings.Split(tmpFilterStr, "},{") {
				tmpFilterList := strings.Split(v, " ")
				if len(tmpFilterList) > 2 {
					tmpCondValue := tmpFilterList[2]
					for i, v := range filterParams {
						tmpCondValue = strings.ReplaceAll(tmpCondValue, fmt.Sprintf("$%d", i), v)
					}
					eso.Filters = append(eso.Filters, &models.EntityQueryObj{AttrName: tmpFilterList[0], Op: tmpFilterList[1], Condition: transConditionValue(tmpCondValue, tmpFilterList[1])})
				} else {
					eso.Filters = append(eso.Filters, &models.EntityQueryObj{AttrName: tmpFilterList[0], Op: tmpFilterList[1]})
				}
			}
		}
		if ci[0] == 58 {
			eso.ResultColumn = ci[2 : len(ci)-1]
		}
		if ci[0] == 46 {
			eso.RefColumn = ci[1:]
		}
		expressionSqlList = append(expressionSqlList, &eso)
	}
	eLen := len(expressionSqlList)
	if eLen == 0 {
		return
	}
	var tmpRows []map[string]interface{}
	startRow := make(map[string]interface{})
	for k, v := range param.FilterMap {
		startRow[k] = v
	}
	tmpRows = append(tmpRows, startRow)
	tmpLength := len(expressionSqlList) - 1
	for i, v := range expressionSqlList {
		log.Logger.Info("expressionSqlList", log.Int("index", i), log.JsonObj("v", v))
		log.Logger.Info("tmpRows", log.Int("len", len(tmpRows)), log.JsonObj("data", tmpRows))
		tmpGuidList := []string{}
		tmpParam := models.EntityQueryParam{}
		if len(v.Filters) > 0 {
			tmpParam.AdditionalFilters = v.Filters
		}
		if v.RightJoinColumn != "" {
			for _, tmpRow := range tmpRows {
				tmpGuidList = append(tmpGuidList, fmt.Sprintf("%s", tmpRow["guid"]))
			}
			tmpParam.AdditionalFilters = []*models.EntityQueryObj{{AttrName: v.RightJoinColumn, Op: "in", Condition: tmpGuidList}}
			tmpRows, err = getCiData(tmpParam, v.Table, param.UserToken, param.NewData)
			if err != nil {
				break
			}
		}
		if i < tmpLength {
			if expressionSqlList[i+1].LeftJoinColumn != "" {
				for _, tmpRow := range tmpRows {
					tmpGuidList = append(tmpGuidList, fmt.Sprintf("%s", tmpRow[v.RefColumn]))
				}
				tmpParam.Criteria.AttrName = "guid"
				tmpParam.Criteria.Op = "in"
				tmpParam.Criteria.Condition = tmpGuidList
				tmpRows, err = getCiData(tmpParam, expressionSqlList[i+1].Table, param.UserToken, param.NewData)
				if err != nil {
					break
				}
			}
		}
		if v.ResultColumn != "" {
			for _, tmpRow := range tmpRows {
				result = append(result, fmt.Sprintf("%s", tmpRow[v.ResultColumn]))
			}
		}
	}
	return
}

func getCiData(param models.EntityQueryParam, ciType, userToken string, newData map[string]map[string]interface{}) (result []map[string]interface{}, err error) {
	paramBytes, tmpErr := json.Marshal(param)
	if tmpErr != nil {
		err = fmt.Errorf("Json marshal param data fail,%s ", tmpErr.Error())
		return
	}
	log.Logger.Info("getCiData", log.String("ciType", ciType), log.String("param", string(paramBytes)))
	req, newReqErr := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/wecmdb/entities/%s/query", models.Config.Wecube.BaseUrl, ciType), bytes.NewReader(paramBytes))
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
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Try to request cmdb reference data fail,statusCode:%d ", resp.StatusCode)
		return
	}
	var response models.EntityResponse
	responseBody, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	log.Logger.Info("getCiDataRemote", log.String("ciType", ciType), log.String("response", string(responseBody)))
	json.Unmarshal(responseBody, &response)
	result = response.Data
	newGuidList := checkQueryParamContainNewData(param)
	if len(newGuidList) > 0 {
		for _, v := range newGuidList {
			if _, b := newData[v]; b {
				result = append(result, newData[v])
			}
		}
	}
	log.Logger.Info("getCiDataResult", log.StringList("newGuid", newGuidList), log.JsonObj("result", result))
	return
}

func checkQueryParamContainNewData(param models.EntityQueryParam) []string {
	var guidList []string
	tmpMap := make(map[string]int)
	if param.Criteria.AttrName == "id" || param.Criteria.AttrName == "guid" {
		if param.Criteria.Op == "in" {
			for _, v := range param.Criteria.Condition.([]string) {
				if _, b := tmpMap[v]; !b {
					guidList = append(guidList, v)
					tmpMap[v] = 1
				}
			}
		}
	}
	for _, filters := range param.AdditionalFilters {
		if filters.AttrName == "id" || filters.AttrName == "guid" {
			if filters.Op == "in" {
				for _, v := range filters.Condition.([]string) {
					if _, b := tmpMap[v]; !b {
						guidList = append(guidList, v)
						tmpMap[v] = 1
					}
				}
			}
		}
	}
	return guidList
}

func transConditionValue(input, op string) interface{} {
	if op == "in" {
		input = input[1 : len(input)-1]
		input = strings.ReplaceAll(input, "'", "")
		return strings.Split(input, ",")
	}
	input = strings.ReplaceAll(input, "'", "")
	return input
}

func getRequestNewData(requestId string) (result map[string]map[string]interface{}, err error) {
	cacheDataObj, cacheErr := GetRequestCache(requestId, "data")
	if cacheErr != nil {
		return result, cacheErr
	}
	cacheData := cacheDataObj.(models.RequestPreDataDto)
	// TODO update data
	result = make(map[string]map[string]interface{})
	for _, entity := range cacheData.Data {
		for _, valueObj := range entity.Value {
			if !strings.HasPrefix(valueObj.Id, "tmp") {
				continue
			}
			result[valueObj.Id] = valueObj.EntityData
			log.Logger.Info("getRequestNewData", log.String("id", valueObj.Id))
		}
	}
	return
}

func FilterInSideData(input []*models.EntityDataObj, attrId, requestId string) (output []*models.EntityDataObj) {
	output = input
	attrSplit := strings.Split(attrId, models.SysTableIdConnector)
	var formItemTemplate []*models.FormItemTemplateTable
	x.SQL("select * from form_item_template where entity=? and name=? and form_template in (select form_template from request_template where id in (select request_template from request where id=?))", attrSplit[0], attrSplit[1], requestId).Find(&formItemTemplate)
	if len(formItemTemplate) == 0 {
		return output
	}
	if formItemTemplate[0].IsRefInside == "no" {
		return output
	}
	var formItems []*models.FormItemTable
	x.SQL("select distinct row_data_id from form_item where form in (select form from request where id=?)", requestId).Find(&formItems)
	rowDataMap := make(map[string]int)
	for _, v := range formItems {
		tmpV := v.RowDataId
		if strings.Contains(tmpV, ":") {
			tmpV = tmpV[strings.LastIndex(tmpV, ":")+1:]
		}
		rowDataMap[tmpV] = 1
	}
	output = []*models.EntityDataObj{}
	for _, v := range input {
		if _, b := rowDataMap[v.Id]; b {
			output = append(output, v)
		}
	}
	return output
}

func getCMDBAttributeOptions(entity, attribute, userToken string) (result []*models.EntityDataObj, err error) {
	result = []*models.EntityDataObj{}
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
		if v.PropertyName == attribute {
			result, err = getAttrCat(v.SelectList, userToken)
			break
		}
	}
	return
}

func getRemoteEntityOptions(url, userToken string, inputMap map[string]string) (result []*models.EntityDataObj, err error) {
	method := http.MethodGet
	if strings.HasPrefix(strings.ToLower(url), "post") {
		method = http.MethodPost
		url = url[strings.Index(url, "http"):]
	}
	for k, v := range inputMap {
		url = strings.ReplaceAll(url, fmt.Sprintf("{{%s}}", k), v)
	}
	reqParam := ""
	if method == http.MethodPost && strings.Contains(url, "=") {
		reqParam = url[strings.Index(url, "=")+1:]
		url = url[:strings.Index(url, "?")]
	}
	log.Logger.Info("curl remote entity options", log.String("url", url), log.String("method", method), log.String("param", reqParam))
	req, reqErr := http.NewRequest(method, url, strings.NewReader(reqParam))
	if reqErr != nil {
		err = fmt.Errorf("Try to new request fail,%s ", reqErr.Error())
		return
	}
	req.Header.Set("Authorization", userToken)
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = fmt.Errorf("Try to do request fail,%s ", respErr.Error())
		return
	}
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 {
		return result, fmt.Errorf("Do request fail,response statusCode:%d resp:%s ", resp.StatusCode, string(b))
	}
	var response models.WorkflowEntityQuery
	err = json.Unmarshal(b, &response)
	if err != nil {
		return result, fmt.Errorf("Json unmarshal response body fail,%s ", err.Error())
	}
	for _, v := range response.Data {
		result = append(result, &models.EntityDataObj{Id: v.Id, DisplayName: v.DisplayName})
	}
	return result, nil
}
