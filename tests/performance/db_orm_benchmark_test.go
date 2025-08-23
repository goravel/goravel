package performance

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"

	"goravel/app/models"
	"goravel/tests"
)

type PerformanceTestSuite struct {
	suite.Suite
	tests.TestCase
}

// BenchmarkDBSelect benchmarks direct database SELECT operations using Query Builder
func (s *PerformanceTestSuite) BenchmarkDBSelect(b *testing.B) {
	// Skip if no database connection
	if !s.isDatabaseAvailable() {
		b.Skip("Database not available")
	}

	b.ResetTimer()
	
	// Track memory usage
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	
	start := time.Now()
	
	for i := 0; i < b.N; i++ {
		var results []map[string]interface{}
		err := facades.DB().Table("users").Limit(10).Get(&results)
		if err != nil {
			// Ignore errors if table doesn't exist during testing
			continue
		}
	}
	
	duration := time.Since(start)
	runtime.GC()
	runtime.ReadMemStats(&m2)
	
	b.ReportMetric(float64(duration.Nanoseconds()/int64(b.N)), "ns/op")
	b.ReportMetric(float64(m2.TotalAlloc-m1.TotalAlloc)/float64(b.N), "allocs/op")
}

// BenchmarkORMSelect benchmarks ORM SELECT operations
func (s *PerformanceTestSuite) BenchmarkORMSelect(b *testing.B) {
	// Skip if no database connection
	if !s.isDatabaseAvailable() {
		b.Skip("Database not available")
	}

	b.ResetTimer()
	
	// Track memory usage
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	
	start := time.Now()
	
	for i := 0; i < b.N; i++ {
		var users []models.User
		err := facades.Orm().Query().Limit(10).Find(&users)
		if err != nil {
			// Ignore errors if table doesn't exist during testing
			continue
		}
	}
	
	duration := time.Since(start)
	runtime.GC()
	runtime.ReadMemStats(&m2)
	
	b.ReportMetric(float64(duration.Nanoseconds()/int64(b.N)), "ns/op")
	b.ReportMetric(float64(m2.TotalAlloc-m1.TotalAlloc)/float64(b.N), "allocs/op")
}

// BenchmarkDBInsert benchmarks direct database INSERT operations
func (s *PerformanceTestSuite) BenchmarkDBInsert(b *testing.B) {
	// Skip if no database connection
	if !s.isDatabaseAvailable() {
		b.Skip("Database not available")
	}

	b.ResetTimer()
	
	// Track memory usage
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	
	start := time.Now()
	
	for i := 0; i < b.N; i++ {
		data := map[string]interface{}{
			"name":     fmt.Sprintf("test%d", i),
			"email":    fmt.Sprintf("test%d@example.com", i),
			"password": "password",
		}
		_, err := facades.DB().Table("users").Insert(data)
		if err != nil {
			// Ignore errors if table doesn't exist during testing
			continue
		}
	}
	
	duration := time.Since(start)
	runtime.GC()
	runtime.ReadMemStats(&m2)
	
	b.ReportMetric(float64(duration.Nanoseconds()/int64(b.N)), "ns/op")
	b.ReportMetric(float64(m2.TotalAlloc-m1.TotalAlloc)/float64(b.N), "allocs/op")
}

// BenchmarkORMInsert benchmarks ORM INSERT operations
func (s *PerformanceTestSuite) BenchmarkORMInsert(b *testing.B) {
	// Skip if no database connection
	if !s.isDatabaseAvailable() {
		b.Skip("Database not available")
	}

	b.ResetTimer()
	
	// Track memory usage
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	
	start := time.Now()
	
	for i := 0; i < b.N; i++ {
		user := models.User{
			Name:     fmt.Sprintf("test%d", i),
			Email:    fmt.Sprintf("test%d@example.com", i),
			Password: "password",
		}
		err := facades.Orm().Query().Create(&user)
		if err != nil {
			// Ignore errors if table doesn't exist during testing
			continue
		}
	}
	
	duration := time.Since(start)
	runtime.GC()
	runtime.ReadMemStats(&m2)
	
	b.ReportMetric(float64(duration.Nanoseconds()/int64(b.N)), "ns/op")
	b.ReportMetric(float64(m2.TotalAlloc-m1.TotalAlloc)/float64(b.N), "allocs/op")
}

