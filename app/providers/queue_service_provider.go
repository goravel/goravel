package providers

import (
	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/support/facades"
)

type QueueServiceProvider struct {
}

func (receiver *QueueServiceProvider) Boot() {
	facades.Queue.Register(receiver.Jobs())
}

func (receiver *QueueServiceProvider) Register() {

}

func (receiver *QueueServiceProvider) Jobs() []queue.Job {
	return []queue.Job{}
}
