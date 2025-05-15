package providers

import (
	"github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/color"
)

type QueueServiceProvider struct {
}

func (receiver *QueueServiceProvider) Register(app foundation.Application) {
	queueFacade := app.MakeQueue()
	if queueFacade == nil {
		color.Errorln(errors.QueueFacadeNotSet.SetModule(errors.ModuleQueue))
		return
	}

	queueFacade.Register(receiver.Jobs())
}

func (receiver *QueueServiceProvider) Boot(app foundation.Application) {

}

func (receiver *QueueServiceProvider) Jobs() []queue.Job {
	return []queue.Job{}
}
