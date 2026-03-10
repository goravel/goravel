package dbresolver

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ResolverModeKey string
type ResolverMode string

const resolverModeKey ResolverModeKey = "dbresolver:resolver_mode_key"
const (
	ResolverModeSource  ResolverMode = "source"
	ResolverModeReplica ResolverMode = "replica"
)

type resolverModeLogger struct {
	logger.Interface
}

func (l resolverModeLogger) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if filter, ok := l.Interface.(gorm.ParamsFilter); ok {
		sql, params = filter.ParamsFilter(ctx, sql, params...)
	}
	return sql, params
}

func (l resolverModeLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.Interface = l.Interface.LogMode(level)
	return l
}

func (l resolverModeLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	var splitFn = func() (sql string, rowsAffected int64) {
		sql, rowsAffected = fc()
		op := ctx.Value(resolverModeKey)
		if op != nil {
			sql = fmt.Sprintf("[%s] %s", op, sql)
			return
		}

		// the situation that dbresolver does not handle
		// such as transactions, or some resolvers do not enable MarkResolverMode.
		return
	}
	l.Interface.Trace(ctx, begin, splitFn, err)
}

func NewResolverModeLogger(l logger.Interface) logger.Interface {
	if _, ok := l.(resolverModeLogger); ok {
		return l
	}
	return resolverModeLogger{
		Interface: l,
	}
}

func markStmtResolverMode(stmt *gorm.Statement, mode ResolverMode) {
	if _, ok := stmt.Logger.(resolverModeLogger); ok {
		stmt.Context = context.WithValue(stmt.Context, resolverModeKey, mode)
	}
}
