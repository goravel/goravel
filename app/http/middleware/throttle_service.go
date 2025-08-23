package middleware

import (
	"fmt"

	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/support/carbon"

	"goravel/app/http/contracts"
)

// ThrottleService implements the throttle service interface
type ThrottleService struct{}

// NewThrottleService creates a new throttle service
func NewThrottleService() *ThrottleService {
	return &ThrottleService{}
}

// Check performs a throttle check for the given rule and key
func (s *ThrottleService) Check(ctx contractshttp.Context, rule contracts.RateLimitRule, key string) (contracts.ThrottleResponse, error) {
	store := rule.GetStore()
	if store == nil {
		return contracts.ThrottleResponse{Allowed: true}, nil
	}

	tokens, remaining, reset, ok, err := store.Take(ctx, key)
	if err != nil {
		return contracts.ThrottleResponse{}, err
	}

	resetTime := carbon.FromTimestampNano(int64(reset)).SetTimezone(carbon.UTC)
	retryAfter := carbon.Now().DiffInSeconds(resetTime)

	return contracts.ThrottleResponse{
		Allowed:    ok,
		Tokens:     tokens,
		Remaining:  remaining,
		Reset:      reset,
		RetryAfter: int(retryAfter),
	}, nil
}

// BuildKey builds a throttle key from context and rule
func (s *ThrottleService) BuildKey(ctx contractshttp.Context, rule contracts.RateLimitRule, name string, index int) string {
	// if no key is set, use the path and ip address as the default key
	key := rule.GetKey()
	if len(key) == 0 && ctx.Request() != nil {
		return fmt.Sprintf("throttle:%s:%d:%s:%s", name, index, ctx.Request().Ip(), ctx.Request().Path())
	}

	return fmt.Sprintf("throttle:%s:%d:%s", name, index, key)
}