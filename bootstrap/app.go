package bootstrap

import (
	"github.com/goravel/framework/foundation"

	"goravel/config"
	"goravel/routes"
)

func Boot() {
	foundation.Configure().
		WithConfig(config.Boot).
		WithProviders(Providers()).
		WithRouting(
			routes.Web,
			routes.Api,
			routes.Grpc,
		).
		Run()
}
