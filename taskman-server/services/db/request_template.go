package db

import (
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

func QueryRequestTemplateGroup(param *models.QueryRequestParam, userRoles []string) (pageInfo models.PageInfo, rowData []*models.RequestTemplateGroupTable, err error) {
	rowData = []*models.RequestTemplateGroupTable{}
	filterSql, queryColumn, queryParam := transFiltersToSQL(param, &models.TransFiltersParam{IsStruct: true, StructObj: models.RequestTemplateGroupTable{}, PrimaryKey: "id"})
	baseSql := fmt.Sprintf("SELECT %s FROM request_template_group WHERE manage_role in ('"+strings.Join(userRoles, "','")+"') and del_flag=0 %s ", queryColumn, filterSql)
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
		roleMap, _ := getRoleMap()
		for _, row := range rowData {
			if row.ManageRole != "" {
				row.ManageRoleObj = models.RoleTable{Id: row.ManageRole, DisplayName: roleMap[row.ManageRole].DisplayName}
			}
		}
	}
	return
}

func CreateRequestTemplateGroup(param *models.RequestTemplateGroupTable) error {
	param.Id = guid.CreateGuid()
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err := x.Exec("insert into request_template_group(id,name,description,manage_role,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?,?)",
		param.Id, param.Name, param.Description, param.ManageRole, param.CreatedBy, nowTime, param.CreatedBy, nowTime)
	if err != nil {
		err = fmt.Errorf("Insert database error:%s ", err.Error())
	}
	return err
}

func UpdateRequestTemplateGroup(param *models.RequestTemplateGroupTable) error {
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err := x.Exec("update request_template_group set name=?,description =?,manage_role=?,updated_by=?,updated_time=? where id=?",
		param.Name, param.Description, param.ManageRole, param.UpdatedBy, nowTime, param.Id)
	if err != nil {
		err = fmt.Errorf("Update database error:%s ", err.Error())
	}
	return err
}

func DeleteRequestTemplateGroup(id string) error {
	_, err := x.Exec("update request_template_group set del_flag=1 where id=?", id)
	if err != nil {
		err = fmt.Errorf("Delete database error:%s ", err.Error())
	}
	return err
}

func CheckRequestTemplateGroupRoles(id string, roles []string) error {
	rowData, err := x.QueryString("select id from request_template_group where id=? and manage_role in ('" + strings.Join(roles, "','") + "')")
	if err != nil {
		return fmt.Errorf("Try to query database data fail,%s ", err.Error())
	}
	if len(rowData) == 0 {
		return fmt.Errorf(models.RowDataPermissionErr)
	}
	return nil
}

func GetCoreProcessListNew(userToken string) (processList []*models.ProcDefObj, err error) {
	processList = []*models.ProcDefObj{}
	req, reqErr := http.NewRequest(http.MethodGet, models.Config.Wecube.BaseUrl+"/platform/v1/public/process/definitions", nil)
	if reqErr != nil {
		err = fmt.Errorf("Try to new http request to core fail,%s ", reqErr.Error())
		return
	}
	req.Header.Set("Authorization", userToken)
	http.DefaultClient.CloseIdleConnections()
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = fmt.Errorf("Try to do request to core fail,%s ", respErr.Error())
		return
	}
	var respObj models.ProcQueryResponse
	respBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	err = json.Unmarshal(respBytes, &respObj)
	log.Logger.Debug("Get core process list", log.String("body", string(respBytes)))
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	if respObj.Status != "OK" {
		err = fmt.Errorf(respObj.Message)
		return
	}
	processList = respObj.Data
	return
}

func GetCoreProcessListAll(userToken string) (processList []*models.ProcAllDefObj, err error) {
	processList = []*models.ProcAllDefObj{}
	req, reqErr := http.NewRequest(http.MethodGet, models.Config.Wecube.BaseUrl+"/platform/v1/process/definitions?includeDraft=0&permission=USE", nil)
	if reqErr != nil {
		err = fmt.Errorf("Try to new http request to core fail,%s ", reqErr.Error())
		return
	}
	req.Header.Set("Authorization", userToken)
	http.DefaultClient.CloseIdleConnections()
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = fmt.Errorf("Try to do request to core fail,%s ", respErr.Error())
		return
	}
	var respObj models.ProcAllQueryResponse
	respBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	err = json.Unmarshal(respBytes, &respObj)
	log.Logger.Debug("Get core process list", log.String("body", string(respBytes)))
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	if respObj.Status != "OK" {
		err = fmt.Errorf(respObj.Message)
		return
	}
	processList = respObj.Data
	return
}

func GetProcessNodesByProc(requestTemplateId, userToken string, filterType string) (nodeList models.ProcNodeObjList, err error) {
	requestTemplateObj, tmpErr := getSimpleRequestTemplate(requestTemplateId)
	if tmpErr != nil {
		err = tmpErr
		return
	}
	if requestTemplateObj.ProcDefId == "" {
		err = fmt.Errorf("Request template proDefId illegal ")
		return
	}
	nodeList = []*models.ProcNodeObj{}
	req, reqErr := http.NewRequest(http.MethodGet, models.Config.Wecube.BaseUrl+"/platform/v1/public/process/definitions/"+requestTemplateObj.ProcDefId+"/tasknodes", nil)
	if reqErr != nil {
		err = fmt.Errorf("Try to new http request to core fail,%s ", reqErr.Error())
		return
	}
	req.Header.Set("Authorization", userToken)
	http.DefaultClient.CloseIdleConnections()
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = fmt.Errorf("Try to do request to core fail,%s ", respErr.Error())
		return
	}
	var respObj models.ProcNodeQueryResponse
	respBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	err = json.Unmarshal(respBytes, &respObj)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	if respObj.Status != "OK" {
		err = fmt.Errorf(respObj.Message)
		return
	}
	for _, v := range respObj.Data {
		if v.NodeType != "subProcess" {
			continue
		}
		if filterType == "template" {
			if v.TaskCategory != "SUTN" {
				continue
			}
		} else if filterType == "bind" {
			if v.DynamicBind == "Y" {
				continue
			}
		}
		if v.OrderedNo == "" {
			v.OrderedNum = 0
		} else {
			v.OrderedNum, _ = strconv.Atoi(v.OrderedNo)
		}
		nodeList = append(nodeList, v)
	}
	sort.Sort(nodeList)
	return
}

