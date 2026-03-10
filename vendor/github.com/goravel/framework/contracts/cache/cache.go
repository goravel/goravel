package cache

import (
	"context"
	"time"

	"github.com/goravel/framework/contracts/testing/docker"
)

type Cache interface {
	Driver
	Store(name string) Driver
}

type Driver interface {
	// Add an item in the cache if the key does not exist.
	Add(key string, value any, t time.Duration) bool
	// Decrement decrements the value of an item in the cache.
	Decrement(key string, value ...int64) (int64, error)
	// Docker gets the docker driver.
	Docker() (docker.CacheDriver, error)
	// Forever add an item in the cache indefinitely.
	Forever(key string, value any) bool
	// Forget removes an item from the cache.
	Forget(key string) bool
	// Flush remove all items from the cache.
	Flush() bool
	// Get retrieve an item from the cache by key.
	Get(key string, def ...any) any
	// GetBool retrieves an item from the cache by key as a boolean.
	GetBool(key string, def ...bool) bool
	// GetInt retrieves an item from the cache by key as an integer.
	GetInt(key string, def ...int) int
	// GetInt64 retrieves an item from the cache by key as a 64-bit integer.
	GetInt64(key string, def ...int64) int64
	// GetString retrieves an item from the cache by key as a string.
	GetString(key string, def ...string) string
	// Has check an item exists in the cache.
	Has(key string) bool
	// Increment increments the value of an item in the cache.
	Increment(key string, value ...int64) (int64, error)
	// Lock get a lock instance, the lock will not be expired if the second parameter is not set.
	Lock(key string, t ...time.Duration) Lock
	// Put Driver an item in the cache for a given time.
	Put(key string, value any, t time.Duration) error
	// Pull retrieve an item from the cache and delete it.
	Pull(key string, def ...any) any
	// Remember gets an item from the cache, or execute the given Closure and store the result.
	Remember(key string, ttl time.Duration, callback func() (any, error)) (any, error)
	// RememberForever get an item from the cache, or execute the given Closure and store the result forever.
	RememberForever(key string, callback func() (any, error)) (any, error)
	// WithContext returns a new Cache instance with the given context.
	WithContext(ctx context.Context) Driver
}

type Lock interface {
	// Block attempt to acquire the lock for the given number of seconds.
	Block(t time.Duration, callback ...func()) bool
	// BlockWithTicker attempt to acquire the lock for the given number of seconds,
	// the second parameter is a ticker that will be used to check if the lock is available.
	BlockWithTicker(t time.Duration, ticker time.Duration, callback ...func()) bool
	// Get attempts to acquire the lock.
	Get(callback ...func()) bool
	// Release the lock.
	Release() bool
	// ForceRelease releases the lock in disregard of ownership.
	ForceRelease() bool
}
