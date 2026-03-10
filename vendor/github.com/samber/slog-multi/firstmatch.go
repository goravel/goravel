package slogmulti

import (
	"context"
	"log/slog"
	"slices"

	"github.com/samber/lo"
)

// Ensure FirstMatchHandler implements the slog.Handler interface at compile time
var _ slog.Handler = (*FirstMatchHandler)(nil)

type FirstMatchHandler struct {
	handlers []*RoutableHandler
}

func FirstMatch(handlers ...*RoutableHandler) *FirstMatchHandler {
	return &FirstMatchHandler{handlers: lo.Map(handlers, func(h *RoutableHandler, _ int) *RoutableHandler {
		return &RoutableHandler{
			predicates:     h.predicates,
			handler:        h.handler,
			groups:         slices.Clone(h.groups),
			attrs:          slices.Clone(h.attrs),
			skipPredicates: true, // prevent double matching
		}
	})}
}

// Enabled checks if any of the underlying handlers are enabled for the given log level.
// This method implements the slog.Handler interface requirement.
// See FanoutHandler.WithAttrs for details.
func (h *FirstMatchHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for i := range h.handlers {
		if h.handlers[i].Enabled(ctx, l) {
			return true
		}
	}

	return false
}

// Handle distributes a log record to the first matching handler.
// This method implements the slog.Handler interface requirement.
//
// The method:
// 1. Iterates through each child handler.
// 2. Checks if the handler's predicates match the record.
// 3. If a match is found, it checks if the handler is enabled for the record's level.
// 4. If enabled, it forwards the record to that handler and returns.
// 5. If no handlers match, it returns nil.
func (h *FirstMatchHandler) Handle(ctx context.Context, r slog.Record) error {
	for i := range h.handlers {
		record, ok := h.handlers[i].isMatch(ctx, r)
		if ok {
			if h.handlers[i].Enabled(ctx, record.Level) {
				return try(func() error {
					return h.handlers[i].Handle(ctx, r)
				})
			}

			return nil // Handler matched but is not enabled; do not proceed further
		}
	}

	return nil
}

// WithAttrs creates a new FirstMatchHandler with additional attributes added to all child handlers.
// This method implements the slog.Handler interface requirement.
// See FanoutHandler.WithAttrs for details.
func (h *FirstMatchHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := lo.Map(h.handlers, func(h *RoutableHandler, _ int) *RoutableHandler {
		return h.WithAttrs(slices.Clone(attrs)).(*RoutableHandler)
	})
	return newFirstMatch(handlers...)
}

// WithGroup creates a new FirstMatchHandler with a group name applied to all child handlers.
// This method implements the slog.Handler interface requirement.
// See FanoutHandler.WithGroup for details.
func (h *FirstMatchHandler) WithGroup(name string) slog.Handler {
	// https://cs.opensource.google/go/x/exp/+/46b07846:slog/handler.go;l=247
	if name == "" {
		return h
	}

	handlers := lo.Map(h.handlers, func(h *RoutableHandler, _ int) *RoutableHandler {
		return h.WithGroup(name).(*RoutableHandler)
	})
	return newFirstMatch(handlers...)
}

func newFirstMatch(handlers ...*RoutableHandler) *FirstMatchHandler {
	return &FirstMatchHandler{handlers: handlers}
}
