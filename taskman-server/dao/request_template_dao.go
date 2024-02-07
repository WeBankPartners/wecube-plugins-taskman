package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type RequestTemplateDao struct {
	DB *xorm.Engine
}

// Add 添加模板,由于直接调用对象插入,formTemplate为""又是外键所以插入报错,所以采用SQL插入
func (d RequestTemplateDao) Add(session *xorm.Session, requestTemplate *models.RequestTemplateTable) (affected int64, err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	return session.InsertOne(requestTemplate)
}

func (d RequestTemplateDao) Update(requestTemplate models.RequestTemplateTable) (err error) {
	_, err = d.DB.ID(requestTemplate.Id).Update(&requestTemplate)
	return
}

func (d RequestTemplateDao) Get(requestTemplateId string) (requestTemplate *models.RequestTemplateTable, err error) {
	_, err = d.DB.ID(requestTemplateId).Get(requestTemplate)
	return
}
