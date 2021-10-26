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
	isContainNewMap := false
	if input.Param.Dialect != nil {
		for _, v := range input.Param.Dialect.AssociatedData {
			if strings.HasPrefix(v, "tmp") {
				isContainNewMap = true
				break
			}
		}
	}
	log.Logger.Info("isContainNewMap", log.String("isContainNewMap", fmt.Sprintf("%v", isContainNewMap)))
	if !isContainNewMap {
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

func getRequestCacheNewData(requestId, attrId string) (result []*models.CiReferenceDataQueryObj, err error) {
	var formItemTemplates []*models.FormItemTemplateTable
	attrSplit := strings.Split(attrId, models.SysTableIdConnector)
	x.SQL("select id,ref_package_name,ref_entity from form_item_template where entity=? and name=? and form_template in (select form_template from request_template where id in (select request_template from request where id=?))", attrSplit[0], attrSplit[1], requestId).Find(&formItemTemplates)
	if len(formItemTemplates) == 0 {
		return result, fmt.Errorf("Can not find form item template with entity:%s name:%s ", attrSplit[0], attrSplit[1])
	}
	cacheDataObj, cacheErr := GetRequestCache(requestId, "data")
	if cacheErr != nil {
		return result, cacheErr
	}
	cacheData := cacheDataObj.(models.RequestPreDataDto)
	for _, v := range cacheData.Data {
		if v.Entity == formItemTemplates[0].RefEntity {
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
	log.Logger.Info("analyzeFilterData", log.JsonObj("input", input))
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
	ciType := strings.Split(input.AttrId, models.SysTableIdConnector)[0]
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
	queryResult, queryErr := getCiData(filterQueryParam, ciType, input.UserToken, newData)
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
	Table           string
	IndexTableName  string
	LeftJoinColumn  string
	RightJoinColumn string
	WhereSql        string
	Filters         []*models.EntityQueryObj
	ResultColumn    string
	RefColumn       string
	MultiRefTable   string
}

func getExpressResultList(param *models.GetExpressResultParam) (result []string, err error) {
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
