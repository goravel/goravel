package slogmulti

import (
	"context"

	"log/slog"
)

// NewWithAttrsInlineMiddleware is a shortcut to a middleware that implements only the `WithAttrs` method.
func NewWithAttrsInlineMiddleware(withAttrsFunc func(attrs []slog.Attr, next func([]slog.Attr) slog.Handler) slog.Handler) Middleware {
	return func(next slog.Handler) slog.Handler {
		if withAttrsFunc == nil {
			panic("slog-multi: withAttrsFunc is required")
		}
		if next == nil {
			panic("slog-multi: next is required")
		}

		return &WithAttrsInlineMiddleware{
			next:          next,
			withAttrsFunc: withAttrsFunc,
		}
	}
}

var _ slog.Handler = (*WithAttrsInlineMiddleware)(nil)

type WithAttrsInlineMiddleware struct {
	next          slog.Handler
	withAttrsFunc func(attrs []slog.Attr, next func([]slog.Attr) slog.Handler) slog.Handler
}

// Implements slog.Handler
func (h *WithAttrsInlineMiddleware) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

// Implements slog.Handler
func (h *WithAttrsInlineMiddleware) Handle(ctx context.Context, record slog.Record) error {
	return h.next.Handle(ctx, record)
}

// Implements slog.Handler
func (h *WithAttrsInlineMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewWithAttrsInlineMiddleware(h.withAttrsFunc)(h.withAttrsFunc(attrs, h.next.WithAttrs))
}

// Implements slog.Handler
func (h *WithAttrsInlineMiddleware) WithGroup(name string) slog.Handler {
	// https://cs.opensource.google/go/x/exp/+/46b07846:slog/handler.go;l=247
	if name == "" {
		return h
	}

	return NewWithAttrsInlineMiddleware(h.withAttrsFunc)(h.next.WithGroup(name))
}
