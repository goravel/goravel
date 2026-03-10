package gorm

import (
	"context"
	"database/sql"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/spf13/cast"
	gormio "gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/goravel/framework/contracts/config"
	contractsdatabase "github.com/goravel/framework/contracts/database"
	contractsdb "github.com/goravel/framework/contracts/database/db"
	contractsdriver "github.com/goravel/framework/contracts/database/driver"
	contractsorm "github.com/goravel/framework/contracts/database/orm"
	"github.com/goravel/framework/contracts/log"
	"github.com/goravel/framework/database/db"
	databasedriver "github.com/goravel/framework/database/driver"
	"github.com/goravel/framework/database/utils"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/collect"
	"github.com/goravel/framework/support/database"
	"github.com/goravel/framework/support/deep"
	"github.com/goravel/framework/support/str"
)

const Associations = clause.Associations

type Query struct {
	config          config.Config
	ctx             context.Context
	grammar         contractsdriver.Grammar
	log             log.Log
	instance        *gormio.DB
	queries         map[string]*Query
	modelToObserver []contractsorm.ModelToObserver
	conditions      Conditions
	dbConfig        contractsdatabase.Config
	mutex           sync.Mutex
}

func NewQuery(
	ctx context.Context,
	config config.Config,
	dbConfig contractsdatabase.Config,
	db *gormio.DB,
	grammar contractsdriver.Grammar,
	log log.Log,
	modelToObserver []contractsorm.ModelToObserver,
	conditions *Conditions,
) *Query {
	queryImpl := &Query{
		config:          config,
		ctx:             ctx,
		dbConfig:        dbConfig,
		instance:        db,
		grammar:         grammar,
		log:             log,
		modelToObserver: modelToObserver,
		queries:         make(map[string]*Query),
	}

	if conditions != nil {
		queryImpl.conditions = *conditions
	}

	return queryImpl
}

func BuildQuery(ctx context.Context, config config.Config, connection string, log log.Log, modelToObserver []contractsorm.ModelToObserver) (*Query, contractsdatabase.Config, error) {
	driverCallback, exist := config.Get(fmt.Sprintf("database.connections.%s.via", connection)).(func() (contractsdriver.Driver, error))
	if !exist {
		return nil, contractsdatabase.Config{}, errors.DatabaseConfigNotFound
	}

	driver, err := driverCallback()
	if err != nil {
		return nil, contractsdatabase.Config{}, err
	}

	pool := driver.Pool()
	logger := db.NewLogger(config, log).ToGorm()
	gorm, err := databasedriver.BuildGorm(config, logger, pool, connection)
	if err != nil {
		return nil, pool.Writers[0], err
	}

	return NewQuery(ctx, config, pool.Writers[0], gorm, driver.Grammar(), log, modelToObserver, nil), pool.Writers[0], nil
}

func (r *Query) Association(association string) contractsorm.Association {
	query := r.buildConditions()

	return query.instance.Association(association)
}

// DEPRECATED Use BeginTransaction instead.
func (r *Query) Begin() (contractsorm.Query, error) {
	return r.BeginTransaction()
}

func (r *Query) BeginTransaction() (contractsorm.Query, error) {
	if r.InTransaction() {
		return r, nil
	}

	tx := r.instance.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return r.new(tx), nil
}

func (r *Query) Commit() error {
	return r.instance.Commit().Error
}

func (r *Query) Count() (int64, error) {
	query := buildSelectForCount(r)

	var count int64
	if err := query.instance.Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Query) Create(value any) error {
	query := r.dest(value).buildConditions()

	if len(query.instance.Statement.Selects) > 0 && len(query.instance.Statement.Omits) > 0 {
		return errors.OrmQuerySelectAndOmitsConflict
	}

	if len(query.instance.Statement.Selects) > 0 {
		return query.selectCreate(value)
	}

	if len(query.instance.Statement.Omits) > 0 {
		return query.omitCreate(value)
	}

	return query.create(value)
}

func (r *Query) Cursor() chan contractsdb.Row {
	with := r.conditions.with
	query := r.addGlobalScopes().buildConditions()
	r.conditions.with = with

	cursorChan := make(chan contractsdb.Row)
	go func() {
		var (
			err  error
			rows *sql.Rows
		)
		defer func() {
			if err != nil {
				cursorChan <- &Row{query: r, err: err}
			}

			close(cursorChan)
		}()

		if rows, err = query.instance.Rows(); err != nil {
			return
		}
		defer errors.Ignore(rows.Close)

		for rows.Next() {
			val := make(map[string]any)
			if err = query.instance.ScanRows(rows, val); err != nil {
				return
			}
			cursorChan <- &Row{query: r, row: val}
		}

		if err = rows.Err(); err != nil {
			return
		}
	}()
	return cursorChan
}

func (r *Query) DB() (*sql.DB, error) {
	return r.instance.DB()
}

func (r *Query) Delete(dests ...any) (*contractsdb.Result, error) {
	var (
		dest  any
		query *Query
	)

	if len(dests) > 0 {
		dest = dests[0]
		query = r.dest(dest).addGlobalScopes().buildConditions()
	} else {
		query = r.addGlobalScopes().buildConditions()
	}

	if err := query.deleting(dest); err != nil {
		return nil, err
	}

	res := query.instance.Delete(dest)
	if res.Error != nil {
		return nil, res.Error
	}

	if err := query.deleted(dest); err != nil {
		return nil, err
	}

	return &contractsdb.Result{
		RowsAffected: res.RowsAffected,
	}, nil
}

func (r *Query) Distinct(columns ...string) contractsorm.Query {
	if len(columns) == 0 {
		columns = []string{"*"}
	}

	query := r.Select(columns...).(*Query)
	query.conditions.distinct = true

	return r.setConditions(query.conditions)
}

func (r *Query) Driver() string {
	return r.dbConfig.Driver
}

func (r *Query) Exec(sql string, values ...any) (*contractsdb.Result, error) {
	query := r.buildConditions()
	result := query.instance.Exec(sql, values...)

	return &contractsdb.Result{
		RowsAffected: result.RowsAffected,
	}, result.Error
}

func (r *Query) Exists() (bool, error) {
	query := r.addGlobalScopes().buildConditions()

	var exists bool
	err := query.instance.Select("1").Limit(1).Find(&exists).Error
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *Query) Find(dest any, conds ...any) error {
	query := r.dest(dest).addGlobalScopes().buildConditions()

	if err := filterFindConditions(conds...); err != nil {
		return err
	}
	if err := query.instance.Find(dest, conds...).Error; err != nil {
		return err
	}

	return query.retrieved(dest)
}

func (r *Query) FindOrFail(dest any, conds ...any) error {
	query := r.dest(dest).addGlobalScopes().buildConditions()

	if err := filterFindConditions(conds...); err != nil {
		return err
	}

	res := query.instance.Find(dest, conds...)
	if err := res.Error; err != nil {
		return err
	}

	if res.RowsAffected == 0 {
		return errors.OrmRecordNotFound
	}

	return query.retrieved(dest)
}

