package slogmulti

import (
	"context"
	"errors"
	"log/slog"
	"slices"

	"github.com/samber/lo"
)

// Ensure FanoutHandler implements the slog.Handler interface at compile time
var _ slog.Handler = (*FanoutHandler)(nil)

// FanoutHandler distributes log records to multiple slog.Handler instances in parallel.
// It implements the slog.Handler interface and forwards all logging operations to all
// registered handlers that are enabled for the given log level.
type FanoutHandler struct {
	// handlers contains the list of slog.Handler instances to which log records will be distributed
	handlers []slog.Handler
}

// Fanout creates a new FanoutHandler that distributes records to multiple slog.Handler instances.
// If exactly one handler is provided, it returns that handler unmodified.
// If you pass a FanoutHandler as an argument, its handlers are flattened into the new FanoutHandler.
// This function is the primary entry point for creating a multi-handler setup.
//
// Example usage:
//
//	handler := slogmulti.Fanout(
//	    slog.NewJSONHandler(os.Stdout, nil),
//	    slogdatadog.NewDatadogHandler(...),
//	)
//	logger := slog.New(handler)
//
// Args:
//
//	handlers: Variable number of slog.Handler instances to distribute logs to
//
// Returns:
//
//	A slog.Handler that forwards all operations to the provided handlers
func Fanout(handlers ...slog.Handler) slog.Handler {
	var flat []slog.Handler
	for _, handler := range handlers {
		if fan, ok := handler.(*FanoutHandler); ok {
			flat = append(flat, fan.handlers...)
		} else {
			flat = append(flat, handler)
		}
	}

	if len(flat) == 1 {
		return flat[0]
	}
	return &FanoutHandler{
		handlers: flat,
	}
}

// Enabled checks if any of the underlying handlers are enabled for the given log level.
// This method implements the slog.Handler interface requirement.
//
// The handler is considered enabled if at least one of its child handlers
// is enabled for the specified level. This ensures that if any handler
// can process the log, the fanout handler will attempt to distribute it.
//
// Args:
//
//	ctx: The context for the logging operation
//	l: The log level to check
//
// Returns:
//
//	true if at least one handler is enabled for the level, false otherwise
func (h *FanoutHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for i := range h.handlers {
		if h.handlers[i].Enabled(ctx, l) {
			return true
		}
	}

	return false
}

// Handle distributes a log record to all enabled handlers.
// This method implements the slog.Handler interface requirement.
//
// The method:
// 1. Iterates through all registered handlers
// 2. Checks if each handler is enabled for the record's level
// 3. For enabled handlers, calls their Handle method with a cloned record
// 4. Collects any errors that occur during handling
// 5. Returns a combined error if any handlers failed
//
// Note: Each handler receives a cloned record to prevent interference between handlers.
// This ensures that one handler cannot modify the record for other handlers.
//
// Args:
//
//	ctx: The context for the logging operation
//	r: The log record to distribute
//
// Returns:
//
//	An error if any handler failed to process the record, nil otherwise
func (h *FanoutHandler) Handle(ctx context.Context, r slog.Record) error {
	var errs []error
	for i := range h.handlers {
		if h.handlers[i].Enabled(ctx, r.Level) {
			err := try(func() error {
				return h.handlers[i].Handle(ctx, r.Clone())
			})
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	// If errs is empty, or contains only nil errors, this returns nil
	return errors.Join(errs...)
}

// WithAttrs creates a new FanoutHandler with additional attributes added to all child handlers.
// This method implements the slog.Handler interface requirement.
//
// The method creates new handler instances for each child handler with the additional
// attributes, ensuring that the attributes are properly propagated to all handlers
// in the fanout chain.
//
// Args:
//
//	attrs: The attributes to add to all handlers
//
// Returns:
//
//	A new FanoutHandler with the attributes added to all child handlers
func (h *FanoutHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := lo.Map(h.handlers, func(h slog.Handler, _ int) slog.Handler {
		return h.WithAttrs(slices.Clone(attrs))
	})
	return Fanout(handlers...)
}

// WithGroup creates a new FanoutHandler with a group name applied to all child handlers.
// This method implements the slog.Handler interface requirement.
//
// The method follows the same pattern as the standard slog implementation:
// - If the group name is empty, returns the original handler unchanged
// - Otherwise, creates new handler instances for each child handler with the group name
//
// Args:
//
//	name: The group name to apply to all handlers
//
// Returns:
//
//	A new FanoutHandler with the group name applied to all child handlers,
//	or the original handler if the group name is empty
func (h *FanoutHandler) WithGroup(name string) slog.Handler {
	// https://cs.opensource.google/go/x/exp/+/46b07846:slog/handler.go;l=247
	if name == "" {
		return h
	}

	handlers := lo.Map(h.handlers, func(h slog.Handler, _ int) slog.Handler {
		return h.WithGroup(name)
	})
	return Fanout(handlers...)
}
