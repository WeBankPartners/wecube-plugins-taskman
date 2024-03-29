package service

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/rpc"
	"github.com/tealeg/xlsx"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	WaitCommit           = "waitCommit"           // 等待提交
	SendRequest          = "sendRequest"          // 发送请求
	RequestPending       = "requestPending"       // 请求定版
	Approval             = "approval"             // 审批
	Task                 = "task"                 // 任务
	Confirm              = "confirm"              // 请求确认
	CurNodeCompleted     = "Completed"            // 完成
	RequestComplete      = "requestComplete"      // 请求完成
	AutoExit             = "autoExit"             // 自动退出
	InternallyTerminated = "internallyTerminated" // 手动终止
	AutoNode             = "autoNode"             // 自动节点
)

// sceneTypeMap 场景数组
var sceneTypeMap = map[models.PlatTab]models.SceneType{
	models.PlatTabRequest:       models.SceneTypeRequest,
	models.PlatTabRelease:       models.SceneTypeRelease,
	models.PlatTabProblem:       models.SceneTypeProblem,
	models.PlatTabProblemEvent:  models.SceneTypeEvent,
	models.PlatTabProblemChange: models.SceneTypeChange,
}

func getTaskTypeByType(uiType int) models.TaskType {
	// 0 所有,1表示请求定版,2 任务处理,3 审批 4确认请求
	switch uiType {
	case 1:
		return models.TaskTypeCheck
	case 2:
		return models.TaskTypeImplement
	case 3:
		return models.TaskTypeApprove
	case 4:
		return models.TaskTypeConfirm
	}
	return models.TaskTypeCheck
}

// GetPlatformCount 工作台数量统计
func GetPlatformCount(param models.CountPlatformParam, user string, userRoles []string) (platformData models.PlatformData, err error) {
	switch param.Tab {
	case "all":
		platformData.Pending = GetPendingCount(param, userRoles)
		platformData.MyPending = GetMyPendingCount(param, user, userRoles)
		platformData.HasProcessed = GetHasProcessedCount(param.QueryTimeStart, param.QueryTimeEnd, user)
		platformData.Submit = GetSubmitCount(param.QueryTimeStart, param.QueryTimeEnd, user)
		platformData.Draft = GetDraftCount(param.QueryTimeStart, param.QueryTimeEnd, user)
	case "pending":
		platformData.Pending = GetPendingCount(param, userRoles)
		platformData.MyPending = GetMyPendingCount(param, user, userRoles)
	case "myPending":
		platformData.Pending = GetPendingCount(param, userRoles)
		platformData.MyPending = GetMyPendingCount(param, user, userRoles)
	case "hasProcessed":
		platformData.HasProcessed = GetHasProcessedCount(param.QueryTimeStart, param.QueryTimeEnd, user)
	case "submit":
		platformData.Submit = GetSubmitCount(param.QueryTimeStart, param.QueryTimeEnd, user)
	case "draft":
		platformData.Draft = GetDraftCount(param.QueryTimeStart, param.QueryTimeEnd, user)
	}
	return
}

func GetSimpleRequest(requestId string) (request models.RequestTable, err error) {
	var requestTable []*models.RequestTable
	err = dao.X.SQL("select * from request where id=?", requestId).Find(&requestTable)
	if err != nil {
		return
	}
	if len(requestTable) == 0 {
		return request, fmt.Errorf("Can not find any request with id:%s ", requestId)
	}
	request = *requestTable[0]
	return
}

// UpdateRequestHandler 请求/发布 认领&转给我逻辑
func UpdateRequestHandler(requestId, user string) (err error) {
	var actions []*dao.ExecAction
	nowTime := time.Now().Format(models.DateTimeFormat)
	actions = append(actions, &dao.ExecAction{Sql: "update request set handler= ?,updated_by= ?,updated_time= ? where id= ?",
		Param: []interface{}{user, user, nowTime, requestId}})
	err = dao.Transaction(actions)
	return
}

// GetPendingCount 统计待处理,包括:请求、发布,以及下面的请求提交、任务、审批、请求定版、请求确认
func GetPendingCount(param models.CountPlatformParam, userRoles []string) (resultMap map[string]int) {
	var pendingTask, pendingApprove, pendingCheck, pendingConfirm, pending int
	resultMap = make(map[string]int)
	for sceneName, sceneType := range sceneTypeMap {
		pendingTask, pendingApprove, pendingCheck, pendingConfirm, pending = GetPendingCountByScene(param, int(sceneType), userRoles)
		resultMap[string(sceneName)] = pending
		resultMap[string(sceneName)+string(models.PlatTabProblemApprove)] = pendingApprove
		resultMap[string(sceneName)+string(models.PlatTabProblemTask)] = pendingTask
		resultMap[string(sceneName)+string(models.PlatTabProblemCheck)] = pendingCheck
		resultMap[string(sceneName)+string(models.PlatTabProblemConfirm)] = pendingConfirm
	}
	return
}

// GetMyPendingCount 统计本人待处理
func GetMyPendingCount(param models.CountPlatformParam, user string, userRoles []string) (resultMap map[string]int) {
	var pendingTask, pendingApprove, pendingCheck, pendingConfirm, pending int
	resultMap = make(map[string]int)
	for sceneName, sceneType := range sceneTypeMap {
		pendingTask, pendingApprove, pendingCheck, pendingConfirm, pending = GetMyPendingCountByScene(param, int(sceneType), user, userRoles)
		resultMap[string(sceneName)] = pending
		resultMap[string(sceneName)+string(models.PlatTabProblemApprove)] = pendingApprove
		resultMap[string(sceneName)+string(models.PlatTabProblemTask)] = pendingTask
		resultMap[string(sceneName)+string(models.PlatTabProblemCheck)] = pendingCheck
		resultMap[string(sceneName)+string(models.PlatTabProblemConfirm)] = pendingConfirm
	}
	return
}

// GetPendingCountByScene 统计本组待处理,包括:请求、发布,以及下面的请求提交、任务、审批、请求定版、请求确认
func GetPendingCountByScene(param models.CountPlatformParam, scene int, userRoles []string) (pendingTask, pendingApprove, pendingCheck, pendingConfirm, pending int) {
	userRolesFilterSql, userRolesFilterParam := dao.CreateListParams(userRoles, "")
	var pendingTaskParam, pendingApproveParam, pendingCheckParam, pendingConfirmParam []interface{}
	var pTaskSQL, pendingApproveSQL, pendingCheckSQL, pendingConfirmSQL string
	pTaskSQL, pendingTaskParam = pendingTaskSQL(param.QueryTimeStart, param.QueryTimeEnd, scene, userRolesFilterSql, userRolesFilterParam, models.TaskTypeImplement)
	pendingTask = dao.QueryCount(pTaskSQL, pendingTaskParam...)

	pendingApproveSQL, pendingApproveParam = pendingTaskSQL(param.QueryTimeStart, param.QueryTimeEnd, scene, userRolesFilterSql, userRolesFilterParam, models.TaskTypeApprove)
	pendingApprove = dao.QueryCount(pendingApproveSQL, pendingApproveParam...)

	pendingCheckSQL, pendingCheckParam = pendingTaskSQL(param.QueryTimeStart, param.QueryTimeEnd, scene, userRolesFilterSql, userRolesFilterParam, models.TaskTypeCheck)
	pendingCheck = dao.QueryCount(pendingCheckSQL, pendingCheckParam...)

	pendingConfirmSQL, pendingConfirmParam = pendingTaskSQL(param.QueryTimeStart, param.QueryTimeEnd, scene, userRolesFilterSql, userRolesFilterParam, models.TaskTypeConfirm)
	pendingConfirm = dao.QueryCount(pendingConfirmSQL, pendingConfirmParam...)

	pending = pendingTask + pendingApprove + pendingCheck + pendingConfirm
	return
}

// GetMyPendingCountByScene 根据场景查询我的待处理
func GetMyPendingCountByScene(param models.CountPlatformParam, scene int, user string, userRoles []string) (pendingTask, pendingApprove, pendingCheck, pendingConfirm, pending int) {
	userRolesFilterSql, userRolesFilterParam := dao.CreateListParams(userRoles, "")
	var pendingTaskParam, pendingApproveParam, pendingCheckParam, pendingConfirmParam []interface{}
	var pTaskSQL, pendingApproveSQL, pendingCheckSQL, pendingConfirmSQL string
	pTaskSQL, pendingTaskParam = pendingMyTaskSQL(param.QueryTimeStart, param.QueryTimeEnd, scene, user, userRolesFilterSql, userRolesFilterParam, models.TaskTypeImplement)
	pendingTask = dao.QueryCount(pTaskSQL, pendingTaskParam...)

	pendingApproveSQL, pendingApproveParam = pendingMyTaskSQL(param.QueryTimeStart, param.QueryTimeEnd, scene, user, userRolesFilterSql, userRolesFilterParam, models.TaskTypeApprove)
	pendingApprove = dao.QueryCount(pendingApproveSQL, pendingApproveParam...)

	pendingCheckSQL, pendingCheckParam = pendingMyTaskSQL(param.QueryTimeStart, param.QueryTimeEnd, scene, user, userRolesFilterSql, userRolesFilterParam, models.TaskTypeCheck)
	pendingCheck = dao.QueryCount(pendingCheckSQL, pendingCheckParam...)

	pendingConfirmSQL, pendingConfirmParam = pendingMyTaskSQL(param.QueryTimeStart, param.QueryTimeEnd, scene, user, userRolesFilterSql, userRolesFilterParam, models.TaskTypeConfirm)
	pendingConfirm = dao.QueryCount(pendingConfirmSQL, pendingConfirmParam...)

	pending = pendingTask + pendingApprove + pendingCheck + pendingConfirm
	return
}

func getPlatRequestSQL(where, sql string) string {
	return fmt.Sprintf("select * from (select r.id,r.name,r.cache,r.report_time,r.del_flag,rt.id as template_id,rt.name as template_name,rt.parent_id,"+
		"r.proc_instance_id,r.operator_obj,rt.proc_def_id,r.type as type,rt.proc_def_key,rt.operator_obj_type,r.role,r.status,r.rollback_desc,r.created_by,r.handler,r.created_time,r.updated_time,rt.proc_def_name,rt.proc_def_version,"+
		"r.expect_time,r.revoke_flag,r.confirm_time as approval_time,rt.expire_day from request r join request_template rt on r.request_template = rt.id ) t %s and id in (%s) ", where, sql)
}

func getPlatTaskSQL(where, sql string) string {
	return fmt.Sprintf("select * from (select r.id,r.name,r.cache,r.report_time,r.del_flag,rt.id as template_id,rt.name as template_name,rt.parent_id,r.proc_instance_id,r.operator_obj,rt.proc_def_id,r.type as type,rt.proc_def_key,rt.operator_obj_type,r.role,r.status,r.rollback_desc,r.created_by,r.created_time,r.updated_time,rt.proc_def_name,r.expect_time,r.revoke_flag,t.id as task_id,t.name as task_name,t.task_handle_role,t.task_created_time,t.task_approval_time as task_approval_time,t.updated_time as task_updated_time,t.status as task_status,t.expire_time as task_expect_time,t.task_handler as task_handler,t.task_handle_id,t.task_handle_created_time,t.task_handle_updated_time from (%s) t left join request r on t.request=r.id join request_template rt on r.request_template = rt.id) temp %s", sql, where)
}

