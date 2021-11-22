package foundation

import (
	"github.com/goravel/framework/config"
	"github.com/goravel/framework/support"
)

func init() {
	app := Application{}
	app.registerBaseServiceProviders()
}

type Application struct {
	BasePath string
}

func (app *Application) Boot() {
}

func (app *Application) GetBasePath() string {
	return app.BasePath
}

func (app *Application) register(serviceProvider support.ServiceProvider) {
	serviceProvider.Register()
	serviceProvider.Boot()
}

func (app *Application) registerBaseServiceProviders() {
	app.register(&config.ConfigServiceProvider{})
}
