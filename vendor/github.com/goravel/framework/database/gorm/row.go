package gorm

import (
	"github.com/goravel/framework/database/db"
)

type Row struct {
	err   error
	query *Query
	row   map[string]any
}

func (r *Row) Err() error {
	return r.err
}

func (r *Row) Scan(value any) error {
	row := db.NewRow(r.row, r.err)
	if err := row.Scan(value); err != nil {
		return err
	}

	for _, item := range r.query.conditions.with {
		// Need to new a query, avoid to clear the conditions
		query := r.query.new(r.query.instance)
		// The new query must be cleared
		query.clearConditions()
		if err := query.Load(value, item.query, item.args...); err != nil {
			return err
		}
	}

	return nil
}
