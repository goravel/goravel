package event

import (
	"github.com/goravel/framework/contracts/event"
	"github.com/goravel/framework/errors"
)

type TestEvent struct{}

func (receiver *TestEvent) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}

type TestEventNoRegister struct{}

func (receiver *TestEventNoRegister) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}

type TestEventHandleError struct{}

func (receiver *TestEventHandleError) Handle(args []event.Arg) ([]event.Arg, error) {
	return nil, errors.New("some errors")
}

type TestCancelEvent struct{}

func (receiver *TestCancelEvent) Handle(args []event.Arg) ([]event.Arg, error) {
	return args, nil
}

type TestListener struct{}

func (receiver *TestListener) Signature() string {
	return "test_listener"
}

func (receiver *TestListener) Queue(args ...any) event.Queue {
	return event.Queue{
		Enable:     false,
		Connection: "",
		Queue:      "",
	}
}

func (receiver *TestListener) Handle(args ...any) error {
	return nil
}

type TestListenerHandleError struct{}

func (receiver *TestListenerHandleError) Signature() string {
	return "test_listener"
}

func (receiver *TestListenerHandleError) Queue(args ...any) event.Queue {
	return event.Queue{
		Enable:     false,
		Connection: "",
		Queue:      "",
	}
}

func (receiver *TestListenerHandleError) Handle(args ...any) error {
	return errors.New("error")
}
