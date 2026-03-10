package process

import (
	"context"
	"os"
	"sync"
	"time"

	contractsprocess "github.com/goravel/framework/contracts/process"
)

var _ contractsprocess.RunningPool = (*RunningPool)(nil)

type RunningPool struct {
	mu             sync.RWMutex
	ctx            context.Context
	cancel         context.CancelFunc
	done           chan struct{}
	processes      map[string]contractsprocess.Running
	results        map[string]contractsprocess.Result
	keys           []string
	loading        bool
	loadingMessage string
}

func NewRunningPool(
	ctx context.Context,
	cancel context.CancelFunc,
	keys []string,
	done chan struct{},
	loading bool,
	loadingMessage string,
) *RunningPool {
	processes := make(map[string]contractsprocess.Running, len(keys))
	results := make(map[string]contractsprocess.Result, len(keys))

	return &RunningPool{
		ctx:            ctx,
		cancel:         cancel,
		keys:           keys,
		done:           done,
		processes:      processes,
		results:        results,
		loading:        loading,
		loadingMessage: loadingMessage,
	}
}

func (r *RunningPool) PIDs() map[string]int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m := make(map[string]int, len(r.keys))
	for _, key := range r.keys {
		pid := 0
		if proc, ok := r.processes[key]; ok {
			pid = proc.PID()
		}
		m[key] = pid
	}
	return m
}

func (r *RunningPool) Running() bool {
	select {
	case <-r.done:
		return false
	default:
		return true
	}
}

func (r *RunningPool) Done() <-chan struct{} {
	return r.done
}

func (r *RunningPool) Wait() map[string]contractsprocess.Result {
	_ = r.spinner(func() error {
		<-r.Done()
		return nil
	})

	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.results
}

func (r *RunningPool) Stop(timeout time.Duration, sig ...os.Signal) error {
	if r.cancel != nil {
		r.cancel()
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var firstErr error
	for _, proc := range r.processes {
		if proc == nil {
			continue
		}
		if err := proc.Stop(timeout, sig...); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (r *RunningPool) Signal(sig os.Signal) error {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var firstErr error
	for _, proc := range r.processes {
		if proc == nil {
			continue
		}
		if err := proc.Signal(sig); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (r *RunningPool) setProcess(key string, proc contractsprocess.Running) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.processes[key] = proc
}

func (r *RunningPool) setResult(key string, res contractsprocess.Result) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.results[key] = res
}

func (r *RunningPool) spinner(fn func() error) error {
	loadingMessage := r.loadingMessage
	if loadingMessage == "" {
		loadingMessage = "Running..."
	}

	return spinner(r.ctx, r.loading, loadingMessage, fn)
}
