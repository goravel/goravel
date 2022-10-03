package providers

import (
	contractevent "github.com/goravel/framework/contracts/events"
	"github.com/goravel/framework/facades"

	"goravel/app/events"
	"goravel/app/listeners"
)

type EventServiceProvider struct {
}

func (receiver *EventServiceProvider) Boot() {

}

func (receiver *EventServiceProvider) Register() {
	facades.Event.Register(receiver.listen())
}

func (receiver *EventServiceProvider) listen() map[contractevent.Event][]contractevent.Listener {
	return map[contractevent.Event][]contractevent.Listener{
		&events.Test{}: {&listeners.TestListener{}},
	}
}