func GetCoreProcessList(userToken string) (processList []*models.CodeProcessQueryObj, err error) {
	req, reqErr := http.NewRequest(http.MethodGet, models.Config.Wecube.BaseUrl+"/platform/v1/process/definitions?includeDraft=0&permission=USE&tags="+models.ProcessFetchTabs, nil)
	if reqErr != nil {
		err = fmt.Errorf("Try to new http request to core fail,%s ", reqErr.Error())
		return
	}
	req.Header.Set("Authorization", userToken)
	http.DefaultClient.CloseIdleConnections()
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		err = fmt.Errorf("Try to do request to core fail,%s ", respErr.Error())
		return
	}
	var respObj models.CoreProcessQueryResponse
	respBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	err = json.Unmarshal(respBytes, &respObj)
	log.Logger.Debug("Get core process list", log.String("body", string(respBytes)))
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	if len(respObj.Data) == 0 {
		processList = []*models.CodeProcessQueryObj{}
		return
	}
	procMap := make(map[string]*models.CodeProcessQueryObj)
	for _, v := range respObj.Data {
		tmpT, tmpErr := time.Parse(models.DateTimeFormat, v.CreatedTime)
		if tmpErr == nil {
			v.CreatedUnixTime = tmpT.Unix()
		}
		if oldProc, b := procMap[v.ProcDefKey]; b {
			if oldProc.CreatedUnixTime < v.CreatedUnixTime {
				procMap[v.ProcDefKey] = v
			}
		} else {
			procMap[v.ProcDefKey] = v
		}
	}
	for _, v := range respObj.Data {
		if procMap[v.ProcDefKey].ProcDefId == v.ProcDefId {
			processList = append(processList, v)
		}
	}
	return
}

func SyncCoreRole() {
	if models.CoreToken.BaseUrl == "" {
		return
	}
	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/platform/v1/roles/retrieve", models.CoreToken.BaseUrl), strings.NewReader(""))
	if err != nil {
		log.Logger.Error("Get core role key new request fail", log.Error(err))
		return
	}
	request.Header.Set("Authorization", models.CoreToken.GetCoreToken())
	res, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Logger.Error("Get core role key ctxhttp request fail", log.Error(err))
		return
	}
	b, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	//log.Logger.Debug("Get core role response", log.String("body", string(b)))
	var result models.CoreRoleDto
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Logger.Error("Get core role key json unmarshal result", log.Error(err))
		return
	}
	if len(result.Data) == 0 {
		log.Logger.Warn("Get core role key fail with no data")
		return
	}
	var roleTable, addRoleList, delRoleList []*models.RoleTable
	err = x.SQL("select * from role").Find(&roleTable)
	if err != nil {
		log.Logger.Error("Try to sync core role fail", log.Error(err))
		return
	}
	for _, v := range result.Data {
		existFlag := false
		for _, vv := range roleTable {
			if v.Name == vv.Id {
				existFlag = true
				break
			}
		}
		if !existFlag {
			addRoleList = append(addRoleList, &models.RoleTable{Id: v.Name, DisplayName: v.DisplayName, CoreId: v.Id})
		}
	}
	for _, v := range roleTable {
		existFlag := false
		for _, vv := range result.Data {
			if v.Id == vv.Name {
				existFlag = true
				break
			}
		}
		if !existFlag {
			delRoleList = append(delRoleList, &models.RoleTable{Id: v.Id})
		}
	}
	var actions []*execAction
	for _, role := range addRoleList {
		actions = append(actions, &execAction{Sql: "insert into `role`(id,display_name,core_id,updated_time) value (?,?,?,NOW())", Param: []interface{}{role.Id, role.DisplayName, role.CoreId}})
	}
	if len(delRoleList) > 0 {
		roleIdList := []string{}
		for _, role := range delRoleList {
			actions = append(actions, &execAction{Sql: "delete from `role` where id=?", Param: []interface{}{role.Id}})
			roleIdList = append(roleIdList, role.Id)
		}
		actions = append(actions, &execAction{Sql: "update form_template set `role`=NULL where `role` in ('" + strings.Join(roleIdList, "','") + "')"})
		actions = append(actions, &execAction{Sql: "update request_template_group set manage_role=NULL where manage_role in ('" + strings.Join(roleIdList, "','") + "')"})
	}
	if len(actions) > 0 {
		err = transactionWithoutForeignCheck(actions)
		if err != nil {
			log.Logger.Error("Sync core role fail", log.Error(err))
		}
	}
}

func GetRoleList(ids []string) (result []*models.RoleTable, err error) {
	result = []*models.RoleTable{}
	if len(ids) == 0 {
		err = x.SQL("select * from role").Find(&result)
	} else {
		err = x.SQL("select * from role where id in ('" + strings.Join(ids, "','") + "')").Find(&result)
	}
	return
}

