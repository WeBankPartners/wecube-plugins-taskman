package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type TaskHandleTemplateDao struct {
	DB *xorm.Engine
}

func (d *TaskHandleTemplateDao) Add(session *xorm.Session, taskTemplate *models.TaskHandleTemplateTable) (affected int64, err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	affected, err = session.Insert(taskTemplate)
	// 打印日志
	logExecuteSql(session, "TaskHandleTemplateDao", "Add", taskTemplate, affected, err)
	return
}