func pendingTaskSQL(queryTimeStart, queryTimeEnd string, templateType int, userRolesFilterSql string, userRolesFilterParam []interface{}, taskType models.TaskType) (sql string, queryParam []interface{}) {
	queryParam = []interface{}{}
	if taskType == models.TaskTypeNone {
		sql = "select * from (select t.id,t.request,t.template_type,t.name,t.type,t.created_time as task_created_time,th.updated_time as task_approval_time,t.updated_time,t.status,t.expire_time,th.role as task_handle_role,th.id as task_handle_id,th.handler as task_handler,t.del_flag,th.latest_flag,th.handle_status,th.handle_result,th.created_time as task_handle_created_time,th.updated_time as task_handle_updated_time,t.request_created_time from task t right join task_handle th ON t.id = th.task) tha where del_flag = 0 and status <> 'done' and template_type = ? and latest_flag = 1 and handle_result is null and request_created_time >= ? and request_created_time <= ? and task_handle_role in (" + userRolesFilterSql + ")"
		queryParam = append([]interface{}{templateType, queryTimeStart, queryTimeEnd}, userRolesFilterParam...)
	} else {
		sql = "select * from (select t.id,t.request,t.template_type,t.name,t.type,t.created_time as task_created_time,th.updated_time as task_approval_time,t.updated_time,t.status,t.expire_time,th.role as task_handle_role,th.id as task_handle_id,th.handler as task_handler,t.del_flag,th.latest_flag,th.handle_status,th.handle_result,th.created_time as task_handle_created_time,th.updated_time as task_handle_updated_time,t.request_created_time from task t right join task_handle th ON t.id = th.task) tha where del_flag = 0 and status <> 'done' and template_type = ? and type = ? and latest_flag = 1 and handle_result is null and request_created_time >= ? and request_created_time <= ?  and task_handle_role in (" + userRolesFilterSql + ")"
		queryParam = append([]interface{}{templateType, taskType, queryTimeStart, queryTimeEnd}, userRolesFilterParam...)
	}
	return
}

func pendingMyTaskSQL(queryTimeStart, queryTimeEnd string, templateType int, user, userRolesFilterSql string, userRolesFilterParam []interface{}, taskType models.TaskType) (sql string, queryParam []interface{}) {
	queryParam = []interface{}{}
	sql = "select * from (select t.id,t.request,t.template_type,t.name,t.type,t.created_time as task_created_time,th.updated_time as task_approval_time,t.updated_time,t.status,t.expire_time,th.role as task_handle_role,th.id as task_handle_id,th.handler as task_handler,t.del_flag,th.latest_flag,th.handle_status,th.handle_result,th.created_time as task_handle_created_time,th.updated_time as task_handle_updated_time,t.request_created_time from task t right join task_handle th ON t.id = th.task) tha where del_flag = 0 and status <> 'done' and template_type = ? and type = ? and latest_flag = 1 and handle_result is null and (task_handler =? or task_handler is null) and request_created_time >= ? and request_created_time <= ? and task_handle_role in (" + userRolesFilterSql + ")"
	queryParam = append([]interface{}{templateType, taskType, user, queryTimeStart, queryTimeEnd}, userRolesFilterParam...)
	return
}

func hasProcessedTaskSQL(queryTimeStart, queryTimeEnd string, templateType int, user string, taskType models.TaskType) (sql string, queryParam []interface{}) {
	queryParam = []interface{}{}
	if taskType == models.TaskTypeNone {
		sql = "select * from (select t.id,t.request,t.template_type,t.name,t.type,t.created_time as task_created_time,th.updated_time as task_approval_time,t.updated_time,t.status,t.expire_time,th.role as task_handle_role,th.id as task_handle_id,th.handler as task_handler,t.del_flag,th.latest_flag,th.handle_status,th.handle_result,th.created_time as task_handle_created_time,th.updated_time as task_handle_updated_time,t.request_created_time from task t right join task_handle th ON t.id = th.task) tha where del_flag = 0 and handle_result is not null and template_type = ? and type != ? and task_handler =? and latest_flag = 1 and request_created_time >= ? and request_created_time <= ?"
		queryParam = append([]interface{}{templateType, models.TaskTypeSubmit, user, queryTimeStart, queryTimeEnd})
	} else {
		sql = "select * from (select t.id,t.request,t.template_type,t.name,t.type,t.created_time as task_created_time,th.updated_time as task_approval_time,t.updated_time,t.status,t.expire_time,th.role as task_handle_role,th.id as task_handle_id,th.handler as task_handler,t.del_flag,th.latest_flag,th.handle_status,th.handle_result,th.created_time as task_handle_created_time,th.updated_time as task_handle_updated_time,t.request_created_time from task t right join task_handle th ON t.id = th.task) tha where del_flag = 0 and handle_result is not null and template_type = ? and type = ? and task_handler =? and latest_flag = 1  and request_created_time >= ? and request_created_time <= ?"
		queryParam = append([]interface{}{templateType, taskType, user, queryTimeStart, queryTimeEnd})
	}
	return
}

// GetHasProcessedCount 统计已处理,包括:(1)处理定版 (2) 任务已审批
func GetHasProcessedCount(queryTimeStart, queryTimeEnd string, user string) map[string]int {
	var resultMap = make(map[string]int)
	for sceneName, sceneType := range sceneTypeMap {
		hasProcessedSQL, hasProcessedParam := hasProcessedTaskSQL(queryTimeStart, queryTimeEnd, int(sceneType), user, models.TaskTypeNone)
		resultMap[string(sceneName)] = dao.QueryCount(hasProcessedSQL, hasProcessedParam...)
	}
	return resultMap
}

// GetSubmitCount  统计用户提交
func GetSubmitCount(queryTimeStart, queryTimeEnd, user string) map[string]int {
	var resultMap = make(map[string]int)
	var sql string
	var queryParam []interface{}
	var count int
	for sceneName, sceneType := range sceneTypeMap {
		sql, queryParam = submitSQL(0, int(sceneType), queryTimeStart, queryTimeEnd, user)
		count = dao.QueryCount(sql, queryParam...)
		resultMap[string(sceneName)] = count
	}
	return resultMap
}

func submitSQL(rollback, templateType int, queryTimeStart, queryTimeEnd, user string) (sql string, queryParam []interface{}) {
	sql = "select id from request where del_flag=0 and created_by = ? and type = ? and created_time >= ? and created_time <= ? and (status != 'Draft' or ( status = 'Draft' and rollback_desc is not null ) or (status = 'Draft' and revoke_flag = 1))"
	if rollback == 1 {
		// 被退回
		sql = "select id from request where del_flag=0 and created_by = ? and type = ? and status = 'Draft' and rollback_desc is not null and created_time >= ? and created_time <= ?"
	} else if rollback == 2 {
		// 其他
		sql = "select id from request where del_flag=0 and created_by = ? and type = ? and status != 'Draft' and created_time >= ? and created_time <= ?"
	} else if rollback == 3 {
		// 撤销
		sql = "select id from request where del_flag=0 and created_by = ? and type = ? and status = 'Draft' and revoke_flag = 1 and created_time >= ? and created_time <= ?"
	}
	queryParam = append([]interface{}{user, templateType, queryTimeStart, queryTimeEnd})
	return
}

// GetDraftCount 统计用户暂存
func GetDraftCount(queryTimeStart, queryTimeEnd, user string) map[string]int {
	var resultMap = make(map[string]int)
	var sql string
	var queryParam []interface{}
	var count int
	for sceneName, sceneType := range sceneTypeMap {
		sql, queryParam = draftSQL(int(sceneType), queryTimeStart, queryTimeEnd, user)
		count = dao.QueryCount(sql, queryParam...)
		resultMap[string(sceneName)] = count
	}
	return resultMap
}

func draftSQL(templateType int, queryTimeStart, queryTimeEnd, user string) (sql string, queryParam []interface{}) {
	sql = "select id from request where del_flag=0 and created_by = ? and status = 'Draft' and type = ? and rollback_desc is null and revoke_flag = 0 and created_time >= ? and created_time <= ?"
	queryParam = append([]interface{}{user, templateType, queryTimeStart, queryTimeEnd})
	return
}

// DataList 首页工作台数据列表
func DataList(param *models.PlatformRequestParam, userRoles []string, userToken, user, language string) (pageInfo models.PageInfo, rowData []*models.PlatformDataObj, err error) {
	// 先拼接查询条件
	var sql, execSql string
	var queryParam []interface{}
	var taskType = getTaskTypeByType(param.Type)
	where := transPlatConditionToSQL(param)
	userRolesFilterSql, userRolesFilterParam := dao.CreateListParams(userRoles, "")
	switch param.Tab {
	case "myPending":
		sql, queryParam = pendingMyTaskSQL(param.QueryTimeStart, param.QueryTimeEnd, param.Action, user, userRolesFilterSql, userRolesFilterParam, taskType)
		execSql = getPlatTaskSQL(where, sql)
	case "pending":
		sql, queryParam = pendingTaskSQL(param.QueryTimeStart, param.QueryTimeEnd, param.Action, userRolesFilterSql, userRolesFilterParam, taskType)
		execSql = getPlatTaskSQL(where, sql)
	case "hasProcessed":
		sql, queryParam = hasProcessedTaskSQL(param.QueryTimeStart, param.QueryTimeEnd, param.Action, user, taskType)
		execSql = getPlatTaskSQL(where, sql)
	case "submit":
		sql, queryParam = submitSQL(param.Rollback, param.Action, param.QueryTimeStart, param.QueryTimeEnd, user)
		execSql = getPlatRequestSQL(where, sql)
	case "draft":
		sql, queryParam = draftSQL(param.Action, param.QueryTimeStart, param.QueryTimeEnd, user)
		execSql = getPlatRequestSQL(where, sql)
	default:
		err = fmt.Errorf("request param err,tab:%s", param.Tab)
		return
	}
	pageInfo, rowData, err = getPlatData(models.PlatDataParam{
		Param:      param.CommonRequestParam,
		QueryParam: queryParam,
		User:       user,
		UserToken:  userToken,
		Tab:        param.Tab,
	}, execSql, language, true)
	return
}

// HistoryList 发布历史
func HistoryList(param *models.RequestHistoryParam, userRoles []string, userToken, user, language string) (pageInfo models.PageInfo, rowsData []*models.PlatformDataObj, err error) {
	var sql = "select id from request"
	var queryParam []interface{}
	where := transHistoryConditionToSQL(param)
	// 查看本组数据
	if param.Permission == "group" {
		userRolesFilterSql, userRolesFilterParam := dao.CreateListParams(userRoles, "")
		sql = "select id from request where  `role` in (" + userRolesFilterSql + ")"
		queryParam = append(queryParam, userRolesFilterParam...)
	}
	pageInfo, rowsData, err = getPlatData(models.PlatDataParam{Param: param.CommonRequestParam, QueryParam: queryParam, User: user, UserToken: userToken}, getPlatRequestSQL(where, sql), language, true)
	return
}

// Export 数据导出
func Export(w http.ResponseWriter, param *models.RequestHistoryParam, userToken, language, user string) (err error) {
	var rowsData []*models.PlatformDataObj
	var sql = "select id from request"
	var queryParam []interface{}
	where := transHistoryConditionToSQL(param)
	_, rowsData, err = getPlatData(models.PlatDataParam{Param: param.CommonRequestParam, QueryParam: queryParam, User: user, UserToken: userToken}, getPlatRequestSQL(where, sql), language, false)
	if len(rowsData) > 0 {
		fileName, titles, dataArr := getRequestExportData(language, rowsData)
		file := xlsx.NewFile()
		sheet, _ := file.AddSheet("Sheet1")
		row := sheet.AddRow()

		var cell *xlsx.Cell
		// 列标题赋值
		for _, title := range titles {
			cell = row.AddCell()
			cell.Value = title
		}

		for _, v := range dataArr {
			row = sheet.AddRow()
			for _, vv := range v {
				cell = row.AddCell()
				cell.Value = vv
			}
		}
		// 设置 HTTP 响应头
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Transfer-Encoding", "binary")
		err = file.Write(w)
	}
	return
}

