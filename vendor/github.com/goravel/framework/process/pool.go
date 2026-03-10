package process

import (
	"context"
	"io"
	"maps"
	"strconv"
	"sync"
	"time"

	contractsprocess "github.com/goravel/framework/contracts/process"
	"github.com/goravel/framework/errors"
)

var _ contractsprocess.PoolBuilder = (*PoolBuilder)(nil)
var _ contractsprocess.Pool = (*Pool)(nil)
var _ contractsprocess.PoolCommand = (*PoolCommand)(nil)

type PoolBuilder struct {
	concurrency    int
	ctx            context.Context
	loading        bool
	loadingMessage string
	onOutput       contractsprocess.OnPoolOutputFunc
	timeout        time.Duration

	poolConfigurer func(pool contractsprocess.Pool)
}

func NewPool() *PoolBuilder {
	return &PoolBuilder{ctx: context.Background()}
}

func (r *PoolBuilder) Concurrency(n int) contractsprocess.PoolBuilder {
	r.concurrency = n
	return r
}

func (r *PoolBuilder) OnOutput(handler contractsprocess.OnPoolOutputFunc) contractsprocess.PoolBuilder {
	r.onOutput = handler
	return r
}

func (r *PoolBuilder) Pool(configurer func(pool contractsprocess.Pool)) contractsprocess.PoolBuilder {
	r.poolConfigurer = configurer
	return r
}

func (r *PoolBuilder) Run() (map[string]contractsprocess.Result, error) {
	run, err := r.start(r.poolConfigurer)
	if err != nil {
		return nil, err
	}
	return run.Wait(), nil
}

func (r *PoolBuilder) Start() (contractsprocess.RunningPool, error) {
	return r.start(r.poolConfigurer)
}

func (r *PoolBuilder) Timeout(timeout time.Duration) contractsprocess.PoolBuilder {
	r.timeout = timeout
	return r
}

func (r *PoolBuilder) WithContext(ctx context.Context) contractsprocess.PoolBuilder {
	if ctx == nil {
		ctx = context.Background()
	}
	r.ctx = ctx
	return r
}

func (r *PoolBuilder) WithSpinner(message ...string) contractsprocess.PoolBuilder {
	r.loading = true
	if len(message) > 0 {
		r.loadingMessage = message[0]
	}
	return r
}

type result struct {
	key string
	res contractsprocess.Result
}

