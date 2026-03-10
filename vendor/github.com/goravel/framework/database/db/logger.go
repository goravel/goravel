package db

import (
	"context"
	"time"

	gormlogger "gorm.io/gorm/logger"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/database/logger"
	"github.com/goravel/framework/contracts/log"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/support/str"
)

var (
	traceStr     = "[%.3fms] [rows:%v] %s"
	traceWarnStr = "[%.3fms] [rows:%v] [SLOW] %s"
	traceErrStr  = "[%.3fms] [rows:%v] %s\t%s"
)

func NewLogger(config config.Config, log log.Log) logger.Logger {
	level := logger.Warn
	if config.GetBool("app.debug") {
		level = logger.Info
	}

	slowThreshold := config.GetInt("database.slow_threshold", 200)
	if slowThreshold <= 0 {
		slowThreshold = 200
	}

	return &Logger{
		log:           log,
		level:         level,
		slowThreshold: time.Duration(slowThreshold) * time.Millisecond,
	}
}

type Logger struct {
	log           log.Log
	level         logger.Level
	slowThreshold time.Duration
}

func (r *Logger) Log() log.Log {
	return r.log
}

func (r *Logger) Level(level logger.Level) logger.Logger {
	r.level = level

	return r
}

func (r *Logger) Infof(ctx context.Context, msg string, data ...any) {
	if r.level >= logger.Info {
		r.log.WithContext(ctx).Infof(msg, data...)
	}
}

func (r *Logger) Warningf(ctx context.Context, msg string, data ...any) {
	if r.level >= logger.Warn {
		r.log.WithContext(ctx).Warningf(msg, data...)
	}
}

func (r *Logger) Errorf(ctx context.Context, msg string, data ...any) {
	for _, item := range data {
		if tempItem, ok := item.(error); ok {
			if str.Of(tempItem.Error()).Contains("Access denied", "connection refused") {
				return
			}
		}
	}

	if r.level >= logger.Error {
		r.log.WithContext(ctx).Errorf(msg, data...)
	}
}

func (r *Logger) Panicf(ctx context.Context, msg string, data ...any) {
	r.log.WithContext(ctx).Panicf(msg, data...)
}

func (r *Logger) Trace(ctx context.Context, begin *carbon.Carbon, sql string, rowsAffected int64, err error) {
	if r.level <= logger.Silent {
		return
	}

	duration := begin.DiffInDuration()
	elapsed := float64(duration.Nanoseconds()) / 1e6

	addQueryLogToContext(ctx, sql, elapsed)

	switch {
	case err != nil && r.level >= logger.Error && !errors.Is(err, gormlogger.ErrRecordNotFound):
		if rowsAffected == -1 {
			r.Errorf(ctx, traceErrStr, elapsed, "-", sql, err)
		} else {
			r.Errorf(ctx, traceErrStr, elapsed, rowsAffected, sql, err)
		}
	case duration > r.slowThreshold && r.slowThreshold != 0 && r.level >= logger.Warn:
		if rowsAffected == -1 {
			r.Warningf(ctx, traceWarnStr, elapsed, "-", sql)
		} else {
			r.Warningf(ctx, traceWarnStr, elapsed, rowsAffected, sql)
		}
	case r.level == logger.Info:
		if rowsAffected == -1 {
			r.Infof(ctx, traceStr, elapsed, "-", sql)
		} else {
			r.Infof(ctx, traceStr, elapsed, rowsAffected, sql)
		}
	}
}

func (r *Logger) ToGorm() gormlogger.Interface {
	return NewGorm(r)
}

type Gorm struct {
	logger logger.Logger
}

func NewGorm(logger logger.Logger) *Gorm {
	return &Gorm{
		logger: logger,
	}
}

func (r *Gorm) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	_ = r.logger.Level(gormLevelToLevel(level))

	return r
}

func (r *Gorm) Info(ctx context.Context, msg string, data ...any) {
	r.logger.Infof(ctx, msg, data...)
}

func (r *Gorm) Warn(ctx context.Context, msg string, data ...any) {
	r.logger.Warningf(ctx, msg, data...)
}

func (r *Gorm) Error(ctx context.Context, msg string, data ...any) {
	r.logger.Errorf(ctx, msg, data...)
}

func (r *Gorm) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rowsAffected := fc()
	r.logger.Trace(ctx, carbon.FromStdTime(begin), sql, rowsAffected, err)
}

func addQueryLogToContext(ctx context.Context, sql string, time float64) {
	value := ctx.Value(queryLogKey{})
	if value == nil {
		return
	}

	queryLogValue := value.(*queryLogValue)
	if !queryLogValue.enabled {
		return
	}

	queryLogValue.queryLogs = append(queryLogValue.queryLogs, QueryLog{
		Query: sql,
		Time:  time,
	})

	_ = context.WithValue(ctx, queryLogKey{}, queryLogValue)
}

func gormLevelToLevel(level gormlogger.LogLevel) logger.Level {
	switch level {
	case gormlogger.Silent:
		return logger.Silent
	case gormlogger.Error:
		return logger.Error
	case gormlogger.Warn:
		return logger.Warn
	default:
		return logger.Info
	}
}