func getRequestExportData(language string, rowsData []*models.PlatformDataObj) (fileName string, titles []string, dataArr [][]string) {
	dataArr = make([][]string, len(rowsData))
	var days string
	var version string
	fileName = "Wecube-RequestAudit-Export-" + time.Now().Format("20060102150405") + ".xlsx"
	if strings.Contains(language, "zh-CN") {
		titles = []string{
			"请求ID",
			"请求名称",
			"请求状态",
			"当前节点",
			"当前处理人",
			"当前处理人角色",
			"进展",
			"请求停留时长",
			"期望完成时间",
			"使用模板",
			"使用编排",
			"操作对象类型",
			"操作对象",
			"创建人",
			"创建人角色",
			"请求提交时间",
		}
		days = "日"
	} else {
		titles = []string{
			"Request ID",
			"Request Name",
			"Request Status",
			"Current Node",
			"Current Approver",
			"Current Approver Role",
			"progress",
			"Request Stay Duration",
			"Expected Completion",
			"Using Template",
			"Using Process",
			"Target Object Type",
			"Operation Object",
			"Creator",
			"Creator Role",
			"Request Submission",
		}
		days = "days"
	}
	for i, row := range rowsData {
		version = row.Version
		if version == "" {
			version = "beta"
		}
		dataArr[i] = []string{
			row.Id,
			row.Name,
			getInternationalizationStatus(language, row.Status),
			getInternationalizationCurNode(language, row.CurNode),
			row.Handler,
			row.HandleRoleDisplay,
			strconv.Itoa(row.Progress) + "%",
			getRequestStayTime(row, days),
			row.ExpectTime,
			row.TemplateName + "【" + version + "】",
			row.ProcDefName,
			row.OperatorObjType,
			row.OperatorObj,
			row.CreatedBy,
			row.RoleDisplay,
			row.ReportTime,
		}
	}
	return
}

func getInternationalizationCurNode(language, node string) string {
	if strings.Contains(language, "zh-CN") {
		switch node {
		case CurNodeCompleted:
			return "请求完成"
		case WaitCommit:
			return "等待提交"
		case SendRequest:
			return "发送请求"
		case RequestPending:
			return "请求定版"
		case RequestComplete:
			return "请求完成"
		case Confirm:
			return "请求确认"
		}
	}
	return node
}

func getInternationalizationStatus(language string, status string) string {
	if strings.Contains(language, "zh-CN") {
		switch status {
		case string(models.RequestStatusDraft):
			return "草稿"
		case string(models.RequestStatusPending):
			return "等待定版"
		case string(models.RequestStatusInProgress):
			return "执行中"
		case string(models.RequestStatusInProgressFaulted):
			return "节点报错"
		case string(models.RequestStatusTermination):
			return "手动终止"
		case string(models.RequestStatusCompleted):
			return "成功"
		case string(models.RequestStatusInApproval):
			return "审批中"
		case string(models.RequestStatusConfirm):
			return "请求确认"
		case string(models.RequestStatusInProgressTimeOuted):
			return "节点超时"
		case string(models.RequestStatusFaulted):
			return "自动退出"
		}
	}
	return status
}

func getRequestStayTime(dataObject *models.PlatformDataObj, format string) string {
	return fmt.Sprintf("%s%s/%d%s", dataObject.RequestStayTime, format, dataObject.RequestStayTimeTotal, format)
}

// calcRequestStayTime 计算请求/任务停留时长
func calcRequestStayTime(dataObject *models.PlatformDataObj) {
	var err error
	var reportTime, requestExpectTime, taskCreateTime, taskExpectTime, taskApprovalTime time.Time
	loc, _ := time.LoadLocation("Local")
	// 设置个默认值
	dataObject.RequestStayTime = "0.0"
	dataObject.TaskStayTime = "0.0"
	// 计算任务停留时长
	if dataObject.TaskId != "" && dataObject.TaskExpectTime != "" && dataObject.TaskCreatedTime != "" {
		taskExpectTime, _ = time.ParseInLocation(models.DateTimeFormat, dataObject.TaskExpectTime, loc)
		taskCreateTime, _ = time.ParseInLocation(models.DateTimeFormat, dataObject.TaskCreatedTime, loc)
		if dataObject.TaskApprovalTime != "" && dataObject.TaskStatus == "done" && dataObject.Status != string(models.RequestStatusDraft) {
			taskApprovalTime, _ = time.ParseInLocation(models.DateTimeFormat, dataObject.TaskApprovalTime, loc)
			dataObject.TaskStayTime = fmt.Sprintf("%.1f", taskApprovalTime.Sub(taskCreateTime).Hours()*1.00/24.00)
		} else {
			dataObject.TaskStayTime = fmt.Sprintf("%.1f", time.Now().Local().Sub(taskCreateTime).Hours()*1.00/24.00)
		}
		dataObject.TaskStayTimeTotal = int(math.Ceil(taskExpectTime.Sub(taskCreateTime).Hours() * 1.00 / 24.00))
	}
	if dataObject.Status == string(models.RequestStatusDraft) {
		dataObject.TaskStayTime = "0.0"
		dataObject.RequestStayTimeTotal = dataObject.ExpireDay
		return
	}
	reportTime, err = time.ParseInLocation(models.DateTimeFormat, dataObject.ReportTime, loc)
	if err != nil {
		log.Logger.Error("getRequestRemainDays ReportTime err", log.Error(err))
		return
	}
	requestExpectTime, err = time.ParseInLocation(models.DateTimeFormat, dataObject.ExpectTime, loc)
	if err != nil {
		log.Logger.Error("getRequestRemainDays ExpectTime err", log.Error(err))
		return
	}
	// 计算请求停留时长
	if dataObject.Status == string(models.RequestStatusCompleted) || dataObject.Status == string(models.RequestStatusTermination) || dataObject.Status == string(models.RequestStatusFaulted) {
		updateTime, err := time.ParseInLocation(models.DateTimeFormat, dataObject.UpdatedTime, loc)
		if err != nil {
			log.Logger.Error("getRequestRemainDays UpdatedTime err", log.Error(err))
			return
		}

		dataObject.RequestStayTime = fmt.Sprintf("%.1f", updateTime.Sub(reportTime).Hours()*1.00/24.00)
	} else {
		dataObject.RequestStayTime = fmt.Sprintf("%.1f", time.Now().Local().Sub(reportTime).Hours()*1.00/24.00)
	}
	// 向上取整
	dataObject.RequestStayTimeTotal = int(math.Ceil(requestExpectTime.Sub(reportTime).Hours() * 1.00 / 24.00))
}

func getPlatData(req models.PlatDataParam, newSQL, language string, page bool) (pageInfo models.PageInfo, rowsData []*models.PlatformDataObj, err error) {
	var operatorObjTypeMap = make(map[string]string)
	var roleDtoMap map[string]*models.SimpleLocalRoleDto
	var roleDisplayMap = make(map[string]string)
	// 请求已处理(防止同一个请求重复处理)
	var processedRequestMap = make(map[string]bool)
	// 排序处理
	if req.Param.Sorting != nil {
		hashMap, _ := dao.GetJsonToXormMap(models.PlatformDataObj{})
		if len(hashMap) > 0 {
			if req.Param.Sorting.Asc {
				newSQL += fmt.Sprintf(" ORDER BY %s ASC ", hashMap[req.Param.Sorting.Field])
			} else {
				newSQL += fmt.Sprintf(" ORDER BY %s DESC ", hashMap[req.Param.Sorting.Field])
			}
		}
	}
	roleDisplayMap, _ = GetRoleService().GetRoleDisplayName(req.UserToken, language)
	// 分页处理
	if page {
		pageInfo.StartIndex = req.Param.StartIndex
		pageInfo.PageSize = req.Param.PageSize
		pageInfo.TotalRows = dao.QueryCount(newSQL, req.QueryParam...)
		pageSQL := newSQL + " limit ?,? "
		req.QueryParam = append(req.QueryParam, req.Param.StartIndex, req.Param.PageSize)
		err = dao.X.SQL(pageSQL, req.QueryParam...).Find(&rowsData)
	} else {
		err = dao.X.SQL(newSQL, req.QueryParam...).Find(&rowsData)
	}
	if len(rowsData) > 0 {
		// 操作对象类型,新增模板是录入.历史模板操作对象类型为空,需要全量处理下
		operatorObjTypeMap = GetRequestTemplateService().GetAllCoreProcess(req.UserToken, language)
		// 查询当前用户所有收藏模板记录
		collectMap, _ := QueryAllTemplateCollect(req.User)
		templateMap, _ := GetRequestTemplateService().getAllRequestTemplate()
		// 查询所有角色管理员
		if roleDtoMap, _ = rpc.QueryAllRoles("Y", req.UserToken, language); len(roleDtoMap) == 0 {
			roleDtoMap = make(map[string]*models.SimpleLocalRoleDto)
		}
		var actions, confirmActions []*dao.ExecAction
		for _, platformDataObj := range rowsData {
			// 获取 使用编排
			if len(templateMap) > 0 && templateMap[platformDataObj.TemplateId] != nil {
				template := templateMap[platformDataObj.TemplateId]
				platformDataObj.ProcDefName = template.ProcDefName
				platformDataObj.TemplateName = template.Name
				if template.Status != "confirm" {
					platformDataObj.Version = "beta"
				} else {
					platformDataObj.Version = template.Version
				}
			}
			// 获取当前节点和计算请求进度
			platformDataObj.Progress, platformDataObj.CurNode = CalcRequestProgressAndCurNode(platformDataObj.Id, platformDataObj.ProcInstanceId, req.UserToken, language)
			if strings.Contains(platformDataObj.Status, "InProgress") && platformDataObj.ProcInstanceId != "" {
				newStatus := getInstanceStatus(platformDataObj.ProcInstanceId, req.UserToken, language)
				if newStatus == "InternallyTerminated" {
					newStatus = "Termination"
				}
				if newStatus != "" && newStatus != platformDataObj.Status {
					// 编排的完成,并不表示 请求完成
					if newStatus == string(models.RequestStatusCompleted) {
						// 防止一个请求重复调用
						if _, ok := processedRequestMap[platformDataObj.Id]; ok {
							continue
						}
						taskSort := GetTaskService().GenerateTaskOrderByRequestId(platformDataObj.Id)
						confirmActions, _ = GetRequestService().CreateRequestConfirm(models.RequestTable{Id: platformDataObj.Id,
							RequestTemplate: platformDataObj.TemplateId, Type: platformDataObj.Type, Role: platformDataObj.Role, CreatedBy: platformDataObj.CreatedBy,
							CreatedTime: platformDataObj.CreatedTime, Name: platformDataObj.Name}, taskSort, req.UserToken, language)
						if len(confirmActions) > 0 {
							actions = append(actions, confirmActions...)
						}
					} else {
						// 防止一个请求重复调用
						if _, ok := processedRequestMap[platformDataObj.Id]; ok {
							platformDataObj.Status = newStatus
							continue
						}
						// 只处理自动退出&手动终止终止情况,需要发邮件
						if newStatus == string(models.RequestStatusFaulted) || newStatus == string(models.RequestStatusTermination) {
							NotifyTaskWorkflowFailMail(platformDataObj.Name, platformDataObj.ProcDefName, newStatus, platformDataObj.CreatedBy, req.UserToken, language)
						}
						actions = append(actions, &dao.ExecAction{Sql: "update request set status=?,updated_time=? where id=?",
							Param: []interface{}{newStatus, time.Now().Format(models.DateTimeFormat), platformDataObj.Id}})
						platformDataObj.Status = newStatus
					}
					processedRequestMap[platformDataObj.Id] = true
				}
			}
			if collectMap[platformDataObj.ParentId] {
				platformDataObj.CollectFlag = 1
			}
			if platformDataObj.OperatorObjType == "" {
				platformDataObj.OperatorObjType = operatorObjTypeMap[platformDataObj.ProcDefKey]
			}
			// 人员设置方式: 模板指定,提交人指定, 只能: 处理人(处理)和角色管理员(转给我),需要RoleAdministrator、HandlerType返回给前端判断
			// 设置角色管理员
			if v, ok := roleDtoMap[platformDataObj.Role]; ok {
				platformDataObj.RoleAdministrator = v.Administrator
			}
			// 设置人员设置方式
			if platformDataObj.TaskHandleId != "" {
				taskHandle, _ := GetTaskHandleService().Get(platformDataObj.TaskHandleId)
				if taskHandle != nil {
					platformDataObj.HandlerType = taskHandle.HandlerType
				}
			}
			if platformDataObj.OperatorObj == "" && platformDataObj.Cache != "" {
				result, _ := GetEntityData(platformDataObj.Id, req.UserToken, language)
				if len(result.Data) > 0 {
					var cacheObj models.RequestPreDataDto
					err = json.Unmarshal([]byte(platformDataObj.Cache), &cacheObj)
					for _, entity := range result.Data {
						if entity.Id == cacheObj.RootEntityId {
							platformDataObj.OperatorObj = entity.DisplayName
							actions = append(actions, &dao.ExecAction{Sql: "update request set operator_obj=? where id=?", Param: []interface{}{platformDataObj.OperatorObj, platformDataObj.Id}})
						}
					}
				}
			}
			platformDataObj.HandleRole, platformDataObj.Handler = getRequestHandler(models.RequestTable{Id: platformDataObj.Id, Status: platformDataObj.Status, RequestTemplate: platformDataObj.TemplateId, CreatedBy: platformDataObj.CreatedBy, Role: platformDataObj.Role})
			if platformDataObj.Status == string(models.RequestStatusDraft) {
				// 草稿状态,定版处理人和角色需要读取模版和定版配置
				platformDataObj.CheckHandleRole, platformDataObj.CheckHandler = GetTaskTemplateService().GetCheckRoleAndHandler(platformDataObj.TemplateId)
				if platformDataObj.CheckHandleRole != "" {
					platformDataObj.CheckHandleRoleDisplay = roleDisplayMap[platformDataObj.CheckHandleRole]
				}
			}
			//如果是待处理tab, 会出现同一个人,2个处理角色,采用2条记录返回,同时每个处理角色和人与每条记录适配
			if req.Tab == "pending" || req.Tab == "myPending" {
				platformDataObj.HandleRole = platformDataObj.TaskHandleRole
				platformDataObj.HandleRoleDisplay = roleDisplayMap[platformDataObj.HandleRole]
				platformDataObj.Handler = platformDataObj.TaskHandler
			}
			// 计算请求/任务停留时长
			calcRequestStayTime(platformDataObj)
			// 计算是否出撤回按钮
			platformDataObj.RevokeBtn = calcShowRequestRevokeButton(platformDataObj.Id, platformDataObj.Status)
			// 设置角色显示名
			if v, ok := roleDisplayMap[platformDataObj.Role]; ok {
				platformDataObj.RoleDisplay = v
			}
			if strings.Contains(platformDataObj.HandleRole, ",") {
				var newRoleDisplayArr []string
				roleArr := strings.Split(platformDataObj.HandleRole, ",")
				for _, role := range roleArr {
					if v, ok := roleDisplayMap[role]; ok {
						newRoleDisplayArr = append(newRoleDisplayArr, v)
					}
				}
				platformDataObj.HandleRoleDisplay = strings.Join(newRoleDisplayArr, ",")
			} else {
				platformDataObj.HandleRoleDisplay = roleDisplayMap[platformDataObj.HandleRole]
			}
			if strings.Contains(platformDataObj.TaskHandleRole, ",") {
				var newRoleDisplayArr []string
				roleArr := strings.Split(platformDataObj.TaskHandleRole, ",")
				for _, role := range roleArr {
					if v, ok := roleDisplayMap[role]; ok {
						newRoleDisplayArr = append(newRoleDisplayArr, v)
					}
				}
				platformDataObj.TaskHandleRoleDisplay = strings.Join(newRoleDisplayArr, ",")
			} else {
				platformDataObj.TaskHandleRoleDisplay = roleDisplayMap[platformDataObj.TaskHandleRole]
			}
		}
		if len(actions) > 0 {
			updateRequestErr := dao.Transaction(actions)
			if updateRequestErr != nil {
				log.Logger.Error("Try to update request status fail", log.Error(updateRequestErr))
			}
		}
	}
	return
}

