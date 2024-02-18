package service

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"xorm.io/xorm"

	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/exterror"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/rpc"
)

type RequestTemplateService struct {
	requestTemplateDao     dao.RequestTemplateDao
	requestTemplateRoleDao dao.RequestTemplateRoleDao
	operationLogDao        dao.OperationLogDao
}

func (s RequestTemplateService) QueryRequestTemplate(param *models.QueryRequestParam, commonParam models.CommonParam) (pageInfo models.PageInfo, result []*models.RequestTemplateQueryObj, err error) {
	var roleMap = make(map[string]*models.SimpleLocalRoleDto)
	extFilterSql := ""
	result = []*models.RequestTemplateQueryObj{}
	isQueryMessage := false
	if len(param.Filters) > 0 {
		var newFilters []*models.QueryRequestFilterObj
		for _, v := range param.Filters {
			if v.Name == "id" {
				isQueryMessage = true
			}
			if v.Name == "mgmtRoles" || v.Name == "useRoles" {
				inValueList := v.Value.([]interface{})
				var inValueStringList []string
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
				roleFilterSql, roleFilterParam := dao.CreateListParams(inValueStringList, "")
				if v.Name == "mgmtRoles" {
					tmpIds, tmpErr = s.getRequestTemplateIdsBySql("select t1.id from request_template t1 left join request_template_role t2 on t1.id=t2.request_template where t2.role_type='MGMT' and t2.role in ("+roleFilterSql+")", roleFilterParam)
				} else {
					tmpIds, tmpErr = s.getRequestTemplateIdsBySql("select t1.id from request_template t1 left join request_template_role t2 on t1.id=t2.request_template where t2.role_type='USE' and t2.role in ("+roleFilterSql+")", roleFilterParam)
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
	var rowData []*models.RequestTemplateTable
	filterSql, queryColumn, queryParam := dao.TransFiltersToSQL(param, &models.TransFiltersParam{IsStruct: true, StructObj: models.RequestTemplateTable{}, PrimaryKey: "id", Prefix: "t1"})
	userRolesFilterSql, userRolesFilterParam := dao.CreateListParams(commonParam.Roles, "")
	queryParam = append(userRolesFilterParam, queryParam...)
	baseSql := fmt.Sprintf("SELECT %s FROM (select * from request_template where del_flag=0 or (del_flag=2 and id not in (select record_id from request_template where del_flag=2 and record_id<>''))) t1 WHERE t1.id in (select request_template from request_template_role where role_type='MGMT' and `role` in ("+userRolesFilterSql+")) %s %s ", queryColumn, extFilterSql, filterSql)
	if param.Paging {
		pageInfo.StartIndex = param.Pageable.StartIndex
		pageInfo.PageSize = param.Pageable.PageSize
		pageInfo.TotalRows = dao.QueryCount(baseSql, queryParam...)
		pageSql, pageParam := dao.TransPageInfoToSQL(*param.Pageable)
		baseSql += pageSql
		queryParam = append(queryParam, pageParam...)
	}
	err = dao.X.SQL(baseSql, queryParam...).Find(&rowData)
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
	err = dao.X.SQL(queryRoleSql).Find(&requestTemplateRows)
	if err != nil {
		return
	}
	var mgmtRoleMap = make(map[string][]*models.RoleTable)
	var useRoleMap = make(map[string][]*models.RoleTable)
	for _, v := range requestTemplateRows {
		var tmpRoles []*models.RoleTable
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
			tmpErr := s.SyncProcDefId(v.Id, v.ProcDefId, v.ProcDefName, v.ProcDefKey, commonParam.Token, commonParam.Language)
			if tmpErr != nil {
				err = fmt.Errorf("Try to sync proDefId fail,%s ", tmpErr.Error())
				break
			}
		}
		if err != nil {
			return
		}
	}
	roleMap, err = rpc.QueryAllRoles("Y", commonParam.Token, commonParam.Language)
	if err != nil {
		return
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
		tmpObj.ModifyType = s.getRequestTemplateModifyType(v)
		// 模板管理角色收敛,只能唯一,读取角色管理员
		if len(tmpObj.MGMTRoles) > 0 && roleMap[tmpObj.MGMTRoles[0].Id] != nil {
			tmpObj.Administrator = roleMap[tmpObj.MGMTRoles[0].Id].Administrator
		}
		result = append(result, &tmpObj)
	}
	return
}

func (s RequestTemplateService) UpdateRequestTemplateStatus(requestTemplateId, user, status, reason string) (err error) {
	requestTemplate := &models.RequestTemplateTable{Id: requestTemplateId, Status: status, UpdatedBy: user,
		UpdatedTime: time.Now().Format(models.DateTimeFormat)}
	// 状态更新到草稿,需要退回
	if status == string(models.RequestStatusDraft) {
		requestTemplate.RollbackDesc = reason
	}
	return s.requestTemplateDao.Update(nil, requestTemplate)
}

func (s RequestTemplateService) UpdateRequestTemplateHandler(requestTemplateId, handler string) (err error) {
	return s.requestTemplateDao.Update(nil, &models.RequestTemplateTable{Id: requestTemplateId, Handler: handler, UpdatedBy: handler,
		UpdatedTime: time.Now().Format(models.DateTimeFormat)})
}

func (s RequestTemplateService) UpdateRequestTemplateBase(session *xorm.Session, requestTemplateId, formTemplate, description, updatedBy string, expireDay int) (err error) {
	now := time.Now().Format(models.DateTimeFormat)
	requestTemplate := &models.RequestTemplateTable{Id: requestTemplateId, FormTemplate: formTemplate, Description: description, ExpireDay: expireDay, UpdatedBy: updatedBy, UpdatedTime: now}
	return s.requestTemplateDao.Update(session, requestTemplate)
}
func (s RequestTemplateService) UpdateRequestTemplateUpdatedBy(session *xorm.Session, requestTemplateId, updatedBy string) (err error) {
	now := time.Now().Format(models.DateTimeFormat)
	requestTemplate := &models.RequestTemplateTable{Id: requestTemplateId, UpdatedBy: updatedBy, UpdatedTime: now}
	return s.requestTemplateDao.Update(session, requestTemplate)
}

func (s RequestTemplateService) UpdateRequestTemplateDataForm(session *xorm.Session, requestTemplateId, dataFormTemplate, updatedBy string) (err error) {
	now := time.Now().Format(models.DateTimeFormat)
	requestTemplate := &models.RequestTemplateTable{Id: requestTemplateId, DataFormTemplate: dataFormTemplate, UpdatedBy: updatedBy, UpdatedTime: now}
	return s.requestTemplateDao.Update(session, requestTemplate)
}
func (s RequestTemplateService) UpdateRequestTemplateStatusToCreated(id, operator string) (err error) {
	nowTime := time.Now().Format(models.DateTimeFormat)
	return s.requestTemplateDao.Update(nil, &models.RequestTemplateTable{Id: id, Status: "created", UpdatedBy: operator, UpdatedTime: nowTime})
}

func (s RequestTemplateService) GetRequestTemplate(requestTemplateId string) (requestTemplate *models.RequestTemplateTable, err error) {
	return s.requestTemplateDao.Get(requestTemplateId)
}
func (s RequestTemplateService) QueryRequestTemplateEntity(requestTemplateId, userToken, language string) (entityList []*models.RequestTemplateEntityDto, err error) {
	var requestTemplate *models.RequestTemplateTable
	var procDefEntities []*models.ProcEntity
	var workflowEntityMap = make(map[string]bool)
	var nodesList []*models.DataModel
	entityList = []*models.RequestTemplateEntityDto{}
	requestTemplate, err = s.GetRequestTemplate(requestTemplateId)
	if err != nil {
		return
	}
	if requestTemplate == nil {
		err = fmt.Errorf("requestTemplate not exist")
		return
	}
	entityList = append(entityList, &models.RequestTemplateEntityDto{
		FormType: "1.自定义表单",
		Entities: nil,
	})
	// 配置了编排数据
	if strings.TrimSpace(requestTemplate.ProcDefId) != "" {
		procDefEntities, err = s.ListRequestTemplateEntityAttrs(requestTemplateId, userToken, language)
		if err != nil {
			return
		}
		if len(procDefEntities) > 0 {
			var entities []string
			for _, entity := range procDefEntities {
				entityStr := fmt.Sprintf("%s:%s", entity.PackageName, entity.Name)
				entities = append(entities, entityStr)
				workflowEntityMap[entityStr] = true
			}
			entityList = append(entityList, &models.RequestTemplateEntityDto{FormType: "2.编排数据项表单", Entities: entities})
		}
	}
	// 自选数据项表单
	nodesList, err = rpc.QueryAllModels(userToken, language)
	if err != nil {
		return
	}
	if len(nodesList) > 0 {
		var entities []string
		for _, model := range nodesList {
			if len(model.Entities) > 0 {
				for _, entity := range model.Entities {
					entityStr := fmt.Sprintf("%s:%s", entity.PackageName, entity.Name)
					if !workflowEntityMap[entityStr] {
						entities = append(entities, entityStr)
					}
				}
			}
		}
		entityList = append(entityList, &models.RequestTemplateEntityDto{FormType: "3.自选数据项表单", Entities: entities})
	}
	return
}

func (s RequestTemplateService) CheckRequestTemplateRoles(requestTemplateId string, userRoles []string) error {
	has, err := s.requestTemplateRoleDao.CheckRequestTemplateRoles(requestTemplateId, userRoles)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf(models.RowDataPermissionErr)
	}
	return nil
}

func (s RequestTemplateService) CreateRequestTemplate(param models.RequestTemplateUpdateParam) (result models.RequestTemplateQueryObj, err error) {
	newGuid := guid.CreateGuid()
	param.Id = newGuid
	result = models.RequestTemplateQueryObj{RequestTemplateTable: param.RequestTemplateTable, MGMTRoles: []*models.RoleTable{}, USERoles: []*models.RoleTable{}}
	result.Id = newGuid
	err = transaction(func(session *xorm.Session) error {
		var err error
		_, err = s.requestTemplateDao.AddBasicInfo(session, models.ConvertRequestTemplateUpdateParam2RequestTemplate(param))
		if err != nil {
			return err
		}
		for _, role := range param.MGMTRoles {
			result.MGMTRoles = append(result.MGMTRoles, &models.RoleTable{Id: role})
			_, err = s.requestTemplateRoleDao.Add(session, models.CreateRequestTemplateRoleTable(newGuid, role, models.RolePermissionMGMT))
			if err != nil {
				return err
			}
		}
		for _, role := range param.USERoles {
			result.USERoles = append(result.USERoles, &models.RoleTable{Id: role})
			_, err = s.requestTemplateRoleDao.Add(session, models.CreateRequestTemplateRoleTable(newGuid, role, models.RolePermissionUse))
			if err != nil {
				return err
			}
		}
		return nil
	})
	return
}

func (s RequestTemplateService) GetAllCoreProcess(userToken, language string) map[string]string {
	var processMap = make(map[string]string)
	// 查询全部流程
	result, _ := GetProcDefService().GetCoreProcessListAll(userToken, language)
	if len(result) > 0 {
		for _, procDef := range result {
			if procDef != nil {
				processMap[procDef.ProcDefKey] = procDef.RootEntity.DisplayName
			}
		}
	}
	return processMap
}

// getRequestTemplateModifyType 模板版本 > v1表示 模板有多个版本,不允许多个版本都去修改模板类型,要求保持一致
func (s RequestTemplateService) getRequestTemplateModifyType(requestTemplate *models.RequestTemplateTable) bool {
	if strings.Compare(requestTemplate.Version, "v1") > 0 {
		return false
	}
	return true
}

func (s RequestTemplateService) getUpdateNodeDefIdActions(requestTemplateId, userToken, language string) (actions []*dao.ExecAction) {
	actions = []*dao.ExecAction{}
	var taskTemplate []*models.TaskTemplateTable
	dao.X.SQL("select * from task_template where request_template=?", requestTemplateId).Find(&taskTemplate)
	if len(taskTemplate) == 0 {
		return actions
	}
	nodeList, _ := GetProcDefService().GetProcessDefineTaskNodes(&models.RequestTemplateTable{Id: requestTemplateId}, userToken, language, "template")
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
				actions = append(actions, &dao.ExecAction{Sql: "update task_template set node_def_id=? where node_id=?", Param: []interface{}{nowDefId, v.NodeId}})
				actions = append(actions, &dao.ExecAction{Sql: "update task set node_def_id=? where task_template=?", Param: []interface{}{nowDefId, v.Id}})
			}
		}
	}
	return actions
}

