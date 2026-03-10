package logger

import (
	"context"

	"gorm.io/gorm/logger"

	"github.com/goravel/framework/contracts/log"
	"github.com/goravel/framework/support/carbon"
)

// Level log level
type Level int

const (
	Silent Level = iota + 1
	Error
	Warn
	Info
)

type Logger interface {
	Level(Level) Logger
	Log() log.Log
	Infof(context.Context, string, ...any)
	Warningf(context.Context, string, ...any)
	Errorf(context.Context, string, ...any)
	Panicf(context.Context, string, ...any)
	Trace(ctx context.Context, begin *carbon.Carbon, sql string, rowsAffected int64, err error)
	ToGorm() logger.Interface
}
