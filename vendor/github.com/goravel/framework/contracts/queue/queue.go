package queue

type Queue interface {
	// Connection gets a driver instance by connection name
	Connection(name string) (Driver, error)
	// Chain creates a chain of jobs to be processed one by one, passing
	Chain(jobs []ChainJob) PendingJob
	// Failer gets failed jobs
	Failer() Failer
	// GetJob gets job by signature
	GetJob(signature string) (Job, error)
	// GetJobs gets all jobs
	GetJobs() []Job
	// JobStorer gets job storer
	JobStorer() JobStorer
	// Job add a job to queue
	Job(job Job, args ...[]Arg) PendingJob
	// Register register jobs
	Register(jobs []Job)
	// Worker create a queue worker
	Worker(payloads ...Args) Worker
}

type Worker interface {
	Run() error
	Shutdown() error
}

type Args struct {
	// Specify connection
	Connection string
	// Specify queue
	Queue string
	// Concurrent num
	Concurrent int
	// Tries maximum attempts
	Tries int
}

type Arg struct {
	Value any    `json:"value"`
	Type  string `json:"type"`
}
