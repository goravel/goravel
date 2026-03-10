package db

import (
	"context"
	databasesql "database/sql"
	"fmt"
	"reflect"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"github.com/goravel/framework/contracts/config"
	contractsdb "github.com/goravel/framework/contracts/database/db"
	contractsdriver "github.com/goravel/framework/contracts/database/driver"
	contractslogger "github.com/goravel/framework/contracts/database/logger"
	"github.com/goravel/framework/contracts/log"
	databasedriver "github.com/goravel/framework/database/driver"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/carbon"
)

type DB struct {
	contractsdb.Tx
	config  config.Config
	ctx     context.Context
	driver  contractsdriver.Driver
	gorm    *gorm.DB
	logger  contractslogger.Logger
	queries map[string]contractsdb.DB
}

func NewDB(ctx context.Context, config config.Config, driver contractsdriver.Driver, logger contractslogger.Logger, gormDB *gorm.DB) (*DB, error) {
	return &DB{
		Tx:      NewTx(ctx, driver, logger, gormDB, nil, &[]TxLog{}),
		ctx:     ctx,
		config:  config,
		driver:  driver,
		gorm:    gormDB,
		logger:  logger,
		queries: make(map[string]contractsdb.DB),
	}, nil
}

func BuildDB(ctx context.Context, config config.Config, log log.Log, connection string) (*DB, error) {
	driverCallback, exist := config.Get(fmt.Sprintf("database.connections.%s.via", connection)).(func() (contractsdriver.Driver, error))
	if !exist {
		return nil, errors.DatabaseConfigNotFound
	}

	driver, err := driverCallback()
	if err != nil {
		return nil, err
	}

	pool := driver.Pool()
	logger := NewLogger(config, log)
	gorm, err := databasedriver.BuildGorm(config, logger.ToGorm(), pool, connection)
	if err != nil {
		return nil, err
	}

	return NewDB(ctx, config, driver, logger, gorm)
}

func (r *DB) BeginTransaction() (contractsdb.Tx, error) {
	driverName := r.driver.Pool().Writers[0].Driver
	txBuilder, err := NewTxBuilder(r.gorm.Clauses(dbresolver.Write), driverName)
	if err != nil {
		return nil, err
	}

	return NewTx(r.ctx, r.driver, r.logger, nil, txBuilder, &[]TxLog{}), nil
}

func (r *DB) Connection(name string) contractsdb.DB {
	if name == "" {
		name = r.config.GetString("database.default")
	}

	if _, ok := r.queries[name]; !ok {
		db, err := BuildDB(r.ctx, r.config, r.logger.Log(), name)
		if err != nil {
			r.logger.Panicf(r.ctx, err.Error())
			return nil
		}
		r.queries[name] = db
		db.queries = r.queries
	}

	return r.queries[name]
}

func (r *DB) Transaction(callback func(tx contractsdb.Tx) error) (err error) {
	tx, err := r.BeginTransaction()
	if err != nil {
		return err
	}

	defer func() {
		if re := recover(); re != nil {
			err = fmt.Errorf("panic: %v", re)
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = errors.Join(err, rollbackErr)
			}
		}
	}()

	if err = callback(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	return tx.Commit()
}

func (r *DB) WithContext(ctx context.Context) contractsdb.DB {
	db, err := NewDB(ctx, r.config, r.driver, r.logger, r.gorm)
	if err != nil {
		r.logger.Panicf(r.ctx, err.Error())
		return nil
	}
	return db
}

type Tx struct {
	ctx        context.Context
	grammar    contractsdriver.Grammar
	logger     contractslogger.Logger
	txBuilder  contractsdb.TxBuilder
	gormDB     *gorm.DB
	txLogs     *[]TxLog
	driverName string
}

func NewTx(
	ctx context.Context,
	driver contractsdriver.Driver,
	logger contractslogger.Logger,
	gormDB *gorm.DB,
	txBuilder contractsdb.TxBuilder,
	txLogs *[]TxLog,
) *Tx {
	pool := driver.Pool()
	driverName := pool.Writers[0].Driver

	return &Tx{
		ctx:        ctx,
		driverName: driverName,
		gormDB:     gormDB,
		grammar:    driver.Grammar(),
		logger:     logger,
		txBuilder:  txBuilder,
		txLogs:     txLogs,
	}
}

