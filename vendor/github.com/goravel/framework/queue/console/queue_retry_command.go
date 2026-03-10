package console

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/foundation"
	contractsqueue "github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/carbon"
)

type QueueRetryCommand struct {
	json  foundation.Json
	queue contractsqueue.Queue
}

func NewQueueRetryCommand(queue contractsqueue.Queue, json foundation.Json) *QueueRetryCommand {
	return &QueueRetryCommand{
		json:  json,
		queue: queue,
	}
}

// Signature The name and signature of the console command.
func (r *QueueRetryCommand) Signature() string {
	return "queue:retry"
}

// Description The console command description.
func (r *QueueRetryCommand) Description() string {
	return "Retry a failed queue job"
}

// Extend The console command extend.
func (r *QueueRetryCommand) Extend() command.Extend {
	return command.Extend{
		Category: "queue",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:    "connection",
				Aliases: []string{"c"},
				Usage:   "Retry all of the failed jobs for the specified connection",
			},
			&command.BoolFlag{
				Name:    "queue",
				Aliases: []string{"q"},
				Usage:   "Retry all of the failed jobs for the specified queue",
			},
		},
	}
}

// Handle Execute the console command.
func (r *QueueRetryCommand) Handle(ctx console.Context) error {
	failedJobs, err := r.queue.Failer().Get(ctx.Option("connection"), ctx.Option("queue"), ctx.Arguments())
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	if len(failedJobs) == 0 {
		ctx.Info(errors.QueueNoRetryableJobsFound.Error())
		return nil
	}

	ctx.Info(errors.QueuePushingFailedJob.Error())
	ctx.Line("")

	for _, failedJob := range failedJobs {
		now := carbon.Now()

		if err := failedJob.Retry(); err != nil {
			ctx.Error(errors.QueueFailedToRetryJob.Args(failedJob, err).Error())
			continue
		}

		r.printSuccess(ctx, failedJob.UUID(), now.DiffAbsInDuration().String())
	}

	ctx.Line("")

	return nil
}

func (r *QueueRetryCommand) printSuccess(ctx console.Context, uuid, duration string) {
	status := "<fg=green;op=bold>DONE</>"
	first := uuid
	second := duration + " " + status

	ctx.TwoColumnDetail(first, second)
}
