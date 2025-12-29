package bootstrap

import (
	"github.com/goravel/framework/foundation"

	"goravel/config"
	"goravel/routes"
)

func Boot() {
	foundation.Setup().
		WithMigrations(Migrations()).
		WithRouting([]func(){
			routes.Web,
			routes.Grpc,
		}).
		WithProviders(Providers()).
		WithConfig(config.Boot).
		Run()
}
