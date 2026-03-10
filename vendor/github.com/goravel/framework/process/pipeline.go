package process

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"strconv"
	"time"

	contractsprocess "github.com/goravel/framework/contracts/process"
	"github.com/goravel/framework/errors"
)

var _ contractsprocess.Pipeline = (*Pipeline)(nil)
var _ contractsprocess.Pipe = (*Pipe)(nil)
var _ contractsprocess.PipeCommand = (*PipeCommand)(nil)

func NewPipe() *Pipeline {
	return &Pipeline{
		ctx:       context.Background(),
		buffering: true,
	}
}

type Pipeline struct {
	ctx            context.Context
	input          io.Reader
	env            []string
	timeout        time.Duration
	onOutput       contractsprocess.OnPipeOutputFunc
	quietly        bool
	path           string
	buffering      bool
	loading        bool
	loadingMessage string

	pipeConfigurer func(pipe contractsprocess.Pipe)
}

func (r *Pipeline) DisableBuffering() contractsprocess.Pipeline {
	r.buffering = false
	return r
}

func (r *Pipeline) Input(in io.Reader) contractsprocess.Pipeline {
	r.input = in
	return r
}

func (r *Pipeline) Env(vars map[string]string) contractsprocess.Pipeline {
	for k, v := range vars {
		r.env = append(r.env, k+"="+v)
	}
	return r
}

func (r *Pipeline) Path(path string) contractsprocess.Pipeline {
	r.path = path
	return r
}

func (r *Pipeline) Pipe(configurer func(pipe contractsprocess.Pipe)) contractsprocess.Pipeline {
	r.pipeConfigurer = configurer
	return r
}

func (r *Pipeline) Timeout(timeout time.Duration) contractsprocess.Pipeline {
	r.timeout = timeout
	return r
}

func (r *Pipeline) Quietly() contractsprocess.Pipeline {
	r.quietly = true
	return r
}

func (r *Pipeline) OnOutput(onOutput contractsprocess.OnPipeOutputFunc) contractsprocess.Pipeline {
	r.onOutput = onOutput
	return r
}

func (r *Pipeline) Run() contractsprocess.Result {
	run, err := r.start(r.pipeConfigurer)
	if err != nil {
		return NewResult(err, 1, "", "", "")
	}

	return run.Wait()
}

func (r *Pipeline) Start() (contractsprocess.RunningPipe, error) {
	return r.start(r.pipeConfigurer)
}

func (r *Pipeline) WithContext(ctx context.Context) contractsprocess.Pipeline {
	if ctx == nil {
		ctx = context.Background()
	}

	r.ctx = ctx
	return r
}

func (r *Pipeline) WithSpinner(message ...string) contractsprocess.Pipeline {
	r.loading = true
	if len(message) > 0 {
		r.loadingMessage = message[0]
	}

	return r
}