func QueryRequestTemplate(param *models.QueryRequestParam, userToken string, userRoles []string) (pageInfo models.PageInfo, result []*models.RequestTemplateQueryObj, err error) {
	extFilterSql := ""
	result = []*models.RequestTemplateQueryObj{}
	isQueryMessage := false
	if len(param.Filters) > 0 {
		newFilters := []*models.QueryRequestFilterObj{}
		for _, v := range param.Filters {
			if v.Name == "id" {
				isQueryMessage = true
			}
			if v.Name == "mgmtRoles" || v.Name == "useRoles" {
				inValueList := v.Value.([]interface{})
				inValueStringList := []string{}
				for _, inValueInterfaceObj := range inValueList {
					if inValueInterfaceObj == nil {
						inValueStringList = append(inValueStringList, "")
					} else {
						inValueStringList = append(inValueStringList, inValueInterfaceObj.(string))
					}
				}
				if len(inValueStringList) == 0 {
					continue
				}
				var tmpIds []string
				var tmpErr error
				if v.Name == "mgmtRoles" {
					tmpIds, tmpErr = getRequestTemplateIdsBySql("select t1.id from request_template t1 left join request_template_role t2 on t1.id=t2.request_template where t2.role_type='MGMT' and t2.role in ('" + strings.Join(inValueStringList, "','") + "')")
				} else {
					tmpIds, tmpErr = getRequestTemplateIdsBySql("select t1.id from request_template t1 left join request_template_role t2 on t1.id=t2.request_template where t2.role_type='USE' and t2.role in ('" + strings.Join(inValueStringList, "','") + "')")
				}
				if tmpErr != nil {
					err = fmt.Errorf("Try to query filter role id fail,%s ", tmpErr.Error())
					break
				}
				extFilterSql += " and id in ('" + strings.Join(tmpIds, "','") + "') "
			} else {
				newFilters = append(newFilters, v)
			}
		}
		if err != nil {
			return models.PageInfo{}, nil, err
		}
		param.Filters = newFilters
	}
	rowData := []*models.RequestTemplateTable{}
	filterSql, queryColumn, queryParam := transFiltersToSQL(param, &models.TransFiltersParam{IsStruct: true, StructObj: models.RequestTemplateTable{}, PrimaryKey: "id", Prefix: "t1"})
	baseSql := fmt.Sprintf("SELECT %s FROM (select * from request_template where del_flag=0 or (del_flag=2 and id not in (select record_id from request_template where del_flag=2 and record_id<>''))) t1 WHERE t1.id in (select request_template from request_template_role where role_type='MGMT' and `role` in ('"+strings.Join(userRoles, "','")+"')) %s %s ", queryColumn, extFilterSql, filterSql)
	if param.Paging {
		pageInfo.StartIndex = param.Pageable.StartIndex
		pageInfo.PageSize = param.Pageable.PageSize
		pageInfo.TotalRows = queryCount(baseSql, queryParam...)
		pageSql, pageParam := transPageInfoToSQL(*param.Pageable)
		baseSql += pageSql
		queryParam = append(queryParam, pageParam...)
	}
	err = x.SQL(baseSql, queryParam...).Find(&rowData)
	if len(rowData) == 0 || err != nil {
		return
	}
	var rtIds []string
	for _, row := range rowData {
		rtIds = append(rtIds, row.Id)
	}
	queryRoleSql := "select t4.id,GROUP_CONCAT(t4.role_obj) as 'role','mgmt' as 'role_type' from ("
	queryRoleSql += "select t1.id,CONCAT(t2.role,'::',t3.display_name) as 'role_obj' from request_template t1 left join request_template_role t2 on t1.id=t2.request_template left join role t3 on t2.role=t3.id where t1.id in ('" + strings.Join(rtIds, "','") + "') and t2.role_type='MGMT'"
	queryRoleSql += ") t4 group by t4.id"
	queryRoleSql += " UNION "
	queryRoleSql += "select t4.id,GROUP_CONCAT(t4.role_obj) as 'role','use' as 'role_type' from ("
	queryRoleSql += "select t1.id,CONCAT(t2.role,'::',t3.display_name) as 'role_obj' from request_template t1 left join request_template_role t2 on t1.id=t2.request_template left join role t3 on t2.role=t3.id where t1.id in ('" + strings.Join(rtIds, "','") + "') and t2.role_type='USE'"
	queryRoleSql += ") t4 group by t4.id"
	var requestTemplateRows []*models.RequestTemplateRoleTable
	err = x.SQL(queryRoleSql).Find(&requestTemplateRows)
	if err != nil {
		return
	}
	var mgmtRoleMap = make(map[string][]*models.RoleTable)
	var useRoleMap = make(map[string][]*models.RoleTable)
	for _, v := range requestTemplateRows {
		tmpRoles := []*models.RoleTable{}
		for _, vv := range strings.Split(v.Role, ",") {
			tmpSplit := strings.Split(vv, "::")
			tmpRoles = append(tmpRoles, &models.RoleTable{Id: tmpSplit[0], DisplayName: tmpSplit[1]})
		}
		if v.RoleType == "mgmt" {
			mgmtRoleMap[v.Id] = tmpRoles
		} else {
			useRoleMap[v.Id] = tmpRoles
		}
	}
	if isQueryMessage {
		actions := []*execAction{}
		for _, v := range rowData {
			tmpExist, newProDefId, tmpErr := checkProDefId(v.ProcDefId, v.ProcDefName, v.ProcDefKey, userToken)
			if tmpErr != nil {
				err = fmt.Errorf("Try to sync new proDefId fail,%s ", tmpErr.Error())
				break
			}
			if !tmpExist {
				v.ProcDefId = newProDefId
				actions = append(actions, &execAction{Sql: "update request_template set proc_def_id=? where id=?", Param: []interface{}{newProDefId, v.Id}})
			}
			tmpActions := getUpdateNodeDefIdActions(v.Id, userToken)
			actions = append(actions, tmpActions...)
		}
		if err != nil {
			return
		}
		if len(actions) > 0 {
			err = transaction(actions)
		}
		if err != nil {
			return
		}
	}
	for _, v := range rowData {
		tmpObj := models.RequestTemplateQueryObj{RequestTemplateTable: *v, MGMTRoles: []*models.RoleTable{}, USERoles: []*models.RoleTable{}}
		if _, b := mgmtRoleMap[v.Id]; b {
			tmpObj.MGMTRoles = mgmtRoleMap[v.Id]
		}
		if _, b := useRoleMap[v.Id]; b {
			tmpObj.USERoles = useRoleMap[v.Id]
		}
		if v.Status == "confirm" {
			tmpObj.OperateOptions = []string{"query", "fork"}
		} else if v.Status == "created" {
			tmpObj.OperateOptions = []string{"edit", "delete"}
		}
		result = append(result, &tmpObj)
	}
	return
}

