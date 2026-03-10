package postgres

import (
	"fmt"
	"strconv"
	"time"

	contractsprocess "github.com/goravel/framework/contracts/process"
	contractsdocker "github.com/goravel/framework/contracts/testing/docker"
	"github.com/goravel/framework/support/color"
	supportdocker "github.com/goravel/framework/support/docker"
	testingdocker "github.com/goravel/framework/testing/docker"
	"github.com/spf13/cast"
	"gorm.io/driver/postgres"
	gormio "gorm.io/gorm"

	"github.com/goravel/postgres/contracts"
)

type Docker struct {
	config         contracts.ConfigBuilder
	databaseConfig contractsdocker.DatabaseConfig
	imageDriver    contractsdocker.ImageDriver
	process        contractsprocess.Process
}

func NewDocker(config contracts.ConfigBuilder, process contractsprocess.Process, database, username, password string) *Docker {
	return &Docker{
		config: config,
		databaseConfig: contractsdocker.DatabaseConfig{
			Database: database,
			Driver:   Name,
			Host:     "127.0.0.1",
			Password: password,
			Port:     5432,
			Username: username,
		},
		imageDriver: testingdocker.NewImageDriver(contractsdocker.Image{
			Repository: "postgres",
			Tag:        "latest",
			Env: []string{
				"POSTGRES_USER=" + username,
				"POSTGRES_PASSWORD=" + password,
				"POSTGRES_DB=" + database,
			},
			ExposedPorts: []string{"5432"},
			Args:         []string{"-c max_connections=1000"},
		}, process),
		process: process,
	}
}

func (r *Docker) Build() error {
	if err := r.imageDriver.Build(); err != nil {
		return err
	}

	config := r.imageDriver.Config()
	r.databaseConfig.ContainerID = config.ContainerID
	r.databaseConfig.Port = cast.ToInt(supportdocker.ExposedPort(config.ExposedPorts, strconv.Itoa(r.databaseConfig.Port)))

	return nil
}

func (r *Docker) Config() contractsdocker.DatabaseConfig {
	return r.databaseConfig
}

func (r *Docker) Database(name string) (contractsdocker.DatabaseDriver, error) {
	go func() {
		gormDB, err := r.connect()
		if err != nil {
			color.Errorf("connect Postgres error: %v", err)
			return
		}

		res := gormDB.Exec(fmt.Sprintf(`CREATE DATABASE "%s";`, name))
		if res.Error != nil {
			color.Errorf("create Postgres database error: %v", res.Error)
		}

		if err := r.close(gormDB); err != nil {
			color.Errorf("close Postgres connection error: %v", err)
		}
	}()

	docker := NewDocker(r.config, r.process, name, r.databaseConfig.Username, r.databaseConfig.Password)
	docker.databaseConfig.ContainerID = r.databaseConfig.ContainerID
	docker.databaseConfig.Port = r.databaseConfig.Port

	return docker, nil
}

func (r *Docker) Driver() string {
	return Name
}

func (r *Docker) Fresh() error {
	gormDB, err := r.connect()
	if err != nil {
		return fmt.Errorf("connect Postgres error when clearing: %v", err)
	}

	if res := gormDB.Exec("DROP SCHEMA public CASCADE;"); res.Error != nil {
		return fmt.Errorf("drop schema of Postgres error: %v", res.Error)
	}

	if res := gormDB.Exec("CREATE SCHEMA public;"); res.Error != nil {
		return fmt.Errorf("create schema of Postgres error: %v", res.Error)
	}

	return r.close(gormDB)
}

func (r *Docker) Image(image contractsdocker.Image) {
	r.imageDriver = testingdocker.NewImageDriver(image, r.process)
}

func (r *Docker) Ready() error {
	gormDB, err := r.connect()
	if err != nil {
		return err
	}

	r.resetConfigPort()

	return r.close(gormDB)
}

func (r *Docker) Reuse(containerID string, port int) error {
	r.databaseConfig.ContainerID = containerID
	r.databaseConfig.Port = port

	return nil
}

func (r *Docker) Shutdown() error {
	return r.imageDriver.Shutdown()
}

func (r *Docker) connect() (*gormio.DB, error) {
	var (
		instance *gormio.DB
		err      error
	)

	// docker compose need time to start
	for i := 0; i < 60; i++ {
		instance, err = gormio.Open(postgres.New(postgres.Config{
			DSN: fmt.Sprintf("postgres://%s:%s@%s:%d/%s", r.databaseConfig.Username, r.databaseConfig.Password, r.databaseConfig.Host, r.databaseConfig.Port, r.databaseConfig.Database),
		}))

		if err == nil {
			break
		}

		time.Sleep(1 * time.Second)
	}

	return instance, err
}

func (r *Docker) close(gormDB *gormio.DB) error {
	db, err := gormDB.DB()
	if err != nil {
		return err
	}

	return db.Close()
}

func (r *Docker) resetConfigPort() {
	writers := r.config.Config().Get(fmt.Sprintf("database.connections.%s.write", r.config.Connection()))
	if writeConfigs, ok := writers.([]contracts.Config); ok {
		writeConfigs[0].Port = r.databaseConfig.Port
		r.config.Config().Add(fmt.Sprintf("database.connections.%s.write", r.config.Connection()), writeConfigs)

		return
	}

	r.config.Config().Add(fmt.Sprintf("database.connections.%s.port", r.config.Connection()), r.databaseConfig.Port)
}
