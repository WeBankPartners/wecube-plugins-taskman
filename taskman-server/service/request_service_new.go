package service

import (
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
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
)

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

// GetRequestCount 工作台请求统计
func GetRequestCount(user string, userRoles []string) (platformData models.PlatformData, err error) {
	var pendingTask, pendingApprove, pendingCheck, pendingConfirm, pending []string
	var hasProcessedTask, hasProcessedApprove, hasProcessedCheck, hasProcessedConfirm, hasProcessed []string
	pendingTask, pendingApprove, pendingCheck, pendingConfirm, pending = GetPendingCount(userRoles)
	hasProcessedTask, hasProcessedApprove, hasProcessedCheck, hasProcessedConfirm, hasProcessed = GetHasProcessedCount(user)
	platformData.Pending = strings.Join(pending, ";")
	platformData.PendingTask = strings.Join(pendingTask, ";")
	platformData.PendingApprove = strings.Join(pendingApprove, ";")
	platformData.PendingCheck = strings.Join(pendingCheck, ";")
	platformData.PendingConfirm = strings.Join(pendingConfirm, ";")
	platformData.HasProcessed = strings.Join(hasProcessed, ";")
	platformData.HasProcessedTask = strings.Join(hasProcessedTask, ";")
	platformData.HasProcessedApprove = strings.Join(hasProcessedApprove, ";")
	platformData.HasProcessedCheck = strings.Join(hasProcessedCheck, ";")
	platformData.HasProcessedConfirm = strings.Join(hasProcessedConfirm, ";")
	platformData.Submit = strings.Join(GetSubmitCount(user), ";")
	platformData.Draft = strings.Join(GetDraftCount(user), ";")
	platformData.Collect = strings.Join(GetCollectCount(user), ";")
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
func GetPendingCount(userRoles []string) (pendingTask, pendingApprove, pendingCheck, pendingConfirm, pending []string) {
	userRolesFilterSql, userRolesFilterParam := dao.CreateListParams(userRoles, "")
	var pendingTaskParam, pendingApproveParam, pendingCheckParam, pendingConfirmParam []interface{}
	var pTaskSQL, pendingApproveSQL, pendingCheckSQL, pendingConfirmSQL string
	var pTaskCount, pendingApproveCount, pendingCheckCount, pendingConfirmCount int
	for i := 0; i < len(templateTypeArr); i++ {
		pTaskSQL, pendingTaskParam = pendingTaskSQL(templateTypeArr[i], userRolesFilterSql, userRolesFilterParam, models.TaskTypeImplement)
		pTaskCount = dao.QueryCount(pTaskSQL, pendingTaskParam...)
		pendingTask = append(pendingTask, strconv.Itoa(pTaskCount))

		pendingApproveSQL, pendingApproveParam = pendingTaskSQL(templateTypeArr[i], userRolesFilterSql, userRolesFilterParam, models.TaskTypeApprove)
		pendingApproveCount = dao.QueryCount(pendingApproveSQL, pendingApproveParam)
		pendingApprove = append(pendingApprove, strconv.Itoa(pendingApproveCount))

		pendingCheckSQL, pendingCheckParam = pendingTaskSQL(templateTypeArr[i], userRolesFilterSql, userRolesFilterParam, models.TaskTypeCheck)
		pendingCheckCount = dao.QueryCount(pendingCheckSQL, pendingCheckParam)
		pendingCheck = append(pendingCheck, strconv.Itoa(pendingCheckCount))

		pendingConfirmSQL, pendingConfirmParam = pendingTaskSQL(templateTypeArr[i], userRolesFilterSql, userRolesFilterParam, models.TaskTypeConfirm)
		pendingConfirmCount = dao.QueryCount(pendingConfirmSQL, pendingConfirmParam)
		pendingConfirm = append(pendingConfirm, strconv.Itoa(pendingConfirmCount))

		pending = append(pending, strconv.Itoa(pTaskCount+pendingApproveCount+pendingCheckCount+pendingConfirmCount))
	}
	return
}

func hasProcessedRequestSQL(templateType int, user string) (sql string, queryParam []interface{}) {
	sql = "select id from request where del_flag= 0 and type = ? and handler = ?  and status not in ('Pending','Draft') " +
		"union select id from request where del_flag = 0 and type = ? and handler = ? and status='Draft' and rollback_desc is not null"
	queryParam = append([]interface{}{templateType, user, templateType, user})
	return
}

/*func hasProcessedTaskSQL(templateType int, user string) (sql string, queryParam []interface{}) {
	sql = "select id from task where handler= ? and del_flag = 0 and status ='done' and template_type = ? and request is not null"
	queryParam = append([]interface{}{user, templateType})
	return
}*/

func getPlatRequestSQL(where, sql string) string {
	return fmt.Sprintf("select * from (select r.id,r.name,r.cache,r.report_time,r.del_flag,rt.id as template_id,rt.name as template_name,rt.parent_id,"+
		"r.proc_instance_id,r.operator_obj,rt.proc_def_id,r.type as type,rt.proc_def_key,rt.operator_obj_type,r.role,r.status,r.rollback_desc,r.created_by,r.handler,r.created_time,r.updated_time,rt.proc_def_name,"+
		"r.expect_time,r.revoke_flag,r.confirm_time as approval_time from request r join request_template rt on r.request_template = rt.id ) t %s and id in (%s) ", where, sql)
}

func getPlatTaskSQL(where, sql string) string {
	return fmt.Sprintf("select * from (select r.id,r.name,r.cache,r.report_time,r.del_flag,rt.id as template_id,rt.name as template_name,rt.parent_id,r.proc_instance_id,r.operator_obj,rt.proc_def_id,r.type as type,rt.proc_def_key,rt.operator_obj_type,r.role,r.status,r.rollback_desc,r.created_by,r.created_time,r.updated_time,rt.proc_def_name,r.expect_time,r.revoke_flag,t.id as task_id,t.name as task_name,t.task_created_time,t.task_approval_time as task_approval_time,t.updated_time as task_updated_time,t.status as task_status,t.expire_time as task_expect_time,t.task_handler as task_handler,t.task_handle_id from (%s) t left join request r on t.request=r.id join request_template rt on r.request_template = rt.id) temp %s", sql, where)
}

func pendingTaskSQL(templateType int, userRolesFilterSql string, userRolesFilterParam []interface{}, taskType models.TaskType) (sql string, queryParam []interface{}) {
	queryParam = []interface{}{}
	sql = "select * from (select t.id,t.request,t.template_type,t.name,t.type,t.created_time as task_created_time,th.updated_time as task_approval_time,t.updated_time,t.status,t.expire_time,th.role as task_handle_role,th.id as task_handle_id,th.handler as task_handler,t.del_flag from task t right join task_handle th ON t.id = th.task) tha where del_flag = 0 and status <> 'done' and template_type = ? and type = ? and task_handle_role in (" + userRolesFilterSql + ")"
	queryParam = append([]interface{}{templateType, taskType}, userRolesFilterParam...)
	return
}

func hasProcessedTaskSQL(templateType int, user string, taskType models.TaskType) (sql string, queryParam []interface{}) {
	queryParam = []interface{}{}
	sql = "select * from (select t.id,t.request,t.template_type,t.name,t.type,t.created_time as task_created_time,th.updated_time as task_approval_time,t.updated_time,t.status,t.expire_time,th.role as task_handle_role,th.id as task_handle_id,th.handler as task_handler,t.del_flag from task t right join task_handle th ON t.id = th.task) tha where del_flag = 0 and status = 'done' and template_type = ? and type = ? and task_handler =?"
	queryParam = append([]interface{}{templateType, taskType, user})
	return
}

// GetHasProcessedCount 统计已处理,包括:(1)处理定版 (2) 任务已审批
func GetHasProcessedCount(user string) (hasProcessedTask, hasProcessedApprove, hasProcessedCheck, hasProcessedConfirm, hasProcessed []string) {
	var hasProcessedTaskParam, hasProcessedApproveParam, hasProcessedCheckParam, hasProcessedConfirmParam []interface{}
	var hpTaskSQL, hasProcessedApproveSQL, hasProcessedCheckSQL, hasProcessedConfirmSQL string
	var hpTaskCount, hasProcessedApproveCount, hasProcessedCheckCount, hasProcessedConfirmCount int
	for i := 0; i < len(templateTypeArr); i++ {
		hpTaskSQL, hasProcessedTaskParam = hasProcessedTaskSQL(templateTypeArr[i], user, models.TaskTypeImplement)
		hpTaskCount = dao.QueryCount(hpTaskSQL, hasProcessedTaskParam...)
		hasProcessedTask = append(hasProcessedTask, strconv.Itoa(hpTaskCount))

		hasProcessedApproveSQL, hasProcessedApproveParam = hasProcessedTaskSQL(templateTypeArr[i], user, models.TaskTypeApprove)
		hasProcessedApproveCount = dao.QueryCount(hasProcessedApproveSQL, hasProcessedApproveParam...)
		hasProcessedApprove = append(hasProcessedApprove, strconv.Itoa(hasProcessedApproveCount))

		hasProcessedCheckSQL, hasProcessedCheckParam = hasProcessedTaskSQL(templateTypeArr[i], user, models.TaskTypeCheck)
		hasProcessedCheckCount = dao.QueryCount(hasProcessedCheckSQL, hasProcessedCheckParam...)
		hasProcessedCheck = append(hasProcessedCheck, strconv.Itoa(hasProcessedCheckCount))

		hasProcessedConfirmSQL, hasProcessedConfirmParam = hasProcessedTaskSQL(templateTypeArr[i], user, models.TaskTypeConfirm)
		hasProcessedConfirmCount = dao.QueryCount(hasProcessedConfirmSQL, hasProcessedConfirmParam...)
		hasProcessedConfirm = append(hasProcessedConfirm, strconv.Itoa(hasProcessedConfirmCount))

		hasProcessed = append(hasProcessed, strconv.Itoa(hpTaskCount+hasProcessedApproveCount+hasProcessedCheckCount+hasProcessedConfirmCount))
	}
	return
}

// GetSubmitCount  统计用户提交
func GetSubmitCount(user string) (resultArr []string) {
	var queryParam []interface{}
	var sql string
	for i := 0; i < len(templateTypeArr); i++ {
		sql, queryParam = submitSQL(0, templateTypeArr[i], user)
		resultArr = append(resultArr, strconv.Itoa(dao.QueryCount(sql, queryParam...)))
	}
	return
}

func submitSQL(rollback, templateType int, user string) (sql string, queryParam []interface{}) {
	sql = "select id from request where del_flag=0 and created_by = ? and type = ? and (status != 'Draft' or ( status = 'Draft' and rollback_desc is not null ) or (status = 'Draft' and revoke_flag = 1))"
	if rollback == 1 {
		// 被退回
		sql = "select id from request where del_flag=0 and created_by = ? and type = ? and status = 'Draft' and rollback_desc is not null"
	} else if rollback == 2 {
		// 其他
		sql = "select id from request where del_flag=0 and created_by = ? and type = ? and status != 'Draft'"
	} else if rollback == 3 {
		// 撤销
		sql = "select id from request where del_flag=0 and created_by = ? and type = ? and status = 'Draft' and revoke_flag = 1 "
	}
	queryParam = append([]interface{}{user, templateType})
	return
}

// GetDraftCount 统计用户暂存
func GetDraftCount(user string) (resultArr []string) {
	var queryParam []interface{}
	var sql string
	for i := 0; i < len(templateTypeArr); i++ {
		sql, queryParam = draftSQL(templateTypeArr[i], user)
		resultArr = append(resultArr, strconv.Itoa(dao.QueryCount(sql, queryParam...)))
	}
	return
}

func draftSQL(templateType int, user string) (sql string, queryParam []interface{}) {
	sql = "select id from request where del_flag=0 and created_by = ? and status = 'Draft' and type = ? and rollback_desc is null and revoke_flag = 0"
	queryParam = append([]interface{}{user, templateType})
	return
}

// GetCollectCount 统计用户收藏
func GetCollectCount(user string) (resultArr []string) {
	var queryParam []interface{}
	var sql string
	for i := 0; i < len(templateTypeArr); i++ {
		sql, queryParam = collectSQL(templateTypeArr[i], user)
		resultArr = append(resultArr, strconv.Itoa(dao.QueryCount(sql, queryParam...)))
	}
	return
}

func collectSQL(templateType int, user string) (sql string, queryParam []interface{}) {
	sql = "select id from collect_template where user = ? and type= ?"
	queryParam = append([]interface{}{user, templateType})
	return
}

// DataList 首页工作台数据列表
func DataList(param *models.PlatformRequestParam, userRoles []string, userToken, user, language string) (pageInfo models.PageInfo, rowData []*models.PlatformDataObj, err error) {
	// 先拼接查询条件
	var sql string
	var queryParam []interface{}
	var taskType = getTaskTypeByType(param.Type)
	where := transPlatConditionToSQL(param)
	userRolesFilterSql, userRolesFilterParam := dao.CreateListParams(userRoles, "")
	switch param.Tab {
	case "pending":
		sql, queryParam = pendingTaskSQL(param.Action, userRolesFilterSql, userRolesFilterParam, taskType)
		pageInfo, rowData, err = getPlatData(models.PlatDataParam{Param: param.CommonRequestParam, QueryParam: queryParam, UserToken: userToken}, getPlatTaskSQL(where, sql), language, true)
		return
	case "hasProcessed":
		sql, queryParam = hasProcessedTaskSQL(param.Action, user, taskType)
		pageInfo, rowData, err = getPlatData(models.PlatDataParam{Param: param.CommonRequestParam, QueryParam: queryParam, UserToken: userToken}, getPlatTaskSQL(where, sql), language, true)
		return
	case "submit":
		sql, queryParam = submitSQL(param.Rollback, param.Action, user)
	case "draft":
		sql, queryParam = draftSQL(param.Action, user)
	default:
		err = fmt.Errorf("request param err,tab:%s", param.Tab)
		return
	}
	pageInfo, rowData, err = getPlatData(models.PlatDataParam{Param: param.CommonRequestParam, QueryParam: queryParam, UserToken: userToken}, getPlatRequestSQL(where, sql), language, true)
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
			row.HandleRole,
			strconv.Itoa(row.Progress) + "%",
			getRequestStayTime(row, days),
			row.ExpectTime,
			row.TemplateName + "【" + version + "】",
			row.ProcDefName,
			row.OperatorObjType,
			row.OperatorObj,
			row.CreatedBy,
			row.Role,
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
		case string(models.RequestStatusInProgressTimeOuted):
			return "节点超时"
		case string(models.RequestStatusFaulted):
			return "自动退出"
		}
	}
	return status
}

