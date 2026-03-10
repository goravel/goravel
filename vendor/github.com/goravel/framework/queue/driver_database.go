package queue

import (
	"fmt"
	"time"

	contractscache "github.com/goravel/framework/contracts/cache"
	contractsdb "github.com/goravel/framework/contracts/database/db"
	contractsfoundation "github.com/goravel/framework/contracts/foundation"
	contractsqueue "github.com/goravel/framework/contracts/queue"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/queue/models"
	"github.com/goravel/framework/queue/utils"
	"github.com/goravel/framework/support/carbon"
)

var (
	_ contractsqueue.Driver = &Database{}
)

type Database struct {
	cache     contractscache.Cache
	db        contractsdb.DB
	jobStorer contractsqueue.JobStorer
	json      contractsfoundation.Json

	jobsTable  string
	retryAfter int
}

func NewDatabase(
	config contractsqueue.Config,
	cache contractscache.Cache,
	db contractsdb.DB,
	jobStorer contractsqueue.JobStorer,
	json contractsfoundation.Json,
	connection string) (*Database, error) {
	if cache == nil {
		return nil, errors.CacheFacadeNotSet.SetModule(errors.ModuleQueue)
	}

	dbConnection := config.GetString(fmt.Sprintf("queue.connections.%s.connection", connection))
	if dbConnection == "" {
		return nil, errors.QueueInvalidDatabaseConnection.Args(connection)
	}

	return &Database{
		cache:     cache,
		db:        db.Connection(dbConnection),
		jobStorer: jobStorer,
		json:      json,

		jobsTable:  config.GetString(fmt.Sprintf("queue.connections.%s.table", connection), "jobs"),
		retryAfter: config.GetInt(fmt.Sprintf("queue.connections.%s.retry_after", connection), 60),
	}, nil
}

func (r *Database) Driver() string {
	return contractsqueue.DriverDatabase
}

func (r *Database) Pop(queue string) (contractsqueue.ReservedJob, error) {
	var job models.Job

	cacheLock := fmt.Sprintf("goravel:queue-database-%s:lock", queue)
	lock := r.cache.Lock(cacheLock, 1*time.Minute)
	if !lock.Block(1 * time.Minute) {
		return nil, errors.QueuePopIsLocked.Args(queue, cacheLock)
	}

	defer lock.Release()

	if err := r.db.Transaction(func(tx contractsdb.Tx) error {
		if err := tx.Table(r.jobsTable).LockForUpdate().Where("queue", queue).Where(func(q contractsdb.Query) contractsdb.Query {
			return q.Where(func(q1 contractsdb.Query) contractsdb.Query {
				return r.isAvailable(q1)
			})

			// TODO: Add the retry logic in another PR
			// .OrWhere(func(q1 contractsdb.Query) contractsdb.Query {
			// 	return r.isReservedButExpired(q1)
			// })
		}).OrderBy("id").First(&job); err != nil {
			return err
		}

		if job.ID == 0 {
			return errors.QueueDriverNoJobFound.Args(queue)
		}

		job.Increment()
		job.Touch()

		_, err := tx.Table(r.jobsTable).Where("id", job.ID).Update(map[string]any{
			"attempts":    job.Attempts,
			"reserved_at": job.ReservedAt,
		})
		if err != nil {
			return errors.QueueFailedToReserveJob.Args(job, err)
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return NewDatabaseReservedJob(&job, r.db, r.jobStorer, r.json, r.jobsTable)
}

func (r *Database) Push(task contractsqueue.Task, queue string) error {
	now := carbon.Now()
	availableAt := carbon.NewDateTime(now)

	if !task.Delay.IsZero() {
		availableAt = carbon.NewDateTime(carbon.FromStdTime(task.Delay))
	}

	payload, err := utils.TaskToJson(task, r.json)
	if err != nil {
		return err
	}

	job := models.Job{
		Queue:       queue,
		Payload:     payload,
		AvailableAt: availableAt,
		CreatedAt:   carbon.NewDateTime(now.Copy()),
	}

	result, err := r.db.Table(r.jobsTable).Insert(&job)
	if err != nil {
		return errors.QueueFailedToInsertJobToDatabase.Args(job, err)
	}
	if result.RowsAffected == 0 {
		return errors.QueueFailedToInsertJobToDatabase.Args(job, nil)
	}

	return nil
}

func (r *Database) isAvailable(query contractsdb.Query) contractsdb.Query {
	return query.WhereNull("reserved_at").Where("available_at <= ?", carbon.Now())
}

// TODO: Add the retry logic in another PR
// func (r *Database) isReservedButExpired(query contractsdb.Query) contractsdb.Query {
// 	return query.Where("reserved_at", "<=", carbon.Now().AddSeconds(r.retryAfter))
// }
