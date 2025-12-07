package facades

import (
	"github.com/goravel/framework/contracts/config"
)

func Config() config.Config {
	return App().MakeConfig()
}
