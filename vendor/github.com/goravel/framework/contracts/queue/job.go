package queue

import (
	"time"
)

type Job interface {
	// Signature set the unique signature of the job.
	Signature() string
	// Handle executes the job.
	Handle(args ...any) error
}

type PendingJob interface {
	// Delay dispatches the task after the given delay.
	Delay(time time.Time) PendingJob
	// Dispatch dispatches the task.
	Dispatch() error
	// DispatchSync dispatches the task synchronously.
	DispatchSync() error
	// OnConnection sets the connection of the task.
	OnConnection(connection string) PendingJob
	// OnQueue sets the queue of the task.
	OnQueue(queue string) PendingJob
}

type ReservedJob interface {
	Delete() error
	Task() Task
}

type JobStorer interface {
	All() []Job
	Call(signature string, args []any) error
	Get(signature string) (Job, error)
	Register(jobs []Job)
}

// Deprecated: Use ChainJob instead.
type Jobs = ChainJob

type ChainJob struct {
	Delay time.Time `json:"delay"`
	Job   Job       `json:"job"`
	Args  []Arg     `json:"args"`
}

type JobWithShouldRetry interface {
	// ShouldRetry determines if the job should be retried based on the error.
	ShouldRetry(err error, attempt int) (retryable bool, delay time.Duration)
}
