package performance

import (
	"testing"

	"github.com/goravel/framework/facades"
	"github.com/stretchr/testify/suite"
)

func TestPerformanceTestSuite(t *testing.T) {
	suite.Run(t, new(PerformanceTestSuite))
}

// TestBasicSetup validates that the performance test setup works
func (s *PerformanceTestSuite) TestBasicSetup() {
	// This test verifies that our test suite is properly set up
	// and can access the framework facades without errors
	
	// Even without a database connection, we should be able to access the facades
	s.NotNil(facades.DB())
	s.NotNil(facades.Orm())
	
	// Log that performance tests are ready
	s.T().Log("Performance test suite initialized successfully")
	s.T().Log("To run benchmarks: go test -bench=. ./tests/performance/")
	s.T().Log("To run with database: ensure PostgreSQL is running and database exists")
}