func (r *Pipeline) start(configurer func(contractsprocess.Pipe)) (contractsprocess.RunningPipe, error) {
	if configurer == nil {
		return nil, errors.ProcessPipeNilConfigurer
	}

	pipe := &Pipe{}
	configurer(pipe)

	pipeCommands := pipe.commands
	if len(pipeCommands) == 0 {
		return nil, errors.ProcessPipelineEmpty
	}

	ctx := r.ctx
	if ctx == nil {
		ctx = context.Background()
	}

	var cancel context.CancelFunc
	if r.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, r.timeout)
	}

	commands := make([]*exec.Cmd, len(pipeCommands))
	for i, step := range pipeCommands {
		cmd := exec.CommandContext(ctx, step.name, step.args...)
		if r.path != "" {
			cmd.Dir = r.path
		}
		setSysProcAttr(cmd)

		if len(r.env) > 0 {
			cmd.Env = append(os.Environ(), r.env...)
		}

		commands[i] = cmd
	}

	// Prepare pipe connections between adjacent commands and configure stdout/stderr writers.
	// For i < len(commands)-1: command[i].Stdout -> pipeWriter -> pipeReader -> command[i+1].Stdin
	// Also, each command's stdout/stderr may also be copied to buffers, os.Stdout/os.Stderr and
	// an onOutput callback via MultiWriter.

	interReaders := make([]*io.PipeReader, 0, len(commands)-1)
	interWriters := make([]*io.PipeWriter, 0, len(commands)-1)

	stdoutBuffers := make([]*bytes.Buffer, len(commands))
	stderrBuffers := make([]*bytes.Buffer, len(commands))

	if len(commands) > 0 && r.input != nil {
		commands[0].Stdin = r.input
	}

	for i, cmd := range commands {
		var stdoutBuffer, stderrBuffer *bytes.Buffer
		var stdoutWriters []io.Writer
		var stderrWriters []io.Writer

		if r.buffering {
			stdoutBuffer = &bytes.Buffer{}
			stderrBuffer = &bytes.Buffer{}
			stdoutWriters = append(stdoutWriters, stdoutBuffer)
			stderrWriters = append(stderrWriters, stderrBuffer)
			stdoutBuffers[i] = stdoutBuffer
			stderrBuffers[i] = stderrBuffer
		}

		if !r.quietly {
			stdoutWriters = append(stdoutWriters, os.Stdout)
			stderrWriters = append(stderrWriters, os.Stderr)
		}

		if r.onOutput != nil {
			stdoutWriters = append(stdoutWriters, NewOutputWriterForPipe(contractsprocess.OutputTypeStdout, pipeCommands[i].key, r.onOutput))
			stderrWriters = append(stderrWriters, NewOutputWriterForPipe(contractsprocess.OutputTypeStderr, pipeCommands[i].key, r.onOutput))
		}

		// If this is not the last command, create a pipe to the next command and include the pipe writer
		// in this command's stdout MultiWriter — but ONLY if the next command does not already have stdin set.
		if i < len(commands)-1 {
			if commands[i+1].Stdin == nil {
				pr, pw := io.Pipe()
				interReaders = append(interReaders, pr)
				interWriters = append(interWriters, pw)
				stdoutWriters = append(stdoutWriters, pw)
				// set next command's stdin to the pipe reader
				commands[i+1].Stdin = pr
			}
		}

		if len(stdoutWriters) > 0 {
			cmd.Stdout = io.MultiWriter(stdoutWriters...)
		}

		if len(stderrWriters) > 0 {
			cmd.Stderr = io.MultiWriter(stderrWriters...)
		}
	}

	started := 0
	for i, cmd := range commands {
		if err := cmd.Start(); err != nil {
			if cancel != nil {
				cancel()
			}

			for j := 0; j < started; j++ {
				if commands[j].Process != nil {
					_ = kill(commands[j].Process)
				}
			}
			for _, w := range interWriters {
				_ = w.Close()
			}
			for _, r := range interReaders {
				_ = r.Close()
			}
			return nil, errors.ProcessPipelineStartFailed.Args(err)
		}
		started = i + 1
	}

	return NewRunningPipe(ctx, commands, pipeCommands, cancel, interReaders, interWriters, stdoutBuffers, stderrBuffers, r.loading, r.loadingMessage), nil
}

type Pipe struct {
	commands []*PipeCommand
}

func (r *Pipe) Command(name string, args ...string) contractsprocess.PipeCommand {
	name, args = formatCommand(name, args)
	command := NewPipeCommand(strconv.Itoa(len(r.commands)), name, args)
	r.commands = append(r.commands, command)
	return command
}

type PipeCommand struct {
	key            string
	name           string
	args           []string
	loading        bool
	loadingMessage string
}

func NewPipeCommand(key, name string, args []string) *PipeCommand {
	return &PipeCommand{
		key:  key,
		name: name,
		args: args,
	}
}

func (r *PipeCommand) As(key string) contractsprocess.PipeCommand {
	r.key = key
	return r
}

func (r *PipeCommand) WithSpinner(message ...string) contractsprocess.PipeCommand {
	r.loading = true
	if len(message) > 0 {
		r.loadingMessage = message[0]
	}

	return r
}
