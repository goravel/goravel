package filesystem

import (
	"fmt"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/filesystem"
	"github.com/goravel/framework/errors"
)

type Driver string

const (
	DriverLocal  Driver = "local"
	DriverCustom Driver = "custom"
)

type Storage struct {
	filesystem.Driver
	config  config.Config
	drivers map[string]filesystem.Driver
}

func NewStorage(config config.Config) (*Storage, error) {
	defaultDisk := config.GetString("filesystems.default")
	if defaultDisk == "" {
		return nil, errors.FilesystemDefaultDiskNotSet.SetModule(errors.ModuleFilesystem)
	}

	driver, err := NewDriver(config, defaultDisk)
	if err != nil {
		return nil, err
	}

	drivers := make(map[string]filesystem.Driver)
	drivers[defaultDisk] = driver
	return &Storage{
		Driver:  driver,
		config:  config,
		drivers: drivers,
	}, nil
}

func NewDriver(config config.Config, disk string) (filesystem.Driver, error) {
	driver := Driver(config.GetString(fmt.Sprintf("filesystems.disks.%s.driver", disk)))
	switch driver {
	case DriverLocal:
		return NewLocal(config, disk)
	case DriverCustom:
		driver, ok := config.Get(fmt.Sprintf("filesystems.disks.%s.via", disk)).(filesystem.Driver)
		if ok {
			return driver, nil
		}

		driverCallback, ok := config.Get(fmt.Sprintf("filesystems.disks.%s.via", disk)).(func() (filesystem.Driver, error))
		if ok {
			return driverCallback()
		}

		return nil, errors.FilesystemInvalidCustomDriver.Args(disk)
	}

	return nil, errors.FilesystemDriverNotSupported.Args(driver)
}

func (r *Storage) Disk(disk string) filesystem.Driver {
	if driver, exist := r.drivers[disk]; exist {
		return driver
	}

	driver, err := NewDriver(r.config, disk)
	if err != nil {
		panic(err)
	}

	r.drivers[disk] = driver

	return driver
}
