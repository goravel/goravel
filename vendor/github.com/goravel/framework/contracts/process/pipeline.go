package process

import (
	"context"
	"io"
	"time"
)

// OnPipeOutputFunc is a callback function invoked when any command in the pipeline produces output.
// typ(OutputType) indicates whether the data came from stdout or stderr,
// line contains the raw output bytes,
// and key identifies which command in the pipeline produced the output.
type OnPipeOutputFunc func(typ OutputType, line []byte, key string)

// Pipeline defines a builder-style API for constructing and running a sequence
// of commands connected via OS pipes. Implementations are mutable and should
// not be used concurrently. Each configuration method returns the same
// Pipeline instance to allow fluent chaining. The Run/Start methods spawn the
// processes according to the provided builder.
type Pipeline interface {
	// DisableBuffering prevents capture of stdout/stderr into memory buffers.
	// When disabled, Result.Output and Result.ErrorOutput will be empty strings.
	DisableBuffering() Pipeline

	// Env adds or overrides environment variables for all steps.
	Env(vars map[string]string) Pipeline

	// Input sets the stdin source for the first step in the pipeline.
	Input(in io.Reader) Pipeline

	// Path sets the working directory for all steps.
	Path(path string) Pipeline

	// Pipe adds commands to the pipeline using the provided configurer function.
	// This method allows for fluent chaining of pipeline configuration methods.
	Pipe(configurer func(Pipe)) Pipeline

	// Quietly discards live stdout/stderr instead of mirroring to os.Stdout/err.
	Quietly() Pipeline

	// OnOutput registers a handler that receives line-delimited output produced
	// by each step while the pipeline runs.
	OnOutput(handler OnPipeOutputFunc) Pipeline

	// Run executes, waits for completion, and returns the final Result.
	Run() Result

	// Start starts execution asynchronously, returning a RunningPipe.
	Start() (RunningPipe, error)

	// Timeout sets a maximum duration for the entire pipeline execution.
	Timeout(timeout time.Duration) Pipeline

	// WithContext binds pipeline execution to the provided context.
	WithContext(ctx context.Context) Pipeline

	// WithSpinner enables a loading spinner in the terminal while the pipeline is running.
	WithSpinner(message ...string) Pipeline
}

// Pipe defines an interface for adding commands to a pipeline.
type Pipe interface {
	// Command creates a new command to be added to the pipeline.
	// The command's stdout will be connected to the stdin of the next command in the pipeline.
	// If only name is provided, and the name contains special characters (like spaces, &, |),
	// the name will be added a `/bin/sh -c` or `cmd /c` wrapper to ensure correct execution.
	// This feature provides a convenient way to run complex shell commands that don't need to add the wrapper manually.
	Command(name string, arg ...string) PipeCommand
}

// PipeCommand defines an interface for configuring a command within a pipeline.
type PipeCommand interface {
	// As assigns a unique string key to the command.
	// This key is used to identify the command in the output handler.
	As(key string) PipeCommand

	// WithSpinner enables a loading spinner in the terminal while the process is running.
	WithSpinner(message ...string) PipeCommand
}