func (s RequestTemplateService) SyncProcDefId(requestTemplateId, proDefId, proDefName, proDefKey, userToken, language string) error {
	log.Logger.Info("Start sync process def id")
	proExistFlag, newProDefId, err := GetProcDefService().CheckProDefId(proDefId, proDefName, proDefKey, userToken, language)
	if err != nil {
		return err
	}
	var actions []*dao.ExecAction
	if !proExistFlag {
		if proDefKey != "" {
			actions = append(actions, &dao.ExecAction{Sql: "update request_template set proc_def_id=? where id=?", Param: []interface{}{newProDefId, requestTemplateId}})
		} else {
			actions = append(actions, &dao.ExecAction{Sql: "update request_template set proc_def_id=? where proc_def_name=?", Param: []interface{}{newProDefId, proDefName}})
		}
		err = dao.Transaction(actions)
		if err != nil {
			return fmt.Errorf("Update requestTemplate procDefId fail,%s ", err.Error())
		}
		log.Logger.Info("Update requestTemplate proDefId done")
		actions = []*dao.ExecAction{}
	}
	tmpActions := s.getUpdateNodeDefIdActions(requestTemplateId, userToken, language)
	if len(tmpActions) > 0 {
		err = dao.Transaction(tmpActions)
		if err != nil {
			return fmt.Errorf("Update template node def id fail,%s ", err.Error())
		}
		log.Logger.Info("Update taskTemplate nodeDefId done")
	}
	return nil
}

