package configuration

import "github.com/goravel/framework/contracts/http"

type Middleware interface {
	// Append adds middleware to the end of the middleware stack.
	Append(middleware ...http.Middleware) Middleware
	// GetGlobalMiddleware returns the global middleware stack.
	GetGlobalMiddleware() []http.Middleware
	// GetRecover returns the recover function for handling panics.
	GetRecover() func(ctx http.Context, err any)
	// Prepend adds middleware to the beginning of the middleware stack.
	Prepend(middleware ...http.Middleware) Middleware
	// Recover sets the recover function for handling panics.
	Recover(fn func(ctx http.Context, err any)) Middleware
	// Use sets the middleware stack to the provided middleware.
	Use(middleware ...http.Middleware) Middleware
}