func getInstanceStatus(instanceId, userToken, language string) string {
	processInstance, err := GetProcDefService().GetProcessDefineInstance(instanceId, userToken, language)
	if err != nil {
		return ""
	}
	if processInstance == nil {
		return ""
	}
	if processInstance.Status != string(models.RequestStatusInProgress) {
		return processInstance.Status
	}
	status := string(models.RequestStatusInProgress)
	if len(processInstance.TaskNodeInstances) > 0 {
		for _, v := range processInstance.TaskNodeInstances {
			if v.Status == string(models.RequestStatusFaulted) {
				status = string(models.RequestStatusInProgressFaulted)
				break
			}
			if v.Status == models.ProcDefStatusTimeout {
				status = string(models.RequestStatusInProgressTimeOuted)
				break
			}
		}
	}
	return status
}

// CalcRequestProgressAndCurNode 获取当前节点和计算请求进度
func CalcRequestProgressAndCurNode(requestId, instanceId, userToken, language string) (progress int, curNode string) {
	var taskTemplateList []*models.TaskTemplateTable
	var doneTaskList []*models.TaskTable
	var task *models.TaskTable
	var request models.RequestTable
	taskTemplateList, _ = GetTaskTemplateService().GetTaskTemplateListByRequestId(requestId)
	request, _ = GetSimpleRequest(requestId)
	doneTaskList, _ = GetTaskService().GetDoneTaskByRequestId(request)
	if len(taskTemplateList) > 0 {
		// 分母+1,统计已完成节点
		progress = int(math.Floor(float64(len(doneTaskList))/float64(len(taskTemplateList)+1)*100 + 0.5))
	}
	if instanceId == "" {
		// 无编排实例
		switch request.Status {
		case string(models.RequestStatusDraft):
			curNode = WaitCommit
			progress = 0
		case string(models.RequestStatusPending):
			curNode = RequestPending
		case string(models.RequestStatusConfirm):
			curNode = Confirm
		case string(models.RequestStatusCompleted):
			curNode = RequestComplete
			progress = 100
		}
		if curNode == "" {
			task, _ = GetTaskService().GetDoingTask(requestId, request.RequestTemplate)
			if task != nil {
				curNode = task.Name
			}
		}
		return
	}
	processInstance, err := GetProcDefService().GetProcessDefineInstance(instanceId, userToken, language)
	if err != nil || processInstance == nil || len(processInstance.TaskNodeInstances) == 0 {
		return
	}
	switch processInstance.Status {
	case "Completed":
		if request.Status == string(models.RequestStatusCompleted) {
			curNode = RequestComplete
			progress = 100
		} else if request.Status == string(models.RequestStatusConfirm) {
			curNode = Confirm
		}
		return
	case "InProgress":
		for _, v := range processInstance.TaskNodeInstances {
			if v.Status == "InProgress" || v.Status == "Timeouted" || v.Status == "Faulted" {
				curNode = v.NodeName
				return
			}
		}
	case "NotStarted":
		curNode = "NotStarted"
	case "Faulted":
		// 失败状态,筛选成功并且有序号的节点
		list := filterSuccessNode(processInstance.TaskNodeInstances)
		if len(list) == 0 {
			return
		}
		// 按 orderNo排序,将有 orderNo的节点按小到大排序,找到最后一个节点状态返回
		sort.Sort(models.QueryNodeSort(list))
		curNode = list[len(list)-1].NodeName
		return
	default:
		// 失败状态,显示具体执行失败的节点. filterNode 过滤orderNo为空大节点
		list := filterNode(processInstance.TaskNodeInstances)
		if len(list) == 0 {
			// 如果都没有序号,找一个NotStarted节点,找不到返回 Completed
			for _, v := range processInstance.TaskNodeInstances {
				if v.Status == "NotStarted" {
					curNode = v.NodeName
					return
				}
			}
			log.Logger.Error("filterNode list is empty fail,instanceId", log.String("instanceId", instanceId))
			return
		}
		// 按 orderNo排序,将有 orderNo的节点按小到大排序,查找第一个非完成的节点状态返回
		sort.Sort(models.QueryNodeSort(list))
		for _, item := range list {
			if item.Status != "Completed" {
				curNode = item.NodeName
				return
			}
		}
	}
	return
}

func filterNode(instances []*models.TaskNodeInstances) []*models.TaskNodeInstances {
	var list []*models.TaskNodeInstances
	for _, node := range instances {
		if node.OrderedNo != "" {
			list = append(list, node)
		}
	}
	return list
}

func filterSuccessNode(instances []*models.TaskNodeInstances) []*models.TaskNodeInstances {
	var list []*models.TaskNodeInstances
	for _, node := range instances {
		if node.OrderedNo != "" && node.Status == "Completed" {
			list = append(list, node)
		}
	}
	return list
}

func transPlatConditionToSQL(param *models.PlatformRequestParam) string {
	return "where 1 = 1 " + transCommonRequestToSQL(param.CommonRequestParam)
}

func transHistoryConditionToSQL(param *models.RequestHistoryParam) (where string) {
	where = "where del_flag = 0 "
	if param.Tab == "commit" {
		where = where + " and status <> 'Draft'"
	} else if param.Tab == "draft" {
		where = where + " and status = 'Draft' and rollback_desc is  null and revoke_flag = 0"
	} else if param.Tab == "rollback" {
		where = where + " and status = 'Draft' and rollback_desc is  not null "
	} else if param.Tab == "revoke" {
		where = where + " and status = 'Draft' and revoke_flag = 1"
	}
	if param.Action != 0 {
		where = where + fmt.Sprintf(" and type = %d", param.Action)
	}
	where = where + transCommonRequestToSQL(param.CommonRequestParam)
	return
}

func transCommonRequestToSQL(param models.CommonRequestParam) (where string) {
	if param.Id != "" {
		where = where + " and ( id like '%" + param.Id + "%') "
	}
	if param.Name != "" {
		where = where + " and ( name like '%" + param.Name + "%') "
	}
	if len(param.TemplateId) > 0 {
		where = where + " and template_id in (" + getSQL(param.TemplateId) + ")"
	}
	if len(param.Status) > 0 {
		where = where + " and status in (" + getSQL(param.Status) + ")"
	}
	if param.OperatorObj != "" {
		where = where + " and operator_obj = '" + param.OperatorObj + "'"
	}
	if len(param.CreatedBy) > 0 {
		where = where + " and created_by in (" + getSQL(param.CreatedBy) + ")"
	}
	if len(param.OperatorObjType) > 0 {
		where = where + " and operator_obj_type in (" + getSQL(param.OperatorObjType) + ")"
	}
	if len(param.ProcDefName) > 0 {
		where = where + " and proc_def_name in (" + getSQL(param.ProcDefName) + ")"
	}
	if len(param.Handler) > 0 {
		where = where + " and handler in (" + getSQL(param.Handler) + ")"
	}
	if param.TaskName != "" {
		where = where + " and ( task_name like '%" + param.TaskName + "%') "
	}
	if param.CreatedStartTime != "" && param.CreatedEndTime != "" {
		where = where + " and created_time >= '" + param.CreatedStartTime + "' and created_time <= '" + param.CreatedEndTime + "'"
	}
	if param.ReportStartTime != "" && param.ReportEndTime != "" {
		where = where + " and report_time >= '" + param.ReportStartTime + "' and report_time <= '" + param.ReportEndTime + "'"
	}
	if param.UpdatedStartTime != "" && param.UpdatedEndTime != "" {
		where = where + " and updated_time >= '" + param.UpdatedStartTime + "' and updated_time <= '" + param.UpdatedEndTime + "'"
	}
	if param.ExpectStartTime != "" && param.ExpectEndTime != "" {
		where = where + " and expect_time >= '" + param.ExpectStartTime + "' and expect_time <= '" + param.ExpectEndTime + "'"
	}
	if param.ApprovalStartTime != "" && param.ApprovalEndTime != "" {
		where = where + " and approval_time >= '" + param.ApprovalStartTime + "' and approval_time <= '" + param.ApprovalEndTime + "'"
	}
	if param.TaskCreatedStartTime != "" && param.TaskCreatedEndTime != "" {
		where = where + " and task_created_time >= '" + param.TaskCreatedStartTime + "' and task_created_time <= '" + param.TaskCreatedEndTime + "'"
	}
	if param.TaskApprovalStartTime != "" && param.TaskApprovalEndTime != "" {
		where = where + " and task_approval_time >= '" + param.TaskApprovalStartTime + "' and task_approval_time <= '" + param.TaskApprovalEndTime + "'"
	}
	if param.TaskExpectStartTime != "" && param.TaskExpectEndTime != "" {
		where = where + " and task_expect_time >= '" + param.TaskExpectStartTime + "' and task_expect_time <= '" + param.TaskExpectEndTime + "'"
	}
	if param.TaskHandleUpdatedStartTime != "" && param.TaskHandleUpdatedEndTime != "" {
		where = where + " and task_handle_updated_time >= '" + param.TaskHandleUpdatedStartTime + "' and task_handle_updated_time <= '" + param.TaskHandleUpdatedEndTime + "'"
	}
	return
}

