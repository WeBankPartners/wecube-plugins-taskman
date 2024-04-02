package service

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/exterror"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/rpc"
	"sort"
	"strconv"
	"strings"
	"time"
	"xorm.io/xorm"
)

type RequestTemplateService struct {
	requestTemplateDao     *dao.RequestTemplateDao
	requestTemplateRoleDao *dao.RequestTemplateRoleDao
	operationLogDao        *dao.OperationLogDao
	taskTemplateDao        *dao.TaskTemplateDao
	taskHandleTemplateDao  *dao.TaskHandleTemplateDao
	formTemplateDao        *dao.FormTemplateDao
}

func (s *RequestTemplateService) GetDtoByRequestTemplate(requestTemplate *models.RequestTemplateTable) (requestTemplateDto *models.RequestTemplateDto) {
	var taskTemplateList []*models.TaskTemplateTable
	var checkTaskHandleTemplateList []*models.TaskHandleTemplateTable
	var err error
	requestTemplateDto = models.ConvertRequestTemplateModel2Dto(requestTemplate)
	err = s.taskTemplateDao.DB.SQL("select * from task_template where request_template=?", requestTemplate.Id).Find(&taskTemplateList)
	if err != nil {
		log.Logger.Error("query task_template error", log.Error(err))
		return
	}
	if len(taskTemplateList) > 0 {
		for _, taskTemplate := range taskTemplateList {
			if taskTemplate.Type == string(models.TaskTypeCheck) {
				// 读取定版配置
				requestTemplateDto.CheckExpireDay = taskTemplate.ExpireDay
				err = s.taskTemplateDao.DB.SQL("select * from task_handle_template where task_template = ?", taskTemplate.Id).Find(&checkTaskHandleTemplateList)
				if err != nil {
					log.Logger.Error("query task_handle_template error", log.Error(err))
					return
				}
				if len(checkTaskHandleTemplateList) > 0 {
					requestTemplateDto.CheckRole = checkTaskHandleTemplateList[0].Role
					requestTemplateDto.CheckHandler = checkTaskHandleTemplateList[0].Handler
				}
			} else if taskTemplate.Type == string(models.TaskTypeConfirm) {
				// 读取请求确认配置
				requestTemplateDto.ConfirmExpireDay = taskTemplate.ExpireDay
			}
			// 前端展示
			if requestTemplateDto.ConfirmExpireDay == 0 {
				requestTemplateDto.ConfirmExpireDay = 1
			}
			if requestTemplateDto.CheckExpireDay == 0 {
				requestTemplateDto.CheckExpireDay = 1
			}
		}
	}
	return
}

