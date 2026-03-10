package process

import (
	"github.com/goravel/framework/contracts/binding"
	"github.com/goravel/framework/contracts/foundation"
)

type ServiceProvider struct {
}

func (r *ServiceProvider) Relationship() binding.Relationship {
	return binding.Relationship{
		Bindings: []string{
			binding.Process,
		},
		Dependencies: binding.Bindings[binding.Process].Dependencies,
		ProvideFor:   []string{},
	}
}

func (r *ServiceProvider) Register(app foundation.Application) {
	app.Bind(binding.Process, func(app foundation.Application) (any, error) {
		return New(), nil
	})
}

func (r *ServiceProvider) Boot(_ foundation.Application) {
}