func transCollectConditionToSQL(param *models.QueryCollectTemplateParam) (where string) {
	where = "where 1 = 1 "
	if param.Id != "" {
		where = where + " and ( id like '%" + param.Id + "%' or id like '%" + param.Id + "%')"
	}
	if param.Name != "" {
		where = where + " and ( id like '%" + param.Name + "%' or name like '%" + param.Name + "%')"
	}
	if len(param.TemplateGroupId) > 0 {
		where = where + " and template_group_id in (" + getSQL(param.TemplateGroupId) + ")"
	}
	if len(param.OperatorObjType) > 0 {
		where = where + " and operator_obj_type in (" + getSQL(param.OperatorObjType) + ")"
	}
	if len(param.ProcDefName) > 0 {
		where = where + " and proc_def_name in (" + getSQL(param.ProcDefName) + ")"
	}
	if len(param.ManageRole) > 0 {
		where = where + " and manage_role in (" + getSQL(param.ManageRole) + ")"
	}
	if len(param.Owner) > 0 {
		where = where + " and owner in (" + getSQL(param.Owner) + ")"
	}
	if len(param.Tags) > 0 {
		where = where + " and tags in (" + getSQL(param.Tags) + ")"
	}
	if param.CreatedStartTime != "" && param.CreatedEndTime != "" {
		where = where + " and created_time >= '" + param.CreatedStartTime + "' and created_time <= '" + param.CreatedEndTime + "'"
	}
	if param.UpdatedStartTime != "" && param.UpdatedEndTime != "" {
		where = where + " and updated_time >= '" + param.UpdatedStartTime + "' and updated_time <= '" + param.UpdatedEndTime + "'"
	}
	if len(param.UseRole) > 0 {
		var templateIdList []string
		dao.X.SQL("select request_template from request_template_role where role_type='USE' and role in (" + getSQL(param.UseRole) + ")").Find(&templateIdList)
		if len(templateIdList) > 0 {
			where = where + " and id in (" + getSQL(templateIdList) + ")"
		}
	}
	return
}

func getSQL(status []string) string {
	var sql string
	for i := 0; i < len(status); i++ {
		if i == len(status)-1 {
			sql = sql + "'" + status[i] + "'"
		} else {
			sql = sql + "'" + status[i] + "',"
		}
	}
	return sql
}

// GetRequestProgress  请求已创建时,获取请求进度
func GetRequestProgress(requestId, userToken, language string) (rowData *models.RequestProgress, err error) {
	var taskTemplateList []*models.TaskTemplateTable
	var taskTemplateProgressList []*models.TaskTemplateProgressDto
	var taskApproveTemplateList []*models.TaskTemplateTable
	var taskImplementTemplateList []*models.TaskTemplateTable
	var taskMap map[string]*models.TaskTable
	var taskHandleTemplateList []*models.TaskHandleTemplateTable
	var taskHandleList []*models.TaskHandleTable
	var taskApproveList, taskImplementList []*models.TaskTable
	// 处理人和角色
	var handler, role string
	var taskSort, status int
	var roleDisplayMap = make(map[string]string)
	var request models.RequestTable
	var approvalProgress, taskProgress []*models.TaskProgressNode
	var requestTemplateRoleList []*models.RequestTemplateRoleTable
	var requestTemplate *models.RequestTemplateTable
	var requestTaskHandleMap = make(map[string]*models.TaskHandleTemplateDto)
	var taskTemplateDtoList []*models.TaskTemplateDto
	// 初始化成未开始
	approveCompleteStatus := int(models.TaskExecStatusNotStart)
	taskCompleteStatus := int(models.TaskExecStatusNotStart)
	completeStatus := int(models.TaskExecStatusNotStart)
	rowData = &models.RequestProgress{RequestProgress: []*models.RequestProgressNode{}, ApprovalProgress: []*models.TaskProgressNode{}, TaskProgress: []*models.TaskProgressNode{}}
	taskTemplateList, err = GetTaskTemplateService().GetTaskTemplateListByRequestId(requestId)
	if err != nil {
		return
	}
	request, err = GetSimpleRequest(requestId)
	if err != nil {
		return
	}
	taskMap, err = GetTaskService().GetTaskMapByRequestId(request)
	if err != nil {
		return
	}
	roleDisplayMap, err = GetRoleService().GetRoleDisplayName(userToken, language)
	if err != nil {
		return
	}
	// 获取用户请求提交的审批配置(提交人指定)
	if request.TaskApprovalCache != "" {
		json.Unmarshal([]byte(request.TaskApprovalCache), &taskTemplateDtoList)
		if len(taskTemplateDtoList) > 0 {
			for _, dto := range taskTemplateDtoList {
				if len(dto.HandleTemplates) > 0 {
					for _, template := range dto.HandleTemplates {
						if template.Id != "" {
							requestTaskHandleMap[template.Id] = template
						}
					}
				}
			}
		}
	}
	// 读取审批任务状态&任务状态
	for _, task := range taskMap {
		if task.Type == string(models.TaskTypeApprove) {
			taskApproveList = append(taskApproveList, task)
		} else if task.Type == string(models.TaskTypeImplement) {
			taskImplementList = append(taskImplementList, task)
		}
	}
	if len(taskApproveList) > 0 {
		// 有创建审批,则状态设置为完成
		approveCompleteStatus = int(models.TaskExecStatusCompleted)
		for _, taskApprove := range taskApproveList {
			// 有任务未完成,则设置成 进行中
			if taskApprove.Status != string(models.TaskStatusDone) || taskApprove.TaskResult == string(models.TaskHandleResultTypeDeny) {
				approveCompleteStatus = int(models.TaskExecStatusDoing)
				break
			}
		}
	}
	if len(taskImplementList) > 0 {
		// 有创建任务,则状态设置为完成
		taskCompleteStatus = int(models.TaskExecStatusCompleted)
		for _, taskImplement := range taskImplementList {
			// 有任务未完成,则设置成 进行中
			if taskImplement.Status != string(models.TaskStatusDone) {
				taskCompleteStatus = int(models.TaskExecStatusDoing)
				break
			}
		}
		// 根据请求状态重新设置
		if request.Status == string(models.RequestStatusInProgressFaulted) {
			taskCompleteStatus = int(models.TaskExecStatusFail)
		}
	}
	if len(taskTemplateList) > 0 {
		for _, taskTemplate := range taskTemplateList {
			taskSort = 0
			taskHandleTemplateList = []*models.TaskHandleTemplateTable{}
			handler = ""
			role = ""
			status = int(models.TaskExecStatusNotStart)
			switch taskTemplate.Type {
			case string(models.TaskTypeSubmit):
				taskSort = 1
				role = request.Role
				handler = request.CreatedBy
				status = int(models.TaskExecStatusDoing)
			case string(models.TaskTypeCheck):
				taskSort = 2
			case string(models.TaskTypeConfirm):
				taskSort = 5
				role = request.Role
				handler = request.CreatedBy
			case string(models.TaskTypeApprove):
				taskApproveTemplateList = append(taskApproveTemplateList, taskTemplate)
				continue
			case string(models.TaskTypeImplement):
				taskImplementTemplateList = append(taskImplementTemplateList, taskTemplate)
				continue
			default:
			}
			err = dao.X.SQL("select * from task_handle_template where task_template = ?", taskTemplate.Id).Find(&taskHandleTemplateList)
			if err != nil {
				return
			}
			if len(taskHandleTemplateList) > 0 {
				handler = taskHandleTemplateList[0].Handler
				role = taskHandleTemplateList[0].Role
				// 定版模板处理角色和处理人如果为空,则设置为属主
				if taskTemplate.Type == string(models.TaskTypeCheck) && role == "" {
					requestTemplate, _ = GetRequestTemplateService().GetRequestTemplate(request.RequestTemplate)
					if requestTemplate != nil {
						handler = requestTemplate.Handler
					}
					requestTemplateRoleList, _ = GetRequestTemplateService().GetRequestTemplateRole(request.RequestTemplate)
					if len(requestTemplateRoleList) > 0 {
						for _, requestTemplateRole := range requestTemplateRoleList {
							if requestTemplateRole.RoleType == string(models.RolePermissionMGMT) {
								role = requestTemplateRole.Role
								break
							}
						}
					}
				}
			}
			// 设置为显示名
			if v, ok := roleDisplayMap[role]; ok {
				role = v
			}
			taskTemplateProgressList = append(taskTemplateProgressList, &models.TaskTemplateProgressDto{
				Id:          taskTemplate.Id,
				Type:        taskTemplate.Type,
				Node:        taskTemplate.Name,
				Handler:     handler,
				Role:        role,
				Status:      status, //初始化状态
				ApproveType: taskTemplate.HandleMode,
				Sort:        taskSort,
			})
		}
	}
	if len(taskApproveTemplateList) > 0 {
		taskTemplateProgressList = append(taskTemplateProgressList, &models.TaskTemplateProgressDto{
			Type:   string(models.TaskTypeApprove),
			Node:   Approval,
			Status: approveCompleteStatus,
			Sort:   3,
		})
	}
	if len(taskImplementTemplateList) > 0 {
		taskTemplateProgressList = append(taskTemplateProgressList, &models.TaskTemplateProgressDto{
			Type:   string(models.TaskTypeImplement),
			Node:   Task,
			Status: taskCompleteStatus,
			Sort:   4,
		})
	}

	if request.Status == string(models.RequestStatusFaulted) {
		// 自动退出
		taskTemplateProgressList = append(taskTemplateProgressList, &models.TaskTemplateProgressDto{
			Node:   AutoExit,
			Status: int(models.TaskExecStatusAutoExitStatus),
			Sort:   6,
		})
	} else {
		// 添加请求完成
		if request.Status == string(models.RequestStatusCompleted) {
			completeStatus = int(models.TaskExecStatusCompleted)
		}
		// 非自动退出,都会有请求完成状态
		taskTemplateProgressList = append(taskTemplateProgressList, &models.TaskTemplateProgressDto{
			Node:   RequestComplete,
			Status: completeStatus,
			Sort:   6,
		})
	}
	sort.Sort(models.TaskTemplateProgressDtoSort(taskTemplateProgressList))

	for _, taskTemplateProgress := range taskTemplateProgressList {
		requestProgress := &models.RequestProgressNode{}
		requestProgress.Status = taskTemplateProgress.Status
		requestProgress.Node = taskTemplateProgress.Node
		requestProgress.Role = taskTemplateProgress.Role
		requestProgress.Handler = taskTemplateProgress.Handler
		if v, ok := taskMap[taskTemplateProgress.Id]; ok {
			taskHandleList = []*models.TaskHandleTable{}
			// 查询到对应任务,表示任务已经创建,拿取最新的处理人,并且更新任务状态
			if v.Status == string(models.TaskStatusDone) {
				requestProgress.Status = int(models.TaskExecStatusCompleted)
			} else {
				requestProgress.Status = int(models.TaskExecStatusDoing)
			}
			// 请求进度里面包含: 提交、定版、确认节点,这些都只会有一个处理节点
			taskHandleList, err = GetTaskHandleService().GetTaskHandleListByTaskId(v.Id)
			if err != nil {
				return
			}
			if len(taskHandleList) > 0 {
				if taskHandleList[0].Handler != "" {
					requestProgress.Handler = taskHandleList[0].Handler
				}
				if taskHandleList[0].Role != "" {
					// 设置为显示名
					if v, ok := roleDisplayMap[taskHandleList[0].Role]; ok {
						requestProgress.Role = v
					}
				}
			}
		}
		rowData.RequestProgress = append(rowData.RequestProgress, requestProgress)
	}

	// 添加审批进度
	approvalProgress, err = getTaskProgress(request.Role, userToken, language, taskApproveTemplateList, taskMap, requestTaskHandleMap, roleDisplayMap)
	if err != nil {
		return
	}
	if len(approvalProgress) > 0 {
		rowData.ApprovalProgress = approvalProgress
	}
	// 添加任务进度
	taskProgress, err = getTaskProgress(request.Role, userToken, language, taskImplementTemplateList, taskMap, requestTaskHandleMap, roleDisplayMap)
	if err != nil {
		return
	}
	if len(taskProgress) > 0 {
		if request.ProcInstanceId != "" && request.Status != string(models.RequestStatusCompleted) && request.Status != string(models.RequestStatusFaulted) {
			response, err := rpc.GetProcessInstance(userToken, language, request.ProcInstanceId)
			if err != nil {
				log.Logger.Error("http getProcessInstances error", log.Error(err))
			}
			if response != nil {
				if response.Status == InternallyTerminated {
					taskProgress = append(taskProgress, &models.TaskProgressNode{Node: InternallyTerminated, Status: int(models.TaskExecStatusInternallyTerminated)})
				}
				// 记录错误节点,如果实例运行中有错误节点,则需要把运行节点展示在列表中并展示对应位置
				var exist bool
				var sort int
				for _, v := range response.TaskNodeInstances {
					exist = false
					// 任务节点根据 编排任务orderNo排序
					for _, rowData := range taskProgress {
						if rowData.NodeDefId == v.NodeDefId || rowData.NodeId == v.NodeId {
							s, _ := strconv.Atoi(v.OrderedNo)
							rowData.Sort = s
							break
						}
					}
					if v.Status == string(models.RequestStatusFaulted) || v.Status == "Timeouted" {
						for _, rowData := range taskProgress {
							if rowData.NodeDefId == v.NodeDefId || rowData.NodeId == v.NodeId {
								exist = true
								rowData.Status = int(models.TaskExecStatusFail)
								break
							}
						}
						if !exist {
							sort, _ = strconv.Atoi(v.OrderedNo)
							taskProgress = append(taskProgress, &models.TaskProgressNode{
								NodeId:         v.NodeId,
								Node:           v.NodeName,
								NodeDefId:      v.NodeDefId,
								Status:         int(models.TaskExecStatusFail),
								Sort:           sort,
								TaskHandleList: []*models.TaskHandleNode{{Handler: AutoNode}},
								NodeType:       "auto",
							})
						}
					}
				}
			}
			sort.Sort(models.TaskProgressNodeSort(taskProgress))
		}
		rowData.TaskProgress = taskProgress
	}
	return
}

