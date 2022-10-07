package events

import "github.com/goravel/framework/contracts/events"

type TestEvent struct {
}

func (receiver *TestEvent) Handle(args []events.Arg) ([]events.Arg, error) {
	return args, nil
}
