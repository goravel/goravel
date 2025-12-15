package facades

import (
	"github.com/goravel/framework/contracts/cache"
)

func Cache() cache.Cache {
	return App().MakeCache()
}
