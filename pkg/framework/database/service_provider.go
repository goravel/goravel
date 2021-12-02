package database

import (
	"github.com/goravel/framework/console/support"
	"github.com/goravel/framework/database/console/migrations"
	"github.com/goravel/framework/support/facades"
)

type ServiceProvider struct {
}

func (database *ServiceProvider) Boot() {

}

func (database *ServiceProvider) Register() {
	app := Application{}
	facades.DB = app.Init()

	database.registerCommand()
}

func (database *ServiceProvider) registerCommand() {
	facades.Artisan.Register([]support.Command{
		migrations.MigrateMakeCommand{},
		migrations.MigrateCommand{},
	})
}
