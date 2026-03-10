package http

import (
	"github.com/goravel/framework/contracts/http"
)

type RateLimiter struct {
	limiters map[string]func(ctx http.Context) []http.Limit
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]func(ctx http.Context) []http.Limit),
	}
}

func (r *RateLimiter) For(name string, callback func(ctx http.Context) http.Limit) {
	r.limiters[name] = func(ctx http.Context) []http.Limit {
		return []http.Limit{callback(ctx)}
	}
}

func (r *RateLimiter) ForWithLimits(name string, callback func(ctx http.Context) []http.Limit) {
	r.limiters[name] = callback
}

func (r *RateLimiter) Limiter(name string) func(ctx http.Context) []http.Limit {
	return r.limiters[name]
}
