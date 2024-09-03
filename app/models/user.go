package models

import (
	"github.com/goravel/framework/database/orm"
	"github.com/goravel/framework/support/carbon"
)

type User struct {
	orm.Model
	Name            string
	Email           string
	Password        string
	EmailVerifiedAt *carbon.DateTime
	orm.SoftDeletes
}
