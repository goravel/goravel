package slogmulti

import (
	"context"

	"log/slog"

	"github.com/samber/lo"
)

// Ensure FailoverHandler implements the slog.Handler interface at compile time
var _ slog.Handler = (*FailoverHandler)(nil)

// FailoverHandler implements a high-availability logging pattern.
// It attempts to forward log records to handlers in order until one succeeds.
// This is useful for scenarios where you want primary and backup logging destinations.
//
// @TODO: implement round robin strategy for load balancing across multiple handlers
type FailoverHandler struct {
	// handlers contains the list of slog.Handler instances in priority order
	// The first handler that successfully processes a record will be used
	handlers []slog.Handler
}

// Failover creates a failover handler factory function.
// This function returns a closure that can be used to create failover handlers
// with different sets of handlers.
//
// Example usage:
//
//	handler := slogmulti.Failover()(
//	    primaryHandler,   // First choice
//	    secondaryHandler, // Fallback if primary fails
//	    backupHandler,    // Last resort
//	)
//	logger := slog.New(handler)
//
// Returns:
//
//	A function that creates FailoverHandler instances with the provided handlers
func Failover() func(...slog.Handler) slog.Handler {
	return func(handlers ...slog.Handler) slog.Handler {
		return &FailoverHandler{
			handlers: handlers,
		}
	}
}

// Enabled checks if any of the underlying handlers are enabled for the given log level.
// This method implements the slog.Handler interface requirement.
//
// The handler is considered enabled if at least one of its child handlers
// is enabled for the specified level. This ensures that if any handler
// can process the log, the failover handler will attempt to distribute it.
//
// Args:
//
//	ctx: The context for the logging operation
//	l: The log level to check
//
// Returns:
//
//	true if at least one handler is enabled for the level, false otherwise
func (h *FailoverHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for i := range h.handlers {
		if h.handlers[i].Enabled(ctx, l) {
			return true
		}
	}

	return false
}

// Handle attempts to process a log record using handlers in priority order.
// This method implements the slog.Handler interface requirement.
//
// This implements a "fail-fast" strategy where the first successful handler
// prevents further attempts, making it efficient for high-availability scenarios.
//
// Args:
//
//	ctx: The context for the logging operation
//	r: The log record to process
//
// Returns:
//
//	nil if any handler successfully processed the record, or the last error encountered
func (h *FailoverHandler) Handle(ctx context.Context, r slog.Record) error {
	var err error

	for i := range h.handlers {
		if h.handlers[i].Enabled(ctx, r.Level) {
			err = try(func() error {
				return h.handlers[i].Handle(ctx, r.Clone())
			})
			if err == nil {
				return nil
			}
		}
	}

	return err
}

// WithAttrs creates a new FailoverHandler with additional attributes added to all child handlers.
// This method implements the slog.Handler interface requirement.
//
// The method creates new handler instances for each child handler with the additional
// attributes, ensuring that the attributes are properly propagated to all handlers
// in the failover chain.
//
// Args:
//
//	attrs: The attributes to add to all handlers
//
// Returns:
//
//	A new FailoverHandler with the attributes added to all child handlers
func (h *FailoverHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handers := lo.Map(h.handlers, func(h slog.Handler, _ int) slog.Handler {
		return h.WithAttrs(attrs)
	})
	return Failover()(handers...)
}

// WithGroup creates a new FailoverHandler with a group name applied to all child handlers.
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
//	A new FailoverHandler with the group name applied to all child handlers,
//	or the original handler if the group name is empty
func (h *FailoverHandler) WithGroup(name string) slog.Handler {
	// https://cs.opensource.google/go/x/exp/+/46b07846:slog/handler.go;l=247
	if name == "" {
		return h
	}

	handers := lo.Map(h.handlers, func(h slog.Handler, _ int) slog.Handler {
		return h.WithGroup(name)
	})
	return Failover()(handers...)
}
