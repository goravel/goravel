package docker

import (
	contractscache "github.com/goravel/framework/contracts/cache"
	contractsconfig "github.com/goravel/framework/contracts/config"
	contractsconsole "github.com/goravel/framework/contracts/console"
	contractsorm "github.com/goravel/framework/contracts/database/orm"
	contractsprocess "github.com/goravel/framework/contracts/process"
	"github.com/goravel/framework/contracts/testing/docker"
	"github.com/goravel/framework/errors"
)

type Docker struct {
	artisan contractsconsole.Artisan
	cache   contractscache.Cache
	config  contractsconfig.Config
	orm     contractsorm.Orm
	process contractsprocess.Process
}

func NewDocker(
	artisan contractsconsole.Artisan,
	cache contractscache.Cache,
	config contractsconfig.Config,
	orm contractsorm.Orm,
	process contractsprocess.Process,
) *Docker {
	return &Docker{
		artisan: artisan,
		cache:   cache,
		config:  config,
		orm:     orm,
		process: process,
	}
}

func (r *Docker) Cache(store ...string) (docker.CacheDriver, error) {
	if r.cache == nil {
		return nil, errors.CacheFacadeNotSet.SetModule(errors.ModuleTesting)
	}
	if len(store) == 0 {
		store = append(store, r.config.GetString("cache.default"))
	}

	return r.cache.Store(store[0]).Docker()
}

func (r *Docker) Database(connection ...string) (docker.Database, error) {
	if r.artisan == nil {
		return nil, errors.ConsoleFacadeNotSet.SetModule(errors.ModuleTesting)
	}
	if r.orm == nil {
		return nil, errors.OrmFacadeNotSet.SetModule(errors.ModuleTesting)
	}

	if len(connection) == 0 {
		return NewDatabase(r.artisan, r.config, r.orm, "")
	} else {
		return NewDatabase(r.artisan, r.config, r.orm, connection[0])
	}
}

func (r *Docker) Image(image docker.Image) docker.ImageDriver {
	return NewImageDriver(image, r.process)
}