func getTaskProgress(role, userToken, language string, taskTemplateList []*models.TaskTemplateTable, taskMap map[string]*models.TaskTable, requestTaskHandleMap map[string]*models.TaskHandleTemplateDto, roleDisplayMap map[string]string) ([]*models.TaskProgressNode, error) {
	var taskHandleTemplateList []*models.TaskHandleTemplateTable
	var taskHandleList []*models.TaskHandleTable
	var taskProgressList []*models.TaskProgressNode
	var err error
	// 任务排序
	if len(taskTemplateList) > 0 {
		sort.Sort(models.TaskTemplateTableSort(taskTemplateList))
		for _, taskTemplate := range taskTemplateList {
			taskHandleTemplateList = []*models.TaskHandleTemplateTable{}
			requestProgress := &models.TaskProgressNode{TaskHandleList: []*models.TaskHandleNode{}}
			requestProgress.Status = int(models.TaskExecStatusNotStart)
			requestProgress.Node = taskTemplate.Name
			requestProgress.ApproveType = taskTemplate.HandleMode
			requestProgress.NodeId = taskTemplate.NodeId
			requestProgress.NodeDefId = taskTemplate.NodeDefId
			err = dao.X.SQL("select * from task_handle_template where task_template = ?", taskTemplate.Id).Find(&taskHandleTemplateList)
			if err != nil {
				return nil, err
			}
			if len(taskHandleTemplateList) > 0 {
				var tempHandler, tempRole string
				for _, taskHandleTemplate := range taskHandleTemplateList {
					tempHandler = ""
					tempRole = ""
					// 用户提交的处理人优先级最高
					if v, ok := requestTaskHandleMap[taskHandleTemplate.Id]; ok {
						tempHandler = v.Handler
						tempRole = v.Role
					}
					if tempHandler == "" && taskHandleTemplate.Handler != "" {
						tempHandler = taskHandleTemplate.Handler
					}
					if tempRole == "" && taskHandleTemplate.Role != "" {
						tempRole = taskHandleTemplate.Role
					}
					// 设置显示名
					if v, ok := roleDisplayMap[tempRole]; ok {
						tempRole = v
					}
					requestProgress.TaskHandleList = append(requestProgress.TaskHandleList, &models.TaskHandleNode{
						Handler:     tempHandler,
						Role:        tempRole,
						HandlerType: taskHandleTemplate.HandlerType,
					})
				}
			}
			// 如果是角色管理员审批,用户提交匹配不上,需要查询接口拿到
			if taskTemplate.HandleMode == string(models.TaskTemplateHandleModeAdmin) {
				result, _ := GetRoleService().GetRoleAdministrators(role, userToken, language)
				if len(result) > 0 && result[0] != "" {
					// 设置显示名
					if v, ok := roleDisplayMap[role]; ok {
						role = v
					}
					requestProgress.TaskHandleList = append(requestProgress.TaskHandleList, &models.TaskHandleNode{
						Handler: result[0],
						Role:    role,
					})
				}
			}

			if v, ok := taskMap[taskTemplate.Id]; ok {
				taskHandleList = []*models.TaskHandleTable{}
				// 查询到对应任务,表示任务已经创建,拿取最新的处理人,并且更新任务状态
				if v.Status == string(models.TaskStatusDone) {
					requestProgress.Status = int(models.TaskExecStatusCompleted)
				} else {
					requestProgress.Status = int(models.TaskExecStatusDoing)
				}
				taskHandleList, err = GetTaskHandleService().GetTaskHandleListByTaskId(v.Id)
				if err != nil {
					return nil, err
				}
				if len(taskHandleList) > 0 {
					var tempTaskHandleNodeList []*models.TaskHandleNode
					for _, taskHandle := range taskHandleList {
						// 设置显示名
						if v, ok := roleDisplayMap[taskHandle.Role]; ok {
							taskHandle.Role = v
						}
						tempTaskHandleNodeList = append(tempTaskHandleNodeList, &models.TaskHandleNode{
							Handler:     taskHandle.Handler,
							Role:        taskHandle.Role,
							HandlerType: taskHandle.HandlerType,
						})
					}
					requestProgress.TaskHandleList = tempTaskHandleNodeList
				}
			}
			taskProgressList = append(taskProgressList, requestProgress)
		}
	}
	return taskProgressList, nil
}

func GetProcessDefinitions(templateId, userToken, language string) (rowData *models.DefinitionsData, err error) {
	var template *models.RequestTemplateTable
	template, err = GetRequestTemplateService().GetRequestTemplate(templateId)
	if err != nil {
		return
	}
	if template == nil {
		err = fmt.Errorf("requestTemplate not exist")
		return
	}
	return GetProcDefService().GetProcessDefine(template.ProcDefId, userToken, language)
}

// getRequestForm 获取请求信息
func getRequestForm(request *models.RequestTable, userToken, language string) (form models.RequestForm) {
	if request == nil {
		return
	}
	var tmpTemplate []*models.RequestTemplateTmp
	var version string
	var customForm models.CustomForm
	var roleDisplayMap = make(map[string]string)
	err := dao.X.SQL("select rt.name as  template_name,rt.status,rtg.name as template_group_name,rt.version,rt.proc_def_id,rt.expire_day from request_template rt join "+
		"request_template_group rtg on rt.group = rtg.id where rt.id= ?", request.RequestTemplate).Find(&tmpTemplate)
	if err != nil {
		return
	}
	if len(tmpTemplate) == 0 {
		err = fmt.Errorf("can not find request_template with id:%s ", request.Id)
		return
	}
	roleDisplayMap, err = GetRoleService().GetRoleDisplayName(userToken, language)
	if err != nil {
		return
	}
	template := tmpTemplate[0]
	form.Id = request.Id
	form.Name = request.Name
	form.RequestType = request.Type
	form.CreatedTime = request.CreatedTime
	form.ExpectTime = request.ExpectTime
	form.CreatedBy = request.CreatedBy
	if v, ok := roleDisplayMap[request.Role]; ok {
		form.Role = v
	} else {
		form.Role = request.Role
	}

	form.Description = request.Description
	form.Status = request.Status
	form.ProcInstanceId = request.ProcInstanceId
	form.ExpireDay = template.ExpireDay
	if template.Status != "confirm" {
		version = "beta"
	} else {
		version = template.Version
	}
	form.TemplateName = template.TemplateName
	form.Version = version
	form.TemplateGroupName = template.TemplateGroupName
	if request.Status == "Pending" {
		form.CurNode = RequestPending
	}
	form.Progress, form.CurNode = CalcRequestProgressAndCurNode(request.Id, request.ProcInstanceId, userToken, language)
	if template.ProcDefId != "" {
		form.AssociationWorkflow = true
	}
	_, form.Handler = getRequestHandler(*request)
	if request.CustomFormCache != "" {
		err = json.Unmarshal([]byte(request.CustomFormCache), &customForm)
		if err != nil {
			log.Logger.Error("json Unmarshal", log.Error(err), log.String("CustomFormCache", request.CustomFormCache))
			return
		}
	} else {
		var items []*models.FormItemTemplateTable
		dao.X.SQL("select * from form_item_template where form_template in (select id from form_template  where request_template=? and"+
			" request_form_type = ?) order by item_group,sort", request.RequestTemplate, models.RequestFormTypeMessage).Find(&items)
		customForm.Title = items
	}
	form.CustomForm = customForm

	form.AttachFiles = GetRequestAttachFileList(request.Id)
	if request.Cache != "" {
		var cacheObj models.RequestPreDataDto
		err = json.Unmarshal([]byte(request.Cache), &cacheObj)
		if err != nil {
			err = fmt.Errorf("Try to json unmarshal cache data fail,%s ", err.Error())
			return
		}
		form.FormData = cacheObj.Data
		form.RootEntityId = cacheObj.RootEntityId
		form.OperatorObj = cacheObj.EntityName
	}
	form.RevokeBtn = calcShowRequestRevokeButton(request.Id, request.Status)
	return
}

