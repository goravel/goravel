package middleware

import (
	contractshttp "github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	httplimit "github.com/goravel/framework/http/limit"

	"goravel/app/http/contracts"
)

// Integration examples showing how to use the optimized throttle middleware in a real application

// CreateAPIThrottle creates a throttle middleware for API endpoints
func CreateAPIThrottle() contractshttp.Middleware {
	// Use the optimized middleware with default service
	return OptimizedThrottle("api", nil)
}

// CreateCustomAPIThrottle creates a throttle middleware with custom response
func CreateCustomAPIThrottle() contractshttp.Middleware {
	// Custom response for rate limited requests
	customResponse := func(ctx contractshttp.Context) {
		ctx.Response().Json(contractshttp.StatusTooManyRequests, map[string]any{
			"error":   "Rate limit exceeded",
			"message": "You have made too many requests. Please try again later.",
			"retry_after": ctx.Response().Writer().Header().Get(HeaderRetryAfter),
		})
	}

	return ThrottleWithCustomResponse("api", customResponse)
}

// CreateTestableThrottle creates a throttle middleware that's easily testable
func CreateTestableThrottle(service contracts.ThrottleService) contractshttp.Middleware {
	return OptimizedThrottle("api", service)
}

// Example of how to configure rate limiting in the RouteServiceProvider
func ExampleRateLimiting() {
	// This would go in your RouteServiceProvider.configureRateLimiting() method
	
	// Configure a basic rate limiter
	facades.RateLimiter().For("api", func(ctx contractshttp.Context) contractshttp.Limit {
		return httplimit.PerMinute(60)
	})

	// Configure a user-specific rate limiter
	facades.RateLimiter().For("user", func(ctx contractshttp.Context) contractshttp.Limit {
		return httplimit.PerMinute(100).By(ctx.Request().Header("X-User-ID"))
	})

	// Configure a stricter rate limiter for guest users
	facades.RateLimiter().For("guest", func(ctx contractshttp.Context) contractshttp.Limit {
		return httplimit.PerMinute(10).By(ctx.Request().Ip())
	})
}

// Example usage in routes
func ExampleRouteUsage() {
	// In your routes/api.go file:
	
	// Apply optimized throttle to API routes
	// facades.Route().Middleware(CreateAPIThrottle()).Group(func(router route.Router) {
	//     router.Get("/users", controllers.UserController{}.Index)
	//     router.Post("/users", controllers.UserController{}.Store)
	// })
	
	// Apply custom throttle with better error messages  
	// facades.Route().Middleware(CreateCustomAPIThrottle()).Group(func(router route.Router) {
	//     router.Get("/premium", controllers.PremiumController{}.Index)
	// })
}