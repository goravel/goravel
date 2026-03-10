package event

import (
	"github.com/goravel/framework/contracts/event"
	contractsqueue "github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/errors"
)

type Task struct {
	event     event.Event
	queue     contractsqueue.Queue
	args      []event.Arg
	listeners []event.Listener
}

func NewTask(queue contractsqueue.Queue, args []event.Arg, event event.Event, listeners []event.Listener) *Task {
	return &Task{
		args:      args,
		event:     event,
		listeners: listeners,
		queue:     queue,
	}
}

func (receiver *Task) Dispatch() error {
	if len(receiver.listeners) == 0 {
		return errors.EventListenerNotBind.Args(receiver.event)
	}

	handledArgs, err := receiver.event.Handle(receiver.args)
	if err != nil {
		return err
	}

	var mapArgs []any
	for _, arg := range handledArgs {
		mapArgs = append(mapArgs, arg.Value)
	}

	for _, listener := range receiver.listeners {
		var err error
		task := receiver.queue.Job(listener, eventArgsToQueueArgs(handledArgs))
		queue := listener.Queue(mapArgs...)
		if queue.Connection != "" {
			task.OnConnection(queue.Connection)
		}
		if queue.Queue != "" {
			task.OnQueue(queue.Queue)
		}
		if queue.Enable {
			err = task.Dispatch()
		} else {
			err = task.DispatchSync()
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func eventArgsToQueueArgs(args []event.Arg) []contractsqueue.Arg {
	var queueArgs []contractsqueue.Arg
	for _, arg := range args {
		queueArgs = append(queueArgs, contractsqueue.Arg{
			Type:  arg.Type,
			Value: arg.Value,
		})
	}

	return queueArgs
}
