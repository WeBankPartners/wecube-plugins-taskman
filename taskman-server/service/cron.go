package service

import (
	"github.com/WeBankPartners/go-common-lib/guid"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"time"
)

const periods = 10

func StartCornJob() {
	go startNotifyCronJob()
	select {}
}

func startNotifyCronJob() {
	t := time.NewTicker(periods * time.Minute).C
	for {
		<-t
		go notifyAction()
		go notifyAndUpdateWorkflowResult()
	}
}

func notifyAction() {
	log.Logger.Info("Start notify action")
	var taskTable []*models.TaskTable
	var actions []*dao.ExecAction
	var taskNotifyList []*models.TaskNotifyTable
	var yesterday time.Time
	var doingNotifyCount, timeoutNotifyCount int
	now := time.Now().Format(models.DateTimeFormat)
	yesterday = time.Now().AddDate(0, 0, -1)
	err := dao.X.SQL("select id,name,created_time,expire_time,notify_count,type,request from task where status<>'done' and created_time >= ? and created_time <= ?", yesterday.Format(models.DateTimeFormat), now).Find(&taskTable)
	if err != nil {
		log.Logger.Error("notify action fail,query task error", log.Error(err))
		return
	}
	if len(taskTable) == 0 {
		return
	}
	for _, v := range taskTable {
		taskNotifyList = []*models.TaskNotifyTable{}
		beforeMinutes := time.Now().Add(-periods * time.Minute).Format(models.DateTimeFormat)
		dao.X.SQL("select doing_notify_count,timeout_notify_count,updated_time from task_notify where task = ? and updated_time > ?", v.Id, beforeMinutes).Find(&taskNotifyList)
		if len(taskNotifyList) > 0 {
			continue
		}
		// 自动确认
		if v.Type == string(models.TaskTypeConfirm) && v.ExpireTime != "" {
			expireT, _ := time.Parse(models.DateTimeFormat, v.ExpireTime)
			nowT, _ := time.Parse(models.DateTimeFormat, now)
			if nowT.Sub(expireT) > 0 {
				taskHandleList, _ := GetTaskHandleService().GetTaskHandleListByTaskId(v.Id)
				if len(taskHandleList) > 0 {
					actions = append(actions, &dao.ExecAction{Sql: "update task_handle set handle_result =?,handle_status=?,updated_time=? where id=?", Param: []interface{}{models.TaskHandleResultTypeApprove, models.TaskHandleResultTypeComplete, now, taskHandleList[0].Id}})
				}
				actions = append(actions, &dao.ExecAction{Sql: "update task set status =?,task_result=?,updated_by=?,updated_time=? where id=?", Param: []interface{}{models.TaskStatusDone, models.TaskHandleResultTypeComplete, "system", now, v.Id}})
				if v.Request != "" {
					actions = append(actions, &dao.ExecAction{Sql: "update request set status =?,updated_by=?,updated_time=?,complete_status=? where id=?", Param: []interface{}{models.RequestStatusCompleted, "system", now, models.TaskHandleResultTypeComplete, v.Request}})
				}
				if v.Request != "" {
					var requestList []*models.RequestTable
					dao.X.SQL("select name,created_by from request where id = ?", v.Request).Find(&requestList)
					if len(requestList) > 0 {
						// 请求完成,给创建人发邮件
						NotifyRequestCompleteMail(requestList[0].Name, requestList[0].CreatedBy, models.CoreToken.GetCoreToken(), "")
					}
				}
			}
		}
		doingNotifyCount = 0
		timeoutNotifyCount = 0
		dao.X.SQL("select id,doing_notify_count,timeout_notify_count,updated_time from task_notify where task = ? order by updated_time desc limit 0,1", v.Id).Find(&taskNotifyList)
		if len(taskNotifyList) > 0 {
			doingNotifyCount = taskNotifyList[0].DoingNotifyCount
			timeoutNotifyCount = taskNotifyList[0].TimeoutNotifyCount
		}
		// 通知过一次,则直接跳过
		if doingNotifyCount >= 1 && timeoutNotifyCount >= 1 {
			continue
		}
		tmpExpireObj := &models.ExpireObj{
			ReportTime:         v.CreatedTime,
			ExpireTime:         v.ExpireTime,
			NowTime:            time.Now().Format(models.DateTimeFormat),
			DoingNotifyCount:   doingNotifyCount,
			TimeoutNotifyCount: timeoutNotifyCount,
		}
		calcExpireObj(tmpExpireObj)
		if (tmpExpireObj.Percent >= 75) && (doingNotifyCount == 0) || ((tmpExpireObj.Percent >= 100) && (timeoutNotifyCount == 0)) {
			tmpErr := NotifyTaskExpireMail(v, tmpExpireObj, models.CoreToken.GetCoreToken(), "")
			var errMsg string
			if tmpErr != nil {
				errMsg = tmpErr.Error()
				log.Logger.Error("notify task mail fail", log.String("taskId", v.Id), log.Error(tmpErr))
			}
			if len(taskNotifyList) > 0 {
				actions = append(actions, &dao.ExecAction{Sql: "update notity_task set doing_notify_count = ?,timeout_notify_count = ?,err_msg = ?,updated_time = ? where id=?", Param: []interface{}{tmpExpireObj.DoingNotifyCount, tmpExpireObj.TimeoutNotifyCount, errMsg, time.Now().Format(models.DateTimeFormat), taskNotifyList[0].Id}})
			} else {
				actions = append(actions, &dao.ExecAction{Sql: "insert into notity_task(id,task,doing_notify_count,timeout_notify_count,err_msg,updated_time)values(?,?,?,?,?,?)", Param: []interface{}{guid.CreateGuid(), tmpExpireObj.DoingNotifyCount, tmpExpireObj.TimeoutNotifyCount, errMsg, time.Now().Format(models.DateTimeFormat)}})
			}
		}
	}
	if len(actions) > 0 {
		err = dao.Transaction(actions)
		if err != nil {
			log.Logger.Error("notify action error", log.Error(err))
		}
	}
}

