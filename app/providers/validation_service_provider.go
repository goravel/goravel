package providers

import (
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/facades"
)

type ValidationServiceProvider struct {
}

func (receiver *ValidationServiceProvider) Register(app foundation.Application) {

}

func (receiver *ValidationServiceProvider) Boot(app foundation.Application) {
	if err := facades.Validation().AddRules(receiver.rules()); err != nil {
		facades.Log().Errorf("add rules error: %+v", err)
	}
}

func (receiver *ValidationServiceProvider) rules() []validation.Rule {
	return []validation.Rule{}
}
