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

func (d *TaskTemplateDao) Update(session *xorm.Session, taskTemplate *models.TaskTemplateTable) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if taskTemplate == nil || taskTemplate.Id == "" {
		return
	}
	affected, err := session.ID(taskTemplate.Id).Update(taskTemplate)
	// 打印日志
	logExecuteSql(session, "TaskTemplateDao", "Update", taskTemplate, affected, err)
	return
}

func (d *TaskTemplateDao) Delete(session *xorm.Session, id string) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	_, err = d.DB.ID(id).Delete(&models.TaskTemplateTable{})
	return
}

func (d *TaskTemplateDao) Deletes(session *xorm.Session, ids []string) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	_, err = d.DB.In("id", ids).Delete(&models.TaskTemplateTable{})
	return
}

func (d *TaskTemplateDao) Get(id string) (*models.TaskTemplateTable, error) {
	taskTemplate := &models.TaskTemplateTable{}
	found, err := d.DB.ID(id).Get(taskTemplate)
	if err != nil {
		return nil, err
	}
	if found {
		return taskTemplate, nil
	}
	return nil, nil
}

func (d *TaskTemplateDao) QueryByRequestTemplate(requestTemplateId string) (list []*models.TaskTemplateTable, err error) {
	err = d.DB.Where("request_template=?", requestTemplateId).Asc("sort").Find(&list)
	return
}

func (d *TaskTemplateDao) QueryByRequestTemplateAndType(requestTemplateId, typ string) (list []*models.TaskTemplateTable, err error) {
	err = d.DB.Where("request_template=?", requestTemplateId).And("type=?", typ).Asc("sort").Find(&list)
	return
}

func (d *TaskTemplateDao) GetProc(requestTemplateId, nodeDefId string) (*models.TaskTemplateTable, error) {
	taskTemplate := &models.TaskTemplateTable{}
	found, err := d.DB.Where("request_template=?", requestTemplateId).And("node_def_id=?", nodeDefId).Get(taskTemplate)
	if err != nil {
		return nil, err
	}
	if found {
		return taskTemplate, nil
	}
	return nil, nil
}
