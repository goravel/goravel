package facades

import (
	"github.com/goravel/framework/contracts/event"
)

func Event() event.Instance {
	return App().MakeEvent()
}
