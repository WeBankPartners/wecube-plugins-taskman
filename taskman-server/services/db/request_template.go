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

func QueryRequestTemplateGroup(param *models.QueryRequestParam) (pageInfo models.PageInfo, rowData []*models.RequestTemplateGroupTable, err error) {
	rowData = []*models.RequestTemplateGroupTable{}
	filterSql, queryColumn, queryParam := transFiltersToSQL(param, &models.TransFiltersParam{IsStruct: true, StructObj: models.RequestTemplateGroupTable{}, PrimaryKey: "id"})
	baseSql := fmt.Sprintf("SELECT %s FROM request_template_group WHERE del_flag=0 %s ", queryColumn, filterSql)
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

func GetProcessNodesByProc(requestTemplateId, userToken string) (nodeList models.ProcNodeObjList, err error) {
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
	log.Logger.Debug("Get process node list", log.String("body", string(respBytes)))
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	if respObj.Status != "OK" {
		err = fmt.Errorf(respObj.Message)
		return
	}
	nodeList = respObj.Data
	for _, v := range nodeList {
		if v.OrderedNo == "" {
			v.OrderedNum = 0
			continue
		}
		v.OrderedNum, _ = strconv.Atoi(v.OrderedNo)
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
			addRoleList = append(addRoleList, &models.RoleTable{Id: v.Name, DisplayName: v.DisplayName})
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
		actions = append(actions, &execAction{Sql: "insert into `role`(id,display_name,updated_time) value (?,?,NOW())", Param: []interface{}{role.Id, role.DisplayName}})
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

func QueryRequestTemplate(param *models.QueryRequestParam) (pageInfo models.PageInfo, result []*models.RequestTemplateQueryObj, err error) {
	extFilterSql := ""
	result = []*models.RequestTemplateQueryObj{}
	if len(param.Filters) > 0 {
		newFilters := []*models.QueryRequestFilterObj{}
		for _, v := range param.Filters {
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
	filterSql, queryColumn, queryParam := transFiltersToSQL(param, &models.TransFiltersParam{IsStruct: true, StructObj: models.RequestTemplateTable{}, PrimaryKey: "id"})
	baseSql := fmt.Sprintf("SELECT %s FROM request_template WHERE del_flag=0 %s %s ", queryColumn, extFilterSql, filterSql)
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
	for _, v := range rowData {
		tmpObj := models.RequestTemplateQueryObj{RequestTemplateTable: *v, MGMTRoles: []*models.RoleTable{}, USERoles: []*models.RoleTable{}}
		if _, b := mgmtRoleMap[v.Id]; b {
			tmpObj.MGMTRoles = mgmtRoleMap[v.Id]
		}
		if _, b := useRoleMap[v.Id]; b {
			tmpObj.USERoles = useRoleMap[v.Id]
		}
		result = append(result, &tmpObj)
	}
	return
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
	insertAction := execAction{Sql: "insert into request_template(id,`group`,name,description,tags,package_name,entity_name,proc_def_key,proc_def_id,proc_def_name,expire_day,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
	insertAction.Param = []interface{}{newGuid, param.Group, param.Name, param.Description, param.Tags, param.PackageName, param.EntityName, param.ProcDefKey, param.ProcDefId, param.ProcDefName, param.ExpireDay, param.CreatedBy, nowTime, param.CreatedBy, nowTime}
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
	updateAction := execAction{Sql: "update request_template set status='created',`group`=?,name=?,description=?,tags=?,package_name=?,entity_name=?,proc_def_key=?,proc_def_id=?,proc_def_name=?,updated_by=?,updated_time=? where id=?"}
	updateAction.Param = []interface{}{param.Group, param.Name, param.Description, param.Tags, param.PackageName, param.EntityName, param.ProcDefKey, param.ProcDefId, param.ProcDefName, param.UpdatedBy, nowTime, param.Id}
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
	_, err := x.Exec("update request_template set del_flag=1 where id=?", id)
	return err
}

func ListRequestTemplateEntityAttrs(id, userToken string) (result []*models.ProcEntity, err error) {
	result = []*models.ProcEntity{}
	nodes, getNodesErr := GetProcessNodesByProc(id, userToken)
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

func UpdateRequestTemplateEntityAttrs(id string, attrs []*models.ProcEntityAttributeObj) error {
	b, _ := json.Marshal(attrs)
	_, err := x.Exec("update request_template set entity_attrs=? where id=?", string(b), id)
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

func ConfirmRequestTemplate(requestTemplateId string) error {
	requestTemplateObj, err := getSimpleRequestTemplate(requestTemplateId)
	if err != nil {
		return err
	}
	if requestTemplateObj.FormTemplate == "" {
		return fmt.Errorf("Please config request template form ")
	}
	err = validateConfirm(requestTemplateId)
	if err != nil {
		return err
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	version := buildVersionNum(requestTemplateObj.Version)
	newRequestTemplateId := guid.CreateGuid()
	newRequestFormTemplateId := guid.CreateGuid()
	var actions []*execAction
	actions = append(actions, &execAction{Sql: "update request_template set status='confirm',`version`=?,confirm_time=?,del_flag=2 where id=?", Param: []interface{}{version, nowTime, requestTemplateObj.Id}})
	actions = append(actions, &execAction{Sql: fmt.Sprintf("insert into request_template(id,`group`,name,description,form_template,tags,status,package_name,entity_name,proc_def_key,proc_def_id,proc_def_name,created_by,created_time,updated_by,updated_time,entity_attrs,record_id,`version`,confirm_time) select '%s' as id,`group`,name,description,'%s' as form_template,tags,'confirm' as status,package_name,entity_name,proc_def_key,proc_def_id,proc_def_name,created_by,created_time,updated_by,updated_time,entity_attrs,'%s' as record_id,'%s' as `version`,'%s' as confirm_time from request_template where id='%s'", newRequestTemplateId, newRequestFormTemplateId, requestTemplateObj.Id, version, nowTime, requestTemplateObj.Id)})
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
		actions = append(actions, &execAction{Sql: fmt.Sprintf("insert into task_template(id,name,description,form_template,request_template,node_def_id,node_name,created_by,created_time,updated_by,updated_time) select '%s' as id,name,description,'%s' as form_template,'%s' as request_template,node_def_id,node_name,created_by,created_time,updated_by,updated_time from task_template where id='%s'", newTaskGuids[i], newTaskFormGuids[i], newRequestTemplateId, task.Id)})
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

func getFormCopyActions(oldFormTemplateId, newFormTemplateId string) (actions []*execAction, err error) {
	var itemRows []*models.FormItemTemplateTable
	err = x.SQL("select id from form_item_template where form_template=?", oldFormTemplateId).Find(&itemRows)
	if err != nil {
		return
	}
	actions = append(actions, &execAction{Sql: fmt.Sprintf("insert into form_template(id,name,description,created_by,created_time,updated_by,updated_time) select '%s' as id,name,description,created_by,created_time,updated_by,updated_time from form_template where id='%s'", newFormTemplateId, oldFormTemplateId)})
	newGuidList := guid.CreateGuidList(len(itemRows))
	for i, item := range itemRows {
		actions = append(actions, &execAction{Sql: fmt.Sprintf("insert into form_item_template(id,form_template,name,description,default_value,sort,package_name,entity,attr_def_id,attr_def_name,attr_def_data_type,element_type,title,width,ref_package_name,ref_entity,data_options,required,regular,is_edit,is_view,is_output,item_group,item_group_name,in_display_name) select '%s' as id,'%s' as form_template,name,description,default_value,sort,package_name,entity,attr_def_id,attr_def_name,attr_def_data_type,element_type,title,width,ref_package_name,ref_entity,data_options,required,regular,is_edit,is_view,is_output,item_group,item_group_name,in_display_name from form_item_template where id='%s'", newGuidList[i], newFormTemplateId, item.Id)})
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
	var requestTemplateTable []*models.RequestTemplateTable
	err = x.SQL("select * from request_template where del_flag=0 and id in (select request_template from request_template_role where `role` in ('" + strings.Join(userRoles, "','") + "')) order by `group`,id").Find(&requestTemplateTable)
	if err != nil {
		return
	}
	if len(requestTemplateTable) == 0 {
		return
	}
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
	return
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
