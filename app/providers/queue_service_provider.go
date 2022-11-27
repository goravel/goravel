package providers

import (
	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/facades"
)

type QueueServiceProvider struct {
}

func (receiver *QueueServiceProvider) Register() {
	facades.Queue.Register(receiver.Jobs())
}

func (receiver *QueueServiceProvider) Boot() {

}

func (receiver *QueueServiceProvider) Jobs() []queue.Job {
	return []queue.Job{}
}
