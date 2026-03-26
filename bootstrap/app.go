package bootstrap

import (
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/foundation"

	"goravel/config"
	"goravel/routes"
)

func Boot() contractsfoundation.Application {
	return foundation.Setup().
		WithRouting(func() {
			routes.Web()
		}).
		WithProviders(Providers).
		WithConfig(config.Boot).
		Create()
}
