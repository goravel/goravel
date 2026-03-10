package cache

import (
	"fmt"

	"github.com/goravel/framework/contracts/cache"
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/errors"
)

type Driver struct {
	config config.Config
}

func NewDriver(config config.Config) *Driver {
	return &Driver{
		config: config,
	}
}

func (d *Driver) New(store string) (cache.Driver, error) {
	driver := d.config.GetString(fmt.Sprintf("cache.stores.%s.driver", store))
	switch driver {
	case "memory":
		return d.memory()
	case "custom":
		return d.custom(store)
	default:
		return nil, errors.CacheDriverNotSupported.Args(driver)
	}
}

func (d *Driver) memory() (cache.Driver, error) {
	return NewMemory(d.config)
}

func (d *Driver) custom(store string) (cache.Driver, error) {
	if custom, ok := d.config.Get(fmt.Sprintf("cache.stores.%s.via", store)).(cache.Driver); ok {
		return custom, nil
	}
	if custom, ok := d.config.Get(fmt.Sprintf("cache.stores.%s.via", store)).(func() (cache.Driver, error)); ok {
		return custom()
	}

	return nil, errors.CacheStoreContractNotFulfilled.Args(store)
}
