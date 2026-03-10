package http

import (
	"context"
)

type Middleware func(Context)

type HandlerFunc func(Context) Response

type ResourceController interface {
	// Index method for controller
	Index(Context) Response
	// Show method for controller
	Show(Context) Response
	// Store method for controller
	Store(Context) Response
	// Update method for controller
	Update(Context) Response
	// Destroy method for controller
	Destroy(Context) Response
}

type Context interface {
	context.Context
	// Context returns the Context
	Context() context.Context
	// WithContext adds a new context to an existing one
	WithContext(ctx context.Context)
	// WithValue add value associated with key in context
	WithValue(key any, value any)
	// Request returns the ContextRequest
	Request() ContextRequest
	// Response returns the ContextResponse
	Response() ContextResponse
}
