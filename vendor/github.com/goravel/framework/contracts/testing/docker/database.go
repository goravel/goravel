package docker

import "github.com/goravel/framework/contracts/database/seeder"

type Database interface {
	DatabaseDriver
	// Migrate runs the database migrations.
	Migrate() error
	// Seed runs the database seeds.
	Seed(seeders ...seeder.Seeder) error
}

type DatabaseDriver interface {
	// Build a database container, it doesn't wait for the database to be ready, the Ready method needs to be called if
	// you want to check the container status.
	Build() error
	// Config get database configuration.
	Config() DatabaseConfig
	// Database returns a new instance with a new database, the Build method needs to be called first.
	Database(name string) (DatabaseDriver, error)
	// Driver gets the database driver name.
	Driver() string
	// Fresh the database.
	Fresh() error
	// Image gets the database image.
	Image(image Image)
	// Ready checks if the database is ready, the Build method needs to be called first.
	Ready() error
	// Reuse the existing database container.
	Reuse(containerID string, port int) error
	// Shutdown the database.
	Shutdown() error
}

type DatabaseConfig struct {
	Driver      string
	Host        string
	Database    string
	Username    string
	Password    string
	ContainerID string
	Port        int
}
