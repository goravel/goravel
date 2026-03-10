package slogmulti

import (
	"log/slog"
)

// PipeBuilder provides a fluent API for building middleware chains.
// It allows you to compose multiple middleware functions that will be applied
// to log records in the order they are added (last-in, first-out).
type PipeBuilder struct {
	// middlewares contains the list of middleware functions to be applied
	// The middlewares are applied in reverse order (LIFO) when building the final handler
	middlewares []Middleware
}

// Pipe creates a new PipeBuilder with the provided middleware functions.
// This function is the entry point for building middleware chains.
//
// Middleware functions are applied in reverse order (last-in, first-out),
// which means the last middleware added will be the first one applied to incoming records.
// This allows for intuitive composition where you can think of the chain as
// "transform A, then transform B, then send to handler".
//
// Example usage:
//
//	handler := slogmulti.Pipe(
//	    RewriteLevel(slog.LevelWarn, slog.LevelInfo),
//	    RewriteMessage("prefix: %s"),
//	    RedactPII(),
//	).Handler(finalHandler)
//
// Args:
//
//	middlewares: Variable number of middleware functions to chain together
//
// Returns:
//
//	A new PipeBuilder instance ready for further configuration
func Pipe(middlewares ...Middleware) *PipeBuilder {
	return &PipeBuilder{middlewares: middlewares}
}

// Pipe adds an additional middleware to the chain.
// This method provides a fluent API for building middleware chains incrementally.
//
// Args:
//
//	middleware: The middleware function to add to the chain
//
// Returns:
//
//	The PipeBuilder instance for method chaining
func (h *PipeBuilder) Pipe(middleware Middleware) *PipeBuilder {
	h.middlewares = append(h.middlewares, middleware)
	return h
}

// Handler creates a slog.Handler by applying all middleware to the provided handler.
// This method finalizes the middleware chain and returns a handler that can be used with slog.New().
//
// This LIFO approach ensures that the middleware chain is applied in the intuitive order:
// the first middleware in the chain is applied first to incoming records.
//
// Args:
//
//	handler: The final slog.Handler that will receive the transformed records
//
// Returns:
//
//	A slog.Handler that applies all middleware transformations before forwarding to the final handler
func (h *PipeBuilder) Handler(handler slog.Handler) slog.Handler {
	for len(h.middlewares) > 0 {
		middleware := h.middlewares[len(h.middlewares)-1]
		h.middlewares = h.middlewares[0 : len(h.middlewares)-1]
		handler = middleware(handler)
	}

	return handler
}
