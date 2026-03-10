package process

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"time"

	contractsprocess "github.com/goravel/framework/contracts/process"
	"github.com/goravel/framework/support/env"
	"github.com/goravel/framework/support/str"
)

var _ contractsprocess.Process = (*Process)(nil)

type Process struct {
	buffering      bool
	ctx            context.Context
	env            []string
	input          io.Reader
	loading        bool
	loadingMessage string
	onOutput       contractsprocess.OnOutputFunc
	path           string
	quietly        bool
	timeout        time.Duration
	tty            bool
}

func New() *Process {
	return &Process{
		ctx:       context.Background(),
		buffering: true,
	}
}

func (r *Process) DisableBuffering() contractsprocess.Process {
	r.buffering = false
	return r
}

func (r *Process) Env(vars map[string]string) contractsprocess.Process {
	if r.env == nil {
		r.env = make([]string, 0, len(vars))
	}
	for k, v := range vars {
		r.env = append(r.env, k+"="+v)
	}
	return r
}

func (r *Process) Input(in io.Reader) contractsprocess.Process {
	r.input = in
	return r
}

func (r *Process) OnOutput(handler contractsprocess.OnOutputFunc) contractsprocess.Process {
	r.onOutput = handler
	return r
}

func (r *Process) Path(path string) contractsprocess.Process {
	r.path = path
	return r
}

func (r *Process) Pipe(configurer func(pipe contractsprocess.Pipe)) contractsprocess.Pipeline {
	return NewPipe().Pipe(configurer)
}

func (r *Process) Pool(configurer func(pool contractsprocess.Pool)) contractsprocess.PoolBuilder {
	return NewPool().Pool(configurer)
}

func (r *Process) Quietly() contractsprocess.Process {
	r.quietly = true
	return r
}

func (r *Process) Run(name string, args ...string) contractsprocess.Result {
	name, args = formatCommand(name, args)
	run, err := r.start(name, args...)
	if err != nil {
		return NewResult(err, 1, "", "", "")
	}

	return run.Wait()
}

func (r *Process) Start(name string, args ...string) (contractsprocess.Running, error) {
	return r.start(name, args...)
}

func (r *Process) Timeout(timeout time.Duration) contractsprocess.Process {
	r.timeout = timeout
	return r
}

func (r *Process) TTY() contractsprocess.Process {
	r.tty = true
	return r
}

func (r *Process) WithContext(ctx context.Context) contractsprocess.Process {
	if ctx == nil {
		ctx = context.Background()
	}

	r.ctx = ctx
	return r
}

func (r *Process) WithSpinner(message ...string) contractsprocess.Process {
	r.loading = true
	if len(message) > 0 {
		r.loadingMessage = message[0]
	}

	return r
}

func (r *Process) start(name string, args ...string) (contractsprocess.Running, error) {
	ctx := r.ctx

	var cancel context.CancelFunc
	if r.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, r.timeout)
	}

	cmd := exec.CommandContext(ctx, name, args...)
	if r.path != "" {
		cmd.Dir = r.path
	}

	if len(r.env) > 0 {
		cmd.Env = append(os.Environ(), r.env...)
	}

	if r.input != nil {
		cmd.Stdin = r.input
	}

	var stdoutBuffer, stderrBuffer *bytes.Buffer

	if r.tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		setSysProcAttr(cmd)

		var stdoutWriters []io.Writer
		var stderrWriters []io.Writer

		if r.buffering {
			stdoutBuffer = &bytes.Buffer{}
			stderrBuffer = &bytes.Buffer{}
			stdoutWriters = append(stdoutWriters, stdoutBuffer)
			stderrWriters = append(stderrWriters, stderrBuffer)
		}

		if !r.quietly {
			stdoutWriters = append(stdoutWriters, os.Stdout)
			stderrWriters = append(stderrWriters, os.Stderr)
		}
		if r.onOutput != nil {
			stdoutWriters = append(stdoutWriters, NewOutputWriterForProcess(contractsprocess.OutputTypeStdout, r.onOutput))
			stderrWriters = append(stderrWriters, NewOutputWriterForProcess(contractsprocess.OutputTypeStderr, r.onOutput))
		}

		if len(stdoutWriters) > 0 {
			cmd.Stdout = io.MultiWriter(stdoutWriters...)
		}
		if len(stderrWriters) > 0 {
			cmd.Stderr = io.MultiWriter(stderrWriters...)
		}
	}

	if err := cmd.Start(); err != nil {
		if cancel != nil {
			cancel()
		}

		return nil, err
	}

	return NewRunning(ctx, cmd, cancel, stdoutBuffer, stderrBuffer, r.loading, r.loadingMessage), nil
}

func formatCommand(name string, args []string) (string, []string) {
	if len(args) == 0 && str.Of(name).Contains(" ", "&", "|") {
		if env.IsWindows() {
			args = []string{"/c", name}
			name = "cmd"
		} else {
			args = []string{"-c", name}
			name = "/bin/sh"
		}
	}

	return name, args
}
