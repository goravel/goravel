package dbresolver

import (
	"context"
	"time"

	"gorm.io/gorm"
)

func (dr *DBResolver) SetConnMaxIdleTime(d time.Duration) *DBResolver {
	dr.Call(func(connPool gorm.ConnPool) error {
		if conn, ok := connPool.(interface{ SetConnMaxIdleTime(time.Duration) }); ok {
			conn.SetConnMaxIdleTime(d)
		} else {
			dr.DB.Logger.Error(context.Background(), "SetConnMaxIdleTime not implemented for %#v, please use golang v1.15+", conn)
		}
		return nil
	})

	return dr
}

func (dr *DBResolver) SetConnMaxLifetime(d time.Duration) *DBResolver {
	dr.Call(func(connPool gorm.ConnPool) error {
		if conn, ok := connPool.(interface{ SetConnMaxLifetime(time.Duration) }); ok {
			conn.SetConnMaxLifetime(d)
		} else {
			dr.DB.Logger.Error(context.Background(), "SetConnMaxLifetime not implemented for %#v", conn)
		}
		return nil
	})

	return dr
}

func (dr *DBResolver) SetMaxIdleConns(n int) *DBResolver {
	dr.Call(func(connPool gorm.ConnPool) error {
		if conn, ok := connPool.(interface{ SetMaxIdleConns(int) }); ok {
			conn.SetMaxIdleConns(n)
		} else {
			dr.DB.Logger.Error(context.Background(), "SetMaxIdleConns not implemented for %#v", conn)
		}
		return nil
	})

	return dr
}

func (dr *DBResolver) SetMaxOpenConns(n int) *DBResolver {
	dr.Call(func(connPool gorm.ConnPool) error {

		if conn, ok := connPool.(interface{ SetMaxOpenConns(int) }); ok {
			conn.SetMaxOpenConns(n)
		} else {
			dr.DB.Logger.Error(context.Background(), "SetMaxOpenConns not implemented for %#v", conn)
		}
		return nil
	})

	return dr
}

func (dr *DBResolver) Call(fc func(connPool gorm.ConnPool) error) error {
	if dr.DB != nil {
		for _, r := range dr.resolvers {
			if err := r.call(fc); err != nil {
				return err
			}
		}

		if dr.global != nil {
			if err := dr.global.call(fc); err != nil {
				return err
			}
		}
	} else {
		dr.compileCallbacks = append(dr.compileCallbacks, fc)
	}

	return nil
}
