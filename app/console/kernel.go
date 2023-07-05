package console

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/schedule"
	"goravel/config"
)

type Kernel struct {
	app foundation.Application
}

// var _ foundation.Kernel = (*Kernel)(nil) // Ensure that Kernel implements foundation.Kernel.
// The foundation.Kernel interface is defined in contracts/foundation/kernel.go:
// Example:
// 	type Kernel interface {
// 		Bootstrap()
// 		Handle()
// 		Terminate()
// 	}

func Newkernel(app foundation.Application) *Kernel {
	return &Kernel{
		app: app,
	}
}

func (kernel *Kernel) Bootstrap() {
	// Bootstrap the application
	kernel.app.Boot()

	// Bootstrap the config.
	config.Boot()
}

func (kernel *Kernel) Handle() {
	defer func() {
		if err := recover(); err != nil {
			// ...call recovery handler...
		}
	}()

	kernel.Bootstrap()
	defer kernel.Terminate()

	// facades.Console().Run() // run console commands and the http server bind to Console commands
}

func (kernel *Kernel) Terminate() {
	// kernel.app.Terminate() // if app defined
}

func (kernel *Kernel) Schedule() []schedule.Event {
	return []schedule.Event{}
}

func (kernel *Kernel) Commands() []console.Command {
	return []console.Command{}
}
