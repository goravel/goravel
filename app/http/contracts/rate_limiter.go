package contracts

import (
	contractshttp "github.com/goravel/framework/contracts/http"
)

// RateLimitRule defines the interface for rate limiting rules that can be tested easily
type RateLimitRule interface {
	// GetStore returns the store instance for rate limiting
	GetStore() contractshttp.Store
	// GetKey returns the rate limit signature key
	GetKey() string
	// GetResponseCallback returns the response generator callback
	GetResponseCallback() func(ctx contractshttp.Context)
	// SetKey sets the rate limit signature key
	SetKey(key string) RateLimitRule
	// SetResponseCallback sets the response generator callback
	SetResponseCallback(callback func(ctx contractshttp.Context)) RateLimitRule
}

// ThrottleResponse represents the result of a throttle check
type ThrottleResponse struct {
	Allowed    bool
	Tokens     uint64
	Remaining  uint64
	Reset      uint64
	RetryAfter int
}

// ThrottleService defines the interface for throttle operations
type ThrottleService interface {
	// Check performs a throttle check for the given rule and key
	Check(ctx contractshttp.Context, rule RateLimitRule, key string) (ThrottleResponse, error)
	// BuildKey builds a throttle key from context and rule
	BuildKey(ctx contractshttp.Context, rule RateLimitRule, name string, index int) string
}