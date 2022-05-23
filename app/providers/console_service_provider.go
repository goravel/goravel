package providers

import (
	"github.com/goravel/framework/support/facades"
	"goravel/app/console"
)

type ConsoleServiceProvider struct {
}

func (receiver *ConsoleServiceProvider) Boot() {
	facades.Schedule.Register(console.Kernel{}.Schedule())
	facades.Artisan.Register(console.Kernel{}.Commands())
}

func (receiver *ConsoleServiceProvider) Register() {

}
