package db

import (
	"context"
	databasesql "database/sql"
	"fmt"
	"maps"
	"reflect"
	"sort"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"

	"github.com/goravel/framework/contracts/database/db"
	contractsdriver "github.com/goravel/framework/contracts/database/driver"
	"github.com/goravel/framework/contracts/database/logger"
	"github.com/goravel/framework/database/utils"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/support/collect"
	"github.com/goravel/framework/support/convert"
	"github.com/goravel/framework/support/deep"
	"github.com/goravel/framework/support/str"
)

type Query struct {
	ctx          context.Context
	err          error
	grammar      contractsdriver.Grammar
	logger       logger.Logger
	readBuilder  db.CommonBuilder
	writeBuilder db.CommonBuilder
	txLogs       *[]TxLog
	conditions   contractsdriver.Conditions
}

func NewQuery(ctx context.Context, readBuilder db.CommonBuilder, writeBuilder db.CommonBuilder, grammar contractsdriver.Grammar, logger logger.Logger, table string, txLogs *[]TxLog) *Query {
	return &Query{
		conditions: contractsdriver.Conditions{
			Table: table,
		},
		ctx:          ctx,
		grammar:      grammar,
		logger:       logger,
		readBuilder:  readBuilder,
		txLogs:       txLogs,
		writeBuilder: writeBuilder,
	}
}

func (r *Query) Chunk(size uint64, callback func(rows []db.Row) error) error {
	offset := uint64(0)

	for {
		rows := r.clone().Offset(offset).Limit(size).Cursor()

		var destSlice []db.Row
		for row := range rows {
			if err := row.Err(); err != nil {
				return err
			}
			destSlice = append(destSlice, row)
		}

		if len(destSlice) == 0 {
			break
		}

		if err := callback(destSlice); err != nil {
			return err
		}

		if len(destSlice) < int(size) {
			break
		}

		offset += size
	}

	return nil
}

func (r *Query) Count() (int64, error) {
	if err := buildSelectForCount(r); err != nil {
		return 0, err
	}

	sql, args, err := r.buildSelect()
	if err != nil {
		return 0, err
	}

	var count int64
	now := carbon.Now()
	err = r.readBuilder.GetContext(r.ctx, &count, sql, args...)
	if err != nil {
		r.trace(r.readBuilder, sql, args, now, -1, err)

		return 0, err
	}

	r.trace(r.readBuilder, sql, args, now, -1, nil)

	return count, nil
}

func (r *Query) CrossJoin(query string, args ...any) db.Query {
	q := r.clone()
	q.conditions.CrossJoin = deep.Append(q.conditions.CrossJoin, contractsdriver.Join{
		Query: query,
		Args:  args,
	})

	return q
}

func (r *Query) Cursor() chan db.Row {
	ch := make(chan db.Row)
	go func() {
		var (
			args  []any
			count int64 = -1
			err   error
			now   = carbon.Now()
			rows  *sqlx.Rows
			sql   string
		)

		defer func() {
			if err != nil {
				ch <- NewRow(nil, err)
			}

			if len(sql) > 0 {
				r.trace(r.readBuilder, sql, args, now, count, err)
			}

			close(ch)
		}()

		if sql, args, err = r.buildSelect(); err != nil {
			return
		}

		if rows, err = r.readBuilder.QueryxContext(r.ctx, sql, args...); err != nil {
			return
		}
		defer errors.Ignore(rows.Close)

		for rows.Next() {
			row := make(map[string]any)
			if err = rows.MapScan(row); err != nil {
				return
			}

			ch <- NewRow(row, nil)
			count++
		}

		if err = rows.Err(); err != nil {
			return
		}
	}()

	return ch
}

func (r *Query) Decrement(column string, value ...uint64) error {
	v := uint64(1)
	if len(value) > 0 {
		v = value[0]
	}

	_, err := r.Update(column, sq.Expr(fmt.Sprintf("%s - ?", column), v))
	if err != nil {
		return err
	}

	return nil
}

func (r *Query) Delete() (*db.Result, error) {
	sql, args, err := r.buildDelete()
	if err != nil {
		return nil, err
	}

	now := carbon.Now()
	result, err := r.writeBuilder.ExecContext(r.ctx, sql, args...)
	if err != nil {
		r.trace(r.writeBuilder, sql, args, now, -1, err)
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.trace(r.writeBuilder, sql, args, now, -1, err)
		return nil, err
	}

	r.trace(r.writeBuilder, sql, args, now, rowsAffected, nil)

	return &db.Result{
		RowsAffected: rowsAffected,
	}, nil
}

func (r *Query) Distinct(columns ...string) db.Query {
	var query *Query

	if len(columns) > 0 {
		query = r.Select(columns...).(*Query)
	} else {
		query = r.clone()
	}

	query.conditions.Distinct = convert.Pointer(true)

	return query
}

