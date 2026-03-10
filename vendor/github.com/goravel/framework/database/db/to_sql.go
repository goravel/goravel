package db

import (
	"github.com/goravel/framework/contracts/database/db"
	"github.com/goravel/framework/errors"
)

type ToSql struct {
	query *Query
	raw   bool
}

func NewToSql(query *Query, raw bool) *ToSql {
	return &ToSql{query: query, raw: raw}
}

func (r *ToSql) Count() string {
	if err := buildSelectForCount(r.query); err != nil {
		return r.generate(r.query.readBuilder, "", nil, err)
	}

	sql, args, err := r.query.buildSelect()

	return r.generate(r.query.readBuilder, sql, args, err)
}

func (r *ToSql) Delete() string {
	sql, args, err := r.query.buildDelete()

	return r.generate(r.query.writeBuilder, sql, args, err)
}

func (r *ToSql) First() string {
	sql, args, err := r.query.buildSelect()

	return r.generate(r.query.readBuilder, sql, args, err)
}

func (r *ToSql) Get() string {
	sql, args, err := r.query.buildSelect()

	return r.generate(r.query.readBuilder, sql, args, err)
}

func (r *ToSql) Insert(data any) string {
	mapData, err := convertToSliceMap(data)
	if err != nil {
		return r.generate(r.query.writeBuilder, "", nil, err)
	}
	if len(mapData) == 0 {
		return r.generate(r.query.writeBuilder, "", nil, errors.DatabaseDataIsEmpty)
	}

	sql, args, err := r.query.buildInsert(mapData)

	return r.generate(r.query.writeBuilder, sql, args, err)
}

func (r *ToSql) Pluck(column string, dest any) string {
	r.query.conditions.Selects = []string{column}
	sql, args, err := r.query.buildSelect()

	return r.generate(r.query.readBuilder, sql, args, err)
}

func (r *ToSql) Update(column any, value ...any) string {
	columnStr, ok := column.(string)
	if ok {
		if len(value) != 1 {
			return r.generate(r.query.writeBuilder, "", nil, errors.DatabaseInvalidArgumentNumber.Args(len(value), "1"))
		}

		return r.Update(map[string]any{columnStr: value[0]})
	}

	mapData, err := convertToMap(column)
	if err != nil {
		return r.generate(r.query.writeBuilder, "", nil, err)
	}

	sql, args, err := r.query.buildUpdate(mapData)

	return r.generate(r.query.writeBuilder, sql, args, err)
}

func (r *ToSql) generate(builder db.CommonBuilder, sql string, args []any, err error) string {
	if err != nil {
		r.query.logger.Errorf(r.query.ctx, errors.DatabaseFailedToGetSql.Args(err).Error())

		return ""
	}

	if r.raw {
		return builder.Explain(sql, args...)
	}

	return sql
}
