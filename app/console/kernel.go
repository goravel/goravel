package console

import (
	"github.com/goravel/framework/contracts/schedule"
)

type Kernel struct {
}

func (kernel Kernel) Schedule() []schedule.Event {
	return []schedule.Event{}
}
