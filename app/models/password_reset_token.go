package models

import "github.com/goravel/framework/database/orm"

type PasswordResetToken struct {
	Email string `gorm:"primaryKey"`
	Token string
	orm.Timestamps
}
