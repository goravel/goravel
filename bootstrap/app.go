package bootstrap

import (
	"github.com/goravel/framework/foundation"
	"goravel/app/http"
	"goravel/config"
)

func Boot() {
	//Create the application
	app := foundation.Application{}

	//Bootstrap the application
	app.Boot()

	//Bootstrap the http kernel, add http middlewares.
	app.BootHttpKernel(&http.Kernel{})

	//Bootstrap the config.
	config.Boot()
}
