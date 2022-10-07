package jobs

import "github.com/goravel/framework/facades"

type TestAsyncJob struct {
}

//Signature The name and signature of the job.
func (receiver *TestAsyncJob) Signature() string {
	return "test_async_job"
}

//Handle Execute the job.
func (receiver *TestAsyncJob) Handle(args ...interface{}) error {
	facades.Log.Infof("test_async_job: %s, %d", args[0], args[1])

	return nil
}
