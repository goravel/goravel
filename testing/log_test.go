package testing

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/testing/file"
)

func TestLog(t *testing.T) {
	facades.Log.Debug("Goravel")
	facades.Log.Debugf("%s Goravel", "Hello")
	facades.Log.Info("Goravel")
	facades.Log.Infof("%s Goravel", "Hello")
	facades.Log.Warning("Goravel")
	facades.Log.Warningf("%s Goravel", "Hello")
	facades.Log.Error("Goravel")
	facades.Log.Errorf("%s Goravel", "Hello")

	assert.Equal(t, 9, file.GetLineNum(fmt.Sprintf("storage/logs/goravel-%s.log", time.Now().Format("2006-01-02"))))
	assert.Equal(t, 9, file.GetLineNum("storage/logs/goravel.log"))
	assert.Equal(t, 9, file.GetLineNum("storage/logs/test.log"))

	assert.Nil(t, os.RemoveAll("storage"))
}

func TestLogFatal(t *testing.T) {
	assert.NotPanics(t, func() {
		facades.Log.Testing()
		facades.Log.Fatal("Goravel")
		facades.Log.Fatalf("%s Goravel", "Hello")
	})
}

func TestLogPanic(t *testing.T) {
	assert.Panics(t, func() {
		facades.Log.Panic("Goravel")
	})
	assert.Panics(t, func() {
		facades.Log.Panicf("%s Goravel", "Hello")
	})

	assert.Nil(t, os.RemoveAll("storage"))
}
