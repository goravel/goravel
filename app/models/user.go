package models

import (
	"time"
)

type User struct {
	ID        uint64
	Name      string
	Phone     string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