// getRequestHandler 获取请求处理人,如果处于任务执行状态
func getRequestHandler(request models.RequestTable) (role, handler string) {
	var task *models.TaskTable
	var taskHandleList []*models.TaskHandleTable
	var roleArr, handlerArr []string
	if request.Status == string(models.RequestStatusDraft) {
		// 草稿状态,当前处理人为自己
		return request.Role, request.CreatedBy
	}
	task, _ = GetTaskService().GetDoingTask(request.Id, request.RequestTemplate)
	if task != nil {
		// 根据任务查询 任务处理人
		dao.X.SQL("select * from task_handle where task = ? and latest_flag = 1", task.Id).Find(&taskHandleList)
		if len(taskHandleList) > 0 {
			for _, taskHandle := range taskHandleList {
				// 待处理 任务节点和角色都要统计
				if taskHandle.HandleStatus == string(models.TaskHandleResultTypeUncompleted) {
					roleArr = append(roleArr, taskHandle.Role)
					handlerArr = append(handlerArr, taskHandle.Handler)
				}
			}
		}
	}
	if len(roleArr) > 0 {
		role = strings.Join(roleArr, ",")
	}
	if len(handlerArr) > 0 {
		handler = strings.Join(handlerArr, ",")
	}
	return
}

func GetFilterItem(param models.FilterRequestParam) (data *models.FilterItem, err error) {
	data = &models.FilterItem{
		RequestTemplateList: make([]*models.KeyValuePair, 0),
		ReleaseTemplateList: make([]*models.KeyValuePair, 0),
	}
	var templateMap = make(map[string]bool, 0)
	var pairList []*models.KeyValuePair
	var dataList []*models.FilterObj
	var operatorObjTypeMap = make(map[string]bool)
	var procDefNameMap = make(map[string]bool)
	var createdByMap = make(map[string]bool)
	var handlerMap = make(map[string]bool)
	var sql = "select rt.id as template_id,rt.name as template_name,rt.version,rt.type as template_type,rt.operator_obj_type,rt.proc_def_name,r.created_by," +
		"r.handler from request r join request_template rt on r.request_template = rt.id where r.created_time > ?"
	err = dao.X.SQL(sql, param.StartTime).Find(&dataList)
	var handlerList []string
	dao.X.SQL("select handler from task_template").Find(&handlerList)
	if err != nil {
		return
	}
	if len(handlerList) > 0 {
		for _, handler := range handlerList {
			handlerMap[handler] = true
		}
	}
	for _, item := range dataList {
		operatorObjTypeMap[item.OperatorObjType] = true
		procDefNameMap[item.ProcDefName] = true
		createdByMap[item.CreatedBy] = true
		handlerMap[item.Handler] = true
		m := &models.KeyValuePair{TemplateId: item.TemplateId, TemplateName: item.TemplateName, Version: item.Version}
		if m.Version == "" {
			m.Version = "beta"
		}
		if templateMap[m.TemplateId+m.TemplateName+m.Version] {
			continue
		} else {
			templateMap[m.TemplateId+m.TemplateName+m.Version] = true
		}
		pairList = append(pairList, m)
		if item.TemplateType == 0 {
			data.RequestTemplateList = append(data.RequestTemplateList, m)
		} else {
			data.ReleaseTemplateList = append(data.ReleaseTemplateList, m)
		}
	}
	data.TemplateList = pairList
	if len(data.TemplateList) > 0 {
		sort.Sort(models.KeyValueSort(data.TemplateList))
	}
	if len(data.RequestTemplateList) > 0 {
		sort.Sort(models.KeyValueSort(data.RequestTemplateList))
	}
	if len(data.ReleaseTemplateList) > 0 {
		sort.Sort(models.KeyValueSort(data.ReleaseTemplateList))
	}
	data.OperatorObjTypeList = convertMap2Array(operatorObjTypeMap)
	if len(data.OperatorObjTypeList) > 0 {
		sort.Strings(data.OperatorObjTypeList)
	}
	data.ProcDefNameList = convertMap2Array(procDefNameMap)
	if len(data.ProcDefNameList) > 0 {
		sort.Strings(data.ProcDefNameList)
	}
	data.CreatedByList = convertMap2Array(createdByMap)
	if len(data.CreatedByList) > 0 {
		sort.Strings(data.CreatedByList)
	}
	data.HandlerList = convertMap2Array(handlerMap)
	if len(data.HandlerList) > 0 {
		sort.Strings(data.HandlerList)
	}
	return
}

// RequestConfirm 请求确认
func RequestConfirm(param models.RequestConfirmParam, user, userToken, language string) error {
	var taskList []*models.TaskTable
	var actions []*dao.ExecAction
	var action *dao.ExecAction
	var request models.RequestTable
	var err error
	var markTaskIdMap = convertArray2Map(param.MarkTaskId)
	now := time.Now().Format(models.DateTimeFormat)
	dao.X.SQL("select * from task where request = ? and type = ?", param.Id, string(models.TaskTypeImplement)).Find(&taskList)
	if len(taskList) > 0 {
		// 更新 处理任务表
		for _, task := range taskList {
			if markTaskIdMap[task.Id] {
				action = &dao.ExecAction{Sql: " update task set confirm_result = ? where id = ?"}
				action.Param = []interface{}{param.CompleteStatus, task.Id}
				actions = append(actions, action)
			}
		}
	}
	// 更新处理节点
	actions = append(actions, &dao.ExecAction{Sql: "update task_handle set handle_result = ?,updated_time = ?,handle_status = ?,result_desc = ? where task = ?", Param: []interface{}{models.TaskHandleResultTypeApprove, now, models.TaskHandleResultTypeComplete, param.Notes, param.TaskId}})
	// 更新请求确认任务设置为已完成
	actions = append(actions, &dao.ExecAction{Sql: "update task set status = ?,task_result = ?,description = ? ,updated_by = ?,updated_time = ? where id = ?", Param: []interface{}{models.TaskStatusDone, models.TaskHandleResultTypeComplete, param.Notes, user, now, param.TaskId}})
	// 更新请求表状态,设置为完成
	actions = append(actions, &dao.ExecAction{Sql: "update request set status = ?,complete_status = ?,updated_by = ?,updated_time= ? where id = ?", Param: []interface{}{models.RequestStatusCompleted, param.CompleteStatus, user, now, param.Id}})
	if request, err = GetSimpleRequest(param.Id); err != nil {
		return err
	}
	NotifyRequestCompleteMail(request.Name, request.CreatedBy, userToken, language)
	return dao.Transaction(actions)
}

func convertMap2Array(hashMap map[string]bool) (arr []string) {
	for key, _ := range hashMap {
		if key != "" {
			arr = append(arr, key)
		}
	}
	return
}

func convertArray2Map(arr []string) map[string]bool {
	hashMap := make(map[string]bool, 0)
	if len(arr) > 0 {
		for _, str := range arr {
			hashMap[str] = true
		}
	}
	return hashMap
}

