package db

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
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
	err := x.SQL("select id,created_time,expire_time,notify_count from task where status<>'done'").Find(&taskTable)
	if err != nil {
		log.Logger.Error("notify action fail,query task error", log.Error(err))
		return
	}
	if len(taskTable) == 0 {
		return
	}
	for _, v := range taskTable {
		if v.NotifyCount >= 1 {
			continue
		}
		tmpExpireObj := models.ExpireObj{ReportTime: v.CreatedTime, ExpireTime: v.ExpireTime, NowTime: time.Now().Format(models.DateTimeFormat)}
		calcExpireObj(&tmpExpireObj)
		if tmpExpireObj.Percent > 75 {
			tmpErr := NotifyTaskMail(v.Id)
			if tmpErr != nil {
				log.Logger.Error("notify task mail fail", log.String("taskId", v.Id), log.Error(tmpErr))
			} else {
				x.Exec("update task set notify_count=? where id=?", v.NotifyCount+1, v.Id)
			}
		}
	}
}
