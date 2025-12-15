package facades

import (
	"github.com/goravel/framework/contracts/console"
)

func Artisan() console.Artisan {
	return App().MakeArtisan()
}
