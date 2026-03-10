package slogmulti

import (
	"context"

	"log/slog"
)

// NewHandleInlineHandler is a shortcut to a middleware that implements only the `Handle` method.
func NewHandleInlineHandler(handleFunc func(ctx context.Context, groups []string, attrs []slog.Attr, record slog.Record) error) slog.Handler {
	return &HandleInlineHandler{
		groups:     []string{},
		attrs:      []slog.Attr{},
		handleFunc: handleFunc,
	}
}

var _ slog.Handler = (*HandleInlineHandler)(nil)

type HandleInlineHandler struct {
	groups     []string
	attrs      []slog.Attr
	handleFunc func(ctx context.Context, groups []string, attrs []slog.Attr, record slog.Record) error
}

// Implements slog.Handler
func (h *HandleInlineHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

// Implements slog.Handler
func (h *HandleInlineHandler) Handle(ctx context.Context, record slog.Record) error {
	return h.handleFunc(ctx, h.groups, h.attrs, record)
}

// Implements slog.Handler
func (h *HandleInlineHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := []slog.Attr{}
	newAttrs = append(newAttrs, h.attrs...)
	newAttrs = append(newAttrs, attrs...)

	return &HandleInlineHandler{
		groups:     h.groups,
		attrs:      newAttrs,
		handleFunc: h.handleFunc,
	}
}

// Implements slog.Handler
func (h *HandleInlineHandler) WithGroup(name string) slog.Handler {
	// https://cs.opensource.google/go/x/exp/+/46b07846:slog/handler.go;l=247
	if name == "" {
		return h
	}

	newGroups := []string{}
	newGroups = append(newGroups, h.groups...)
	newGroups = append(newGroups, name)
	return &HandleInlineHandler{
		groups:     newGroups,
		attrs:      h.attrs,
		handleFunc: h.handleFunc,
	}
}