func (r *Tx) Commit() error {
	if r.txBuilder == nil {
		return errors.DatabaseTransactionNotStarted
	}

	if err := r.txBuilder.Commit(); err != nil {
		return err
	}

	for _, log := range *r.txLogs {
		r.logger.Trace(log.ctx, log.begin, log.sql, log.rowsAffected, log.err)
	}

	return nil
}

func (r *Tx) Delete(sql string, args ...any) (*contractsdb.Result, error) {
	return r.exec(sql, args...)
}

func (r *Tx) Insert(sql string, args ...any) (*contractsdb.Result, error) {
	return r.exec(sql, args...)
}

func (r *Tx) Rollback() error {
	if r.txBuilder == nil {
		return errors.DatabaseTransactionNotStarted
	}

	return r.txBuilder.Rollback()
}

func (r *Tx) Select(dest any, sql string, args ...any) error {
	var (
		builder contractsdb.CommonBuilder
		realSql string
		err     error
	)

	if r.txBuilder != nil {
		builder = r.txBuilder
	} else {
		builder, err = r.readBuilder()
		if err != nil {
			return err
		}
	}

	realSql = builder.Explain(sql, args...)

	destValue := reflect.Indirect(reflect.ValueOf(dest))

	rowsAffected := int64(1)
	if destValue.Kind() == reflect.Slice {
		if err = builder.SelectContext(r.ctx, dest, realSql, args...); err != nil {
			r.logger.Trace(r.ctx, carbon.Now(), realSql, -1, err)

			return err
		}

		rowsAffected = int64(destValue.Len())
	} else {
		if err = builder.GetContext(r.ctx, dest, realSql, args...); err != nil {
			r.logger.Trace(r.ctx, carbon.Now(), realSql, -1, err)

			return err
		}
	}

	r.logger.Trace(r.ctx, carbon.Now(), realSql, rowsAffected, nil)

	return nil
}

func (r *Tx) Statement(sql string, args ...any) error {
	_, err := r.exec(sql, args...)

	return err
}

func (r *Tx) Table(name string) contractsdb.Query {
	if r.txBuilder != nil {
		return NewQuery(r.ctx, r.txBuilder, r.txBuilder, r.grammar, r.logger, name, r.txLogs)
	}

	readBuilder, err := r.readBuilder()
	if err != nil {
		r.logger.Panicf(r.ctx, err.Error())
		return nil
	}

	writeBuilder, err := r.writeBuilder()
	if err != nil {
		r.logger.Panicf(r.ctx, err.Error())
		return nil
	}

	return NewQuery(r.ctx, readBuilder, writeBuilder, r.grammar, r.logger, name, nil)
}

func (r *Tx) Update(sql string, args ...any) (*contractsdb.Result, error) {
	return r.exec(sql, args...)
}

func (r *Tx) exec(sql string, args ...any) (*contractsdb.Result, error) {
	var (
		builder contractsdb.CommonBuilder
		realSql string
		result  databasesql.Result
		err     error
	)

	if r.txBuilder != nil {
		builder = r.txBuilder
	} else {
		builder, err = r.writeBuilder()
		if err != nil {
			return nil, err
		}
	}

	realSql = builder.Explain(sql, args...)
	result, err = builder.ExecContext(r.ctx, sql, args...)
	if err != nil {
		r.logger.Trace(r.ctx, carbon.Now(), realSql, -1, err)
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Trace(r.ctx, carbon.Now(), realSql, -1, err)
		return nil, err
	}

	r.logger.Trace(r.ctx, carbon.Now(), realSql, rowsAffected, nil)

	return &contractsdb.Result{RowsAffected: rowsAffected}, nil
}

func (r *Tx) readBuilder() (contractsdb.Builder, error) {
	builder, err := NewBuilder(r.gormDB.Clauses(dbresolver.Read), r.driverName)
	if err != nil {
		return nil, err
	}

	return builder, nil
}

func (r *Tx) writeBuilder() (contractsdb.Builder, error) {
	builder, err := NewBuilder(r.gormDB.Clauses(dbresolver.Write), r.driverName)
	if err != nil {
		return nil, err
	}

	return builder, nil
}