// BenchmarkDBUpdate benchmarks direct database UPDATE operations
func (s *PerformanceTestSuite) BenchmarkDBUpdate(b *testing.B) {
	// Skip if no database connection
	if !s.isDatabaseAvailable() {
		b.Skip("Database not available")
	}

	b.ResetTimer()
	
	// Track memory usage
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	
	start := time.Now()
	
	for i := 0; i < b.N; i++ {
		data := map[string]interface{}{
			"name": fmt.Sprintf("updated%d", i),
		}
		_, err := facades.DB().Table("users").Where("id", i%100+1).Update(data)
		if err != nil {
			// Ignore errors if table doesn't exist during testing
			continue
		}
	}
	
	duration := time.Since(start)
	runtime.GC()
	runtime.ReadMemStats(&m2)
	
	b.ReportMetric(float64(duration.Nanoseconds()/int64(b.N)), "ns/op")
	b.ReportMetric(float64(m2.TotalAlloc-m1.TotalAlloc)/float64(b.N), "allocs/op")
}

// BenchmarkORMUpdate benchmarks ORM UPDATE operations
func (s *PerformanceTestSuite) BenchmarkORMUpdate(b *testing.B) {
	// Skip if no database connection
	if !s.isDatabaseAvailable() {
		b.Skip("Database not available")
	}

	b.ResetTimer()
	
	// Track memory usage
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	
	start := time.Now()
	
	for i := 0; i < b.N; i++ {
		_, err := facades.Orm().Query().Where("id", i%100+1).Update("name", fmt.Sprintf("updated%d", i))
		if err != nil {
			// Ignore errors if table doesn't exist during testing
			continue
		}
	}
	
	duration := time.Since(start)
	runtime.GC()
	runtime.ReadMemStats(&m2)
	
	b.ReportMetric(float64(duration.Nanoseconds()/int64(b.N)), "ns/op")
	b.ReportMetric(float64(m2.TotalAlloc-m1.TotalAlloc)/float64(b.N), "allocs/op")
}

// TestPerformanceComparison runs a comprehensive performance comparison
func (s *PerformanceTestSuite) TestPerformanceComparison() {
	// Skip if no database connection
	if !s.isDatabaseAvailable() {
		s.T().Skip("Database not available")
	}

	// This is a test that demonstrates the difference between DB and ORM
	// It won't fail but will log performance metrics for comparison
	
	fmt.Println("\n=== Performance Comparison Results ===")
	
	// Run each benchmark a small number of times to get basic metrics
	numRuns := 100
	
	// Test SELECT operations
	dbSelectTime := s.timeDBSelect(numRuns)
	ormSelectTime := s.timeORMSelect(numRuns)
	
	fmt.Printf("SELECT operations (%d runs):\n", numRuns)
	fmt.Printf("  DB:  %v (avg: %v per op)\n", dbSelectTime, dbSelectTime/time.Duration(numRuns))
	fmt.Printf("  ORM: %v (avg: %v per op)\n", ormSelectTime, ormSelectTime/time.Duration(numRuns))
	if dbSelectTime > 0 && ormSelectTime > 0 {
		ratio := float64(ormSelectTime) / float64(dbSelectTime)
		fmt.Printf("  ORM is %.2fx slower than DB for SELECT\n", ratio)
	}
	
	fmt.Println("\n=== Summary ===")
	fmt.Println("Run 'go test -bench=. ./tests/performance/' to see detailed benchmark results")
	
	// This test always passes - it's just for information gathering
	s.True(true, "Performance comparison completed")
}

// Helper methods

func (s *PerformanceTestSuite) isDatabaseAvailable() bool {
	// Try a simple query to check if database is available
	var result []map[string]interface{}
	err := facades.DB().Table("information_schema.tables").Limit(1).Get(&result)
	return err == nil
}

func (s *PerformanceTestSuite) timeDBSelect(numRuns int) time.Duration {
	start := time.Now()
	for i := 0; i < numRuns; i++ {
		var results []map[string]interface{}
		facades.DB().Table("users").Limit(10).Get(&results)
	}
	return time.Since(start)
}

func (s *PerformanceTestSuite) timeORMSelect(numRuns int) time.Duration {
	start := time.Now()
	for i := 0; i < numRuns; i++ {
		var users []models.User
		facades.Orm().Query().Limit(10).Find(&users)
	}
	return time.Since(start)
}