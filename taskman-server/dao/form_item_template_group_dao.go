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

func (d FormItemTemplateGroupDao) Get(formItemTemplateGroupId string) (*models.FormItemTemplateTable, error) {
	var formItemTemplate *models.FormItemTemplateTable
	var found bool
	var err error
	formItemTemplate = &models.FormItemTemplateTable{}
	found, err = d.DB.ID(formItemTemplateGroupId).Get(formItemTemplate)
	if err != nil {
		return nil, err
	}
	if found {
		return formItemTemplate, nil
	}
	return nil, nil
}
