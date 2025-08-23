package providers

import (
	"github.com/goravel/framework/contracts/foundation"

	"goravel/routes"
)

type RouteServiceProvider struct {
}

func (receiver *RouteServiceProvider) Register(app foundation.Application) {
}

func (receiver *RouteServiceProvider) Boot(app foundation.Application) {
	receiver.configureRateLimiting()

	// Add routes
	routes.Web()
	// routes.Api() // Uncomment to enable API routes
}

func (receiver *RouteServiceProvider) configureRateLimiting() {

}
