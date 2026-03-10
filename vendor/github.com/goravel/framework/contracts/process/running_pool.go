package process

import (
	"os"
	"time"
)

// RunningPool is a handle to a collection of concurrently running processes
// spawned by PoolBuilder. It provides process introspection (PIDs, Running status),
// coordinated lifecycle controls (Signal, Stop), and a completion mechanism (Done, Wait).
type RunningPool interface {
	// Done returns a channel that is closed when all pool processes have finished
	// and their results have been collected. This allows non-blocking checks or
	// select-based waits across multiple pools.
	Done() <-chan struct{}

	// PIDs returns a map of process IDs keyed by the command keys supplied during
	// pool configuration. If a process failed to start, its PID will be 0.
	PIDs() map[string]int

	// Running reports whether the pool is still executing. It uses a non-blocking
	// select on the Done channel to avoid races or blocking calls.
	Running() bool

	// Signal sends an OS signal to each running process in the pool. The first
	// error encountered is returned, though all processes are still signaled.
	// Processes that are nil or have already terminated are skipped silently.
	Signal(sig os.Signal) error

	// Stop attempts to gracefully terminate each process in the pool using the
	// specified timeout and optional signal(s). The first error encountered is
	// returned, but all processes receive the stop request. If a process does
	// not exit within the timeout, it may be forcefully killed depending on the
	// platform implementation.
	Stop(timeout time.Duration, sig ...os.Signal) error

	// Wait blocks until the Done channel is closed and returns the final results
	// map, keyed by the command keys provided during pool configuration.
	Wait() map[string]Result
}
