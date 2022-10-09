package providers

import (
	"github.com/goravel/framework/facades"
	"goravel/app/http"
	"goravel/routes"
)

type RouteServiceProvider struct {
}

func (receiver *RouteServiceProvider) Boot() {

}

func (receiver *RouteServiceProvider) Register() {
	//Add HTTP middlewares.
	kernel := http.Kernel{}
	facades.Route.Middleware(kernel.Middleware()...)

	//Add routes
	routes.Web()
}
