package bootstrap

import (
	"github.com/goravel/framework/contracts/foundation/configuration"
	"github.com/goravel/framework/foundation"

	"goravel/config"
	"goravel/routes"
)

func Boot() {
	foundation.Configure().
		WithConfig(config.Boot).
		WithProviders(Providers()).
		WithRouting([]func(){
			routes.Web,
			routes.Api,
			routes.Grpc,
		}).
		WithMiddleware(func(middleware configuration.Middleware) {
			middleware.Use()
		}).
		Run()
}
