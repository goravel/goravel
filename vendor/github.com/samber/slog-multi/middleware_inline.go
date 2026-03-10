package slogmulti

import (
	"context"

	"log/slog"
)

// NewInlineMiddleware is a shortcut to a middleware that implements all methods.
func NewInlineMiddleware(
	enabledFunc func(ctx context.Context, level slog.Level, next func(context.Context, slog.Level) bool) bool,
	handleFunc func(ctx context.Context, record slog.Record, next func(context.Context, slog.Record) error) error,
	withAttrsFunc func(attrs []slog.Attr, next func([]slog.Attr) slog.Handler) slog.Handler,
	withGroupFunc func(name string, next func(string) slog.Handler) slog.Handler,
) Middleware {
	return func(next slog.Handler) slog.Handler {
		if next == nil {
			panic("slog-multi: next is required")
		}
		if enabledFunc == nil {
			panic("slog-multi: enabledFunc is required")
		}
		if handleFunc == nil {
			panic("slog-multi: handleFunc is required")
		}
		if withAttrsFunc == nil {
			panic("slog-multi: withAttrsFunc is required")
		}
		if withGroupFunc == nil {
			panic("slog-multi: withGroupFunc is required")
		}

		return &InlineMiddleware{
			next:          next,
			enabledFunc:   enabledFunc,
			handleFunc:    handleFunc,
			withAttrsFunc: withAttrsFunc,
			withGroupFunc: withGroupFunc,
		}
	}
}

var _ slog.Handler = (*InlineMiddleware)(nil)

type InlineMiddleware struct {
	next          slog.Handler
	enabledFunc   func(ctx context.Context, level slog.Level, next func(context.Context, slog.Level) bool) bool
	handleFunc    func(ctx context.Context, record slog.Record, next func(context.Context, slog.Record) error) error
	withAttrsFunc func(attrs []slog.Attr, next func([]slog.Attr) slog.Handler) slog.Handler
	withGroupFunc func(name string, next func(string) slog.Handler) slog.Handler
}

// Implements slog.Handler
func (h *InlineMiddleware) Enabled(ctx context.Context, level slog.Level) bool {
	return h.enabledFunc(ctx, level, h.next.Enabled)
}

// Implements slog.Handler
func (h *InlineMiddleware) Handle(ctx context.Context, record slog.Record) error {
	return h.handleFunc(ctx, record, h.next.Handle)
}

// Implements slog.Handler
func (h *InlineMiddleware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewInlineMiddleware(
		h.enabledFunc,
		h.handleFunc,
		h.withAttrsFunc,
		h.withGroupFunc,
	)(h.withAttrsFunc(attrs, h.next.WithAttrs))
}

// Implements slog.Handler
func (h *InlineMiddleware) WithGroup(name string) slog.Handler {
	// https://cs.opensource.google/go/x/exp/+/46b07846:slog/handler.go;l=247
	if name == "" {
		return h
	}

	return NewInlineMiddleware(
		h.enabledFunc,
		h.handleFunc,
		h.withAttrsFunc,
		h.withGroupFunc,
	)(h.withGroupFunc(name, h.next.WithGroup))
}
