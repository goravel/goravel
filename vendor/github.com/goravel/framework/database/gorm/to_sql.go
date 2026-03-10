package gorm

import (
	"gorm.io/gorm"

	"github.com/goravel/framework/contracts/log"
	"github.com/goravel/framework/support/database"
)

type ToSql struct {
	log   log.Log
	query *Query
	raw   bool
}

func NewToSql(query *Query, log log.Log, raw bool) *ToSql {
	return &ToSql{
		log:   log,
		query: query,
		raw:   raw,
	}
}

func (r *ToSql) Count() string {
	query := buildSelectForCount(r.query)
	var count int64

	return r.sql(query.instance.Session(&gorm.Session{DryRun: true}).Count(&count))
}

func (r *ToSql) Create(value any) string {
	query := r.query.dest(value).buildConditions()

	return r.sql(query.instance.Session(&gorm.Session{DryRun: true}).Create(value))
}

func (r *ToSql) Delete(dests ...any) string {
	var (
		dest  any
		query *Query
	)

	if len(dests) > 0 {
		dest = dests[0]
		query = r.query.dest(dest).addGlobalScopes().buildConditions()
	} else {
		query = r.query.addGlobalScopes().buildConditions()
	}

	return r.sql(query.instance.Session(&gorm.Session{DryRun: true}).Delete(dest))
}

func (r *ToSql) Find(dest any, conds ...any) string {
	query := r.query.dest(dest).addGlobalScopes().buildConditions()

	return r.sql(query.instance.Session(&gorm.Session{DryRun: true}).Find(dest, conds...))
}

func (r *ToSql) First(dest any) string {
	query := r.query.dest(dest).addGlobalScopes().buildConditions()

	return r.sql(query.instance.Session(&gorm.Session{DryRun: true}).First(dest))
}

func (r *ToSql) ForceDelete(dests ...any) string {
	var (
		dest  any
		query *Query
	)

	if len(dests) > 0 {
		dest = dests[0]
		query = r.query.dest(dest).addGlobalScopes().buildConditions()
	} else {
		query = r.query.addGlobalScopes().buildConditions()
	}

	return r.sql(query.instance.Session(&gorm.Session{DryRun: true}).Unscoped().Delete(dest))
}

func (r *ToSql) Get(dest any) string {
	query := r.query.dest(dest).addGlobalScopes().buildConditions()

	return r.sql(query.instance.Session(&gorm.Session{DryRun: true}).Find(dest))
}

func (r *ToSql) Pluck(column string, dest any) string {
	query := r.query.addGlobalScopes().buildConditions()

	return r.sql(query.instance.Session(&gorm.Session{DryRun: true}).Pluck(column, dest))
}

func (r *ToSql) Save(dest any) string {
	id := database.GetID(dest)
	update := id != nil

	var query *Query
	if update {
		query = r.query.dest(dest).addGlobalScopes().buildConditions()
	} else {
		query = r.query.dest(dest).buildConditions()
	}

	return r.sql(query.instance.Session(&gorm.Session{DryRun: true}).Save(dest))
}

func (r *ToSql) Sum(column string, dest any) string {
	query := r.query.addGlobalScopes().buildConditions()

	return r.sql(query.instance.Session(&gorm.Session{DryRun: true}).Select("SUM(" + column + ")").Find(dest))
}

func (r *ToSql) Update(column any, value ...any) string {
	query := r.query.addGlobalScopes().buildConditions()
	if _, ok := column.(string); !ok && len(value) > 0 {
		return ""
	}

	if c, ok := column.(string); ok && len(value) > 0 {
		query.instance.Statement.Dest = map[string]any{c: value[0]}
	}
	if len(value) == 0 {
		query.instance.Statement.Dest = column
	}

	return r.sql(query.instance.Session(&gorm.Session{DryRun: true}).Updates(query.instance.Statement.Dest))
}

func (r *ToSql) sql(db *gorm.DB) string {
	sql := db.Statement.SQL.String()
	if db.Statement.Error != nil {
		r.log.Errorf("failed to get sql: %v", db.Statement.Error)
	}
	if !r.raw {
		return sql
	}

	return r.query.instance.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Raw(sql, db.Statement.Vars...)
	})
}
