package slogmulti

import (
	"context"
	"fmt"
	"log/slog"
)

// RecoveryFunc is a callback function that handles errors and panics from logging handlers.
// It receives the context, the log record that caused the error, and the error itself.
// This function can be used to log the error, send alerts, or perform any other
// error handling logic without affecting the main application flow.
type RecoveryFunc func(ctx context.Context, record slog.Record, err error)

// Ensure HandlerErrorRecovery implements the slog.Handler interface at compile time
var _ slog.Handler = (*HandlerErrorRecovery)(nil)

// HandlerErrorRecovery wraps a slog.Handler to provide panic and error recovery.
// It catches both panics and errors from the underlying handler and calls
// a recovery function to handle them gracefully.
type HandlerErrorRecovery struct {
	// recovery is the function called when an error or panic occurs
	recovery RecoveryFunc
	// handler is the underlying slog.Handler that this recovery wrapper protects
	handler slog.Handler
}

// RecoverHandlerError creates a middleware that adds error recovery to a slog.Handler.
// This function returns a closure that can be used to wrap handlers with recovery logic.
//
// The recovery handler provides fault tolerance by:
// 1. Catching panics from the underlying handler
// 2. Catching errors returned by the underlying handler
// 3. Calling the recovery function with the error details
// 4. Propagating the original error to maintain logging semantics
//
// Example usage:
//
//	recovery := slogmulti.RecoverHandlerError(func(ctx context.Context, record slog.Record, err error) {
//	    fmt.Printf("Logging error: %v\n", err)
//	})
//	safeHandler := recovery(riskyHandler)
//	logger := slog.New(safeHandler)
//
// Args:
//
//	recovery: The function to call when an error or panic occurs
//
// Returns:
//
//	A function that wraps handlers with recovery logic
func RecoverHandlerError(recovery RecoveryFunc) func(slog.Handler) slog.Handler {
	return func(handler slog.Handler) slog.Handler {
		return &HandlerErrorRecovery{
			recovery: recovery,
			handler:  handler,
		}
	}
}

// Enabled checks if the underlying handler is enabled for the given log level.
// This method implements the slog.Handler interface requirement.
//
// Args:
//
//	ctx: The context for the logging operation
//	l: The log level to check
//
// Returns:
//
//	true if the underlying handler is enabled for the level, false otherwise
func (h *HandlerErrorRecovery) Enabled(ctx context.Context, l slog.Level) bool {
	return h.handler.Enabled(ctx, l)
}

// Handle processes a log record with error recovery.
// This method implements the slog.Handler interface requirement.
//
// This ensures that logging errors don't crash the application while still
// allowing the error to be handled appropriately by the calling code.
//
// Args:
//
//	ctx: The context for the logging operation
//	record: The log record to process
//
// Returns:
//
//	The error from the underlying handler (never nil if an error occurred)
func (h *HandlerErrorRecovery) Handle(ctx context.Context, record slog.Record) error {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				h.recovery(ctx, record, e)
			} else {
				h.recovery(ctx, record, fmt.Errorf("%+v", r))
			}
		}
	}()

	err := h.handler.Handle(ctx, record)
	if err != nil {
		h.recovery(ctx, record, err)
	}

	// propagate error
	return err
}

// WithAttrs creates a new HandlerErrorRecovery with additional attributes.
// This method implements the slog.Handler interface requirement.
//
// Args:
//
//	attrs: The attributes to add to the underlying handler
//
// Returns:
//
//	A new HandlerErrorRecovery with the additional attributes
func (h *HandlerErrorRecovery) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerErrorRecovery{
		recovery: h.recovery,
		handler:  h.handler.WithAttrs(attrs),
	}
}

// WithGroup creates a new HandlerErrorRecovery with a group name.
// This method implements the slog.Handler interface requirement.
//
// The method follows the same pattern as the standard slog implementation:
// - If the group name is empty, returns the original handler unchanged
// - Otherwise, creates a new handler with the group name applied to the underlying handler
//
// Args:
//
//	name: The group name to apply to the underlying handler
//
// Returns:
//
//	A new HandlerErrorRecovery with the group name, or the original handler if the name is empty
func (h *HandlerErrorRecovery) WithGroup(name string) slog.Handler {
	// https://cs.opensource.google/go/x/exp/+/46b07846:slog/handler.go;l=247
	if name == "" {
		return h
	}

	return &HandlerErrorRecovery{
		recovery: h.recovery,
		handler:  h.handler.WithGroup(name),
	}
}
