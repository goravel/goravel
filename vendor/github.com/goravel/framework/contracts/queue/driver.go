package queue

const (
	DriverSync     string = "sync"
	DriverDatabase string = "database"
	DriverCustom   string = "custom"
)

type DriverCreator interface {
	Create(connection string) (Driver, error)
}

type Driver interface {
	// Driver returns the driver name for the driver.
	Driver() string
	// Pop pops the next job off of the queue.
	Pop(queue string) (ReservedJob, error)
	// Push pushes the job onto the queue.
	Push(task Task, queue string) error
}