func (r *Query) First(dest any) error {
	query := r.dest(dest).addGlobalScopes().buildConditions()

	res := query.instance.First(dest)
	if res.Error != nil {
		if errors.Is(res.Error, gormio.ErrRecordNotFound) {
			return nil
		}

		return res.Error
	}

	return query.retrieved(dest)
}

func (r *Query) FirstOr(dest any, callback func() error) error {
	query := r.dest(dest).addGlobalScopes().buildConditions()

	if err := query.instance.First(dest).Error; err != nil {
		if errors.Is(err, gormio.ErrRecordNotFound) {
			return callback()
		}

		return err
	}

	return query.retrieved(dest)
}

func (r *Query) FirstOrCreate(dest any, conds ...any) error {
	query := r.dest(dest).addGlobalScopes().buildConditions()

	if len(conds) == 0 {
		return errors.OrmQueryConditionRequired
	}

	var res *gormio.DB
	if len(conds) > 1 {
		res = query.instance.Attrs(conds[1]).FirstOrInit(dest, conds[0])
	} else {
		res = query.instance.FirstOrInit(dest, conds[0])
	}

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected > 0 {
		return query.retrieved(dest)
	}

	return r.Create(dest)
}

func (r *Query) FirstOrFail(dest any) error {
	query := r.dest(dest).addGlobalScopes().buildConditions()

	if err := query.instance.First(dest).Error; err != nil {
		if errors.Is(err, gormio.ErrRecordNotFound) {
			return errors.OrmRecordNotFound
		}

		return err
	}

	return query.retrieved(dest)
}

func (r *Query) FirstOrNew(dest any, attributes any, values ...any) error {
	query := r.dest(dest).addGlobalScopes().buildConditions()

	var res *gormio.DB
	if len(values) > 0 {
		res = query.instance.Attrs(values[0]).FirstOrInit(dest, attributes)
	} else {
		res = query.instance.FirstOrInit(dest, attributes)
	}

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected > 0 {
		return query.retrieved(dest)
	}

	return nil
}

func (r *Query) ForceDelete(dests ...any) (*contractsdb.Result, error) {
	var (
		dest  any
		query *Query
	)

	if len(dests) > 0 {
		dest = dests[0]
		query = r.dest(dest).addGlobalScopes().buildConditions()
	} else {
		query = r.addGlobalScopes().buildConditions()
	}

	if err := query.forceDeleting(dest); err != nil {
		return nil, err
	}

	res := query.instance.Unscoped().Delete(dest)
	if res.Error != nil {
		return nil, res.Error
	}

	if res.RowsAffected > 0 {
		if err := query.forceDeleted(dest); err != nil {
			return nil, err
		}
	}

	return &contractsdb.Result{
		RowsAffected: res.RowsAffected,
	}, res.Error
}

func (r *Query) Get(dest any) error {
	return r.Find(dest)
}

// DEPRECATED Use GroupBy instead.
func (r *Query) Group(column string) contractsorm.Query {
	return r.GroupBy(column)
}

func (r *Query) GroupBy(column ...string) contractsorm.Query {
	if len(column) == 0 {
		return r
	}

	conditions := r.conditions
	conditions.groupBy = column

	return r.setConditions(conditions)
}

func (r *Query) Having(query any, args ...any) contractsorm.Query {
	conditions := r.conditions
	conditions.having = &contractsdriver.Having{
		Query: query,
		Args:  args,
	}

	return r.setConditions(conditions)
}

func (r *Query) Join(query string, args ...any) contractsorm.Query {
	conditions := r.conditions
	conditions.join = deep.Append(conditions.join, contractsdriver.Join{
		Query: query,
		Args:  args,
	})

	return r.setConditions(conditions)
}

func (r *Query) Limit(limit int) contractsorm.Query {
	conditions := r.conditions
	conditions.limit = &limit

	return r.setConditions(conditions)
}

func (r *Query) Load(model any, relation string, args ...any) error {
	if relation == "" {
		return errors.OrmQueryEmptyRelation
	}

	destType := reflect.TypeOf(model)
	if destType.Kind() != reflect.Pointer {
		return errors.OrmQueryModelNotPointer
	}

	if id := database.GetID(model); id == nil {
		return errors.OrmQueryEmptyId
	}

	copyDest := copyStruct(model)
	err := r.With(relation, args...).Find(model)

	relationRoot := relation
	if dotIndex := strings.Index(relation, "."); dotIndex > 0 {
		relationRoot = relation[:dotIndex]
	}

	t := destType.Elem()
	v := reflect.ValueOf(model).Elem()
	for i := 0; i < t.NumField(); i++ {
		if !t.Field(i).IsExported() {
			continue
		}
		if t.Field(i).Name != relationRoot {
			v.Field(i).Set(copyDest.Field(i))
		}
	}

	return err
}

func (r *Query) LoadMissing(model any, relation string, args ...any) error {
	destType := reflect.TypeOf(model)
	if destType.Kind() != reflect.Pointer {
		return errors.OrmQueryModelNotPointer
	}

	t := reflect.TypeOf(model).Elem()
	v := reflect.ValueOf(model).Elem()
	for i := 0; i < t.NumField(); i++ {
		if !t.Field(i).IsExported() {
			continue
		}
		if t.Field(i).Name == relation {
			var id any
			if v.Field(i).Kind() == reflect.Pointer {
				if !v.Field(i).IsNil() {
					id = database.GetIDByReflect(v.Field(i).Type().Elem(), v.Field(i).Elem())
				}
			} else if v.Field(i).Kind() == reflect.Slice {
				if v.Field(i).Len() > 0 {
					return nil
				}
			} else {
				id = database.GetIDByReflect(v.Field(i).Type(), v.Field(i))
			}
			if cast.ToString(id) != "" {
				return nil
			}
		}
	}

	return r.Load(model, relation, args...)
}

func (r *Query) LockForUpdate() contractsorm.Query {
	conditions := r.conditions
	conditions.lockForUpdate = true

	return r.setConditions(conditions)
}

func (r *Query) Model(value any) contractsorm.Query {
	if r.conditions.model != nil {
		return r
	}

	conditions := r.conditions
	conditions.model = value

	return r.setConditions(conditions)
}

func (r *Query) Observe(model any, observer contractsorm.Observer) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.modelToObserver = append(r.modelToObserver, contractsorm.ModelToObserver{
		Model:    model,
		Observer: observer,
	})
}

func (r *Query) Offset(offset int) contractsorm.Query {
	conditions := r.conditions
	conditions.offset = &offset

	return r.setConditions(conditions)
}

func (r *Query) Omit(columns ...string) contractsorm.Query {
	conditions := r.conditions
	conditions.omit = columns

	return r.setConditions(conditions)
}

// DEPRECATED: Use OrderByRaw instead
func (r *Query) Order(value any) contractsorm.Query {
	return r.OrderByRaw(fmt.Sprintf("%s", value))
}

func (r *Query) OrderBy(column string, direction ...string) contractsorm.Query {
	var orderDirection string
	if len(direction) > 0 {
		orderDirection = direction[0]
	} else {
		orderDirection = "ASC"
	}
	return r.OrderByRaw(fmt.Sprintf("%s %s", column, orderDirection))
}

