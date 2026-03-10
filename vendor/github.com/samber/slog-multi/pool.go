package slogmulti

import (
	"context"
	"math/rand"
	"time"

	"log/slog"

	"github.com/samber/lo"
)

// Ensure PoolHandler implements the slog.Handler interface at compile time
var _ slog.Handler = (*PoolHandler)(nil)

// PoolHandler implements a load balancing strategy for logging handlers.
// It distributes log records across multiple handlers using a round-robin approach
// with randomization to ensure even distribution and prevent hot-spotting.
type PoolHandler struct {
	// randSource provides a thread-safe random number generator for load balancing
	randSource rand.Source
	// handlers contains the list of slog.Handler instances to distribute records across
	handlers []slog.Handler
}

// Pool creates a load balancing handler factory function.
// This function returns a closure that can be used to create pool handlers
// with different sets of handlers for load balancing.
//
// The pool uses a round-robin strategy with randomization to distribute
// log records evenly across all available handlers. This is useful for:
// - Increasing logging throughput by parallelizing handler operations
// - Providing redundancy by having multiple handlers process the same records
// - Load balancing across multiple logging destinations
//
// Example usage:
//
//	handler := slogmulti.Pool()(
//	    handler1, // Will receive ~33% of records
//	    handler2, // Will receive ~33% of records
//	    handler3, // Will receive ~33% of records
//	)
//	logger := slog.New(handler)
//
// Returns:
//
//	A function that creates PoolHandler instances with the provided handlers
func Pool() func(...slog.Handler) slog.Handler {
	return func(handlers ...slog.Handler) slog.Handler {
		return &PoolHandler{
			randSource: rand.NewSource(time.Now().UnixNano()),
			handlers:   handlers,
		}
	}
}

// Enabled checks if any of the underlying handlers are enabled for the given log level.
// This method implements the slog.Handler interface requirement.
//
// The handler is considered enabled if at least one of its child handlers
// is enabled for the specified level. This ensures that if any handler
// can process the log, the pool handler will attempt to distribute it.
//
// Args:
//
//	ctx: The context for the logging operation
//	l: The log level to check
//
// Returns:
//
//	true if at least one handler is enabled for the level, false otherwise
func (h *PoolHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for i := range h.handlers {
		if h.handlers[i].Enabled(ctx, l) {
			return true
		}
	}

	return false
}

// Handle distributes a log record to a handler selected using round-robin with randomization.
// This method implements the slog.Handler interface requirement.
//
// This approach ensures even distribution of load while providing fault tolerance
// through the failover behavior when a handler is unavailable.
//
// Args:
//
//	ctx: The context for the logging operation
//	r: The log record to distribute
//
// Returns:
//
//	nil if any handler successfully processed the record, or the last error encountered
func (h *PoolHandler) Handle(ctx context.Context, r slog.Record) error {
	if len(h.handlers) == 0 {
		return nil
	}

	// round robin with randomization
	rand := h.randSource.Int63() % int64(len(h.handlers))
	handlers := append(h.handlers[rand:], h.handlers[:rand]...)

	var err error

	for i := range handlers {
		if handlers[i].Enabled(ctx, r.Level) {
			err = try(func() error {
				return handlers[i].Handle(ctx, r.Clone())
			})
			if err == nil {
				return nil
			}
		}
	}

	return err
}

// WithAttrs creates a new PoolHandler with additional attributes added to all child handlers.
// This method implements the slog.Handler interface requirement.
//
// The method creates new handler instances for each child handler with the additional
// attributes, ensuring that the attributes are properly propagated to all handlers
// in the pool.
//
// Args:
//
//	attrs: The attributes to add to all handlers
//
// Returns:
//
//	A new PoolHandler with the attributes added to all child handlers
func (h *PoolHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handers := lo.Map(h.handlers, func(h slog.Handler, _ int) slog.Handler {
		return h.WithAttrs(attrs)
	})
	return Pool()(handers...)
}

// WithGroup creates a new PoolHandler with a group name applied to all child handlers.
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
//	A new PoolHandler with the group name applied to all child handlers,
//	or the original handler if the group name is empty
func (h *PoolHandler) WithGroup(name string) slog.Handler {
	// https://cs.opensource.google/go/x/exp/+/46b07846:slog/handler.go;l=247
	if name == "" {
		return h
	}

	handers := lo.Map(h.handlers, func(h slog.Handler, _ int) slog.Handler {
		return h.WithGroup(name)
	})
	return Pool()(handers...)
}