func (r *Query) DoesntExist() (bool, error) {
	count, err := r.Count()
	if err != nil {
		return false, err
	}

	return count == 0, nil
}

func (r *Query) Each(callback func(row db.Row) error) error {
	return r.Chunk(1, func(rows []db.Row) error {
		for _, row := range rows {
			if err := callback(row); err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *Query) Exists() (bool, error) {
	count, err := r.Count()
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *Query) Find(dest any, conds ...any) error {
	q, err := r.findQuery(conds)
	if err != nil {
		return err
	}

	destValue := reflect.Indirect(reflect.ValueOf(dest))
	if destValue.Kind() == reflect.Slice {
		return q.Get(dest)
	}

	return q.First(dest)
}

func (r *Query) FindOrFail(dest any, conds ...any) error {
	q, err := r.findQuery(conds)
	if err != nil {
		return err
	}

	destValue := reflect.Indirect(reflect.ValueOf(dest))
	if destValue.Kind() == reflect.Slice {
		return q.Get(dest)
	}

	return q.FirstOrFail(dest)
}

func (r *Query) First(dest any) error {
	sql, args, err := r.buildSelect()
	if err != nil {
		return err
	}

	now := carbon.Now()
	err = r.readBuilder.GetContext(r.ctx, dest, sql, args...)
	if err != nil {
		if errors.Is(err, databasesql.ErrNoRows) {
			r.trace(r.readBuilder, sql, args, now, 0, nil)
			return nil
		}

		r.trace(r.readBuilder, sql, args, now, -1, err)

		return err
	}

	r.trace(r.readBuilder, sql, args, now, 1, nil)

	return nil
}

func (r *Query) FirstOr(dest any, callback func() error) error {
	sql, args, err := r.buildSelect()
	if err != nil {
		return err
	}

	now := carbon.Now()
	err = r.readBuilder.GetContext(r.ctx, dest, sql, args...)
	if err != nil {
		if errors.Is(err, databasesql.ErrNoRows) {
			r.trace(r.readBuilder, sql, args, now, 0, nil)

			return callback()
		}

		r.trace(r.readBuilder, sql, args, now, -1, err)

		return err
	}

	r.trace(r.readBuilder, sql, args, now, 1, nil)

	return nil
}

func (r *Query) FirstOrFail(dest any) error {
	sql, args, err := r.buildSelect()
	if err != nil {
		return err
	}

	now := carbon.Now()
	err = r.readBuilder.GetContext(r.ctx, dest, sql, args...)
	if err != nil {
		r.trace(r.readBuilder, sql, args, now, -1, err)

		return err
	}

	r.trace(r.readBuilder, sql, args, now, 1, nil)

	return nil
}

func (r *Query) Get(dest any) error {
	sql, args, err := r.buildSelect()
	if err != nil {
		return err
	}

	now := carbon.Now()
	err = r.readBuilder.SelectContext(r.ctx, dest, sql, args...)
	if err != nil {
		r.trace(r.readBuilder, sql, args, now, -1, err)
		return err
	}

	destValue := reflect.Indirect(reflect.ValueOf(dest))
	rowsAffected := int64(-1)
	if destValue.Kind() == reflect.Slice {
		rowsAffected = int64(destValue.Len())
	}

	r.trace(r.readBuilder, sql, args, now, rowsAffected, nil)

	return nil
}

func (r *Query) GroupBy(column ...string) db.Query {
	if len(column) == 0 {
		return r
	}

	q := r.clone()
	q.conditions.GroupBy = column

	return q
}

func (r *Query) Having(query any, args ...any) db.Query {
	q := r.clone()
	q.conditions.Having = &contractsdriver.Having{
		Query: query,
		Args:  args,
	}

	return q
}

func (r *Query) Join(query string, args ...any) db.Query {
	q := r.clone()
	q.conditions.Join = deep.Append(q.conditions.Join, contractsdriver.Join{
		Query: query,
		Args:  args,
	})

	return q
}

func (r *Query) Increment(column string, value ...uint64) error {
	v := uint64(1)
	if len(value) > 0 {
		v = value[0]
	}

	_, err := r.Update(column, sq.Expr(fmt.Sprintf("%s + ?", column), v))
	if err != nil {
		return err
	}

	return nil
}

func (r *Query) InRandomOrder() db.Query {
	q := r.clone()
	q.conditions.InRandomOrder = convert.Pointer(true)

	return q
}

func (r *Query) Insert(data any) (*db.Result, error) {
	mapData, err := convertToSliceMap(data)
	if err != nil {
		return nil, err
	}
	if len(mapData) == 0 {
		return nil, errors.DatabaseDataIsEmpty
	}

	sql, args, err := r.buildInsert(mapData)
	if err != nil {
		return nil, err
	}

	now := carbon.Now()
	result, err := r.writeBuilder.ExecContext(r.ctx, sql, args...)
	if err != nil {
		r.trace(r.writeBuilder, sql, args, now, -1, err)
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.trace(r.writeBuilder, sql, args, now, -1, err)
		return nil, err
	}

	r.trace(r.writeBuilder, sql, args, now, rowsAffected, nil)

	return &db.Result{
		RowsAffected: rowsAffected,
	}, nil
}

func (r *Query) InsertGetID(data any) (int64, error) {
	mapData, err := convertToMap(data)
	if err != nil {
		return 0, err
	}
	if len(mapData) == 0 {
		return 0, errors.DatabaseUnsupportedType.Args("nil", "struct, map[string]any").SetModule("DB")
	}

	sql, args, err := r.buildInsert([]map[string]any{mapData})
	if err != nil {
		return 0, err
	}

	now := carbon.Now()
	result, err := r.writeBuilder.ExecContext(r.ctx, sql, args...)
	if err != nil {
		r.trace(r.writeBuilder, sql, args, now, -1, err)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		r.trace(r.writeBuilder, sql, args, now, -1, err)
		return 0, err
	}

	r.trace(r.writeBuilder, sql, args, now, 1, nil)

	return id, nil
}

func (r *Query) Latest(column ...string) db.Query {
	col := "created_at"
	if len(column) > 0 {
		col = column[0]
	}

	return r.OrderByDesc(col)
}

func (r *Query) LeftJoin(query string, args ...any) db.Query {
	q := r.clone()
	q.conditions.LeftJoin = deep.Append(q.conditions.LeftJoin, contractsdriver.Join{
		Query: query,
		Args:  args,
	})

	return q
}

func (r *Query) Limit(limit uint64) db.Query {
	q := r.clone()
	q.conditions.Limit = &limit

	return q
}

func (r *Query) LockForUpdate() db.Query {
	q := r.clone()
	q.conditions.LockForUpdate = convert.Pointer(true)

	return q
}

func (r *Query) Offset(offset uint64) db.Query {
	q := r.clone()
	q.conditions.Offset = &offset

	return q
}

func (r *Query) OrderBy(column string, directions ...string) db.Query {
	direction := "ASC"
	if len(directions) > 0 {
		direction = directions[0]
	}

	q := r.clone()
	q.conditions.OrderBy = deep.Append(q.conditions.OrderBy, column+" "+direction)

	return q
}

func (r *Query) OrderByDesc(column string) db.Query {
	q := r.clone()
	q.conditions.OrderBy = deep.Append(q.conditions.OrderBy, column+" DESC")

	return q
}

func (r *Query) OrderByRaw(raw string) db.Query {
	q := r.clone()
	q.conditions.OrderBy = deep.Append(q.conditions.OrderBy, raw)

	return q
}

func (r *Query) OrWhere(query any, args ...any) db.Query {
	return r.addWhere(contractsdriver.Where{
		Query: query,
		Args:  args,
		Or:    true,
	})
}

func (r *Query) OrWhereBetween(column string, x, y any) db.Query {
	return r.OrWhere(sq.Expr(fmt.Sprintf("%s BETWEEN ? AND ?", column), x, y))
}

func (r *Query) OrWhereColumn(column1 string, column2 ...string) db.Query {
	if len(column2) == 0 || len(column2) > 2 {
		r.err = errors.DatabaseInvalidArgumentNumber.Args(len(column2), "1 or 2")
		return r
	}

	if len(column2) == 1 {
		return r.OrWhere(sq.Expr(fmt.Sprintf("%s = %s", column1, column2[0])))
	}

	return r.OrWhere(sq.Expr(fmt.Sprintf("%s %s %s", column1, column2[0], column2[1])))
}

func (r *Query) OrWhereIn(column string, values []any) db.Query {
	return r.OrWhere(column, values)
}

func (r *Query) OrWhereJsonContains(column string, value any) db.Query {
	return r.addWhere(contractsdriver.Where{
		Type:  contractsdriver.WhereTypeJsonContains,
		Query: column,
		Args:  []any{value},
		Or:    true,
	})
}

func (r *Query) OrWhereJsonContainsKey(column string) db.Query {
	return r.addWhere(contractsdriver.Where{
		Type:  contractsdriver.WhereTypeJsonContainsKey,
		Query: column,
		Or:    true,
	})
}

func (r *Query) OrWhereJsonDoesntContain(column string, value any) db.Query {
	return r.addWhere(contractsdriver.Where{
		Type:  contractsdriver.WhereTypeJsonContains,
		Query: column,
		Args:  []any{value},
		IsNot: true,
		Or:    true,
	})
}

func (r *Query) OrWhereJsonDoesntContainKey(column string) db.Query {
	return r.addWhere(contractsdriver.Where{
		Type:  contractsdriver.WhereTypeJsonContainsKey,
		Query: column,
		IsNot: true,
		Or:    true,
	})
}

func (r *Query) OrWhereJsonLength(column string, length int) db.Query {
	return r.addWhere(contractsdriver.Where{
		Type:  contractsdriver.WhereTypeJsonLength,
		Query: column,
		Args:  []any{length},
		Or:    true,
	})
}

func (r *Query) OrWhereLike(column string, value string) db.Query {
	return r.OrWhere(sq.Like{column: value})
}

func (r *Query) OrWhereNot(query any, args ...any) db.Query {
	query, args, err := r.buildWhere(contractsdriver.Where{
		Query: query,
		Args:  args,
	})
	if err != nil {
		r.err = err
		return r
	}

	sqlizer, err := r.toSqlizer(query, args)
	if err != nil {
		r.err = err
		return r
	}

	sql, args, err := sqlizer.ToSql()
	if err != nil {
		r.err = err
		return r
	}

	return r.OrWhere(sq.Expr(fmt.Sprintf("NOT (%s)", sql), args...))
}

func (r *Query) OrWhereNotBetween(column string, x, y any) db.Query {
	return r.OrWhere(sq.Expr(fmt.Sprintf("%s NOT BETWEEN ? AND ?", column), x, y))
}

func (r *Query) OrWhereNotIn(column string, values []any) db.Query {
	return r.OrWhere(sq.NotEq{column: values})
}

func (r *Query) OrWhereNotLike(column string, value string) db.Query {
	return r.OrWhere(sq.NotLike{column: value})
}

func (r *Query) OrWhereNotNull(column string) db.Query {
	return r.OrWhere(sq.NotEq{column: nil})
}

func (r *Query) OrWhereNull(column string) db.Query {
	return r.OrWhere(sq.Eq{column: nil})
}

func (r *Query) OrWhereRaw(raw string, args []any) db.Query {
	return r.OrWhere(sq.Expr(raw, args...))
}

func (r *Query) Paginate(page, limit int, dest any, total *int64) error {
	offset := (page - 1) * limit

	q := r.clone()
	count, err := q.Count()
	if err != nil {
		return err
	}

	*total = count

	return r.Offset(uint64(offset)).Limit(uint64(limit)).Get(dest)
}

func (r *Query) Pluck(column string, dest any) error {
	r.conditions.Selects = []string{column}

	return r.Get(dest)
}

func (r *Query) RightJoin(query string, args ...any) db.Query {
	q := r.clone()
	q.conditions.RightJoin = deep.Append(q.conditions.RightJoin, contractsdriver.Join{
		Query: query,
		Args:  args,
	})

	return q
}

func (r *Query) Select(columns ...string) db.Query {
	q := r.clone()
	q.conditions.Selects = deep.Append(q.conditions.Selects, columns...)
	q.conditions.Selects = collect.Unique(q.conditions.Selects)

	// * may be added along with other columns, remove it.
	if len(q.conditions.Selects) > 1 {
		q.conditions.Selects = collect.Filter(q.conditions.Selects, func(column string, _ int) bool {
			return column != "*"
		})
	}

	return q
}

func (r *Query) SharedLock() db.Query {
	q := r.clone()
	q.conditions.SharedLock = convert.Pointer(true)

	return q
}

func (r *Query) Sum(column string, dest any) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr {
		return errors.DatabaseUnsupportedType.Args(destValue.Kind(), "pointer")
	}

	return r.Select(fmt.Sprintf("SUM(%s)", column)).First(dest)
}

func (r *Query) Avg(column string, dest any) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr {
		return errors.DatabaseUnsupportedType.Args(destValue.Kind(), "pointer")
	}

	return r.Select(fmt.Sprintf("AVG(%s)", column)).First(dest)
}

func (r *Query) Min(column string, dest any) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr {
		return errors.DatabaseUnsupportedType.Args(destValue.Kind(), "pointer")
	}

	return r.Select(fmt.Sprintf("MIN(%s)", column)).First(dest)
}

