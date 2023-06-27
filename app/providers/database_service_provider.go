package providers

import (
	"goravel/database/seeders"

	"github.com/goravel/framework/contracts/database/seeder"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/facades"
)

type DatabaseServiceProvider struct {
}

func (receiver *DatabaseServiceProvider) Register(app foundation.Application) {

}

func (receiver *DatabaseServiceProvider) Boot(app foundation.Application) {
	facades.Seeder().Register([]seeder.Seeder{
		&seeders.DatabaseSeeder{},
	})
}
