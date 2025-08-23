package performance

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"

	"goravel/app/models"
	"goravel/tests"
)

type AdvancedPerformanceTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestAdvancedPerformanceTestSuite(t *testing.T) {
	suite.Run(t, new(AdvancedPerformanceTestSuite))
}

// BenchmarkDBSelectSingle benchmarks single record selection with DB
func (s *AdvancedPerformanceTestSuite) BenchmarkDBSelectSingle(b *testing.B) {
	if !s.isDatabaseAvailable() {
		b.Skip("Database not available")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result map[string]interface{}
		facades.DB().Table("users").Where("id", 1).First(&result)
	}
}

// BenchmarkORMSelectSingle benchmarks single record selection with ORM  
func (s *AdvancedPerformanceTestSuite) BenchmarkORMSelectSingle(b *testing.B) {
	if !s.isDatabaseAvailable() {
		b.Skip("Database not available")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var user models.User
		facades.Orm().Query().Where("id", 1).First(&user)
	}
}

// BenchmarkDBSelectLarge benchmarks large result set selection with DB
func (s *AdvancedPerformanceTestSuite) BenchmarkDBSelectLarge(b *testing.B) {
	if !s.isDatabaseAvailable() {
		b.Skip("Database not available")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var results []map[string]interface{}
		facades.DB().Table("users").Limit(1000).Get(&results)
	}
}

// BenchmarkORMSelectLarge benchmarks large result set selection with ORM
func (s *AdvancedPerformanceTestSuite) BenchmarkORMSelectLarge(b *testing.B) {
	if !s.isDatabaseAvailable() {
		b.Skip("Database not available")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []models.User
		facades.Orm().Query().Limit(1000).Find(&users)
	}
}

// BenchmarkDBSelectWithWhere benchmarks complex WHERE queries with DB
func (s *AdvancedPerformanceTestSuite) BenchmarkDBSelectWithWhere(b *testing.B) {
	if !s.isDatabaseAvailable() {
		b.Skip("Database not available")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var results []map[string]interface{}
		facades.DB().Table("users").
			Where("email", "LIKE", "%@example.com").
			Where("id", ">", 1).
			Limit(10).
			Get(&results)
	}
}

// BenchmarkORMSelectWithWhere benchmarks complex WHERE queries with ORM
func (s *AdvancedPerformanceTestSuite) BenchmarkORMSelectWithWhere(b *testing.B) {
	if !s.isDatabaseAvailable() {
		b.Skip("Database not available")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var users []models.User
		facades.Orm().Query().
			Where("email", "LIKE", "%@example.com").
			Where("id", ">", 1).
			Limit(10).
			Find(&users)
	}
}

// BenchmarkDBCount benchmarks COUNT queries with DB
func (s *AdvancedPerformanceTestSuite) BenchmarkDBCount(b *testing.B) {
	if !s.isDatabaseAvailable() {
		b.Skip("Database not available")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count, _ := facades.DB().Table("users").Count()
		_ = count // Use the count variable to avoid unused variable error
	}
}

// BenchmarkORMCount benchmarks COUNT queries with ORM
func (s *AdvancedPerformanceTestSuite) BenchmarkORMCount(b *testing.B) {
	if !s.isDatabaseAvailable() {
		b.Skip("Database not available")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		count, _ := facades.Orm().Query().Count()
		_ = count // Use the count variable to avoid unused variable error
	}
}

// BenchmarkDBBulkInsert benchmarks bulk insert operations with DB
func (s *AdvancedPerformanceTestSuite) BenchmarkDBBulkInsert(b *testing.B) {
	if !s.isDatabaseAvailable() {
		b.Skip("Database not available")
	}

	// Setup test data
	batchSize := 10
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data := make([]map[string]interface{}, batchSize)
		for j := 0; j < batchSize; j++ {
			data[j] = map[string]interface{}{
				"name":     fmt.Sprintf("bulk_test_%d_%d", i, j),
				"email":    fmt.Sprintf("bulk_test_%d_%d@example.com", i, j),
				"password": "password",
			}
		}
		facades.DB().Table("users").Insert(data)
	}
}

// BenchmarkORMBulkInsert benchmarks bulk insert operations with ORM
func (s *AdvancedPerformanceTestSuite) BenchmarkORMBulkInsert(b *testing.B) {
	if !s.isDatabaseAvailable() {
		b.Skip("Database not available")
	}

	// Setup test data
	batchSize := 10
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		users := make([]models.User, batchSize)
		for j := 0; j < batchSize; j++ {
			users[j] = models.User{
				Name:     fmt.Sprintf("bulk_test_%d_%d", i, j),
				Email:    fmt.Sprintf("bulk_test_%d_%d@example.com", i, j),
				Password: "password",
			}
		}
		facades.Orm().Query().Create(&users)
	}
}

// TestMemoryUsageComparison compares memory usage patterns
func (s *AdvancedPerformanceTestSuite) TestMemoryUsageComparison() {
	if !s.isDatabaseAvailable() {
		s.T().Skip("Database not available")
	}

	fmt.Println("\n=== Memory Usage Comparison ===")
	
	// Test small result sets
	dbMem := s.measureMemoryUsage(func() {
		var results []map[string]interface{}
		facades.DB().Table("users").Limit(10).Get(&results)
	})
	
	ormMem := s.measureMemoryUsage(func() {
		var users []models.User
		facades.Orm().Query().Limit(10).Find(&users)
	})
	
	fmt.Printf("Small result set (10 records):\n")
	fmt.Printf("  DB:  %d bytes allocated\n", dbMem)
	fmt.Printf("  ORM: %d bytes allocated\n", ormMem)
	if ormMem > 0 && dbMem > 0 {
		ratio := float64(ormMem) / float64(dbMem)
		fmt.Printf("  ORM uses %.2fx more memory than DB\n", ratio)
	}
	
	// Test larger result sets
	dbMemLarge := s.measureMemoryUsage(func() {
		var results []map[string]interface{}
		facades.DB().Table("users").Limit(100).Get(&results)
	})
	
	ormMemLarge := s.measureMemoryUsage(func() {
		var users []models.User
		facades.Orm().Query().Limit(100).Find(&users)
	})
	
	fmt.Printf("\nLarge result set (100 records):\n")
	fmt.Printf("  DB:  %d bytes allocated\n", dbMemLarge)
	fmt.Printf("  ORM: %d bytes allocated\n", ormMemLarge)
	if ormMemLarge > 0 && dbMemLarge > 0 {
		ratio := float64(ormMemLarge) / float64(dbMemLarge)
		fmt.Printf("  ORM uses %.2fx more memory than DB\n", ratio)
	}
	
	s.True(true, "Memory usage comparison completed")
}

// Helper methods

func (s *AdvancedPerformanceTestSuite) isDatabaseAvailable() bool {
	var result []map[string]interface{}
	err := facades.DB().Table("information_schema.tables").Limit(1).Get(&result)
	return err == nil
}

func (s *AdvancedPerformanceTestSuite) measureMemoryUsage(fn func()) uint64 {
	runtime.GC()
	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)
	
	fn()
	
	runtime.GC()
	runtime.ReadMemStats(&m2)
	
	return m2.TotalAlloc - m1.TotalAlloc
}