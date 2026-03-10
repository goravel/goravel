package testing

import (
	"github.com/goravel/framework/contracts/testing/docker"
)

type TestingT interface {
	Errorf(format string, args ...any)
	FailNow()
}

type Testing interface {
	// Docker get the Docker instance.
	Docker() Docker
}

type Docker interface {
	// Cache gets a cache connection instance.
	Cache(store ...string) (docker.CacheDriver, error)
	// Database gets a database connection instance.
	Database(connection ...string) (docker.Database, error)
	// Image gets a image instance.
	Image(image docker.Image) docker.ImageDriver
}
