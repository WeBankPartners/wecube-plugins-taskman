package dao

import "xorm.io/xorm"

type RequestDao struct {
	DB *xorm.Engine
}
