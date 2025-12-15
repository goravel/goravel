package facades

import (
	"github.com/goravel/framework/contracts/filesystem"
)

func Storage() filesystem.Storage {
	return App().MakeStorage()
}
