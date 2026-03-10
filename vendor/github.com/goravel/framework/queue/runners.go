package queue

import (
	"fmt"

	"github.com/goravel/framework/contracts/config"
	contractsqueue "github.com/goravel/framework/contracts/queue"
)

type QueueRunner struct {
	config config.Config
	worker contractsqueue.Worker
}

func NewQueueRunner(config config.Config, queue contractsqueue.Queue) *QueueRunner {
	var worker contractsqueue.Worker
	if queue != nil {
		worker = queue.Worker()
	}

	return &QueueRunner{
		config: config,
		worker: worker,
	}
}

func (r *QueueRunner) Signature() string {
	return "queue"
}

func (r *QueueRunner) ShouldRun() bool {
	connection := r.config.GetString("queue.default")

	return r.worker != nil &&
		connection != "" &&
		r.config.GetString(fmt.Sprintf("queue.connections.%s.driver", connection)) != SyncDriverName &&
		r.config.GetBool("app.auto_run", true)
}

func (r *QueueRunner) Run() error {
	return r.worker.Run()
}

func (r *QueueRunner) Shutdown() error {
	return r.worker.Shutdown()
}