func (r *Query) Max(column string, dest any) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr {
		return errors.DatabaseUnsupportedType.Args(destValue.Kind(), "pointer")
	}

	return r.Select(fmt.Sprintf("MAX(%s)", column)).First(dest)
}

func (r *Query) ToSql() db.ToSql {
	q := r.clone()
	return NewToSql(q, false)
}

func (r *Query) ToRawSql() db.ToSql {
	q := r.clone()
	return NewToSql(q, true)
}

func (r *Query) Update(column any, value ...any) (*db.Result, error) {
	columnStr, ok := column.(string)
	if ok {
		if len(value) != 1 {
			return nil, errors.DatabaseInvalidArgumentNumber.Args(len(value), "1")
		}

		return r.Update(map[string]any{columnStr: value[0]})
	}

	mapData, err := convertToMap(column)
	if err != nil {
		return nil, err
	}

	sql, args, err := r.buildUpdate(mapData)
	if err != nil {
		return nil, err
	}

	now := carbon.Now()
	result, err := r.writeBuilder.ExecContext(r.ctx, sql, args...)
	if err != nil {
		r.trace(r.writeBuilder, sql, args, now, -1, err)
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.trace(r.writeBuilder, sql, args, now, -1, err)
		return nil, err
	}

	r.trace(r.writeBuilder, sql, args, now, rowsAffected, nil)

	return &db.Result{
		RowsAffected: rowsAffected,
	}, nil
}

