package dbresolver

import (
	"math/rand"
	"sync/atomic"

	"gorm.io/gorm"
)

type Policy interface {
	Resolve([]gorm.ConnPool) gorm.ConnPool
}

type PolicyFunc func([]gorm.ConnPool) gorm.ConnPool

func (f PolicyFunc) Resolve(connPools []gorm.ConnPool) gorm.ConnPool {
	return f(connPools)
}

type RandomPolicy struct {
}

func (RandomPolicy) Resolve(connPools []gorm.ConnPool) gorm.ConnPool {
	return connPools[rand.Intn(len(connPools))]
}

func RoundRobinPolicy() Policy {
	var i int
	return PolicyFunc(func(connPools []gorm.ConnPool) gorm.ConnPool {
		i = (i + 1) % len(connPools)
		return connPools[i]
	})
}

func StrictRoundRobinPolicy() Policy {
	var i int64
	return PolicyFunc(func(connPools []gorm.ConnPool) gorm.ConnPool {
		return connPools[int(atomic.AddInt64(&i, 1))%len(connPools)]
	})
}
