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
	_, err = d.DB.ID(formTemplate.Id).Update(formTemplate)
	return
}

func (d FormTemplateDao) Get(formTemplateId string) (formTemplate *models.FormTemplateTable, err error) {
	_, err = d.DB.ID(formTemplateId).Get(&formTemplate)
	return
}
