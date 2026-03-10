package models

import "github.com/goravel/framework/support/carbon"

type Job struct {
	ReservedAt  *carbon.DateTime `db:"reserved_at"`
	AvailableAt *carbon.DateTime `db:"available_at"`
	CreatedAt   *carbon.DateTime `db:"created_at"`
	Queue       string           `db:"queue"`
	Payload     string           `db:"payload"`
	ID          uint             `db:"id"`
	Attempts    int              `db:"attempts"`
}

func (r *Job) Increment() int {
	r.Attempts++

	return r.Attempts
}

func (r *Job) Touch() *carbon.DateTime {
	r.ReservedAt = carbon.NewDateTime(carbon.Now())

	return r.ReservedAt
}

type FailedJob struct {
	FailedAt   *carbon.DateTime `db:"failed_at"`
	UUID       string           `db:"uuid"`
	Connection string           `db:"connection"`
	Queue      string           `db:"queue"`
	Payload    string           `db:"payload"`
	Exception  string           `db:"exception"`
	ID         uint             `db:"id"`
}
