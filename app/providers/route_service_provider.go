package providers

import (
	"goravel/app/http"
	"goravel/routes"

	"github.com/goravel/framework/facades"
)

type RouteServiceProvider struct {
}

func (receiver *RouteServiceProvider) Boot() {

}

func (receiver *RouteServiceProvider) Register() {
	//Add HTTP middlewares.
	kernel := http.Kernel{}
	facades.Route.GlobalMiddleware(kernel.Middleware()...)

	//Add routes
	routes.Web()
}
