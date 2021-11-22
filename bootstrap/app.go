package bootstrap

import (
	"github.com/goravel/framework/foundation"
	"goravel/app/http"
	"goravel/config"
)

func Boot() {
	app := foundation.Application{}
	app.Boot()
	app.BootHttpKernel(&http.Kernel{})
	config.Boot()

}