func (r *Query) UpdateOrInsert(attributes any, values any) (*db.Result, error) {
	mapAttributes, err := convertToMap(attributes)
	if err != nil {
		return nil, err
	}

	mapValues, err := convertToMap(values)
	if err != nil {
		return nil, err
	}

	exist, err := r.clone().Where(mapAttributes).Exists()
	if err != nil {
		return nil, err
	}

	if exist {
		return r.Where(mapAttributes).Update(values)
	}

	maps.Copy(mapAttributes, mapValues)

	return r.Insert(mapAttributes)
}

func (r *Query) Value(column string, dest any) error {
	return r.Select(column).Limit(1).First(dest)
}

func (r *Query) When(condition bool, callback func(query db.Query) db.Query, falseCallback ...func(query db.Query) db.Query) db.Query {
	if condition {
		return callback(r)
	}

	if len(falseCallback) > 0 {
		return falseCallback[0](r)
	}

	return r
}

func (r *Query) Where(query any, args ...any) db.Query {
	return r.addWhere(contractsdriver.Where{
		Query: query,
		Args:  args,
	})
}

func (r *Query) WhereAll(columns []string, args ...any) db.Query {
	op, value, err := utils.PrepareWhereOperatorAndValue(args...)
	if err != nil {
		r.err = err
		return r
	}

	var conditions []string
	var conditionArgs []any
	for _, column := range columns {
		conditions = append(conditions, fmt.Sprintf("%s %v ?", column, op))
		conditionArgs = append(conditionArgs, value)
	}

	query := strings.Join(conditions, " AND ")
	where := contractsdriver.Where{
		Query: sq.Expr(query, conditionArgs...),
	}
	r.conditions.Where = deep.Append(r.conditions.Where, where)

	return r
}

