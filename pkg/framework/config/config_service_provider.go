package config

import (
	"github.com/goravel/framework/support/facades"
)

type ConfigServiceProvider struct {
	instance facades.ConfigFacade
}

func (config *ConfigServiceProvider) Boot() {
	facades.Config = config.instance
}

func (config *ConfigServiceProvider) Register() {
	viper := Viper{}
	config.instance = viper.Init()
}
