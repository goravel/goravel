package facades

import (
	"github.com/goravel/framework/contracts/database/seeder"
)

func Seeder() seeder.Facade {
	return App().MakeSeeder()
}