func (r *Query) OrderByDesc(column string) contractsorm.Query {
	return r.OrderByRaw(fmt.Sprintf("%s DESC", column))
}

func (r *Query) OrderByRaw(raw string) contractsorm.Query {
	var rawAny any = raw

	conditions := r.conditions
	conditions.order = deep.Append(r.conditions.order, rawAny)

	return r.setConditions(conditions)
}

func (r *Query) Instance() *gormio.DB {
	return r.instance
}

func (r *Query) InRandomOrder() contractsorm.Query {
	return r.OrderByRaw(r.grammar.CompileRandomOrderForGorm())
}

func (r *Query) InTransaction() bool {
	committer, ok := r.Instance().Statement.ConnPool.(gormio.TxCommitter)

	return ok && committer != nil
}

func (r *Query) OrWhere(query any, args ...any) contractsorm.Query {
	return r.addWhere(contractsdriver.Where{
		Query: query,
		Args:  args,
		Or:    true,
	})
}

func (r *Query) Paginate(page, limit int, dest any, total *int64) error {
	offset := (page - 1) * limit
	if total != nil {
		if r.conditions.table == nil && r.conditions.model == nil {
			count, err := r.Model(dest).Count()
			if err != nil {
				return err
			}
			*total = count
		} else {
			count, err := r.Count()
			if err != nil {
				return err
			}
			*total = count
		}
	}

	return r.Offset(offset).Limit(limit).Find(dest)
}

func (r *Query) Pluck(column string, dest any) error {
	query := r.addGlobalScopes().buildConditions()

	return query.instance.Pluck(column, dest).Error
}

func (r *Query) Raw(sql string, values ...any) contractsorm.Query {
	return r.new(r.instance.Raw(sql, values...))
}

func (r *Query) Restore(dests ...any) (*contractsdb.Result, error) {
	var (
		dest  any
		query *Query
	)

	if len(dests) > 0 {
		dest = dests[0]
		query = r.dest(dest).addGlobalScopes().buildConditions()
	} else {
		query = r.addGlobalScopes().buildConditions()
	}

	var (
		deletedAtColumnName string

		tx = query.instance
	)
	if dest != nil {
		deletedAtColumnName = getDeletedAtColumn(dest)
		tx = query.instance.Model(dest)
	} else if query.conditions.model != nil {
		deletedAtColumnName = getDeletedAtColumn(query.conditions.model)
	}
	if deletedAtColumnName == "" {
		return nil, errors.OrmDeletedAtColumnNotFound
	}

	if err := r.restoring(dest); err != nil {
		return nil, err
	}

	res := tx.Update(deletedAtColumnName, nil)
	if res.Error != nil {
		return nil, res.Error
	}

	if err := r.restored(dest); err != nil {
		return nil, err
	}

	return &contractsdb.Result{
		RowsAffected: res.RowsAffected,
	}, res.Error
}

func (r *Query) Rollback() error {
	return r.instance.Rollback().Error
}

func (r *Query) Save(dest any) error {
	if len(r.conditions.selectColumns) > 0 && len(r.conditions.omit) > 0 {
		return errors.OrmQuerySelectAndOmitsConflict
	}

	id := database.GetID(dest)
	update := id != nil

	var query *Query
	if update {
		query = r.dest(dest).addGlobalScopes().buildConditions()
	} else {
		query = r.dest(dest).buildConditions()
	}

	if err := query.saving(dest); err != nil {
		return err
	}
	if update {
		if err := query.updating(dest); err != nil {
			return err
		}
	} else {
		if err := query.creating(dest); err != nil {
			return err
		}
	}

	if len(query.instance.Statement.Selects) > 0 {
		if err := query.selectSave(dest); err != nil {
			return err
		}
	} else if len(query.instance.Statement.Omits) > 0 {
		if err := query.omitSave(dest); err != nil {
			return err
		}
	} else {
		if err := query.save(dest); err != nil {
			return err
		}
	}

	if update {
		if err := query.updated(dest); err != nil {
			return err
		}
	} else {
		if err := query.created(dest); err != nil {
			return err
		}
	}
	if err := query.saved(dest); err != nil {
		return err
	}

	return nil
}

func (r *Query) SaveQuietly(value any) error {
	return r.WithoutEvents().Save(value)
}

func (r *Query) Scan(dest any) error {
	query := r.dest(dest).buildConditions()

	return query.instance.Scan(dest).Error
}

func (r *Query) Scopes(funcs ...func(contractsorm.Query) contractsorm.Query) contractsorm.Query {
	conditions := r.conditions
	conditions.scopes = deep.Append(r.conditions.scopes, funcs...)

	return r.setConditions(conditions)
}

func (r *Query) Select(columns ...string) contractsorm.Query {
	conditions := r.conditions
	conditions.selectColumns = append(conditions.selectColumns, columns...)
	conditions.selectColumns = collect.Unique(conditions.selectColumns)

	// * may be added along with other columns, remove it.
	if len(conditions.selectColumns) > 1 {
		conditions.selectColumns = collect.Filter(conditions.selectColumns, func(column string, _ int) bool {
			return column != "*"
		})
	}

	return r.setConditions(conditions)
}

func (r *Query) SelectRaw(query any, args ...any) contractsorm.Query {
	conditions := r.conditions
	conditions.selectRaw = &Select{
		query: query,
		args:  args,
	}

	return r.setConditions(conditions)
}

func (r *Query) WithContext(ctx context.Context) contractsorm.Query {
	instance := r.instance.WithContext(ctx)

	return NewQuery(ctx, r.config, r.dbConfig, instance, r.grammar, r.log, r.modelToObserver, nil)
}

func (r *Query) SharedLock() contractsorm.Query {
	conditions := r.conditions
	conditions.sharedLock = true

	return r.setConditions(conditions)
}

func (r *Query) Sum(column string, dest any) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr {
		return errors.DatabaseUnsupportedType.Args(destValue.Kind(), "pointer")
	}

	query := r.addGlobalScopes().buildConditions()
	return query.instance.Select("SUM(" + column + ")").Row().Scan(dest)
}

func (r *Query) Avg(column string, dest any) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr {
		return errors.DatabaseUnsupportedType.Args(destValue.Kind(), "pointer")
	}

	query := r.addGlobalScopes().buildConditions()
	return query.instance.Select("AVG(" + column + ")").Row().Scan(dest)
}

func (r *Query) Min(column string, dest any) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr {
		return errors.DatabaseUnsupportedType.Args(destValue.Kind(), "pointer")
	}

	query := r.addGlobalScopes().buildConditions()
	return query.instance.Select("MIN(" + column + ")").Row().Scan(dest)
}

func (r *Query) Max(column string, dest any) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr {
		return errors.DatabaseUnsupportedType.Args(destValue.Kind(), "pointer")
	}

	query := r.addGlobalScopes().buildConditions()
	return query.instance.Select("MAX(" + column + ")").Row().Scan(dest)
}

func (r *Query) Table(name string, args ...any) contractsorm.Query {
	conditions := r.conditions
	conditions.table = &Table{
		name: r.dbConfig.Prefix + name,
		args: args,
	}

	return r.setConditions(conditions)
}

