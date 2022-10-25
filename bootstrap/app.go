package bootstrap

import (
	"goravel/config"

	"github.com/goravel/framework/foundation"
)

func Boot() {
	app := foundation.Application{}

	//Bootstrap the application
	app.Boot()

	//Bootstrap the config.
	config.Boot()
}
