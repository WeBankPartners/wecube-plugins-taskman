package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type RequestTemplateRoleDao struct {
	DB *xorm.Engine
}

func (d *RequestTemplateRoleDao) CheckRequestTemplateRoles(requestTemplateId string, userRoles []string) (bool, error) {
	return d.DB.Table(models.RequestTemplateRoleTable{}.TableName()).Where("request_template=?", requestTemplateId).And("role_type=?",
		models.RolePermissionMGMT).In("role", userRoles).Exist()
}

func (d *RequestTemplateRoleDao) Add(session *xorm.Session, requestTemplateRole *models.RequestTemplateRoleTable) (affected int64, err error) {
	if session == nil {
		session = d.DB.NewSession()
		defer session.Close()
	}
	affected, err = session.Insert(requestTemplateRole)
	// 打印日志
	logExecuteSql(session, "RequestTemplateRoleDao", "Add", requestTemplateRole, affected, err)
	return
}
