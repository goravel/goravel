package testing

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	"github.com/gookit/color"
	"github.com/stretchr/testify/assert"
)

func TestKeyGenerate(t *testing.T) {
	Command(t, "go run ../main.go artisan key:generate", "Application key set successfully")
	Command(t, "go run ../main.go artisan key:generate", "Exist application key")
}

func Command(t *testing.T, command, output string) {
	outStr, errStr, err := RunCommand(command)

	assert.Equal(t, output, ClearOutput(outStr))
	assert.Empty(t, errStr)
	assert.NoError(t, err)
}

func RunCommand(command string) (outStr string, errStr string, err error) {
	commandArr := strings.Split(command, " ")
	cmd := exec.Command(commandArr[0], commandArr[1:]...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	outStr, errStr = string(stdout.Bytes()), string(stderr.Bytes())

	return
}

func ClearOutput(output string) string {
	return strings.ReplaceAll(color.ClearCode(output), "\n", "")
}
