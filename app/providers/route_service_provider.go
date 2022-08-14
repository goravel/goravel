package providers

import (
	"github.com/goravel/framework/support/facades"
	"goravel/app/http"
	"goravel/routes"
)

type RouteServiceProvider struct {
}

func (receiver *RouteServiceProvider) Boot() {

}

func (receiver *RouteServiceProvider) Register() {
	//Add HTTP middlewares.
	facades.Route.Use(http.Kernel{}.Middleware()...)

	//Add routes
	routes.Web()
}