func (r *Query) WhereAny(columns []string, args ...any) db.Query {
	op, value, err := utils.PrepareWhereOperatorAndValue(args...)
	if err != nil {
		r.err = err
		return r
	}

	var orConditions []sq.Sqlizer
	for _, column := range columns {
		orConditions = append(orConditions, sq.Expr(fmt.Sprintf("%s %v ?", column, op), value))
	}

	where := contractsdriver.Where{
		Query: sq.Or(orConditions),
	}
	r.conditions.Where = deep.Append(r.conditions.Where, where)

	return r
}

func (r *Query) WhereBetween(column string, x, y any) db.Query {
	return r.Where(sq.Expr(fmt.Sprintf("%s BETWEEN ? AND ?", column), x, y))
}

func (r *Query) WhereColumn(column1 string, column2 ...string) db.Query {
	if len(column2) == 0 || len(column2) > 2 {
		r.err = errors.DatabaseInvalidArgumentNumber.Args(len(column2), "1 or 2")
		return r
	}

	if len(column2) == 1 {
		return r.Where(sq.Expr(fmt.Sprintf("%s = %s", column1, column2[0])))
	}

	return r.Where(sq.Expr(fmt.Sprintf("%s %s %s", column1, column2[0], column2[1])))
}

func (r *Query) WhereExists(query func() db.Query) db.Query {
	subQuery := query()
	sql, args, err := subQuery.(*Query).buildSelect()
	if err != nil {
		r.err = err
		return r
	}

	sql = r.readBuilder.Explain(sql, args...)

	return r.Where(sq.Expr(fmt.Sprintf("EXISTS (%s)", sql)))
}

func (r *Query) WhereIn(column string, values []any) db.Query {
	return r.Where(column, values)
}

func (r *Query) WhereJsonContains(column string, value any) db.Query {
	return r.addWhere(contractsdriver.Where{
		Type:  contractsdriver.WhereTypeJsonContains,
		Query: column,
		Args:  []any{value},
	})
}

func (r *Query) WhereJsonContainsKey(column string) db.Query {
	return r.addWhere(contractsdriver.Where{
		Type:  contractsdriver.WhereTypeJsonContainsKey,
		Query: column,
	})
}

func (r *Query) WhereJsonDoesntContain(column string, value any) db.Query {
	return r.addWhere(contractsdriver.Where{
		Type:  contractsdriver.WhereTypeJsonContains,
		Query: column,
		Args:  []any{value},
		IsNot: true,
	})
}

func (r *Query) WhereJsonDoesntContainKey(column string) db.Query {
	return r.addWhere(contractsdriver.Where{
		Type:  contractsdriver.WhereTypeJsonContainsKey,
		Query: column,
		IsNot: true,
	})
}

