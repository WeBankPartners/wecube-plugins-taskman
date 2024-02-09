package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type FormItemTemplateDao struct {
	DB *xorm.Engine
}

func (d FormItemTemplateDao) Add(session *xorm.Session, formItemTemplate *models.FormItemTemplateTable) (affected int64, err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	return session.InsertOne(formItemTemplate)
}

func (d FormItemTemplateDao) Update(session *xorm.Session, formItemTemplate *models.FormItemTemplateTable) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	_, err = d.DB.ID(formItemTemplate.Id).Update(formItemTemplate)
	return
}

func (d FormItemTemplateDao) Get(formItemTemplateId string) (formItemTemplate *models.FormItemTemplateTable, err error) {
	formItemTemplate = &models.FormItemTemplateTable{}
	_, err = d.DB.ID(formItemTemplateId).Get(formItemTemplate)
	return
}

func (d FormItemTemplateDao) QueryByFormTemplate(formTemplate string) (formItemTemplate []*models.FormItemTemplateTable, err error) {
	formItemTemplate = []*models.FormItemTemplateTable{}
	err = d.DB.Where("form_template = ?", formTemplate).Find(&formItemTemplate)
	return
}

func (d FormItemTemplateDao) Delete(session *xorm.Session, id string) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	_, err = d.DB.Where("id=?", id).Delete(&models.FormItemTemplateTable{})
	return
}
