package slogmulti

import (
	"log/slog"
)

// Middleware is a function type that transforms one slog.Handler into another.
// It follows the standard middleware pattern where a function takes a handler
// and returns a new handler that wraps the original with additional functionality.
//
// Middleware functions can be used to:
// - Transform log records (e.g., add timestamps, modify levels)
// - Filter records based on conditions
// - Add context or attributes to records
// - Implement cross-cutting concerns like error recovery or sampling
//
// Example usage:
//
//	  gdprMiddleware := NewGDPRMiddleware()
//	  sink := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{})
//
//		 logger := slog.New(
//			slogmulti.
//				Pipe(gdprMiddleware).
//				// ...
//				Handler(sink),
//		  )
type Middleware func(slog.Handler) slog.Handler
