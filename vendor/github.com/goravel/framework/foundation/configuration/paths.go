package configuration

import (
	"github.com/goravel/framework/contracts/foundation/configuration"
	"github.com/goravel/framework/support"
)

type Paths struct {
}

func NewPaths() *Paths {
	return &Paths{}
}

func (r *Paths) App(path string) configuration.Paths {
	support.Config.Paths.App = path

	return r
}

func (r *Paths) Bootstrap(path string) configuration.Paths {
	support.Config.Paths.Bootstrap = path

	return r
}

func (r *Paths) Commands(path string) configuration.Paths {
	support.Config.Paths.Commands = path

	return r
}

func (r *Paths) Config(path string) configuration.Paths {
	support.Config.Paths.Config = path

	return r
}

func (r *Paths) Controllers(path string) configuration.Paths {
	support.Config.Paths.Controllers = path

	return r
}

func (r *Paths) Database(path string) configuration.Paths {
	support.Config.Paths.Database = path

	return r
}

func (r *Paths) Events(path string) configuration.Paths {
	support.Config.Paths.Events = path

	return r
}

func (r *Paths) Facades(path string) configuration.Paths {
	support.Config.Paths.Facades = path

	return r
}

func (r *Paths) Factories(path string) configuration.Paths {
	support.Config.Paths.Factories = path

	return r
}

func (r *Paths) Filters(path string) configuration.Paths {
	support.Config.Paths.Filters = path

	return r
}

func (r *Paths) Jobs(path string) configuration.Paths {
	support.Config.Paths.Jobs = path

	return r
}

func (r *Paths) Lang(path string) configuration.Paths {
	support.Config.Paths.Lang = path

	return r
}

func (r *Paths) Listeners(path string) configuration.Paths {
	support.Config.Paths.Listeners = path

	return r
}

func (r *Paths) Mails(path string) configuration.Paths {
	support.Config.Paths.Mails = path

	return r
}

func (r *Paths) Middleware(path string) configuration.Paths {
	support.Config.Paths.Middleware = path

	return r
}

func (r *Paths) Migrations(path string) configuration.Paths {
	support.Config.Paths.Migrations = path

	return r
}

func (r *Paths) Models(path string) configuration.Paths {
	support.Config.Paths.Models = path

	return r
}

func (r *Paths) Observers(path string) configuration.Paths {
	support.Config.Paths.Observers = path

	return r
}

func (r *Paths) Packages(path string) configuration.Paths {
	support.Config.Paths.Packages = path

	return r
}

func (r *Paths) Policies(path string) configuration.Paths {
	support.Config.Paths.Policies = path

	return r
}

func (r *Paths) Providers(path string) configuration.Paths {
	support.Config.Paths.Providers = path

	return r
}

func (r *Paths) Public(path string) configuration.Paths {
	support.Config.Paths.Public = path

	return r
}

func (r *Paths) Requests(path string) configuration.Paths {
	support.Config.Paths.Requests = path

	return r
}

func (r *Paths) Resources(path string) configuration.Paths {
	support.Config.Paths.Resources = path

	return r
}

func (r *Paths) Routes(path string) configuration.Paths {
	support.Config.Paths.Routes = path

	return r
}

func (r *Paths) Rules(path string) configuration.Paths {
	support.Config.Paths.Rules = path

	return r
}

func (r *Paths) Seeders(path string) configuration.Paths {
	support.Config.Paths.Seeders = path

	return r
}

func (r *Paths) Storage(path string) configuration.Paths {
	support.Config.Paths.Storage = path

	return r
}

func (r *Paths) Tests(path string) configuration.Paths {
	support.Config.Paths.Tests = path

	return r
}

func (r *Paths) Views(path string) configuration.Paths {
	support.Config.Paths.Views = path

	return r
}
