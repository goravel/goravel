package events

import "github.com/goravel/framework/contracts/events"

type TestCancelEvent struct {
}

func (receiver *TestCancelEvent) Handle(args []events.Arg) ([]events.Arg, error) {
	return args, nil
}
