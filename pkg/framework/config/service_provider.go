package config

import (
	"github.com/goravel/framework/support/facades"
)

type ServiceProvider struct {
}

func (config *ServiceProvider) Boot() {

}

func (config *ServiceProvider) Register() {
	app := Application{}
	facades.Config = app.Init()
}