func getRequestStayTime(dataObject *models.PlatformDataObj, format string) string {
	return fmt.Sprintf("%d%s/%d%s", dataObject.RequestStayTime, format, dataObject.RequestStayTimeTotal, format)
}

// calcRequestStayTime 计算请求/任务停留时长
func calcRequestStayTime(dataObject *models.PlatformDataObj) {
	var err error
	var reportTime, requestExpectTime, taskCreateTime, taskExpectTime, taskApprovalTime time.Time
	loc, _ := time.LoadLocation("Local")
	if dataObject.Status == string(models.RequestStatusDraft) {
		return
	}
	// 计算任务停留时长
	if dataObject.TaskId != "" && dataObject.TaskExpectTime != "" && dataObject.TaskCreatedTime != "" {
		taskExpectTime, _ = time.ParseInLocation(models.DateTimeFormat, dataObject.TaskExpectTime, loc)
		taskCreateTime, _ = time.ParseInLocation(models.DateTimeFormat, dataObject.TaskCreatedTime, loc)
		if dataObject.TaskApprovalTime != "" && dataObject.TaskStatus == "done" {
			taskApprovalTime, _ = time.ParseInLocation(models.DateTimeFormat, dataObject.TaskApprovalTime, loc)
			dataObject.TaskStayTime = int(math.Ceil(taskApprovalTime.Sub(taskCreateTime).Hours() * 1.00 / 24.00))
		} else {
			dataObject.TaskStayTime = int(math.Ceil(time.Now().Local().Sub(taskCreateTime).Hours() * 1.00 / 24.00))
		}
		dataObject.TaskStayTimeTotal = int(math.Ceil(taskExpectTime.Sub(taskCreateTime).Hours() * 1.00 / 24.00))
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
		// 向上取整
		dataObject.RequestStayTime = int(math.Ceil(updateTime.Sub(reportTime).Hours() * 1.00 / 24.00))
	} else {
		dataObject.RequestStayTime = int(math.Ceil(time.Now().Local().Sub(reportTime).Hours() * 1.00 / 24.00))
	}
	dataObject.RequestStayTimeTotal = int(math.Ceil(requestExpectTime.Sub(reportTime).Hours() * 1.00 / 24.00))
}

