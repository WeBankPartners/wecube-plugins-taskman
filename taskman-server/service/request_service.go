package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/exterror"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type RequestService struct {
	requestDao dao.RequestDao
}

const (
	InProgress                 ProgressStatus = 1 // 进行中
	NotStart                   ProgressStatus = 2 // 未开始
	Completed                  ProgressStatus = 3 // 已完成
	Fail                       ProgressStatus = 4 // 报错失败,被拒绝了
	AutoExitStatus             ProgressStatus = 5 // 自动退出
	InternallyTerminatedStatus ProgressStatus = 6 // 自动退出
)

type ProgressStatus int

const (
	WaitCommit           = "waitCommit"           // 等待提交
	SendRequest          = "sendRequest"          // 发送请求
	RequestPending       = "requestPending"       // 请求定版
	CurNodeCompleted     = "Completed"            // 完成
	RequestComplete      = "requestComplete"      // 请求完成
	AutoExit             = "autoExit"             // 自动退出
	InternallyTerminated = "internallyTerminated" // 手动终止
)

const (
	AutoNode = "autoNode" //自动节点
)

var (
	requestIdLock   = new(sync.RWMutex)
	templateTypeArr = []int{1, 0} // 模版类型: 1表示发布,0表示请求
)

func GetEntityData(requestId, userToken string) (result models.EntityQueryResult, err error) {
	var byteArr []byte
	requestTemplateId, tmpErr := getRequestTemplateByRequest(requestId)
	if tmpErr != nil {
		return result, tmpErr
	}
	requestTemplateObj, getTemplateErr := GetSimpleRequestTemplate(requestTemplateId)
	if getTemplateErr != nil {
		err = getTemplateErr
		return
	}
	if requestTemplateObj.PackageName == "" || requestTemplateObj.EntityName == "" {
		err = fmt.Errorf("RequestTemplate packageName or entityName illegal ")
		return
	}
	url := fmt.Sprintf("%s/platform/v1/process/definitions/%s/root-entities", models.Config.Wecube.BaseUrl, requestTemplateObj.ProcDefId)
	byteArr, err = HttpGet(url, userToken)
	var responseObj models.WorkflowEntityQuery
	err = json.Unmarshal(byteArr, &responseObj)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
	} else {
		result.Status = responseObj.Status
		result.Message = responseObj.Message
		for _, v := range responseObj.Data {
			result.Data = append(result.Data, &models.EntityDataObj{Id: v.Id, DisplayName: v.DisplayName, PackageName: requestTemplateObj.PackageName, Entity: requestTemplateObj.EntityName})
		}
	}
	return
}

func ProcessDataPreview(requestTemplateId, entityDataId, userToken string) (result models.EntityTreeResult, err error) {
	requestTemplateObj, getTemplateErr := GetSimpleRequestTemplate(requestTemplateId)
	if getTemplateErr != nil {
		err = getTemplateErr
		return
	}
	if requestTemplateObj.ProcDefId == "" {
		err = fmt.Errorf("RequestTemplate proDefId illegal ")
		return
	}
	req, newReqErr := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/platform/v1/public/process/definitions/%s/preview/entities/%s", models.Config.Wecube.BaseUrl, requestTemplateObj.ProcDefId, entityDataId), nil)
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
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	err = json.Unmarshal(b, &result)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
	}
	return
}

// GetRequestCount 工作台请求统计
func GetRequestCount(user string, userRoles []string) (platformData models.PlatformData, err error) {
	platformData.Pending = strings.Join(GetPendingCount(userRoles), ";")
	platformData.HasProcessed = strings.Join(GetHasProcessedCount(user), ";")
	platformData.Submit = strings.Join(GetSubmitCount(user), ";")
	platformData.Draft = strings.Join(GetDraftCount(user), ";")
	platformData.Collect = strings.Join(GetCollectCount(user), ";")
	return
}

// GetPendingCount 统计待处理,包括:(1)待定版 (2)任务待审批 分配给本组、本人的所有请求,此处按照角色去统计
func GetPendingCount(userRoles []string) (resultArr []string) {
	userRolesFilterSql, userRolesFilterParam := dao.CreateListParams(userRoles, "")
	var requestQueryParam, taskQueryParam []interface{}
	var roleFilterList []string
	var requestSQL, taskSQL string
	roleFilterSql := "1=1"
	if len(userRoles) > 0 {
		for _, v := range userRoles {
			roleFilterList = append(roleFilterList, "report_role like '%,"+v+",%'")
		}
		roleFilterSql = strings.Join(roleFilterList, " or ")
	}
	for i := 0; i < len(templateTypeArr); i++ {
		requestSQL, requestQueryParam = pendingRequestSQL(templateTypeArr[i], userRolesFilterSql, userRolesFilterParam)
		taskSQL, taskQueryParam = pendingTaskSQL(templateTypeArr[i], userRolesFilterSql, userRolesFilterParam, roleFilterSql)
		resultArr = append(resultArr, strconv.Itoa(dao.QueryCount(requestSQL, requestQueryParam...)+dao.QueryCount(taskSQL, taskQueryParam...)))
	}
	return
}

