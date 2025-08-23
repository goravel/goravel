package providers

import (
	"github.com/goravel/framework/contracts/foundation"
)

type ConsoleServiceProvider struct {
}

func (receiver *ConsoleServiceProvider) Register(app foundation.Application) {
	// Commands and schedules can be registered in routes/console.go
	// or auto-discovered by the framework
}

func (receiver *ConsoleServiceProvider) Boot(app foundation.Application) {

}