func (r *Query) ToSql() contractsorm.ToSql {
	return NewToSql(r.setConditions(r.conditions), r.log, false)
}

func (r *Query) ToRawSql() contractsorm.ToSql {
	return NewToSql(r.setConditions(r.conditions), r.log, true)
}

func (r *Query) Update(column any, value ...any) (*contractsdb.Result, error) {
	query := r.addGlobalScopes().buildConditions()

	if _, ok := column.(string); !ok && len(value) > 0 {
		return nil, errors.OrmQueryInvalidParameter
	}

	var singleUpdate bool
	model := query.instance.Statement.Model
	if model != nil {
		id := database.GetID(model)
		singleUpdate = id != nil
	}

	if c, ok := column.(string); ok && len(value) > 0 {
		query.instance.Statement.Dest = map[string]any{c: value[0]}
	}
	if len(value) == 0 {
		query.instance.Statement.Dest = column
	}

	if singleUpdate {
		if err := query.saving(query.instance.Statement.Dest); err != nil {
			return nil, err
		}
		if err := query.updating(query.instance.Statement.Dest); err != nil {
			return nil, err
		}
	}

	res, err := query.update(query.instance.Statement.Dest)

	if singleUpdate && err == nil {
		if err := query.updated(query.instance.Statement.Dest); err != nil {
			return nil, err
		}
		if err := query.saved(query.instance.Statement.Dest); err != nil {
			return nil, err
		}
	}

	return res, err
}

func (r *Query) UpdateOrCreate(dest any, attributes any, values any) error {
	query := r.dest(dest).addGlobalScopes().buildConditions()

	res := query.instance.Assign(values).FirstOrInit(dest, attributes)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected > 0 {
		return r.Save(dest)
	}

	return r.Create(dest)
}

func (r *Query) Where(query any, args ...any) contractsorm.Query {
	return r.addWhere(contractsdriver.Where{
		Query: query,
		Args:  args,
	})
}

func (r *Query) WhereAll(columns []string, args ...any) contractsorm.Query {
	op, value, err := utils.PrepareWhereOperatorAndValue(args...)
	if err != nil {
		query := r.new(r.instance.Session(&gormio.Session{}))
		_ = query.instance.AddError(err)
		return query
	}

	var conditions []string
	var conditionArgs []any
	for _, column := range columns {
		conditions = append(conditions, fmt.Sprintf("%s %v ?", column, op))
		conditionArgs = append(conditionArgs, value)
	}

	query := strings.Join(conditions, " AND ")
	r = r.addWhere(contractsdriver.Where{
		Query: query,
		Args:  conditionArgs,
	}).(*Query)

	return r
}

func (r *Query) WhereAny(columns []string, args ...any) contractsorm.Query {
	op, value, err := utils.PrepareWhereOperatorAndValue(args...)
	if err != nil {
		query := r.new(r.instance.Session(&gormio.Session{}))
		_ = query.instance.AddError(err)
		return query
	}

	var conditions []string
	var conditionArgs []any
	for _, column := range columns {
		conditions = append(conditions, fmt.Sprintf("%s %v ?", column, op))
		conditionArgs = append(conditionArgs, value)
	}

	query := fmt.Sprintf("(%s)", strings.Join(conditions, " OR "))
	r = r.addWhere(contractsdriver.Where{
		Query: query,
		Args:  conditionArgs,
	}).(*Query)

	return r
}

func (r *Query) WhereIn(column string, values []any) contractsorm.Query {
	return r.Where(fmt.Sprintf("%s IN ?", column), values)
}

func (r *Query) WhereJsonContains(column string, value any) contractsorm.Query {
	return r.addWhere(contractsdriver.Where{
		Type:  contractsdriver.WhereTypeJsonContains,
		Query: column,
		Args:  []any{value},
	})
}

func (r *Query) WhereJsonContainsKey(column string) contractsorm.Query {
	return r.addWhere(contractsdriver.Where{
		Type:  contractsdriver.WhereTypeJsonContainsKey,
		Query: column,
	})
}

func (r *Query) WhereJsonDoesntContain(column string, value any) contractsorm.Query {
	return r.addWhere(contractsdriver.Where{
		Type:  contractsdriver.WhereTypeJsonContains,
		Query: column,
		Args:  []any{value},
		IsNot: true,
	})
}

func (r *Query) WhereJsonDoesntContainKey(column string) contractsorm.Query {
	return r.addWhere(contractsdriver.Where{
		Type:  contractsdriver.WhereTypeJsonContainsKey,
		Query: column,
		IsNot: true,
	})
}

func (r *Query) WhereJsonLength(column string, length int) contractsorm.Query {
	return r.addWhere(contractsdriver.Where{
		Type:  contractsdriver.WhereTypeJsonLength,
		Query: column,
		Args:  []any{length},
	})
}

func (r *Query) OrWhereIn(column string, values []any) contractsorm.Query {
	return r.OrWhere(fmt.Sprintf("%s IN ?", column), values)
}

func (r *Query) OrWhereJsonContains(column string, value any) contractsorm.Query {
	return r.addWhere(contractsdriver.Where{
		Type:  contractsdriver.WhereTypeJsonContains,
		Query: column,
		Args:  []any{value},
		Or:    true,
	})
}

func (r *Query) OrWhereJsonContainsKey(column string) contractsorm.Query {
	return r.addWhere(contractsdriver.Where{
		Type:  contractsdriver.WhereTypeJsonContainsKey,
		Query: column,
		Or:    true,
	})
}

func (r *Query) OrWhereJsonDoesntContain(column string, value any) contractsorm.Query {
	return r.addWhere(contractsdriver.Where{
		Type:  contractsdriver.WhereTypeJsonContains,
		Query: column,
		Args:  []any{value},
		IsNot: true,
		Or:    true,
	})
}

func (r *Query) OrWhereJsonDoesntContainKey(column string) contractsorm.Query {
	return r.addWhere(contractsdriver.Where{
		Type:  contractsdriver.WhereTypeJsonContainsKey,
		Query: column,
		IsNot: true,
		Or:    true,
	})
}

func (r *Query) OrWhereJsonLength(column string, length int) contractsorm.Query {
	return r.addWhere(contractsdriver.Where{
		Type:  contractsdriver.WhereTypeJsonLength,
		Query: column,
		Args:  []any{length},
		Or:    true,
	})
}

func (r *Query) WhereNotIn(column string, values []any) contractsorm.Query {
	return r.Where(fmt.Sprintf("%s NOT IN ?", column), values)
}

func (r *Query) OrWhereNotIn(column string, values []any) contractsorm.Query {
	return r.OrWhere(fmt.Sprintf("%s NOT IN ?", column), values)
}

func (r *Query) WhereBetween(column string, x, y any) contractsorm.Query {
	return r.Where(fmt.Sprintf("%s BETWEEN ? AND ?", column), x, y)
}

func (r *Query) WhereNotBetween(column string, x, y any) contractsorm.Query {
	return r.Where(fmt.Sprintf("%s NOT BETWEEN ? AND ?", column), x, y)
}

