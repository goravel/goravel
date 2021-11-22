package database

import "github.com/goravel/framework/support/facades"

type ServiceProvider struct {
}

func (config *ServiceProvider) Boot() {

}

func (config *ServiceProvider) Register() {
	gorm := Gorm{}
	facades.DB = gorm.Init()
}
