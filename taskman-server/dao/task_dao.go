package dao

import "xorm.io/xorm"

type TaskDao struct {
	DB *xorm.Engine
}
