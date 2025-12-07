package bootstrap

import (
	"github.com/goravel/framework/foundation"

	"goravel/config"
	"goravel/routes"
)

func Boot() {
	foundation.Setup().
		WithRouting([]func(){
			routes.Grpc,
			routes.Web,
		}).
		WithMigrations(Migrations()).
		WithConfig(config.Boot).
		WithProviders(Providers()).
		Run()
}
