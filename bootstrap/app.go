package bootstrap

import (
	"github.com/goravel/framework/foundation"
	"goravel/config"
)

func Boot() {
	app := foundation.Application{}
	app.Boot()
	config.Boot()
}
