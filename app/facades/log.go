package facades

import (
	"github.com/goravel/framework/contracts/log"
)

func Log() log.Log {
	return App().MakeLog()
}
