package models

import (
	"github.com/goravel/framework/database/orm"
)

type User struct {
	orm.Model
	Name     string
	Email    string
	Password string
	orm.SoftDeletes
}
