package config

import (
	"github.com/goravel/framework/support/facades"
)

type RouteServiceProvider struct {
	instance facades.ConfigFacade
}

func (config *RouteServiceProvider) Boot() {
	facades.Config = config.instance
}

func (config *RouteServiceProvider) Register() {
}
