package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type TaskDao struct {
	DB *xorm.Engine
}

func (d *TaskDao) Add(session *xorm.Session, task *models.TaskTable) (affected int64, err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	affected, err = session.Insert(task)
	// 打印日志
	logExecuteSql(session, "TaskDao", "Add", task, affected, err)
	return
}

func (d *TaskDao) Update(session *xorm.Session, task *models.TaskTable) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if task == nil || task.Id == "" {
		return
	}
	affected, err := session.ID(task.Id).Update(task)
	// 打印日志
	logExecuteSql(session, "TaskDao", "Update", task, affected, err)
	return
}

func (d *TaskDao) Delete(session *xorm.Session, id string) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	_, err = d.DB.ID(id).Delete(&models.TaskTable{})
	return
}

func (d *TaskDao) Deletes(session *xorm.Session, ids []string) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if len(ids) == 0 {
		return
	}
	_, err = d.DB.In("id", ids).Delete(&models.TaskTable{})
	return
}

func (d *TaskDao) QueryByRequest(requestId string) (list []*models.TaskTable, err error) {
	err = d.DB.Where("request=?", requestId).Asc("sort").Find(&list)
	return
}

func (d *TaskDao) QueryByRequestAndType(requestId, typ string) (list []*models.TaskTable, err error) {
	err = d.DB.Where("request=?", requestId).And("type=?", typ).Asc("sort").Find(&list)
	return
}
