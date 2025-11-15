package bootstrap

import (
	"github.com/goravel/framework/foundation"

	"goravel/config"
	"goravel/routes"
)

func Boot() {
	foundation.Setup().
		WithConfig(config.Boot).
		WithProviders(Providers()).
		WithRouting([]func(){
			routes.Web,
			routes.Api,
			routes.Grpc,
		}).
		Run()
}
