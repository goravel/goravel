package foundation

import (
	"github.com/goravel/framework/config"
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
	BasePath string
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
		&config.ConfigServiceProvider{},
		&route.RouteServiceProvider{},
	}
}

func (app *Application) getConfiguredServiceProviders() []support.ServiceProvider {
	return facades.Config.Get("app.providers").([]support.ServiceProvider)
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
