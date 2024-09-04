package console

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/schedule"
)

type Kernel struct {
}

func (kernel *Kernel) Schedule() []schedule.Event {
	// When you add some commands to run scheduled tasks
	// Remove comments from the 'main.go' file
	return []schedule.Event{}
}

func (kernel *Kernel) Commands() []console.Command {
	return []console.Command{}
}
