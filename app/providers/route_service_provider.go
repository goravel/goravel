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
	//Add HTTP middlewares
	facades.Route().GlobalMiddleware(http.Kernel{}.Middleware()...)
}

func (receiver *RouteServiceProvider) Boot(app foundation.Application) {
	receiver.configureRateLimiting()

	routes.Web()
}

func (receiver *RouteServiceProvider) configureRateLimiting() {

}
