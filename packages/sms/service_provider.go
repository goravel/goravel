package sms

import (
	configcontract "github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/support"
)

var App foundation.Application

type ServiceProvider struct {
	support.ServiceProvider
}

func (receiver *ServiceProvider) Register(app foundation.Application) {
	//receiver.Publishes(map[string]string{"a": "b"})
	App = app

	app.Bind("sms", func() (any, error) {
		config, err := app.Make("config")
		if err != nil {
			return nil, err
		}

		return NewSms(config.(configcontract.Config)), nil
	})
}

func (receiver *ServiceProvider) Boot(app foundation.Application) {

}
