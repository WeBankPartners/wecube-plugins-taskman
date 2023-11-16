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
	userRoleFilterSql, userRoleFilterParams := createListParams(userRoles, "")
	baseSql := fmt.Sprintf("SELECT %s FROM request_template_group WHERE manage_role in ("+userRoleFilterSql+") and del_flag=0 %s ", queryColumn, filterSql)
	queryParam = append(userRoleFilterParams, queryParam...)
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
	rolesFilterSql, rolesFilterParam := createListParams(roles, "")
	var requestTemplateGroupRows []*models.RequestTemplateGroupTable
	rolesFilterParam = append([]interface{}{id}, rolesFilterParam...)
	err := x.SQL("select id from request_template_group where id=? and manage_role in ("+rolesFilterSql+")", rolesFilterParam...).Find(&requestTemplateGroupRows)
	if err != nil {
		return fmt.Errorf("Try to query database data fail,%s ", err.Error())
	}
	if len(requestTemplateGroupRows) == 0 {
		return fmt.Errorf(models.RowDataPermissionErr)
	}
	return nil
}

func GetCoreProcessListNew(userToken string) (processList []*models.ProcDefObj, err error) {
	processList = []*models.ProcDefObj{}
	req, reqErr := http.NewRequest(http.MethodGet, models.Config.Wecube.BaseUrl+"/platform/v1/public/process/definitions?tags="+models.ProcessFetchTabs, nil)
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

func GetProcessNodesByProc(requestTemplateObj models.RequestTemplateTable, userToken string, filterType string) (nodeList models.ProcNodeObjList, err error) {
	if requestTemplateObj.ProcDefId == "" {
		requestTemplateObj, err = getSimpleRequestTemplate(requestTemplateObj.Id)
		if err != nil {
			return
		}
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
	log.Logger.Debug("platform process task node return", log.String("body", string(respBytes)))
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

func getRoleMail(roleList []*models.RoleTable) (mailList []string) {
	if models.CoreToken.BaseUrl == "" || len(roleList) == 0 {
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
	var result models.CoreRoleDto
	err = json.Unmarshal(b, &result)
	if err != nil {
		log.Logger.Error("Get core role key json unmarshal result", log.Error(err))
		return
	}
	for _, v := range roleList {
		tmpMail := ""
		for _, vv := range result.Data {
			if vv.Name == v.Id {
				tmpMail = vv.Email
				break
			}
		}
		if tmpMail != "" {
			mailList = append(mailList, tmpMail)
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
			addRoleList = append(addRoleList, &models.RoleTable{Id: v.Name, DisplayName: v.DisplayName, CoreId: v.Id, Email: v.Email})
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
		actions = append(actions, &execAction{Sql: "insert into `role`(id,display_name,core_id,email,updated_time) value (?,?,?,?,NOW())", Param: []interface{}{role.Id, role.DisplayName, role.CoreId, role.Email}})
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
		idFilterSql, idFilterParam := createListParams(ids, "")
		err = x.SQL("select * from role where id in ("+idFilterSql+")", idFilterParam...).Find(&result)
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
				roleFilterSql, roleFilterParam := createListParams(inValueStringList, "")
				if v.Name == "mgmtRoles" {
					tmpIds, tmpErr = getRequestTemplateIdsBySql("select t1.id from request_template t1 left join request_template_role t2 on t1.id=t2.request_template where t2.role_type='MGMT' and t2.role in ("+roleFilterSql+")", roleFilterParam)
				} else {
					tmpIds, tmpErr = getRequestTemplateIdsBySql("select t1.id from request_template t1 left join request_template_role t2 on t1.id=t2.request_template where t2.role_type='USE' and t2.role in ("+roleFilterSql+")", roleFilterParam)
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
	userRolesFilterSql, userRolesFilterParam := createListParams(userRoles, "")
	queryParam = append(userRolesFilterParam, queryParam...)
	baseSql := fmt.Sprintf("SELECT %s FROM (select * from request_template where del_flag=0 or (del_flag=2 and id not in (select record_id from request_template where del_flag=2 and record_id<>''))) t1 WHERE t1.id in (select request_template from request_template_role where role_type='MGMT' and `role` in ("+userRolesFilterSql+")) %s %s ", queryColumn, extFilterSql, filterSql)
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
		for _, v := range rowData {
			tmpErr := SyncProcDefId(v.Id, v.ProcDefId, v.ProcDefName, v.ProcDefKey, userToken)
			if tmpErr != nil {
				err = fmt.Errorf("Try to sync proDefId fail,%s ", tmpErr.Error())
				break
			}
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
			tmpObj.OperateOptions = []string{"query", "fork", "export", "disable"}
		} else if v.Status == "created" {
			tmpObj.OperateOptions = []string{"edit", "delete"}
		} else if v.Status == "disable" {
			tmpObj.OperateOptions = []string{"query", "enable"}
		}
		result = append(result, &tmpObj)
	}
	return
}

func CheckRequestTemplateRoles(requestTemplateId string, userRoles []string) error {
	var requestTemplateRoleRows []*models.RequestTemplateRoleTable
	userRolesFilterSql, userRolesFilterParam := createListParams(userRoles, "")
	userRolesFilterParam = append([]interface{}{requestTemplateId}, userRolesFilterParam...)
	err := x.SQL("select request_template from request_template_role where request_template=? and role_type='MGMT' and `role` in ("+userRolesFilterSql+")", userRolesFilterParam...).Find(&requestTemplateRoleRows)
	if err != nil {
		return fmt.Errorf("Try to query database data fail:%s ", err.Error())
	}
	if len(requestTemplateRoleRows) == 0 {
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
	nodeList, _ := GetProcessNodesByProc(models.RequestTemplateTable{Id: requestTemplateId}, userToken, "template")
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

func SyncProcDefId(requestTemplateId, proDefId, proDefName, proDefKey, userToken string) error {
	log.Logger.Info("Start sync process def id")
	proExistFlag, newProDefId, err := checkProDefId(proDefId, proDefName, proDefKey, userToken)
	if err != nil {
		return err
	}
	var actions []*execAction
	if !proExistFlag {
		if proDefKey != "" {
			actions = append(actions, &execAction{Sql: "update request_template set proc_def_id=? where id=?", Param: []interface{}{newProDefId, requestTemplateId}})
		} else {
			actions = append(actions, &execAction{Sql: "update request_template set proc_def_id=? where proc_def_name=?", Param: []interface{}{newProDefId, proDefName}})
		}
		err = transaction(actions)
		if err != nil {
			return fmt.Errorf("Update requestTemplate procDefId fail,%s ", err.Error())
		}
		log.Logger.Info("Update requestTemplate proDefId done")
		actions = []*execAction{}
	}
	tmpActions := getUpdateNodeDefIdActions(requestTemplateId, userToken)
	if len(tmpActions) > 0 {
		err = transaction(tmpActions)
		if err != nil {
			return fmt.Errorf("Update template node def id fail,%s ", err.Error())
		}
		log.Logger.Info("Update taskTemplate nodeDefId done")
	}
	return nil
}

func getRequestTemplateIdsBySql(sql string, param []interface{}) (ids []string, err error) {
	var requestTemplateTables []*models.RequestTemplateTable
	err = x.SQL(sql, param...).Find(&requestTemplateTables)
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
	insertAction := execAction{Sql: "insert into request_template(id,`group`,name,description,tags,package_name,entity_name,proc_def_key,proc_def_id,proc_def_name,expire_day,handler,created_by,created_time,updated_by,updated_time,type) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
	insertAction.Param = []interface{}{newGuid, param.Group, param.Name, param.Description, param.Tags, param.PackageName, param.EntityName, param.ProcDefKey, param.ProcDefId, param.ProcDefName, param.ExpireDay, param.Handler, param.CreatedBy, nowTime, param.CreatedBy, nowTime, param.Type}
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

func DeleteRequestTemplate(id string, getActionFlag bool) (actions []*execAction, err error) {
	rtObj, err := getSimpleRequestTemplate(id)
	if err != nil {
		return actions, err
	}
	if rtObj.Status == "confirm" {
		return actions, fmt.Errorf("confirm status can not delete")
	}
	var taskTemplateTable []*models.TaskTemplateTable
	x.SQL("select id,form_template from task_template where request_template=?", id).Find(&taskTemplateTable)
	formTemplateIds := []string{rtObj.FormTemplate}
	for _, v := range taskTemplateTable {
		formTemplateIds = append(formTemplateIds, v.FormTemplate)
	}
	actions = []*execAction{}
	var requestTable []*models.RequestTable
	x.SQL("select id,name from request where request_template=?", id).Find(&requestTable)
	if len(requestTable) > 0 {
		var formTable []*models.FormTable
		x.SQL("select id from form where form_template in ('" + strings.Join(formTemplateIds, "','") + "')").Find(&formTable)
		formIds := []string{}
		for _, v := range formTable {
			formIds = append(formIds, v.Id)
		}
		actions = append(actions, &execAction{Sql: "delete from operation_log where task in (select id from task where task_template in (select id from task_template where request_template=?))", Param: []interface{}{id}})
		actions = append(actions, &execAction{Sql: "delete from operation_log where request in (select id from request where request_template=?)", Param: []interface{}{id}})
		actions = append(actions, &execAction{Sql: "delete from task where task_template in (select id from task_template where request_template=?)", Param: []interface{}{id}})
		actions = append(actions, &execAction{Sql: "delete from request where request_template=?", Param: []interface{}{id}})
		actions = append(actions, &execAction{Sql: "delete from form_item where form in ('" + strings.Join(formIds, "','") + "')", Param: []interface{}{}})
		actions = append(actions, &execAction{Sql: "delete from form where form_template in ('" + strings.Join(formTemplateIds, "','") + "')", Param: []interface{}{}})
	}
	actions = append(actions, &execAction{Sql: "delete from task_template_role where task_template in (select id from task_template where request_template=?)", Param: []interface{}{id}})
	actions = append(actions, &execAction{Sql: "delete from task_template where request_template=?", Param: []interface{}{id}})
	actions = append(actions, &execAction{Sql: "delete from request_template_role where request_template=?", Param: []interface{}{id}})
	actions = append(actions, &execAction{Sql: "delete from request_template where id=?", Param: []interface{}{id}})
	actions = append(actions, &execAction{Sql: "delete from form_item_template where form_template in ('" + strings.Join(formTemplateIds, "','") + "')", Param: []interface{}{}})
	actions = append(actions, &execAction{Sql: "delete from form_template where id in ('" + strings.Join(formTemplateIds, "','") + "')", Param: []interface{}{}})
	if !getActionFlag {
		err = transaction(actions)
	}
	return actions, err
}

func ListRequestTemplateEntityAttrs(id, userToken string) (result []*models.ProcEntity, err error) {
	result = []*models.ProcEntity{}
	nodes, getNodesErr := GetProcessNodesByProc(models.RequestTemplateTable{Id: id}, userToken, "all")
	if getNodesErr != nil {
		err = getNodesErr
		return
	}
	if len(nodes) == 0 {
		return
	}
	entityMap := make(map[string]int)
	existAttrMap := make(map[string]int)
	existAttrs, _ := GetRequestTemplateEntityAttrs(id)
	for _, attr := range existAttrs {
		existAttrMap[attr.Id] = 1
	}
	for _, node := range nodes {
		if node == nil {
			continue
		}
		for _, entity := range node.BoundEntities {
			if entity == nil {
				continue
			}
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

func getAllRequestTemplate() (templateMap map[string]*models.RequestTemplateTable, err error) {
	templateMap = make(map[string]*models.RequestTemplateTable)
	var requestTemplateTable []*models.RequestTemplateTable
	err = x.SQL("select * from request_template").Find(&requestTemplateTable)
	if err != nil {
		err = fmt.Errorf("Try to query database fail,%s ", err.Error())
		return
	}
	for _, template := range requestTemplateTable {
		templateMap[template.Id] = template
	}
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

func GetRequestTemplateByUser(user string, userRoles []string) (result []*models.UserRequestTemplateQueryObjNew, err error) {
	var roleTemplateGroupMap = make(map[string]map[string][]*models.RequestTemplateTableObj)
	var groupIdMap = make(map[string]string)
	result = []*models.UserRequestTemplateQueryObjNew{}
	var requestTemplateTable, tmpTemplateTable []*models.RequestTemplateTable
	userRolesFilterSql, userRolesFilterParam := createListParams(userRoles, "")
	queryParam := append(userRolesFilterParam, userRolesFilterParam...)
	err = x.SQL("select * from request_template where (del_flag=2 and id in (select request_template from request_template_role where role_type='USE' and `role` in ("+userRolesFilterSql+"))) or (del_flag=0 and id in (select request_template from request_template_role where role_type='MGMT' and `role` in ("+userRolesFilterSql+"))) order by `group`,tags,status,id", queryParam...).Find(&requestTemplateTable)
	if err != nil {
		return
	}
	if len(requestTemplateTable) == 0 {
		return
	}
	recordIdMap := make(map[string]int)
	disableNameMap := make(map[string]int)
	for _, v := range requestTemplateTable {
		if v.Status == "disable" {
			if v.Version == "" {
				disableNameMap[v.Name] = 1
			} else {
				tmpV, _ := strconv.Atoi(v.Version[1:])
				disableNameMap[v.Name] = tmpV
			}
		}
	}
	for _, v := range requestTemplateTable {
		if disVersion, isDisable := disableNameMap[v.Name]; isDisable {
			if v.Version == "" {
				continue
			}
			tmpV, _ := strconv.Atoi(v.Version[1:])
			if tmpV <= disVersion {
				continue
			}
		}
		if v.Status == "confirm" {
			if v.RecordId != "" {
				recordIdMap[v.RecordId] = 1
			}
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
		if disVersion, isDisable := disableNameMap[v.Name]; isDisable {
			if v.Version == "" {
				continue
			}
			tmpV, _ := strconv.Atoi(v.Version[1:])
			if tmpV <= disVersion {
				continue
			}
		}
		tmpTemplateTable = append(tmpTemplateTable, v)
		groupIdMap[v.Group] = "groupName"
	}
	requestTemplateTable = tmpTemplateTable
	// 查询当前用户所有收藏模板记录
	collectMap, _ := QueryAllTemplateCollect(user)
	var collectFlag int
	var allGroupTable []*models.RequestTemplateGroupTable
	err = x.SQL("select id,name,description,manage_role from request_template_group").Find(&allGroupTable)
	if err != nil {
		return
	}
	for _, group := range allGroupTable {
		if groupIdMap[group.Id] != "" {
			// 展示的 group
			if roleTemplateGroupMap[group.ManageRole] == nil {
				roleTemplateGroupMap[group.ManageRole] = make(map[string][]*models.RequestTemplateTableObj)
			}
			roleTemplateGroupMap[group.ManageRole][group.Id] = make([]*models.RequestTemplateTableObj, 0)
			// 设置 groupIdMap value 为group name
			groupIdMap[group.Id] = group.Name
		}
	}
	// 组装数据
	if len(roleTemplateGroupMap) > 0 && len(requestTemplateTable) > 0 {
		for _, template := range requestTemplateTable {
			collectFlag = 0
			if len(collectMap) > 0 && collectMap[template.Id] {
				collectFlag = 1
			}
			for _, roleGroupMap := range roleTemplateGroupMap {
				for groupId, templateArr := range roleGroupMap {
					if template.Group == groupId {
						templateArr = append(templateArr, &models.RequestTemplateTableObj{
							Id:          template.Id,
							Name:        template.Name,
							Tags:        template.Tags,
							Status:      template.Status,
							UpdatedBy:   template.UpdatedBy,
							UpdatedTime: template.UpdatedTime,
							CollectFlag: collectFlag,
						})
					}
					roleGroupMap[groupId] = templateArr
				}
			}
		}
	}
	for role, roleGroupMap := range roleTemplateGroupMap {
		groups := make([]*models.TemplateGroupObj, 0)
		for groupId, templateArr := range roleGroupMap {
			groups = append(groups, &models.TemplateGroupObj{
				GroupId:   groupId,
				GroupName: groupIdMap[groupId],
				Templates: templateArr,
			})
		}
		result = append(result, &models.UserRequestTemplateQueryObjNew{
			ManageRole: role,
			Groups:     groups,
		})
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
	rolesFilterSql, rolesFilterParam := createListParams(roles, "")
	err = x.SQL("select * from `role` where id in ("+rolesFilterSql+")", rolesFilterParam...).Find(&roleTable)
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

func RequestTemplateExport(requestTemplateId string) (result models.RequestTemplateExport, err error) {
	var requestTemplateTable []*models.RequestTemplateTable
	result.RequestTemplateRole = []*models.RequestTemplateRoleTable{}
	result.TaskTemplate = []*models.TaskTemplateTable{}
	result.TaskTemplateRole = []*models.TaskTemplateRoleTable{}
	result.FormTemplate = []*models.FormTemplateTable{}
	result.FormItemTemplate = []*models.FormItemTemplateTable{}
	err = x.SQL("select * from request_template where id=?", requestTemplateId).Find(&requestTemplateTable)
	if err != nil {
		return
	}
	if len(requestTemplateTable) == 0 {
		err = fmt.Errorf("Can not find requestTemplate with id:%s ", requestTemplateId)
		return
	}
	result.RequestTemplate = *requestTemplateTable[0]
	x.SQL("select * from request_template_role where request_template=?", requestTemplateId).Find(&result.RequestTemplateRole)
	x.SQL("select * from task_template where request_template=?", requestTemplateId).Find(&result.TaskTemplate)
	x.SQL("select * from task_template_role where task_template in (select id from task_template where request_template=?)", requestTemplateId).Find(&result.TaskTemplateRole)
	x.SQL("select * from form_template where id in (select form_template from request_template where id=? union select form_template from task_template where request_template=?)", requestTemplateId, requestTemplateId).Find(&result.FormTemplate)
	x.SQL("select * from form_item_template where form_template in (select id from form_template where id in (select form_template from request_template where id=? union select form_template from task_template where request_template=?))", requestTemplateId, requestTemplateId).Find(&result.FormItemTemplate)
	var requestTemplateGroupTable []*models.RequestTemplateGroupTable
	x.SQL("select * from request_template_group where id=?", result.RequestTemplate.Group).Find(&requestTemplateGroupTable)
	if len(requestTemplateGroupTable) > 0 {
		result.RequestTemplateGroup = *requestTemplateGroupTable[0]
	}
	return
}

func RequestTemplateImport(input models.RequestTemplateExport, userToken, confirmToken, operator string) (backToken string, err error) {
	var actions []*execAction
	if confirmToken == "" {
		existFlag, err := checkImportExist(input.RequestTemplate.Id)
		if err != nil {
			return backToken, err
		}
		if existFlag {
			backToken = guid.CreateGuid()
			models.RequestTemplateImportMap[backToken] = input
			return backToken, nil
		}
	} else {
		if inputCache, b := models.RequestTemplateImportMap[confirmToken]; b {
			input = inputCache
		} else {
			return backToken, fmt.Errorf("Fetch input cache fail,please refersh and try again ")
		}
		delete(models.RequestTemplateImportMap, confirmToken)
		delActions, delErr := DeleteRequestTemplate(input.RequestTemplate.Id, true)
		if delErr != nil {
			return backToken, delErr
		}
		actions = append(actions, delActions...)
	}
	if input.RequestTemplate.Id == "" {
		return backToken, fmt.Errorf("RequestTemplate id illegal ")
	}
	allProcessList, processErr := GetCoreProcessListAll(userToken)
	if processErr != nil {
		return backToken, fmt.Errorf("Get core process list fail,%s ", processErr.Error())
	}
	processExistFlag := false
	for _, v := range allProcessList {
		if v.ProcDefName == input.RequestTemplate.ProcDefName {
			processExistFlag = true
			input.RequestTemplate.ProcDefId = v.ProcDefId
			input.RequestTemplate.ProcDefKey = v.ProcDefKey
		}
	}
	if !processExistFlag {
		return backToken, fmt.Errorf("Reqeust process:%s can not find! ", input.RequestTemplate.ProcDefName)
	}
	nodeList, _ := GetProcessNodesByProc(input.RequestTemplate, userToken, "template")
	for i, v := range input.TaskTemplate {
		existFlag := false
		for _, node := range nodeList {
			if v.NodeId == node.NodeId {
				existFlag = true
				input.TaskTemplate[i].NodeDefId = node.NodeDefId
				break
			}
		}
		if !existFlag {
			err = fmt.Errorf("Node:%s can not find in exist process:%s ", v.NodeName, input.RequestTemplate.ProcDefName)
			break
		}
	}
	if err != nil {
		return backToken, err
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	for _, v := range input.FormTemplate {
		actions = append(actions, &execAction{Sql: "insert into form_template(id,name,description,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?)", Param: []interface{}{v.Id, v.Name, v.Description, operator, nowTime, operator, nowTime}})
	}
	for _, v := range input.FormItemTemplate {
		tmpAction := execAction{Sql: "insert into form_item_template(id,form_template,name,description,item_group,item_group_name,default_value,sort,package_name,entity,attr_def_id,attr_def_name,attr_def_data_type,element_type,title,width,ref_package_name,ref_entity,data_options,required,regular,is_edit,is_view,is_output,in_display_name,is_ref_inside,multiple) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
		tmpAction.Param = []interface{}{v.Id, v.FormTemplate, v.Name, v.Description, v.ItemGroup, v.ItemGroupName, v.DefaultValue, v.Sort, v.PackageName, v.Entity, v.AttrDefId, v.AttrDefName, v.AttrDefDataType, v.ElementType, v.Title, v.Width, v.RefPackageName, v.RefEntity, v.DataOptions, v.Required, v.Regular, v.IsEdit, v.IsView, v.IsOutput, v.InDisplayName, v.IsRefInside, v.Multiple}
		actions = append(actions, &tmpAction)
	}
	var roleTable []*models.RoleTable
	x.SQL("select id from `role`").Find(&roleTable)
	roleMap := make(map[string]int)
	for _, v := range roleTable {
		roleMap[v.Id] = 1
	}
	var requestTemplateGroupTable []*models.RequestTemplateGroupTable
	x.SQL("select id from request_template_group where id=?", input.RequestTemplate.Group).Find(&requestTemplateGroupTable)
	if len(requestTemplateGroupTable) == 0 {
		if _, b := roleMap[input.RequestTemplateGroup.ManageRole]; !b {
			input.RequestTemplateGroup.ManageRole = models.AdminRole
		}
		actions = append(actions, &execAction{Sql: "insert into request_template_group(id,name,description,manage_role,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?,?)", Param: []interface{}{input.RequestTemplateGroup.Id, input.RequestTemplateGroup.Name, input.RequestTemplateGroup.Description, input.RequestTemplateGroup.ManageRole, operator, nowTime, operator, nowTime}})
	}
	input.RequestTemplate.Status = "created"
	input.RequestTemplate.ConfirmTime = ""
	rtAction := execAction{Sql: "insert into request_template(id,`group`,name,description,form_template,tags,record_id,`version`,confirm_time,status,package_name,entity_name,proc_def_key,proc_def_id,proc_def_name,expire_day,created_by,created_time,updated_by,updated_time,entity_attrs,handler) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
	rtAction.Param = []interface{}{input.RequestTemplate.Id, input.RequestTemplate.Group, input.RequestTemplate.Name, input.RequestTemplate.Description, input.RequestTemplate.FormTemplate, input.RequestTemplate.Tags, input.RequestTemplate.RecordId, input.RequestTemplate.Version, input.RequestTemplate.ConfirmTime, input.RequestTemplate.Status, input.RequestTemplate.PackageName, input.RequestTemplate.EntityName,
		input.RequestTemplate.ProcDefKey, input.RequestTemplate.ProcDefId, input.RequestTemplate.ProcDefName, input.RequestTemplate.ExpireDay, operator, nowTime, operator, nowTime, input.RequestTemplate.EntityAttrs, input.RequestTemplate.Handler}
	actions = append(actions, &rtAction)
	for _, v := range input.TaskTemplate {
		tmpAction := execAction{Sql: "insert into task_template(id,name,description,form_template,request_template,node_id,node_def_id,node_name,expire_day,handler,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
		tmpAction.Param = []interface{}{v.Id, v.Name, v.Description, v.FormTemplate, v.RequestTemplate, v.NodeId, v.NodeDefId, v.NodeName, v.ExpireDay, v.Handler, operator, nowTime, operator, nowTime}
		actions = append(actions, &tmpAction)
	}
	rtRoleFetch := false
	for _, v := range input.RequestTemplateRole {
		if _, b := roleMap[v.Role]; b {
			rtRoleFetch = true
			actions = append(actions, &execAction{Sql: "insert into request_template_role(id,request_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{v.Id, v.RequestTemplate, v.Role, v.RoleType}})
		}
	}
	if !rtRoleFetch {
		actions = append(actions, &execAction{Sql: "insert into request_template_role(id,request_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{guid.CreateGuid() + models.SysTableIdConnector + models.AdminRole + models.SysTableIdConnector + "MGMT", input.RequestTemplate.Id, models.AdminRole, "MGMT"}})
	}
	for _, v := range input.TaskTemplateRole {
		if _, b := roleMap[v.Role]; b {
			actions = append(actions, &execAction{Sql: "insert into task_template_role(id,task_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{v.Id, v.TaskTemplate, v.Role, v.RoleType}})
		}
	}
	return backToken, transaction(actions)
}

func checkImportExist(requestTemplateId string) (exist bool, err error) {
	exist = false
	var requestTemplateTable []*models.RequestTemplateTable
	x.SQL("select id,name,version,status from request_template where id=?", requestTemplateId).Find(&requestTemplateTable)
	if len(requestTemplateTable) == 0 {
		return
	}
	exist = true
	if requestTemplateTable[0].Status == "confirm" {
		err = fmt.Errorf("RequestTemplate:%s %s already confirm ", requestTemplateTable[0].Name, requestTemplateTable[0].Version)
	}
	return
}

func DisableRequestTemplate(requestTemplateId, operator string) (err error) {
	queryRows, queryErr := x.QueryString("select status from request_template where id=?", requestTemplateId)
	if queryErr != nil {
		err = queryErr
		return
	}
	if len(queryRows) == 0 {
		err = fmt.Errorf("can not find template with id: %s ", requestTemplateId)
		return
	}
	if queryRows[0]["status"] != "confirm" {
		err = fmt.Errorf("only confirm status template can disable")
		return
	}
	_, err = x.Exec("update request_template set status='disable' where id=?", requestTemplateId)
	return
}

func EnableRequestTemplate(requestTemplateId, operator string) (err error) {
	queryRows, queryErr := x.QueryString("select status from request_template where id=?", requestTemplateId)
	if queryErr != nil {
		err = queryErr
		return
	}
	if len(queryRows) == 0 {
		err = fmt.Errorf("can not find template with id: %s ", requestTemplateId)
		return
	}
	if queryRows[0]["status"] != "disable" {
		err = fmt.Errorf("only disable status template can enable")
		return
	}
	_, err = x.Exec("update request_template set status='confirm' where id=?", requestTemplateId)
	return
}
