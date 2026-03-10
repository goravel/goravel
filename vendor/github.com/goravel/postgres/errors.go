package postgres

import "github.com/goravel/framework/errors"

var (
	FailedToGenerateDSN = errors.New("failed to generate DSN, please check the database configuration")
	ConfigNotFound      = errors.New("not found database configuration")
)
