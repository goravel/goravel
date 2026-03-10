package db

import (
	"context"
)

type queryLogKey struct{}

type queryLogValue struct {
	queryLogs []QueryLog
	enabled   bool
}

type QueryLog struct {
	Query string  `json:"query"`
	Time  float64 `json:"time"`
}

func DisableQueryLog(ctx context.Context) context.Context {
	return context.WithValue(ctx, queryLogKey{}, &queryLogValue{
		enabled:   false,
		queryLogs: nil,
	})
}

func EnableQueryLog(ctx context.Context) context.Context {
	return context.WithValue(ctx, queryLogKey{}, &queryLogValue{
		enabled:   true,
		queryLogs: nil,
	})
}

func GetQueryLog(ctx context.Context) []QueryLog {
	value := ctx.Value(queryLogKey{})
	if value == nil {
		return nil
	}

	queryLogValue := value.(*queryLogValue)
	if !queryLogValue.enabled {
		return nil
	}

	return queryLogValue.queryLogs
}
