package queue

import "github.com/goravel/framework/support/carbon"

type Failer interface {
	All() ([]FailedJob, error)
	Get(connection, queue string, uuids []string) ([]FailedJob, error)
}

type FailedJob interface {
	Connection() string
	Queue() string
	FailedAt() *carbon.DateTime
	Retry() error
	Signature() string
	UUID() string
}
