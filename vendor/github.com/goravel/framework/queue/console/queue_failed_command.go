package console

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	contractsqueue "github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/support/color"
)

type QueueFailedCommand struct {
	queue contractsqueue.Queue
}

func NewQueueFailedCommand(queue contractsqueue.Queue) *QueueFailedCommand {
	return &QueueFailedCommand{
		queue: queue,
	}
}

// Signature The name and signature of the console command.
func (r *QueueFailedCommand) Signature() string {
	return "queue:failed"
}

// Description The console command description.
func (r *QueueFailedCommand) Description() string {
	return "List all of the failed queue jobs"
}

// Extend The console command extend.
func (r *QueueFailedCommand) Extend() command.Extend {
	return command.Extend{
		Category: "queue",
	}
}

// Handle Execute the console command.
func (r *QueueFailedCommand) Handle(ctx console.Context) error {
	failedJobs, err := r.queue.Failer().All()
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	if len(failedJobs) == 0 {
		ctx.Info(errors.QueueNoFailedJobsFound.Error())
		return nil
	}

	for _, failedJob := range failedJobs {
		r.printJob(ctx, failedJob.UUID(), failedJob.Connection(), failedJob.Queue())
	}

	ctx.Line("")

	return nil
}

func (r *QueueFailedCommand) printJob(ctx console.Context, uuid, connection, queue string) {
	datetime := color.Gray().Sprint(carbon.Now().ToDateTimeString())
	status := connection + "@" + queue
	first := datetime + " " + uuid
	second := status

	ctx.TwoColumnDetail(first, second)
}
