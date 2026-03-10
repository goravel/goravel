package postgres

import (
	"fmt"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/database"
	"github.com/goravel/framework/contracts/database/driver"
	"github.com/goravel/framework/contracts/log"
	"github.com/goravel/framework/contracts/process"
	"github.com/goravel/framework/contracts/testing/docker"
	"github.com/goravel/framework/errors"
	"github.com/goravel/postgres/contracts"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var _ driver.Driver = &Postgres{}

type Postgres struct {
	config  contracts.ConfigBuilder
	log     log.Log
	process process.Process
}

func NewPostgres(config config.Config, log log.Log, process process.Process, connection string) *Postgres {
	return &Postgres{
		config:  NewConfig(config, connection),
		log:     log,
		process: process,
	}
}

func (r *Postgres) Docker() (docker.DatabaseDriver, error) {
	if r.process == nil {
		return nil, errors.ProcessFacadeNotSet.SetModule(Name)
	}

	writers := r.config.Writers()
	if len(writers) == 0 {
		return nil, errors.DatabaseConfigNotFound
	}

	return NewDocker(r.config, r.process, writers[0].Database, writers[0].Username, writers[0].Password), nil
}

func (r *Postgres) Grammar() driver.Grammar {
	return NewGrammar(r.config.Writers()[0].Prefix)
}

func (r *Postgres) Pool() database.Pool {
	return database.Pool{
		Readers: r.fullConfigsToConfigs(r.config.Readers()),
		Writers: r.fullConfigsToConfigs(r.config.Writers()),
	}
}

func (r *Postgres) Processor() driver.Processor {
	return NewProcessor()
}

func (r *Postgres) fullConfigsToConfigs(fullConfigs []contracts.FullConfig) []database.Config {
	configs := make([]database.Config, len(fullConfigs))
	for i, fullConfig := range fullConfigs {
		configs[i] = database.Config{
			Connection:   fullConfig.Connection,
			Dsn:          fullConfig.Dsn,
			Database:     fullConfig.Database,
			Dialector:    fullConfigToDialector(fullConfig),
			Driver:       Name,
			Host:         fullConfig.Host,
			NameReplacer: fullConfig.NameReplacer,
			NoLowerCase:  fullConfig.NoLowerCase,
			Password:     fullConfig.Password,
			Port:         fullConfig.Port,
			Prefix:       fullConfig.Prefix,
			Schema:       fullConfig.Schema,
			Singular:     fullConfig.Singular,
			Sslmode:      fullConfig.Sslmode,
			Timezone:     fullConfig.Timezone,
			Username:     fullConfig.Username,
		}
	}

	return configs
}

func dsn(fullConfig contracts.FullConfig) string {
	if fullConfig.Dsn != "" {
		return fullConfig.Dsn
	}
	if fullConfig.Host == "" {
		return ""
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s&timezone=%s&search_path=%s",
		fullConfig.Username, fullConfig.Password, fullConfig.Host, fullConfig.Port, fullConfig.Database, fullConfig.Sslmode, fullConfig.Timezone, fullConfig.Schema)
}

func fullConfigToDialector(fullConfig contracts.FullConfig) gorm.Dialector {
	dsn := dsn(fullConfig)
	if dsn == "" {
		return nil
	}

	return postgres.New(postgres.Config{
		DSN: dsn,
		// When running a migration to add or remove columns, the driver will panic with cached plan must not change result type.
		// So PreferSimpleProtocol should be set to true to avoid this issue. The performance will be reduced a little bit,
		// but it's worth it to provide full migration support.
		PreferSimpleProtocol: true,
	})
}
