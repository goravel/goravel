package listeners

import (
	"github.com/goravel/framework/contracts/events"
	"github.com/goravel/framework/facades"
)

type TestAsyncListener struct {
}

func (receiver *TestAsyncListener) Signature() string {
	return "test_async_listener"
}

func (receiver *TestAsyncListener) Queue(args ...interface{}) events.Queue {
	return events.Queue{
		Enable:     true,
		Connection: "",
		Queue:      "",
	}
}

func (receiver *TestAsyncListener) Handle(args ...interface{}) error {
	facades.Log.Infof("test_async_listener: %s, %d", args[0], args[1])

	return nil
}