func (r *Query) OrWhereBetween(column string, x, y any) contractsorm.Query {
	return r.OrWhere(fmt.Sprintf("%s BETWEEN ? AND ?", column), x, y)
}

func (r *Query) OrWhereNotBetween(column string, x, y any) contractsorm.Query {
	return r.OrWhere(fmt.Sprintf("%s NOT BETWEEN ? AND ?", column), x, y)
}

func (r *Query) OrWhereNull(column string) contractsorm.Query {
	return r.OrWhere(fmt.Sprintf("%s IS NULL", column))
}

func (r *Query) WhereNone(columns []string, args ...any) contractsorm.Query {
	op, value, err := utils.PrepareWhereOperatorAndValue(args...)
	if err != nil {
		query := r.new(r.instance.Session(&gormio.Session{}))
		_ = query.instance.AddError(err)
		return query
	}

	var conditions []string
	var conditionArgs []any
	for _, column := range columns {
		if op == "=" {
			conditions = append(conditions, fmt.Sprintf("%s <> ?", column))
		} else {
			conditions = append(conditions, fmt.Sprintf("NOT (%s %v ?)", column, op))
		}
		conditionArgs = append(conditionArgs, value)
	}

	query := strings.Join(conditions, " AND ")
	r = r.addWhere(contractsdriver.Where{
		Query: query,
		Args:  conditionArgs,
	}).(*Query)

	return r
}

func (r *Query) WhereNotNull(column string) contractsorm.Query {
	return r.Where(fmt.Sprintf("%s IS NOT NULL", column))
}

func (r *Query) WhereNull(column string) contractsorm.Query {
	return r.Where(fmt.Sprintf("%s IS NULL", column))
}

func (r *Query) With(query string, args ...any) contractsorm.Query {
	conditions := r.conditions
	conditions.with = deep.Append(r.conditions.with, With{
		query: query,
		args:  args,
	})

	return r.setConditions(conditions)
}

func (r *Query) WithoutEvents() contractsorm.Query {
	conditions := r.conditions
	conditions.withoutEvents = true

	return r.setConditions(conditions)
}

func (r *Query) WithoutGlobalScopes(names ...string) contractsorm.Query {
	conditions := r.conditions

	if len(names) == 0 {
		names = []string{"*"}
	}
	conditions.withoutGlobalScopes = append(conditions.withoutGlobalScopes, names...)

	return r.setConditions(conditions)
}

func (r *Query) WithTrashed() contractsorm.Query {
	conditions := r.conditions
	conditions.withTrashed = true

	return r.setConditions(conditions)
}

func (r *Query) addGlobalScopes() *Query {
	if slices.Contains(r.conditions.withoutGlobalScopes, "*") {
		return r
	}

	var model any

	if r.conditions.model != nil {
		model = r.conditions.model
	} else if r.conditions.dest != nil {
		model = r.conditions.dest
	} else {
		return r
	}

	model, err := modelToStruct(model)
	if err != nil {
		return r
	}

	modelWithGlobalScopes, ok := model.(contractsorm.ModelWithGlobalScopes)
	if !ok {
		return r
	}

	nameToGlobalScopes := modelWithGlobalScopes.GlobalScopes()
	if len(nameToGlobalScopes) == 0 {
		return r
	}

	var globalScopes []func(contractsorm.Query) contractsorm.Query

	names := slices.Sorted(maps.Keys(nameToGlobalScopes))
	for _, name := range names {
		if slices.Contains(r.conditions.withoutGlobalScopes, name) {
			continue
		}

		globalScopes = append(globalScopes, nameToGlobalScopes[name])
	}

	return r.Scopes(globalScopes...).(*Query)
}

func (r *Query) addWhere(where contractsdriver.Where) contractsorm.Query {
	conditions := r.conditions
	conditions.where = deep.Append(conditions.where, where)

	return r.setConditions(conditions)
}

func (r *Query) buildConditions() *Query {
	query, err := r.refreshConnection()
	if err != nil {
		query = r.new(r.instance.Session(&gormio.Session{}))
		_ = query.instance.AddError(err)

		return query
	}

	query = query.buildModel()
	db := query.instance
	db = query.buildDistinct(db)
	db = query.buildGroup(db)
	db = query.buildHaving(db)
	db = query.buildJoin(db)
	db = query.buildLockForUpdate(db)
	db = query.buildLimit(db)
	db = query.buildOrder(db)
	db = query.buildOffset(db)
	db = query.buildOmit(db)
	db = query.buildScopes(db)
	db = query.buildSelectColumns(db)
	db = query.buildSharedLock(db)
	db = query.buildTable(db)
	db = query.buildWith(db)
	db = query.buildWithTrashed(db)
	db = query.buildWhere(db)

	return query.new(db)
}

func (r *Query) buildDistinct(db *gormio.DB) *gormio.DB {
	if !r.conditions.distinct {
		return db
	}

	db = db.Distinct()
	r.conditions.distinct = false

	return db
}

func (r *Query) buildGroup(db *gormio.DB) *gormio.DB {
	if len(r.conditions.groupBy) == 0 {
		return db
	}

	db = db.Group(strings.Join(r.conditions.groupBy, ", "))
	r.conditions.groupBy = nil

	return db
}

func (r *Query) buildHaving(db *gormio.DB) *gormio.DB {
	if r.conditions.having == nil {
		return db
	}

	db = db.Having(r.conditions.having.Query, r.conditions.having.Args...)
	r.conditions.having = nil

	return db
}

func (r *Query) buildJoin(db *gormio.DB) *gormio.DB {
	if r.conditions.join == nil {
		return db
	}

	for _, item := range r.conditions.join {
		db = db.Joins(item.Query, item.Args...)
	}

	r.conditions.join = nil

	return db
}

func (r *Query) buildLimit(db *gormio.DB) *gormio.DB {
	if r.conditions.limit == nil {
		return db
	}

	db = db.Limit(*r.conditions.limit)
	r.conditions.limit = nil

	return db
}

func (r *Query) buildLockForUpdate(db *gormio.DB) *gormio.DB {
	if !r.conditions.lockForUpdate {
		return db
	}

	lockForUpdate := r.grammar.CompileLockForUpdateForGorm()
	if lockForUpdate != nil {
		return db.Clauses(lockForUpdate)
	}

	r.conditions.lockForUpdate = false

	return db
}

func (r *Query) buildModel() *Query {
	if r.conditions.model == nil {
		return r
	}

	return r.new(r.instance.Model(r.conditions.model))
}

func (r *Query) buildOffset(db *gormio.DB) *gormio.DB {
	if r.conditions.offset == nil {
		return db
	}

	db = db.Offset(*r.conditions.offset)
	r.conditions.offset = nil

	return db
}

func (r *Query) buildOmit(db *gormio.DB) *gormio.DB {
	if len(r.conditions.omit) == 0 {
		return db
	}

	db = db.Omit(r.conditions.omit...)
	r.conditions.omit = nil

	return db
}

