package facades

import (
	"github.com/goravel/framework/contracts/auth/access"
	"github.com/goravel/framework/contracts/http"
)

func Gate(ctx ...http.Context) access.Gate {
	return App().MakeGate()
}
