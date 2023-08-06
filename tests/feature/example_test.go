package feature

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"goravel/tests"
)

type ExampleTestSuite struct {
	suite.Suite
	tests.TestCase
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(ExampleTestSuite))
}

// SetupTest will run before each test in the suite.
func (s *ExampleTestSuite) SetupTest() {
}

// TearDownTest will run after each test in the suite.
func (s *ExampleTestSuite) TearDownTest() {
}

func (s *ExampleTestSuite) TestIndex() {
	s.True(true)
}
