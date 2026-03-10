package driver

import (
	"sync"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/database"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/support/color"
)

var (
	connectionToDB     sync.Map
	connectionToDBLock sync.Mutex
	pingWarning        sync.Once
)

func BuildGorm(config config.Config, logger logger.Interface, pool database.Pool, connection string) (*gorm.DB, error) {
	if db, ok := connectionToDB.Load(connection); ok {
		return db.(*gorm.DB), nil
	}

	if len(pool.Writers) == 0 {
		return nil, errors.DatabaseConfigNotFound
	}

	// If the database is empty, it means the database is not configured, we don't want to return an error or print a warning here.
	if pool.Writers[0].Database == "" {
		return nil, nil
	}

	connectionToDBLock.Lock()
	defer connectionToDBLock.Unlock()

	if db, ok := connectionToDB.Load(connection); ok {
		return db.(*gorm.DB), nil
	}

	gormConfig := &gorm.Config{
		DisableAutomaticPing:                     true,
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   true,
		Logger:                                   logger,
		NowFunc: func() time.Time {
			return carbon.Now().StdTime()
		},
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   pool.Writers[0].Prefix,
			SingularTable: pool.Writers[0].Singular,
			NoLowerCase:   pool.Writers[0].NoLowerCase,
			NameReplacer:  pool.Writers[0].NameReplacer,
		},
	}

	instance, err := gorm.Open(pool.Writers[0].Dialector, gormConfig)
	if err != nil {
		return nil, err
	}
	if pinger, ok := instance.ConnPool.(interface{ Ping() error }); ok {
		if err = pinger.Ping(); err != nil {
			pingWarning.Do(func() {
				color.Warningln(err.Error())
			})
		}
	}

	maxIdleConns := config.GetInt("database.pool.max_idle_conns", 10)
	maxOpenConns := config.GetInt("database.pool.max_open_conns", 100)
	connMaxIdleTime := config.GetDuration("database.pool.conn_max_idletime", 3600)
	connMaxLifetime := config.GetDuration("database.pool.conn_max_lifetime", 3600)

	if len(pool.Writers) == 1 && len(pool.Readers) == 0 {
		db, err := instance.DB()
		if err != nil {
			return nil, err
		}

		db.SetMaxIdleConns(maxIdleConns)
		db.SetMaxOpenConns(maxOpenConns)
		db.SetConnMaxIdleTime(connMaxIdleTime * time.Second)
		db.SetConnMaxLifetime(connMaxLifetime * time.Second)

		connectionToDB.Store(connection, instance)

		return instance, nil
	}

	var (
		writeDialectors []gorm.Dialector
		readDialectors  []gorm.Dialector
	)

	for _, writer := range pool.Writers {
		writeDialectors = append(writeDialectors, writer.Dialector)
	}

	for _, reader := range pool.Readers {
		readDialectors = append(readDialectors, reader.Dialector)
	}

	if err := instance.Use(dbresolver.Register(dbresolver.Config{
		Sources:           writeDialectors,
		Replicas:          readDialectors,
		Policy:            dbresolver.RandomPolicy{},
		TraceResolverMode: true,
	}).SetMaxIdleConns(maxIdleConns).
		SetMaxOpenConns(maxOpenConns).
		SetConnMaxLifetime(connMaxLifetime * time.Second).
		SetConnMaxIdleTime(connMaxIdleTime * time.Second)); err != nil {
		return nil, err
	}

	connectionToDB.Store(connection, instance)

	return instance, nil
}

func ResetConnections() {
	connectionToDBLock.Lock()
	defer connectionToDBLock.Unlock()

	connectionToDB.Range(func(key, value any) bool {
		connectionToDB.Delete(key)
		return true
	})
}
