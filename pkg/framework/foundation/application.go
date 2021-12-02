package foundation

import (
	"github.com/goravel/framework/config"
	"github.com/goravel/framework/console"
	"github.com/goravel/framework/database"
	"github.com/goravel/framework/route"
	"github.com/goravel/framework/support"
	"github.com/goravel/framework/support/facades"
)

func init() {
	app := Application{}
	app.registerBaseServiceProviders()
	app.bootBaseServiceProviders()
}

type Application struct {
}

func (app *Application) Boot() {
}

func (app *Application) BootHttpKernel(kernel support.Kernel) {
	app.registerConfiguredServiceProviders()
	app.bootConfiguredServiceProviders()
	facades.Route.Use(kernel.Middleware()...)
}

func (app *Application) register(serviceProvider support.ServiceProvider) {
	serviceProvider.Register()
}

func (app *Application) boot(serviceProvider support.ServiceProvider) {
	serviceProvider.Boot()
}

func (app *Application) getBaseServiceProviders() []support.ServiceProvider {
	return []support.ServiceProvider{
		&config.ServiceProvider{},
	}
}

func (app *Application) getConfiguredServiceProviders() []support.ServiceProvider {
	configuredServiceProviders := []support.ServiceProvider{
		&database.ServiceProvider{},
		&console.ServiceProvider{},
		&route.ServiceProvider{},
	}

	configuredServiceProviders = append(configuredServiceProviders, facades.Config.Get("app.providers").([]support.ServiceProvider)...)

	// Last load
	//configuredServiceProviders = append(configuredServiceProviders, &console.ServiceProvider{})

	return configuredServiceProviders
}

func (app *Application) registerBaseServiceProviders() {
	app.registerServiceProviders(app.getBaseServiceProviders())
}

func (app *Application) bootBaseServiceProviders() {
	app.bootServiceProviders(app.getBaseServiceProviders())
}

func (app *Application) registerConfiguredServiceProviders() {
	app.registerServiceProviders(app.getConfiguredServiceProviders())
}

func (app *Application) bootConfiguredServiceProviders() {
	app.bootServiceProviders(app.getConfiguredServiceProviders())
}

func (app *Application) registerServiceProviders(serviceProviders []support.ServiceProvider) {
	for _, serviceProvider := range serviceProviders {
		app.register(serviceProvider)
	}
}

func (app *Application) bootServiceProviders(serviceProviders []support.ServiceProvider) {
	for _, serviceProvider := range serviceProviders {
		app.boot(serviceProvider)
	}
}
