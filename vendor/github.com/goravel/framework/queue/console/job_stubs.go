package console

type JobStubs struct {
}

// Job Create a job.
func (receiver JobStubs) Job() string {
	return `package DummyPackage

type DummyJob struct {
}

// Signature The name and signature of the job.
func (receiver *DummyJob) Signature() string {
	return "DummySignature"
}

// Handle Execute the job.
func (receiver *DummyJob) Handle(args ...any) error {
	return nil
}
`
}