func pendingRequestSQL(templateType int, userRolesFilterSql string, userRolesFilterParam []interface{}) (sql string, queryParam []interface{}) {
	sql = "select id from request where del_flag = 0 and status = 'Pending' and type = ? and request_template in (select id " +
		"from request_template where  id in (select request_template from request_template_role where role_type= 'MGMT' " +
		"and `role` in (" + userRolesFilterSql + "))) "
	queryParam = append([]interface{}{templateType}, userRolesFilterParam...)
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
func DataList(param *models.PlatformRequestParam, userRoles []string, userToken, user string) (pageInfo models.PageInfo, rowData []*models.PlatformDataObj, err error) {
	// 先拼接查询条件
	var templateType int
	var sql string
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
			sql, queryParam = pendingRequestSQL(templateType, userRolesFilterSql, userRolesFilterParam)
		} else if param.Type == 2 {
			sql, queryParam = pendingTaskSQL(templateType, userRolesFilterSql, userRolesFilterParam, roleFilterSql)
			pageInfo, rowData, err = getPlatData(models.PlatDataParam{Param: param.CommonRequestParam, QueryParam: queryParam, UserToken: userToken}, getPlatTaskSQL(where, sql), true)
			return
		}
	case "hasProcessed":
		if param.Type == 1 {
			sql, queryParam = hasProcessedRequestSQL(templateType, user)
		} else if param.Type == 2 {
			sql, queryParam = hasProcessedTaskSQL(templateType, user)
			pageInfo, rowData, err = getPlatData(models.PlatDataParam{Param: param.CommonRequestParam, QueryParam: queryParam, UserToken: userToken}, getPlatTaskSQL(where, sql), true)
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
	pageInfo, rowData, err = getPlatData(models.PlatDataParam{Param: param.CommonRequestParam, QueryParam: queryParam, UserToken: userToken}, getPlatRequestSQL(where, sql), true)
	return
}

// HistoryList 发布历史
func HistoryList(param *models.RequestHistoryParam, userRoles []string, userToken, user string) (pageInfo models.PageInfo, rowsData []*models.PlatformDataObj, err error) {
	var sql = "select id from request"
	var queryParam []interface{}
	where := transHistoryConditionToSQL(param)
	// 查看本组数据
	if param.Permission == "group" {
		userRolesFilterSql, userRolesFilterParam := dao.CreateListParams(userRoles, "")
		sql = "select id from request where  `role` in (" + userRolesFilterSql + ")"
		queryParam = append(queryParam, userRolesFilterParam...)
	}
	pageInfo, rowsData, err = getPlatData(models.PlatDataParam{Param: param.CommonRequestParam, QueryParam: queryParam, User: user, UserToken: userToken}, getPlatRequestSQL(where, sql), true)
	return
}

// Export 数据导出
func Export(w http.ResponseWriter, param *models.RequestHistoryParam, userToken, language, user string) (err error) {
	var rowsData []*models.PlatformDataObj
	var sql = "select id from request"
	var queryParam []interface{}
	where := transHistoryConditionToSQL(param)
	_, rowsData, err = getPlatData(models.PlatDataParam{Param: param.CommonRequestParam, QueryParam: queryParam, User: user, UserToken: userToken}, getPlatRequestSQL(where, sql), false)
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
		case string(models.Draft):
			return "草稿"
		case string(models.Pending):
			return "等待定版"
		case string(models.InProgress):
			return "执行中"
		case string(models.InProgressFaulted):
			return "节点报错"
		case string(models.Termination):
			return "手动终止"
		case string(models.Completed):
			return "成功"
		case string(models.InProgressTimeOuted):
			return "节点超时"
		case string(models.Faulted):
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
	if dataObject.Status == string(models.Draft) {
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
	if dataObject.Status == string(models.Completed) || dataObject.Status == string(models.Termination) || dataObject.Status == string(models.Faulted) {
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

func getPlatData(req models.PlatDataParam, newSQL string, page bool) (pageInfo models.PageInfo, rowsData []*models.PlatformDataObj, err error) {
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
		operatorObjTypeMap = getAllCoreProcess(req.UserToken)
		// 查询当前用户所有收藏模板记录
		collectMap, _ := QueryAllTemplateCollect(req.User)
		templateMap, _ := getAllRequestTemplate()
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
				platformDataObj.Progress, platformDataObj.CurNode = getCurNodeName(platformDataObj.ProcInstanceId, req.UserToken)
			}
			if strings.Contains(platformDataObj.Status, "InProgress") && platformDataObj.ProcInstanceId != "" {
				newStatus := getInstanceStatus(platformDataObj.ProcInstanceId, req.UserToken)
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
				result, _ := GetEntityData(platformDataObj.Id, req.UserToken)
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

func ListRequest(param *models.QueryRequestParam, userRoles []string, userToken, permission, operator string) (pageInfo models.PageInfo, rowData []*models.RequestTable, err error) {
	rowData = []*models.RequestTable{}
	if strings.ToLower(permission) == "mgmt" {
		permission = "MGMT"
	} else {
		permission = "USE"
	}
	filterSql, _, queryParam := dao.TransFiltersToSQL(param, &models.TransFiltersParam{IsStruct: true, StructObj: models.RequestTable{}, PrimaryKey: "id"})
	userRolesFilterSql, userRolesFilterParam := dao.CreateListParams(userRoles, "")
	baseSql := fmt.Sprintf("select id,name,form,request_template,proc_instance_id,proc_instance_key,reporter,handler,report_time,emergency,status,expire_time,expect_time,confirm_time,created_by,created_time,updated_by,updated_time,rollback_desc from request where del_flag=0 and (created_by=? or request_template in (select id from request_template where id in (select request_template from request_template_role where role_type=? and `role` in ("+userRolesFilterSql+")))) %s ", filterSql)
	queryParam = append(append([]interface{}{operator, permission}, userRolesFilterParam...), queryParam...)
	if param.Paging {
		pageInfo.StartIndex = param.Pageable.StartIndex
		pageInfo.PageSize = param.Pageable.PageSize
		pageInfo.TotalRows = dao.QueryCount(baseSql, queryParam...)
		pageSql, pageParam := dao.TransPageInfoToSQL(*param.Pageable)
		baseSql += pageSql
		queryParam = append(queryParam, pageParam...)
	}
	err = dao.X.SQL(baseSql, queryParam...).Find(&rowData)
	if len(rowData) > 0 {
		var requestTemplateTable []*models.RequestTemplateTable
		dao.X.SQL("select id,name,status,version from request_template").Find(&requestTemplateTable)
		rtMap := make(map[string]string)
		for _, v := range requestTemplateTable {
			if v.Status != "confirm" {
				rtMap[v.Id] = fmt.Sprintf("%s(beta)", v.Name)
			} else {
				rtMap[v.Id] = fmt.Sprintf("%s(%s)", v.Name, v.Version)
			}
		}
		rtRoleMap := getRequestTemplateMGMTRole()
		var actions []*dao.ExecAction
		for _, v := range rowData {
			v.RequestTemplateName = rtMap[v.RequestTemplate]
			if tmpRoles, b := rtRoleMap[v.RequestTemplate]; b {
				v.HandleRoles = tmpRoles
			} else {
				v.HandleRoles = []string{}
			}
			if strings.Contains(v.Status, "InProgress") && v.ProcInstanceId != "" {
				newStatus := getInstanceStatus(v.ProcInstanceId, userToken)
				if newStatus == "InternallyTerminated" {
					newStatus = "Termination"
				}
				if newStatus != "" && newStatus != v.Status {
					actions = append(actions, &dao.ExecAction{Sql: "update request set status=? where id=?", Param: []interface{}{newStatus, v.Id}})
					v.Status = newStatus
				}
			}
			if v.Status == "Completed" {
				v.CompletedTime = v.UpdatedTime
			}
		}
		if len(actions) > 0 {
			updateStatusErr := dao.Transaction(actions)
			if updateStatusErr != nil {
				log.Logger.Error("Try to update request status fail", log.Error(updateStatusErr))
			}
		}
	}
	return
}

func getRequestTemplateMGMTRole() (result map[string][]string) {
	result = make(map[string][]string)
	var requestTemplateRole []*models.RequestTemplateRoleTable
	dao.X.SQL("select * from request_template_role where role_type='MGMT' order by request_template").Find(&requestTemplateRole)
	if len(requestTemplateRole) == 0 {
		return result
	}
	tmpTemplate := requestTemplateRole[0].RequestTemplate
	tmpRoles := []string{}
	for _, v := range requestTemplateRole {
		if v.RequestTemplate != tmpTemplate {
			result[tmpTemplate] = tmpRoles
			tmpTemplate = v.RequestTemplate
			tmpRoles = []string{}
		}
		tmpRoles = append(tmpRoles, v.Role)
	}
	if len(tmpRoles) > 0 {
		tmpTemplate = requestTemplateRole[len(requestTemplateRole)-1].RequestTemplate
		result[tmpTemplate] = tmpRoles
	}
	return result
}

func calcExpireObj(param *models.ExpireObj) {
	if param.ReportTime == "" || param.ExpireTime == "" {
		return
	}
	reportT, _ := time.Parse(models.DateTimeFormat, param.ReportTime)
	expireT, _ := time.Parse(models.DateTimeFormat, param.ExpireTime)
	nowT, _ := time.Parse(models.DateTimeFormat, param.NowTime)
	max := expireT.Sub(reportT).Seconds()
	use := nowT.Sub(reportT).Seconds()
	param.Percent = (use / max) * 100
	param.TotalDay = max / 86400
	param.UseDay = use / 86400
	param.LeftDay = (max - use) / 86400
	param.Percent, _ = strconv.ParseFloat(fmt.Sprintf("%.0f", param.Percent), 64)
	param.TotalDay, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", param.TotalDay), 64)
	param.LeftDay, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", param.LeftDay), 64)
	param.UseDay, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", param.UseDay), 64)
	return
}

func getInstanceStatus(instanceId, userToken string) string {
	response, err := getProcessInstances(instanceId, userToken)
	if err != nil {
		return ""
	}
	if response.Data.Status != "InProgress" {
		return response.Data.Status
	}
	status := "InProgress"
	for _, v := range response.Data.TaskNodeInstances {
		if v.Status == "Faulted" {
			status = "InProgress(Faulted)"
			break
		}
		if v.Status == "Timeouted" {
			status = "InProgress(Timeouted)"
			break
		}
	}
	return status
}

func getCurNodeName(instanceId, userToken string) (progress int, curNode string) {
	var total int
	response, err := getProcessInstances(instanceId, userToken)
	if err != nil || len(response.Data.TaskNodeInstances) == 0 {
		return
	}
	// 统计完成进度 ,已完成/总数, 编号不为空
	for _, v := range response.Data.TaskNodeInstances {
		if v.OrderedNo != "" {
			if v.Status == "Completed" {
				progress++
			}
			total++
		}
	}
	progress = int(math.Floor(float64(progress)/float64(total)*100 + 0.5))
	switch response.Data.Status {
	case "Completed":
		curNode = CurNodeCompleted
		return
	case "InProgress":
		for _, v := range response.Data.TaskNodeInstances {
			if v.Status == "InProgress" || v.Status == "Timeouted" || v.Status == "Faulted" {
				curNode = v.NodeName
				return
			}
		}
	case "NotStarted":
		curNode = "NotStarted"
	case "Faulted":
		// 失败状态,筛选成功并且有序号的节点
		list := filterSuccessNode(response.Data.TaskNodeInstances)
		if len(list) == 0 {
			return
		}
		// 按 orderNo排序,将有 orderNo的节点按小到大排序,找到最后一个节点状态返回
		sort.Sort(models.QueryNodeSort(list))
		curNode = list[len(list)-1].NodeName
		return
	default:
		// 失败状态,显示具体执行失败的节点. filterNode 过滤orderNo为空大节点
		list := filterNode(response.Data.TaskNodeInstances)
		if len(list) == 0 {
			// 如果都没有序号,找一个NotStarted节点,找不到返回 Completed
			for _, v := range response.Data.TaskNodeInstances {
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

func filterNode(instances []*models.InstanceStatusQueryNode) []*models.InstanceStatusQueryNode {
	var list []*models.InstanceStatusQueryNode
	for _, node := range instances {
		if node.OrderedNo != "" {
			list = append(list, node)
		}
	}
	return list
}

func filterSuccessNode(instances []*models.InstanceStatusQueryNode) []*models.InstanceStatusQueryNode {
	var list []*models.InstanceStatusQueryNode
	for _, node := range instances {
		if node.OrderedNo != "" && node.Status == "Completed" {
			list = append(list, node)
		}
	}
	return list
}

// getProcessInstances 获取编排
func getProcessInstances(instanceId, userToken string) (response models.InstanceStatusQuery, err error) {
	var req *http.Request
	var resp *http.Response
	req, err = http.NewRequest(http.MethodGet, fmt.Sprintf("%s/platform/v1/process/instances/%s", models.Config.Wecube.BaseUrl, instanceId), nil)
	if err != nil {
		log.Logger.Error("GetInstanceStatus fail", log.String("msg", "new http request fail"), log.Error(err))
		return
	}
	req.Header.Set("Authorization", userToken)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Logger.Error("GetInstanceStatus fail", log.String("msg", "Try to do http request fail"), log.Error(err))
		return
	}
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	err = json.Unmarshal(b, &response)
	if err != nil {
		log.Logger.Error("GetInstanceStatus fail", log.String("msg", "Try to json unmarshal body fail"), log.Error(err))
		return
	}
	if response.Status != "OK" {
		log.Logger.Error("GetInstanceStatus fail", log.String("msg", response.Message))
		err = fmt.Errorf("GetInstanceStatus fail")
		return
	}
	return
}

func calcExpireTime(reportTime string, expireDay int) (expire string) {
	t, err := time.Parse(models.DateTimeFormat, reportTime)
	if err != nil {
		return
	}
	expire = t.Add(time.Duration(expireDay*24) * time.Hour).Format(models.DateTimeFormat)
	return
}

func RevokeRequest(requestId, user string) (err error) {
	var request models.RequestTable
	nowTime := time.Now().Format(models.DateTimeFormat)
	request, err = GetSimpleRequest(requestId)
	if err != nil {
		return
	}
	if request.Id == "" {
		err = fmt.Errorf("param requestId:%s not exist", requestId)
		return
	}
	if request.Status != "Pending" {
		err = exterror.New().RevokeRequestError
		return
	}
	if request.CreatedBy != user {
		err = fmt.Errorf("request not yours")
		return
	}
	_, err = dao.X.Exec("update request set status ='Draft',revoke_flag=1,updated_time=? where id=?", nowTime, request.Id)
	return
}

func GetRequestWithRoot(requestId string) (result models.RequestTable, err error) {
	result = models.RequestTable{}
	var requestTable []*models.RequestTable
	err = dao.X.SQL("select id,name,form,request_template,proc_instance_id,proc_instance_key,reporter,report_time,emergency,status,cache,expire_time,expect_time,confirm_time,handler from request where id=?", requestId).Find(&requestTable)
	if err != nil {
		return
	}
	if len(requestTable) == 0 {
		err = fmt.Errorf("Can not find any request with id:%s ", requestId)
		return
	}
	result = *requestTable[0]
	if result.Cache != "" {
		var cacheObj models.RequestPreDataDto
		err = json.Unmarshal([]byte(result.Cache), &cacheObj)
		if err != nil {
			err = fmt.Errorf("Try to json unmarshal cache data fail,%s ", err.Error())
			return
		}
		result.Cache = cacheObj.RootEntityId
	}
	result.AttachFiles = GetRequestAttachFileList(requestId)
	return
}

func GetRequest(requestId string) (result models.RequestTable, err error) {
	result = models.RequestTable{}
	var requestTable []*models.RequestTable
	err = dao.X.SQL("select id,name,form,request_template,proc_instance_id,proc_instance_key,reporter,report_time,emergency,status from request where id=?", requestId).Find(&requestTable)
	if err != nil {
		return
	}
	if len(requestTable) == 0 {
		err = fmt.Errorf("Can not find any request with id:%s ", requestId)
		return
	}
	result = *requestTable[0]
	return
}

func CreateRequest(param *models.RequestTable, operatorRoles []string, userToken string) error {
	requestTemplateObj, err := GetSimpleRequestTemplate(param.RequestTemplate)
	if err != nil {
		return err
	}
	var actions []*dao.ExecAction
	err = SyncProcDefId(requestTemplateObj.Id, requestTemplateObj.ProcDefId, requestTemplateObj.ProcDefName, "", userToken)
	if err != nil {
		return fmt.Errorf("Try to sync proDefId fail,%s ", err.Error())
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	formGuid := guid.CreateGuid()
	param.Id = newRequestId()
	formInsertAction := dao.ExecAction{Sql: "insert into form(id,name,description,form_template,created_time,created_by,updated_time,updated_by) value (?,?,?,?,?,?,?,?)"}
	formInsertAction.Param = []interface{}{formGuid, param.Name + models.SysTableIdConnector + "form", "", requestTemplateObj.FormTemplate,
		nowTime, param.CreatedBy, nowTime, param.CreatedBy}
	actions = append(actions, &formInsertAction)
	requestInsertAction := dao.ExecAction{Sql: "insert into request(id,name,form,request_template,reporter,emergency,report_role,status,expire_time," +
		"expect_time,handler,created_by,created_time,updated_by,updated_time,type,role) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
	requestInsertAction.Param = []interface{}{param.Id, param.Name, formGuid, param.RequestTemplate, param.CreatedBy, param.Emergency,
		strings.Join(operatorRoles, ","), "Draft", "", param.ExpectTime, requestTemplateObj.Handler, param.CreatedBy, nowTime,
		param.CreatedBy, nowTime, param.Type, param.Role}
	actions = append(actions, &requestInsertAction)
	return dao.TransactionWithoutForeignCheck(actions)
}

func UpdateRequest(param *models.RequestTable) error {
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err := dao.X.Exec("update request set name=?,expect_time=?,emergency=?,handler=?,updated_by=?,updated_time=? where id=?", param.Name, param.ExpectTime, param.Emergency, param.Handler, param.UpdatedBy, nowTime, param.Id)
	return err
}

func DeleteRequest(requestId, operator string) error {
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err := dao.X.Exec("update request set del_flag=1,updated_by=?,updated_time=? where id=?", operator, nowTime, requestId)
	return err
}

func SaveRequestCacheNew(requestId, operator, userToken string, param *models.RequestPreDataDto) error {
	err := ValidateRequestForm(param.Data, userToken)
	if err != nil {
		return err
	}
	var formItemNameQuery []*models.FormItemTemplateTable
	err = dao.X.SQL("select item_group,group_concat(name,',') as name from form_item_template where in_display_name='yes' and form_template in (select form_template from request_template where id in (select request_template from request where id=?)) group by item_group", requestId).Find(&formItemNameQuery)
	itemGroupNameMap := make(map[string][]string)
	for _, v := range formItemNameQuery {
		itemGroupNameMap[v.ItemGroup] = strings.Split(v.Name, ",")
	}
	for _, v := range param.Data {
		nameList := itemGroupNameMap[v.ItemGroup]
		for _, value := range v.Value {
			if value.Id == "" {
				value.PackageName = v.PackageName
				value.EntityName = v.Entity
				value.EntityDataOp = "create"
				value.Id = fmt.Sprintf("tmp%s%s", models.SysTableIdConnector, guid.CreateGuid())
				value.DisplayName = concatItemDisplayName(value.EntityData, nameList)
			}
		}
	}
	paramBytes, err := json.Marshal(param)
	if err != nil {
		return fmt.Errorf("Try to json marshal param fail,%s ", err.Error())
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	actions := UpdateRequestFormItem(requestId, param)
	actions = append(actions, &dao.ExecAction{Sql: "update request set cache=?,updated_by=?,updated_time=?,operator_obj=? where id=?", Param: []interface{}{string(paramBytes), operator, nowTime, param.EntityName, requestId}})
	return dao.Transaction(actions)
}

func SaveRequestCacheV2(requestId, operator, userToken string, param *models.RequestProDataV2Dto) error {
	err := ValidateRequestForm(param.Data, userToken)
	if err != nil {
		return err
	}
	var formItemNameQuery []*models.FormItemTemplateTable
	err = dao.X.SQL("select item_group,group_concat(name,',') as name from form_item_template where in_display_name='yes' and form_template in (select form_template from request_template where id in (select request_template from request where id=?)) group by item_group", requestId).Find(&formItemNameQuery)
	itemGroupNameMap := make(map[string][]string)
	for _, v := range formItemNameQuery {
		itemGroupNameMap[v.ItemGroup] = strings.Split(v.Name, ",")
	}
	for _, v := range param.Data {
		nameList := itemGroupNameMap[v.ItemGroup]
		for _, value := range v.Value {
			if value.Id == "" {
				value.PackageName = v.PackageName
				value.EntityName = v.Entity
				value.EntityDataOp = "create"
				value.Id = fmt.Sprintf("tmp%s%s", models.SysTableIdConnector, guid.CreateGuid())
				value.DisplayName = concatItemDisplayName(value.EntityData, nameList)
			}
		}
	}
	newParam := &models.RequestPreDataDto{
		RootEntityId: param.RootEntityId,
		Data:         param.Data,
	}
	paramBytes, err := json.Marshal(newParam)
	if err != nil {
		return fmt.Errorf("Try to json marshal param fail,%s ", err.Error())
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	actions := UpdateRequestFormItem(requestId, newParam)
	actions = append(actions, &dao.ExecAction{Sql: "update request set cache=?,updated_by=?,updated_time=?,name=?,description=?,expect_time=?,operator_obj=?" +
		" where id=?", Param: []interface{}{string(paramBytes), operator, nowTime, param.Name, param.Description, param.ExpectTime, param.EntityName, requestId}})
	return dao.Transaction(actions)
}

func SaveRequestBindCache(requestId, operator string, param *models.RequestCacheData) error {
	cacheBytes, _ := json.Marshal(param)
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err := dao.X.Exec("update request set bind_cache=?,updated_by=?,updated_time=? where id=?", string(cacheBytes), operator, nowTime, requestId)
	return err
}

func concatItemDisplayName(rowData map[string]interface{}, nameList []string) string {
	displayNameList := []string{}
	for _, v := range nameList {
		if v == "" {
			continue
		}
		if _, b := rowData[v]; b {
			tmpValue := fmt.Sprintf("%s", rowData[v])
			if rowData[v] == nil {
				tmpValue = ""
			}
			displayNameList = append(displayNameList, tmpValue)
		}
	}
	return strings.Join(displayNameList, "__")
}

func UpdateRequestFormItem(requestId string, param *models.RequestPreDataDto) []*dao.ExecAction {
	var actions []*dao.ExecAction
	requestObj, _ := GetRequest(requestId)
	actions = append(actions, &dao.ExecAction{Sql: "delete from form_item where form in (select form from request where id=?)", Param: []interface{}{requestId}})
	for _, v := range param.Data {
		for _, valueObj := range v.Value {
			tmpGuidList := guid.CreateGuidList(len(v.Title))
			for i, title := range v.Title {
				if title.Multiple == "Y" {
					if tmpV, b := valueObj.EntityData[title.Name]; b {
						tmpStringV := []string{}
						for _, interfaceV := range tmpV.([]interface{}) {
							tmpStringV = append(tmpStringV, fmt.Sprintf("%s", interfaceV))
						}
						actions = append(actions, &dao.ExecAction{Sql: "insert into form_item(id,form,form_item_template,name,value,item_group,row_data_id) value (?,?,?,?,?,?,?)", Param: []interface{}{tmpGuidList[i], requestObj.Form, title.Id, title.Name, strings.Join(tmpStringV, ","), title.ItemGroup, valueObj.Id}})
					}
				} else {
					actions = append(actions, &dao.ExecAction{Sql: "insert into form_item(id,form,form_item_template,name,value,item_group,row_data_id) value (?,?,?,?,?,?,?)", Param: []interface{}{tmpGuidList[i], requestObj.Form, title.Id, title.Name, valueObj.EntityData[title.Name], title.ItemGroup, valueObj.Id}})
				}
			}
		}
	}
	return actions
}

func GetRequestCache(requestId, cacheType string) (result interface{}, err error) {
	var requestTable []*models.RequestTable
	if cacheType == "data" {
		err = dao.X.SQL("select cache from request where id=?", requestId).Find(&requestTable)
		if err != nil {
			return
		}
		if len(requestTable) == 0 {
			err = fmt.Errorf("Can not find any request with id:%s ", requestId)
			return
		}
		if requestTable[0].Cache == "" {
			return models.RequestPreDataDto{Data: []*models.RequestPreDataTableObj{}}, nil
		}
		var dataCache models.RequestPreDataDto
		err = json.Unmarshal([]byte(requestTable[0].Cache), &dataCache)
		return dataCache, err
	} else {
		err = dao.X.SQL("select bind_cache from request where id=?", requestId).Find(&requestTable)
		if err != nil {
			return
		}
		if len(requestTable) == 0 {
			err = fmt.Errorf("Can not find any request with id:%s ", requestId)
			return
		}
		if requestTable[0].BindCache == "" {
			return models.RequestCacheData{TaskNodeBindInfos: []*models.RequestCacheTaskNodeBindObj{}}, nil
		}
		var bindCache models.RequestCacheData
		err = json.Unmarshal([]byte(requestTable[0].BindCache), &bindCache)
		return bindCache, err
	}
}

func getRequestTemplateByRequest(requestId string) (templateId string, err error) {
	var requestTable []*models.RequestTable
	err = dao.X.SQL("select request_template from request where id=?", requestId).Find(&requestTable)
	if err != nil {
		return
	}
	if len(requestTable) == 0 {
		err = fmt.Errorf("Can not find request with id:%s ", requestId)
		return
	}
	templateId = requestTable[0].RequestTemplate
	return
}

func GetRequestRootForm(requestId string) (result models.RequestTemplateFormStruct, err error) {
	result = models.RequestTemplateFormStruct{}
	requestTemplateId, tmpErr := getRequestTemplateByRequest(requestId)
	if tmpErr != nil {
		return result, tmpErr
	}
	requestTemplateObj, _ := GetSimpleRequestTemplate(requestTemplateId)
	result.Id = requestTemplateObj.Id
	result.Name = requestTemplateObj.Name
	result.PackageName = requestTemplateObj.PackageName
	result.EntityName = requestTemplateObj.EntityName
	result.ProcDefName = requestTemplateObj.ProcDefName
	result.ProcDefId = requestTemplateObj.ProcDefId
	result.ProcDefKey = requestTemplateObj.ProcDefKey
	var items []*models.FormItemTemplateTable
	dao.X.SQL("select * from form_item_template where form_template=?", requestTemplateObj.FormTemplate).Find(&items)
	result.FormItems = items
	return
}

func GetRequestPreData(requestId, entityDataId, userToken string) (result []*models.RequestPreDataTableObj, err error) {
	var requestTables []*models.RequestTable
	err = dao.X.SQL("select cache from request where id=?", requestId).Find(&requestTables)
	if err != nil {
		return
	}
	if len(requestTables) == 0 {
		return result, fmt.Errorf("Can not find requestId:%s ", requestId)
	}
	if requestTables[0].Cache != "" {
		var cacheObj models.RequestPreDataDto
		err = json.Unmarshal([]byte(requestTables[0].Cache), &cacheObj)
		if err != nil {
			return result, fmt.Errorf("Try to json unmarshal cache data fail,%s ", err.Error())
		}
		if cacheObj.RootEntityId == entityDataId {
			result = cacheObj.Data
			return
		}
	}
	result = []*models.RequestPreDataTableObj{}
	requestTemplateId, tmpErr := getRequestTemplateByRequest(requestId)
	if tmpErr != nil {
		return result, tmpErr
	}
	var items []*models.FormItemTemplateTable
	err = dao.X.SQL("select * from form_item_template where form_template in (select form_template from request_template where id=?) order by item_group,sort", requestTemplateId).Find(&items)
	if err != nil {
		return
	}
	if len(items) == 0 {
		err = exterror.New().GetRequestPreviewDataError
		return
	}
	result = getItemTemplateTitle(items)
	if entityDataId == "" {
		return
	}
	previewData, previewErr := ProcessDataPreview(requestTemplateId, entityDataId, userToken)
	if previewErr != nil {
		return result, previewErr
	}
	if len(previewData.Data.EntityTreeNodes) == 0 {
		return
	}
	for _, entity := range result {
		for _, tmpData := range previewData.Data.EntityTreeNodes {
			if tmpData.EntityName == entity.Entity {
				tmpValueData := make(map[string]interface{})
				for _, title := range entity.Title {
					tmpValueData[title.Name] = tmpData.EntityData[title.Name]
				}
				entity.Value = append(entity.Value, &models.EntityTreeObj{Id: tmpData.Id, PackageName: tmpData.PackageName, EntityName: tmpData.EntityName, DataId: tmpData.DataId, PreviousIds: tmpData.PreviousIds, SucceedingIds: tmpData.SucceedingIds, DisplayName: tmpData.DisplayName, FullDataId: tmpData.FullDataId, EntityData: tmpValueData})
			}
		}
	}
	return
}

func getItemTemplateTitle(items []*models.FormItemTemplateTable) []*models.RequestPreDataTableObj {
	result := []*models.RequestPreDataTableObj{}
	tmpPackageName := items[0].PackageName
	tmpEntity := items[0].Entity
	tmpItemGroup := items[0].ItemGroup
	tmpItemGroupName := items[0].ItemGroupName
	tmpRefEntity := []string{}
	tmpItems := []*models.FormItemTemplateTable{}
	existItemMap := make(map[string]int)
	for _, v := range items {
		tmpKey := fmt.Sprintf("%s__%s", v.ItemGroup, v.Name)
		if _, b := existItemMap[tmpKey]; b {
			continue
		} else {
			existItemMap[tmpKey] = 1
		}
		if v.ItemGroup != tmpItemGroup {
			if tmpItemGroup != "" {
				result = append(result, &models.RequestPreDataTableObj{Entity: tmpEntity, ItemGroup: tmpItemGroup, ItemGroupName: tmpItemGroupName, PackageName: tmpPackageName, Title: tmpItems, RefEntity: tmpRefEntity, Value: []*models.EntityTreeObj{}})
			}
			tmpItems = []*models.FormItemTemplateTable{}
			tmpEntity = v.Entity
			tmpPackageName = v.PackageName
			tmpItemGroup = v.ItemGroup
			tmpItemGroupName = v.ItemGroupName
			tmpRefEntity = []string{}
		} else {
			if tmpEntity == "" && v.Entity != "" {
				tmpEntity = v.Entity
				tmpPackageName = v.PackageName
			}
		}
		tmpItems = append(tmpItems, v)
		if v.RefEntity != "" {
			existFlag := false
			for _, vv := range tmpRefEntity {
				if vv == v.RefEntity {
					existFlag = true
					break
				}
			}
			if !existFlag {
				tmpRefEntity = append(tmpRefEntity, v.RefEntity)
			}
		}
	}
	if len(tmpItems) > 0 {
		if items[len(items)-1].Entity != "" {
			tmpEntity = items[len(items)-1].Entity
		}
		if items[len(items)-1].PackageName != "" {
			tmpPackageName = items[len(items)-1].PackageName
		}
		tmpItemGroup = items[len(items)-1].ItemGroup
		tmpItemGroupName = items[len(items)-1].ItemGroupName
		result = append(result, &models.RequestPreDataTableObj{Entity: tmpEntity, ItemGroup: tmpItemGroup, ItemGroupName: tmpItemGroupName, PackageName: tmpPackageName, Title: tmpItems, RefEntity: tmpRefEntity, Value: []*models.EntityTreeObj{}})
	}
	// sort result by dependence
	result = sortRequestEntity(result)
	//var err error
	//result, err = getCMDBSelectList(result, models.CoreToken.GetCoreToken())
	//if err != nil {
	//	log.Logger.Error("Try to get selectList fail", log.Error(err))
	//}
	for _, v := range result {
		for _, vv := range v.Title {
			if vv.SelectList == nil {
				vv.SelectList = []*models.EntityDataObj{}
			}
		}
	}
	return result
}

func sortRequestEntity(param []*models.RequestPreDataTableObj) models.RequestPreDataSort {
	var result models.RequestPreDataSort
	entityMap := make(map[string][]string)
	entityNumMap := make(map[string]int)
	for _, v := range param {
		for _, vv := range v.RefEntity {
			if vv == v.Entity {
				continue
			}
			if _, b := entityMap[vv]; b {
				entityMap[vv] = append(entityMap[vv], v.Entity)
			} else {
				entityMap[vv] = []string{v.Entity}
			}
		}
	}
	countNum := 0
	levelNum := 1
	for _, v := range param {
		if v.Entity == "" {
			v.SortLevel = levelNum
			entityNumMap[v.Entity] = v.SortLevel
			countNum = countNum + 1
			continue
		}
		if _, b := entityMap[v.Entity]; !b {
			v.SortLevel = levelNum
			entityNumMap[v.Entity] = v.SortLevel
			countNum = countNum + 1
		}
	}
	for countNum < len(param) {
		levelNum = levelNum + 1
		for k, v := range entityMap {
			ready := true
			for _, vv := range v {
				if _, b := entityNumMap[vv]; !b {
					ready = false
					break
				} else {
					if entityNumMap[vv] == levelNum {
						ready = false
						break
					}
				}
			}
			if ready {
				for _, vv := range param {
					if vv.Entity == k {
						vv.SortLevel = levelNum
						entityNumMap[vv.Entity] = vv.SortLevel
						countNum = countNum + 1
					}
				}
				delete(entityMap, k)
			}
		}
	}
	for _, v := range param {
		result = append(result, v)
	}
	sort.Sort(result)
	return result
}

func StartRequest(requestId, operator, userToken string, cacheData models.RequestCacheData) (result models.StartInstanceResultData, err error) {
	var requestTemplateTable []*models.RequestTemplateTable
	dao.X.SQL("select * from request_template where id in (select request_template from request where id=?)", requestId).Find(&requestTemplateTable)
	if len(requestTemplateTable) == 0 {
		return result, fmt.Errorf("Can not find requestTemplate with request:%s ", requestId)
	}
	cacheData.ProcDefId = requestTemplateTable[0].ProcDefId
	cacheData.ProcDefKey = requestTemplateTable[0].ProcDefKey
	entityDepMap, tmpErr := AppendUselessEntity(requestTemplateTable[0].Id, userToken, &cacheData)
	if tmpErr != nil {
		return result, fmt.Errorf("Try to append useless entity fail,%s ", tmpErr.Error())
	}
	fillBindingWithRequestData(requestId, userToken, &cacheData, entityDepMap)
	cacheBytes, _ := json.Marshal(cacheData)
	log.Logger.Info("cacheByte", log.String("cacheBytes", string(cacheBytes)))
	startParam := BuildRequestProcessData(cacheData)
	startParamBytes, tmpErr := json.Marshal(startParam)
	if tmpErr != nil {
		err = fmt.Errorf("Json marshal cache data fail,%s ", tmpErr.Error())
		return
	}
	log.Logger.Info("Start request", log.String("param", string(startParamBytes)))
	req, newReqErr := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/platform/v1/public/process/instances", models.Config.Wecube.BaseUrl), bytes.NewReader(startParamBytes))
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
	var respResult models.StartInstanceResponse
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	err = json.Unmarshal(b, &respResult)
	if err != nil {
		err = fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
		return
	}
	if respResult.Status != "OK" {
		err = fmt.Errorf("Start instance fail,%s ", respResult.Message)
		return
	}
	result = respResult.Data
	nowTime := time.Now().Format(models.DateTimeFormat)
	expireTime := calcExpireTime(nowTime, requestTemplateTable[0].ExpireDay)
	_, err = dao.X.Exec("update request set handler=?,proc_instance_id=?,proc_instance_key=?,confirm_time=?,expire_time=?,status=?,bind_cache=?,updated_by=?,updated_time=? where id=?", operator, strconv.Itoa(result.Id), result.ProcInstKey, nowTime, expireTime, respResult.Data.Status, string(cacheBytes), operator, nowTime, requestId)
	return
}

func UpdateRequestStatus(requestId, status, operator, userToken, description string) error {
	var err error
	var request models.RequestTable
	nowTime := time.Now().Format(models.DateTimeFormat)
	if status == "Pending" {
		bindData, bindErr := GetRequestPreBindData(requestId, userToken)
		if bindErr != nil {
			return fmt.Errorf("Try to build bind data fail,%s ", bindErr.Error())
		}
		bindCacheBytes, _ := json.Marshal(bindData)
		bindCache := string(bindCacheBytes)
		_, err = dao.X.Exec("update request set status=?,reporter=?,report_time=?,bind_cache=?,updated_by=?,updated_time=?,rollback_desc=null,revoke_flag=0 where id=?", status, operator, nowTime, bindCache, operator, nowTime, requestId)
		if err == nil {
			notifyRoleMail(requestId)
		}
	} else if status == "Draft" {
		request, err = GetSimpleRequest(requestId)
		if request.Handler != operator {
			err = exterror.New().UpdateRequestHandlerStatusError
			return err
		}
		_, err = dao.X.Exec("update request set status=?,rollback_desc=?,updated_by=?,handler=?,updated_time=?,confirm_time=? where id=?", status, description, operator, operator, nowTime, nowTime, requestId)
	} else {
		_, err = dao.X.Exec("update request set status=?,updated_by=?,updated_time=? where id=?", status, operator, nowTime, requestId)
	}
	return err
}

func fillBindingWithRequestData(requestId, userToken string, cacheData *models.RequestCacheData, existDepMap map[string][]string) {
	var items []*models.FormItemTemplateTable
	dao.X.SQL("select * from form_item_template where form_template in (select form_template from request_template where id in (select request_template from request where id=?)) order by entity,sort", requestId).Find(&items)
	itemMap := make(map[string][]string)
	for _, item := range items {
		if item.RefEntity == "" {
			continue
		}
		if nowRefList, b := itemMap[item.Entity]; b {
			existFlag := false
			for _, v := range nowRefList {
				if v == item.RefEntity {
					existFlag = true
				}
			}
			if !existFlag {
				itemMap[item.Entity] = append(itemMap[item.Entity], item.RefEntity)
			}
		} else {
			itemMap[item.Entity] = []string{item.RefEntity}
		}
	}
	// itemMap -> entity:[refEntity]
	for k, v := range itemMap {
		log.Logger.Info("itemMap", log.String("key", k), log.StringList("value", v))
	}
	entityNewMap := make(map[string][]string)
	for k, v := range existDepMap {
		log.Logger.Info("existDepMap", log.String("k", k), log.StringList("v", v))
		if k != "" && len(v) > 0 {
			entityNewMap[k] = v
		}
	}
	entityOidMap := make(map[string]int)
	dataIdOidMap := make(map[string]string)
	matchEntityRoot(requestId, userToken, cacheData)
	if cacheData.RootEntityValue.EntityDataOp != "create" {
		dataIdOidMap[cacheData.RootEntityValue.EntityDataId] = cacheData.RootEntityValue.Oid
	}
	for _, taskNode := range cacheData.TaskNodeBindInfos {
		for _, entityValue := range taskNode.BoundEntityValues {
			if entityValue.EntityDataId != "" {
				dataIdOidMap[entityValue.EntityDataId] = entityValue.Oid
			}
			entityOidMap[entityValue.Oid] = 1
			if _, b := entityNewMap[entityValue.Oid]; b {
				continue
			}
			if entityRefs, b := itemMap[entityValue.EntityName]; b {
				findEntityRefByItemRef(entityValue, entityRefs, entityNewMap, dataIdOidMap)
			}
		}
	}
	for k, v := range dataIdOidMap {
		log.Logger.Debug("dataIdOidMap", log.String("k", k), log.String("v", v))
	}
	for k, v := range entityNewMap {
		tmpRefOidList := []string{}
		for _, vv := range v {
			if oid, b := dataIdOidMap[vv]; b {
				tmpRefOidList = append(tmpRefOidList, oid)
			} else {
				tmpRefOidList = append(tmpRefOidList, vv)
			}
		}
		entityNewMap[k] = tmpRefOidList
		log.Logger.Info("entityNewMap", log.String("key", k), log.StringList("value", tmpRefOidList))
	}
	if len(entityNewMap) > 0 {
		rebuildEntityRefOids(&cacheData.RootEntityValue, entityNewMap, entityOidMap)
		for _, taskNode := range cacheData.TaskNodeBindInfos {
			for _, entityValue := range taskNode.BoundEntityValues {
				rebuildEntityRefOids(entityValue, entityNewMap, entityOidMap)
			}
		}
	}
}

func rebuildEntityRefOids(entityValue *models.RequestCacheEntityValue, entityNewMap map[string][]string, entityOidMap map[string]int) {
	if refOids, b := entityNewMap[entityValue.Oid]; b {
		entityValue.PreviousOids = append(entityValue.PreviousOids, refOids...)
	}
	for tmpOid, refOids := range entityNewMap {
		inRefFlag := false
		for _, refOid := range refOids {
			if entityValue.Oid == refOid {
				inRefFlag = true
			}
		}
		if inRefFlag {
			entityValue.SucceedingOids = append(entityValue.SucceedingOids, tmpOid)
		}
	}
	entityValue.PreviousOids = listToSet(entityValue.PreviousOids, entityOidMap)
	entityValue.SucceedingOids = listToSet(entityValue.SucceedingOids, entityOidMap)
}

func findEntityRefByItemRef(entityValue *models.RequestCacheEntityValue, entityRefs []string, entityNewMap map[string][]string, dataIdOidMap map[string]string) {
	if entityValue.EntityDataOp == "create" {
		log.Logger.Debug("findEntityRefByItemRef create", log.String("oid", entityValue.Oid))
		tmpRefOidList := []string{}
		for _, attrValueObj := range entityValue.AttrValues {
			tmpAttrEntity := getEntityNameFromAttrDefId(attrValueObj.AttrDefId, attrValueObj.AttrName)
			for _, entityRef := range entityRefs {
				if tmpAttrEntity == entityRef && attrValueObj.DataType == "ref" {
					tmpV := fmt.Sprintf("%s", attrValueObj.DataValue)
					if strings.Contains(tmpV, ",") {
						tmpRefOidList = append(tmpRefOidList, strings.Split(tmpV, ",")...)
					} else {
						tmpRefOidList = append(tmpRefOidList, tmpV)
					}
				}
			}
		}
		entityNewMap[entityValue.Oid] = tmpRefOidList
	} else {
		log.Logger.Debug("findEntityRefByItemRef exist", log.String("oid", entityValue.Oid), log.String("EntityDataId", entityValue.EntityDataId))
		dataIdOidMap[entityValue.EntityDataId] = entityValue.Oid
		tmpRefOidList := []string{}
		for _, attrValueObj := range entityValue.AttrValues {
			tmpAttrEntity := getEntityNameFromAttrDefId(attrValueObj.AttrDefId, attrValueObj.AttrName)
			for _, entityRef := range entityRefs {
				if tmpAttrEntity == entityRef && attrValueObj.DataType == "ref" {
					valueString := fmt.Sprintf("%s", attrValueObj.DataValue)
					log.Logger.Debug("findEntityRefByItemRef ref", log.String("oid", entityValue.Oid), log.String("valueString", valueString))
					if strings.Contains(valueString, ",") {
						for _, tmpV := range strings.Split(valueString, ",") {
							if strings.HasPrefix(tmpV, "tmp") {
								tmpRefOidList = append(tmpRefOidList, tmpV)
							}
						}
					} else {
						if strings.HasPrefix(valueString, "tmp") {
							tmpRefOidList = append(tmpRefOidList, valueString)
						}
					}
				}
			}
		}
		entityNewMap[entityValue.Oid] = tmpRefOidList
	}
}

func getEntityNameFromAttrDefId(attrDefId, attrName string) string {
	stringSplit := strings.Split(attrDefId, ":")
	if len(stringSplit) == 3 {
		return stringSplit[2]
	}
	return attrName
}

func matchEntityRoot(requestId, userToken string, cacheData *models.RequestCacheData) {
	for _, taskNode := range cacheData.TaskNodeBindInfos {
		existFlag := false
		for _, entityValue := range taskNode.BoundEntityValues {
			if entityValue.EntityDataId == cacheData.RootEntityValue.Oid {
				existFlag = true
				cacheData.RootEntityValue.Oid = entityValue.Oid
				cacheData.RootEntityValue.EntityName = entityValue.EntityName
				cacheData.RootEntityValue.EntityDataOp = entityValue.EntityDataOp
				cacheData.RootEntityValue.AttrValues = entityValue.AttrValues
				cacheData.RootEntityValue.PackageName = entityValue.PackageName
				cacheData.RootEntityValue.EntityDisplayName = entityValue.EntityDisplayName
				cacheData.RootEntityValue.BindFlag = entityValue.BindFlag
				cacheData.RootEntityValue.EntityDataId = entityValue.EntityDataId
				cacheData.RootEntityValue.EntityDataState = entityValue.EntityDataState
				cacheData.RootEntityValue.EntityDefId = entityValue.EntityDefId
				cacheData.RootEntityValue.FullEntityDataId = entityValue.FullEntityDataId
				cacheData.RootEntityValue.PreviousOids = entityValue.PreviousOids
				cacheData.RootEntityValue.SucceedingOids = entityValue.SucceedingOids
				break
			}
		}
		if existFlag {
			break
		}
	}
	if cacheData.RootEntityValue.Oid != "" && cacheData.RootEntityValue.EntityName == "" {
		entityQueryResult, entityQueryErr := GetEntityData(requestId, userToken)
		if entityQueryErr != nil {
			log.Logger.Error("Try to fill root entity data fail", log.Error(entityQueryErr))
		} else {
			for _, v := range entityQueryResult.Data {
				if cacheData.RootEntityValue.Oid == v.Id {
					cacheData.RootEntityValue.PackageName = v.PackageName
					cacheData.RootEntityValue.EntityName = v.Entity
					cacheData.RootEntityValue.BindFlag = "N"
					cacheData.RootEntityValue.EntityDataId = v.Id
					cacheData.RootEntityValue.Oid = fmt.Sprintf("%s:%s:%s", v.PackageName, v.Entity, v.Id)
					cacheData.RootEntityValue.EntityDisplayName = v.DisplayName
					cacheData.RootEntityValue.Processed = false
					cacheData.RootEntityValue.SucceedingOids = []string{}
					cacheData.RootEntityValue.PreviousOids = []string{}
					break
				}
			}
		}
	}
}

func listToSet(input []string, itemMap map[string]int) []string {
	result := []string{}
	tmpMap := make(map[string]int)
	for _, v := range input {
		if v == "" {
			continue
		}
		if _, b := tmpMap[v]; !b {
			if _, bb := itemMap[v]; bb {
				result = append(result, v)
				tmpMap[v] = 1
			}
		}
	}
	return result
}

func RequestTermination(requestId, operator, userToken string) error {
	requestObj, err := GetRequest(requestId)
	if err != nil {
		return err
	}
	param := models.TerminateInstanceParam{ProcInstId: requestObj.ProcInstanceId, ProcInstKey: requestObj.ProcInstanceKey}
	paramBytes, tmpErr := json.Marshal(param)
	if tmpErr != nil {
		return fmt.Errorf("Json marshal param data fail,%s ", tmpErr.Error())
	}
	req, newReqErr := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/platform/v1/public/process/instances/%s/terminations", models.Config.Wecube.BaseUrl, requestObj.ProcInstanceId), bytes.NewReader(paramBytes))
	if newReqErr != nil {
		return fmt.Errorf("Try to new http request fail,%s ", newReqErr.Error())
	}
	req.Header.Set("Authorization", userToken)
	req.Header.Set("Content-Type", "application/json")
	resp, respErr := http.DefaultClient.Do(req)
	if respErr != nil {
		return fmt.Errorf("Try to do http request fail,%s ", respErr.Error())
	}
	var respResult models.StartInstanceResponse
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	err = json.Unmarshal(b, &respResult)
	if err != nil {
		return fmt.Errorf("Try to json unmarshal response body fail,%s ", err.Error())
	}
	if respResult.Status != "OK" {
		return fmt.Errorf("Terminate instance fail,%s ", respResult.Message)
	}
	nowTime := time.Now().Format(models.DateTimeFormat)
	_, err = dao.X.Exec("update request set status='Termination',updated_by=?,updated_time=? where id=?", operator, nowTime, requestId)
	return err
}

func GetCmdbReferenceData(attrId, userToken string, param models.QueryRequestParam) (result []byte, statusCode int, err error) {
	paramBytes, tmpErr := json.Marshal(param)
	if tmpErr != nil {
		err = fmt.Errorf("Json marshal param data fail,%s ", tmpErr.Error())
		return
	}
	req, newReqErr := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/wecmdb/api/v1/ci-data/reference-data/query/%s", models.Config.Wecube.BaseUrl, attrId), bytes.NewReader(paramBytes))
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
	statusCode = resp.StatusCode
	result, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return
}

func GetRequestPreBindData(requestId, userToken string) (result models.RequestCacheData, err error) {
	var requestTable []*models.RequestTable
	err = dao.X.SQL("select * from request where id=?", requestId).Find(&requestTable)
	if err != nil {
		return result, fmt.Errorf("Try to query request fail,%s ", err.Error())
	}
	if len(requestTable) == 0 {
		return result, fmt.Errorf("Can not find request with id:%s ", requestId)
	}
	if requestTable[0].Cache == "" {
		return result, fmt.Errorf("Can not find request cache data with id:%s ", requestId)
	}
	processNodes, processErr := GetProcessNodesByProc(models.RequestTemplateTable{Id: requestTable[0].RequestTemplate}, userToken, "bind")
	if processErr != nil {
		return result, processErr
	}
	entityDefIdMap := make(map[string]string)
	entityBindMap := make(map[string][]string)
	for _, v := range processNodes {
		tmpBoundEntities := []string{}
		for _, vv := range v.BoundEntities {
			if vv == nil {
				continue
			}
			entityDefIdMap[fmt.Sprintf("%s:%s", vv.PackageName, vv.Name)] = vv.Id
			tmpBoundEntities = append(tmpBoundEntities, fmt.Sprintf("%s:%s", vv.PackageName, vv.Name))
		}
		if v.TaskCategory != "SUTN" {
			entityBindMap[v.NodeDefId] = tmpBoundEntities
		}
	}
	var dataCache models.RequestPreDataDto
	json.Unmarshal([]byte(requestTable[0].Cache), &dataCache)
	entityOidMap := make(map[string][]string)
	entityValueMap := make(map[string]*models.RequestCacheEntityValue)
	for i, v := range dataCache.Data {
		for ii, vv := range v.Value {
			if vv.EntityName == "" {
				continue
			}
			tmpValueObj := models.RequestCacheEntityValue{Oid: vv.Id, BindFlag: "Y", EntityName: vv.EntityName, EntityDisplayName: vv.DisplayName, EntityDataOp: vv.EntityDataOp, EntityDataId: vv.DataId, FullEntityDataId: vv.FullDataId, PackageName: vv.PackageName, PreviousOids: vv.PreviousIds, SucceedingOids: vv.SucceedingIds}
			tmpEntityName := fmt.Sprintf("%s:%s", vv.PackageName, vv.EntityName)
			if _, b := entityDefIdMap[tmpEntityName]; b {
				tmpValueObj.EntityDefId = entityDefIdMap[tmpEntityName]
			}
			tmpValueObj.AttrValues = buildEntityValueAttrData(v.Title, vv.EntityData)
			if dataCache.RootEntityId == "" && i == 0 && ii == 0 {
				result.RootEntityValue = tmpValueObj
			}
			entityValueMap[vv.Id] = &tmpValueObj
			if _, b := entityOidMap[v.ItemGroup]; b {
				entityOidMap[v.ItemGroup] = append(entityOidMap[v.ItemGroup], vv.Id)
			} else {
				entityOidMap[v.ItemGroup] = []string{vv.Id}
			}
		}
	}
	var entityNodeBind []*models.EntityNodeBindQueryObj
	dao.X.SQL("select distinct t1.node_def_id,t2.item_group from task_template t1 left join form_item_template t2 on t1.form_template=t2.form_template where t1.request_template=?", requestTable[0].RequestTemplate).Find(&entityNodeBind)
	for _, v := range entityNodeBind {
		if _, b := entityBindMap[v.NodeDefId]; b {
			entityBindMap[v.NodeDefId] = append(entityBindMap[v.NodeDefId], v.ItemGroup)
		} else {
			entityBindMap[v.NodeDefId] = []string{v.ItemGroup}
		}
	}
	for _, v := range processNodes {
		if v.TaskCategory == "" {
			continue
		}
		tmpNodeBindInfo := models.RequestCacheTaskNodeBindObj{NodeId: v.NodeId, NodeDefId: v.NodeDefId, BoundEntityValues: []*models.RequestCacheEntityValue{}}
		if entities, b := entityBindMap[v.NodeDefId]; b {
			for _, entity := range entities {
				if oids, entityExist := entityOidMap[entity]; entityExist {
					for _, oid := range oids {
						if _, oidExist := entityValueMap[oid]; oidExist {
							tmpNodeBindInfo.BoundEntityValues = append(tmpNodeBindInfo.BoundEntityValues, entityValueMap[oid])
						}
					}
				}
			}
		}
		result.TaskNodeBindInfos = append(result.TaskNodeBindInfos, &tmpNodeBindInfo)
	}
	return
}

func buildEntityValueAttrData(titles []*models.FormItemTemplateTable, entityData map[string]interface{}) (result []*models.RequestCacheEntityAttrValue) {
	result = []*models.RequestCacheEntityAttrValue{}
	titleMap := make(map[string]*models.FormItemTemplateTable)
	for _, v := range titles {
		titleMap[v.Name] = v
	}
	for k, v := range entityData {
		if vv, b := titleMap[k]; b {
			if vv.Multiple == "Y" {
				tmpV := []string{}
				for _, interfaceV := range v.([]interface{}) {
					tmpV = append(tmpV, fmt.Sprintf("%s", interfaceV))
				}
				result = append(result, &models.RequestCacheEntityAttrValue{AttrDefId: vv.AttrDefId, AttrName: k, DataType: vv.AttrDefDataType, DataValue: strings.Join(tmpV, ",")})
			} else {
				result = append(result, &models.RequestCacheEntityAttrValue{AttrDefId: vv.AttrDefId, AttrName: k, DataType: vv.AttrDefDataType, DataValue: v})
			}
		}
	}
	return
}

func GetRequestTaskList(requestId string) (result models.TaskQueryResult, err error) {
	var taskTable []*models.TaskTable
	err = dao.X.SQL("select id from task where request=? order by created_time desc", requestId).Find(&taskTable)
	if err != nil {
		return
	}
	if len(taskTable) > 0 {
		result, err = GetTask(taskTable[0].Id)
		return
	}
	// get request
	var requests []*models.RequestTable
	dao.X.SQL("select * from request where id=?", requestId).Find(&requests)
	if len(requests) == 0 {
		return result, fmt.Errorf("Can not find request with id:%s ", requestId)
	}
	var requestCache models.RequestPreDataDto
	err = json.Unmarshal([]byte(requests[0].Cache), &requestCache)
	if err != nil {
		return result, fmt.Errorf("Try to json unmarshal request cache fail,%s ", err.Error())
	}
	var requestTemplateTable []*models.RequestTemplateTable
	dao.X.SQL("select * from request_template where id in (select request_template from request where id=?)", requestId).Find(&requestTemplateTable)
	requestQuery := models.TaskQueryObj{RequestId: requestId, RequestName: requests[0].Name, Reporter: requests[0].Reporter, ReportTime: requests[0].ReportTime, Comment: requests[0].Result, Editable: false}
	requestQuery.FormData = requestCache.Data
	requestQuery.AttachFiles = GetRequestAttachFileList(requestId)
	requestQuery.ExpireTime = requests[0].ExpireTime
	requestQuery.ExpectTime = requests[0].ExpectTime
	requestQuery.ProcInstanceId = requests[0].ProcInstanceId
	if len(requestTemplateTable) > 0 {
		requestQuery.RequestTemplate = requestTemplateTable[0].Name
	}
	result.Data = []*models.TaskQueryObj{&requestQuery}
	result.TimeStep, err = getRequestTimeStep(requests[0].RequestTemplate)
	if len(result.TimeStep) > 0 {
		result.TimeStep[0].Active = true
	}
	return
}

func GetRequestTaskListV2(requestId string) (taskQueryList []*models.TaskQueryObj, err error) {
	var taskTable []*models.TaskTable
	err = dao.X.SQL("select id from task where request=? order by created_time desc", requestId).Find(&taskTable)
	if err != nil {
		return
	}
	if len(taskTable) > 0 {
		taskQueryList, err = GetTaskV2(taskTable[0].Id)
		return
	}
	// get request
	var requests []*models.RequestTable
	dao.X.SQL("select * from request where id=?", requestId).Find(&requests)
	if len(requests) == 0 {
		err = fmt.Errorf("Can not find request with id:%s ", requestId)
		return
	}
	requestQuery := models.TaskQueryObj{RequestId: requestId, RequestName: requests[0].Name, Reporter: requests[0].Reporter, ReportTime: requests[0].ReportTime, Comment: requests[0].Result, Editable: false, RollbackDesc: requests[0].RollbackDesc}
	if requests[0].Cache != "" {
		var requestCache models.RequestPreDataDto
		err = json.Unmarshal([]byte(requests[0].Cache), &requestCache)
		if err != nil {
			return
		}
		requestQuery.FormData = requestCache.Data
	}
	requestQuery.AttachFiles = GetRequestAttachFileList(requestId)
	requestQuery.ExpireTime = requests[0].ExpireTime
	requestQuery.ExpectTime = requests[0].ExpectTime
	requestQuery.ProcInstanceId = requests[0].ProcInstanceId
	requestQuery.CreatedTime = requests[0].CreatedTime
	requestQuery.HandleTime = requests[0].ReportTime
	requestQuery.HandleRoleName = requests[0].Role
	requestQuery.Handler = requests[0].CreatedBy
	taskQueryList = append(taskQueryList, []*models.TaskQueryObj{&requestQuery, getPendingRequestData(requests[0])}...)
	return
}

func getPendingRequestData(request *models.RequestTable) *models.TaskQueryObj {
	var role string
	// 请求在定版状态,从模板角色表中读取
	rtRoleMap := getRequestTemplateMGMTRole()
	roles := rtRoleMap[request.RequestTemplate]
	if len(roles) > 0 {
		role = roles[0]
	}
	taskQueryObj := &models.TaskQueryObj{
		RequestId:      request.Id,
		RequestName:    request.Name,
		Editable:       false,
		Status:         "",
		ExpireTime:     request.ExpireTime,
		ExpectTime:     request.ExpectTime,
		Handler:        request.Handler,
		HandleTime:     request.UpdatedTime,
		FormData:       nil,
		IsHistory:      false,
		HandleRoleName: role,
		CreatedTime:    request.ReportTime,
		RollbackDesc:   request.RollbackDesc,
	}
	if request.Status != "Draft" && request.Status != "Pending" {
		taskQueryObj.HandleTime = request.UpdatedTime
	}
	return taskQueryObj
}

func GetRequestDetailV2(requestId, userToken string) (result models.RequestDetail, err error) {
	// get request
	var requests []*models.RequestTable
	var taskQueryList []*models.TaskQueryObj
	var actions []*dao.ExecAction
	dao.X.SQL("select * from request where id=?", requestId).Find(&requests)
	if len(requests) == 0 {
		return result, fmt.Errorf("Can not find request with id:%s ", requestId)
	}
	if strings.Contains(requests[0].Status, "InProgress") && requests[0].ProcInstanceId != "" {
		newStatus := getInstanceStatus(requests[0].ProcInstanceId, userToken)
		if newStatus == "InternallyTerminated" {
			newStatus = "Termination"
		}
		if newStatus != "" && newStatus != requests[0].Status {
			actions = append(actions, &dao.ExecAction{Sql: "update request set status=?,updated_time=? where id=?",
				Param: []interface{}{newStatus, time.Now().Format(models.DateTimeFormat), requests[0].Id}})
			requests[0].Status = newStatus
		}
		if len(actions) > 0 {
			updateRequestErr := dao.Transaction(actions)
			if updateRequestErr != nil {
				log.Logger.Error("Try to update request status fail", log.Error(updateRequestErr))
			}
		}
	}
	result.Request = getRequestForm(requests[0], userToken)
	taskQueryList, err = GetRequestTaskListV2(requestId)
	if err != nil {
		return
	}
	result.Data = taskQueryList
	return
}

func getCMDBSelectList(input []*models.RequestPreDataTableObj, userToken string) (output []*models.RequestPreDataTableObj, err error) {
	ciAttrMap := make(map[string][]string)
	ciAttrSelectMap := make(map[string][]*models.EntityDataObj)
	for _, v := range input {
		if v.Entity == "" {
			continue
		}
		for _, vv := range v.Title {
			if vv.ElementType == "select" && vv.RefEntity == "" {
				if _, b := ciAttrMap[v.Entity]; b {
					ciAttrMap[v.Entity] = append(ciAttrMap[v.Entity], vv.Name)
				} else {
					ciAttrMap[v.Entity] = []string{vv.Name}
				}
			}
			vv.SelectList = []*models.EntityDataObj{}
		}
	}
	output = input
	if len(ciAttrMap) <= 0 {
		return
	}
	for k, v := range ciAttrMap {
		tmpV, tmpErr := getCMDBAttributeSelectList(k, userToken, v)
		if tmpErr != nil {
			err = tmpErr
			break
		}
		for kk, vv := range tmpV {
			ciAttrSelectMap[kk] = vv
		}
	}
	if err != nil {
		return
	}
	for _, v := range output {
		if v.Entity == "" {
			continue
		}
		for _, vv := range v.Title {
			tmpKey := v.Entity + "_" + vv.Name
			if tmpSelectList, b := ciAttrSelectMap[tmpKey]; b {
				vv.SelectList = tmpSelectList
			}
		}
	}
	return
}

func getCMDBAttributes(entity, userToken string) (result []*models.EntityAttributeObj, err error) {
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
	result = attrQueryResp.Data
	return
}

func getCMDBAttributeSelectList(entity, userToken string, attributes []string) (result map[string][]*models.EntityDataObj, err error) {
	result = make(map[string][]*models.EntityDataObj)
	attrList, queryErr := getCMDBAttributes(entity, userToken)
	if queryErr != nil {
		err = queryErr
		return
	}
	for _, v := range attrList {
		existFlag := false
		for _, vv := range attributes {
			if v.PropertyName == vv {
				existFlag = true
				break
			}
		}
		if existFlag && v.SelectList != "" {
			tmpSelectList, tmpErr := getAttrCat(v.SelectList, userToken)
			if tmpErr != nil {
				err = tmpErr
				break
			}
			result[entity+"_"+v.PropertyName] = tmpSelectList
		}
	}
	return
}

func getAttrCat(catId, userToken string) (result []*models.EntityDataObj, err error) {
	result = []*models.EntityDataObj{}
	req, newReqErr := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/wecmdb/api/v1/base-key/categories/%s", models.Config.Wecube.BaseUrl, catId), nil)
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
		err = fmt.Errorf("Request cmdb categories fail,%s ", string(responseBytes))
		return
	}
	var response models.CMDBCategoriesResponse
	err = json.Unmarshal(responseBytes, &response)
	if err != nil {
		err = fmt.Errorf("Json unmarshal categories response fail,%s ", err.Error())
		return
	}
	for _, v := range response.Data {
		result = append(result, &models.EntityDataObj{Id: v.Code, DisplayName: v.Value})
	}
	return
}

func BuildRequestProcessData(input models.RequestCacheData) (result models.RequestProcessData) {
	result.ProcDefId = input.ProcDefId
	result.ProcDefKey = input.ProcDefKey
	result.RootEntityOid = input.RootEntityValue.Oid
	result.Entities = []*models.RequestCacheEntityValue{}
	result.Bindings = []*models.RequestProcessTaskNodeBindObj{}
	entityExistMap := make(map[string]int)
	for _, node := range input.TaskNodeBindInfos {
		for _, entity := range node.BoundEntityValues {
			if _, b := entityExistMap[entity.Oid]; !b {
				result.Entities = append(result.Entities, entity)
				entityExistMap[entity.Oid] = 1
			}
			if node.NodeId != "" {
				result.Bindings = append(result.Bindings, &models.RequestProcessTaskNodeBindObj{Oid: entity.Oid, NodeId: node.NodeId, NodeDefId: node.NodeDefId, EntityDataId: entity.EntityDataId, BindFlag: entity.BindFlag})
			}
		}
	}
	if _, b := entityExistMap[input.RootEntityValue.Oid]; !b {
		result.Entities = append(result.Entities, &input.RootEntityValue)
	}
	if len(result.Entities) == 0 {
		tmpEntityValue := models.RequestCacheEntityValue{Oid: result.RootEntityOid, PackageName: "pseudo", EntityName: "pseudo", BindFlag: "N"}
		result.Entities = append(result.Entities, &tmpEntityValue)
	}
	return result
}

func AppendUselessEntity(requestTemplateId, userToken string, cacheData *models.RequestCacheData) (entityDepMap map[string][]string, err error) {
	entityDepMap = make(map[string][]string)
	if cacheData.RootEntityValue.Oid == "" || strings.HasPrefix(cacheData.RootEntityValue.Oid, "tmp") {
		return entityDepMap, nil
	}
	// get core preview data list
	preData, preErr := ProcessDataPreview(requestTemplateId, cacheData.RootEntityValue.Oid, userToken)
	if preErr != nil {
		return entityDepMap, fmt.Errorf("Try to get process preview data fail,%s ", preErr.Error())
	}
	// get binding entity data
	entityList := []*models.RequestCacheEntityValue{}
	entityExistMap := make(map[string]int)
	for _, v := range cacheData.TaskNodeBindInfos {
		for _, vv := range v.BoundEntityValues {
			if _, b := entityExistMap[vv.Oid]; !b {
				entityList = append(entityList, vv)
				entityExistMap[vv.Oid] = 1
			}
		}
	}
	// preEntityList is other entity data
	preEntityList := []*models.EntityTreeObj{}
	rootParent := models.RequestCacheEntityAttrValue{}
	rootSucceeding := []string{}
	for _, v := range preData.Data.EntityTreeNodes {
		if _, b := entityExistMap[v.Id]; !b {
			preEntityList = append(preEntityList, v)
		}
		if v.DataId == cacheData.RootEntityValue.Oid {
			rootSucceeding = v.SucceedingIds
			rootParent.DataType = "ref"
			rootParent.AttrName = v.EntityName
			rootParent.DataValue = v.DataId
		}
	}
	// preEntityList -> in preData but no int boundValues
	if len(preEntityList) == 0 {
		return entityDepMap, nil
	}
	dependEntityMap := make(map[string]*models.RequestCacheEntityAttrValue)
	log.Logger.Info("getDependEntity", log.StringList("rootSucceeding", rootSucceeding), log.Int("preLen", len(preEntityList)), log.Int("entityLen", len(entityList)))
	// entityList -> in boundValue entity
	getDependEntity(rootSucceeding, rootParent, preEntityList, entityList, dependEntityMap)
	for k, refAttr := range dependEntityMap {
		refDataValue := fmt.Sprintf("%s", refAttr.DataValue)
		if _, b := entityDepMap[k]; b {
			entityDepMap[k] = append(entityDepMap[k], refDataValue)
		} else {
			entityDepMap[k] = []string{refDataValue}
		}
		log.Logger.Info("dependEntityMap", log.String("id", k), log.String("refValue", refDataValue))
	}
	if len(dependEntityMap) > 0 {
		newNode := models.RequestCacheTaskNodeBindObj{NodeId: "", NodeDefId: "", BoundEntityValues: []*models.RequestCacheEntityValue{}}
		for _, v := range preEntityList {
			if refAttr, b := dependEntityMap[v.Id]; b {
				newNode.BoundEntityValues = append(newNode.BoundEntityValues, &models.RequestCacheEntityValue{Oid: v.Id, EntityDataId: v.DataId, PackageName: v.PackageName, EntityName: v.EntityName, EntityDisplayName: v.DisplayName, FullEntityDataId: v.FullDataId, AttrValues: []*models.RequestCacheEntityAttrValue{refAttr}, PreviousOids: v.PreviousIds, SucceedingOids: v.SucceedingIds})
			}
		}
		cacheData.TaskNodeBindInfos = append(cacheData.TaskNodeBindInfos, &newNode)
	}
	return entityDepMap, nil
}

func getDependEntity(succeeding []string, parent models.RequestCacheEntityAttrValue, preEntityList []*models.EntityTreeObj, entityList []*models.RequestCacheEntityValue, dependEntityMap map[string]*models.RequestCacheEntityAttrValue) {
	if len(succeeding) == 0 {
		return
	}
	// check if use by entity
	for _, v := range succeeding {
		tmpRefAttr := models.RequestCacheEntityAttrValue{}
		vDataId := v
		if strings.Contains(vDataId, ":") {
			vDataId = vDataId[strings.LastIndex(vDataId, ":")+1:]
		}
		useFlag := false
		for _, vv := range entityList {
			for _, attr := range vv.AttrValues {
				if attr.DataType == "ref" {
					tmpV := fmt.Sprintf("%s", attr.DataValue)
					if strings.Contains(tmpV, vDataId) {
						tmpRefAttr = parent
						useFlag = true
						break
					}
				}
			}
			if useFlag {
				break
			}
		}
		if !useFlag {
			// find other succeeding
			tmpParent := models.RequestCacheEntityAttrValue{}
			nextSucceeding := []string{}
			for _, vv := range preEntityList {
				if vv.DataId == v {
					tmpParent.DataType = "ref"
					tmpParent.DataValue = v
					tmpParent.AttrName = vv.EntityName
				}
				if listContains(vv.PreviousIds, v) {
					nextSucceeding = append(nextSucceeding, vv.Id)
				}
			}
			if len(nextSucceeding) > 0 {
				getDependEntity(nextSucceeding, tmpParent, preEntityList, entityList, dependEntityMap)
				for _, vv := range nextSucceeding {
					if _, b := dependEntityMap[vv]; b {
						useFlag = true
						break
					}
				}
			}
		}
		if useFlag {
			for _, vv := range preEntityList {
				if vv.Id == v {
					dependEntityMap[v] = &tmpRefAttr
				}
			}
		}
	}
}

func listContains(inputList []string, element string) bool {
	result := false
	for _, v := range inputList {
		if v == element {
			result = true
			break
		}
	}
	return result
}

func notifyRoleMail(requestId string) error {
	if !models.MailEnable {
		return nil
	}
	log.Logger.Info("Start notify request mail", log.String("requestId", requestId))
	var roleTable []*models.RoleTable
	err := dao.X.SQL("select id,email from `role` where id in (select `role` from request_template_role where role_type='MGMT' and request_template in (select request_template from request where id=?))", requestId).Find(&roleTable)
	if err != nil {
		return fmt.Errorf("Notify role mail query roles fail,%s ", err.Error())
	}
	if len(roleTable) == 0 {
		return nil
	}
	mailList := getRoleMail(roleTable)
	if len(mailList) == 0 {
		log.Logger.Warn("Notify role mail break,email is empty", log.String("role", roleTable[0].Id))
		return nil
	}
	var requestTable []*models.RequestTable
	dao.X.SQL("select t1.id,t1.name,t2.name as request_template,t1.reporter,t1.report_time,t1.emergency from request t1 left join request_template t2 on t1.request_template=t2.id where t1.id=?", requestId).Find(&requestTable)
	if len(requestTable) == 0 {
		return nil
	}
	var subject, content string
	subject = fmt.Sprintf("Taskman Request [%s] %s[%s]", models.PriorityLevelMap[requestTable[0].Emergency], requestTable[0].Name, requestTable[0].RequestTemplate)
	content = fmt.Sprintf("Taskman Request \nID:%s \nPriority:%s \nName:%s \nTemplate:%s \nReporter:%s \nReportTime:%s\n", requestTable[0].Id, models.PriorityLevelMap[requestTable[0].Emergency], requestTable[0].Name, requestTable[0].RequestTemplate, requestTable[0].Reporter, requestTable[0].ReportTime)
	err = models.MailSender.Send(subject, content, mailList)
	if err != nil {
		log.Logger.Error("Notify role mail fail", log.Error(err))
		return fmt.Errorf("Notify role mail fail,%s ", err.Error())
	}
	return nil
}

func CopyRequest(requestId, createdBy string) (result models.RequestTable, err error) {
	parentRequest := &models.RequestTable{}
	var requestTable []*models.RequestTable
	var requestTemplate models.RequestTemplateTable
	err = dao.X.SQL("select * from request where id=?", requestId).Find(&requestTable)
	if err != nil {
		return
	}
	if len(requestTable) == 0 {
		err = fmt.Errorf("Can not find any request with id:%s ", requestId)
		return
	}
	parentRequest = requestTable[0]
	requestTemplate, err = GetSimpleRequestTemplate(parentRequest.RequestTemplate)
	if err != nil {
		return
	}
	// 重新设置请求名称
	parentRequest.Name = fmt.Sprintf("%s-%s-%s", requestTemplate.Name, requestTemplate.OperatorObjType, time.Now().Format("060102150405"))
	// 重新设置期望时间
	d, _ := time.ParseDuration(fmt.Sprintf("%dh", 24*requestTemplate.ExpireDay))
	parentRequest.ExpectTime = time.Now().Add(d).Format(models.DateTimeFormat)
	result = *parentRequest
	var formTable []*models.FormTable
	err = dao.X.SQL("select * from form where id=?", parentRequest.Form).Find(&formTable)
	if err != nil {
		return
	}
	if len(formTable) == 0 {
		err = fmt.Errorf("Can not find form %s from parent request:%s ", parentRequest.Form, parentRequest.Id)
		return
	}
	parentForm := formTable[0]
	var formItemTable []*models.FormItemTable
	err = dao.X.SQL("select * from form_item where form=?", parentRequest.Form).Find(&formItemTable)
	if err != nil {
		return
	}
	var actions []*dao.ExecAction
	nowTime := time.Now().Format(models.DateTimeFormat)
	formGuid := guid.CreateGuid()
	newRequestId := newRequestId()
	result.Id = newRequestId
	formInsertAction := dao.ExecAction{Sql: "insert into form(id,name,description,form_template,created_time,created_by,updated_time,updated_by) value (?,?,?,?,?,?,?,?)"}
	formInsertAction.Param = []interface{}{formGuid, parentRequest.Name + models.SysTableIdConnector + "form", "", parentForm.FormTemplate, nowTime, createdBy, nowTime, createdBy}
	actions = append(actions, &formInsertAction)
	requestInsertAction := dao.ExecAction{Sql: "insert into request(id,name,form,request_template,reporter,emergency,report_role,status," +
		"cache,expire_time,expect_time,handler,created_by,created_time,updated_by,updated_time,parent,type,role) value (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"}
	requestInsertAction.Param = []interface{}{newRequestId, parentRequest.Name, formGuid, parentRequest.RequestTemplate, createdBy, parentRequest.Emergency, parentRequest.ReportRole, "Draft", parentRequest.Cache,
		"", parentRequest.ExpectTime, parentRequest.Handler, createdBy, nowTime, createdBy, nowTime, parentRequest.Id, parentRequest.Type, parentRequest.Role}
	actions = append(actions, &requestInsertAction)
	for _, formItemRow := range formItemTable {
		actions = append(actions, &dao.ExecAction{Sql: "INSERT INTO form_item (id,form,form_item_template,name,value,item_group,row_data_id) VALUES (?,?,?,?,?,?,?)", Param: []interface{}{
			guid.CreateGuid(), formGuid, formItemRow.FormItemTemplate, formItemRow.Name, formItemRow.Value, formItemRow.ItemGroup, formItemRow.RowDataId,
		}})
	}
	// copy attach file
	var attachFileRows []*models.AttachFileTable
	err = dao.X.SQL("select * from attach_file where request=?", requestId).Find(&attachFileRows)
	if err != nil {
		err = fmt.Errorf("query attach file table fail:%s ", err.Error())
		return
	}
	for _, v := range attachFileRows {
		actions = append(actions, &dao.ExecAction{Sql: "insert into attach_file(id,name,s3_bucket_name,s3_key_name,request,created_by,created_time,updated_by,updated_time) value (?,?,?,?,?,?,?,?,?)", Param: []interface{}{
			guid.CreateGuid(), v.Name, v.S3BucketName, v.S3KeyName, newRequestId, createdBy, nowTime, createdBy, nowTime}})
	}
	err = dao.TransactionWithoutForeignCheck(actions)
	return
}

func GetRequestParent(requestId string) (parentRequestId string, err error) {
	var requestTable []*models.RequestTable
	err = dao.X.SQL("select `parent` from request where id=?", requestId).Find(&requestTable)
	if err != nil {
		return
	}
	if len(requestTable) == 0 {
		err = fmt.Errorf("can not find request with id:%s ", requestId)
		return
	}
	parentRequestId = requestTable[0].Parent
	return
}

func newRequestId() (requestId string) {
	dateString := time.Now().Format("2006-01-02")
	requestId = fmt.Sprintf("%s", time.Now().Format("20060102"))
	requestIdLock.Lock()
	defer requestIdLock.Unlock()
	result, err := dao.X.QueryString(fmt.Sprintf("select count(1) as num from request where created_time>='%s 00:00:00'", dateString))
	if err != nil {
		log.Logger.Error("try to new request id fail with count table num", log.Error(err))
		requestId = fmt.Sprintf("%s-%s", requestId, guid.CreateGuid())
		return
	}
	countNum, _ := strconv.Atoi(result[0]["num"])
	countNumString := fmt.Sprintf("%d", countNum+1)
	subId := countNumString
	for i := 0; i < 6-len(countNumString); i++ {
		subId = "0" + subId
	}
	requestId = fmt.Sprintf("%s-%s", requestId, subId)
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
func GetRequestProgress(requestId, userToken string) (rowsData []*models.RequestProgressObj, err error) {
	var request models.RequestTable
	var pendingHandler string
	var status = int(Completed) //初始化为已完成
	request, err = GetSimpleRequest(requestId)
	if err != nil {
		return
	}
	if request.Handler != "" {
		pendingHandler = request.Handler
	} else {
		pendingHandler = GetRequestTemplateManageRole(request.RequestTemplate)
	}
	subRowsData := getCommonRequestProgress(requestId, request.RequestTemplate, userToken)
	switch request.Status {
	case "Draft":
		status = int(NotStart)
		rowsData = append(append([]*models.RequestProgressObj{
			{
				Node:    SendRequest,
				Handler: request.CreatedBy,
				Status:  int(InProgress),
			},
			{
				Node:    RequestPending,
				Handler: pendingHandler,
				Status:  int(NotStart),
			},
		}), subRowsData...)
		return
	case "Pending":
		status = int(InProgress)
	case "Completed":
		for _, row := range subRowsData {
			row.Status = int(Completed)
		}
	default:
		// 请求状态值,从编排接口读取最新状态值
		if request.ProcInstanceId != "" {
			// 非定版&非完成状态,需要查询 任务节点状态
			taskMap, _ := getTaskMapByRequestId(requestId)
			if len(taskMap) > 0 && len(subRowsData) > 0 {
				for i := len(subRowsData) - 1; i >= 0; i-- {
					if v, ok := taskMap[subRowsData[i].NodeDefId]; ok {
						// 当前任务节点存在,表示前面任务节点都已经完成(任务节点顺序创建,前面完成才会创建后面节点)
						for j := i - 1; j >= 0; j-- {
							subRowsData[j].Status = int(Completed)
						}
						if v.Status == "done" {
							// 任务还在进行中,节点还在审批
							subRowsData[i].Status = int(Completed)
						} else {
							// 任务还在进行中,节点还在审批
							subRowsData[i].Status = int(InProgress)
						}
						break
					}
				}
			}
			response, err := getProcessInstances(request.ProcInstanceId, userToken)
			if err != nil {
				log.Logger.Error("http getProcessInstances error", log.Error(err))
			}

			// 自动退出
			if response.Data.Status == "Faulted" {
				subRowsData[len(subRowsData)-1].Node = AutoExit
				subRowsData[len(subRowsData)-1].Status = int(AutoExitStatus)
			} else {
				// 手动终止
				if response.Data.Status == "InternallyTerminated" {
					subRowsData[len(subRowsData)-1].Node = InternallyTerminated
					subRowsData[len(subRowsData)-1].Status = int(InternallyTerminatedStatus)
				}
				// 记录错误节点,如果实例运行中有错误节点,则需要把运行节点展示在列表中并展示对应位置
				var exist bool
				for _, v := range response.Data.TaskNodeInstances {
					exist = false
					if v.Status == "Faulted" || v.Status == "Timeouted" {
						for _, rowData := range subRowsData {
							if rowData.NodeDefId == v.NodeDefId || rowData.NodeId == v.NodeId {
								exist = true
								rowData.Status = int(Fail)
								break
							}
						}
						if !exist {
							subRowsData = append(subRowsData, &models.RequestProgressObj{
								NodeId:    v.NodeId,
								NodeDefId: v.NodeDefId,
								Node:      v.NodeName,
								Handler:   AutoNode,
								Status:    int(Fail),
								OrderedNo: v.OrderedNo,
							})
						}
						sort.Sort(models.RequestProgressObjSort(subRowsData))
					}
				}
			}
		} else {
			status = int(InProgress)
		}
	}
	rowsData = append(append([]*models.RequestProgressObj{
		{
			Node:    SendRequest,
			Handler: request.CreatedBy,
			Status:  int(Completed),
		},
		{
			Node:    RequestPending,
			Handler: pendingHandler,
			Status:  status,
		},
	}), subRowsData...)
	return
}

func getCommonRequestProgress(requestId, templateId, userToken string) (rowsData []*models.RequestProgressObj) {
	var nodeList models.ProcNodeObjList
	var err error
	nodeList, err = GetProcessNodesByProc(models.RequestTemplateTable{Id: templateId}, userToken, "template")
	if err != nil {
		log.Logger.Error("GetProcessNodesByProc err", log.Error(err))
		rowsData = append(rowsData, &models.RequestProgressObj{
			Node:    RequestComplete,
			Handler: "",
			Status:  int(NotStart),
		})
		return
	}
	if len(nodeList) > 0 {
		// 通过 nodeDefId 和task_template匹配节点
		for _, node := range nodeList {
			dto, err := GetTaskTemplate(templateId, node.NodeDefId, "")
			if err != nil {
				continue
			}
			if dto.Id == "" {
				continue
			}
			rowsData = append(rowsData, &models.RequestProgressObj{
				NodeId:    node.NodeId,
				OrderedNo: node.OrderedNo,
				NodeDefId: node.NodeDefId,
				Node:      node.NodeName,
				Handler:   getTaskApproveHandler(requestId, dto),
				Status:    int(NotStart),
			})
		}
		// 由于编排的编辑会导致 nodeDefId匹配不上,则通过 nodeId 与task_template匹配一遍
		if len(rowsData) == 0 {
			for _, node := range nodeList {
				dto, err := GetTaskTemplate(templateId, "", node.NodeId)
				if err != nil {
					continue
				}
				if dto.Id == "" {
					continue
				}
				rowsData = append(rowsData, &models.RequestProgressObj{
					NodeId:    node.NodeId,
					OrderedNo: node.OrderedNo,
					NodeDefId: node.NodeDefId,
					Node:      node.NodeName,
					Handler:   getTaskApproveHandler(requestId, dto),
					Status:    int(NotStart),
				})
			}
		}
	}
	rowsData = append(rowsData, &models.RequestProgressObj{
		Node:      RequestComplete,
		Handler:   "",
		Status:    int(NotStart),
		OrderedNo: "9999",
	})
	return
}

func GetProcessDefinitions(templateId, userToken string) (rowData *models.DefinitionsData, err error) {
	var template models.RequestTemplateTable
	var response models.ProcessDefinitionsResponse
	var url string
	var byteArr []byte
	template, err = GetSimpleRequestTemplate(templateId)
	if err != nil {
		return
	}
	url = fmt.Sprintf("%s/platform/v1/process/definitions/%s/outline", models.Config.Wecube.BaseUrl, template.ProcDefId)
	byteArr, err = HttpGet(url, userToken)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteArr, &response)
	if err != nil {
		return
	}
	rowData = response.Data
	return
}

func GetProcessInstance(instanceId, userToken string) (rowData *models.ProcessInstance, err error) {
	var byteArr []byte
	var response models.ProcessInstanceResponse
	url := fmt.Sprintf("%s/platform/v1/process/instances/%s", models.Config.Wecube.BaseUrl, instanceId)
	byteArr, err = HttpGet(url, userToken)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteArr, &response)
	if err != nil {
		return
	}
	rowData = response.Data
	return
}

func HttpGet(url, userToken string) (byteArr []byte, err error) {
	req, newReqErr := http.NewRequest(http.MethodGet, url, strings.NewReader(""))
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
	byteArr, _ = ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return
}

// getRequestForm 获取请求信息
func getRequestForm(request *models.RequestTable, userToken string) (form models.RequestForm) {
	if request == nil {
		return
	}
	var tmpTemplate []*models.RequestTemplateTmp
	var version string
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
		form.Progress, form.CurNode = getCurNodeName(request.ProcInstanceId, userToken)
	}
	_, form.Handler = getRequestHandler(request.Id)
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

// getRequestRemainTime 获取请求停留时长
func getRequestRemainTime(requestId string) (startTime string, effectiveDays int) {
	request, _ := GetSimpleRequest(requestId)
	if request.Status == "Draft" || request.Status == "Pending" || request.Status == "Completed" {
		requestTemplate, _ := GetSimpleRequestTemplate(request.RequestTemplate)
		startTime = request.ReportTime
		effectiveDays = requestTemplate.ExpireDay
		return
	}
	// 请求在任务状态
	taskTemplateMap, _ := GetTaskTemplateMapByRequestTemplate(request.RequestTemplate)
	if len(taskTemplateMap) > 0 {
		taskMap, _ := getTaskMapByRequestId(requestId)
		if len(taskMap) > 0 {
			for _, task := range taskMap {
				if task.Status != "done" && taskTemplateMap[task.Name] != 0 {
					startTime = task.CreatedTime
					effectiveDays = taskTemplateMap[task.Name]
					return
				}
			}
		}
	}
	requestTemplate, _ := GetSimpleRequestTemplate(request.RequestTemplate)
	startTime = request.CreatedTime
	effectiveDays = requestTemplate.ExpireDay
	return
}

func GetExecutionNodes(userToken string, procInstanceId, nodeInstanceId string) (data interface{}, err error) {
	var response models.ExecutionResponse
	var byteArr []byte
	var url = fmt.Sprintf("%s/platform/v1/process/instances/%s/tasknodes/%s/context", models.Config.Wecube.BaseUrl, procInstanceId, nodeInstanceId)
	byteArr, err = HttpGet(url, userToken)
	if err != nil {
		return
	}
	err = json.Unmarshal(byteArr, &response)
	if err != nil {
		return
	}
	data = response.Data
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

func getAllCoreProcess(userToken string) map[string]string {
	var processMap = make(map[string]string)
	// 查询全部流程
	result, _ := GetCoreProcessListNew(userToken)
	if len(result) > 0 {
		for _, procDef := range result {
			if procDef != nil {
				processMap[procDef.ProcDefKey] = procDef.RootEntity.DisplayName
			}
		}
	}
	return processMap
}