// CreateRequestCheck 创建确认定版
func (s *RequestService) CreateRequestCheck(request models.RequestTable, operator, bindCache, userToken, language string) (err error) {
	now := time.Now().Format(models.DateTimeFormat)
	var actions []*dao.ExecAction
	var approvalActions []*dao.ExecAction
	var action *dao.ExecAction
	var requestTemplate *models.RequestTemplateTable
	var submitTaskTemplateList, checkTaskTemplateList []*models.TaskTemplateTable
	var taskHandleTemplateList []*models.TaskHandleTemplateTable
	var checkRole, checkHandler string
	var taskSort int
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(request.RequestTemplate)
	if err != nil {
		return err
	}
	if requestTemplate == nil {
		err = fmt.Errorf("requestTemplate is empty")
		return
	}
	if requestTemplate.ProcDefId != "" {
		// 关联编排
		request.AssociationWorkflow = true
	}
	submitTaskTemplateList, err = GetTaskTemplateService().QueryTaskTemplateListByRequestTemplateAndType(requestTemplate.Id, string(models.TaskTypeSubmit))
	if err != nil {
		return
	}
	if len(submitTaskTemplateList) == 0 {
		err = fmt.Errorf("taskTemplate not find submit template")
		return
	}
	expireTime := calcExpireTime(now, requestTemplate.ExpireDay)
	taskSort = GetTaskService().GenerateTaskOrderByRequestId(request.Id)
	// 新增提交请求任务
	submitTaskId := "su_" + guid.CreateGuid()
	action = &dao.ExecAction{Sql: "insert into task(id,name,expire_time,template_type,status,request,task_template,type,task_result,created_by,created_time,updated_by,updated_time,sort,request_created_time) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
	action.Param = []interface{}{submitTaskId, "submit", expireTime, request.Type, models.TaskStatusDone, request.Id, submitTaskTemplateList[0].Id, models.TaskTypeSubmit, models.TaskHandleResultTypeComplete, "system", now, "system", now, taskSort, request.CreatedTime}
	taskSort++
	actions = append(actions, action)

	// 新增提交请求处理人
	action = &dao.ExecAction{Sql: "insert into task_handle(id,task,role,handler,handle_result,handle_status,created_time,updated_time) values (?,?,?,?,?,?,?,?)"}
	action.Param = []interface{}{guid.CreateGuid(), submitTaskId, request.Role, request.CreatedBy, models.TaskHandleResultTypeApprove, models.TaskHandleResultTypeComplete, now, now}
	actions = append(actions, action)

	// 先更新请求状态为定版
	action = &dao.ExecAction{Sql: "update request set status=?,reporter=?,report_time=?,bind_cache=?,updated_by=?,updated_time=?,rollback_desc=null,revoke_flag=0 where id=?"}
	action.Param = []interface{}{models.RequestStatusPending, operator, now, bindCache, operator, now, request.Id}
	actions = append(actions, action)
	// 根据读取配置判断是否跳过定版
	if requestTemplate.CheckSwitch {
		// 新增确认定版任务
		checkTaskId := "ch_" + guid.CreateGuid()
		checkTaskTemplateList, err = GetTaskTemplateService().QueryTaskTemplateListByRequestTemplateAndType(requestTemplate.Id, string(models.TaskTypeCheck))
		if err != nil {
			return
		}
		if len(checkTaskTemplateList) == 0 {
			err = fmt.Errorf("taskTemplate not find check template")
		}
		checkExpireTime := calcExpireTime(now, checkTaskTemplateList[0].ExpireDay)
		action = &dao.ExecAction{Sql: "insert into task(id,name,expire_time,template_type,status,request,task_template,type,created_by,created_time,updated_by,updated_time,sort,request_created_time) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
		action.Param = []interface{}{checkTaskId, "check", checkExpireTime, request.Type, models.TaskStatusCreated, request.Id, checkTaskTemplateList[0].Id, models.TaskTypeCheck, operator, now, operator, now, taskSort, request.CreatedTime}
		actions = append(actions, action)
		taskSort++

		// 新增确认定版处理人
		dao.X.SQL("select * from task_handle_template where task_template = ?", checkTaskTemplateList[0].Id).Find(&taskHandleTemplateList)
		if len(taskHandleTemplateList) > 0 {
			checkRole = taskHandleTemplateList[0].Role
			checkHandler = taskHandleTemplateList[0].Handler
			var requestTemplateRole []*models.RequestTemplateRoleTable
			if checkRole == "" {
				// 定版配置角色为空,取模版属主角色
				dao.X.SQL("select * from request_template_role where request_template = ? and role_type = 'MGMT'", requestTemplate.Id).Find(&requestTemplateRole)
				if len(requestTemplateRole) > 0 {
					checkRole = requestTemplateRole[0].Role
					checkHandler = requestTemplate.Handler
				}
			}
			if checkHandler != "" {
				// 给对应处理人发送邮件
				NotifyTaskAssignMail(request.Name, RequestPending, checkExpireTime, checkHandler, userToken, language)
			} else {
				// 给角色发送邮件
				NotifyTaskRoleMail(request.Name, RequestPending, checkExpireTime, checkRole, userToken, language)
			}
			action = &dao.ExecAction{Sql: "insert into task_handle(id,task_handle_template,task,role,handler,created_time,updated_time) values (?,?,?,?,?,?,?)"}
			action.Param = []interface{}{guid.CreateGuid(), taskHandleTemplateList[0].Id, checkTaskId, checkRole, checkHandler, now, now}
			actions = append(actions, action)
		}
		err = dao.Transaction(actions)
		return
	}

	// 没有配置定版,请求继续往后面走
	approvalActions, err = s.CreateRequestApproval(request, "", userToken, language, taskSort, true)
	if err != nil {
		return
	}
	if len(approvalActions) > 0 {
		actions = append(actions, approvalActions...)
	}
	err = dao.Transaction(actions)
	return
}

// CreateRequestApproval 创建请求审批, submitFlag 当前是否正在提交请求,如果是 taskList直接为空
func (s *RequestService) CreateRequestApproval(request models.RequestTable, curTaskId, userToken, language string, taskSort int, submitFlag bool) (actions []*dao.ExecAction, err error) {
	var requestTaskActions []*dao.ExecAction
	var taskTemplateList []*models.TaskTemplateTable
	var taskList []*models.TaskTable
	var action *dao.ExecAction
	var newTaskId string
	actions = []*dao.ExecAction{}
	now := time.Now().Format(models.DateTimeFormat)
	err = dao.X.SQL("select * from task_template where request_template = ? and type = ? order by sort asc", request.RequestTemplate, string(models.TaskTypeApprove)).Find(&taskTemplateList)
	if err != nil {
		return
	}
	// 没有审批,直接跳过到下一步,到任务
	if len(taskTemplateList) == 0 {
		return s.CreateRequestTask(request, "", userToken, language, taskSort)
	}
	for _, taskTemplate := range taskTemplateList {
		taskList = []*models.TaskTable{}
		if !submitFlag {
			// 查询
			taskList, _ = GetTaskService().GetLatestTaskListByRequestIdAndTaskTemplateId(request.Id, taskTemplate.Id)
		}
		if len(taskList) > 0 {
			// 取最新的任务
			if taskList[0].Status == string(models.TaskStatusDone) || taskList[0].Id == curTaskId {
				// 任务已完成,或者当前正在处理任务
				continue
			}
			// 有未处理完成任务,直接return
			return
		} else {
			// 模版没有对应的的任务,需要创建当前任务并设置审批角色和人,同时根据审批方式设置审批状态
			newTaskId = "ap_" + guid.CreateGuid()
			taskExpireTime := calcExpireTime(now, taskTemplate.ExpireDay)
			// 审批模板配置自动通过,设置当前审批完成,并且直接下一个审批
			if taskTemplate.HandleMode == string(models.TaskTemplateHandleModeAuto) {
				action = &dao.ExecAction{Sql: "insert into task (id,name,expire_time,template_type,description,status,request,task_template,type,task_result,created_by,created_time,updated_by,updated_time,sort,request_created_time) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
				action.Param = []interface{}{newTaskId, taskTemplate.Name, taskExpireTime, request.Type, taskTemplate.Description, models.TaskStatusDone, request.Id, taskTemplate.Id, taskTemplate.Type, models.TaskHandleResultTypeApprove, "system", now, "system", now, taskSort, request.CreatedTime}
				actions = append(actions, action)
				taskSort++
				continue
			}

			// 新增任务
			action = &dao.ExecAction{Sql: "insert into task (id,name,expire_time,template_type,description,status,request,task_template,type,created_by,created_time,updated_by,updated_time,sort,request_created_time) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
			action.Param = []interface{}{newTaskId, taskTemplate.Name, taskExpireTime, request.Type, taskTemplate.Description, models.TaskStatusCreated, request.Id, taskTemplate.Id, taskTemplate.Type, "system", now, "system", now, taskSort, request.CreatedTime}
			actions = append(actions, action)
			taskSort++

			// 根据任务审批模版表&请求人指定,设置审批处理
			createTaskHandleAction := GetTaskHandleService().CreateTaskHandleByTemplate(newTaskId, userToken, language, &request, taskTemplate)
			if len(createTaskHandleAction) > 0 {
				actions = append(actions, createTaskHandleAction...)
			}

			// 更新请求表为审批状态
			action = &dao.ExecAction{Sql: "update request set status=?,updated_time=? where id=?"}
			action.Param = []interface{}{string(models.RequestStatusInApproval), now, request.Id}
			actions = append(actions, action)
			return
		}
	}
	// 所有审批都处理完成,走请求任务处理
	requestTaskActions, err = s.CreateRequestTask(request, "", userToken, language, taskSort)
	if err != nil {
		return
	}
	if len(requestTaskActions) > 0 {
		actions = append(actions, requestTaskActions...)
	}
	return
}

// CreateRequestTask 创建任务
func (s *RequestService) CreateRequestTask(request models.RequestTable, curTaskId, userToken, language string, taskSort int) (actions []*dao.ExecAction, err error) {
	log.Logger.Debug("CreateRequestTask", log.String("taskId", curTaskId))
	var taskTemplateList []*models.TaskTemplateTable
	var requestTemplate *models.RequestTemplateTable
	var taskList []*models.TaskTable
	var newTaskId string
	var action *dao.ExecAction
	var requestConfirmActions, workflowActions []*dao.ExecAction
	now := time.Now().Format(models.DateTimeFormat)
	actions = []*dao.ExecAction{}
	if request.AssociationWorkflow && request.ProcInstanceId == "" && request.BindCache != "" {
		// 关联编排,调用编排启动
		var bindCache models.RequestCacheData
		json.Unmarshal([]byte(request.BindCache), &bindCache)
		workflowActions, err = StartRequestNew(request, userToken, language, bindCache)
		if err != nil {
			return
		}
		if len(workflowActions) > 0 {
			actions = append(actions, workflowActions...)
		}
		return
	}
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(request.RequestTemplate)
	if err != nil {
		return
	}
	if requestTemplate == nil {
		err = fmt.Errorf("requestTemplate is empty")
		return
	}
	err = dao.X.SQL("select * from task_template where request_template = ? and type = ? order by sort asc", request.RequestTemplate, models.TaskTypeImplement).Find(&taskTemplateList)
	if err != nil {
		return
	}
	// 没有任务,直接跳过到下一步,到请求确认
	if len(taskTemplateList) == 0 {
		return s.CreateRequestConfirm(request, taskSort, userToken, language)
	}
	for _, taskTemplate := range taskTemplateList {
		if taskTemplate.NodeDefId != "" {
			// 编排关联的任务,不会在这里触发
			continue
		}
		taskExpireTime := calcExpireTime(now, taskTemplate.ExpireDay)
		taskList, _ = GetTaskService().GetLatestTaskListByRequestIdAndTaskTemplateId(request.Id, taskTemplate.Id)
		if len(taskList) > 0 {
			// 取最新的任务
			if taskList[0].Status == string(models.TaskStatusDone) || taskList[0].Id == curTaskId {
				// 任务已完成,或者当前任务正在处理,直接下一步
				continue
			}
			// 任务状态未完成,直接return
			return
		} else {
			// 模版没有对应的的任务,需要创建当前任务
			newTaskId = "im_" + guid.CreateGuid()
			// 新增任务
			action = &dao.ExecAction{Sql: "insert into task (id,name,expire_time,template_type,description,status,request,task_template,type,created_by,created_time,updated_by,updated_time,sort,request_created_time) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
			action.Param = []interface{}{newTaskId, taskTemplate.Name, taskExpireTime, request.Type, taskTemplate.Description, models.TaskStatusCreated, request.Id, taskTemplate.Id, taskTemplate.Type, "system", now, "system", now, taskSort, request.CreatedTime}
			actions = append(actions, action)
			taskSort++

			// 根据任务审批模版表&请求人指定,设置审批处理
			createTaskHandleAction := GetTaskHandleService().CreateTaskHandleByTemplate(newTaskId, userToken, language, &request, taskTemplate)
			if len(createTaskHandleAction) > 0 {
				actions = append(actions, createTaskHandleAction...)
			}

			// 更新请求表为审批状态
			actions = append(actions, &dao.ExecAction{Sql: "update request set status=?,updated_time=? where id=?", Param: []interface{}{string(models.RequestStatusInProgress), now, request.Id}})
			return
		}
	}
	// 所有审批都处理完成,走请求任务处理
	requestConfirmActions, err = s.CreateRequestConfirm(request, taskSort, userToken, language)
	if err != nil {
		return
	}
	if len(requestConfirmActions) > 0 {
		actions = append(actions, requestConfirmActions...)
	}
	return
}

// CreateRequestConfirm 创建请求确认
func (s *RequestService) CreateRequestConfirm(request models.RequestTable, taskSort int, userToken, language string) (actions []*dao.ExecAction, err error) {
	var newTaskId string
	var taskTemplateList []*models.TaskTemplateTable
	actions = []*dao.ExecAction{}
	now := time.Now().Format(models.DateTimeFormat)
	// 创建请求确认任务
	newTaskId = "co_" + guid.CreateGuid()
	err = dao.X.SQL("select * from task_template where request_template = ? and type = ? order by sort asc", request.RequestTemplate, string(models.TaskTypeConfirm)).Find(&taskTemplateList)
	if err != nil {
		return
	}
	if len(taskTemplateList) == 0 {
		// 请求没有开启请求确认,直接更新请求到完成
		actions = append(actions, &dao.ExecAction{Sql: "update request set status=?,updated_time=? where id=?", Param: []interface{}{string(models.RequestStatusCompleted), now, request.Id}})
		return
	}
	// 新增任务
	taskExpireTime := calcExpireTime(now, taskTemplateList[0].ExpireDay)
	actions = append(actions, &dao.ExecAction{Sql: "insert into task (id,name,expire_time,template_type,status,request,task_template,type,created_by,created_time,updated_by,updated_time,sort,request_created_time) values(?,?,?,?,?,?,?,?,?,?,?,?,?,?)", Param: []interface{}{
		newTaskId, "confirm", taskExpireTime, request.Type, models.TaskStatusCreated, request.Id, taskTemplateList[0].Id, models.TaskTypeConfirm, "system", now, "system", now, taskSort, request.CreatedTime}})
	actions = append(actions, &dao.ExecAction{Sql: "insert into task_handle (id,task,role,handler,created_time,updated_time) values(?,?,?,?,?,?)", Param: []interface{}{guid.CreateGuid(), newTaskId, request.Role, request.CreatedBy, now, now}})
	// 更新请求表状态为请求确认
	actions = append(actions, &dao.ExecAction{Sql: "update request set status=?,updated_time=? where id=?", Param: []interface{}{string(models.RequestStatusConfirm), now, request.Id}})
	// 发送请求确认邮件
	NotifyTaskAssignMail(request.Name, Confirm, taskExpireTime, request.CreatedBy, userToken, language)
	return
}

func (s *RequestService) CreateProcessTask(request models.RequestTable, task *models.TaskTable, userToken, language string) (actions []*dao.ExecAction, err error) {
	log.Logger.Debug("CreateProcessTask", log.String("taskId", task.Id))
	var workflowActions []*dao.ExecAction
	actions = []*dao.ExecAction{}
	if request.AssociationWorkflow && request.ProcInstanceId == "" && request.BindCache != "" {
		// 关联编排,调用编排启动
		var bindCache models.RequestCacheData
		json.Unmarshal([]byte(request.BindCache), &bindCache)
		workflowActions, err = StartRequestNew(request, userToken, language, bindCache)
		if err != nil {
			return
		}
		if len(workflowActions) > 0 {
			actions = append(actions, workflowActions...)
		}
		return
	}
	return
}

// calcShowRequestRevokeButton 计算是否出请求撤回按钮
func calcShowRequestRevokeButton(requestId, requestStatus string) bool {
	var taskList []*models.TaskTable
	if requestStatus == string(models.RequestStatusDraft) || requestStatus == string(models.RequestStatusConfirm) || requestStatus == string(models.RequestStatusCompleted) {
		return false
	}
	// 当前正在请求定版
	if requestStatus == string(models.RequestStatusPending) {
		return true
	}
	taskList, _ = GetTaskService().GetDoneTaskByRequestId(models.RequestTable{Id: requestId, Status: requestStatus})
	if len(taskList) == 1 {
		// 只有一个任务完成,这个任务只能是任务提交,可以展示撤回按钮
		return true
	}
	return false
}
