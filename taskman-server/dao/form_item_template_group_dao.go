package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type FormItemTemplateGroupDao struct {
	DB *xorm.Engine
}

func (d FormItemTemplateGroupDao) Add(session *xorm.Session, formItemTemplate *models.FormItemTemplateGroupTable) (affected int64, err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	affected, err = session.Insert(formItemTemplate)
	// 打印日志
	logExecuteSql(session, "FormItemTemplateGroupDao", "Add", formItemTemplate, affected, err)
	return
}

func (d FormItemTemplateGroupDao) Update(session *xorm.Session, formItemTemplateGroup *models.FormItemTemplateGroupTable) (err error) {
	var affected int64
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if formItemTemplateGroup.Id == "" {
		return
	}
	affected, err = session.ID(formItemTemplateGroup.Id).Update(formItemTemplateGroup)
	// 打印日志
	logExecuteSql(session, "FormItemTemplateDao", "Update", formItemTemplateGroup, affected, err)
	return
}

func (d FormItemTemplateGroupDao) Get(formItemTemplateGroupId string) (*models.FormItemTemplateGroupTable, error) {
	var formItemTemplateGroup *models.FormItemTemplateGroupTable
	var found bool
	var err error
	formItemTemplateGroup = &models.FormItemTemplateGroupTable{}
	found, err = d.DB.ID(formItemTemplateGroupId).Get(formItemTemplateGroup)
	if err != nil {
		return nil, err
	}
	if found {
		return formItemTemplateGroup, nil
	}
	return nil, nil
}

func (d FormItemTemplateGroupDao) QueryFormTemplate(formTemplated string) (formItemTemplateGroupList []*models.FormItemTemplateGroupTable, err error) {
	formItemTemplateGroupList = []*models.FormItemTemplateGroupTable{}
	err = d.DB.Where("form_template = ?", formTemplated).Find(&formItemTemplateGroupList)
	if err != nil {
		return
	}
	return
}

func (d FormItemTemplateGroupDao) DeleteByIdOrCopyId(session *xorm.Session, id string) (err error) {
	var affected int64
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if id == "" {
		return
	}
	affected, err = session.Where("id = ? or copy_id = ?", id, id).Delete(&models.FormItemTemplateGroupTable{})
	// 打印日志
	logExecuteSql(session, "FormItemTemplateGroupDao", "DeleteByIdOrCopyId", id, affected, err)
	return
}
