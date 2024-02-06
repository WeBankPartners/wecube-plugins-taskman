package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type RequestTemplateDao struct {
	DB *xorm.Engine
}

func (d RequestTemplateDao) Update(requestTemplate models.RequestTemplateTable) (err error) {
	_, err = d.DB.ID(requestTemplate.Id).Update(&requestTemplate)
	return
}

func (d RequestTemplateDao) UpdateByTransaction(session *xorm.Session, requestTemplate models.RequestTemplateTable) (err error) {
	_, err = session.ID(requestTemplate.Id).Update(&requestTemplate)
	return
}

func (d RequestTemplateDao) Get(requestTemplateId string) (requestTemplate *models.RequestTemplateTable, err error) {
	_, err = d.DB.ID(requestTemplateId).Get(requestTemplate)
	return
}
