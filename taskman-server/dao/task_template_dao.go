package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type TaskTemplateDao struct {
	DB *xorm.Engine
}

func (d *TaskTemplateDao) Add(session *xorm.Session, taskTemplate *models.TaskTemplateTable) (affected int64, err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	affected, err = session.Insert(taskTemplate)
	// 打印日志
	logExecuteSql(session, "TaskTemplateDao", "Add", taskTemplate, affected, err)
	return
}

func (d *TaskTemplateDao) QueryByRequestTemplate(requestTemplateId string) (list []*models.TaskTemplateTable, err error) {
	err = d.DB.Where("request_template=?", requestTemplateId).Find(&list)
	return
}
