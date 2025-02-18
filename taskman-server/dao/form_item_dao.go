package dao

import (
	"github.com/WeBankPartners/wecube-plugins-taskman/taskman-server/models"
	"xorm.io/xorm"
)

type FormItemDao struct {
	DB *xorm.Engine
}

func (d *FormItemDao) QueryByRequestId(requestId string) (formItemList []*models.FormItemTable, err error) {
	err = d.DB.Where("request = ?", requestId).Find(&formItemList)
	return
}
