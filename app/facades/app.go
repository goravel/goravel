package facades

import (
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	"github.com/goravel/framework/foundation"
)

func App() contractsfoundation.Application {
	return foundation.App
}
