package slogmulti

import (
	"context"

	"log/slog"
)

// NewInlineHandler is a shortcut to a handler that implements all methods.
func NewInlineHandler(
	enabledFunc func(ctx context.Context, groups []string, attrs []slog.Attr, level slog.Level) bool,
	handleFunc func(ctx context.Context, groups []string, attrs []slog.Attr, record slog.Record) error,
) slog.Handler {
	if enabledFunc == nil {
		panic("slog-multi: enabledFunc is required")
	}
	if handleFunc == nil {
		panic("slog-multi: handleFunc is required")
	}

	return &InlineHandler{
		groups:      []string{},
		attrs:       []slog.Attr{},
		enabledFunc: enabledFunc,
		handleFunc:  handleFunc,
	}
}

var _ slog.Handler = (*InlineHandler)(nil)

type InlineHandler struct {
	groups      []string
	attrs       []slog.Attr
	enabledFunc func(ctx context.Context, groups []string, attrs []slog.Attr, level slog.Level) bool
	handleFunc  func(ctx context.Context, groups []string, attrs []slog.Attr, record slog.Record) error
}

// Implements slog.Handler
func (h *InlineHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.enabledFunc(ctx, h.groups, h.attrs, level)
}

// Implements slog.Handler
func (h *InlineHandler) Handle(ctx context.Context, record slog.Record) error {
	return h.handleFunc(ctx, h.groups, h.attrs, record)
}

// Implements slog.Handler
func (h *InlineHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := []slog.Attr{}
	newAttrs = append(newAttrs, h.attrs...)
	newAttrs = append(newAttrs, attrs...)

	return &InlineHandler{
		groups:      h.groups,
		attrs:       newAttrs,
		enabledFunc: h.enabledFunc,
		handleFunc:  h.handleFunc,
	}
}

// Implements slog.Handler
func (h *InlineHandler) WithGroup(name string) slog.Handler {
	// https://cs.opensource.google/go/x/exp/+/46b07846:slog/handler.go;l=247
	if name == "" {
		return h
	}

	newGroups := []string{}
	newGroups = append(newGroups, h.groups...)
	newGroups = append(newGroups, name)
	return &InlineHandler{
		groups:      newGroups,
		attrs:       h.attrs,
		enabledFunc: h.enabledFunc,
		handleFunc:  h.handleFunc,
	}
}
