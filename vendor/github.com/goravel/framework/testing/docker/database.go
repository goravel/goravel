package docker

import (
	"fmt"

	contractsconfig "github.com/goravel/framework/contracts/config"
	contractsconsole "github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/database/driver"
	contractsorm "github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/database/seeder"
	"github.com/goravel/framework/contracts/testing/docker"
	"github.com/goravel/framework/errors"
)

type Database struct {
	docker.DatabaseDriver
	artisan    contractsconsole.Artisan
	config     contractsconfig.Config
	orm        contractsorm.Orm
	connection string
}

func NewDatabase(artisan contractsconsole.Artisan, config contractsconfig.Config, orm contractsorm.Orm, connection string) (*Database, error) {
	if artisan == nil {
		return nil, errors.ConsoleFacadeNotSet
	}
	if config == nil {
		return nil, errors.ConfigFacadeNotSet
	}
	if orm == nil {
		return nil, errors.OrmFacadeNotSet
	}

	if connection == "" {
		connection = config.GetString("database.default")
	}

	databaseDriverCallback, exist := config.Get(fmt.Sprintf("database.connections.%s.via", connection)).(func() (driver.Driver, error))
	if !exist {
		return nil, errors.DatabaseConfigNotFound
	}
	databaseDriver, err := databaseDriverCallback()
	if err != nil {
		return nil, err
	}

	databaseDocker, err := databaseDriver.Docker()
	if err != nil {
		return nil, err
	}

	return &Database{
		DatabaseDriver: databaseDocker,
		artisan:        artisan,
		config:         config,
		connection:     connection,
		orm:            orm,
	}, nil
}

func (r *Database) Migrate() error {
	return r.artisan.Call("--no-ansi migrate")
}

func (r *Database) Ready() error {
	if err := r.DatabaseDriver.Ready(); err != nil {
		return err
	}

	r.orm.Fresh()

	return nil
}

func (r *Database) Seed(seeders ...seeder.Seeder) error {
	command := "db:seed"
	if len(seeders) > 0 {
		command += " --seeder"
		for _, seed := range seeders {
			command += fmt.Sprintf(" %s", seed.Signature())
		}
	}

	return r.artisan.Call(command)
}