func (r *Query) buildOrder(db *gormio.DB) *gormio.DB {
	if len(r.conditions.order) == 0 {
		return db
	}

	for _, order := range r.conditions.order {
		db = db.Order(order)
	}

	r.conditions.order = nil

	return db
}

func (r *Query) buildSelectColumns(db *gormio.DB) *gormio.DB {
	if len(r.conditions.selectColumns) == 0 && r.conditions.selectRaw == nil {
		return db
	}

	if len(r.conditions.selectColumns) > 0 {
		var selectColumns []any
		for _, column := range r.conditions.selectColumns {
			selectColumns = append(selectColumns, column)
		}

		db = db.Select(selectColumns[0], selectColumns[1:]...)
	} else if r.conditions.selectRaw != nil {
		db = db.Select(r.conditions.selectRaw.query, r.conditions.selectRaw.args...)
	}

	r.conditions.selectColumns = nil
	r.conditions.selectRaw = nil

	return db
}

func (r *Query) buildScopes(db *gormio.DB) *gormio.DB {
	if len(r.conditions.scopes) == 0 {
		return db
	}

	var gormFuncs []func(*gormio.DB) *gormio.DB
	for _, scope := range r.conditions.scopes {
		currentScope := scope
		gormFuncs = append(gormFuncs, func(tx *gormio.DB) *gormio.DB {
			queryImpl := r.new(tx)
			query := currentScope(queryImpl)
			queryImpl = query.(*Query)
			queryImpl = queryImpl.buildConditions()

			return queryImpl.instance
		})
	}

	db = db.Scopes(gormFuncs...)
	r.conditions.scopes = nil

	return db
}

func (r *Query) buildSharedLock(db *gormio.DB) *gormio.DB {
	if !r.conditions.sharedLock {
		return db
	}

	sharedLock := r.grammar.CompileSharedLockForGorm()
	if sharedLock != nil {
		return db.Clauses(sharedLock)
	}

	r.conditions.sharedLock = false

	return db
}

func (r *Query) buildSubquery(sub func(contractsorm.Query) contractsorm.Query) *gormio.DB {
	db := r.instance.Session(&gormio.Session{NewDB: true, Initialized: true})
	queryImpl := NewQuery(r.ctx, r.config, r.dbConfig, db, r.grammar, r.log, r.modelToObserver, nil)
	query := sub(queryImpl)
	var ok bool
	if queryImpl, ok = query.(*Query); ok {
		return queryImpl.buildWhere(db)
	}

	return db
}

func (r *Query) buildTable(db *gormio.DB) *gormio.DB {
	if r.conditions.table == nil {
		return db
	}

	db = db.Table(r.conditions.table.name, r.conditions.table.args...)
	r.conditions.table = nil

	return db
}

func (r *Query) buildWhere(db *gormio.DB) *gormio.DB {
	if len(r.conditions.where) == 0 {
		return db
	}

	for _, item := range r.conditions.where {
		switch item.Type {
		case contractsdriver.WhereTypeJsonContains:
			query, value, err := r.grammar.CompileJsonContains(item.Query.(string), item.Args[0], item.IsNot)
			if err != nil {
				_ = r.instance.AddError(errors.OrmJsonContainsInvalidBinding.Args(err))
				continue
			}

			item.Query = query
			item.Args = value
		case contractsdriver.WhereTypeJsonContainsKey:
			item.Query = r.grammar.CompileJsonContainsKey(item.Query.(string), item.IsNot)

		case contractsdriver.WhereTypeJsonLength:
			segments := strings.SplitN(item.Query.(string), " ", 2)
			segments[0] = r.grammar.CompileJsonLength(segments[0])
			item.Query = r.buildWherePlaceholder(strings.Join(segments, " "), item.Args...)
		default:
			switch query := item.Query.(type) {
			case func(contractsorm.Query) contractsorm.Query:
				item.Query = r.buildSubquery(query)
				item.Args = nil
			case string:
				if strings.Contains(query, "->") {
					segments := strings.Split(query, " ")
					for i := range segments {
						if strings.Contains(segments[i], "->") {
							segments[i] = r.grammar.CompileJsonSelector(segments[i])
						}
					}
					item.Query = r.buildWherePlaceholder(strings.Join(segments, " "), item.Args...)
					item.Args = r.grammar.CompileJsonValues(item.Args...)
				}
			}
		}

		if item.Or {
			db = db.Or(item.Query, item.Args...)
		} else {
			db = db.Where(item.Query, item.Args...)
		}
	}

	r.conditions.where = nil

	return db
}

func (r *Query) buildWherePlaceholder(query string, args ...any) string {
	// if query does not contain a placeholder,it might be incorrectly quoted or treated as an expression
	// to avoid errors, append a manual placeholder
	if !strings.Contains(query, "?") && len(args) == 1 {
		if val := reflect.ValueOf(args[0]); val.Kind() == reflect.Slice || val.Kind() == reflect.Array {
			return query + " IN (?)"
		}
		return query + " = ?"
	}

	return query
}

func (r *Query) buildWith(db *gormio.DB) *gormio.DB {
	if len(r.conditions.with) == 0 {
		return db
	}

	for _, item := range r.conditions.with {
		isSet := false
		if len(item.args) == 1 {
			if arg, ok := item.args[0].(func(contractsorm.Query) contractsorm.Query); ok {
				newArgs := []any{
					func(tx *gormio.DB) *gormio.DB {
						queryImpl := NewQuery(r.ctx, r.config, r.dbConfig, tx, r.grammar, r.log, r.modelToObserver, nil)
						query := arg(queryImpl)
						queryImpl = query.(*Query)
						queryImpl = queryImpl.buildConditions()

						return queryImpl.instance
					},
				}

				db = db.Preload(item.query, newArgs...)
				isSet = true
			}
		}

		if !isSet {
			db = db.Preload(item.query, item.args...)
		}
	}

	r.conditions.with = nil

	return db
}

func (r *Query) buildWithTrashed(db *gormio.DB) *gormio.DB {
	if !r.conditions.withTrashed {
		return db
	}

	db = db.Unscoped()
	r.conditions.withTrashed = false

	return db
}

func (r *Query) clearConditions() {
	r.conditions = Conditions{}
}

func (r *Query) create(dest any) error {
	if err := r.saving(dest); err != nil {
		return err
	}
	if err := r.creating(dest); err != nil {
		return err
	}

	if err := r.instance.Omit(Associations).Create(dest).Error; err != nil {
		return err
	}

	if err := r.created(dest); err != nil {
		return err
	}
	if err := r.saved(dest); err != nil {
		return err
	}

	return nil
}

func (r *Query) created(dest any) error {
	if isSlice(dest) {
		return nil
	}

	return r.event(contractsorm.EventCreated, r.conditions.model, dest)
}

func (r *Query) creating(dest any) error {
	if isSlice(dest) {
		return nil
	}

	return r.event(contractsorm.EventCreating, r.conditions.model, dest)
}

