package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type TaskHandleTemplateDao struct {
	DB *xorm.Engine
}

func (d *TaskHandleTemplateDao) Add(session *xorm.Session, taskHandleTemplate *models.TaskHandleTemplateTable) (affected int64, err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	affected, err = session.Insert(taskHandleTemplate)
	// 打印日志
	logExecuteSql(session, "TaskHandleTemplateDao", "Add", taskHandleTemplate, affected, err)
	return
}

func (d *TaskHandleTemplateDao) Update(session *xorm.Session, taskHandleTemplate *models.TaskHandleTemplateTable) (err error) {
	var affected int64
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if taskHandleTemplate == nil || taskHandleTemplate.Id == "" {
		return
	}
	affected, err = session.ID(taskHandleTemplate.Id).Update(taskHandleTemplate)
	// 打印日志
	logExecuteSql(session, "TaskHandleTemplateDao", "Update", taskHandleTemplate, affected, err)
	if err != nil {
		return
	}
	return
}

func (d *TaskHandleTemplateDao) Delete(session *xorm.Session, id string) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	_, err = d.DB.ID(id).Delete(&models.TaskHandleTemplateTable{})
	return
}

func (d *TaskHandleTemplateDao) Deletes(session *xorm.Session, ids []string) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if len(ids) == 0 {
		return
	}
	_, err = d.DB.In("id", ids).Delete(&models.TaskHandleTemplateTable{})
	return
}

func (d *TaskHandleTemplateDao) DeleteByTaskTemplate(session *xorm.Session, taskTemplateId string) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	_, err = d.DB.Where("task_template=?", taskTemplateId).Delete(&models.TaskHandleTemplateTable{})
	return
}

func (d *TaskHandleTemplateDao) DeleteByTaskTemplates(session *xorm.Session, taskTemplateIds []string) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if len(taskTemplateIds) == 0 {
		return
	}
	_, err = d.DB.In("task_template", taskTemplateIds).Delete(&models.TaskHandleTemplateTable{})
	return
}

func (d *TaskHandleTemplateDao) Get(id string) (*models.TaskHandleTemplateTable, error) {
	taskHandleTemplate := &models.TaskHandleTemplateTable{}
	found, err := d.DB.ID(id).Get(taskHandleTemplate)
	if err != nil {
		return nil, err
	}
	if found {
		return taskHandleTemplate, nil
	}
	return nil, nil
}

func (d *TaskHandleTemplateDao) QueryByTaskTemplate(taskTemplateId string) (list []*models.TaskHandleTemplateTable, err error) {
	err = d.DB.Where("task_template=?", taskTemplateId).Asc("sort").Find(&list)
	return
}
