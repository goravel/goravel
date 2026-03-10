package queue

import (
	"time"

	"github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/support/carbon"
)

var (
	SyncDriverName              = "sync"
	_              queue.Driver = &Sync{}
)

type Sync struct {
}

func NewSync() *Sync {
	return &Sync{}
}

func (r *Sync) Driver() string {
	return queue.DriverSync
}

func (r *Sync) Pop(_ string) (queue.ReservedJob, error) {
	// sync driver does not support pop
	return nil, nil
}

func (r *Sync) Push(task queue.Task, _ string) error {
	if err := push(task.ChainJob); err != nil {
		return err
	}

	if len(task.Chain) > 0 {
		for _, chain := range task.Chain {
			if err := push(chain); err != nil {
				return err
			}
		}
	}

	return nil
}

func push(job queue.ChainJob) error {
	if !job.Delay.IsZero() {
		time.Sleep(carbon.FromStdTime(job.Delay).DiffAbsInDuration())
	}

	var realArgs []any
	for _, arg := range job.Args {
		realArgs = append(realArgs, arg.Value)
	}

	if err := job.Job.Handle(realArgs...); err != nil {
		return err
	}

	return nil
}
