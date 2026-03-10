package process

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	contractsprocess "github.com/goravel/framework/contracts/process"
)

type RunningPipe struct {
	ctx            context.Context
	commands       []*exec.Cmd
	cancel         context.CancelFunc
	pipeCommands   []*PipeCommand
	loading        bool
	loadingMessage string

	interReaders []*io.PipeReader
	interWriters []*io.PipeWriter

	stdOutputBuffers []*bytes.Buffer
	stdErrorBuffers  []*bytes.Buffer

	doneChan chan struct{}
	result   contractsprocess.Result
}

func NewRunningPipe(
	ctx context.Context,
	commands []*exec.Cmd,
	pipeCommands []*PipeCommand,
	cancel context.CancelFunc,
	interReaders []*io.PipeReader,
	interWriters []*io.PipeWriter,
	stdout, stderr []*bytes.Buffer,
	loading bool,
	loadingMessage string,
) *RunningPipe {
	pipeRunner := &RunningPipe{
		ctx:              ctx,
		commands:         commands,
		cancel:           cancel,
		pipeCommands:     pipeCommands,
		loading:          loading,
		loadingMessage:   loadingMessage,
		interReaders:     interReaders,
		interWriters:     interWriters,
		stdOutputBuffers: stdout,
		stdErrorBuffers:  stderr,
		doneChan:         make(chan struct{}),
	}

	go func(runner *RunningPipe) {
		var (
			lastCmd *exec.Cmd
			lastIdx int
			err     error
		)

		defer func() {
			if err := recover(); err != nil {
				// append panic to the last step's stderr buffer if available
				if len(runner.stdErrorBuffers) > 0 && runner.stdErrorBuffers[lastIdx] != nil {
					runner.stdErrorBuffers[lastIdx].WriteString("panic: ")
					_, _ = fmt.Fprint(runner.stdErrorBuffers[lastIdx], err)
					runner.stdErrorBuffers[lastIdx].WriteString("\n")
				}
			}
			if runner.cancel != nil {
				runner.cancel()
			}
			close(runner.doneChan)
		}()

		for i, cmd := range runner.commands {
			lastCmd = cmd
			lastIdx = i

			// Execute cmd.Wait() with spinner for this specific command
			err = runner.spinnerForCommand(i, func() error {
				return cmd.Wait()
			})

			// Close the writer that fed the next process's stdin.
			// Closing here ensures the next process sees EOF when upstream finishes.
			if i < len(runner.interWriters) {
				_ = runner.interWriters[i].Close()
			}

			if err != nil {
				for j := i + 1; j < len(runner.interWriters); j++ {
					_ = runner.interWriters[j].Close()
				}

				break
			}
		}

		exitCode := getExitCode(lastCmd, err)
		cmdStr := lastCmd.String()

		stdoutStr, stderrStr := "", ""
		if runner.stdOutputBuffers[lastIdx] != nil {
			stdoutStr = runner.stdOutputBuffers[lastIdx].String()
		}
		if runner.stdErrorBuffers[lastIdx] != nil {
			stderrStr = runner.stdErrorBuffers[lastIdx].String()
		}

		runner.result = NewResult(err, exitCode, cmdStr, stdoutStr, stderrStr)

		for _, r := range runner.interReaders {
			_ = r.Close()
		}
	}(pipeRunner)

	return pipeRunner
}

func (r *RunningPipe) PIDs() map[string]int {
	m := make(map[string]int, len(r.commands))
	for i, cmd := range r.commands {
		key := r.pipeCommands[i].key
		pid := 0
		if cmd.Process != nil {
			pid = cmd.Process.Pid
		}
		m[key] = pid
	}
	return m
}

func (r *RunningPipe) Running() bool {
	for _, cmd := range r.commands {
		if running(cmd.Process) {
			return true
		}
	}
	return false
}

func (r *RunningPipe) Done() <-chan struct{} {
	return r.doneChan
}

func (r *RunningPipe) Wait() contractsprocess.Result {
	<-r.Done()
	return r.result
}

func (r *RunningPipe) Signal(sig os.Signal) error {
	var firstErr error
	for _, cmd := range r.commands {
		if running(cmd.Process) {
			if err := signal(cmd.Process, sig); err != nil {
				if firstErr == nil {
					firstErr = err
				}
			}
		}
	}
	return firstErr
}

func (r *RunningPipe) Stop(timeout time.Duration, sig ...os.Signal) error {
	var firstErr error
	for _, cmd := range r.commands {
		if err := stop(cmd.Process, r.doneChan, timeout, sig...); err != nil {
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

func (r *RunningPipe) spinnerForCommand(index int, fn func() error) error {
	pc := r.pipeCommands[index]

	// Determine if loading should be shown for this command
	loading := pc.loading || r.loading

	// Determine the loading message for this command
	loadingMessage := pc.loadingMessage
	if loadingMessage == "" {
		// Use global message if set
		if r.loadingMessage != "" {
			loadingMessage = r.loadingMessage
		} else {
			// Generate default message from this command
			args := append([]string{pc.name}, pc.args...)
			loadingMessage = fmt.Sprintf("> %s", strings.Join(args, " "))
		}
	}

	return spinner(r.ctx, loading, loadingMessage, fn)
}
