package service

import (
	"encoding/json"
	"fmt"
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

type ProgressStatus int

const (
	InProgress                 ProgressStatus = 1 // 进行中
	NotStart                   ProgressStatus = 2 // 未开始
	Completed                  ProgressStatus = 3 // 已完成
	Fail                       ProgressStatus = 4 // 报错失败,被拒绝了
	AutoExitStatus             ProgressStatus = 5 // 自动退出
	InternallyTerminatedStatus ProgressStatus = 6 // 自动退出
)

const (
	WaitCommit           = "waitCommit"           // 等待提交
	SendRequest          = "sendRequest"          // 发送请求
	RequestPending       = "requestPending"       // 请求定版
	Approval             = "approval"             // 审批
	Task                 = "task"                 // 任务
	CurNodeCompleted     = "Completed"            // 完成
	RequestComplete      = "requestComplete"      // 请求完成
	AutoExit             = "autoExit"             // 自动退出
	InternallyTerminated = "internallyTerminated" // 手动终止
)

const (
	AutoNode = "autoNode" //自动节点
)

// GetRequestCount 工作台请求统计
func GetRequestCount(user string, userRoles []string) (platformData models.PlatformData, err error) {
	var pendingTask, pendingApprove, requestPending, requestConfirm, pending []string
	pendingTask, pendingApprove, requestPending, requestConfirm, pending = GetPendingCount(userRoles)
	platformData.Pending = strings.Join(pending, ";")
	platformData.PendingTask = strings.Join(pendingTask, ";")
	platformData.PendingApprove = strings.Join(pendingApprove, ";")
	platformData.RequestPending = strings.Join(requestPending, ";")
	platformData.RequestConfirm = strings.Join(requestConfirm, ";")
	platformData.HasProcessed = strings.Join(GetHasProcessedCount(user), ";")
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

// GetPendingCount 统计待处理,包括:(1)待定版 (2)任务待审批 分配给本组、本人的所有请求,此处按照角色去统计
func GetPendingCount(userRoles []string) (pendingTask, pendingApprove, requestPending, requestConfirm, pending []string) {
	userRolesFilterSql, userRolesFilterParam := dao.CreateListParams(userRoles, "")
	var requestQueryParam, taskQueryParam, approveQueryParam, requestConfirmQueryParam []interface{}
	var roleFilterList []string
	var requestSQL, taskSQL, approveSQL, requestConfirmSQL string
	var requestSQLCount, taskSQLCount, approveSQLCount, requestConfirmSQLCount int
	roleFilterSql := "1=1"
	if len(userRoles) > 0 {
		for _, v := range userRoles {
			roleFilterList = append(roleFilterList, "report_role like '%,"+v+",%'")
		}
		roleFilterSql = strings.Join(roleFilterList, " or ")
	}
	for i := 0; i < len(templateTypeArr); i++ {
		requestSQL, requestQueryParam = pendingRequestSQL(templateTypeArr[i], string(models.RequestStatusPending), userRolesFilterSql, userRolesFilterParam)
		taskSQL, taskQueryParam = pendingTaskSQL(templateTypeArr[i], userRolesFilterSql, userRolesFilterParam, roleFilterSql)
		approveSQL, approveQueryParam = pendingRequestSQL(templateTypeArr[i], string(models.RequestStatusInApproval), userRolesFilterSql, userRolesFilterParam)
		requestConfirmSQL, requestConfirmQueryParam = pendingRequestSQL(templateTypeArr[i], string(models.RequestStatusConfirm), userRolesFilterSql, userRolesFilterParam)
		requestSQLCount = dao.QueryCount(requestSQL, requestQueryParam...)
		taskSQLCount = dao.QueryCount(taskSQL, taskQueryParam...)
		approveSQLCount = dao.QueryCount(approveSQL, approveQueryParam...)
		requestConfirmSQLCount = dao.QueryCount(requestConfirmSQL, requestConfirmQueryParam...)
		requestPending = append(requestPending, strconv.Itoa(requestSQLCount))
		pendingTask = append(pendingTask, strconv.Itoa(taskSQLCount))
		pendingApprove = append(pendingApprove, strconv.Itoa(approveSQLCount))
		requestConfirm = append(requestConfirm, strconv.Itoa(requestConfirmSQLCount))
		pending = append(pending, strconv.Itoa(requestSQLCount+approveSQLCount+requestConfirmSQLCount+taskSQLCount))
	}
	return
}

func pendingRequestSQL(templateType int, status, userRolesFilterSql string, userRolesFilterParam []interface{}) (sql string, queryParam []interface{}) {
	sql = "select id from request where del_flag = 0 and status = ? and type = ? and request_template in (select id " +
		"from request_template where  id in (select request_template from request_template_role where role_type= 'MGMT' " +
		"and `role` in (" + userRolesFilterSql + "))) "
	queryParam = append([]interface{}{status, templateType}, userRolesFilterParam...)
	return
}

func pendingTaskSQL(templateType int, userRolesFilterSql string, userRolesFilterParam []interface{}, roleFilterSql string) (sql string, queryParam []interface{}) {
	sql = fmt.Sprintf("select id from task where del_flag = 0 and status <> 'done' and template_type = ? and task_template "+
		"in (select task_template from task_template_role where role_type='USE' and `role` in ("+userRolesFilterSql+")"+
		" union select id from task where task_template is null and (%s) and del_flag=0 and template_type = ?)", roleFilterSql)
	queryParam = append(append([]interface{}{templateType}, userRolesFilterParam...), templateType)
	return
}

// GetHasProcessedCount 统计已处理,包括:(1)处理定版 (2) 任务已审批
func GetHasProcessedCount(user string) (resultArr []string) {
	var requestSQL, taskSQL string
	var requestQueryParam, taskQueryParam []interface{}
	for i := 0; i < len(templateTypeArr); i++ {
		requestSQL, requestQueryParam = hasProcessedRequestSQL(templateTypeArr[i], user)
		taskSQL, taskQueryParam = hasProcessedTaskSQL(templateTypeArr[i], user)
		resultArr = append(resultArr, strconv.Itoa(dao.QueryCount(requestSQL, requestQueryParam...)+dao.QueryCount(taskSQL, taskQueryParam...)))
	}
	return
}

func hasProcessedRequestSQL(templateType int, user string) (sql string, queryParam []interface{}) {
	sql = "select id from request where del_flag= 0 and type = ? and handler = ?  and status not in ('Pending','Draft') " +
		"union select id from request where del_flag = 0 and type = ? and handler = ? and status='Draft' and rollback_desc is not null"
	queryParam = append([]interface{}{templateType, user, templateType, user})
	return
}

func hasProcessedTaskSQL(templateType int, user string) (sql string, queryParam []interface{}) {
	sql = "select id from task where handler= ? and del_flag = 0 and status ='done' and template_type = ? and request is not null"
	queryParam = append([]interface{}{user, templateType})
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
	var templateType int
	var sql, status string
	var queryParam []interface{}
	where := transPlatConditionToSQL(param)
	if param.Action == 1 {
		templateType = 1
	} else if param.Action == 2 {
		templateType = 0
	}
	userRolesFilterSql, userRolesFilterParam := dao.CreateListParams(userRoles, "")
	switch param.Tab {
	case "pending":
		var roleFilterList []string
		roleFilterSql := "1=1"
		if len(userRoles) > 0 {
			for _, v := range userRoles {
				roleFilterList = append(roleFilterList, "report_role like '%,"+v+",%'")
			}
			roleFilterSql = strings.Join(roleFilterList, " or ")
		}
		if param.Type == 1 {
			sql, queryParam = pendingTaskSQL(templateType, userRolesFilterSql, userRolesFilterParam, roleFilterSql)
			pageInfo, rowData, err = getPlatData(models.PlatDataParam{Param: param.CommonRequestParam, QueryParam: queryParam, UserToken: userToken}, getPlatTaskSQL(where, sql), language, true)
			return
		} else {
			switch param.Type {
			case 1:
				status = string(models.RequestStatusPending)
			case 3:
				status = string(models.RequestStatusInApproval)
			case 4:
				status = string(models.RequestStatusConfirm)
			}
			sql, queryParam = pendingRequestSQL(templateType, status, userRolesFilterSql, userRolesFilterParam)
		}
	case "hasProcessed":
		if param.Type == 1 {
			sql, queryParam = hasProcessedRequestSQL(templateType, user)
		} else if param.Type == 2 {
			sql, queryParam = hasProcessedTaskSQL(templateType, user)
			pageInfo, rowData, err = getPlatData(models.PlatDataParam{Param: param.CommonRequestParam, QueryParam: queryParam, UserToken: userToken}, getPlatTaskSQL(where, sql), language, true)
			return
		}
	case "submit":
		sql, queryParam = submitSQL(param.Rollback, templateType, user)
	case "draft":
		sql, queryParam = draftSQL(templateType, user)
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

func getPlatRequestSQL(where, sql string) string {
	return fmt.Sprintf("select * from (select r.id,r.name,r.cache,r.report_time,r.del_flag,rt.id as template_id,rt.name as template_name,rt.parent_id,"+
		"r.proc_instance_id,r.operator_obj,rt.proc_def_id,r.type as type,rt.proc_def_key,rt.operator_obj_type,r.role,r.status,r.rollback_desc,r.created_by,r.handler,r.created_time,r.updated_time,rt.proc_def_name,"+
		"r.expect_time,r.revoke_flag,r.confirm_time as approval_time from request r join request_template rt on r.request_template = rt.id ) t %s and id in (%s) ", where, sql)
}

func getPlatTaskSQL(where, sql string) string {
	return fmt.Sprintf("select * from (select r.id,r.name,r.cache,r.report_time,r.del_flag,rt.id as template_id,rt.name as template_name,rt.parent_id,"+
		"r.proc_instance_id,r.operator_obj,rt.proc_def_id,r.type as type,rt.proc_def_key,rt.operator_obj_type,r.role,r.status,r.rollback_desc,r.created_by,r.created_time,r.updated_time,rt.proc_def_name,"+
		"r.expect_time,r.revoke_flag,t.id as task_id,t.name as task_name,t.report_time as task_created_time,t.updated_time as task_approval_time,t.updated_time as task_updated_time,t.status as task_status,t.expire_time as task_expect_time,t.handler as task_handler "+
		"from request r join request_template rt on r.request_template = rt.id left join task t on r.id = t.request) temp %s and task_id in (%s) ", where, sql)
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
	if param.Action == 1 {
		where = where + " and type = 1"
	} else if param.Action == 2 {
		where = where + " and type = 0"
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
	if param.TaskReportStartTime != "" && param.TaskReportEndTime != "" {
		where = where + " and task_created_time >= '" + param.TaskReportStartTime + "' and task_created_time <= '" + param.TaskReportEndTime + "'"
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

// getTaskApproveHandler 获取任务审批人,有人返回人,没人返回审批角色
func getTaskApproveHandler(requestId string, result models.TaskTemplateDto) string {
	// 审批人以 任务表审批人为主,任务可以认领转给我会修改任务审批人
	var taskList []*models.TaskTable
	dao.X.SQL("select name,handler,node_def_id,node_name from task where request = ?", requestId).Find(&taskList)
	if len(taskList) > 0 {
		for _, task := range taskList {
			if task.NodeDefId == result.NodeDefId && task.Handler != "" {
				return task.Handler
			}
		}
	}
	if result.Handler != "" {
		return result.Handler
	}
	if len(result.USERoleObjs) > 0 {
		return result.USERoleObjs[0].DisplayName
	}
	if len(result.USERoles) > 0 {
		return result.USERoles[0]
	}
	return ""
}

// GetRequestProgress  请求已创建时,获取请求进度
func GetRequestProgress(requestId, userToken, language string) (rowData *models.RequestProgressObj, err error) {
	var request models.RequestTable
	var requestTemplate *models.RequestTemplateTable
	var pendingRole, pendingHandler string
	var approvalTemplateList []*models.ApprovalTemplateDto
	var approvalList []*models.ApprovalDto
	var taskTemplateList []*models.TaskTemplateDto
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
	if requestTemplate.PendingSwitch {
		pendingRole, pendingHandler = GetRequestTemplateService().GetRequestPendingRoleAndHandler(requestTemplate)
		// 没有处理人,展示处理角色
		if pendingHandler == "" {
			pendingHandler = pendingRole
		}
		rowData.RequestProgress = append(rowData.RequestProgress, &models.ProgressObj{Node: RequestPending, Handler: pendingHandler})
	}
	approvalTemplateList, err = GetApprovalTemplateService().ListApprovalTemplates(requestTemplate.Id)
	if len(approvalTemplateList) > 0 {
		rowData.RequestProgress = append(rowData.RequestProgress, &models.ProgressObj{Node: Approval, Handler: ""})
		approvalList, err = GetApprovalService().ListApprovals(requestId)
		// 新建请求还未创建审批
		if len(approvalList) == 0 {

		}
	}
	// 任务进度
	taskTemplateList, err = GetTaskTemplateService().ListTaskTemplates(requestTemplate.Id)
	if len(taskTemplateList) > 0 {
		rowData.RequestProgress = append(rowData.RequestProgress, &models.ProgressObj{Node: Task, Handler: ""})
	}
	// 请求完成
	rowData.RequestProgress = append(rowData.RequestProgress, &models.ProgressObj{
		Node:    RequestComplete,
		Handler: "",
		Status:  int(NotStart),
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
	_, form.Handler = getRequestHandler(request.Id)
	if request.CustomFormCache != "" {
		err = json.Unmarshal([]byte(request.CustomFormCache), &customForm)
		if err != nil {
			log.Logger.Error("json Unmarshal", log.Error(err), log.String("CustomFormCache", request.CustomFormCache))
			return
		}
		form.CustomForm = customForm
	} else {
		var items []*models.FormItemTemplateTable
		dao.X.SQL("select * from form_item_template where form_template in (select form_template from request_template where id=?) order by item_group,sort", request.RequestTemplate).Find(&items)
		customForm.Title = items
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
	taskTemplateMap, _ := getTaskTemplateHandler(request.RequestTemplate)
	if len(taskTemplateMap) > 0 {
		taskMap, _ := getTaskMapByRequestId(requestId)
		if len(taskMap) > 0 {
			for _, task := range taskMap {
				if task.Status != "done" && taskTemplateMap[task.TaskTemplate] != nil {
					taskTemplate := taskTemplateMap[task.TaskTemplate]
					role = taskTemplate.Role
					// 任务处理人已任务处理为主,可以通过认领转给我修改.空的时候才取模板配置值
					handler = task.Handler
					if handler == "" {
						handler = taskTemplate.Handler
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
	for _, str := range arr {
		hashMap[str] = true
	}
	return hashMap
}