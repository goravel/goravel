package process

import (
	"context"
	"io"
	"time"
)

// OutputType represents the type of output stream produced by a running process.
type OutputType int

const (
	// OutputTypeStdout indicates output written to the standard output stream.
	OutputTypeStdout OutputType = iota

	// OutputTypeStderr indicates output written to the standard error stream.
	OutputTypeStderr
)

// OnOutputFunc is a callback function invoked when the process produces output.
// The typ(OutputType) parameter indicates whether the data came from stdout or stderr,
// and line contains the raw output bytes (typically a line of text).
type OnOutputFunc func(typ OutputType, line []byte)

// Process defines an interface for configuring and running external processes.
//
// Implementations are mutable and should not be reused concurrently.
// Each method modifies the same underlying process configuration.
type Process interface {
	// DisableBuffering prevents the process's stdout and stderr from being buffered
	// in memory. This is a critical optimization for commands that produce a large
	// volume of output, especially when that output is already being handled by
	// an OnOutput streaming callback.
	//
	// CONSEQUENCE: As output is not captured, the following methods on the
	// Running and Result handles will always return an empty string:
	//   - Running.Output()
	//   - Running.ErrorOutput()
	//   - Result.Output()
	//   - Result.ErrorOutput()
	DisableBuffering() Process

	// Env adds or overrides environment variables for the process.
	// Modifies the current process configuration.
	Env(vars map[string]string) Process

	// Input sets the stdin source for the process.
	// By default, processes run without stdin input.
	Input(in io.Reader) Process

	// Path sets the working directory where the process will be executed.
	Path(path string) Process

	// Pipe creates a pipeline of commands where the output of each command is connected to the input of the next command.
	// The configurer function is used to define the sequence of commands in the pipeline.
	//
	// Note: Process configurations (timeout, context, etc.) are NOT inherited by the pipeline.
	// You must configure these settings directly on the returned Pipeline instance.
	Pipe(configurer func(Pipe)) Pipeline

	// Pool creates a pool of concurrent processes that can be executed and managed together.
	// The configurer function is used to define the commands to be executed in the pool.
	//
	// Note: Process configurations (timeout, context, etc.) are NOT inherited by the pool.
	// You must configure these settings directly on the returned PoolBuilder instance.
	Pool(configurer func(Pool)) PoolBuilder

	// Quietly suppresses all process output, discarding both stdout and stderr.
	Quietly() Process

	// OnOutput registers a handler to receive stdout and stderr output
	// while the process runs. Multiple handlers may be chained depending
	// on the implementation.
	OnOutput(handler OnOutputFunc) Process

	// Run starts the process, waits for it to complete, and returns the result.
	// If only name is provided, and the name contains special characters (like spaces, &, |),
	// the name will be added a `/bin/sh -c` or `cmd /c` wrapper to ensure correct execution.
	// This feature provides a convenient way to run complex shell commands that don't need to add the wrapper manually.
	Run(name string, arg ...string) Result

	// Start begins running the process asynchronously and returns a Running
	// handle to monitor and control its execution. The caller must later
	// wait or terminate the process explicitly.
	Start(name string, arg ...string) (Running, error)

	// Timeout sets a maximum execution duration for the process.
	// If the timeout is exceeded, the process will be terminated.
	// A zero duration disables the timeout.
	Timeout(timeout time.Duration) Process

	// TTY runs the command in an interactive TTY mode.
	//
	// This is the method you need when you're running a command that asks for
	// input, requires a password, or shows a TUI menu (like `artisan make:controller`).
	// It essentially "borrows" your terminal and gives it to the subprocess.
	//
	// Be aware of two major side effects:
	//  1. Output is NOT captured. It goes straight to your terminal. The
	//     Result object won't contain any output from the command.
	//  2. The `.Input()` method is ignored. Your live keyboard becomes the
	//     command's standard input.
	TTY() Process

	// WithContext binds the process lifecycle to the provided context.
	// If the context is canceled or reaches its deadline, the process
	// will be terminated. When combined with Timeout, the earlier of
	// the two deadlines takes effect.
	WithContext(ctx context.Context) Process

	// WithSpinner enables a loading spinner in the terminal while the process is running.
	WithSpinner(message ...string) Process
}
