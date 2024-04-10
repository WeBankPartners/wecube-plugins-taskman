package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type OperationLogDao struct {
	DB *xorm.Engine
}

func (d *OperationLogDao) AddOperationLog(record *models.OperationLogTable) (affected int64, err error) {
	return d.DB.Insert(record)
}