func (r *Query) WhereJsonLength(column string, length int) db.Query {
	return r.addWhere(contractsdriver.Where{
		Type:  contractsdriver.WhereTypeJsonLength,
		Query: column,
		Args:  []any{length},
	})
}

func (r *Query) WhereLike(column string, value string) db.Query {
	return r.Where(sq.Like{column: value})
}

func (r *Query) WhereNone(columns []string, args ...any) db.Query {
	op, value, err := utils.PrepareWhereOperatorAndValue(args...)
	if err != nil {
		r.err = err
		return r
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
	where := contractsdriver.Where{
		Query: sq.Expr(query, conditionArgs...),
	}
	r.conditions.Where = deep.Append(r.conditions.Where, where)

	return r
}

func (r *Query) WhereNot(query any, args ...any) db.Query {
	query, args, err := r.buildWhere(contractsdriver.Where{
		Query: query,
		Args:  args,
	})
	if err != nil {
		r.err = err
		return r
	}

	sqlizer, err := r.toSqlizer(query, args)
	if err != nil {
		r.err = err
		return r
	}

	sql, args, err := sqlizer.ToSql()
	if err != nil {
		r.err = err
		return r
	}

	return r.Where(sq.Expr(fmt.Sprintf("NOT (%s)", sql), args...))
}

func (r *Query) WhereNotBetween(column string, x, y any) db.Query {
	return r.Where(sq.Expr(fmt.Sprintf("%s NOT BETWEEN ? AND ?", column), x, y))
}

func (r *Query) WhereNotIn(column string, values []any) db.Query {
	return r.Where(sq.NotEq{column: values})
}

func (r *Query) WhereNotLike(column string, value string) db.Query {
	return r.Where(sq.NotLike{column: value})
}

func (r *Query) WhereNotNull(column string) db.Query {
	return r.Where(sq.NotEq{column: nil})
}

func (r *Query) WhereNull(column string) db.Query {
	return r.Where(sq.Eq{column: nil})
}

func (r *Query) WhereRaw(raw string, args []any) db.Query {
	return r.Where(sq.Expr(raw, args...))
}

func (r *Query) addWhere(where contractsdriver.Where) db.Query {
	q := r.clone()
	q.conditions.Where = deep.Append(q.conditions.Where, where)

	return q
}

func (r *Query) buildDelete() (sql string, args []any, err error) {
	if r.err != nil {
		return "", nil, r.err
	}

	if r.conditions.Table == "" {
		return "", nil, errors.DatabaseTableIsRequired
	}

	builder := sq.Delete(r.conditions.Table)
	if placeholderFormat := r.grammar.CompilePlaceholderFormat(); placeholderFormat != nil {
		builder = builder.PlaceholderFormat(placeholderFormat)
	}

	sqlizer, err := r.buildWheres(r.conditions.Where)
	if err != nil {
		return "", nil, err
	}

	return builder.Where(sqlizer).ToSql()
}

func (r *Query) buildInsert(data []map[string]any) (sql string, args []any, err error) {
	if r.err != nil {
		return "", nil, r.err
	}

	if r.conditions.Table == "" {
		return "", nil, errors.DatabaseTableIsRequired
	}

	builder := sq.Insert(r.conditions.Table)
	if placeholderFormat := r.grammar.CompilePlaceholderFormat(); placeholderFormat != nil {
		builder = builder.PlaceholderFormat(placeholderFormat)
	}

	// Collect all unique columns from all maps to avoid missing columns
	colSet := make(map[string]bool)
	for _, row := range data {
		for col := range row {
			colSet[col] = true
		}
	}

	cols := make([]string, 0, len(colSet))
	for col := range colSet {
		cols = append(cols, col)
	}
	sort.Strings(cols)
	builder = builder.Columns(cols...)

	for _, row := range data {
		vals := make([]any, 0, len(cols))
		for _, col := range cols {
			vals = append(vals, row[col])
		}
		builder = builder.Values(vals...)
	}

	return builder.ToSql()
}

func (r *Query) buildSelect() (sql string, args []any, err error) {
	if r.err != nil {
		return "", nil, r.err
	}

	if r.conditions.Table == "" {
		return "", nil, errors.DatabaseTableIsRequired
	}

	selects := "*"
	if len(r.conditions.Selects) > 0 {
		selects = strings.Join(r.conditions.Selects, ", ")
	}

	builder := sq.Select(selects)

	if r.conditions.Distinct != nil && *r.conditions.Distinct {
		builder = builder.Distinct()
	}

	if placeholderFormat := r.grammar.CompilePlaceholderFormat(); placeholderFormat != nil {
		builder = builder.PlaceholderFormat(placeholderFormat)
	}

	builder = builder.From(r.conditions.Table)

	for _, join := range r.conditions.Join {
		builder = builder.Join(join.Query, join.Args...)
	}

	for _, leftJoin := range r.conditions.LeftJoin {
		builder = builder.LeftJoin(leftJoin.Query, leftJoin.Args...)
	}

	for _, rightJoin := range r.conditions.RightJoin {
		builder = builder.RightJoin(rightJoin.Query, rightJoin.Args...)
	}

	for _, crossJoin := range r.conditions.CrossJoin {
		builder = builder.CrossJoin(crossJoin.Query, crossJoin.Args...)
	}

	sqlizer, err := r.buildWheres(r.conditions.Where)
	if err != nil {
		return "", nil, err
	}

	builder = builder.Where(sqlizer)

	if r.conditions.InRandomOrder != nil && *r.conditions.InRandomOrder {
		builder = r.grammar.CompileInRandomOrder(builder, &r.conditions)
	}

	if len(r.conditions.GroupBy) > 0 {
		builder = builder.GroupBy(r.conditions.GroupBy...)

		if r.conditions.Having != nil {
			builder = builder.Having(r.conditions.Having.Query, r.conditions.Having.Args...)
		}
	}

	compileOrderByGrammar, ok := r.grammar.(contractsdriver.CompileOrderByGrammar)
	if ok {
		builder = compileOrderByGrammar.CompileOrderBy(builder, &r.conditions)
	} else {
		if len(r.conditions.OrderBy) > 0 {
			builder = builder.OrderBy(r.conditions.OrderBy...)
		}
	}

	compileOffsetGrammar, ok := r.grammar.(contractsdriver.CompileOffsetGrammar)
	if ok {
		builder = compileOffsetGrammar.CompileOffset(builder, &r.conditions)
	} else {
		if r.conditions.Offset != nil {
			builder = builder.Offset(*r.conditions.Offset)
		}
	}

	compileLimitGrammar, ok := r.grammar.(contractsdriver.CompileLimitGrammar)
	if ok {
		builder = compileLimitGrammar.CompileLimit(builder, &r.conditions)
	} else {
		if r.conditions.Limit != nil {
			builder = builder.Limit(*r.conditions.Limit)
		}
	}

	if r.conditions.LockForUpdate != nil && *r.conditions.LockForUpdate {
		builder = r.grammar.CompileLockForUpdate(builder, &r.conditions)
	}
	if r.conditions.SharedLock != nil && *r.conditions.SharedLock {
		builder = r.grammar.CompileSharedLock(builder, &r.conditions)
	}

	return builder.ToSql()
}

func (r *Query) buildUpdate(data map[string]any) (sql string, args []any, err error) {
	if r.err != nil {
		return "", nil, r.err
	}

	if r.conditions.Table == "" {
		return "", nil, errors.DatabaseTableIsRequired
	}

	builder := sq.Update(r.conditions.Table)
	if placeholderFormat := r.grammar.CompilePlaceholderFormat(); placeholderFormat != nil {
		builder = builder.PlaceholderFormat(placeholderFormat)
	}

	sqlizer, err := r.buildWheres(r.conditions.Where)
	if err != nil {
		return "", nil, err
	}

	if data, err = r.grammar.CompileJsonColumnsUpdate(data); err != nil {
		return "", nil, errors.OrmJsonColumnUpdateInvalid.Args(err)
	}

	return builder.Where(sqlizer).SetMap(data).ToSql()
}

func (r *Query) buildWhere(where contractsdriver.Where) (any, []any, error) {
	switch query := where.Query.(type) {
	case string:
		switch where.Type {
		case contractsdriver.WhereTypeJsonContains:
			var err error
			query, where.Args, err = r.grammar.CompileJsonContains(query, where.Args[0], where.IsNot)
			if err != nil {
				return nil, nil, errors.OrmJsonContainsInvalidBinding.Args(err)
			}
		case contractsdriver.WhereTypeJsonContainsKey:
			query = str.Of(r.grammar.CompileJsonContainsKey(query, where.IsNot)).Replace("?", "??").String()
		case contractsdriver.WhereTypeJsonLength:
			segments := strings.SplitN(query, " ", 2)
			segments[0] = r.grammar.CompileJsonLength(segments[0])
			query = strings.Join(segments, " ")
		default:
			if str.Of(query).Trim().Contains("->") {
				segments := strings.Split(query, " ")
				for i := range segments {
					if strings.Contains(segments[i], "->") {
						segments[i] = r.grammar.CompileJsonSelector(segments[i])
					}
				}
				query = strings.Join(segments, " ")
				where.Args = r.grammar.CompileJsonValues(where.Args...)
			}
		}
		if !str.Of(query).Trim().Contains("?") {
			if len(where.Args) > 1 {
				return sq.Eq{query: where.Args}, nil, nil
			} else if len(where.Args) == 1 {
				return sq.Eq{query: where.Args[0]}, nil, nil
			}
		}

		return query, where.Args, nil
	case map[string]any:
		return sq.Eq(query), nil, nil
	case func(db.Query) db.Query:
		// Handle nested conditions by creating a new query and applying the callback
		nestedQuery := NewQuery(r.ctx, r.readBuilder, r.writeBuilder, r.grammar, r.logger, r.conditions.Table, r.txLogs)
		nestedQuery = query(nestedQuery).(*Query)

		// Build the nested conditions
		sqlizer, err := r.buildWheres(nestedQuery.conditions.Where)
		if err != nil {
			return nil, nil, err
		}

		return sqlizer, nil, nil
	default:
		return where.Query, where.Args, nil
	}
}

func (r *Query) buildWheres(wheres []contractsdriver.Where) (sq.Sqlizer, error) {
	if len(wheres) == 0 {
		return nil, nil
	}

	var sqlizers []sq.Sqlizer
	for _, where := range wheres {
		query, args, err := r.buildWhere(where)
		if err != nil {
			return nil, err
		}

		sqlizer, err := r.toSqlizer(query, args)
		if err != nil {
			return nil, err
		}

		if where.Or && len(sqlizers) > 0 {
			// If it's an OR condition and we have previous conditions,
			// wrap the previous conditions in an AND and create an OR condition
			if len(sqlizers) == 1 {
				sqlizers = []sq.Sqlizer{
					sq.Or{
						sqlizers[0],
						sqlizer,
					},
				}
			} else {
				sqlizers = []sq.Sqlizer{
					sq.Or{
						sq.And(sqlizers),
						sqlizer,
					},
				}
			}
		} else {
			// For regular WHERE conditions or the first condition
			sqlizers = append(sqlizers, sqlizer)
		}
	}

	if len(sqlizers) == 1 {
		return sqlizers[0], nil
	}

	return sq.And(sqlizers), nil
}

func (r *Query) clone() *Query {
	query := NewQuery(r.ctx, r.readBuilder, r.writeBuilder, r.grammar, r.logger, r.conditions.Table, r.txLogs)
	query.conditions = r.conditions
	query.err = r.err

	return query
}

func (r *Query) findQuery(conds []any) (db.Query, error) {
	var q db.Query
	if len(conds) > 2 {
		return nil, errors.DatabaseInvalidArgumentNumber.Args(len(conds), "1 or 2")
	} else if len(conds) == 1 {
		q = r.Where("id", conds...)
	} else if len(conds) == 2 {
		q = r.Where(conds[0], conds[1])
	} else {
		q = r.clone()
	}

	return q, nil
}

func (r *Query) toSqlizer(query any, args []any) (sq.Sqlizer, error) {
	switch q := query.(type) {
	case map[string]any:
		return sq.Eq(q), nil
	case string:
		return sq.Expr(q, args...), nil
	case sq.Sqlizer:
		return q, nil
	default:
		return nil, errors.DatabaseUnsupportedType.Args(fmt.Sprintf("%T", query), "string-keyed map or string or squirrel.Sqlizer")
	}
}

func (r *Query) trace(builder db.CommonBuilder, sql string, args []any, now *carbon.Carbon, rowsAffected int64, err error) {
	if r.txLogs != nil {
		*r.txLogs = append(*r.txLogs, TxLog{
			ctx:          r.ctx,
			begin:        now,
			sql:          builder.Explain(sql, args...),
			rowsAffected: rowsAffected,
			err:          err,
		})
	} else {
		r.logger.Trace(r.ctx, now, builder.Explain(sql, args...), rowsAffected, err)
	}
}

func buildSelectForCount(query *Query) error {
	distinct := query.conditions.Distinct != nil && *query.conditions.Distinct

	// If selectColumns only contains a raw select with spaces (rename), gorm will fail, but this case will appear when calling Paginate, so use COUNT(*) here.
	// If there are multiple selectColumns, gorm will transform them into *, so no need to handle that case.
	// For example: Select("name as n").Count() will fail, but Select("name", "age as a").Count() will be treated as Select("*").Count()
	if len(query.conditions.Selects) > 1 {
		query.conditions.Selects = []string{"COUNT(*)"}
	} else if len(query.conditions.Selects) == 1 {
		column := query.conditions.Selects[0]
		if str.Of(query.conditions.Selects[0]).Trim().Contains(" ") {
			column = str.Of(query.conditions.Selects[0]).Split(" ")[0]
		}

		if distinct {
			query.conditions.Selects = []string{fmt.Sprintf("COUNT(DISTINCT %s)", column)}
		} else {
			query.conditions.Selects = []string{fmt.Sprintf("COUNT(%s)", column)}
		}
	} else {
		if distinct {
			return errors.DatabaseCountDistinctWithoutColumns
		} else {
			query.conditions.Selects = []string{"COUNT(*)"}
		}
	}

	query.conditions.Distinct = nil

	return nil
}
