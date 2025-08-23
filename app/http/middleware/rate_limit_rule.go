package middleware

import (
	contractshttp "github.com/goravel/framework/contracts/http"

	"goravel/app/http/contracts"
)

// RateLimitRule implements the improved rate limiting rule interface
type RateLimitRule struct {
	store            contractshttp.Store
	key              string
	responseCallback func(ctx contractshttp.Context)
}

// NewRateLimitRule creates a new rate limiting rule
func NewRateLimitRule(store contractshttp.Store) *RateLimitRule {
	return &RateLimitRule{
		store: store,
		responseCallback: func(ctx contractshttp.Context) {
			ctx.Request().Abort(contractshttp.StatusTooManyRequests)
		},
	}
}

// GetStore returns the store instance for rate limiting
func (r *RateLimitRule) GetStore() contractshttp.Store {
	return r.store
}

// GetKey returns the rate limit signature key
func (r *RateLimitRule) GetKey() string {
	return r.key
}

// GetResponseCallback returns the response generator callback
func (r *RateLimitRule) GetResponseCallback() func(ctx contractshttp.Context) {
	return r.responseCallback
}

// SetKey sets the rate limit signature key
func (r *RateLimitRule) SetKey(key string) contracts.RateLimitRule {
	r.key = key
	return r
}

// SetResponseCallback sets the response generator callback
func (r *RateLimitRule) SetResponseCallback(callback func(ctx contractshttp.Context)) contracts.RateLimitRule {
	r.responseCallback = callback
	return r
}