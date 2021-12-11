package providers

import (
	"github.com/goravel/framework/support/facades"
	"goravel/app/http"
	"goravel/routes"
)

type RouteServiceProvider struct {
}

func (router *RouteServiceProvider) Boot() {
	//Add HTTP middlewares.
	kernel := http.Kernel{}
	facades.Route.Use(kernel.Middleware()...)

	//Add routes
	routes.Web()
}

func (router *RouteServiceProvider) Register() {

}
