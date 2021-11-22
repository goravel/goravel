package config

import (
	"github.com/goravel/framework/support/facades"
)

type ServiceProvider struct {
}

func (config *ServiceProvider) Boot() {

}

func (config *ServiceProvider) Register() {
	viper := Viper{}
	facades.Config = viper.Init()
}