func (r *Query) event(event contractsorm.EventType, model, dest any) error {
	if r.conditions.withoutEvents {
		return nil
	}

	instance := NewEvent(r, model, dest)

	if model != nil {
		if dispatchesEvents, exist := model.(contractsorm.DispatchesEvents); exist {
			if dispatchesEvent, exists := dispatchesEvents.DispatchesEvents()[event]; exists {
				return dispatchesEvent(instance)
			}

			return nil
		}

		if observer := r.getObserver(model); observer != nil {
			if observerEvent := getObserverEvent(event, observer); observerEvent != nil {
				return observerEvent(instance)
			}

			return nil
		}
	}

	if dest != nil {
		if dispatchesEvents, exist := dest.(contractsorm.DispatchesEvents); exist {
			if dispatchesEvent, exists := dispatchesEvents.DispatchesEvents()[event]; exists {
				return dispatchesEvent(instance)
			}
			return nil
		}

		if observer := r.getObserver(dest); observer != nil {
			if observerEvent := getObserverEvent(event, observer); observerEvent != nil {
				return observerEvent(instance)
			}
			return nil
		}
	}

	return nil
}

func (r *Query) deleting(dest any) error {
	if !hasID(dest) {
		return nil
	}

	return r.event(contractsorm.EventDeleting, r.conditions.model, dest)
}

func (r *Query) deleted(dest any) error {
	if !hasID(dest) {
		return nil
	}

	return r.event(contractsorm.EventDeleted, r.conditions.model, dest)
}

func (r *Query) dest(value any) *Query {
	conditions := r.conditions
	conditions.dest = value

	return r.setConditions(conditions)
}

func (r *Query) forceDeleting(dest any) error {
	if !hasID(dest) {
		return nil
	}

	return r.event(contractsorm.EventForceDeleting, r.conditions.model, dest)
}

func (r *Query) forceDeleted(dest any) error {
	if !hasID(dest) {
		return nil
	}

	return r.event(contractsorm.EventForceDeleted, r.conditions.model, dest)
}

func (r *Query) getModelConnection() string {
	var (
		model any
		err   error
	)

	if r.conditions.model != nil {
		model, err = modelToStruct(r.conditions.model)
	} else if r.conditions.dest != nil {
		model, err = modelToStruct(r.conditions.dest)
	}
	if err != nil {
		return ""
	}

	connectionModel, ok := model.(contractsorm.ModelWithConnection)
	if !ok {
		return ""
	}

	return connectionModel.Connection()
}

func (r *Query) getObserver(dest any) contractsorm.Observer {
	destType := reflect.TypeOf(dest)
	if destType.Kind() == reflect.Pointer {
		destType = destType.Elem()
	}

	for _, observer := range r.modelToObserver {
		modelType := reflect.TypeOf(observer.Model)
		if modelType.Kind() == reflect.Pointer {
			modelType = modelType.Elem()
		}
		if destType.PkgPath() == modelType.PkgPath() && destType.Name() == modelType.Name() {
			return observer.Observer
		}
	}

	return nil
}

func (r *Query) new(db *gormio.DB) *Query {
	return NewQuery(r.ctx, r.config, r.dbConfig, db, r.grammar, r.log, r.modelToObserver, &r.conditions)
}

func (r *Query) omitCreate(value any) error {
	if len(r.instance.Statement.Omits) > 1 {
		if slices.Contains(r.instance.Statement.Omits, Associations) {
			return errors.OrmQueryAssociationsConflict
		}
	}

	if len(r.instance.Statement.Omits) == 1 && r.instance.Statement.Omits[0] == Associations {
		r.instance.Statement.Selects = []string{}
	}

	if err := r.saving(value); err != nil {
		return err
	}
	if err := r.creating(value); err != nil {
		return err
	}

	if len(r.instance.Statement.Omits) == 1 && r.instance.Statement.Omits[0] == Associations {
		if err := r.instance.Omit(Associations).Create(value).Error; err != nil {
			return err
		}
	} else {
		if err := r.instance.Create(value).Error; err != nil {
			return err
		}
	}

	if err := r.created(value); err != nil {
		return err
	}
	if err := r.saved(value); err != nil {
		return err
	}

	return nil
}

func (r *Query) omitSave(value any) error {
	if slices.Contains(r.instance.Statement.Omits, Associations) {
		return r.instance.Omit(Associations).Save(value).Error
	}

	return r.instance.Save(value).Error
}

func (r *Query) refreshConnection() (*Query, error) {
	connection := r.getModelConnection()
	if connection == "" || connection == r.dbConfig.Connection {
		return r, nil
	}

	query, ok := r.queries[connection]
	if !ok {
		var err error
		query, _, err = BuildQuery(r.ctx, r.config, connection, r.log, r.modelToObserver)
		if err != nil {
			return nil, err
		}

		if r.queries == nil {
			r.queries = make(map[string]*Query)
		}
		r.queries[connection] = query
	}

	query.conditions = r.conditions

	return query, nil
}

func (r *Query) restored(dest any) error {
	return r.event(contractsorm.EventRestored, r.conditions.model, dest)
}

func (r *Query) restoring(dest any) error {
	return r.event(contractsorm.EventRestoring, r.conditions.model, dest)
}

func (r *Query) retrieved(dest any) error {
	if isSlice(dest) {
		return nil
	}

	return r.event(contractsorm.EventRetrieved, r.conditions.model, dest)
}

func (r *Query) save(value any) error {
	return r.instance.Omit(Associations).Save(value).Error
}

func (r *Query) saved(dest any) error {
	if isSlice(dest) {
		return nil
	}

	return r.event(contractsorm.EventSaved, r.conditions.model, dest)
}

func (r *Query) saving(dest any) error {
	if isSlice(dest) {
		return nil
	}

	return r.event(contractsorm.EventSaving, r.conditions.model, dest)
}

func (r *Query) selectCreate(value any) error {
	if len(r.instance.Statement.Selects) > 1 {
		if slices.Contains(r.instance.Statement.Selects, Associations) {
			return errors.OrmQueryAssociationsConflict
		}
	}

	if len(r.instance.Statement.Selects) == 1 && r.instance.Statement.Selects[0] == Associations {
		r.instance.Statement.Selects = []string{}
	}

	if err := r.saving(value); err != nil {
		return err
	}
	if err := r.creating(value); err != nil {
		return err
	}

	if err := r.instance.Create(value).Error; err != nil {
		return err
	}

	if err := r.created(value); err != nil {
		return err
	}
	if err := r.saved(value); err != nil {
		return err
	}

	return nil
}

func (r *Query) selectSave(value any) error {
	if slices.Contains(r.instance.Statement.Selects, Associations) {
		return r.instance.Session(&gormio.Session{FullSaveAssociations: true}).Save(value).Error
	}

	if err := r.instance.Save(value).Error; err != nil {
		return err
	}

	return nil
}

func (r *Query) setConditions(conditions Conditions) *Query {
	query := r.new(r.instance)
	query.conditions = conditions

	return query
}

func (r *Query) updating(dest any) error {
	if isSlice(dest) {
		return nil
	}

	return r.event(contractsorm.EventUpdating, r.conditions.model, dest)
}

func (r *Query) updated(dest any) error {
	if isSlice(dest) {
		return nil
	}

	return r.event(contractsorm.EventUpdated, r.conditions.model, dest)
}