// notifyAndUpdateWorkflowResult 通知并且更新编排结果
func notifyAndUpdateWorkflowResult() {
	var requestList []*models.RequestTable
	var requestTemplate *models.RequestTemplateTable
	var actions []*dao.ExecAction
	var yesterday time.Time
	var err error
	yesterday = time.Now().AddDate(0, 0, -1)
	err = dao.X.SQL("select id,name,proc_instance_id,request_template,created_by from request where status = ? and proc_instance_id is not null and created_time >= ? and created_time <= ?",
		models.RequestStatusInProgress, yesterday.Format(models.DateTimeFormat), time.Now().Format(models.DateTimeFormat)).Find(&requestList)
	if err != nil {
		log.Logger.Error("notifyAndUpdateWorkflowResult fail,query request error", log.Error(err))
		return
	}
	if len(requestList) > 0 {
		for _, request := range requestList {
			newStatus := getInstanceStatus(request.ProcInstanceId, models.CoreToken.GetCoreToken(), "")
			if newStatus == "InternallyTerminated" {
				newStatus = "Termination"
			}
			// 只处理自动退出&手动终止终止情况,需要发邮件
			if newStatus == string(models.RequestStatusFaulted) || newStatus == string(models.RequestStatusTermination) {
				actions = append(actions, &dao.ExecAction{Sql: "update request set status=?,updated_time=? where id=?", Param: []interface{}{newStatus, time.Now().Format(models.DateTimeFormat), request.Id}})
				if requestTemplate, err = GetRequestTemplateService().GetRequestTemplate(request.RequestTemplate); err != nil {
					continue
				}
				NotifyTaskWorkflowFailMail(request.Name, requestTemplate.ProcDefName, newStatus, request.CreatedBy, models.CoreToken.GetCoreToken(), "")
			}
		}
	}
	if len(actions) > 0 {
		if err = dao.Transaction(actions); err != nil {
			log.Logger.Error("notifyAndUpdateWorkflowResult  error", log.Error(err))
		}
	}
}
