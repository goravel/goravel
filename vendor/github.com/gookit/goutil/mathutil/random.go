package mathutil

import (
	"math/rand"
	"time"
)

// RandomInt return a random int at the [min, max)
//
// Usage:
//
//	RandomInt(10, 99)
//	RandomInt(100, 999)
//	RandomInt(1000, 9999)
func RandomInt(min, max int) int {
	rr := rand.New(rand.NewSource(time.Now().UnixNano()))
	return min + rr.Intn(max-min)
}

// RandInt alias of RandomInt()
func RandInt(min, max int) int { return RandomInt(min, max) }

// RandIntWithSeed alias of RandomIntWithSeed()
func RandIntWithSeed(min, max int, seed int64) int {
	return RandomIntWithSeed(min, max, seed)
}

// RandomIntWithSeed return a random int at the [min, max)
//
// Usage:
//
//	seed := time.Now().UnixNano()
//	RandomIntWithSeed(1000, 9999, seed)
func RandomIntWithSeed(min, max int, seed int64) int {
	rr := rand.New(rand.NewSource(seed))
	return min + rr.Intn(max-min)
}