func (s RequestTemplateService) getRequestTemplateIdsBySql(sql string, param []interface{}) (ids []string, err error) {
	var requestTemplateTables []*models.RequestTemplateTable
	err = dao.X.SQL(sql, param...).Find(&requestTemplateTables)
	ids = []string{}
	for _, v := range requestTemplateTables {
		ids = append(ids, v.Id)
	}
	return
}

func (s RequestTemplateService) UpdateRequestTemplate(param *models.RequestTemplateUpdateParam) (result models.RequestTemplateQueryObj, err error) {
	var actions []*dao.ExecAction
	nowTime := time.Now().Format(models.DateTimeFormat)
	result = models.RequestTemplateQueryObj{RequestTemplateTable: param.RequestTemplateTable, MGMTRoles: []*models.RoleTable{}, USERoles: []*models.RoleTable{}}
	updateAction := dao.ExecAction{Sql: "update request_template set status='created',`group`=?,name=?,description=?,tags=?,package_name=?,entity_name=?,proc_def_key=?,proc_def_id=?,proc_def_name=?,expire_day=?,handler=?,updated_by=?,updated_time=?,type=? where id=?"}
	updateAction.Param = []interface{}{param.Group, param.Name, param.Description, param.Tags, param.PackageName, param.EntityName, param.ProcDefKey, param.ProcDefId, param.ProcDefName, param.ExpireDay, param.Handler, param.UpdatedBy, nowTime, param.Type, param.Id}
	actions = append(actions, &updateAction)
	actions = append(actions, &dao.ExecAction{Sql: "delete from request_template_role where request_template=?", Param: []interface{}{param.Id}})
	for _, v := range param.MGMTRoles {
		result.MGMTRoles = append(result.MGMTRoles, &models.RoleTable{Id: v})
		actions = append(actions, &dao.ExecAction{Sql: "insert into request_template_role(id,request_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{param.Id + models.SysTableIdConnector + v + models.SysTableIdConnector + "MGMT", param.Id, v, "MGMT"}})
	}
	for _, v := range param.USERoles {
		result.USERoles = append(result.USERoles, &models.RoleTable{Id: v})
		actions = append(actions, &dao.ExecAction{Sql: "insert into request_template_role(id,request_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{param.Id + models.SysTableIdConnector + v + models.SysTableIdConnector + "USE", param.Id, v, "USE"}})
	}
	err = dao.Transaction(actions)
	return
}

func (s RequestTemplateService) DeleteRequestTemplate(id string, getActionFlag bool) (actions []*dao.ExecAction, err error) {
	rtObj, err := GetRequestTemplateService().GetRequestTemplate(id)
	if err != nil {
		return actions, err
	}
	if rtObj.Status == "confirm" {
		return actions, fmt.Errorf("confirm status can not delete")
	}
	var taskTemplateTable []*models.TaskTemplateTable
	dao.X.SQL("select id,form_template from task_template where request_template=?", id).Find(&taskTemplateTable)
	formTemplateIds := []string{rtObj.FormTemplate}
	for _, v := range taskTemplateTable {
		formTemplateIds = append(formTemplateIds, v.FormTemplate)
	}
	actions = []*dao.ExecAction{}
	var requestTable []*models.RequestTable
	dao.X.SQL("select id,name from request where request_template=?", id).Find(&requestTable)
	if len(requestTable) > 0 {
		var formTable []*models.FormTable
		dao.X.SQL("select id from form where form_template in ('" + strings.Join(formTemplateIds, "','") + "')").Find(&formTable)
		formIds := []string{}
		for _, v := range formTable {
			formIds = append(formIds, v.Id)
		}
		//actions = append(actions, &execAction{Sql: "delete from operation_log where task in (select id from task where task_template in (select id from task_template where request_template=?))", Param: []interface{}{id}})
		//actions = append(actions, &execAction{Sql: "delete from operation_log where request in (select id from request where request_template=?)", Param: []interface{}{id}})
		actions = append(actions, &dao.ExecAction{Sql: "delete from task where task_template in (select id from task_template where request_template=?)", Param: []interface{}{id}})
		actions = append(actions, &dao.ExecAction{Sql: "delete from request where request_template=?", Param: []interface{}{id}})
		actions = append(actions, &dao.ExecAction{Sql: "delete from form_item where form in ('" + strings.Join(formIds, "','") + "')", Param: []interface{}{}})
		actions = append(actions, &dao.ExecAction{Sql: "delete from form where form_template in ('" + strings.Join(formTemplateIds, "','") + "')", Param: []interface{}{}})
	}
	actions = append(actions, &dao.ExecAction{Sql: "delete from task_template_role where task_template in (select id from task_template where request_template=?)", Param: []interface{}{id}})
	actions = append(actions, &dao.ExecAction{Sql: "delete from task_template where request_template=?", Param: []interface{}{id}})
	actions = append(actions, &dao.ExecAction{Sql: "delete from request_template_role where request_template=?", Param: []interface{}{id}})
	actions = append(actions, &dao.ExecAction{Sql: "delete from request_template where id=?", Param: []interface{}{id}})
	actions = append(actions, &dao.ExecAction{Sql: "delete from form_item_template where form_template in ('" + strings.Join(formTemplateIds, "','") + "')", Param: []interface{}{}})
	actions = append(actions, &dao.ExecAction{Sql: "delete from form_template where id in ('" + strings.Join(formTemplateIds, "','") + "')", Param: []interface{}{}})
	if !getActionFlag {
		err = dao.Transaction(actions)
	}
	return actions, err
}

func (s RequestTemplateService) ListRequestTemplateEntityAttrs(id, userToken, language string) (result []*models.ProcEntity, err error) {
	var nodes []*models.ProcNodeObj
	result = []*models.ProcEntity{}
	nodes, err = GetProcDefService().GetProcessDefineTaskNodes(&models.RequestTemplateTable{Id: id}, userToken, language, "all")
	if err != nil {
		return
	}
	if len(nodes) == 0 {
		return
	}
	entityMap := make(map[string]int)
	existAttrMap := make(map[string]int)
	existAttrs, _ := s.GetRequestTemplateEntityAttrs(id)
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

func (s RequestTemplateService) GetRequestTemplateEntityAttrs(id string) (result []*models.ProcEntityAttributeObj, err error) {
	result = []*models.ProcEntityAttributeObj{}
	var requestTemplateTable []*models.RequestTemplateTable
	err = dao.X.SQL("select entity_attrs from request_template where id=?", id).Find(&requestTemplateTable)
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

func (s RequestTemplateService) UpdateRequestTemplateEntityAttrs(id string, attrs []*models.ProcEntityAttributeObj, operator string) error {
	b, _ := json.Marshal(attrs)
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err := dao.X.Exec("update request_template set entity_attrs=?,updated_time=?,updated_by=? where id=?", string(b), nowTime, operator, id)
	return err
}

func (s RequestTemplateService) GetRequestTemplateManageRole(id string) (role string) {
	var roleList []string
	err := dao.X.SQL("select role from request_template_role where request_template=? and role_type='MGMT'", id).Find(&roleList)
	if err != nil {
		err = fmt.Errorf("Try to query database fail,%s ", err.Error())
		return
	}
	if len(roleList) > 0 {
		role = roleList[0]
	}
	return
}

func getRequestTemplateRole(templateId string) (requestTemplateRoleList []*models.RequestTemplateRoleTable, err error) {
	err = dao.X.SQL("select * from request_template_role where request_template=?", templateId).Find(&requestTemplateRoleList)
	if err != nil {
		err = fmt.Errorf("Try to query database fail,%s ", err.Error())
		return
	}
	return
}

func getAllRequestTemplate() (templateMap map[string]*models.RequestTemplateTable, err error) {
	templateMap = make(map[string]*models.RequestTemplateTable)
	var requestTemplateTable []*models.RequestTemplateTable
	err = dao.X.SQL("select * from request_template").Find(&requestTemplateTable)
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
	requestTemplateObj, err := GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if err != nil {
		return err
	}
	existQuery, tmpErr := dao.X.QueryString("select id,name,version from request_template where del_flag!=1 and record_id=?", requestTemplateObj.Id)
	if tmpErr != nil {
		return fmt.Errorf("Query database fail,%s ", tmpErr.Error())
	}
	if len(existQuery) > 0 {
		return fmt.Errorf("RequestTemplate already have a branch %s:%s", existQuery[0]["name"], existQuery[0]["version"])
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	version := common.BuildVersionNum(requestTemplateObj.Version)
	newRequestTemplateId := guid.CreateGuid()
	newRequestFormTemplateId := guid.CreateGuid()
	var actions []*dao.ExecAction
	if requestTemplateObj.ParentId == "" {
		actions = append(actions, &dao.ExecAction{Sql: fmt.Sprintf("insert into request_template(id,`group`,name,description,form_template,"+
			"tags,status,package_name,entity_name,proc_def_key,proc_def_id,proc_def_name,created_by,created_time,updated_by,updated_time,"+
			"entity_attrs,record_id,`version`,confirm_time,expire_day,handler,type,operator_obj_type) select '%s' as id,`group`,name,description,'%s' as form_template,"+
			"tags,'created' as status,package_name,entity_name,proc_def_key,proc_def_id,proc_def_name,'%s' as created_by,'%s' as created_time,"+
			"'%s' as updated_by,'%s' as updated_time,entity_attrs,'%s' as record_id,'%s' as `version`,'' as confirm_time,expire_day,handler, "+
			"type,operator_obj_type from request_template where id='%s'", newRequestTemplateId, newRequestFormTemplateId, operator, nowTime, operator, nowTime,
			requestTemplateObj.Id, version, requestTemplateObj.Id)})
	} else {
		actions = append(actions, &dao.ExecAction{Sql: fmt.Sprintf("insert into request_template(id,`group`,name,description,form_template,"+
			"tags,status,package_name,entity_name,proc_def_key,proc_def_id,proc_def_name,created_by,created_time,updated_by,updated_time,"+
			"entity_attrs,record_id,`version`,confirm_time,expire_day,handler,type,operator_obj_type,parent_id) select '%s' as id,`group`,name,"+
			"description,'%s' as form_template,tags,'created' as status,package_name,entity_name,proc_def_key,proc_def_id,proc_def_name,"+
			"'%s' as created_by,'%s' as created_time,'%s' as updated_by,'%s' as updated_time,entity_attrs,'%s' as record_id,'%s' as `version`,"+
			"'' as confirm_time,expire_day,handler,type,operator_obj_type,'%s' as parent_id from request_template where id='%s'", newRequestTemplateId, newRequestFormTemplateId, operator,
			nowTime, operator, nowTime, requestTemplateObj.Id, version, requestTemplateObj.Id, requestTemplateObj.ParentId)})
	}
	newRequestFormActions, tmpErr := getFormCopyActions(requestTemplateObj.FormTemplate, newRequestFormTemplateId)
	if tmpErr != nil {
		return fmt.Errorf("Try to copy request form fail,%s ", tmpErr.Error())
	}
	actions = append(actions, newRequestFormActions...)
	var requestTemplateRoles []*models.RequestTemplateRoleTable
	dao.X.SQL("select * from request_template_role where request_template=?", requestTemplateObj.Id).Find(&requestTemplateRoles)
	for _, v := range requestTemplateRoles {
		tmpId := newRequestTemplateId + models.SysTableIdConnector + v.Role + models.SysTableIdConnector + v.RoleType
		actions = append(actions, &dao.ExecAction{Sql: "insert into request_template_role(id,request_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{tmpId, newRequestTemplateId, v.Role, v.RoleType}})
	}
	var taskTemplates []*models.TaskTemplateTable
	dao.X.SQL("select id,form_template from task_template where request_template=?", requestTemplateObj.Id).Find(&taskTemplates)
	newTaskGuids := guid.CreateGuidList(len(taskTemplates))
	newTaskFormGuids := guid.CreateGuidList(len(taskTemplates))
	for i, task := range taskTemplates {
		actions = append(actions, &dao.ExecAction{Sql: fmt.Sprintf("insert into task_template(id,name,description,form_template,request_template,node_id,node_def_id,node_name,expire_day,handler,created_by,created_time,updated_by,updated_time) select '%s' as id,name,description,'%s' as form_template,'%s' as request_template,node_id,node_def_id,node_name,expire_day,handler,created_by,created_time,updated_by,updated_time from task_template where id='%s'", newTaskGuids[i], newTaskFormGuids[i], newRequestTemplateId, task.Id)})
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
	return dao.TransactionWithoutForeignCheck(actions)
}

func ConfirmRequestTemplate(requestTemplateId string) error {
	var parentId string
	requestTemplateObj, err := GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
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
	if requestTemplateObj.RecordId != "" {
		prevRequestTemplateObj, _ := GetRequestTemplateService().GetRequestTemplate(requestTemplateObj.RecordId)
		parentId = prevRequestTemplateObj.ParentId
	}
	version := requestTemplateObj.Version
	if version == "" {
		version = "v1"
	}
	var actions []*dao.ExecAction
	actions = append(actions, &dao.ExecAction{Sql: "update request_template set status='confirm',`version`=?,confirm_time=?,del_flag=2,parent_id=? where id=?", Param: []interface{}{version, nowTime, parentId, requestTemplateObj.Id}})
	return dao.Transaction(actions)
}

func getFormCopyActions(oldFormTemplateId, newFormTemplateId string) (actions []*dao.ExecAction, err error) {
	var itemRows []*models.FormItemTemplateTable
	err = dao.X.SQL("select id from form_item_template where form_template=?", oldFormTemplateId).Find(&itemRows)
	if err != nil {
		return
	}
	actions = append(actions, &dao.ExecAction{Sql: fmt.Sprintf("insert into form_template(id,name,description,created_by,created_time,updated_by,updated_time) select '%s' as id,name,description,created_by,created_time,updated_by,updated_time from form_template where id='%s'", newFormTemplateId, oldFormTemplateId)})
	newGuidList := guid.CreateGuidList(len(itemRows))
	for i, item := range itemRows {
		actions = append(actions, &dao.ExecAction{Sql: fmt.Sprintf("insert into form_item_template(id,form_template,name,description,default_value,sort,package_name,entity,attr_def_id,attr_def_name,attr_def_data_type,element_type,title,width,ref_package_name,ref_entity,data_options,required,regular,is_edit,is_view,is_output,item_group,item_group_name,in_display_name,is_ref_inside,multiple,default_clear) select '%s' as id,'%s' as form_template,name,description,default_value,sort,package_name,entity,attr_def_id,attr_def_name,attr_def_data_type,element_type,title,width,ref_package_name,ref_entity,data_options,required,regular,is_edit,is_view,is_output,item_group,item_group_name,in_display_name,is_ref_inside,multiple,default_clear from form_item_template where id='%s'", newGuidList[i], newFormTemplateId, item.Id)})
	}
	return
}

func getTaskTemplateRoleActions(oldTaskTemplateId, newTaskTemplateId string) (actions []*dao.ExecAction, err error) {
	var taskTemplateRoles []*models.TaskTemplateRoleTable
	err = dao.X.SQL("select * from task_template_role where task_template=?", oldTaskTemplateId).Find(&taskTemplateRoles)
	if err != nil {
		return
	}
	for _, v := range taskTemplateRoles {
		tmpId := newTaskTemplateId + models.SysTableIdConnector + v.Role + models.SysTableIdConnector + v.RoleType
		actions = append(actions, &dao.ExecAction{Sql: "insert into task_template_role(id,task_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{tmpId, newTaskTemplateId, v.Role, v.RoleType}})
	}
	return
}

func validateConfirm(requestTemplateId string) error {
	var taskTemplateTable []*models.TaskTemplateTable
	dao.X.SQL("select id from task_template where request_template=? and form_template IS NOT NULL", requestTemplateId).Find(&taskTemplateTable)
	if len(requestTemplateId) == 0 {
		return fmt.Errorf("Please config task template ")
	}
	return nil
}

// GetRequestTemplateByUserV2  新的选择模板接口
func (s RequestTemplateService) GetRequestTemplateByUserV2(user, userToken, language string, userRoles []string) (result []*models.UserRequestTemplateQueryObjNew, err error) {
	var operatorObjTypeMap = make(map[string]string)
	var roleTemplateGroupMap = make(map[string]map[string][]*models.RequestTemplateTableObj)
	var resultMap = make(map[string]*models.UserRequestTemplateQueryObjNew)
	var roleList []string
	var requestTemplateTable, allTemplateTable, tmpTemplateTable []*models.RequestTemplateTable
	var requestTemplateRoleTable []*models.RequestTemplateRoleTable
	var ownerRoleMap = make(map[string]string)
	var requestTemplateLatestMap = make(map[string]*models.RequestTemplateTable)
	var userRoleMap = convertArray2Map(userRoles)
	result = []*models.UserRequestTemplateQueryObjNew{}
	useGroupMap, _ := getAllRequestTemplateGroup()
	userRolesFilterSql, userRolesFilterParam := dao.CreateListParams(userRoles, "")
	err = dao.X.SQL("select * from request_template ").Find(&allTemplateTable)
	if err != nil {
		return
	}
	err = dao.X.SQL("select * from request_template where del_flag=2 and id in (select request_template from request_template_role where role_type='USE' and `role` in ("+userRolesFilterSql+"))  order by `group`,tags,status,id", userRolesFilterParam...).Find(&requestTemplateTable)
	if err != nil {
		return
	}
	if len(requestTemplateTable) == 0 {
		return
	}
	requestTemplateLatestMap = getLatestVersionTemplate(requestTemplateTable, allTemplateTable)
	err = dao.X.SQL("select * from request_template_role where role_type='MGMT'").Find(&requestTemplateRoleTable)
	if err != nil {
		return
	}
	for _, requestTemplateRole := range requestTemplateRoleTable {
		ownerRoleMap[requestTemplateRole.RequestTemplate] = requestTemplateRole.Role
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
			v.Version = "beta"
		}
		// 此处需要查询db判断 用户是否有当前模板的最新发布的权限
		if _, ok := recordIdMap[v.Id]; !ok {
			// 获取最新模板Id,当前模板不是最新模板,直接加入 recordIdMap
			if latestTemplate, ok := requestTemplateLatestMap[v.Id]; ok && latestTemplate.Id != v.Id {
				recordIdMap[v.Id] = 1
			}
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
	}
	requestTemplateTable = tmpTemplateTable
	// 查询当前用户所有收藏模板记录
	collectMap, _ := QueryAllTemplateCollect(user)
	var collectFlag int
	// 组装数据
	// 操作对象类型,新增模板是录入.历史模板操作对象类型为空,需要全量处理下
	operatorObjTypeMap = s.GetAllCoreProcess(userToken, language)
	if len(requestTemplateTable) > 0 {
		for _, template := range requestTemplateTable {
			collectFlag = 0
			var tempRoleArr, roleArr []string
			err = dao.X.SQL("SELECT role FROM request_template_role WHERE request_template = ?  AND role_type = 'USE' ", template.Id).Find(&tempRoleArr)
			if err != nil {
				continue
			}
			if len(tempRoleArr) > 0 {
				for _, role := range tempRoleArr {
					if userRoleMap[role] {
						roleArr = append(roleArr, role)
					}
				}
			}
			if len(collectMap) > 0 && collectMap[template.ParentId] {
				collectFlag = 1
			}
			for _, role := range roleArr {
				if _, ok := roleTemplateGroupMap[role]; !ok {
					roleTemplateGroupMap[role] = make(map[string][]*models.RequestTemplateTableObj)
				}
				if _, ok := roleTemplateGroupMap[role][template.Group]; !ok {
					roleTemplateGroupMap[role][template.Group] = make([]*models.RequestTemplateTableObj, 0)
				}
				if template.OperatorObjType == "" {
					template.OperatorObjType = operatorObjTypeMap[template.ProcDefKey]
				}
				roleTemplateGroupMap[role][template.Group] = append(roleTemplateGroupMap[role][template.Group], &models.RequestTemplateTableObj{
					Id:              template.Id,
					Name:            template.Name,
					Version:         template.Version,
					Tags:            template.Tags,
					Status:          template.Status,
					UpdatedBy:       template.UpdatedBy,
					Handler:         template.Handler,
					Role:            ownerRoleMap[template.Id],
					UpdatedTime:     template.UpdatedTime,
					CollectFlag:     collectFlag,
					Type:            template.Type,
					OperatorObjType: template.OperatorObjType,
				})
			}
		}
	}
	for role, roleGroupMap := range roleTemplateGroupMap {
		groups := make([]*models.TemplateGroupObj, 0)
		for groupId, templateArr := range roleGroupMap {
			useGroup := useGroupMap[groupId]
			groups = append(groups, &models.TemplateGroupObj{
				GroupId:     groupId,
				GroupName:   useGroup.Name,
				CreatedTime: useGroup.CreatedTime,
				UpdatedTime: useGroup.UpdatedTime,
				Templates:   templateArr,
			})
		}
		resultMap[role] = &models.UserRequestTemplateQueryObjNew{
			ManageRole: role,
			Groups:     groups,
		}
		roleList = append(roleList, role)
	}
	// 角色排序
	sort.Strings(roleList)
	for _, role := range roleList {
		group := resultMap[role].Groups
		if len(group) > 0 {
			// 模版组排序
			sort.Sort(models.TemplateGroupSort(group))
			for _, templateObj := range group {
				//模板排序
				sort.Sort(models.RequestTemplateSort(templateObj.Templates))
			}
		}
		result = append(result, resultMap[role])
	}
	return
}

// getLatestVersionTemplate 获取requestTemplateList每个模板的最新发布或者禁用版本模板)
func getLatestVersionTemplate(requestTemplateList, allRequestTemplateList []*models.RequestTemplateTable) map[string]*models.RequestTemplateTable {
	allTemplateMap := make(map[string]*models.RequestTemplateTable)
	allTemplateRecordMap := make(map[string]*models.RequestTemplateTable)
	resultMap := make(map[string]*models.RequestTemplateTable)
	for _, requestTemplate := range allRequestTemplateList {
		allTemplateMap[requestTemplate.Id] = requestTemplate
	}
	for _, requestTemplate := range allRequestTemplateList {
		if requestTemplate.RecordId != "" {
			allTemplateRecordMap[requestTemplate.RecordId] = requestTemplate
		}
	}
	for _, requestTemplate := range requestTemplateList {
		var latestTemplate *models.RequestTemplateTable
		// 有 recordId指向当前模板,需要一直遍历找到最新版本模板
		if v, ok := allTemplateRecordMap[requestTemplate.Id]; ok && v != nil {
			temp := v
			for {
				if t, ok := allTemplateRecordMap[temp.Id]; ok && t != nil {
					temp = t
				} else {
					latestTemplate = temp
					break
				}
			}
		} else {
			// 没有 recordId指向当前模板,表示当前模板就是最新版本模板
			latestTemplate = requestTemplate
		}
		if latestTemplate == nil {
			// 找不到兜底默认值
			log.Logger.Warn("latestTemplate is empty", log.String("requestTemplateId", requestTemplate.Id))
			latestTemplate = requestTemplate
		}
		resultMap[requestTemplate.Id] = latestTemplate
		// 如果最新版本是创建状态,需要记录上一个版本模板
		if latestTemplate.Status == "created" && latestTemplate.RecordId != "" {
			resultMap[latestTemplate.RecordId] = allTemplateMap[latestTemplate.RecordId]
		}
	}
	return resultMap
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

func GetRequestTemplateTags(group string) (result []string, err error) {
	result = []string{}
	var requestTemplates []*models.RequestTemplateTable
	err = dao.X.SQL("select distinct tags from request_template where `group`=?", group).Find(&requestTemplates)
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
	err = dao.X.SQL("select * from request_template where id=?", requestTemplateId).Find(&requestTemplateTable)
	if err != nil {
		return
	}
	if len(requestTemplateTable) == 0 {
		err = fmt.Errorf("Can not find requestTemplate with id:%s ", requestTemplateId)
		return
	}
	result.RequestTemplate = *requestTemplateTable[0]
	dao.X.SQL("select * from request_template_role where request_template=?", requestTemplateId).Find(&result.RequestTemplateRole)
	dao.X.SQL("select * from task_template where request_template=?", requestTemplateId).Find(&result.TaskTemplate)
	dao.X.SQL("select * from task_template_role where task_template in (select id from task_template where request_template=?)", requestTemplateId).Find(&result.TaskTemplateRole)
	dao.X.SQL("select * from form_template where id in (select form_template from request_template where id=? union select form_template from task_template where request_template=?)", requestTemplateId, requestTemplateId).Find(&result.FormTemplate)
	dao.X.SQL("select * from form_item_template where form_template in (select id from form_template where id in (select form_template from request_template where id=? union select form_template from task_template where request_template=?))", requestTemplateId, requestTemplateId).Find(&result.FormItemTemplate)
	var requestTemplateGroupTable []*models.RequestTemplateGroupTable
	dao.X.SQL("select * from request_template_group where id=?", result.RequestTemplate.Group).Find(&requestTemplateGroupTable)
	if len(requestTemplateGroupTable) > 0 {
		result.RequestTemplateGroup = *requestTemplateGroupTable[0]
	}
	return
}

func (s RequestTemplateService) RequestTemplateImport(input models.RequestTemplateExport, userToken, language, confirmToken, operator string) (templateName, backToken string, err error) {
	var actions []*dao.ExecAction
	var inputVersion = getTemplateVersion(&input.RequestTemplate)
	var templateList []*models.RequestTemplateTable
	// 记录重复并且是草稿态的Id
	var repeatTemplateIdList []string
	if confirmToken == "" {
		// 1.判断名称是否重复
		templateName = input.RequestTemplate.Name
		templateList, err = getTemplateListByName(input.RequestTemplate.Name)
		if err != nil {
			return templateName, backToken, err
		}
		if len(templateList) > 0 {
			// 有名称重复数据,判断导入版本是否高于所有模板版本
			for _, template := range templateList {
				// 导入版本 低于同名版本,直接报错
				if inputVersion <= getTemplateVersion(template) {
					err = exterror.New().ImportTemplateVersionConflictError
					return
				}
				if template.Status == "created" {
					repeatTemplateIdList = append(repeatTemplateIdList, template.Id)
					models.RequestTemplateImportMap[template.Id] = input
				}
			}
			if len(repeatTemplateIdList) > 0 {
				backToken = strings.Join(repeatTemplateIdList, ",")
				return
			} else {
				// 有重复数据,但是新导入模板版本最高,直接当成新建处理
				input = createNewImportTemplate(input, input.RequestTemplate.RecordId)
			}
		} else {
			// 无名称重复数据，新建模板id以及模板关联表id都新建
			input = createNewImportTemplate(input, "")
		}
	} else {
		// 删除冲突模板数据
		confirmTokenList := strings.Split(confirmToken, ",")
		for _, ct := range confirmTokenList {
			if inputCache, b := models.RequestTemplateImportMap[ct]; b {
				input = inputCache
			} else {
				err = fmt.Errorf("Fetch input cache fail,please refersh and try again ")
				return
			}
			delete(models.RequestTemplateImportMap, ct)
			delActions, delErr := s.DeleteRequestTemplate(ct, true)
			if delErr != nil {
				err = delErr
				return
			}
			actions = append(actions, delActions...)
		}
		// 新建模板&模板相关表属性
		input = createNewImportTemplate(input, input.RequestTemplate.RecordId)
	}
	if input.RequestTemplate.Id == "" {
		err = fmt.Errorf("RequestTemplate id illegal ")
		return
	}
	allProcessList, processErr := GetProcDefService().GetCoreProcessListAll(userToken, language)
	if processErr != nil {
		err = fmt.Errorf("Get core process list fail,%s ", processErr.Error())
		return
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
		err = fmt.Errorf("Reqeust process:%s can not find! ", input.RequestTemplate.ProcDefName)
		return
	}
	nodeList, _ := GetProcDefService().GetProcessDefineTaskNodes(&input.RequestTemplate, userToken, language, "template")
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
		return
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	for _, v := range input.FormTemplate {
		actions = append(actions, &dao.ExecAction{Sql: "insert into form_template(id,name,description,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?)", Param: []interface{}{v.Id, v.Name, v.Description, operator, nowTime, operator, nowTime}})
	}
	for _, v := range input.FormItemTemplate {
		tmpAction := dao.ExecAction{Sql: "insert into form_item_template(id,form_template,name,description,item_group,item_group_name,default_value,sort,package_name,entity,attr_def_id,attr_def_name,attr_def_data_type,element_type,title,width,ref_package_name,ref_entity,data_options,required,regular,is_edit,is_view,is_output,in_display_name,is_ref_inside,multiple,default_clear) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
		tmpAction.Param = []interface{}{v.Id, v.FormTemplate, v.Name, v.Description, v.ItemGroup, v.ItemGroupName, v.DefaultValue, v.Sort, v.PackageName, v.Entity, v.AttrDefId, v.AttrDefName, v.AttrDefDataType, v.ElementType, v.Title, v.Width, v.RefPackageName, v.RefEntity, v.DataOptions, v.Required, v.Regular, v.IsEdit, v.IsView, v.IsOutput, v.InDisplayName, v.IsRefInside, v.Multiple, v.DefaultClear}
		actions = append(actions, &tmpAction)
	}
	var roleTable []*models.RoleTable
	dao.X.SQL("select id from `role`").Find(&roleTable)
	roleMap := make(map[string]int)
	for _, v := range roleTable {
		roleMap[v.Id] = 1
	}
	var requestTemplateGroupTable []*models.RequestTemplateGroupTable
	dao.X.SQL("select id from request_template_group where id=?", input.RequestTemplate.Group).Find(&requestTemplateGroupTable)
	if len(requestTemplateGroupTable) == 0 {
		if _, b := roleMap[input.RequestTemplateGroup.ManageRole]; !b {
			input.RequestTemplateGroup.ManageRole = models.AdminRole
		}
		actions = append(actions, &dao.ExecAction{Sql: "insert into request_template_group(id,name,description,manage_role,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?,?)", Param: []interface{}{input.RequestTemplateGroup.Id, input.RequestTemplateGroup.Name, input.RequestTemplateGroup.Description, input.RequestTemplateGroup.ManageRole, operator, nowTime, operator, nowTime}})
	}
	input.RequestTemplate.Status = "created"
	input.RequestTemplate.ConfirmTime = ""
	rtAction := dao.ExecAction{Sql: "insert into request_template(id,`group`,name,description,form_template,tags,record_id,`version`,confirm_time,status,package_name,entity_name,proc_def_key,proc_def_id,proc_def_name,expire_day,created_by,created_time,updated_by,updated_time,entity_attrs,handler,type,operator_obj_type) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
	rtAction.Param = []interface{}{input.RequestTemplate.Id, input.RequestTemplate.Group, input.RequestTemplate.Name, input.RequestTemplate.Description, input.RequestTemplate.FormTemplate, input.RequestTemplate.Tags, input.RequestTemplate.RecordId, input.RequestTemplate.Version, input.RequestTemplate.ConfirmTime, input.RequestTemplate.Status, input.RequestTemplate.PackageName, input.RequestTemplate.EntityName,
		input.RequestTemplate.ProcDefKey, input.RequestTemplate.ProcDefId, input.RequestTemplate.ProcDefName, input.RequestTemplate.ExpireDay, operator, nowTime, operator, nowTime, input.RequestTemplate.EntityAttrs, input.RequestTemplate.Handler, input.RequestTemplate.Type, input.RequestTemplate.OperatorObjType}
	actions = append(actions, &rtAction)
	for _, v := range input.TaskTemplate {
		tmpAction := dao.ExecAction{Sql: "insert into task_template(id,name,description,form_template,request_template,node_id,node_def_id,node_name,expire_day,handler,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
		tmpAction.Param = []interface{}{v.Id, v.Name, v.Description, v.FormTemplate, v.RequestTemplate, v.NodeId, v.NodeDefId, v.NodeName, v.ExpireDay, v.Handler, operator, nowTime, operator, nowTime}
		actions = append(actions, &tmpAction)
	}
	rtRoleFetch := false
	for _, v := range input.RequestTemplateRole {
		if _, b := roleMap[v.Role]; b {
			rtRoleFetch = true
			actions = append(actions, &dao.ExecAction{Sql: "insert into request_template_role(id,request_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{v.Id, v.RequestTemplate, v.Role, v.RoleType}})
		}
	}
	if !rtRoleFetch {
		actions = append(actions, &dao.ExecAction{Sql: "insert into request_template_role(id,request_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{guid.CreateGuid() + models.SysTableIdConnector + models.AdminRole + models.SysTableIdConnector + "MGMT", input.RequestTemplate.Id, models.AdminRole, "MGMT"}})
	}
	for _, v := range input.TaskTemplateRole {
		if _, b := roleMap[v.Role]; b {
			actions = append(actions, &dao.ExecAction{Sql: "insert into task_template_role(id,task_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{v.Id, v.TaskTemplate, v.Role, v.RoleType}})
		}
	}
	err = dao.Transaction(actions)
	return
}

func createNewImportTemplate(input models.RequestTemplateExport, recordId string) models.RequestTemplateExport {
	var historyTemplateId = input.RequestTemplate.Id
	input.RequestTemplate.Id = guid.CreateGuid()
	input.RequestTemplate.RecordId = recordId
	// 修改模板角色中模板id,新建角色id
	for _, requestTemplateRole := range input.RequestTemplateRole {
		requestTemplateRole.Id = guid.CreateGuid()
		requestTemplateRole.RequestTemplate = input.RequestTemplate.Id
	}
	// 修改 formTemplate
	for _, formTemplate := range input.FormTemplate {
		historyFormTemplateId := formTemplate.Id
		formTemplate.Id = guid.CreateGuid()
		// 修改模板里面的 formTemplateId
		if input.RequestTemplate.FormTemplate == historyFormTemplateId {
			input.RequestTemplate.FormTemplate = formTemplate.Id
		}
		for _, formItemTemplate := range input.FormItemTemplate {
			formItemTemplate.Id = guid.CreateGuid()
			if formItemTemplate.FormTemplate == historyFormTemplateId {
				formItemTemplate.FormTemplate = formTemplate.Id
			}
		}
		// 修改 taskTemplate中formTemplate,RequestTemplate,以及taskTemplateRole修改
		for _, taskTemplate := range input.TaskTemplate {
			historyTaskTemplateId := taskTemplate.Id
			taskTemplate.Id = guid.CreateGuid()
			if taskTemplate.FormTemplate == historyFormTemplateId {
				taskTemplate.FormTemplate = formTemplate.Id
			}
			if taskTemplate.RequestTemplate == historyTemplateId {
				taskTemplate.RequestTemplate = input.RequestTemplate.Id
			}
			for _, taskTemplateRole := range input.TaskTemplateRole {
				taskTemplateRole.Id = guid.CreateGuid()
				if taskTemplateRole.TaskTemplate == historyTaskTemplateId {
					taskTemplateRole.TaskTemplate = taskTemplate.Id
				}
			}
		}
	}

	return input
}

func getTemplateVersion(template *models.RequestTemplateTable) int {
	var version int
	if template == nil {
		return 0
	}
	if len(template.Version) > 1 {
		version, _ = strconv.Atoi(template.Version[1:])
	}
	return version
}

func getTemplateListByName(templateName string) (requestTemplateTable []*models.RequestTemplateTable, err error) {
	err = dao.X.SQL("select id,name,version,status from request_template where name=?", templateName).Find(&requestTemplateTable)
	return
}

func DisableRequestTemplate(requestTemplateId, operator string) (err error) {
	queryRows, queryErr := dao.X.QueryString("select status from request_template where id=?", requestTemplateId)
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
	_, err = dao.X.Exec("update request_template set status='disable' where id=?", requestTemplateId)
	return
}

func EnableRequestTemplate(requestTemplateId, operator string) (err error) {
	queryRows, queryErr := dao.X.QueryString("select status from request_template where id=?", requestTemplateId)
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
	_, err = dao.X.Exec("update request_template set status='confirm' where id=?", requestTemplateId)
	return
}

func getAllRequestTemplateGroup() (groupMap map[string]*models.RequestTemplateGroupTable, err error) {
	groupMap = make(map[string]*models.RequestTemplateGroupTable)
	var allGroupTable []*models.RequestTemplateGroupTable
	err = dao.X.SQL("select id,name,description,manage_role,created_time,updated_time from request_template_group").Find(&allGroupTable)
	if err != nil {
		return
	}
	for _, group := range allGroupTable {
		groupMap[group.Id] = group
	}
	return
}

func UpdateRequestTemplateParentId(requestTemplate *models.RequestTemplateTable) (parentId string) {
	var actions []*dao.ExecAction
	var templateIds []string
	// 老模板有多个版本,需要更新所有版本,并找到 recordId为空的记录
	requestTemplateMap, _ := getAllRequestTemplate()
	if len(requestTemplateMap) > 0 {
		temp := requestTemplate
		for {
			if temp == nil {
				break
			}
			templateIds = append(templateIds, temp.Id)
			if temp.RecordId == "" {
				parentId = temp.Id
				break
			}
			temp = requestTemplateMap[temp.RecordId]
		}
	}
	if len(templateIds) > 0 && parentId != "" {
		for _, templateId := range templateIds {
			actions = append(actions, &dao.ExecAction{Sql: "update request_template set parent_id=? where id=?", Param: []interface{}{parentId, templateId}})
		}
	}
	if len(actions) > 0 {
		updateErr := dao.Transaction(actions)
		if updateErr != nil {
			log.Logger.Error("Try to update request_template parent_id fail", log.Error(updateErr))
		}
	}
	return
}

func UpdateRequestTemplateParentIdById(templateId, parentId string) (err error) {
	var actions []*dao.ExecAction
	actions = append(actions, &dao.ExecAction{Sql: "update request_template set parent_id=? where id=?", Param: []interface{}{parentId, templateId}})
	if len(actions) > 0 {
		err = dao.Transaction(actions)
		if err != nil {
			log.Logger.Error("Try to update request_template parent_id fail", log.Error(err))
		}
	}
	return
}
