package middleware

import (
	contractshttp "github.com/goravel/framework/contracts/http"
	httplimit "github.com/goravel/framework/http/limit"

	"goravel/app/http/contracts"
)

// Example usage of the optimized throttle middleware

// ThrottleWithCustomStore demonstrates how to create a throttle middleware with a custom store
func ThrottleWithCustomStore(name string, store contractshttp.Store) contractshttp.Middleware {
	service := NewThrottleService()
	return OptimizedThrottle(name, service)
}

// ThrottleWithCustomResponse demonstrates how to create a throttle middleware with custom response handling
func ThrottleWithCustomResponse(name string, responseHandler func(ctx contractshttp.Context)) contractshttp.Middleware {
	return func(ctx contractshttp.Context) {
		// Create a custom service that can modify the response callback
		service := &CustomThrottleService{
			ThrottleService:   NewThrottleService(),
			responseHandler:   responseHandler,
		}
		
		middleware := OptimizedThrottle(name, service)
		middleware(ctx)
	}
}

// CustomThrottleService extends the base throttle service with custom response handling
type CustomThrottleService struct {
	*ThrottleService
	responseHandler func(ctx contractshttp.Context)
}

// Check performs a throttle check and applies custom response handling if needed
func (s *CustomThrottleService) Check(ctx contractshttp.Context, rule contracts.RateLimitRule, key string) (contracts.ThrottleResponse, error) {
	response, err := s.ThrottleService.Check(ctx, rule, key)
	if err != nil {
		return response, err
	}

	// If request is not allowed and we have a custom response handler, set it on the rule
	if !response.Allowed && s.responseHandler != nil {
		rule.SetResponseCallback(s.responseHandler)
	}

	return response, nil
}

// TestableThrottleRule creates a rate limit rule that's easy to test
func TestableThrottleRule(store contractshttp.Store, key string, responseCallback func(ctx contractshttp.Context)) contracts.RateLimitRule {
	rule := NewRateLimitRule(store)
	if key != "" {
		rule.SetKey(key)
	}
	if responseCallback != nil {
		rule.SetResponseCallback(responseCallback)
	}
	return rule
}

// PerMinuteTestable creates a testable per-minute rate limit
func PerMinuteTestable(maxAttempts int, store contractshttp.Store) contracts.RateLimitRule {
	// If no store provided, use the default one like the framework does
	if store == nil {
		limit := httplimit.PerMinute(maxAttempts)
		return ConvertLimitToRule(limit)
	}
	
	return NewRateLimitRule(store)
}