package middleware

import (
	"strconv"

	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/http"
	httplimit "github.com/goravel/framework/http/limit"
	"github.com/goravel/framework/support/carbon"

	"goravel/app/http/contracts"
)

const (
	// HeaderRateLimitLimit, HeaderRateLimitRemaining, and HeaderRateLimitReset
	// are the recommended return header values from IETF on rate limiting. Reset
	// is in UTC time.
	HeaderRateLimitLimit     = "X-RateLimit-Limit"
	HeaderRateLimitRemaining = "X-RateLimit-Remaining"
	HeaderRateLimitReset     = "X-RateLimit-Reset"

	// HeaderRetryAfter is the header used to indicate when a client should retry
	// requests (when the rate limit expires), in UTC time.
	HeaderRetryAfter = "Retry-After"
)

// OptimizedThrottle creates an optimized throttle middleware that uses contracts for better testability
func OptimizedThrottle(name string, service contracts.ThrottleService) contractshttp.Middleware {
	if service == nil {
		service = NewThrottleService()
	}

	return func(ctx contractshttp.Context) {
		if limiter := http.RateLimiterFacade.Limiter(name); limiter != nil {
			if limits := limiter(ctx); len(limits) > 0 {
				for index, limit := range limits {
					// Convert framework limit to our contract
					rule := ConvertLimitToRule(limit)
					if rule == nil {
						continue
					}

					key := service.BuildKey(ctx, rule, name, index)
					response, err := service.Check(ctx, rule, key)
					if err != nil {
						http.LogFacade.Error(errors.HttpRateLimitFailedToCheckThrottle.Args(err))
						break
					}

					resetTime := carbon.FromTimestampNano(int64(response.Reset)).SetTimezone(carbon.UTC)
					ctx.Response().Header(HeaderRateLimitLimit, strconv.FormatUint(response.Tokens, 10))
					ctx.Response().Header(HeaderRateLimitRemaining, strconv.FormatUint(response.Remaining, 10))

					if !response.Allowed {
						ctx.Response().Header(HeaderRateLimitReset, strconv.Itoa(int(resetTime.Timestamp())))
						ctx.Response().Header(HeaderRetryAfter, strconv.Itoa(response.RetryAfter))
						rule.GetResponseCallback()(ctx)
						return
					}
				}
			}
		}

		ctx.Request().Next()
	}
}

// ConvertLimitToRule converts a framework limit to our contract
func ConvertLimitToRule(limit contractshttp.Limit) contracts.RateLimitRule {
	// Try to convert to concrete type to access properties
	if instance, ok := limit.(*httplimit.Limit); ok {
		rule := NewRateLimitRule(instance.Store)
		rule.SetKey(instance.Key)
		if instance.ResponseCallback != nil {
			rule.SetResponseCallback(instance.ResponseCallback)
		}
		return rule
	}
	return nil
}