package service

import (
	"fmt"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/dao"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"time"
)

func StartCornJob() {
	go startNotifyCronJob()
	select {}
}

func startNotifyCronJob() {
	t := time.NewTicker(1 * time.Hour).C
	for {
		<-t
		go notifyAction()
	}
}

func notifyAction() {
	log.Logger.Info("Start notify action")
	var taskTable []*models.TaskTable
	var actions []*dao.ExecAction
	err := dao.X.SQL("select id,created_time,expire_time,notify_count,type,request from task where status<>'done'").Find(&taskTable)
	if err != nil {
		log.Logger.Error("notify action fail,query task error", log.Error(err))
		return
	}
	if len(taskTable) == 0 {
		return
	}
	for _, v := range taskTable {
		// 自动确认
		if v.Type == string(models.TaskTypeConfirm) && v.ExpireTime != "" {
			now := time.Now().Format(models.DateTimeFormat)
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
			}
		}
		if v.NotifyCount >= 2 {
			continue
		}
		tmpExpireObj := models.ExpireObj{ReportTime: v.CreatedTime, ExpireTime: v.ExpireTime, NowTime: time.Now().Format(models.DateTimeFormat)}
		calcExpireObj(&tmpExpireObj)
		if ((tmpExpireObj.Percent >= 75) && (v.NotifyCount == 0)) || ((tmpExpireObj.Percent >= 100) && (v.NotifyCount < 2)) {
			mailSubject := "【任务超时提醒】"
			mailContent := ""
			if (tmpExpireObj.Percent >= 75) && (v.NotifyCount == 0) {
				mailContent = fmt.Sprintf("分配给您的任务[请求:%s-任务:%s]快过期了,有效期到%s,请点击链接处理", v.Request, v.Name, v.ExpireTime)
			} else {
				mailContent = fmt.Sprintf("分配给您的任务[请求:%s-任务:%s]已过期,请点击链接尽快处理", v.Request, v.Name)
			}

			tmpErr := NotifyTaskMail(v.Id, models.CoreToken.GetCoreToken(), "", mailSubject, mailContent)
			if tmpErr != nil {
				log.Logger.Error("notify task mail fail", log.String("taskId", v.Id), log.Error(tmpErr))
			} else {
				actions = append(actions, &dao.ExecAction{Sql: "update task set notify_count=? where id=?", Param: []interface{}{v.NotifyCount + 1, v.Id}})
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