func getPlatData(req models.PlatDataParam, newSQL, language string, page bool) (pageInfo models.PageInfo, rowsData []*models.PlatformDataObj, err error) {
	var operatorObjTypeMap = make(map[string]string)
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
		var actions []*dao.ExecAction
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
			if platformDataObj.Status == "Draft" {
				platformDataObj.CurNode = WaitCommit
			} else if platformDataObj.Status == "Pending" {
				platformDataObj.CurNode = RequestPending
			}
			if platformDataObj.ProcInstanceId != "" {
				platformDataObj.Progress, platformDataObj.CurNode = getCurNodeName(platformDataObj.ProcInstanceId, req.UserToken, language)
			}
			if strings.Contains(platformDataObj.Status, "InProgress") && platformDataObj.ProcInstanceId != "" {
				newStatus := getInstanceStatus(platformDataObj.ProcInstanceId, req.UserToken, language)
				if newStatus == "InternallyTerminated" {
					newStatus = "Termination"
				}
				if newStatus != "" && newStatus != platformDataObj.Status {
					actions = append(actions, &dao.ExecAction{Sql: "update request set status=?,updated_time=? where id=?",
						Param: []interface{}{newStatus, time.Now().Format(models.DateTimeFormat), platformDataObj.Id}})
					platformDataObj.Status = newStatus
				}
			}
			if collectMap[platformDataObj.ParentId] {
				platformDataObj.CollectFlag = 1
			}
			if platformDataObj.OperatorObjType == "" {
				platformDataObj.OperatorObjType = operatorObjTypeMap[platformDataObj.ProcDefKey]
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
			platformDataObj.HandleRole, platformDataObj.Handler = getRequestHandler(platformDataObj.Id)
			// 计算请求/任务停留时长
			calcRequestStayTime(platformDataObj)
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

func getCurNodeName(instanceId, userToken, language string) (progress int, curNode string) {
	var total int
	processInstance, err := GetProcDefService().GetProcessDefineInstance(instanceId, userToken, language)
	if err != nil || processInstance == nil || len(processInstance.TaskNodeInstances) == 0 {
		return
	}
	// 统计完成进度 ,已完成/总数, 编号不为空
	for _, v := range processInstance.TaskNodeInstances {
		if v.OrderedNo != "" {
			if v.Status == string(models.RequestStatusCompleted) {
				progress++
			}
			total++
		}
	}
	progress = int(math.Floor(float64(progress)/float64(total)*100 + 0.5))
	switch processInstance.Status {
	case "Completed":
		curNode = CurNodeCompleted
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
	where = where + fmt.Sprintf(" and type = %d", param.Action)
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
func GetRequestProgress(requestId, userToken, language string) (rowData *models.RequestProgressObj, err error) {
	var request models.RequestTable
	var requestTemplate *models.RequestTemplateTable
	var pendingRole, pendingHandler string
	var requestTemplateService = GetRequestTemplateService()
	//var taskList []*models.TaskTable
	rowData = &models.RequestProgressObj{RequestProgress: []*models.ProgressObj{}, ApprovalProgress: []*models.ProgressObj{}, TaskProgress: []*models.ProgressObj{}}
	request, err = GetSimpleRequest(requestId)
	if err != nil {
		return
	}
	requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(request.RequestTemplate)
	if err != nil {
		return
	}
	if requestTemplate == nil {
		return
	}
	// 添加提交请求
	rowData.RequestProgress = append(rowData.RequestProgress, &models.ProgressObj{Node: SendRequest, Handler: request.CreatedBy})
	// 配置了请求定版
	if requestTemplate.CheckSwitch {
		pendingRole, pendingHandler = requestTemplateService.GetRequestPendingRoleAndHandler(requestTemplateService.GetDtoByRequestTemplate(requestTemplate))
		// 没有处理人,展示处理角色
		if pendingHandler == "" {
			pendingHandler = pendingRole
		}
		rowData.RequestProgress = append(rowData.RequestProgress, &models.ProgressObj{Node: RequestPending, Handler: pendingHandler})
	}
	/*approvalTemplateList, err = GetApprovalTemplateService().ListApprovalTemplates(requestTemplate.Id)
	  if len(approvalTemplateList) > 0 {
	  	rowData.RequestProgress = append(rowData.RequestProgress, &models.ProgressObj{Node: Approval, Handler: ""})
	  	approvalList, err = GetApprovalService().ListApprovals(requestId)
	  	// 新建请求还未创建审批
	  	if len(approvalList) == 0 {

	  	}
	  }*/
	// 任务进度
	/*taskTemplateList, err = GetTaskTemplateService().ListTaskTemplates(requestTemplate.Id)
	  if len(taskTemplateList) > 0 {
	  	rowData.RequestProgress = append(rowData.RequestProgress, &models.ProgressObj{Node: Task, Handler: ""})
	  }*/
	// 请求完成
	rowData.RequestProgress = append(rowData.RequestProgress, &models.ProgressObj{
		Node:    RequestComplete,
		Handler: "",
		Status:  int(models.ProgressStatusNotStart),
	})
	return
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
	err := dao.X.SQL("select rt.name as  template_name,rt.status,rtg.name as template_group_name,rt.version,rt.proc_def_id,rt.expire_day from request_template rt join "+
		"request_template_group rtg on rt.group = rtg.id where rt.id= ?", request.RequestTemplate).Find(&tmpTemplate)
	if err != nil {
		return
	}
	if len(tmpTemplate) == 0 {
		err = fmt.Errorf("can not find request_template with id:%s ", request.Id)
		return
	}
	template := tmpTemplate[0]
	form.Id = request.Id
	form.Name = request.Name
	form.RequestType = request.Type
	form.CreatedTime = request.CreatedTime
	form.ExpectTime = request.ExpectTime
	form.CreatedBy = request.CreatedBy
	form.Role = request.Role
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
	if request.ProcInstanceId != "" {
		form.Progress, form.CurNode = getCurNodeName(request.ProcInstanceId, userToken, language)
	}
	if template.ProcDefId != "" {
		form.AssociationWorkflow = true
	}
	_, form.Handler = getRequestHandler(request.Id)
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
	}
	return
}

// getRequestHandler 获取请求处理人,如果处于任务执行状态,查询任务处理人
func getRequestHandler(requestId string) (role, handler string) {
	request, _ := GetSimpleRequest(requestId)
	if request.Status == "Draft" || request.Status == "Pending" {
		// 请求在定版状态,从模板角色表中读取
		rtRoleMap := getRequestTemplateMGMTRole()
		roles := rtRoleMap[request.RequestTemplate]
		if len(roles) > 0 {
			role = roles[0]
		}
		handler = request.Handler
		return
	}
	// 请求在任务状态,需要从模板配置的任务表中获取
	taskTemplateMap, _ := GetTaskTemplateService().getTaskTemplateHandler(request.RequestTemplate)
	if len(taskTemplateMap) > 0 {
		taskMap, _ := getTaskMapByRequestId(requestId)
		if len(taskMap) > 0 {
			for _, task := range taskMap {
				if task.Status != "done" && taskTemplateMap[task.TaskTemplate] != nil {
					//taskTemplate := taskTemplateMap[task.TaskTemplate]
					//role = taskTemplate.Role
					// 任务处理人已任务处理为主,可以通过认领转给我修改.空的时候才取模板配置值
					handler = task.Handler
					if handler == "" {
						//handler = taskTemplate.Handler
					}
					break
				}
			}
		}
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
func RequestConfirm(param models.RequestConfirmParam, user string) (err error) {
	var taskList []*models.TaskTable
	var actions []*dao.ExecAction
	var action *dao.ExecAction
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
	// 更新请求确认任务设置为已完成
	action = &dao.ExecAction{Sql: "update task set status = ?,description = ? ,updated_by = ?,updated_time = ? where id = ?"}
	action.Param = []interface{}{models.TaskStatusDone, param.Notes, user, now, param.TaskId}
	actions = append(actions, action)
	// 更新请求表状态,设置为完成
	action = &dao.ExecAction{Sql: "update request set status = ?,complete_status = ?,updated_by = ?,updated_time= ? where id = ?"}
	action.Param = []interface{}{models.RequestStatusCompleted, param.CompleteStatus, user, now, param.Id}
	actions = append(actions, action)
	err = dao.Transaction(actions)
	return
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

// HandleRequestCheck 处理确认定版
func (s *RequestService) HandleRequestCheck(request models.RequestTable, operator, bindCache, userToken, language string) (err error) {
	now := time.Now().Format(models.DateTimeFormat)
	var actions []*dao.ExecAction
	var approvalActions []*dao.ExecAction
	var action *dao.ExecAction
	var requestTemplate *models.RequestTemplateTable
	var submitTaskTemplateList, checkTaskTemplateList []*models.TaskTemplateTable
	var taskHandleTemplateList []*models.TaskHandleTemplateTable
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
	// 新增提交请求任务
	submitTaskId := "su_" + guid.CreateGuid()
	action = &dao.ExecAction{Sql: "insert into task(id,name,expire_time,template_type,status,request,task_template,type,created_by,created_time,updated_by,updated_time) values (?,?,?,?,?,?,?,?,?,?,?,?)"}
	action.Param = []interface{}{submitTaskId, "submit", expireTime, request.Type, models.TaskStatusDone, request.Id, submitTaskTemplateList[0].Id, models.TaskTypeSubmit, "system", now, "system", now}
	actions = append(actions, action)

	// 新增提交请求处理人
	action = &dao.ExecAction{Sql: "insert into task_handle(id,task,role,handler,created_time,updated_time) values (?,?,?,?,?,?)"}
	action.Param = []interface{}{guid.CreateGuid(), submitTaskId, request.Role, request.CreatedBy, now, now}
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
		checkTime := time.Now().Add(time.Second * 1).Format(models.DateTimeFormat)
		checkExpireTime := calcExpireTime(checkTime, checkTaskTemplateList[0].ExpireDay)
		action = &dao.ExecAction{Sql: "insert into task(id,name,expire_time,template_type,status,request,task_template,type,created_by,created_time,updated_by,updated_time) values (?,?,?,?,?,?,?,?,?,?,?,?)"}
		action.Param = []interface{}{checkTaskId, "check", checkExpireTime, request.Type, models.TaskStatusCreated, request.Id, checkTaskTemplateList[0].Id, models.TaskTypeCheck, operator, checkTime, operator, checkTime}
		actions = append(actions, action)

		// 新增确认定版处理人
		dao.X.SQL("select * from task_handle_template where task_template = ?", checkTaskTemplateList[0].Id).Find(&taskHandleTemplateList)
		if len(taskHandleTemplateList) > 0 {
			action = &dao.ExecAction{Sql: "insert into task_handle(id,task_handle_template,task,role,handler,created_time,updated_time) values (?,?,?,?,?,?,?)"}
			action.Param = []interface{}{guid.CreateGuid(), taskHandleTemplateList[0].Id, checkTaskId, taskHandleTemplateList[0].Role, taskHandleTemplateList[0].Handler, checkTime, checkTime}
			actions = append(actions, action)
		}
		err = dao.Transaction(actions)
		return
	}

	// 没有配置定版,请求继续往后面走
	approvalActions, err = s.HandleRequestApproval(request, userToken, language)
	if err != nil {
		return
	}
	if len(approvalActions) > 0 {
		actions = append(actions, approvalActions...)
	}
	err = dao.Transaction(actions)
	return
}

// HandleRequestApproval 处理请求审批
func (s *RequestService) HandleRequestApproval(request models.RequestTable, userToken, language string) (actions []*dao.ExecAction, err error) {
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
		return s.HandleRequestTask(request, userToken, language)
	}
	for _, taskTemplate := range taskTemplateList {
		dao.X.SQL("select * from task where request = ? and task_template = ? order by created_time desc", request.Id, taskTemplate.Id).Find(&taskList)
		if len(taskList) > 0 {
			// 取最新的任务
			if taskList[0].Status == string(models.TaskStatusDone) {
				// 任务已完成,continue
				continue
			}
			// 任务状态未完成,break等待当前任务处理完
			break
		} else {
			// 模版没有对应的的任务,需要创建当前任务并设置审批角色和人,同时根据审批方式设置审批状态
			newTaskId = "ap_" + guid.CreateGuid()
			taskExpireTime := calcExpireTime(now, taskTemplate.ExpireDay)
			// 审批模板配置自动通过,设置当前审批完成,并且直接下一个审批
			if taskTemplate.HandleMode == string(models.TaskTemplateHandleModeAuto) {
				action = &dao.ExecAction{Sql: "insert into task (id,name,expire_time,template_type,description,status,request,task_template,type,sort,created_by,created_time) values(?,?,?,?,?,?,?,?,?,?,?,?)"}
				action.Param = []interface{}{newTaskId, taskTemplate.Name, taskExpireTime, request.Type, taskTemplate.Description, models.TaskStatusDone, request.Id, taskTemplate.Id, taskTemplate.Type, taskTemplate.Sort, "system", now}
				actions = append(actions, action)
				continue
			}

			// 新增任务
			action = &dao.ExecAction{Sql: "insert into task (id,name,expire_time,template_type,description,status,request,task_template,type,sort,created_by,created_time) values(?,?,?,?,?,?,?,?,?,?,?,?)"}
			action.Param = []interface{}{newTaskId, taskTemplate.Name, taskExpireTime, request.Type, taskTemplate.Description, models.TaskStatusCreated, request.Id, taskTemplate.Id, taskTemplate.Type, taskTemplate.Sort, "system", now}
			actions = append(actions, action)

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
	return s.HandleRequestTask(request, userToken, language)
}

// HandleRequestTask 处理任务
func (s *RequestService) HandleRequestTask(request models.RequestTable, userToken, language string) (actions []*dao.ExecAction, err error) {
	var taskTemplateList []*models.TaskTemplateTable
	var requestTemplate *models.RequestTemplateTable
	var taskList []*models.TaskTable
	var newTaskId string
	var action *dao.ExecAction
	now := time.Now().Format(models.DateTimeFormat)
	actions = []*dao.ExecAction{}
	if request.AssociationWorkflow && request.ProcInstanceId == "" && request.BindCache != "" {
		// 关联编排,调用编排启动
		var bindCache models.RequestCacheData
		json.Unmarshal([]byte(request.BindCache), &bindCache)
		_, err = StartRequestNew(request, userToken, language, bindCache)
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
		return s.HandleRequestConfirm(request)
	}
	for _, taskTemplate := range taskTemplateList {
		taskExpireTime := calcExpireTime(now, taskTemplate.ExpireDay)
		dao.X.SQL("select * from task where request = ? and task_template = ? order by created_time desc", request.Id, taskTemplate.Id).Find(&taskList)
		if len(taskList) > 0 {
			// 取最新的任务
			if taskList[0].Status == string(models.TaskStatusDone) {
				// 任务已完成,continue
				continue
			}
			// 任务状态未完成,break等待当前任务处理完
			break
		} else {
			// 模版没有对应的的任务,需要创建当前任务
			newTaskId = "im_" + guid.CreateGuid()
			// 新增任务
			action = &dao.ExecAction{Sql: "insert into task (id,name,expire_time,template_type,description,status,request,task_template,type,sort,created_by,created_time) values(?,?,?,?,?,?,?,?,?,?,?,?)"}
			action.Param = []interface{}{newTaskId, taskTemplate.Name, taskExpireTime, request.Type, taskTemplate.Description, models.TaskStatusCreated, request.Id, taskTemplate.Id, taskTemplate.Type, taskTemplate.Sort, "system", now}
			actions = append(actions, action)

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
	return s.HandleRequestConfirm(request)
}

// HandleRequestConfirm 处理请求确认
func (s *RequestService) HandleRequestConfirm(request models.RequestTable) (actions []*dao.ExecAction, err error) {
	var newTaskId string
	var action *dao.ExecAction
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
		err = fmt.Errorf("taskTemplate not find confirm template")
		return
	}
	// 新增任务
	taskExpireTime := calcExpireTime(now, taskTemplateList[0].ExpireDay)
	action = &dao.ExecAction{Sql: "insert into task (id,name,expire_time,template_type,status,request,task_template,type,created_by,created_time) values(?,?,?,?,?,?,?,?,?,?)"}
	action.Param = []interface{}{newTaskId, "confirm", taskExpireTime, request.Type, models.TaskStatusCreated, request.Id, taskTemplateList[0], models.TaskTypeConfirm, "system", now}
	actions = append(actions, action)
	// 更新请求表状态为请求确认
	actions = append(actions, &dao.ExecAction{Sql: "update request set status=?,updated_time=? where id=?", Param: []interface{}{string(models.RequestStatusConfirm), now, request.Id}})
	return
}