func CheckRequestTemplateRoles(requestTemplateId string, userRoles []string) error {
	rowData, err := x.QueryString("select request_template from request_template_role where request_template=? and role_type='MGMT' and `role` in ('"+strings.Join(userRoles, "','")+"')", requestTemplateId)
	if err != nil {
		return fmt.Errorf("Try to query database data fail:%s ", err.Error())
	}
	if len(rowData) == 0 {
		return fmt.Errorf(models.RowDataPermissionErr)
	}
	return nil
}

func checkProDefId(proDefId, proDefName, proDefKey, userToken string) (exist bool, newProDefId string, err error) {
	exist = false
	var processList []*models.ProcDefObj
	if proDefKey != "" {
		tmpProcessList, tmpErr := GetCoreProcessListNew(userToken)
		if tmpErr != nil {
			err = tmpErr
		} else {
			for _, v := range tmpProcessList {
				processList = append(processList, &models.ProcDefObj{ProcDefId: v.ProcDefId, ProcDefName: v.ProcDefName, ProcDefKey: v.ProcDefKey})
			}
		}
	} else {
		allProcessList, tmpErr := GetCoreProcessListAll(userToken)
		if tmpErr != nil {
			err = tmpErr
		} else {
			for _, v := range allProcessList {
				processList = append(processList, &models.ProcDefObj{ProcDefId: v.ProcDefId, ProcDefName: v.ProcDefName, ProcDefKey: v.ProcDefKey})
			}
		}
	}
	if err != nil {
		return
	}
	for _, v := range processList {
		if v.ProcDefId == proDefId {
			exist = true
			break
		}
	}
	if exist {
		return
	}
	count := 0
	for _, v := range processList {
		if proDefKey != "" {
			if proDefKey == v.ProcDefKey {
				count = count + 1
				newProDefId = v.ProcDefId
			}
			continue
		}
		if v.ProcDefName == proDefName {
			count = count + 1
			newProDefId = v.ProcDefId
		}
	}
	if count != 1 {
		if proDefKey != "" {
			err = fmt.Errorf("Find %d record from process list by query proDefKey:%s ", count, proDefKey)
		} else {
			err = fmt.Errorf("Find %d record from process list by query proDefName:%s ", count, proDefName)
		}
	}
	return
}

func getUpdateNodeDefIdActions(requestTemplateId, userToken string) (actions []*execAction) {
	actions = []*execAction{}
	var taskTemplate []*models.TaskTemplateTable
	x.SQL("select * from task_template where request_template=?", requestTemplateId).Find(&taskTemplate)
	if len(taskTemplate) == 0 {
		return actions
	}
	nodeList, _ := GetProcessNodesByProc(requestTemplateId, userToken, "template")
	nodeMap := make(map[string]string)
	for _, v := range nodeList {
		nodeMap[v.NodeId] = v.NodeDefId
	}
	for _, v := range taskTemplate {
		if v.NodeId == "" {
			continue
		}
		if nowDefId, b := nodeMap[v.NodeId]; b {
			if v.NodeDefId != nowDefId {
				actions = append(actions, &execAction{Sql: "update task_template set node_def_id=? where node_id=?", Param: []interface{}{nowDefId, v.NodeId}})
				actions = append(actions, &execAction{Sql: "update task set node_def_id=? where task_template=?", Param: []interface{}{nowDefId, v.Id}})
			}
		}
	}
	return actions
}

func getRequestTemplateIdsBySql(sql string) (ids []string, err error) {
	var requestTemplateTables []*models.RequestTemplateTable
	err = x.SQL(sql).Find(&requestTemplateTables)
	ids = []string{}
	for _, v := range requestTemplateTables {
		ids = append(ids, v.Id)
	}
	return
}

