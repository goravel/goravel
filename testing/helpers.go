package testing

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"

	"github.com/gookit/color"
	"github.com/stretchr/testify/assert"
)

const Artisan = "go run ../main.go artisan "

func NotEmpty(t *testing.T, command string) {
	outStr, errStr, err := RunCommand(Artisan + command)

	assert.NotEmpty(t, outStr)
	assert.Empty(t, errStr)
	assert.NoError(t, err)
}

func Equal(t *testing.T, command, expect string) {
	outStr, _, err := RunCommand(Artisan + command)

	assert.Equal(t, expect, ClearOutput(outStr))
	//assert.Empty(t, errStr)
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
