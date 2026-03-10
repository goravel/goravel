package driver

import (
	"github.com/goravel/framework/contracts/database"
	"github.com/goravel/framework/contracts/testing/docker"
)

type Driver interface {
	// Docker returns the database driver for Docker.
	Docker() (docker.DatabaseDriver, error)
	// Grammar returns the database grammar.
	Grammar() Grammar
	// Pool returns the database pool.
	Pool() database.Pool
	// Processor returns the database processor.
	Processor() Processor
}
