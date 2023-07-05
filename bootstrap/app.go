package bootstrap

import (
	foundationcontract "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/foundation"

	"goravel/config"
)

func App() foundationcontract.Application {
	return foundation.NewApplication()
}

func Boot() {
	app := foundation.NewApplication()

	// Bootstrap the application
	app.Boot()

	// Bootstrap the config.
	config.Boot()
}
