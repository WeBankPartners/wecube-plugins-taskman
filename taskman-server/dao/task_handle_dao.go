package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type TaskHandleDao struct {
	DB *xorm.Engine
}

func (d *TaskHandleDao) Add(session *xorm.Session, taskHandle *models.TaskHandleTable) (affected int64, err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	affected, err = session.Insert(taskHandle)
	// 打印日志
	logExecuteSql(session, "TaskHandleDao", "Add", taskHandle, affected, err)
	return
}

func (d *TaskHandleDao) Update(session *xorm.Session, taskHandle *models.TaskHandleTable) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if taskHandle == nil || taskHandle.Id == "" {
		return
	}
	affected, err := session.ID(taskHandle.Id).Update(taskHandle)
	// 打印日志
	logExecuteSql(session, "TaskHandleDao", "Update", taskHandle, affected, err)
	return
}

func (d *TaskHandleDao) Delete(session *xorm.Session, id string) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	_, err = session.ID(id).Delete(&models.TaskHandleTable{})
	return
}

func (d *TaskHandleDao) Deletes(session *xorm.Session, ids []string) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if len(ids) == 0 {
		return
	}
	_, err = session.In("id", ids).Delete(&models.TaskHandleTable{})
	return
}

func (d *TaskHandleDao) DeleteByTask(session *xorm.Session, taskId string) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	_, err = session.Where("task=?", taskId).Delete(&models.TaskHandleTable{})
	return
}

func (d *TaskHandleDao) DeleteByTasks(session *xorm.Session, taskIds []string) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if len(taskIds) == 0 {
		return
	}
	_, err = session.In("task", taskIds).Delete(&models.TaskHandleTable{})
	return
}

func (d *TaskHandleDao) QueryByTask(taskId string) (list []*models.TaskHandleTable, err error) {
	err = d.DB.Where("task=?", taskId).Asc("sort").Find(&list)
	return
}

func (d *TaskHandleDao) QueryByTasks(taskIds []string) (list []*models.TaskHandleTable, err error) {
	err = d.DB.In("task", taskIds).Asc("task", "sort").Find(&list)
	return
}
