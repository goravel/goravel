package providers

import (
	"github.com/goravel/framework/contracts/events"
	"github.com/goravel/framework/facades"
)

type EventServiceProvider struct {
}

func (receiver *EventServiceProvider) Boot() {

}

func (receiver *EventServiceProvider) Register() {
	facades.Event.Register(receiver.listen())
}

func (receiver *EventServiceProvider) listen() map[events.Event][]events.Listener {
	return map[events.Event][]events.Listener{}
}
