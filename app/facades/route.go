package facades

import (
	"github.com/goravel/framework/contracts/route"
)

func Route() route.Route {
	return App().MakeRoute()
}
