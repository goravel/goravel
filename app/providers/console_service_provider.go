package providers

import (
	"github.com/goravel/framework/facades"

	"goravel/app/console"
)

type ConsoleServiceProvider struct {
}

func (receiver *ConsoleServiceProvider) Register() {
	kernel := console.Kernel{}
	facades.Schedule.Register(kernel.Schedule())
	facades.Artisan.Register(kernel.Commands())
}

func (receiver *ConsoleServiceProvider) Boot() {

}
