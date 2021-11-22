package config

import (
	"github.com/goravel/framework/support/facades"
)

type ConfigServiceProvider struct {
}

func (config *ConfigServiceProvider) Boot() {

}

func (config *ConfigServiceProvider) Register() {
	viper := Viper{}
	facades.Config = viper.Init()
}
