package queue

import (
	"github.com/goravel/framework/contracts/binding"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/errors"
	queueconsole "github.com/goravel/framework/queue/console"
)

type ServiceProvider struct {
}

func (r *ServiceProvider) Relationship() binding.Relationship {
	return binding.Relationship{
		Bindings: []string{
			binding.Queue,
		},
		Dependencies: binding.Bindings[binding.Queue].Dependencies,
		ProvideFor:   []string{},
	}
}

func (r *ServiceProvider) Register(app foundation.Application) {
	app.Singleton(binding.Queue, func(app foundation.Application) (any, error) {
		config := app.MakeConfig()
		if config == nil {
			return nil, errors.ConfigFacadeNotSet.SetModule(errors.ModuleQueue)
		}

		log := app.MakeLog()
		if log == nil {
			return nil, errors.LogFacadeNotSet.SetModule(errors.ModuleQueue)
		}

		queueConfig := NewConfig(config)
		job := NewJobStorer()

		return NewApplication(queueConfig, app.MakeCache(), app.MakeDB(), job, app.Json(), log), nil
	})
}

func (r *ServiceProvider) Boot(app foundation.Application) {
	app.MakeArtisan().Register([]console.Command{
		&queueconsole.JobMakeCommand{},
		queueconsole.NewQueueRetryCommand(app.MakeQueue(), app.GetJson()),
		queueconsole.NewQueueFailedCommand(app.MakeQueue()),
	})
}

func (r *ServiceProvider) Runners(app foundation.Application) []foundation.Runner {
	return []foundation.Runner{NewQueueRunner(app.MakeConfig(), app.MakeQueue())}
}
