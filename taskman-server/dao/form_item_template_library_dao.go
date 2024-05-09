package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type FormItemTemplateLibraryDao struct {
	DB *xorm.Engine
}

func (d *FormItemTemplateLibraryDao) Add(session *xorm.Session, formItemTemplateLibrary *models.FormItemTemplateLibraryTable) (affected int64, err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	affected, err = session.Insert(formItemTemplateLibrary)
	// 打印日志
	logExecuteSql(session, "FormItemTemplateLibraryDao", "Add", formItemTemplateLibrary, affected, err)
	return
}

func (d *FormItemTemplateLibraryDao) QueryByFormTemplateLibrary(formTemplateLibrary string) (formItemTemplateLibraryList []*models.FormItemTemplateLibraryTable, err error) {
	formItemTemplateLibraryList = []*models.FormItemTemplateLibraryTable{}
	err = d.DB.Where("form_template_library = ?", formTemplateLibrary).Find(&formItemTemplateLibraryList)
	return
}

func (d *FormItemTemplateLibraryDao) Delete(session *xorm.Session, id string) (err error) {
	var affected int64
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	if id == "" {
		return
	}
	affected, err = session.ID(id).Delete(&models.FormItemTemplateLibraryTable{})
	// 打印日志
	logExecuteSql(session, "FormItemTemplateLibraryDao", "Delete", id, affected, err)
	return
}
