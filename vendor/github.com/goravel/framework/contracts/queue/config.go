package queue

import (
	"github.com/goravel/framework/contracts/config"
)

type Config interface {
	config.Config
	Debug() bool
	DefaultConnection() string
	DefaultQueue() string
	DefaultConcurrent() int
	Driver(connection string) string
	FailedDatabase() string
	FailedTable() string
	Via(connection string) any
}
