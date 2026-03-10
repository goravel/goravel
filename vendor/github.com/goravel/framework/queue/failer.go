package queue

import (
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/database/db"
	"github.com/goravel/framework/contracts/foundation"
	contractsqueue "github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/queue/models"
	"github.com/goravel/framework/queue/utils"
	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/support/convert"
)

type Failer struct {
	query db.Query
	queue contractsqueue.Queue
	json  foundation.Json
}

func NewFailer(config config.Config, db db.DB, queue contractsqueue.Queue, json foundation.Json) *Failer {
	failedDatabase := config.GetString("queue.failed.database")
	failedTable := config.GetString("queue.failed.table")

	return &Failer{query: db.Connection(failedDatabase).Table(failedTable), queue: queue, json: json}
}

func (r *Failer) All() ([]contractsqueue.FailedJob, error) {
	var modelFailedJobs []models.FailedJob
	if err := r.query.Get(&modelFailedJobs); err != nil {
		return nil, err
	}

	return r.modelFailedJobsToFailedJobs(modelFailedJobs), nil
}

func (r *Failer) Get(connection, queue string, uuids []string) ([]contractsqueue.FailedJob, error) {
	query := r.query

	if connection != "" {
		query = query.Where("connection", connection)
	}

	if queue != "" {
		query = query.Where("queue", queue)
	}

	if len(uuids) > 0 {
		query = query.WhereIn("uuid", convert.ToAnySlice(uuids))
	}

	var modelFailedJobs []models.FailedJob
	if err := query.Get(&modelFailedJobs); err != nil {
		return nil, err
	}

	return r.modelFailedJobsToFailedJobs(modelFailedJobs), nil
}

func (r *Failer) modelFailedJobsToFailedJobs(modelFailedJobs []models.FailedJob) []contractsqueue.FailedJob {
	var failedJobs []contractsqueue.FailedJob
	for _, modelFailedJob := range modelFailedJobs {
		failedJobs = append(failedJobs, NewFailedJob(modelFailedJob, r.query, r.queue, r.json))
	}

	return failedJobs
}

type FailedJob struct {
	query     db.Query
	queue     contractsqueue.Queue
	json      foundation.Json
	failedJob models.FailedJob
}

func NewFailedJob(failedJob models.FailedJob, query db.Query, queue contractsqueue.Queue, json foundation.Json) *FailedJob {
	return &FailedJob{failedJob: failedJob, query: query, queue: queue, json: json}
}

func (r *FailedJob) Connection() string {
	return r.failedJob.Connection
}

func (r *FailedJob) Queue() string {
	return r.failedJob.Queue
}

func (r *FailedJob) Retry() error {
	connection, err := r.queue.Connection(r.failedJob.Connection)
	if err != nil {
		return err
	}

	task, err := utils.JsonToTask(r.failedJob.Payload, r.queue.JobStorer(), r.json)
	if err != nil {
		return err
	}

	if err := connection.Push(task, r.failedJob.Queue); err != nil {
		return err
	}

	_, err = r.query.Where("id", r.failedJob.ID).Delete()

	return err
}

func (r *FailedJob) FailedAt() *carbon.DateTime {
	return r.failedJob.FailedAt
}

func (r *FailedJob) Signature() string {
	task, err := utils.JsonToTask(r.failedJob.Payload, r.queue.JobStorer(), r.json)
	if err != nil {
		return ""
	}

	return task.Job.Signature()
}

func (r *FailedJob) UUID() string {
	return r.failedJob.UUID
}
