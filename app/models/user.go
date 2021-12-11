package models

import "github.com/goravel/framework/database/orm"

type User struct {
	orm.Model
	orm.SoftDeletes
	Name     string
	Phone    string
	Email    string
	Password string
}
