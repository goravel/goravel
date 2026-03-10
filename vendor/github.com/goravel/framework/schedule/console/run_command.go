package console

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/schedule"
)

type Run struct {
	schedule schedule.Schedule
}

func NewRun(schedule schedule.Schedule) *Run {
	return &Run{
		schedule: schedule,
	}
}

// Signature The name and signature of the console command.
func (r *Run) Signature() string {
	return "schedule:run"
}

// Description The console command description.
func (r *Run) Description() string {
	return "Run the scheduled commands"
}

// Extend The console command extend.
func (r *Run) Extend() command.Extend {
	return command.Extend{
		Category: "schedule",
	}
}

// Handle Execute the console command.
func (r *Run) Handle(_ console.Context) error {
	r.schedule.Run()

	return nil
}
