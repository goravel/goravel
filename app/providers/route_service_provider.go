package providers

import (
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"

	"goravel/app/http"
	"goravel/routes"
)

type RouteServiceProvider struct {
}

func (receiver *RouteServiceProvider) Register(app foundation.Application) {
}

func (receiver *RouteServiceProvider) Boot(app foundation.Application) {
	// Add HTTP middleware
	facades.Route().GlobalMiddleware(http.Kernel{}.Middleware()...)

	receiver.configureRateLimiting()

	// Add routes
	routes.Web()
	routes.Api()
}

func (receiver *RouteServiceProvider) configureRateLimiting() {
	// Example of configuring rate limiting with the optimized throttle approach
	// You can uncomment these lines to see the improved middleware in action:
	
	// Basic API rate limiting - 60 requests per minute
	// facades.RateLimiter().For("api", func(ctx contractshttp.Context) contractshttp.Limit {
	//     return httplimit.PerMinute(60)
	// })
	
	// User-specific rate limiting - 100 requests per minute per user
	// facades.RateLimiter().For("user", func(ctx contractshttp.Context) contractshttp.Limit {
	//     userID := ctx.Request().Header("X-User-ID", "guest")
	//     return httplimit.PerMinute(100).By(userID)
	// })
	
	// IP-based rate limiting for guests - 10 requests per minute per IP
	// facades.RateLimiter().For("guest", func(ctx contractshttp.Context) contractshttp.Limit {
	//     return httplimit.PerMinute(10).By(ctx.Request().Ip())
	// })
}
