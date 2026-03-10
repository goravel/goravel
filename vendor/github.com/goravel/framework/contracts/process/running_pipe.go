package process

import (
	"os"
	"time"
)

// RunningPipe represents a running pipeline of OS commands connected together
// via pipes. It exposes methods to observe state (PIDs, Running, Done), wait
// for completion (Wait), and control execution (Stop, Signal). Implementations
// should be safe for read-only use from multiple goroutines.
type RunningPipe interface {
	// PIDs returns a mapping of step keys to their process IDs. Keys are the
	// identifiers assigned to each step (via PipeCommand.As or default numeric index).
	PIDs() map[string]int
	// Running reports whether any process in the pipeline is still running.
	Running() bool
	// Done returns a channel that is closed when the pipeline finishes
	// (successfully or with failure). It is safe to use in select statements.
	Done() <-chan struct{}
	// Wait blocks until the pipeline completes and returns the aggregated
	// result of the last step.
	Wait() Result
	// Stop attempts to gracefully stop all processes. On Unix this typically
	// sends SIGTERM then SIGKILL after the timeout. On Windows, processes are
	// terminated immediately. Returns the first error encountered, if any.
	Stop(timeout time.Duration, sig ...os.Signal) error
	// Signal sends the provided signal to all running processes where supported
	// by the platform. On Windows, some signals may be a no-op.
	Signal(sig os.Signal) error
}
