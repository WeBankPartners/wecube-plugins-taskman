package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type FormTemplateDao struct {
	DB *xorm.Engine
}

func (d *FormTemplateDao) Add(session *xorm.Session, formTemplate *models.FormTemplateTable) (affected int64, err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	affected, err = session.Insert(formTemplate)
	// 打印日志
	logExecuteSql(session, "FormTemplateDao", "Add", formTemplate, affected, err)
	return
}

func (d *FormTemplateDao) Update(session *xorm.Session, formTemplate *models.FormTemplateTable) (err error) {
	var affected int64
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if formTemplate == nil || formTemplate.Id == "" {
		return
	}
	affected, err = session.ID(formTemplate.Id).Update(formTemplate)
	// 打印日志
	logExecuteSql(session, "FormTemplateDao", "Update", formTemplate, affected, err)
	if err != nil {
		return
	}
	return
}

func (d *FormTemplateDao) Get(formTemplateId string) (*models.FormTemplateTable, error) {
	var formTemplate *models.FormTemplateTable
	var found bool
	var err error
	formTemplate = &models.FormTemplateTable{}
	found, err = d.DB.ID(formTemplateId).Get(formTemplate)
	if err != nil {
		return nil, err
	}
	if found {
		return formTemplate, nil
	}
	return nil, nil
}

func (d *FormTemplateDao) Delete(session *xorm.Session, id string) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	_, err = d.DB.ID(id).Delete(&models.FormItemTemplateTable{})
	return
}
