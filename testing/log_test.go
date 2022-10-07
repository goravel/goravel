package testing

import (
	"fmt"
	"testing"
	"time"

	"github.com/goravel/framework/facades"
	supportfile "github.com/goravel/framework/support/file"
	"github.com/goravel/framework/testing/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"goravel/bootstrap"
)

type LogTestSuite struct {
	suite.Suite
}

func TestLogTestSuite(t *testing.T) {
	bootstrap.Boot()
	supportfile.Remove("storage")
	suite.Run(t, new(LogTestSuite))
}

func (s *LogTestSuite) SetupTest() {

}

func (s *LogTestSuite) TestLog() {
	t := s.T()
	facades.Log.Debug("Goravel")
	facades.Log.Debugf("%s Goravel", "Hello")
	facades.Log.Info("Goravel")
	facades.Log.Infof("%s Goravel", "Hello")
	facades.Log.Warning("Goravel")
	facades.Log.Warningf("%s Goravel", "Hello")
	facades.Log.Error("Goravel")
	facades.Log.Errorf("%s Goravel", "Hello")

	file1 := fmt.Sprintf("storage/logs/goravel-%s.log", time.Now().Format("2006-01-02"))
	assert.True(t, supportfile.Exist(file1))
	if supportfile.Exist(file1) {
		assert.Equal(t, 9, file.GetLineNum(file1))
	}
	file2 := "storage/logs/goravel.log"
	assert.True(t, supportfile.Exist(file2))
	if supportfile.Exist(file2) {
		assert.Equal(t, 9, file.GetLineNum(file2))
	}
	file3 := "storage/logs/test.log"
	assert.True(t, supportfile.Exist(file3))
	if supportfile.Exist(file3) {
		assert.Equal(t, 9, file.GetLineNum("storage/logs/test.log"))
	}
}

func (s *LogTestSuite) TestLogFatal() {
	t := s.T()
	assert.NotPanics(t, func() {
		facades.Log.Testing(true)
		facades.Log.Fatal("Goravel")
		facades.Log.Fatalf("%s Goravel", "Hello")
		facades.Log.Testing(false)
	})
}

func (s *LogTestSuite) TestLogPanic() {
	t := s.T()
	assert.Panics(t, func() {
		facades.Log.Panic("Goravel")
	})
	assert.Panics(t, func() {
		facades.Log.Panicf("%s Goravel", "Hello")
	})
}
