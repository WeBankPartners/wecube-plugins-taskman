package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type RequestDao struct {
	DB *xorm.Engine
}

func (d *RequestDao) QueryRequestByProcInstanceId(procInstanceId string) (result []*models.RequestTable, err error) {
	err = d.DB.SQL("select id,name,request_template,status,created_by,created_time,updated_by,updated_time,type from request where proc_instance_id = ?", procInstanceId).Find(&result)
	return
}
