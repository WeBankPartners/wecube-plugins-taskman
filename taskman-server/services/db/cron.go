package db

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/common/log"
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"time"
)

var notifyTaskMap map[string]int

func StartCornJob() {
	go startNotifyCronJob()
	select {}
}

func startNotifyCronJob() {
	notifyTaskMap = make(map[string]int)
	t := time.NewTicker(1 * time.Minute).C
	for {
		go notifyAction()
		<-t
	}
}

func notifyAction() {
	log.Logger.Info("Start notify action")
	var taskTable []*models.TaskTable
	err := x.SQL("select id,created_time,expire_time from task where status<>'done'").Find(&taskTable)
	if err != nil {
		log.Logger.Error("notify action fail,query task error", log.Error(err))
		return
	}
	if len(taskTable) == 0 {
		return
	}
	for _, v := range taskTable {
		if _, b := notifyTaskMap[v.Id]; b {
			continue
		}
		tmpExpireObj := models.ExpireObj{ReportTime: v.CreatedTime, ExpireTime: v.ExpireTime}
		calcExpireObj(&tmpExpireObj)
		if tmpExpireObj.Percent > 75 {
			tmpErr := NotifyTaskMail(v.Id)
			if tmpErr != nil {
				log.Logger.Error("notify task mail fail", log.String("taskId", v.Id), log.Error(tmpErr))
			} else {
				notifyTaskMap[v.Id] = 1
			}
		}
	}
}
