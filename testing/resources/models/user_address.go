package models

import (
	"github.com/goravel/framework/database/orm"
)

type UserAddress struct {
	orm.Model
	UserId   uint
	Name     string
	Province string
}
