package cache

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/spf13/cast"

	contractscache "github.com/goravel/framework/contracts/cache"
	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/testing/docker"
	"github.com/goravel/framework/errors"
)

type Memory struct {
	ctx      context.Context
	prefix   string
	instance sync.Map
}

func NewMemory(config config.Config) (*Memory, error) {
	return &Memory{
		prefix: prefix(config),
	}, nil
}

// Add an item in the cache if the key does not exist.
func (r *Memory) Add(key string, value any, t time.Duration) bool {
	if t != NoExpiration {
		time.AfterFunc(t, func() {
			r.Forget(key)
		})
	}

	_, loaded := r.instance.LoadOrStore(r.key(key), value)
	return !loaded
}

// Decrement decrements the value of an item in the cache.
func (r *Memory) Decrement(key string, value ...int64) (int64, error) {
	if len(value) == 0 {
		value = append(value, 1)
	}

	r.Add(key, new(int64), NoExpiration)
	pv := r.Get(key)
	switch nv := pv.(type) {
	case *atomic.Int64:
		return nv.Add(-value[0]), nil
	case *atomic.Int32:
		return int64(nv.Add(int32(-value[0]))), nil
	case *int64:
		return atomic.AddInt64(nv, -value[0]), nil
	case *int32:
		return int64(atomic.AddInt32(nv, int32(-value[0]))), nil
	default:
		return 0, errors.CacheMemoryInvalidIntValueType.Args(key)
	}
}

func (r *Memory) Docker() (docker.CacheDriver, error) {
	return nil, errors.CacheMemoryDriverNotSupportDocker
}

// Forever Put an item in the cache indefinitely.
func (r *Memory) Forever(key string, value any) bool {
	if err := r.Put(key, value, NoExpiration); err != nil {
		return false
	}

	return true
}

// Forget Remove an item from the cache.
func (r *Memory) Forget(key string) bool {
	r.instance.Delete(r.key(key))

	return true
}

// Flush Remove all items from the cache.
func (r *Memory) Flush() bool {
	r.instance = sync.Map{}
	return true
}

// Get Retrieve an item from the cache by key.
func (r *Memory) Get(key string, def ...any) any {
	val, exist := r.instance.Load(r.key(key))
	if exist {
		return val
	}
	if len(def) == 0 {
		return nil
	}

	switch s := def[0].(type) {
	case func() any:
		return s()
	default:
		return s
	}
}

func (r *Memory) GetBool(key string, def ...bool) bool {
	if len(def) == 0 {
		def = append(def, false)
	}
	res := r.Get(key, def[0])

	return cast.ToBool(res)
}

func (r *Memory) GetInt(key string, def ...int) int {
	if len(def) == 0 {
		def = append(def, 0)
	}

	return cast.ToInt(r.Get(key, def[0]))
}

func (r *Memory) GetInt64(key string, def ...int64) int64 {
	if len(def) == 0 {
		def = append(def, 0)
	}

	return cast.ToInt64(r.Get(key, def[0]))
}

func (r *Memory) GetString(key string, def ...string) string {
	if len(def) == 0 {
		def = append(def, "")
	}

	return cast.ToString(r.Get(key, def[0]))
}

// Has Checks an item exists in the cache.
func (r *Memory) Has(key string) bool {
	_, exist := r.instance.Load(r.key(key))
	return exist
}

func (r *Memory) Increment(key string, value ...int64) (int64, error) {
	if len(value) == 0 {
		value = append(value, 1)
	}

	r.Add(key, new(int64), NoExpiration)
	pv := r.Get(key)
	switch nv := pv.(type) {
	case *atomic.Int64:
		return nv.Add(value[0]), nil
	case *atomic.Int32:
		return int64(nv.Add(int32(value[0]))), nil
	case *int64:
		return atomic.AddInt64(nv, value[0]), nil
	case *int32:
		return int64(atomic.AddInt32(nv, int32(value[0]))), nil
	default:
		return 0, errors.CacheMemoryInvalidIntValueType.Args(key)
	}
}

func (r *Memory) Lock(key string, t ...time.Duration) contractscache.Lock {
	return NewLock(r, key, t...)
}

// Pull Retrieve an item from the cache and delete it.
func (r *Memory) Pull(key string, def ...any) any {
	var res any
	if len(def) == 0 {
		res = r.Get(key)
	} else {
		res = r.Get(key, def[0])
	}
	r.Forget(key)

	return res
}

// Put an item in the cache for a given number of seconds.
func (r *Memory) Put(key string, value any, t time.Duration) error {
	if t != NoExpiration {
		time.AfterFunc(t, func() {
			r.Forget(key)
		})
	}

	r.instance.Store(r.key(key), value)
	return nil
}

// Remember Get an item from the cache, or execute the given Closure and store the result.
func (r *Memory) Remember(key string, seconds time.Duration, callback func() (any, error)) (any, error) {
	val := r.Get(key, nil)
	if val != nil {
		return val, nil
	}

	var err error
	val, err = callback()
	if err != nil {
		return nil, err
	}

	if err := r.Put(key, val, seconds); err != nil {
		return nil, err
	}

	return val, nil
}

// RememberForever Get an item from the cache, or execute the given Closure and store the result forever.
func (r *Memory) RememberForever(key string, callback func() (any, error)) (any, error) {
	val := r.Get(key, nil)
	if val != nil {
		return val, nil
	}

	var err error
	val, err = callback()
	if err != nil {
		return nil, err
	}

	if err = r.Put(key, val, NoExpiration); err != nil {
		return nil, err
	}

	return val, nil
}

func (r *Memory) WithContext(ctx context.Context) contractscache.Driver {
	r.ctx = ctx

	return r
}

func (r *Memory) key(key string) string {
	return r.prefix + key
}
