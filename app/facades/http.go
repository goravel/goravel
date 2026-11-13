package facades

import (
	"github.com/goravel/framework/contracts/http/client"
)

func Http() client.Request {
	return App().MakeHttp()
}
