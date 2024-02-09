package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type FormDao struct {
	DB *xorm.Engine
}

func (d FormDao) Delete(formId string) (err error) {
	_, err = d.DB.Where("id = ?", formId).Delete(&models.FormItemTemplateTable{})
	return
}

func (d FormDao) DeleteByFormItemTemplate(session *xorm.Session, formItemTemplate string) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	_, err = session.Where("form_item_template = ?", formItemTemplate).Delete(&models.FormTemplateTable{})
	return
}
