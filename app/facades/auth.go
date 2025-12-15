package facades

import (
	"github.com/goravel/framework/contracts/auth"
	"github.com/goravel/framework/contracts/http"
)

func Auth(ctx ...http.Context) auth.Auth {
	return App().MakeAuth(ctx...)
}