// start initiates the execution of all configured commands concurrently based on the
// concurrency limit, but does not wait for them to complete.
//
// This method is non-blocking. It returns a RunningPool instance immediately, which
// can be used to wait for the completion of all processes and retrieve their results.
//
// The core concurrency pattern is as follows:
//  1. A job channel (`jobCh`) distributes commands to a pool of worker goroutines.
//  2. Workers pick up a job, start the process, and wait synchronously for it to finish.
//     This synchronous wait ensures the concurrency limit is strictly respected.
//  3. A result channel (`resultCh`) collects the outcome (success/failure) of each command.
//  4. A separate "collector" goroutine safely populates the RunningPool's internal map
//     from the result channel to avoid concurrent map write panics.
//  5. A background orchestrator waits for all workers and results to finish, then
//     cleanly closes resources and signals the `done` channel.
func (r *PoolBuilder) start(configurer func(contractsprocess.Pool)) (contractsprocess.RunningPool, error) {
	if configurer == nil {
		return nil, errors.ProcessPoolNilConfigurer
	}

	pool := &Pool{}
	configurer(pool)

	commands := pool.commands
	if len(commands) == 0 {
		return nil, errors.ProcessPipelineEmpty
	}

	ctx := r.ctx
	var cancel context.CancelFunc
	if r.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, r.timeout)
	} else {
		ctx, cancel = context.WithCancel(ctx)
	}

	jobCh := make(chan *PoolCommand, len(commands))
	resultCh := make(chan result, len(commands))
	done := make(chan struct{})

	keys := make([]string, len(commands))
	for i, cmd := range commands {
		keys[i] = cmd.key
	}

	runningPool := NewRunningPool(ctx, cancel, keys, done, r.loading, r.loadingMessage)

	concurrency := r.concurrency
	if concurrency <= 0 || concurrency > len(commands) {
		concurrency = len(commands)
	}

	var resultsWg sync.WaitGroup
	var workersWg sync.WaitGroup

	// The results collector goroutine centralizes writing to the results map
	// to avoid race conditions, as map writes are not concurrent-safe.
	// It waits for all expected results before exiting.
	resultsWg.Add(len(commands))
	go func() {
		for i := 0; i < len(commands); i++ {
			rc := <-resultCh
			runningPool.setResult(rc.key, rc.res)
			resultsWg.Done()
		}
	}()

	for i := 0; i < concurrency; i++ {
		workersWg.Add(1)
		go func() {
			defer workersWg.Done()
			for command := range jobCh {
				if ctx.Err() != nil {
					// If the pool was stopped (Stop() called or timeout reached), we skip execution.
					// We must still send a result to ensure resultsWg decrements correctly.
					resultCh <- result{
						key: command.key,
						res: NewResult(ctx.Err(), -1, "", "", ""),
					}
					continue
				}
				cmdCtx := command.ctx
				if cmdCtx == nil {
					cmdCtx = ctx
				}

				proc := New().WithContext(cmdCtx).Path(command.path).Env(command.env).Input(command.input)
				if command.quietly {
					proc = proc.Quietly()
				}
				if !command.buffering {
					proc = proc.DisableBuffering()
				}
				if command.timeout > 0 {
					proc = proc.Timeout(command.timeout)
				}
				if r.onOutput != nil {
					proc = proc.OnOutput(func(typ contractsprocess.OutputType, line []byte) {
						r.onOutput(typ, line, command.key)
					})
				}

				run, err := proc.Start(command.name, command.args...)

				if err != nil {
					resultCh <- result{key: command.key, res: NewResult(err, -1, command.name, "", "")}
				} else {
					runningPool.setProcess(command.key, run)
					res := run.Wait()
					resultCh <- result{key: command.key, res: res}
				}
			}
		}()
	}

	for _, command := range commands {
		jobCh <- command
	}
	close(jobCh)

	// This goroutine orchestrates the clean shutdown. It waits for all workers
	// to finish processing jobs, then waits for all results to be collected.
	// Finally, it cancels the context (if a timeout was set) and signals
	// completion by closing the `done` channel.
	go func() {
		workersWg.Wait()
		resultsWg.Wait()
		if cancel != nil {
			cancel()
		}
		close(done)
	}()

	return runningPool, nil
}

type Pool struct {
	commands []*PoolCommand
}

func (r *Pool) Command(name string, args ...string) contractsprocess.PoolCommand {
	name, args = formatCommand(name, args)
	command := NewPoolCommand(strconv.Itoa(len(r.commands)), name, args)
	r.commands = append(r.commands, command)
	return command
}

type PoolCommand struct {
	args      []string
	buffering bool
	ctx       context.Context
	env       map[string]string
	input     io.Reader
	key       string
	name      string
	path      string
	quietly   bool
	timeout   time.Duration
}

func NewPoolCommand(key, name string, args []string) *PoolCommand {
	return &PoolCommand{
		key:       key,
		name:      name,
		args:      args,
		buffering: true,
		env:       make(map[string]string),
	}
}

func (r *PoolCommand) As(key string) contractsprocess.PoolCommand {
	r.key = key
	return r
}

func (r *PoolCommand) DisableBuffering() contractsprocess.PoolCommand {
	r.buffering = false
	return r
}

func (r *PoolCommand) Env(vars map[string]string) contractsprocess.PoolCommand {
	maps.Copy(r.env, vars)
	return r
}

func (r *PoolCommand) Input(in io.Reader) contractsprocess.PoolCommand {
	r.input = in
	return r
}

func (r *PoolCommand) Path(path string) contractsprocess.PoolCommand {
	r.path = path
	return r
}

func (r *PoolCommand) Quietly() contractsprocess.PoolCommand {
	r.quietly = true
	return r
}

func (r *PoolCommand) Timeout(timeout time.Duration) contractsprocess.PoolCommand {
	r.timeout = timeout
	return r
}

func (r *PoolCommand) WithContext(ctx context.Context) contractsprocess.PoolCommand {
	r.ctx = ctx
	return r
}
