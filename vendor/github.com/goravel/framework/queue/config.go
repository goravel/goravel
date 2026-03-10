package queue

import (
	"fmt"

	contractsconfig "github.com/goravel/framework/contracts/config"
)

type Config struct {
	contractsconfig.Config

	appName           string
	defaultConnection string
	defaultQueue      string
	failedDatabase    string
	failedTable       string
	defaultConcurrent int
	debug             bool
}

func NewConfig(config contractsconfig.Config) *Config {
	defaultConnection := config.GetString("queue.default")
	defaultQueue := config.GetString(fmt.Sprintf("queue.connections.%s.queue", defaultConnection), "default")
	defaultConcurrent := max(config.GetInt(fmt.Sprintf("queue.connections.%s.concurrent", defaultConnection), 1), 1)

	c := &Config{
		Config: config,

		appName:           config.GetString("app.name", "goravel"),
		debug:             config.GetBool("app.debug"),
		defaultConnection: defaultConnection,
		defaultQueue:      defaultQueue,
		defaultConcurrent: defaultConcurrent,
		failedDatabase:    config.GetString("queue.failed.database"),
		failedTable:       config.GetString("queue.failed.table"),
	}

	return c
}

func (r *Config) Debug() bool {
	return r.debug
}

func (r *Config) DefaultConnection() string {
	return r.defaultConnection
}

func (r *Config) DefaultQueue() string {
	return r.defaultQueue
}

func (r *Config) DefaultConcurrent() int {
	return r.defaultConcurrent
}

func (r *Config) Driver(connection string) string {
	return r.GetString(fmt.Sprintf("queue.connections.%s.driver", connection))
}

func (r *Config) FailedDatabase() string {
	return r.failedDatabase
}

func (r *Config) FailedTable() string {
	return r.failedTable
}

func (r *Config) Via(connection string) any {
	return r.Get(fmt.Sprintf("queue.connections.%s.via", connection))
}
