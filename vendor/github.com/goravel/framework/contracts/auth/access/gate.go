package access

import "context"

type Gate interface {
	// WithContext returns a new Gate instance with the given context.
	WithContext(ctx context.Context) Gate
	// Allows determines if the given ability should be granted for the current user.
	Allows(ability string, arguments map[string]any) bool
	// Denies determines if the given ability should be denied for the current user.
	Denies(ability string, arguments map[string]any) bool
	// Inspect the given ability against the current user.
	Inspect(ability string, arguments map[string]any) Response
	// Define a new ability.
	Define(ability string, callback func(ctx context.Context, arguments map[string]any) Response)
	// Any one of the given abilities should be granted for the current user.
	Any(abilities []string, arguments map[string]any) bool
	// None of the given abilities should be granted for the current user.
	None(abilities []string, arguments map[string]any) bool
	// Before register a callback to run before all Gate checks.
	Before(callback func(ctx context.Context, ability string, arguments map[string]any) Response)
	// After register a callback to run after all Gate checks.
	After(callback func(ctx context.Context, ability string, arguments map[string]any, result Response) Response)
}

type Response interface {
	// Allowed to determine if the response was allowed.
	Allowed() bool
	// Message to get the response message.
	Message() string
}
