package testing

import (
	"github.com/stretchr/testify/suite"
	"goravel/bootstrap"
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/assert"
)

type CacheTestSuite struct {
	suite.Suite
}

func TestCacheTestSuite(t *testing.T) {
	bootstrap.Boot()

	suite.Run(t, new(CacheTestSuite))
}

func (s *CacheTestSuite) SetupTest() {

}

func (s *CacheTestSuite) TestPut() {
	t := s.T()
	assert.Nil(t, facades.Cache.Put("name", "Goravel", 1*time.Second))
	assert.True(t, facades.Cache.Has("name"))
	assert.Equal(t, "Goravel", facades.Cache.Get("name", "").(string))
	time.Sleep(2 * time.Second)
	assert.False(t, facades.Cache.Has("name"))
}

func (s *CacheTestSuite) TestGet() {
	t := s.T()
	assert.Nil(t, facades.Cache.Put("name", "Goravel", 1*time.Second))
	assert.Equal(t, "Goravel", facades.Cache.Get("name", "").(string))
	assert.Equal(t, "World", facades.Cache.Get("name1", "World").(string))
	assert.Equal(t, "World1", facades.Cache.Get("name2", func() interface{} {
		return "World1"
	}).(string))
	assert.True(t, facades.Cache.Forget("name"))
	assert.True(t, facades.Cache.Flush())
}

func (s *CacheTestSuite) TestAdd() {
	t := s.T()
	assert.Nil(t, facades.Cache.Put("name", "Goravel", 1*time.Second))
	assert.False(t, facades.Cache.Add("name", "World", 1*time.Second))
	assert.True(t, facades.Cache.Add("name1", "World", 1*time.Second))
	assert.True(t, facades.Cache.Has("name1"))
	time.Sleep(2 * time.Second)
	assert.False(t, facades.Cache.Has("name1"))
	assert.True(t, facades.Cache.Flush())
}

func (s *CacheTestSuite) TestRemember() {
	t := s.T()
	assert.Nil(t, facades.Cache.Put("name", "Goravel", 1*time.Second))
	value, err := facades.Cache.Remember("name", 1*time.Second, func() interface{} {
		return "World"
	})
	assert.Nil(t, err)
	assert.Equal(t, "Goravel", value)

	value, err = facades.Cache.Remember("name1", 1*time.Second, func() interface{} {
		return "World1"
	})
	assert.Nil(t, err)
	assert.Equal(t, "World1", value)
	time.Sleep(2 * time.Second)
	assert.False(t, facades.Cache.Has("name1"))
	assert.True(t, facades.Cache.Flush())
}

func (s *CacheTestSuite) TestRememberForever() {
	t := s.T()
	assert.Nil(t, facades.Cache.Put("name", "Goravel", 1*time.Second))
	value, err := facades.Cache.RememberForever("name", func() interface{} {
		return "World"
	})
	assert.Nil(t, err)
	assert.Equal(t, "Goravel", value)

	value, err = facades.Cache.RememberForever("name1", func() interface{} {
		return "World1"
	})
	assert.Nil(t, err)
	assert.Equal(t, "World1", value)
	assert.True(t, facades.Cache.Flush())
}

func (s *CacheTestSuite) TestPull() {
	t := s.T()
	assert.Nil(t, facades.Cache.Put("name", "Goravel", 1*time.Second))
	assert.True(t, facades.Cache.Has("name"))
	assert.Equal(t, "Goravel", facades.Cache.Pull("name", "").(string))
	assert.False(t, facades.Cache.Has("name"))
}

func (s *CacheTestSuite) TestForever() {
	t := s.T()
	assert.True(t, facades.Cache.Forever("name", "Goravel"))
	assert.Equal(t, "Goravel", facades.Cache.Get("name", "").(string))
	assert.True(t, facades.Cache.Flush())
}

func (s *CacheTestSuite) TestCustomDriver() {
	t := s.T()
	facades.Config.Add("cache", map[string]interface{}{
		"default": "store",
		"stores": map[string]interface{}{
			"store": map[string]interface{}{
				"driver": "custom",
				"via":    &Store{},
			},
		},
		"prefix": "goravel_cache",
	})

	assert.Equal(t, "Goravel", facades.Cache.Get("name", "Goravel").(string))

	facades.Config.Add("cache", map[string]interface{}{
		"default": "redis",
		"stores": map[string]interface{}{
			"redis": map[string]interface{}{
				"driver":     "redis",
				"connection": "default",
			},
		},
		"prefix": "goravel_cache",
	})
}

type Store struct {
}

//Get Retrieve an item from the cache by key.
func (r *Store) Get(key string, defaults interface{}) interface{} {
	return defaults
}

//Has Determine if an item exists in the cache.
func (r *Store) Has(key string) bool {
	return true
}

//Put Store an item in the cache for a given number of seconds.
func (r *Store) Put(key string, value interface{}, seconds time.Duration) error {
	return nil
}

//Pull Retrieve an item from the cache and delete it.
func (r *Store) Pull(key string, defaults interface{}) interface{} {
	return defaults
}

//Add Store an item in the cache if the key does not exist.
func (r *Store) Add(key string, value interface{}, seconds time.Duration) bool {
	return true
}

//Remember Get an item from the cache, or execute the given Closure and store the result.
func (r *Store) Remember(key string, ttl time.Duration, callback func() interface{}) (interface{}, error) {
	return "", nil
}

//RememberForever Get an item from the cache, or execute the given Closure and store the result forever.
func (r *Store) RememberForever(key string, callback func() interface{}) (interface{}, error) {
	return "", nil
}

//Forever Store an item in the cache indefinitely.
func (r *Store) Forever(key string, value interface{}) bool {
	return true
}

//Forget Remove an item from the cache.
func (r *Store) Forget(key string) bool {
	return true
}

//Flush Remove all items from the cache.
func (r *Store) Flush() bool {
	return true
}
