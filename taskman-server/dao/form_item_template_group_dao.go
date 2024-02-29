package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type FormItemTemplateGroupDao struct {
	DB *xorm.Engine
}

func (d *FormItemTemplateGroupDao) Add(session *xorm.Session, formTemplate *models.FormTemplateNewTable) (affected int64, err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	affected, err = session.Insert(formTemplate)
	// 打印日志
	logExecuteSql(session, "FormItemTemplateGroupDao", "Add", formTemplate, affected, err)
	return
}

func (d *FormItemTemplateGroupDao) Update(session *xorm.Session, formTemplate *models.FormTemplateNewTable) (err error) {
	var affected int64
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if formTemplate.Id == "" {
		return
	}
	affected, err = session.ID(formTemplate.Id).Update(formTemplate)
	// 打印日志
	logExecuteSql(session, "FormItemTemplateDao", "Update", formTemplate, affected, err)
	return
}

func (d *FormItemTemplateGroupDao) Get(formItemTemplateGroupId string) (*models.FormTemplateNewTable, error) {
	var formItemTemplateGroup *models.FormTemplateNewTable
	var found bool
	var err error
	formItemTemplateGroup = &models.FormTemplateNewTable{}
	found, err = d.DB.ID(formItemTemplateGroupId).Get(formItemTemplateGroup)
	if err != nil {
		return nil, err
	}
	if found {
		return formItemTemplateGroup, nil
	}
	return nil, nil
}

func (d *FormItemTemplateGroupDao) QueryFormTemplate(formTemplated string) (formItemTemplateGroupList []*models.FormTemplateNewTable, err error) {
	formItemTemplateGroupList = []*models.FormTemplateNewTable{}
	err = d.DB.Where("form_template = ?", formTemplated).Find(&formItemTemplateGroupList)
	if err != nil {
		return
	}
	return
}

func (d *FormItemTemplateGroupDao) DeleteByIdOrRefId(session *xorm.Session, id string) (err error) {
	var affected int64
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if id == "" {
		return
	}
	affected, err = session.Where("id = ? or ref_id = ?", id, id).Delete(&models.FormTemplateNewTable{})
	// 打印日志
	logExecuteSql(session, "FormItemTemplateGroupDao", "DeleteByIdOrRefId", id, affected, err)
	return
}
