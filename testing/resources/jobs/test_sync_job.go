package jobs

import "github.com/goravel/framework/facades"

type TestSyncJob struct {
}

//Signature The name and signature of the job.
func (receiver *TestSyncJob) Signature() string {
	return "test_sync_job"
}

//Handle Execute the job.
func (receiver *TestSyncJob) Handle(args ...interface{}) error {
	facades.Log.Infof("test_sync_job: %s, %d", args[0], args[1])

	return nil
}
