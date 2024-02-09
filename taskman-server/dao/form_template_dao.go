package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type FormTemplateDao struct {
	DB *xorm.Engine
}

func (d FormTemplateDao) Add(session *xorm.Session, formTemplate *models.FormTemplateTable) (affected int64, err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	return session.InsertOne(formTemplate)
}

func (d FormTemplateDao) Update(session *xorm.Session, formTemplate *models.FormTemplateTable) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	_, err = d.DB.Where("id=?", formTemplate.Id).Update(formTemplate)
	return
}

func (d FormTemplateDao) Get(formTemplateId string) (*models.FormTemplateTable, error) {
	var formTemplate *models.FormTemplateTable
	var found bool
	var err error
	formTemplate = &models.FormTemplateTable{}
	found, err = d.DB.Where("id=?", formTemplateId).Get(formTemplate)
	if err != nil {
		return nil, err
	}
	if found {
		return formTemplate, nil
	}
	return nil, nil
}
