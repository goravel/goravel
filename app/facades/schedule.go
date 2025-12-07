package facades

import (
	"github.com/goravel/framework/contracts/schedule"
)

func Schedule() schedule.Schedule {
	return App().MakeSchedule()
}