func (s *RequestTemplateService) QueryRequestTemplate(param *models.QueryRequestParam, commonParam models.CommonParam) (pageInfo models.PageInfo, result []*models.RequestTemplateQueryObj, err error) {
	var roleMap = make(map[string]*models.SimpleLocalRoleDto)
	extFilterSql := ""
	result = []*models.RequestTemplateQueryObj{}
	isQueryMessage := false
	if roleMap, err = rpc.QueryAllRoles("Y", commonParam.Token, commonParam.Language); err != nil {
		return
	}
	if len(param.Filters) > 0 {
		var newFilters []*models.QueryRequestFilterObj
		for _, v := range param.Filters {
			if v.Name == "id" {
				isQueryMessage = true
			}
			if v.Name == "mgmtRoles" || v.Name == "useRoles" || v.Name == "type" {
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
				} else if v.Name == "useRoles" {
					tmpIds, tmpErr = s.getRequestTemplateIdsBySql("select t1.id from request_template t1 left join request_template_role t2 on t1.id=t2.request_template where t2.role_type='USE' and t2.role in ("+roleFilterSql+")", roleFilterParam)
				} else if v.Name == "type" {
					tmpIds, tmpErr = s.getRequestTemplateIdsBySql("select id from request_template where type in ("+roleFilterSql+")", roleFilterParam)
				}
				if tmpErr != nil {
					err = fmt.Errorf("try to query filter role id fail,%s ", tmpErr.Error())
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
	queryRoleSql := "select t4.id,GROUP_CONCAT(t4.role) as 'role','mgmt' as 'role_type' from ("
	queryRoleSql += "select t1.id,t2.role from request_template t1 left join request_template_role t2 on t1.id=t2.request_template  where t1.id in ('" + strings.Join(rtIds, "','") + "') and t2.role_type='MGMT'"
	queryRoleSql += ") t4 group by t4.id"
	queryRoleSql += " UNION "
	queryRoleSql += "select t4.id,GROUP_CONCAT(t4.role) as 'role','use' as 'role_type' from ("
	queryRoleSql += "select t1.id,t2.role from request_template t1 left join request_template_role t2 on t1.id=t2.request_template  where t1.id in ('" + strings.Join(rtIds, "','") + "') and t2.role_type='USE'"
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
			if roleDto, ok := roleMap[vv]; ok {
				tmpRoles = append(tmpRoles, &models.RoleTable{Id: roleDto.Name, DisplayName: roleDto.DisplayName})
			}
		}
		if v.RoleType == "mgmt" {
			mgmtRoleMap[v.Id] = tmpRoles
		} else {
			useRoleMap[v.Id] = tmpRoles
		}
	}
	if isQueryMessage {
		for _, v := range rowData {
			if v.ProcDefId != "" {
				tmpErr := s.SyncProcDefId(v.Id, v.ProcDefId, v.ProcDefName, v.ProcDefKey, commonParam.Token, commonParam.Language)
				if tmpErr != nil {
					err = fmt.Errorf("try to sync proDefId fail,%s ", tmpErr.Error())
					break
				}
			}
		}
		if err != nil {
			return
		}
	}
	for _, v := range rowData {
		tmpObj := models.RequestTemplateQueryObj{RequestTemplateDto: *s.GetDtoByRequestTemplate(v), MGMTRoles: []*models.RoleTable{}, USERoles: []*models.RoleTable{}}
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

func (s *RequestTemplateService) UpdateRequestTemplateStatus(requestTemplateId, user, status, reason string) (err error) {
	var rt *models.RequestTemplateTable
	requestTemplate := &models.RequestTemplateTable{Id: requestTemplateId, Status: status, UpdatedBy: user,
		UpdatedTime: time.Now().Format(models.DateTimeFormat)}
	// 状态更新到草稿,需要退回
	if status == string(models.RequestTemplateStatusCreated) {
		rt, err = s.GetRequestTemplate(requestTemplateId)
		if err != nil {
			return
		}
		requestTemplate.BackDesc = reason + "\n" + rt.BackDesc
	}
	return s.requestTemplateDao.Update(nil, requestTemplate)
}

func (s *RequestTemplateService) UpdateRequestTemplateHandler(requestTemplateId, handler string) (err error) {
	return s.requestTemplateDao.Update(nil, &models.RequestTemplateTable{Id: requestTemplateId, Handler: handler, UpdatedBy: handler,
		UpdatedTime: time.Now().Format(models.DateTimeFormat)})
}

func (s *RequestTemplateService) UpdateRequestTemplateUpdatedBy(session *xorm.Session, requestTemplateId, updatedBy string) (err error) {
	now := time.Now().Format(models.DateTimeFormat)
	requestTemplate := &models.RequestTemplateTable{Id: requestTemplateId, UpdatedBy: updatedBy, UpdatedTime: now}
	return s.requestTemplateDao.Update(session, requestTemplate)
}

func (s *RequestTemplateService) UpdateRequestTemplateStatusToCreated(id, operator string) (err error) {
	nowTime := time.Now().Format(models.DateTimeFormat)
	return s.requestTemplateDao.Update(nil, &models.RequestTemplateTable{Id: id, Status: "created", UpdatedBy: operator, UpdatedTime: nowTime})
}

func (s *RequestTemplateService) GetRequestTemplate(requestTemplateId string) (requestTemplate *models.RequestTemplateTable, err error) {
	return s.requestTemplateDao.Get(requestTemplateId)
}
func (s *RequestTemplateService) QueryRequestTemplateEntity(requestTemplateId, userToken, language string) (entityList []*models.RequestTemplateEntityDto, err error) {
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
		FormType: string(models.FormItemGroupTypeCustom),
		Entities: nil,
	})
	// 配置了编排数据
	if strings.TrimSpace(requestTemplate.ProcDefId) != "" {
		procDefEntities, err = s.ListRequestTemplateEntityAttrs(requestTemplateId, requestTemplate.ProcDefId, userToken, language)
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
			entityList = append(entityList, &models.RequestTemplateEntityDto{FormType: string(models.FormItemGroupTypeWorkflow), Entities: entities})
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
		entityList = append(entityList, &models.RequestTemplateEntityDto{FormType: string(models.FormItemGroupTypeOptional), Entities: entities})
	}
	return
}

func (s *RequestTemplateService) CheckRequestTemplateRoles(requestTemplateId string, userRoles []string) error {
	has, err := s.requestTemplateRoleDao.CheckRequestTemplateRoles(requestTemplateId, userRoles)
	if err != nil {
		return err
	}
	if !has {
		return fmt.Errorf(models.RowDataPermissionErr)
	}
	return nil
}

func (s *RequestTemplateService) CreateRequestTemplate(param models.RequestTemplateUpdateParam, userToken, language string) (result models.RequestTemplateQueryObj, err error) {
	newGuid := guid.CreateGuid()
	newCheckTaskId := fmt.Sprintf("ch_%s", guid.CreateGuid())
	newConfirmTaskId := fmt.Sprintf("co_%s", guid.CreateGuid())
	now := time.Now().Format(models.DateTimeFormat)
	param.Id = newGuid
	result = models.RequestTemplateQueryObj{RequestTemplateDto: param.RequestTemplateDto, MGMTRoles: []*models.RoleTable{}, USERoles: []*models.RoleTable{}}
	result.Id = newGuid
	err = transaction(func(session *xorm.Session) error {
		param.Status = "created"
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
		// 添加 提交请求模板
		_, err = s.taskTemplateDao.Add(session, &models.TaskTemplateTable{Id: "su_" + guid.CreateGuid(), Name: "submit", RequestTemplate: newGuid,
			Type: string(models.TaskTypeSubmit), CreatedTime: now, UpdatedTime: now})
		if err != nil {
			return err
		}
		// 任务模板添加定版任务和确认任务
		if param.CheckSwitch {
			_, err = s.taskTemplateDao.Add(session, &models.TaskTemplateTable{Id: newCheckTaskId, Name: "check", RequestTemplate: newGuid,
				ExpireDay: param.CheckExpireDay, Type: string(models.TaskTypeCheck), CreatedTime: now, UpdatedTime: now})
			if err != nil {
				return err
			}
			_, err = s.taskHandleTemplateDao.Add(session, &models.TaskHandleTemplateTable{Id: guid.CreateGuid(), TaskTemplate: newCheckTaskId,
				Role: param.CheckRole, Handler: param.CheckHandler})
			if err != nil {
				return err
			}
		}
		if param.ConfirmSwitch {
			_, err = s.taskTemplateDao.Add(session, &models.TaskTemplateTable{Id: newConfirmTaskId, Name: "confirm", RequestTemplate: newGuid,
				ExpireDay: param.ConfirmExpireDay, Type: string(models.TaskTypeConfirm), CreatedTime: now, UpdatedTime: now})
			if err != nil {
				return err
			}
		}
		// 开启了编排,创建编排任务&数据表单数据初始化
		if param.ProcDefId != "" {
			// 初始化编排表单
			if err = s.CreateWorkflowFormTemplate(session, newGuid, param.ProcDefId, userToken, language); err != nil {
				return err
			}
			if err = GetTaskTemplateService().createProcTaskTemplates(session, param.ProcDefId, newGuid, userToken, language, param.CreatedBy); err != nil {
				return err
			}
		}
		return nil
	})
	return
}

func (s *RequestTemplateService) GetAllCoreProcess(userToken, language string) map[string]string {
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
func (s *RequestTemplateService) getRequestTemplateModifyType(requestTemplate *models.RequestTemplateTable) bool {
	return strings.Compare(requestTemplate.Version, "v1") <= 0
}

func (s *RequestTemplateService) getUpdateNodeDefIdActions(requestTemplateId, userToken, language string) (actions []*dao.ExecAction) {
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

func (s *RequestTemplateService) SyncProcDefId(requestTemplateId, proDefId, proDefName, proDefKey, userToken, language string) error {
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
			return fmt.Errorf("update requestTemplate procDefId fail,%s ", err.Error())
		}
		log.Logger.Info("Update requestTemplate proDefId done")
	}
	tmpActions := s.getUpdateNodeDefIdActions(requestTemplateId, userToken, language)
	if len(tmpActions) > 0 {
		err = dao.Transaction(tmpActions)
		if err != nil {
			return fmt.Errorf("update template node def id fail,%s ", err.Error())
		}
		log.Logger.Info("Update taskTemplate nodeDefId done")
	}
	return nil
}

func (s *RequestTemplateService) getRequestTemplateIdsBySql(sql string, param []interface{}) (ids []string, err error) {
	var requestTemplateTables []*models.RequestTemplateTable
	err = dao.X.SQL(sql, param...).Find(&requestTemplateTables)
	ids = []string{}
	for _, v := range requestTemplateTables {
		ids = append(ids, v.Id)
	}
	return
}

func (s *RequestTemplateService) UpdateRequestTemplate(param *models.RequestTemplateUpdateParam, userToken, language string) (result models.RequestTemplateQueryObj, err error) {
	var actions, insertWorkflowFormTemplateActions, insertTaskTemplateActions, deleteFormActions []*dao.ExecAction
	var taskTemplateList, implementTaskTemplateList []*models.TaskTemplateTable
	var workflowTaskNodeMap = make(map[string]*models.TaskTemplateTable)
	var requestTemplate *models.RequestTemplateTable
	var workflowTaskNodeList []*models.ProcNodeObj
	var updateActions []*dao.ExecAction
	var delWorkFlowFormFlag bool
	nowTime := time.Now().Format(models.DateTimeFormat)
	if taskTemplateList, err = s.taskTemplateDao.QueryByRequestTemplate(param.Id); err != nil {
		return
	}
	if requestTemplate, err = s.GetRequestTemplate(param.Id); err != nil {
		return
	}
	result = models.RequestTemplateQueryObj{RequestTemplateDto: param.RequestTemplateDto, MGMTRoles: []*models.RoleTable{}, USERoles: []*models.RoleTable{}}
	updateAction := dao.ExecAction{Sql: "update request_template set status='created',`group`=?,name=?,description=?,tags=?,package_name=?,entity_name=?," +
		"proc_def_key=?,proc_def_id=?,proc_def_name=?,proc_def_version=?,expire_day=?,handler=?,updated_by=?,updated_time=?,type=?,operator_obj_type=?,approve_by=?,check_switch=?," +
		"confirm_switch=?,back_desc=? where id=?"}
	updateAction.Param = []interface{}{param.Group, param.Name, param.Description, param.Tags, param.PackageName, param.EntityName, param.ProcDefKey,
		param.ProcDefId, param.ProcDefName, param.ProcDefVersion, param.ExpireDay, param.Handler, param.UpdatedBy, nowTime, param.Type, param.OperatorObjType, param.ApproveBy,
		param.CheckSwitch, param.ConfirmSwitch, param.BackDesc, param.Id}
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

	// 删除定版任务和请求确认任务,记录编排任务
	if len(taskTemplateList) > 0 {
		for _, taskTemplate := range taskTemplateList {
			if taskTemplate.Type == string(models.TaskTypeCheck) || taskTemplate.Type == string(models.TaskTypeConfirm) {
				if taskTemplate.Type == string(models.TaskTypeCheck) {
					actions = append(actions, &dao.ExecAction{Sql: "delete from task_handle_template where task_template=?", Param: []interface{}{taskTemplate.Id}})
				}
				actions = append(actions, &dao.ExecAction{Sql: "delete from task_template where id=?", Param: []interface{}{taskTemplate.Id}})
			} else if taskTemplate.Type == string(models.TaskTypeImplement) {
				implementTaskTemplateList = append(implementTaskTemplateList, taskTemplate)
				if taskTemplate.NodeDefId != "" {
					workflowTaskNodeMap[taskTemplate.NodeDefId] = taskTemplate
				}
			}
		}
	}
	actions = append(actions, &dao.ExecAction{Sql: "delete from task_template where request_template=? and type=?", Param: []interface{}{param.Id, string(models.TaskTypeCheck)}})
	actions = append(actions, &dao.ExecAction{Sql: "delete from task_template where request_template=? and type=?", Param: []interface{}{param.Id, string(models.TaskTypeConfirm)}})
	// 根据参数重新任务模板添加定版任务和确认任务
	if param.CheckSwitch {
		newCheckTaskId := fmt.Sprintf("ch_%s", guid.CreateGuid())
		actions = append(actions, &dao.ExecAction{Sql: "insert into task_template(id,name,request_template,expire_day,type,created_time," +
			"updated_time) values (?,?,?,?,?,?,?)", Param: []interface{}{newCheckTaskId, "check", param.Id, param.CheckExpireDay, string(models.TaskTypeCheck), nowTime, nowTime}})
		actions = append(actions, &dao.ExecAction{Sql: "insert into task_handle_template(id,task_template,role,handler) values(?,?,?,?)",
			Param: []interface{}{guid.CreateGuid(), newCheckTaskId, param.CheckRole, param.CheckHandler}})
	}
	if param.ConfirmSwitch {
		newConfirmTaskId := fmt.Sprintf("co_%s", guid.CreateGuid())
		actions = append(actions, &dao.ExecAction{Sql: "insert into task_template(id,name,request_template,expire_day,type,created_time," +
			"updated_time) values (?,?,?,?,?,?,?)", Param: []interface{}{newConfirmTaskId, "confirm", param.Id, param.ConfirmExpireDay, string(models.TaskTypeConfirm), nowTime, nowTime}})
	}
	// 编排处理,(1)新增了编排 (2)编排节点名称有更新 (3)删除编排
	if requestTemplate.ProcDefId == "" && param.ProcDefId != "" {
		// 先删除已有任务
		for _, taskTemplate := range implementTaskTemplateList {
			var deleteTaskTemplateActions []*dao.ExecAction
			if deleteTaskTemplateActions, err = GetTaskTemplateService().deleteTaskTemplateSql(param.Id, taskTemplate.Id); err != nil {
				return
			}
			if len(deleteTaskTemplateActions) > 0 {
				actions = append(actions, deleteTaskTemplateActions...)
			}
		}
		// 新增编排表单&编排任务
		if insertWorkflowFormTemplateActions, err = s.CreateWorkflowFormTemplateSql(param.Id, param.ProcDefId, userToken, language); err != nil {
			return
		}
		if len(insertWorkflowFormTemplateActions) > 0 {
			actions = append(actions, insertWorkflowFormTemplateActions...)
		}
		if insertTaskTemplateActions, err = GetTaskTemplateService().createProcTaskTemplatesSql(param.ProcDefId, param.Id, userToken, language, "system"); err != nil {
			return
		}
		if len(insertTaskTemplateActions) > 0 {
			actions = append(actions, insertTaskTemplateActions...)
		}
		delWorkFlowFormFlag = true
	} else if requestTemplate.ProcDefId != "" {
		// 删除编排
		if param.ProcDefId == "" {
			for _, taskTemplate := range implementTaskTemplateList {
				var deleteTaskTemplateActions []*dao.ExecAction
				if deleteTaskTemplateActions, err = GetTaskTemplateService().deleteTaskTemplateSql(param.Id, taskTemplate.Id); err != nil {
					return
				}
				if len(deleteTaskTemplateActions) > 0 {
					actions = append(actions, deleteTaskTemplateActions...)
				}
			}
			delWorkFlowFormFlag = true
		} else if param.ProcDefId != "" {
			// 换了编排,先删除,再新增
			if requestTemplate.ProcDefId != param.ProcDefId {
				for _, taskTemplate := range implementTaskTemplateList {
					var deleteTaskTemplateActions []*dao.ExecAction
					if deleteTaskTemplateActions, err = GetTaskTemplateService().deleteTaskTemplateSql(param.Id, taskTemplate.Id); err != nil {
						return
					}
					if len(deleteTaskTemplateActions) > 0 {
						actions = append(actions, deleteTaskTemplateActions...)
					}
				}
				// 新增编排表单&编排任务
				if insertWorkflowFormTemplateActions, err = s.CreateWorkflowFormTemplateSql(param.Id, param.ProcDefId, userToken, language); err != nil {
					return
				}
				if len(insertWorkflowFormTemplateActions) > 0 {
					actions = append(actions, insertWorkflowFormTemplateActions...)
				}
				if insertTaskTemplateActions, err = GetTaskTemplateService().createProcTaskTemplatesSql(param.ProcDefId, param.Id, userToken, language, "system"); err != nil {
					return
				}
				if len(insertTaskTemplateActions) > 0 {
					actions = append(actions, insertTaskTemplateActions...)
				}
				delWorkFlowFormFlag = true
			} else {
				// 关联编排无改动,需要查找编排,看编排节点名称是否变更,有变更需要替换
				// 查询编排任务节点
				if workflowTaskNodeList, err = GetTaskTemplateService().getProcTaskTemplateNodes(requestTemplate.ProcDefId, userToken, language); err != nil {
					return
				}
				if len(workflowTaskNodeList) > 0 {
					for _, taskNode := range workflowTaskNodeList {
						taskNodeTmp := workflowTaskNodeMap[taskNode.NodeDefId]
						if taskNodeTmp != nil && (taskNodeTmp.NodeId != taskNode.NodeId || taskNodeTmp.NodeName != taskNode.NodeName) {
							if updateActions, err = GetTaskTemplateService().updateProcTaskTemplatesSql(taskNodeTmp.Id, taskNode.NodeId, taskNode.NodeName); err != nil {
								return
							}
							actions = append(actions, updateActions...)
						}
					}
				}
			}
		}
	}
	// 删除配置的编排 fromTemplateGroup,主要是删除数据表单和审批的,任务表单在删除任务逻辑会删除
	if delWorkFlowFormFlag {
		deleteFormActions, err = GetFormTemplateService().DeleteWorkflowFormTemplateGroupSql(requestTemplate.Id)
		if err != nil {
			return
		}
		if len(deleteFormActions) > 0 {
			actions = append(actions, deleteFormActions...)
		}
	}
	err = dao.Transaction(actions)
	return
}

func (s *RequestTemplateService) DeleteRequestTemplate(id string, getActionFlag bool) (actions []*dao.ExecAction, err error) {
	var taskTemplateList []*models.TaskTemplateTable
	var formTemplateList []*models.FormTemplateTable
	var requestTemplate *models.RequestTemplateTable
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(id)
	if err != nil {
		return actions, err
	}
	if requestTemplate.Status == "confirm" {
		return actions, fmt.Errorf("confirm status can not delete")
	}
	dao.X.SQL("select * from form_template where request_template = ?", id).Find(&formTemplateList)
	if len(formTemplateList) > 0 {
		for _, formTemplate := range formTemplateList {
			// 删除表单项模板表
			actions = append(actions, &dao.ExecAction{Sql: "delete from form_item_template WHERE form_template = ? ", Param: []interface{}{formTemplate.Id}})
		}
	}
	// 删除 表单模板
	actions = append(actions, &dao.ExecAction{Sql: "delete from form_template WHERE request_template = ? ", Param: []interface{}{id}})

	dao.X.SQL("select * from task_template where request_template = ?", id).Find(&taskTemplateList)
	if len(taskTemplateList) > 0 {
		for _, taskTemplate := range taskTemplateList {
			// 删除任务处理模板表
			actions = append(actions, &dao.ExecAction{Sql: "delete from task_handle_template WHERE task_template = ? ", Param: []interface{}{taskTemplate.Id}})
			// 删除表单模板表
			actions = append(actions, &dao.ExecAction{Sql: "delete from form_template WHERE task_template = ?", Param: []interface{}{taskTemplate.Id}})
		}
	}
	// 删除任务模版表
	actions = append(actions, &dao.ExecAction{Sql: "delete from task_template WHERE request_template = ?", Param: []interface{}{id}})

	// 删除模板角色
	actions = append(actions, &dao.ExecAction{Sql: "delete from request_template_role WHERE request_template = ?", Param: []interface{}{id}})
	// 删除请求模板
	actions = append(actions, &dao.ExecAction{Sql: "delete from request_template where id=?", Param: []interface{}{id}})
	if !getActionFlag {
		err = dao.Transaction(actions)
	}
	return actions, err
}

func (s *RequestTemplateService) ListRequestTemplateEntityAttrs(id, procDefId, userToken, language string) (result []*models.ProcEntity, err error) {
	var nodes []*models.ProcNodeObj
	result = []*models.ProcEntity{}
	nodes, err = GetProcDefService().GetProcessDefineTaskNodes(&models.RequestTemplateTable{Id: id, ProcDefId: procDefId}, userToken, language, "all")
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

func (s *RequestTemplateService) GetRequestTemplateEntityAttrs(id string) (result []*models.ProcEntityAttributeObj, err error) {
	result = []*models.ProcEntityAttributeObj{}
	var requestTemplateTable []*models.RequestTemplateTable
	err = dao.X.SQL("select entity_attrs from request_template where id=?", id).Find(&requestTemplateTable)
	if err != nil {
		return
	}
	if len(requestTemplateTable) == 0 {
		err = fmt.Errorf("can not find request template wit id:%s ", id)
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

func (s *RequestTemplateService) UpdateRequestTemplateEntityAttrs(id string, attrs []*models.ProcEntityAttributeObj, operator string) error {
	b, _ := json.Marshal(attrs)
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err := dao.X.Exec("update request_template set entity_attrs=?,updated_time=?,updated_by=? where id=?", string(b), nowTime, operator, id)
	return err
}

func (s *RequestTemplateService) GetRequestTemplateManageRole(id string) (role string) {
	var roleList []string
	err := dao.X.SQL("select role from request_template_role where request_template=? and role_type='MGMT'", id).Find(&roleList)
	if err != nil {
		log.Logger.Error("try to query database fail ", log.Error(err))
		return
	}
	if len(roleList) > 0 {
		role = roleList[0]
	}
	return
}

func (s *RequestTemplateService) GetRequestTemplateRole(templateId string) (requestTemplateRoleList []*models.RequestTemplateRoleTable, err error) {
	err = dao.X.SQL("select * from request_template_role where request_template=?", templateId).Find(&requestTemplateRoleList)
	if err != nil {
		err = fmt.Errorf("try to query database fail,%s ", err.Error())
		return
	}
	return
}

func (s *RequestTemplateService) getAllRequestTemplate() (templateMap map[string]*models.RequestTemplateTable, err error) {
	templateMap = make(map[string]*models.RequestTemplateTable)
	var requestTemplateTable []*models.RequestTemplateTable
	err = dao.X.SQL("select * from request_template").Find(&requestTemplateTable)
	if err != nil {
		err = fmt.Errorf("try to query database fail,%s ", err.Error())
		return
	}
	for _, template := range requestTemplateTable {
		templateMap[template.Id] = template
	}
	return
}

func (s *RequestTemplateService) ForkConfirmRequestTemplate(requestTemplateId, operator string) (err error) {
	var actions []*dao.ExecAction
	var requestTemplateRoles []*models.RequestTemplateRoleTable
	// 查询任务模版处理列表
	var taskHandleTemplateList []*models.TaskHandleTemplateTable
	var taskTemplateList []*models.TaskTemplateTable
	var formItemTemplateList []*models.FormItemTemplateTable
	// 新任务模版ID和老模板ID映射
	var newTaskTemplateIdMap = make(map[string]string)
	// 新表单模版ID 和老模板ID映射
	var newFormTemplateIdMap = make(map[string]string)
	// 新表单项ID和老表单项ID映射
	var newFormItemTemplateIdMap = make(map[string]string)
	var requestTemplate *models.RequestTemplateTable

	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if err != nil {
		return
	}
	if requestTemplate == nil {
		err = fmt.Errorf("requestTemplateId invalid")
		return
	}
	existQuery, tmpErr := dao.X.QueryString("select id,name,version,status from request_template where del_flag!=1 and record_id=?", requestTemplate.Id)
	if tmpErr != nil {
		return fmt.Errorf("query database fail,%s ", tmpErr.Error())
	}
	if len(existQuery) > 0 {
		if existQuery[0]["status"] == string(models.RequestTemplateStatusCreated) {
			err = exterror.New().RequestTemplateHasDraftError
		} else if existQuery[0]["status"] == string(models.RequestTemplateStatusPending) {
			err = exterror.New().RequestTemplateHasPendingError
		}
		return err
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	version := common.BuildVersionNum(requestTemplate.Version)
	newRequestTemplateId := guid.CreateGuid()
	if requestTemplate.ParentId == "" {
		actions = append(actions, &dao.ExecAction{Sql: fmt.Sprintf("insert into request_template(id,`group`,name,description,"+
			"tags,status,package_name,entity_name,proc_def_key,proc_def_id,proc_def_name,created_by,created_time,updated_by,updated_time,"+
			"entity_attrs,record_id,`version`,confirm_time,expire_day,handler,type,operator_obj_type,approve_by,check_switch,confirm_switch,proc_def_version) select '%s' as id,`group`,name,description,"+
			"tags,'created' as status,package_name,entity_name,proc_def_key,proc_def_id,proc_def_name,'%s' as created_by,'%s' as created_time,"+
			"'%s' as updated_by,'%s' as updated_time,entity_attrs,'%s' as record_id,'%s' as `version`,'' as confirm_time,expire_day,handler, "+
			"type,operator_obj_type,approve_by,check_switch,confirm_switch,proc_def_version from request_template where id='%s'", newRequestTemplateId, operator, nowTime, operator, nowTime,
			requestTemplate.Id, version, requestTemplate.Id)})
	} else {
		actions = append(actions, &dao.ExecAction{Sql: fmt.Sprintf("insert into request_template(id,`group`,name,description,"+
			"tags,status,package_name,entity_name,proc_def_key,proc_def_id,proc_def_name,created_by,created_time,updated_by,updated_time,"+
			"entity_attrs,record_id,`version`,confirm_time,expire_day,handler,type,operator_obj_type,parent_id,approve_by,check_switch,confirm_switch,proc_def_version) select '%s' as id,`group`,name,"+
			"description,tags,'created' as status,package_name,entity_name,proc_def_key,proc_def_id,proc_def_name,"+
			"'%s' as created_by,'%s' as created_time,'%s' as updated_by,'%s' as updated_time,entity_attrs,'%s' as record_id,'%s' as `version`,"+
			"'' as confirm_time,expire_day,handler,type,operator_obj_type,'%s' as parent_id,approve_by,check_switch,confirm_switch,proc_def_version from request_template where id='%s'", newRequestTemplateId, operator,
			nowTime, operator, nowTime, requestTemplate.Id, version, requestTemplate.Id, requestTemplate.ParentId)})
	}

	dao.X.SQL("select * from request_template_role where request_template=?", requestTemplate.Id).Find(&requestTemplateRoles)
	for _, v := range requestTemplateRoles {
		tmpId := newRequestTemplateId + models.SysTableIdConnector + v.Role + models.SysTableIdConnector + v.RoleType
		actions = append(actions, &dao.ExecAction{Sql: "insert into request_template_role(id,request_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{tmpId, newRequestTemplateId, v.Role, v.RoleType}})
	}

	// 查询 任务模板列表
	if err = dao.X.SQL("select * from task_template where request_template = ?", requestTemplateId).Find(&taskTemplateList); err != nil {
		return
	}
	// 查询表单模板列表
	var formTemplateList []*models.FormTemplateTable
	if err = dao.X.SQL("select * from form_template where request_template = ?", requestTemplateId).Find(&formTemplateList); err != nil {
		return
	}
	if len(formTemplateList) > 0 {
		var tempFormItemTemplateList []*models.FormItemTemplateTable
		for _, formTemplate := range formTemplateList {
			newFormTemplateIdMap[formTemplate.Id] = guid.CreateGuid()
			if err = dao.X.SQL("select * from form_item_template where form_template = ?", formTemplate.Id).Find(&tempFormItemTemplateList); err != nil {
				return
			}
		}
		if len(tempFormItemTemplateList) > 0 {
			formItemTemplateList = append(formItemTemplateList, tempFormItemTemplateList...)
		}
		if len(formItemTemplateList) > 0 {
			for _, formItemTemplate := range formItemTemplateList {
				newFormItemTemplateIdMap[formItemTemplate.Id] = guid.CreateGuid()
			}
		}
	}

	if len(taskTemplateList) > 0 {
		for _, taskTemplate := range taskTemplateList {
			prefix, _ := taskTemplateService.genTaskIdPrefix(taskTemplate.Type)
			newTaskTemplateIdMap[taskTemplate.Id] = prefix + "_" + guid.CreateGuid()
			var tempTaskHandleTemplateList []*models.TaskHandleTemplateTable
			err = dao.X.SQL("select * from task_handle_template where task_template = ?", taskTemplate.Id).Find(&tempTaskHandleTemplateList)
			if err != nil {
				return
			}
			if len(tempTaskHandleTemplateList) > 0 {
				taskHandleTemplateList = append(taskHandleTemplateList, tempTaskHandleTemplateList...)
			}
			actions = append(actions, &dao.ExecAction{Sql: "insert into task_template(id,name,description,request_template,node_id,node_def_id,node_name,expire_day,handler,created_by,created_time,updated_by,updated_time,del_flag,sort,handle_mode,type) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", Param: []interface{}{newTaskTemplateIdMap[taskTemplate.Id], taskTemplate.Name, taskTemplate.Description, newRequestTemplateId, taskTemplate.NodeId, taskTemplate.NodeDefId, taskTemplate.NodeName, taskTemplate.ExpireDay, taskTemplate.Handler, taskTemplate.CreatedBy, taskTemplate.CreatedTime, taskTemplate.UpdatedBy, taskTemplate.UpdatedTime, taskTemplate.DelFlag, taskTemplate.Sort, taskTemplate.HandleMode, taskTemplate.Type}})
		}
		if len(taskHandleTemplateList) > 0 {
			for _, taskHandleTemplate := range taskHandleTemplateList {
				actions = append(actions, &dao.ExecAction{Sql: "insert into task_handle_template(id,task_template,role,assign,handler_type,handler,handle_mode,sort)values(?,?,?,?,?,?,?,?)", Param: []interface{}{
					guid.CreateGuid(), newTaskTemplateIdMap[taskHandleTemplate.TaskTemplate], taskHandleTemplate.Role, taskHandleTemplate.Assign, taskHandleTemplate.HandlerType, taskHandleTemplate.Handler, taskHandleTemplate.HandleMode, taskHandleTemplate.Sort,
				}})
			}
		}
		if len(formTemplateList) > 0 {
			for _, formTemplate := range formTemplateList {
				if formTemplate.TaskTemplate == "" {
					actions = append(actions, &dao.ExecAction{Sql: "insert into form_template(id,request_template,item_group,item_group_name,item_group_type,item_group_rule,item_group_sort,created_time,ref_id,request_form_type,del_flag)values(?,?,?,?,?,?,?,?,?,?,?)", Param: []interface{}{
						newFormTemplateIdMap[formTemplate.Id], newRequestTemplateId, formTemplate.ItemGroup, formTemplate.ItemGroupName, formTemplate.ItemGroupType, formTemplate.ItemGroupRule, formTemplate.ItemGroupSort, nowTime, newFormTemplateIdMap[formTemplate.RefId], formTemplate.RequestFormType, formTemplate.DelFlag,
					}})
				} else {
					actions = append(actions, &dao.ExecAction{Sql: "insert into form_template(id,request_template,task_template,item_group,item_group_name,item_group_type,item_group_rule,item_group_sort,created_time,ref_id,request_form_type,del_flag)values(?,?,?,?,?,?,?,?,?,?,?,?)", Param: []interface{}{
						newFormTemplateIdMap[formTemplate.Id], newRequestTemplateId, newTaskTemplateIdMap[formTemplate.TaskTemplate], formTemplate.ItemGroup, formTemplate.ItemGroupName, formTemplate.ItemGroupType, formTemplate.ItemGroupRule, formTemplate.ItemGroupSort, nowTime, newFormTemplateIdMap[formTemplate.RefId], formTemplate.RequestFormType, formTemplate.DelFlag,
					}})
				}
			}
		}
		if len(formItemTemplateList) > 0 {
			for _, formItemTemplate := range formItemTemplateList {
				actions = append(actions, &dao.ExecAction{Sql: fmt.Sprintf("insert into form_item_template(id,form_template,name,description,default_value,sort,package_name,entity,attr_def_id,attr_def_name,attr_def_data_type,element_type,title,width,ref_package_name,ref_entity,data_options,required,regular,is_edit,is_view,is_output,item_group,item_group_name,in_display_name,is_ref_inside,multiple,default_clear,ref_id,routine_expression) select '%s' as id,'%s' as form_template,name,description,default_value,sort,package_name,entity,attr_def_id,attr_def_name,attr_def_data_type,element_type,title,width,ref_package_name,ref_entity,data_options,required,regular,is_edit,is_view,is_output,item_group,item_group_name,in_display_name,is_ref_inside,multiple,default_clear,'%s' as ref_id,routine_expression from form_item_template where id='%s'", newFormItemTemplateIdMap[formItemTemplate.Id], newFormTemplateIdMap[formItemTemplate.FormTemplate], newFormItemTemplateIdMap[formItemTemplate.RefId], formItemTemplate.Id)})
			}
		}
	}
	return dao.TransactionWithoutForeignCheck(actions)
}

func (s *RequestTemplateService) CopyConfirmRequestTemplate(requestTemplateId, operator string) (err error) {
	var actions []*dao.ExecAction
	var requestTemplateRoles []*models.RequestTemplateRoleTable
	// 查询任务模版处理列表
	var taskHandleTemplateList []*models.TaskHandleTemplateTable
	var taskTemplateList []*models.TaskTemplateTable
	var formItemTemplateList []*models.FormItemTemplateTable
	// 新任务模版ID和老模板ID映射
	var newTaskTemplateIdMap = make(map[string]string)
	// 新表单模版ID 和老模板ID映射
	var newFormTemplateIdMap = make(map[string]string)
	// 新表单项ID和老表单项ID映射
	var newFormItemTemplateIdMap = make(map[string]string)
	var requestTemplate *models.RequestTemplateTable
	var list []*models.RequestTemplateTable

	if requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(requestTemplateId); err != nil {
		return
	}
	if requestTemplate == nil {
		err = fmt.Errorf("requestTemplateId invalid")
		return
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	newRequestTemplateId := guid.CreateGuid()
	var requestName = requestTemplate.Name + "(1)"
	// 查询 当前requestName是否存在
	if list, err = s.requestTemplateDao.QueryListByName(requestName); err != nil {
		return
	}
	if len(list) > 0 {
		// 表示已经有重名模版,用当前时间
		requestName = fmt.Sprintf("%s-%s", requestTemplate.Name, time.Now().Format(models.NewDateTimeFormat))
	}
	actions = append(actions, &dao.ExecAction{Sql: fmt.Sprintf("insert into request_template(id,`group`,name,description,"+
		"tags,status,package_name,entity_name,proc_def_key,proc_def_id,proc_def_name,created_by,created_time,updated_by,updated_time,"+
		"entity_attrs,`version`,confirm_time,expire_day,handler,type,operator_obj_type,approve_by,check_switch,confirm_switch,back_desc,proc_def_version) select '%s' as id,`group`,'%s' as name,description,"+
		"tags,'created' as status,package_name,entity_name,proc_def_key,proc_def_id,proc_def_name,'%s' as created_by,'%s' as created_time,"+
		"'%s' as updated_by,'%s' as updated_time,entity_attrs,'%s' as `version`,'' as confirm_time,expire_day,handler, "+
		"type,operator_obj_type,approve_by,check_switch,confirm_switch,back_desc,proc_def_version from request_template where id='%s'", newRequestTemplateId, requestName, operator, nowTime, operator, nowTime, "", requestTemplate.Id)})

	dao.X.SQL("select * from request_template_role where request_template=?", requestTemplate.Id).Find(&requestTemplateRoles)
	for _, v := range requestTemplateRoles {
		tmpId := newRequestTemplateId + models.SysTableIdConnector + v.Role + models.SysTableIdConnector + v.RoleType
		actions = append(actions, &dao.ExecAction{Sql: "insert into request_template_role(id,request_template,`role`,role_type) value (?,?,?,?)", Param: []interface{}{tmpId, newRequestTemplateId, v.Role, v.RoleType}})
	}

	// 查询 任务模板列表
	if err = dao.X.SQL("select * from task_template where request_template = ?", requestTemplateId).Find(&taskTemplateList); err != nil {
		return
	}
	// 查询表单模板列表
	var formTemplateList []*models.FormTemplateTable
	if err = dao.X.SQL("select * from form_template where request_template = ?", requestTemplateId).Find(&formTemplateList); err != nil {
		return
	}
	if len(formTemplateList) > 0 {
		var tempFormItemTemplateList []*models.FormItemTemplateTable
		for _, formTemplate := range formTemplateList {
			newFormTemplateIdMap[formTemplate.Id] = guid.CreateGuid()
			if err = dao.X.SQL("select * from form_item_template where form_template = ?", formTemplate.Id).Find(&tempFormItemTemplateList); err != nil {
				return
			}
		}
		if len(tempFormItemTemplateList) > 0 {
			formItemTemplateList = append(formItemTemplateList, tempFormItemTemplateList...)
		}
		if len(formItemTemplateList) > 0 {
			for _, formItemTemplate := range formItemTemplateList {
				newFormItemTemplateIdMap[formItemTemplate.Id] = guid.CreateGuid()
			}
		}
	}

	if len(taskTemplateList) > 0 {
		for _, taskTemplate := range taskTemplateList {
			prefix, _ := taskTemplateService.genTaskIdPrefix(taskTemplate.Type)
			newTaskTemplateIdMap[taskTemplate.Id] = prefix + "_" + guid.CreateGuid()
			var tempTaskHandleTemplateList []*models.TaskHandleTemplateTable
			if err = dao.X.SQL("select * from task_handle_template where task_template = ?", taskTemplate.Id).Find(&tempTaskHandleTemplateList); err != nil {
				return
			}
			if len(tempTaskHandleTemplateList) > 0 {
				taskHandleTemplateList = append(taskHandleTemplateList, tempTaskHandleTemplateList...)
			}
			actions = append(actions, &dao.ExecAction{Sql: "insert into task_template(id,name,description,request_template,node_id,node_def_id,node_name,expire_day,handler,created_by,created_time,updated_by,updated_time,del_flag,sort,handle_mode,type) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)", Param: []interface{}{newTaskTemplateIdMap[taskTemplate.Id], taskTemplate.Name, taskTemplate.Description, newRequestTemplateId, taskTemplate.NodeId, taskTemplate.NodeDefId, taskTemplate.NodeName, taskTemplate.ExpireDay, taskTemplate.Handler, taskTemplate.CreatedBy, taskTemplate.CreatedTime, taskTemplate.UpdatedBy, taskTemplate.UpdatedTime, taskTemplate.DelFlag, taskTemplate.Sort, taskTemplate.HandleMode, taskTemplate.Type}})
		}
		if len(taskHandleTemplateList) > 0 {
			for _, taskHandleTemplate := range taskHandleTemplateList {
				actions = append(actions, &dao.ExecAction{Sql: "insert into task_handle_template(id,task_template,role,assign,handler_type,handler,handle_mode,sort)values(?,?,?,?,?,?,?,?)", Param: []interface{}{
					guid.CreateGuid(), newTaskTemplateIdMap[taskHandleTemplate.TaskTemplate], taskHandleTemplate.Role, taskHandleTemplate.Assign, taskHandleTemplate.HandlerType, taskHandleTemplate.Handler, taskHandleTemplate.HandleMode, taskHandleTemplate.Sort,
				}})
			}
		}
		if len(formTemplateList) > 0 {
			for _, formTemplate := range formTemplateList {
				if formTemplate.TaskTemplate == "" {
					actions = append(actions, &dao.ExecAction{Sql: "insert into form_template(id,request_template,item_group,item_group_name,item_group_type,item_group_rule,item_group_sort,created_time,ref_id,request_form_type,del_flag)values(?,?,?,?,?,?,?,?,?,?,?)", Param: []interface{}{
						newFormTemplateIdMap[formTemplate.Id], newRequestTemplateId, formTemplate.ItemGroup, formTemplate.ItemGroupName, formTemplate.ItemGroupType, formTemplate.ItemGroupRule, formTemplate.ItemGroupSort, nowTime, newFormTemplateIdMap[formTemplate.RefId], formTemplate.RequestFormType, formTemplate.DelFlag,
					}})
				} else {
					actions = append(actions, &dao.ExecAction{Sql: "insert into form_template(id,request_template,task_template,item_group,item_group_name,item_group_type,item_group_rule,item_group_sort,created_time,ref_id,request_form_type,del_flag)values(?,?,?,?,?,?,?,?,?,?,?,?)", Param: []interface{}{
						newFormTemplateIdMap[formTemplate.Id], newRequestTemplateId, newTaskTemplateIdMap[formTemplate.TaskTemplate], formTemplate.ItemGroup, formTemplate.ItemGroupName, formTemplate.ItemGroupType, formTemplate.ItemGroupRule, formTemplate.ItemGroupSort, nowTime, newFormTemplateIdMap[formTemplate.RefId], formTemplate.RequestFormType, formTemplate.DelFlag,
					}})
				}
			}
		}
		if len(formItemTemplateList) > 0 {
			for _, formItemTemplate := range formItemTemplateList {
				actions = append(actions, &dao.ExecAction{Sql: fmt.Sprintf("insert into form_item_template(id,form_template,name,description,default_value,sort,package_name,entity,attr_def_id,attr_def_name,attr_def_data_type,element_type,title,width,ref_package_name,ref_entity,data_options,required,regular,is_edit,is_view,is_output,item_group,item_group_name,in_display_name,is_ref_inside,multiple,default_clear,ref_id,routine_expression) select '%s' as id,'%s' as form_template,name,description,default_value,sort,package_name,entity,attr_def_id,attr_def_name,attr_def_data_type,element_type,title,width,ref_package_name,ref_entity,data_options,required,regular,is_edit,is_view,is_output,item_group,item_group_name,in_display_name,is_ref_inside,multiple,default_clear,'%s' as ref_id,routine_expression from form_item_template where id='%s'", newFormItemTemplateIdMap[formItemTemplate.Id], newFormTemplateIdMap[formItemTemplate.FormTemplate], newFormItemTemplateIdMap[formItemTemplate.RefId], formItemTemplate.Id)})
			}
		}
	}
	return dao.TransactionWithoutForeignCheck(actions)
}

func (s *RequestTemplateService) ConfirmRequestTemplate(requestTemplateId, operator, userToken, language string) error {
	var parentId string
	var requestTemplateRoleList []*models.RequestTemplateRoleTable
	requestTemplateObj, err := GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if err != nil {
		return err
	}
	if requestTemplateObj.Status == "confirm" {
		return fmt.Errorf("request template already confirm ")
	}
	err = s.validateConfirm(requestTemplateId)
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
	// 调用编排新增
	if requestTemplateObj.ProcDefId != "" {
		requestTemplateRoleList, err = s.requestTemplateRoleDao.QueryByRequestTemplateAndType(requestTemplateId, string(models.RolePermissionUse))
		if err != nil {
			return err
		}
		if len(requestTemplateRoleList) > 0 {
			var userRoles []string
			for _, roleTable := range requestTemplateRoleList {
				userRoles = append(userRoles, roleTable.Role)
			}
			_, err = rpc.SyncWorkflowUseRole(models.SyncUseRoleParam{ProcDefId: requestTemplateObj.ProcDefId, UseRoles: userRoles}, userToken, language)
			if err != nil {
				return err
			}
		}
	}
	var actions []*dao.ExecAction
	actions = append(actions, &dao.ExecAction{Sql: "update request_template set status='confirm',`version`=?,confirm_time=?,del_flag=2,parent_id=?,updated_by=?,updated_time=? where id=?", Param: []interface{}{version, nowTime, parentId, operator, nowTime, requestTemplateObj.Id}})
	return dao.Transaction(actions)
}

func (s *RequestTemplateService) validateConfirm(requestTemplateId string) error {
	var taskTemplateTable []*models.TaskTemplateTable
	dao.X.SQL("select id from task_template where request_template=? and form_template IS NOT NULL", requestTemplateId).Find(&taskTemplateTable)
	if len(requestTemplateId) == 0 {
		return fmt.Errorf("please config task template ")
	}
	return nil
}

// GetRequestTemplateByUserV2  新的选择模板接口
func (s *RequestTemplateService) GetRequestTemplateByUserV2(user, userToken, language string, userRoles []string) (result []*models.UserRequestTemplateQueryObjNew, err error) {
	var operatorObjTypeMap = make(map[string]string)
	var roleTemplateGroupMap = make(map[string]map[string][]*models.RequestTemplateTableObj)
	var resultMap = make(map[string]*models.UserRequestTemplateQueryObjNew)
	var roleList []string
	var requestTemplateTable, allTemplateTable, tmpTemplateTable []*models.RequestTemplateTable
	var requestTemplateRoleTable []*models.RequestTemplateRoleTable
	var ownerRoleMap = make(map[string]string)
	var requestTemplateLatestMap = make(map[string]*models.RequestTemplateTable)
	var userRoleMap = convertArray2Map(userRoles)
	var roleDisplayNameMap map[string]string
	roleDisplayNameMap, err = GetRoleService().GetRoleDisplayName(userToken, language)
	if err != nil {
		return
	}
	result = []*models.UserRequestTemplateQueryObjNew{}
	useGroupMap, _ := s.getAllRequestTemplateGroup()
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
	requestTemplateLatestMap = s.getLatestVersionTemplate(requestTemplateTable, allTemplateTable)
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
				if !common.CompareUpdateConfirmTime(v.UpdatedTime, v.ConfirmTime) {
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
				var roleDisplayName string
				if _, ok := roleTemplateGroupMap[role]; !ok {
					roleTemplateGroupMap[role] = make(map[string][]*models.RequestTemplateTableObj)
				}
				if _, ok := roleTemplateGroupMap[role][template.Group]; !ok {
					roleTemplateGroupMap[role][template.Group] = make([]*models.RequestTemplateTableObj, 0)
				}
				if template.OperatorObjType == "" {
					template.OperatorObjType = operatorObjTypeMap[template.ProcDefKey]
				}
				// 采用显示名展示
				if v, ok := roleDisplayNameMap[ownerRoleMap[template.Id]]; ok {
					roleDisplayName = v
				} else {
					roleDisplayName = ownerRoleMap[template.Id]
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
					RoleDisplay:     roleDisplayName,
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
			ManageRole:        role,
			ManageRoleDisplay: roleDisplayNameMap[role],
			Groups:            groups,
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
func (s *RequestTemplateService) getLatestVersionTemplate(requestTemplateList, allRequestTemplateList []*models.RequestTemplateTable) map[string]*models.RequestTemplateTable {
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
		if (latestTemplate.Status == string(models.RequestTemplateStatusCreated) || latestTemplate.Status == string(models.RequestTemplateStatusPending)) && latestTemplate.RecordId != "" {
			resultMap[latestTemplate.RecordId] = allTemplateMap[latestTemplate.RecordId]
		}
	}
	return resultMap
}

func (s *RequestTemplateService) GetRequestTemplateTags(group string) (result []string, err error) {
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

func (s *RequestTemplateService) RequestTemplateExport(requestTemplateId string) (result models.RequestTemplateExport, err error) {
	var requestTemplateTable []*models.RequestTemplateTable
	result.RequestTemplateRole = []*models.RequestTemplateRoleTable{}
	result.TaskTemplate = []*models.TaskTemplateTable{}
	result.TaskHandleTemplate = []*models.TaskHandleTemplateTable{}
	result.FormTemplate = []*models.FormTemplateTable{}
	result.FormItemTemplate = []*models.FormItemTemplateTable{}
	err = dao.X.SQL("select * from request_template where id=?", requestTemplateId).Find(&requestTemplateTable)
	if err != nil {
		return
	}
	if len(requestTemplateTable) == 0 {
		err = fmt.Errorf("can not find requestTemplate with id:%s ", requestTemplateId)
		return
	}
	result.RequestTemplate = *s.GetDtoByRequestTemplate(requestTemplateTable[0])
	dao.X.SQL("select * from request_template_role where request_template=?", requestTemplateId).Find(&result.RequestTemplateRole)
	dao.X.SQL("select * from task_template where request_template=?", requestTemplateId).Find(&result.TaskTemplate)
	dao.X.SQL("select * from task_handle_template where task_template in (select id from task_template where request_template=?)", requestTemplateId).Find(&result.TaskHandleTemplate)
	dao.X.SQL("select * from form_template where request_template = ?", requestTemplateId).Find(&result.FormTemplate)
	dao.X.SQL("select * from form_item_template where form_template in (select id from form_template where request_template = ?)", requestTemplateId).Find(&result.FormItemTemplate)
	var requestTemplateGroupTable []*models.RequestTemplateGroupTable
	dao.X.SQL("select * from request_template_group where id=?", result.RequestTemplate.Group).Find(&requestTemplateGroupTable)
	if len(requestTemplateGroupTable) > 0 {
		result.RequestTemplateGroup = *requestTemplateGroupTable[0]
	}
	return
}

func (s *RequestTemplateService) RequestTemplateImport(input models.RequestTemplateExport, userToken, language, confirmToken, operator string, userRoles []string) (templateName, backToken string, err error) {
	var actions []*dao.ExecAction
	var inputVersion = s.getTemplateVersion(models.ConvertRequestTemplateDto2Model(input.RequestTemplate))
	var templateList []*models.RequestTemplateTable
	var manageRole, maxVersionRecordId string
	// 记录最大版本
	var maxVersion int
	// 记录重复并且是草稿态的Id
	var repeatTemplateIdList []string
	if confirmToken == "" {
		// 1.判断名称是否重复
		templateName = input.RequestTemplate.Name
		templateList, err = s.getTemplateListByName(input.RequestTemplate.Name)
		if err != nil {
			return templateName, backToken, err
		}
		if len(templateList) > 0 {
			// 有名称重复数据,判断导入版本是否高于所有模板版本
			for _, template := range templateList {
				// 导入版本 低于同名版本,直接报错
				version := s.getTemplateVersion(template)
				if inputVersion <= version {
					err = exterror.New().ImportTemplateVersionConflictError
					return
				}
				if template.Status == string(models.RequestTemplateStatusCreated) || template.Status == string(models.RequestTemplateStatusPending) {
					repeatTemplateIdList = append(repeatTemplateIdList, template.Id)
					models.RequestTemplateImportMap[template.Id] = input
					break
				}
				if maxVersion < version {
					maxVersion = version
					maxVersionRecordId = template.Id
				}
			}
			if len(repeatTemplateIdList) > 0 {
				backToken = strings.Join(repeatTemplateIdList, ",")
				return
			} else {
				// 有重复数据,但是新导入模板版本最高,直接当成新建处理,需要记录最新发布版本Id
				input.RequestTemplate.RecordId = maxVersionRecordId
				input = s.createNewImportTemplate(input, operator, input.RequestTemplate.RecordId, userToken, language)
			}
		} else {
			// 无名称重复数据，新建模板id以及模板关联表id都新建
			input = s.createNewImportTemplate(input, operator, "", userToken, language)
		}
	} else {
		// 删除冲突模板数据
		var tempTemplateIdMap = make(map[string]bool)
		confirmTokenList := strings.Split(confirmToken, ",")
		for _, ct := range confirmTokenList {
			if inputCache, b := models.RequestTemplateImportMap[ct]; b {
				input = inputCache
			} else {
				err = fmt.Errorf("fetch input cache fail,please refersh and try again ")
				return
			}
			tempTemplateIdMap[ct] = true
			delete(models.RequestTemplateImportMap, ct)
			delActions, delErr := s.DeleteRequestTemplate(ct, true)
			if delErr != nil {
				err = delErr
				return
			}
			actions = append(actions, delActions...)
		}
		// 新建模板&模板相关表属性
		input = s.createNewImportTemplate(input, operator, input.RequestTemplate.RecordId, userToken, language)
		// 查询 当前最新模版,记录recordId
		templateName = input.RequestTemplate.Name
		templateList, _ = s.getTemplateListByName(input.RequestTemplate.Name)
		if len(templateList) > 0 {
			for _, template := range templateList {
				if _, ok := tempTemplateIdMap[template.Id]; ok {
					continue
				}
				version := s.getTemplateVersion(template)
				if maxVersion < version {
					maxVersion = version
					maxVersionRecordId = template.Id
				}
			}
			// 有重复数据,但是新导入模板版本最高,直接当成新建处理,需要记录最新发布版本Id
			input.RequestTemplate.RecordId = maxVersionRecordId
		}
	}
	if input.RequestTemplate.Id == "" {
		err = fmt.Errorf("RequestTemplate id illegal ")
		return
	}

	// 判断该平台是否有自选表单entity
	// 自选数据项表单
	var nodesList []*models.DataModel
	var entityMap = make(map[string]bool)
	if nodesList, err = rpc.QueryAllModels(userToken, language); err != nil {
		return
	}
	if len(nodesList) > 0 {
		for _, model := range nodesList {
			if len(model.Entities) > 0 {
				for _, entity := range model.Entities {
					entityStr := fmt.Sprintf("%s:%s", entity.PackageName, entity.Name)
					entityMap[entityStr] = true
				}
			}
		}
		if len(input.FormTemplate) > 0 {
			for _, formTemplate := range input.FormTemplate {
				if formTemplate.ItemGroupType == "optional" && !entityMap[formTemplate.ItemGroup] {
					err = exterror.New().TemplateImportNotMatchEntityError.WithParam(formTemplate.ItemGroupName, formTemplate.ItemGroupName)
					return
				}
			}
		}
	}

	if input.RequestTemplate.ProcDefId != "" {
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
				input.RequestTemplate.ProcDefVersion = v.Version
				manageRole = v.ManageRole
			}
		}
		if !processExistFlag {
			err = exterror.New().TemplateImportNotWorkflowError.WithParam(input.RequestTemplate.ProcDefName, input.RequestTemplate.ProcDefVersion, input.RequestTemplate.ProcDefName, input.RequestTemplate.ProcDefVersion)
			return
		}
		roleFlag := false
		for _, role := range userRoles {
			if role == manageRole {
				roleFlag = true
				break
			}
		}
		if !roleFlag {
			err = exterror.New().TemplateImportNotWorkflowRoleError.WithParam(manageRole)
			return
		}
		// 重新设置模版属主角色
		input.RequestTemplateRole = make([]*models.RequestTemplateRoleTable, 0)
		input.RequestTemplateRole = append(input.RequestTemplateRole, &models.RequestTemplateRoleTable{
			Id:              guid.CreateGuid(),
			RequestTemplate: input.RequestTemplate.Id,
			Role:            manageRole,
			RoleType:        string(models.RolePermissionMGMT),
		})
		input.RequestTemplateRole = append(input.RequestTemplateRole, &models.RequestTemplateRoleTable{
			Id:              guid.CreateGuid(),
			RequestTemplate: input.RequestTemplate.Id,
			Role:            manageRole,
			RoleType:        string(models.RolePermissionUse),
		})

		nodeList, _ := GetProcDefService().GetProcessDefineTaskNodes(models.ConvertRequestTemplateDto2Model(input.RequestTemplate), userToken, language, "template")
		for i, v := range input.TaskTemplate {
			if v.NodeId == "" {
				continue
			}
			existFlag := false
			for _, node := range nodeList {
				if v.NodeId == node.NodeId {
					existFlag = true
					input.TaskTemplate[i].NodeDefId = node.NodeDefId
					break
				}
			}
			if !existFlag {
				err = exterror.New().TemplateImportNotMatchWorkflowTaskError.WithParam(input.RequestTemplate.ProcDefName, input.RequestTemplate.ProcDefVersion, input.RequestTemplate.ProcDefName, input.RequestTemplate.ProcDefVersion)
				return
			}
		}
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	var roleTable []*models.RoleTable
	if roleTable, err = GetRoleService().QueryRoleList(userToken, language); err != nil {
		return
	}
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
	rtAction := dao.ExecAction{Sql: "insert into request_template(id,`group`,name,description,tags,record_id,`version`,confirm_time,status,package_name,entity_name,proc_def_key,proc_def_id,proc_def_name,expire_day,created_by,created_time,updated_by,updated_time,entity_attrs,handler,type,operator_obj_type,approve_by,check_switch,confirm_switch,back_desc,proc_def_version) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
	rtAction.Param = []interface{}{input.RequestTemplate.Id, input.RequestTemplate.Group, input.RequestTemplate.Name, input.RequestTemplate.Description, input.RequestTemplate.Tags, input.RequestTemplate.RecordId, input.RequestTemplate.Version, input.RequestTemplate.ConfirmTime, input.RequestTemplate.Status, input.RequestTemplate.PackageName, input.RequestTemplate.EntityName,
		input.RequestTemplate.ProcDefKey, input.RequestTemplate.ProcDefId, input.RequestTemplate.ProcDefName, input.RequestTemplate.ExpireDay, operator, nowTime, operator, nowTime, input.RequestTemplate.EntityAttrs, input.RequestTemplate.Handler, input.RequestTemplate.Type, input.RequestTemplate.OperatorObjType, input.RequestTemplate.ApproveBy, input.RequestTemplate.CheckSwitch, input.RequestTemplate.ConfirmSwitch, input.RequestTemplate.BackDesc, input.RequestTemplate.ProcDefVersion}
	actions = append(actions, &rtAction)
	for _, v := range input.TaskTemplate {
		tmpAction := dao.ExecAction{Sql: "insert into task_template(id,name,description,request_template,node_id,node_def_id,node_name,expire_day,handler,created_by,created_time,updated_by,updated_time,sort,handle_mode,type) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
		tmpAction.Param = []interface{}{v.Id, v.Name, v.Description, v.RequestTemplate, v.NodeId, v.NodeDefId, v.NodeName, v.ExpireDay, v.Handler, operator, nowTime, operator, nowTime, v.Sort, v.HandleMode, v.Type}
		actions = append(actions, &tmpAction)
	}
	rtRoleFetch := false
	for _, v := range input.RequestTemplateRole {
		if _, b := roleMap[v.Role]; b {
			rtRoleFetch = true
			actions = append(actions, &dao.ExecAction{Sql: "insert into request_template_role(id,request_template,`role`,role_type) values (?,?,?,?)", Param: []interface{}{v.Id, v.RequestTemplate, v.Role, v.RoleType}})
		}
	}
	if !rtRoleFetch {
		actions = append(actions, &dao.ExecAction{Sql: "insert into request_template_role(id,request_template,`role`,role_type) values (?,?,?,?)", Param: []interface{}{guid.CreateGuid() + models.SysTableIdConnector + models.AdminRole + models.SysTableIdConnector + "MGMT", input.RequestTemplate.Id, models.AdminRole, "MGMT"}})
	}
	for _, v := range input.TaskHandleTemplate {
		actions = append(actions, &dao.ExecAction{Sql: "insert into task_handle_template(id,task_template,role,assign,handler_type,handler,handle_mode,sort)values(?,?,?,?,?,?,?,?)", Param: []interface{}{
			v.Id, v.TaskTemplate, v.Role, v.Assign, v.HandlerType, v.Handler, v.HandleMode, v.Sort}})
	}
	for _, v := range input.FormTemplate {
		if v.TaskTemplate == "" {
			actions = append(actions, &dao.ExecAction{Sql: "insert into form_template(id,request_template,item_group,item_group_name,item_group_type,item_group_rule,item_group_sort,created_time,request_form_type,ref_id,del_flag) values(?,?,?,?,?,?,?,?,?,?,?)", Param: []interface{}{
				v.Id, v.RequestTemplate, v.ItemGroup, v.ItemGroupName, v.ItemGroupType, v.ItemGroupRule, v.ItemGroupSort, nowTime, v.RequestFormType, v.RefId, v.DelFlag,
			}})
		} else {
			actions = append(actions, &dao.ExecAction{Sql: "insert into form_template(id,request_template,task_template,item_group,item_group_name,item_group_type,item_group_rule,item_group_sort,created_time,request_form_type,ref_id,del_flag) values(?,?,?,?,?,?,?,?,?,?,?,?)", Param: []interface{}{
				v.Id, v.RequestTemplate, v.TaskTemplate, v.ItemGroup, v.ItemGroupName, v.ItemGroupType, v.ItemGroupRule, v.ItemGroupSort, nowTime, v.RequestFormType, v.RefId, v.DelFlag,
			}})
		}
	}
	for _, v := range input.FormItemTemplate {
		tmpAction := dao.ExecAction{Sql: "insert into form_item_template(id,form_template,name,description,item_group,item_group_name,default_value,sort,package_name,entity,attr_def_id,attr_def_name,attr_def_data_type,element_type,title,width,ref_package_name,ref_entity,data_options,required,regular,is_edit,is_view,is_output,in_display_name,is_ref_inside,multiple,default_clear,ref_id,routine_expression) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
		tmpAction.Param = []interface{}{v.Id, v.FormTemplate, v.Name, v.Description, v.ItemGroup, v.ItemGroupName, v.DefaultValue, v.Sort, v.PackageName, v.Entity, v.AttrDefId, v.AttrDefName, v.AttrDefDataType, v.ElementType, v.Title, v.Width, v.RefPackageName, v.RefEntity, v.DataOptions, v.Required, v.Regular, v.IsEdit, v.IsView, v.IsOutput, v.InDisplayName, v.IsRefInside, v.Multiple, v.DefaultClear, v.RefId, v.RoutineExpression}
		actions = append(actions, &tmpAction)
	}
	err = dao.Transaction(actions)
	return
}

func (s *RequestTemplateService) createNewImportTemplate(input models.RequestTemplateExport, operator, recordId, userToken, language string) models.RequestTemplateExport {
	var newTaskTemplateIdMap = make(map[string]string)
	var newFormTemplateIdMap = make(map[string]string)
	var newFormItemTemplateIdMap = make(map[string]string)
	var historyTemplateId = input.RequestTemplate.Id
	var roleList []*models.SimpleLocalRoleDto
	now := time.Now().Format(models.DateTimeFormat)
	input.RequestTemplate.Id = guid.CreateGuid()
	input.RequestTemplate.RecordId = recordId
	input.RequestTemplate.CreatedBy = operator
	input.RequestTemplate.CreatedTime = now
	input.RequestTemplate.UpdatedBy = operator
	input.RequestTemplate.UpdatedTime = now
	input.RequestTemplate.Handler = operator
	input.RequestTemplate.BackDesc = ""
	// 模版导入,模版使用角色和属主角色取当前操作人角色
	roleList, _ = rpc.QueryUserRoles(operator, userToken, language)
	if len(roleList) > 0 {
		role := roleList[0].Name
		input.RequestTemplateRole = make([]*models.RequestTemplateRoleTable, 0)
		input.RequestTemplateRole = append(input.RequestTemplateRole, &models.RequestTemplateRoleTable{
			Id:              guid.CreateGuid(),
			RequestTemplate: input.RequestTemplate.Id,
			Role:            role,
			RoleType:        string(models.RolePermissionMGMT),
		})
		input.RequestTemplateRole = append(input.RequestTemplateRole, &models.RequestTemplateRoleTable{
			Id:              guid.CreateGuid(),
			RequestTemplate: input.RequestTemplate.Id,
			Role:            role,
			RoleType:        string(models.RolePermissionUse),
		})
	}
	// 修改 taskTemplate中formTemplate,RequestTemplate,以及taskTemplateRole修改
	for _, taskTemplate := range input.TaskTemplate {
		historyTaskTemplateId := taskTemplate.Id
		prefix, _ := GetTaskTemplateService().genTaskIdPrefix(taskTemplate.Type)
		taskTemplate.Id = prefix + "_" + guid.CreateGuid()
		if taskTemplate.RequestTemplate == historyTemplateId {
			taskTemplate.RequestTemplate = input.RequestTemplate.Id
		}
		taskTemplate.CreatedBy = operator
		taskTemplate.UpdatedBy = operator
		taskTemplate.CreatedTime = now
		taskTemplate.UpdatedTime = now
		taskTemplate.Handler = ""
		newTaskTemplateIdMap[historyTaskTemplateId] = taskTemplate.Id
	}
	for _, taskHandleTemplate := range input.TaskHandleTemplate {
		taskHandleTemplate.Id = guid.CreateGuid()
		taskHandleTemplate.TaskTemplate = newTaskTemplateIdMap[taskHandleTemplate.TaskTemplate]
		// 清空角色和处理人
		taskHandleTemplate.Handler = ""
		taskHandleTemplate.Role = ""
	}
	// 修改 formTemplate
	for _, formTemplate := range input.FormTemplate {
		historyFormTemplateId := formTemplate.Id
		formTemplate.Id = guid.CreateGuid()
		newFormTemplateIdMap[historyFormTemplateId] = formTemplate.Id
		formTemplate.RequestTemplate = input.RequestTemplate.Id
		formTemplate.TaskTemplate = newTaskTemplateIdMap[formTemplate.TaskTemplate]
		formTemplate.CreatedTime = now
	}
	for _, formItemTemplate := range input.FormItemTemplate {
		historyFormItemTemplateId := formItemTemplate.Id
		formItemTemplate.Id = guid.CreateGuid()
		formItemTemplate.FormTemplate = newFormTemplateIdMap[formItemTemplate.FormTemplate]
		newFormItemTemplateIdMap[historyFormItemTemplateId] = formItemTemplate.Id
	}

	// 注意 formTemplate 和 formItemTemplate有自引用
	for _, formTemplate := range input.FormTemplate {
		if v, ok := newFormTemplateIdMap[formTemplate.RefId]; ok {
			formTemplate.RefId = v
		}
	}
	for _, formItemTemplate := range input.FormItemTemplate {
		if v, ok := newFormItemTemplateIdMap[formItemTemplate.RefId]; ok {
			formItemTemplate.RefId = v
		}
	}
	return input
}

func (s *RequestTemplateService) getTemplateVersion(template *models.RequestTemplateTable) int {
	var version int
	if template == nil {
		return 0
	}
	if len(template.Version) > 1 {
		version, _ = strconv.Atoi(template.Version[1:])
	}
	return version
}

func (s *RequestTemplateService) getTemplateListByName(templateName string) (requestTemplateTable []*models.RequestTemplateTable, err error) {
	err = dao.X.SQL("select id,name,version,status from request_template where name=?", templateName).Find(&requestTemplateTable)
	return
}

func (s *RequestTemplateService) DisableRequestTemplate(requestTemplateId, operator string) (err error) {
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
	_, err = dao.X.Exec("update request_template set status='disable',updated_by = ?,updated_time =? where id=?", operator, time.Now().Format(models.DateTimeFormat), requestTemplateId)
	return
}

func (s *RequestTemplateService) EnableRequestTemplate(requestTemplateId, operator string) (err error) {
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
	_, err = dao.X.Exec("update request_template set status='confirm',updated_by=?,updated_time=? where id=?", operator, time.Now().Format(models.DateTimeFormat), requestTemplateId)
	return
}

func (s *RequestTemplateService) getAllRequestTemplateGroup() (groupMap map[string]*models.RequestTemplateGroupTable, err error) {
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

func (s *RequestTemplateService) UpdateRequestTemplateParentId(requestTemplate *models.RequestTemplateTable) (parentId string) {
	var actions []*dao.ExecAction
	var templateIds []string
	// 老模板有多个版本,需要更新所有版本,并找到 recordId为空的记录
	requestTemplateMap, _ := s.getAllRequestTemplate()
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

func (s *RequestTemplateService) UpdateRequestTemplateParentIdById(templateId, parentId string) (err error) {
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

func (s *RequestTemplateService) CheckPermission(requestTemplateId, user string) (err error) {
	var requestTemplate *models.RequestTemplateTable
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(requestTemplateId)
	if err != nil {
		return
	}
	if requestTemplate == nil {
		err = exterror.Catch(exterror.New().RequestParamValidateError, fmt.Errorf("param requestTemplateId is invalid"))
		return
	}
	// 请求模板的更新人不是当前用户,不允许操作
	if requestTemplate.UpdatedBy != user {
		err = exterror.New().DataPermissionDeny
	}
	return
}

func (s *RequestTemplateService) GetRequestPendingRoleAndHandler(requestTemplate *models.RequestTemplateDto) (role, handler string) {
	if requestTemplate == nil {
		return
	}
	if !requestTemplate.CheckSwitch {
		return
	}

	// 配置了定版角色
	if requestTemplate.CheckRole != "" {
		role = requestTemplate.CheckRole
		handler = requestTemplate.CheckHandler
	} else {
		// 没有配置定版角色和定版人则读取模板属主角色和属主处理人
		role = s.GetRequestTemplateManageRole(requestTemplate.Id)
		handler = requestTemplate.Handler
	}
	return
}

func (s *RequestTemplateService) GetConfirmCount(user, userToken, language string) (count int, err error) {
	var requestTemplateList []*models.RequestTemplateTable
	var roleMap map[string]*models.SimpleLocalRoleDto
	roleMap, err = rpc.QueryAllRoles("Y", userToken, language)
	if err != nil {
		return
	}

	if err = dao.X.SQL("select * from request_template where status = ? and del_flag=0 ", models.RequestTemplateStatusPending).Find(&requestTemplateList); err != nil {
		return
	}
	if len(requestTemplateList) > 0 {
		for _, requestTemplate := range requestTemplateList {
			var requestTemplateRoleList []*models.RequestTemplateRoleTable
			err = dao.X.SQL("select * from request_template_role where request_template=? and role_type=?", requestTemplate.Id, models.RolePermissionMGMT).Find(&requestTemplateRoleList)
			if err != nil {
				return
			}
			if len(requestTemplateRoleList) > 0 {
				if v, ok := roleMap[requestTemplateRoleList[0].Role]; ok {
					if v.Administrator == user {
						count++
					}
				}
			}
		}
	}
	return
}

func (s *RequestTemplateService) QueryListByName(name string) (list []*models.RequestTemplateTable, err error) {
	return s.requestTemplateDao.QueryListByName(name)
}

func (s *RequestTemplateService) CreateWorkflowFormTemplate(session *xorm.Session, requestTemplateId, procDefId, userToken, language string) (err error) {
	var procDefEntities []*models.ProcEntity
	procDefEntities, err = s.ListRequestTemplateEntityAttrs(requestTemplateId, procDefId, userToken, language)
	if err != nil {
		return
	}
	if len(procDefEntities) > 0 {
		for index, entity := range procDefEntities {
			entityStr := fmt.Sprintf("%s:%s", entity.PackageName, entity.Name)
			_, err = session.Exec("insert into form_template(id,request_template,item_group,item_group_name,item_group_type,item_group_rule,item_group_sort,"+
				"created_time,request_form_type) values(?,?,?,?,?,?,?,?,?)", guid.CreateGuid(), requestTemplateId, entityStr, entityStr, models.FormItemGroupTypeWorkflow,
				"exist", index+1, time.Now().Format(models.DateTimeFormat), models.RequestFormTypeData)
			if err != nil {
				return
			}
		}
	}
	return
}

func (s *RequestTemplateService) CreateWorkflowFormTemplateSql(requestTemplateId, procDefId, userToken, language string) (actions []*dao.ExecAction, err error) {
	var procDefEntities []*models.ProcEntity
	actions = []*dao.ExecAction{}
	procDefEntities, err = s.ListRequestTemplateEntityAttrs(requestTemplateId, procDefId, userToken, language)
	if err != nil {
		return
	}
	if len(procDefEntities) > 0 {
		for index, entity := range procDefEntities {
			entityStr := fmt.Sprintf("%s:%s", entity.PackageName, entity.Name)
			actions = append(actions, &dao.ExecAction{Sql: "insert into form_template(id,request_template,item_group,item_group_name,item_group_type,item_group_rule," +
				"item_group_sort,created_time,request_form_type) values(?,?,?,?,?,?,?,?,?)", Param: []interface{}{guid.CreateGuid(), requestTemplateId, entityStr, entityStr,
				models.FormItemGroupTypeWorkflow, "exist", index + 1, time.Now().Format(models.DateTimeFormat), models.RequestFormTypeData}})
		}
	}
	return
}
