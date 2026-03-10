package queue

import (
	"sync"

	contractsqueue "github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/errors"
)

type JobStorer struct {
	jobs sync.Map
}

func NewJobStorer() *JobStorer {
	return &JobStorer{}
}

// All gets all registered jobs
func (r *JobStorer) All() []contractsqueue.Job {
	var jobs []contractsqueue.Job
	r.jobs.Range(func(_, value any) bool {
		jobs = append(jobs, value.(contractsqueue.Job))
		return true
	})

	return jobs
}

// Call calls a registered job using its signature
func (r *JobStorer) Call(signature string, args []any) error {
	job, err := r.Get(signature)
	if err != nil {
		return err
	}

	return job.Handle(args...)
}

// Get gets a registered job using its signature
func (r *JobStorer) Get(signature string) (contractsqueue.Job, error) {
	if job, ok := r.jobs.Load(signature); ok {
		return job.(contractsqueue.Job), nil
	}

	return nil, errors.QueueJobNotFound.Args(signature)
}

// Register registers jobs to the job manager
func (r *JobStorer) Register(jobs []contractsqueue.Job) {
	for _, job := range jobs {
		r.jobs.Store(job.Signature(), job)
	}
}