func (r *Query) update(values any) (*contractsdb.Result, error) {
	if len(r.instance.Statement.Selects) > 0 && len(r.instance.Statement.Omits) > 0 {
		return nil, errors.OrmQuerySelectAndOmitsConflict
	}

	if v, ok := values.(map[string]any); ok {
		var err error
		if values, err = r.grammar.CompileJsonColumnsUpdate(v); err != nil {
			return nil, errors.OrmJsonColumnUpdateInvalid.Args(err)
		}
	}

	if len(r.instance.Statement.Selects) > 0 {
		if slices.Contains(r.instance.Statement.Selects, Associations) {
			result := r.instance.Session(&gormio.Session{FullSaveAssociations: true}).Updates(values)
			return &contractsdb.Result{
				RowsAffected: result.RowsAffected,
			}, result.Error
		}

		result := r.instance.Updates(values)

		return &contractsdb.Result{
			RowsAffected: result.RowsAffected,
		}, result.Error
	}

	if len(r.instance.Statement.Omits) > 0 {
		if slices.Contains(r.instance.Statement.Omits, Associations) {
			result := r.instance.Omit(Associations).Updates(values)

			return &contractsdb.Result{
				RowsAffected: result.RowsAffected,
			}, result.Error
		}
		result := r.instance.Updates(values)

		return &contractsdb.Result{
			RowsAffected: result.RowsAffected,
		}, result.Error
	}
	result := r.instance.Omit(Associations).Updates(values)

	return &contractsdb.Result{
		RowsAffected: result.RowsAffected,
	}, result.Error
}

func buildSelectForCount(query *Query) *Query {
	conditions := query.conditions

	// If selectColumns only contains a raw select with spaces (rename), gorm will fail, but this case will appear when calling Paginate, so use COUNT(*) here.
	// If there are multiple selectColumns, gorm will transform them into *, so no need to handle that case.
	// For example: Select("name as n").Count() will fail, but Select("name", "age as a").Count() will be treated as Select("*").Count()
	if len(conditions.selectColumns) == 1 && str.Of(conditions.selectColumns[0]).Trim().Contains(" ") {
		conditions.selectColumns = []string{str.Of(conditions.selectColumns[0]).Split(" ")[0]}
	}

	return query.setConditions(conditions).addGlobalScopes().buildConditions()
}

func filterFindConditions(conds ...any) error {
	if len(conds) > 0 {
		switch cond := conds[0].(type) {
		case string:
			if cond == "" {
				return errors.OrmMissingWhereClause
			}
		default:
			reflectValue := reflect.Indirect(reflect.ValueOf(cond))
			switch reflectValue.Kind() {
			case reflect.Slice, reflect.Array:
				if reflectValue.Len() == 0 {
					return errors.OrmMissingWhereClause
				}
			}
		}
	}

	return nil
}

func getDeletedAtColumn(model any) string {
	if model == nil {
		return ""
	}

	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	for i := 0; i < t.NumField(); i++ {
		if !t.Field(i).IsExported() {
			continue
		}
		if t.Field(i).Type.Kind() == reflect.Struct {
			if t.Field(i).Type == reflect.TypeOf(gormio.DeletedAt{}) {
				return t.Field(i).Name
			}

			structField := t.Field(i).Type
			for j := 0; j < structField.NumField(); j++ {
				if !structField.Field(j).IsExported() {
					continue
				}

				if structField.Field(j).Type == reflect.TypeOf(gormio.DeletedAt{}) {
					return structField.Field(j).Name
				}
			}
		}
	}

	return ""
}

func getObserverEvent(event contractsorm.EventType, observer contractsorm.Observer) func(contractsorm.Event) error {
	switch event {
	case contractsorm.EventCreated:
		return observer.Created
	case contractsorm.EventCreating:
		if o, ok := observer.(contractsorm.ObserverWithCreating); ok {
			return o.Creating
		}
	case contractsorm.EventDeleted:
		return observer.Deleted
	case contractsorm.EventDeleting:
		if o, ok := observer.(contractsorm.ObserverWithDeleting); ok {
			return o.Deleting
		}
	case contractsorm.EventForceDeleted:
		return observer.ForceDeleted
	case contractsorm.EventForceDeleting:
		if o, ok := observer.(contractsorm.ObserverWithForceDeleting); ok {
			return o.ForceDeleting
		}
	case contractsorm.EventRestored:
		if o, ok := observer.(contractsorm.ObserverWithRestored); ok {
			return o.Restored
		}
	case contractsorm.EventRestoring:
		if o, ok := observer.(contractsorm.ObserverWithRestoring); ok {
			return o.Restoring
		}
	case contractsorm.EventRetrieved:
		if o, ok := observer.(contractsorm.ObserverWithRetrieved); ok {
			return o.Retrieved
		}
	case contractsorm.EventSaved:
		if o, ok := observer.(contractsorm.ObserverWithSaved); ok {
			return o.Saved
		}
	case contractsorm.EventSaving:
		if o, ok := observer.(contractsorm.ObserverWithSaving); ok {
			return o.Saving
		}
	case contractsorm.EventUpdated:
		return observer.Updated
	case contractsorm.EventUpdating:
		if o, ok := observer.(contractsorm.ObserverWithUpdating); ok {
			return o.Updating
		}
	}

	return nil
}

// modelToStruct normalizes a model value for database operations.
// It handles nil pointers by creating new instances, resolves interface types,
// unwraps slice/array/pointer types to their underlying struct type,
// validates that the final type is a struct (rejecting maps and other types),
// and returns a new instance of the struct type ready for database operations.
// Returns an error if the model type is invalid (map, primitive, etc.).
func modelToStruct(model any) (any, error) {
	if model == nil {
		return nil, errors.OrmQueryInvalidModel.Args("nil")
	}

	modelValue := reflect.ValueOf(model)
	if modelValue.Kind() == reflect.Ptr && modelValue.IsNil() {
		// If the model is a pointer and is nil, we will create a new instance of the model
		modelValue = reflect.New(modelValue.Type().Elem())
	}
	modelType := reflect.Indirect(modelValue).Type()

	if modelType.Kind() == reflect.Interface {
		modelType = reflect.Indirect(modelValue).Elem().Type()
	}

	for modelType.Kind() == reflect.Slice || modelType.Kind() == reflect.Array || modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	if modelType.Kind() == reflect.Map {
		return nil, errors.OrmQueryInvalidModel.Args("map")
	}

	if modelType.Kind() != reflect.Struct {
		if modelType.PkgPath() == "" {
			return nil, errors.OrmQueryInvalidModel.Args("")
		}
		return nil, errors.OrmQueryInvalidModel.Args(fmt.Sprintf(": %s.%s", modelType.PkgPath(), modelType.Name()))
	}

	newModel := reflect.New(modelType)

	return newModel.Interface(), nil
}

func isSlice(dest any) bool {
	if dest == nil {
		return false
	}
	destKind := reflect.Indirect(reflect.ValueOf(dest)).Type().Kind()

	return destKind == reflect.Slice || destKind == reflect.Array
}

func hasID(dest any) bool {
	return database.GetID(dest) != nil
}