func CreateRequestTemplate(param *models.RequestTemplateUpdateParam) (result models.RequestTemplateQueryObj, err error) {
	var actions []*execAction
	newGuid := guid.CreateGuid()
	result = models.RequestTemplateQueryObj{RequestTemplateTable: param.RequestTemplateTable, MGMTRoles: []*models.RoleTable{}, USERoles: []*models.RoleTable{}}
	result.Id = newGuid
	nowTime := time.Now().Format(models.DateTimeFormat)
	insertAction := execAction{Sql: "insert into request_template(id,`group`,name,description,tags,package_name,entity_name,proc_def_key,proc_def_id,proc_def_name,expire_day,handler,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
	insertAction.Param = []interface{}{newGuid, param.Group, param.Name, param.Description, param.Tags, param.PackageName, param.EntityName, param.ProcDefKey, param.ProcDefId, param.ProcDefName, param.ExpireDay, param.Handler, param.CreatedBy, nowTime, param.CreatedBy, nowTime}
	actions = append(actions, &insertAction)
	for _, v := range param.MGMTRoles {
		result.MGMTRoles = append(result.MGMTRoles, &models.RoleTable{Id: v})
		actions = append(actions, &execAction{Sql: "insert into request_template_role(id,request_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{newGuid + models.SysTableIdConnector + v + models.SysTableIdConnector + "MGMT", newGuid, v, "MGMT"}})
	}
	for _, v := range param.USERoles {
		result.USERoles = append(result.USERoles, &models.RoleTable{Id: v})
		actions = append(actions, &execAction{Sql: "insert into request_template_role(id,request_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{newGuid + models.SysTableIdConnector + v + models.SysTableIdConnector + "USE", newGuid, v, "USE"}})
	}
	err = transaction(actions)
	return
}

func UpdateRequestTemplate(param *models.RequestTemplateUpdateParam) (result models.RequestTemplateQueryObj, err error) {
	var actions []*execAction
	nowTime := time.Now().Format(models.DateTimeFormat)
	result = models.RequestTemplateQueryObj{RequestTemplateTable: param.RequestTemplateTable, MGMTRoles: []*models.RoleTable{}, USERoles: []*models.RoleTable{}}
	updateAction := execAction{Sql: "update request_template set status='created',`group`=?,name=?,description=?,tags=?,package_name=?,entity_name=?,proc_def_key=?,proc_def_id=?,proc_def_name=?,expire_day=?,handler=?,updated_by=?,updated_time=? where id=?"}
	updateAction.Param = []interface{}{param.Group, param.Name, param.Description, param.Tags, param.PackageName, param.EntityName, param.ProcDefKey, param.ProcDefId, param.ProcDefName, param.ExpireDay, param.Handler, param.UpdatedBy, nowTime, param.Id}
	actions = append(actions, &updateAction)
	actions = append(actions, &execAction{Sql: "delete from request_template_role where request_template=?", Param: []interface{}{param.Id}})
	for _, v := range param.MGMTRoles {
		result.MGMTRoles = append(result.MGMTRoles, &models.RoleTable{Id: v})
		actions = append(actions, &execAction{Sql: "insert into request_template_role(id,request_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{param.Id + models.SysTableIdConnector + v + models.SysTableIdConnector + "MGMT", param.Id, v, "MGMT"}})
	}
	for _, v := range param.USERoles {
		result.USERoles = append(result.USERoles, &models.RoleTable{Id: v})
		actions = append(actions, &execAction{Sql: "insert into request_template_role(id,request_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{param.Id + models.SysTableIdConnector + v + models.SysTableIdConnector + "USE", param.Id, v, "USE"}})
	}
	err = transaction(actions)
	return
}

func DeleteRequestTemplate(id string) error {
	rtObj, err := getSimpleRequestTemplate(id)
	if err != nil {
		return err
	}
	if rtObj.Status == "confirm" {
		return fmt.Errorf("confirm status can not delete")
	}
	var requestTable []*models.RequestTable
	x.SQL("select id,name from request where request_template=?", id).Find(&requestTable)
	if len(requestTable) > 0 {
		return fmt.Errorf("The template used by request:%s ", requestTable[0].Name)
	}
	var taskTemplateTable []*models.TaskTemplateTable
	x.SQL("select id,form_template from task_template where request_template=?", id).Find(&taskTemplateTable)
	formTemplateIds := []string{rtObj.FormTemplate}
	for _, v := range taskTemplateTable {
		formTemplateIds = append(formTemplateIds, v.FormTemplate)
	}
	var actions []*execAction
	actions = append(actions, &execAction{Sql: "delete from task_template_role where task_template in (select id from task_template where request_template=?)", Param: []interface{}{id}})
	actions = append(actions, &execAction{Sql: "delete from task_template where request_template=?", Param: []interface{}{id}})
	actions = append(actions, &execAction{Sql: "delete from request_template_role where request_template=?", Param: []interface{}{id}})
	actions = append(actions, &execAction{Sql: "delete from request_template where id=?", Param: []interface{}{id}})
	actions = append(actions, &execAction{Sql: "delete from form_item_template where form_template in ('" + strings.Join(formTemplateIds, "','") + "')", Param: []interface{}{}})
	actions = append(actions, &execAction{Sql: "delete from form_template where id in ('" + strings.Join(formTemplateIds, "','") + "')", Param: []interface{}{}})
	return transaction(actions)
}

func ListRequestTemplateEntityAttrs(id, userToken string) (result []*models.ProcEntity, err error) {
	result = []*models.ProcEntity{}
	nodes, getNodesErr := GetProcessNodesByProc(id, userToken, "all")
	if getNodesErr != nil {
		err = getNodesErr
		return
	}
	entityMap := make(map[string]int)
	existAttrMap := make(map[string]int)
	existAttrs, _ := GetRequestTemplateEntityAttrs(id)
	for _, attr := range existAttrs {
		existAttrMap[attr.Id] = 1
	}
	for _, node := range nodes {
		for _, entity := range node.BoundEntities {
			if _, b := entityMap[entity.Id]; b {
				continue
			}
			entityMap[entity.Id] = 1
			for _, attribute := range entity.Attributes {
				attribute.Id = fmt.Sprintf("%s:%s:%s", entity.PackageName, entity.Name, attribute.Name)
				attribute.EntityId = entity.Id
				attribute.EntityName = entity.Name
				attribute.EntityDisplayName = entity.DisplayName
				attribute.EntityPackage = entity.PackageName
				if _, b := existAttrMap[attribute.Id]; b {
					attribute.Active = true
				}
			}
			result = append(result, entity)
		}
	}
	return
}

func GetRequestTemplateEntityAttrs(id string) (result []*models.ProcEntityAttributeObj, err error) {
	result = []*models.ProcEntityAttributeObj{}
	var requestTemplateTable []*models.RequestTemplateTable
	err = x.SQL("select entity_attrs from request_template where id=?", id).Find(&requestTemplateTable)
	if err != nil {
		return
	}
	if len(requestTemplateTable) == 0 {
		err = fmt.Errorf("Can not find request template wit id:%s ", id)
		return
	}
	if requestTemplateTable[0].EntityAttrs == "" {
		return
	}
	err = json.Unmarshal([]byte(requestTemplateTable[0].EntityAttrs), &result)
	if err != nil {
		err = fmt.Errorf("json unmarshal data fail:%s ", err.Error())
	}
	return
}

func UpdateRequestTemplateEntityAttrs(id string, attrs []*models.ProcEntityAttributeObj, operator string) error {
	b, _ := json.Marshal(attrs)
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err := x.Exec("update request_template set entity_attrs=?,updated_time=?,updated_by=? where id=?", string(b), nowTime, operator, id)
	return err
}

func getSimpleRequestTemplate(id string) (result models.RequestTemplateTable, err error) {
	var requestTemplateTable []*models.RequestTemplateTable
	err = x.SQL("select * from request_template where id=?", id).Find(&requestTemplateTable)
	if err != nil {
		err = fmt.Errorf("Try to query database fail,%s ", err.Error())
		return
	}
	if len(requestTemplateTable) == 0 {
		err = fmt.Errorf("Can not find request template with id:%s ", id)
		result = models.RequestTemplateTable{}
		return
	}
	result = *requestTemplateTable[0]
	return
}

func ForkConfirmRequestTemplate(requestTemplateId, operator string) error {
	requestTemplateObj, err := getSimpleRequestTemplate(requestTemplateId)
	if err != nil {
		return err
	}
	existQuery, tmpErr := x.QueryString("select id,name,version from request_template where del_flag!=1 and record_id=?", requestTemplateObj.Id)
	if tmpErr != nil {
		return fmt.Errorf("Query database fail,%s ", tmpErr.Error())
	}
	if len(existQuery) > 0 {
		return fmt.Errorf("RequestTemplate already have a branch %s:%s", existQuery[0]["name"], existQuery[0]["version"])
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	version := buildVersionNum(requestTemplateObj.Version)
	newRequestTemplateId := guid.CreateGuid()
	newRequestFormTemplateId := guid.CreateGuid()
	var actions []*execAction
	actions = append(actions, &execAction{Sql: fmt.Sprintf("insert into request_template(id,`group`,name,description,form_template,tags,status,package_name,entity_name,proc_def_key,proc_def_id,proc_def_name,created_by,created_time,updated_by,updated_time,entity_attrs,record_id,`version`,confirm_time,expire_day,handler) select '%s' as id,`group`,name,description,'%s' as form_template,tags,'created' as status,package_name,entity_name,proc_def_key,proc_def_id,proc_def_name,'%s' as created_by,'%s' as created_time,'%s' as updated_by,'%s' as updated_time,entity_attrs,'%s' as record_id,'%s' as `version`,'' as confirm_time,expire_day,handler from request_template where id='%s'", newRequestTemplateId, newRequestFormTemplateId, operator, nowTime, operator, nowTime, requestTemplateObj.Id, version, requestTemplateObj.Id)})
	newRequestFormActions, tmpErr := getFormCopyActions(requestTemplateObj.FormTemplate, newRequestFormTemplateId)
	if tmpErr != nil {
		return fmt.Errorf("Try to copy request form fail,%s ", tmpErr.Error())
	}
	actions = append(actions, newRequestFormActions...)
	var requestTemplateRoles []*models.RequestTemplateRoleTable
	x.SQL("select * from request_template_role where request_template=?", requestTemplateObj.Id).Find(&requestTemplateRoles)
	for _, v := range requestTemplateRoles {
		tmpId := newRequestTemplateId + models.SysTableIdConnector + v.Role + models.SysTableIdConnector + v.RoleType
		actions = append(actions, &execAction{Sql: "insert into request_template_role(id,request_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{tmpId, newRequestTemplateId, v.Role, v.RoleType}})
	}
	var taskTemplates []*models.TaskTemplateTable
	x.SQL("select id,form_template from task_template where request_template=?", requestTemplateObj.Id).Find(&taskTemplates)
	newTaskGuids := guid.CreateGuidList(len(taskTemplates))
	newTaskFormGuids := guid.CreateGuidList(len(taskTemplates))
	for i, task := range taskTemplates {
		actions = append(actions, &execAction{Sql: fmt.Sprintf("insert into task_template(id,name,description,form_template,request_template,node_id,node_def_id,node_name,expire_day,handler,created_by,created_time,updated_by,updated_time) select '%s' as id,name,description,'%s' as form_template,'%s' as request_template,node_id,node_def_id,node_name,expire_day,handler,created_by,created_time,updated_by,updated_time from task_template where id='%s'", newTaskGuids[i], newTaskFormGuids[i], newRequestTemplateId, task.Id)})
		tmpTaskFormActions, tmpErr := getFormCopyActions(task.FormTemplate, newTaskFormGuids[i])
		if tmpErr != nil {
			err = fmt.Errorf("Try to copy task form fail,%s ", tmpErr.Error())
			break
		}
		actions = append(actions, tmpTaskFormActions...)
		tmpTaskRoleActions, tmpErr := getTaskTemplateRoleActions(task.Id, newTaskGuids[i])
		if tmpErr != nil {
			err = fmt.Errorf("Try to copy task role releation fail,%s ", tmpErr.Error())
			break
		}
		actions = append(actions, tmpTaskRoleActions...)
	}
	if err != nil {
		return err
	}
	return transactionWithoutForeignCheck(actions)
}

func ConfirmRequestTemplate(requestTemplateId string) error {
	requestTemplateObj, err := getSimpleRequestTemplate(requestTemplateId)
	if err != nil {
		return err
	}
	if requestTemplateObj.FormTemplate == "" {
		return fmt.Errorf("Please config request template form ")
	}
	if requestTemplateObj.Status == "confirm" {
		return fmt.Errorf("Request template already confirm ")
	}
	err = validateConfirm(requestTemplateId)
	if err != nil {
		return err
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	version := requestTemplateObj.Version
	if version == "" {
		version = "v1"
	}
	var actions []*execAction
	actions = append(actions, &execAction{Sql: "update request_template set status='confirm',`version`=?,confirm_time=?,del_flag=2 where id=?", Param: []interface{}{version, nowTime, requestTemplateObj.Id}})
	return transaction(actions)
}

func getFormCopyActions(oldFormTemplateId, newFormTemplateId string) (actions []*execAction, err error) {
	var itemRows []*models.FormItemTemplateTable
	err = x.SQL("select id from form_item_template where form_template=?", oldFormTemplateId).Find(&itemRows)
	if err != nil {
		return
	}
	actions = append(actions, &execAction{Sql: fmt.Sprintf("insert into form_template(id,name,description,created_by,created_time,updated_by,updated_time) select '%s' as id,name,description,created_by,created_time,updated_by,updated_time from form_template where id='%s'", newFormTemplateId, oldFormTemplateId)})
	newGuidList := guid.CreateGuidList(len(itemRows))
	for i, item := range itemRows {
		actions = append(actions, &execAction{Sql: fmt.Sprintf("insert into form_item_template(id,form_template,name,description,default_value,sort,package_name,entity,attr_def_id,attr_def_name,attr_def_data_type,element_type,title,width,ref_package_name,ref_entity,data_options,required,regular,is_edit,is_view,is_output,item_group,item_group_name,in_display_name,is_ref_inside,multiple) select '%s' as id,'%s' as form_template,name,description,default_value,sort,package_name,entity,attr_def_id,attr_def_name,attr_def_data_type,element_type,title,width,ref_package_name,ref_entity,data_options,required,regular,is_edit,is_view,is_output,item_group,item_group_name,in_display_name,is_ref_inside,multiple from form_item_template where id='%s'", newGuidList[i], newFormTemplateId, item.Id)})
	}
	return
}

func getTaskTemplateRoleActions(oldTaskTemplateId, newTaskTemplateId string) (actions []*execAction, err error) {
	var taskTemplateRoles []*models.TaskTemplateRoleTable
	err = x.SQL("select * from task_template_role where task_template=?", oldTaskTemplateId).Find(&taskTemplateRoles)
	if err != nil {
		return
	}
	for _, v := range taskTemplateRoles {
		tmpId := newTaskTemplateId + models.SysTableIdConnector + v.Role + models.SysTableIdConnector + v.RoleType
		actions = append(actions, &execAction{Sql: "insert into task_template_role(id,task_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{tmpId, newTaskTemplateId, v.Role, v.RoleType}})
	}
	return
}

func validateConfirm(requestTemplateId string) error {
	var taskTemplateTable []*models.TaskTemplateTable
	x.SQL("select id from task_template where request_template=? and form_template IS NOT NULL", requestTemplateId).Find(&taskTemplateTable)
	if len(requestTemplateId) == 0 {
		return fmt.Errorf("Please config task template ")
	}
	return nil
}

func SetRequestTemplateToCreated(id, operator string) {
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err := x.Exec("update request_template set status='created',updated_by=?,updated_time=? where id=?", operator, nowTime, id)
	if err != nil {
		log.Logger.Error("Update request template to created status fail", log.Error(err), log.String("operator", operator), log.String("requestTemplateId", id))
	}
}

func GetRequestTemplateByUser(userRoles []string) (result []*models.UserRequestTemplateQueryObj, err error) {
	result = []*models.UserRequestTemplateQueryObj{}
	var requestTemplateTable, tmpTemplateTable []*models.RequestTemplateTable
	err = x.SQL("select * from request_template where (del_flag=2 and id in (select request_template from request_template_role where role_type='USE' and `role` in ('" + strings.Join(userRoles, "','") + "'))) or (del_flag=0 and id in (select request_template from request_template_role where role_type='MGMT' and `role` in ('" + strings.Join(userRoles, "','") + "'))) order by `group`,tags,status,id").Find(&requestTemplateTable)
	if err != nil {
		return
	}
	if len(requestTemplateTable) == 0 {
		return
	}
	recordIdMap := make(map[string]int)
	for _, v := range requestTemplateTable {
		if v.Status == "confirm" {
			if v.RecordId != "" {
				recordIdMap[v.RecordId] = 1
			}
			//v.Name = fmt.Sprintf("%s(%s)", v.Name, v.Version)
		} else {
			if v.ConfirmTime != "" {
				if !compareUpdateConfirmTime(v.UpdatedTime, v.ConfirmTime) {
					recordIdMap[v.Id] = 1
				}
			}
			v.Name = fmt.Sprintf("%s(beta)", v.Name)
			v.Version = "beta"
		}
	}
	for _, v := range requestTemplateTable {
		if _, b := recordIdMap[v.Id]; b {
			continue
		}
		tmpTemplateTable = append(tmpTemplateTable, v)
	}
	requestTemplateTable = tmpTemplateTable
	var groupTable []*models.RequestTemplateGroupTable
	x.SQL("select id,name,description from request_template_group").Find(&groupTable)
	groupMap := make(map[string]*models.RequestTemplateGroupTable)
	for _, v := range groupTable {
		groupMap[v.Id] = v
	}
	tmpGroup := requestTemplateTable[0].Group
	tmpTemplateList := []*models.RequestTemplateTable{}
	for _, v := range requestTemplateTable {
		if v.Group != tmpGroup {
			result = append(result, &models.UserRequestTemplateQueryObj{GroupId: groupMap[tmpGroup].Id, GroupName: groupMap[tmpGroup].Name, GroupDescription: groupMap[tmpGroup].Description, Templates: tmpTemplateList})
			tmpTemplateList = []*models.RequestTemplateTable{}
			tmpGroup = v.Group
		}
		tmpTemplateList = append(tmpTemplateList, v)
	}
	if len(tmpTemplateList) > 0 {
		lastGroup := requestTemplateTable[len(requestTemplateTable)-1].Group
		result = append(result, &models.UserRequestTemplateQueryObj{GroupId: groupMap[lastGroup].Id, GroupName: groupMap[lastGroup].Name, GroupDescription: groupMap[lastGroup].Description, Templates: tmpTemplateList})
	}
	for _, v := range result {
		v.Tags = groupRequestTemplateByTags(v.Templates)
	}
	return
}

func compareUpdateConfirmTime(updatedTime, confirmTime string) bool {
	// if updatedTime > confirmTime -> true
	result := false
	ut, _ := time.Parse(models.DateTimeFormat, updatedTime)
	ct, _ := time.Parse(models.DateTimeFormat, confirmTime)
	if ut.Unix() > ct.Unix() {
		result = true
	}
	return result
}

func groupRequestTemplateByTags(templates []*models.RequestTemplateTable) []*models.UserRequestTemplateTagObj {
	result := []*models.UserRequestTemplateTagObj{}
	if len(templates) == 0 {
		return result
	}
	tmpTag := templates[0].Tags
	tmpTemplateList := []*models.RequestTemplateTable{}
	for _, v := range templates {
		if v.Tags != tmpTag {
			result = append(result, &models.UserRequestTemplateTagObj{Tag: tmpTag, Templates: tmpTemplateList})
			tmpTemplateList = []*models.RequestTemplateTable{}
			tmpTag = v.Tags
		}
		tmpTemplateList = append(tmpTemplateList, v)
	}
	if len(tmpTemplateList) > 0 {
		lastTag := templates[len(templates)-1].Tags
		result = append(result, &models.UserRequestTemplateTagObj{Tag: lastTag, Templates: tmpTemplateList})
	}
	return result
}

func getRequestTemplateRoles(requestTemplateId, roleType string) []string {
	result := []string{}
	var rtRoles []*models.RequestTemplateRoleTable
	x.SQL("select `role` from request_template_role where request_template=? and role_type=?", requestTemplateId, roleType).Find(&rtRoles)
	for _, v := range rtRoles {
		result = append(result, v.Role)
	}
	return result
}

func QueryUserByRoles(roles []string, userToken string) (result []string, err error) {
	result = []string{}
	var roleTable []*models.RoleTable
	err = x.SQL("select * from `role` where id in ('" + strings.Join(roles, "','") + "')").Find(&roleTable)
	if err != nil || len(roleTable) == 0 {
		return
	}
	existMap := make(map[string]int)
	for _, v := range roleTable {
		if v.CoreId == "" {
			continue
		}
		tmpResult, tmpErr := queryCoreUser(v.CoreId, userToken)
		if tmpErr != nil {
			err = tmpErr
			break
		}
		for _, vv := range tmpResult {
			if _, b := existMap[vv]; !b {
				result = append(result, vv)
				existMap[vv] = 1
			}
		}
	}
	return
}

func queryCoreUser(roleId, userToken string) (result []string, err error) {
	request, newRequestErr := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/platform/v1/roles/%s/users", models.CoreToken.BaseUrl, roleId), strings.NewReader(""))
	if newRequestErr != nil {
		err = fmt.Errorf("Get core role key new request fail:%s ", newRequestErr.Error())
		return
	}
	request.Header.Set("Authorization", userToken)
	res, requestErr := http.DefaultClient.Do(request)
	if requestErr != nil {
		err = fmt.Errorf("Get core user request fail:%s ", requestErr.Error())
		return
	}
	b, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()
	var response models.CoreUserDto
	err = json.Unmarshal(b, &response)
	if err != nil {
		err = fmt.Errorf("Get core user json unmarshal result:%s ", err.Error())
		return
	}
	if response.Status != "OK" {
		err = fmt.Errorf(response.Message)
		return
	}
	for _, v := range response.Data {
		result = append(result, v.Username)
	}
	return
}

func GetRequestTemplateTags(group string) (result []string, err error) {
	result = []string{}
	var requestTemplates []*models.RequestTemplateTable
	err = x.SQL("select distinct tags from request_template where `group`=?", group).Find(&requestTemplates)
	for _, v := range requestTemplates {
		if v.Tags == "" {
			continue
		}
		result = append(result, v.Tags)
	}
	return
}
