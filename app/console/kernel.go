package console

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/schedule/support"
)

type Kernel struct {
}

func (kernel Kernel) Schedule() []*support.Event {
	return []*support.Event{}
}

func (kernel Kernel) Commands() []console.Command {
	return []console.Command{}
}
