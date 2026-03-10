package schedule

import (
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/schedule"
)

type ScheduleRunner struct {
	config   config.Config
	schedule schedule.Schedule
}

func NewScheduleRunner(config config.Config, schedule schedule.Schedule) *ScheduleRunner {
	return &ScheduleRunner{
		config:   config,
		schedule: schedule,
	}
}

func (r *ScheduleRunner) Signature() string {
	return "schedule"
}

func (r *ScheduleRunner) ShouldRun() bool {
	return r.schedule != nil && len(r.schedule.Events()) > 0 && r.config.GetBool("app.auto_run", true)
}

func (r *ScheduleRunner) Run() error {
	r.schedule.Run()

	return nil
}

func (r *ScheduleRunner) Shutdown() error {
	return r.schedule.Shutdown()
}
