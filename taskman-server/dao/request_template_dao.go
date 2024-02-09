package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type RequestTemplateDao struct {
	DB *xorm.Engine
}

// Add 添加模板
func (d RequestTemplateDao) Add(session *xorm.Session, requestTemplate *models.RequestTemplateTable) (affected int64, err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	return session.InsertOne(requestTemplate)
}

func (d RequestTemplateDao) Update(session *xorm.Session, requestTemplate *models.RequestTemplateTable) (err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	_, err = d.DB.ID(requestTemplate.Id).Update(requestTemplate)
	return
}

func (d RequestTemplateDao) Get(requestTemplateId string) (*models.RequestTemplateTable, error) {
	var requestTemplate *models.RequestTemplateTable
	var found bool
	var err error
	requestTemplate = &models.RequestTemplateTable{}
	found, err = d.DB.Where("id=?", requestTemplateId).Get(requestTemplate)
	if err != nil {
		return nil, err
	}
	if found {
		return requestTemplate, nil
	}
	return nil, nil
}
