package facades

import (
	"github.com/goravel/framework/contracts/process"
)

func Process() process.Process {
	return App().MakeProcess()
}
