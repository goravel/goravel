package facades

import (
	"github.com/goravel/framework/contracts/session"
)

func Session() session.Manager {
	return App().MakeSession()
}
