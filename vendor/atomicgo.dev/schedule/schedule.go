package schedule

import "time"

// Task holds information about the running task and can be used to stop running tasks.
type Task struct {
	stop          chan struct{}
	nextExecution time.Time
	startedAt     time.Time
}

// newTask creates a new Task.
func newTask() *Task {
	return &Task{
		stop:      make(chan struct{}),
		startedAt: time.Now(),
	}
}

// StartedAt returns the time when the scheduler was started.
func (s *Task) StartedAt() time.Time {
	return s.startedAt
}

// NextExecutionTime returns the time when the next execution will happen.
func (s *Task) NextExecutionTime() time.Time {
	return s.nextExecution
}

// ExecutesIn returns the duration until the next execution.
func (s *Task) ExecutesIn() time.Duration {
	return time.Until(s.nextExecution)
}

// IsActive returns true if the scheduler is active.
func (s *Task) IsActive() bool {
	select {
	case <-s.stop:
		return false
	default:
		return true
	}
}

// Wait blocks until the scheduler is stopped.
// After and At will stop automatically after the task is executed.
func (s *Task) Wait() {
	<-s.stop
}

// Stop stops the scheduler.
func (s *Task) Stop() {
	close(s.stop)
}

// After executes the task after the given duration.
// The function is non-blocking. If you want to wait for the task to be executed, use the Task.Wait method.
func After(d time.Duration, task func()) *Task {
	scheduler := newTask()
	scheduler.nextExecution = time.Now().Add(d)

	go func() {
		select {
		case <-time.After(d):
			task()
			scheduler.Stop()
		case <-scheduler.stop:
			return
		}
	}()

	return scheduler
}

// At executes the task at the given time.
// The function is non-blocking. If you want to wait for the task to be executed, use the Task.Wait method.
func At(t time.Time, task func()) *Task {
	scheduler := newTask()
	scheduler.nextExecution = t

	go func() {
		select {
		case <-time.After(time.Until(t)):
			task()
			scheduler.Stop()
		case <-scheduler.stop:
			return
		}
	}()

	return scheduler
}

// Every executes the task in the given interval, as long as the task function returns true.
// The function is non-blocking. If you want to wait for the task to be executed, use the Task.Wait method.
func Every(interval time.Duration, task func() bool) *Task {
	scheduler := newTask()
	scheduler.nextExecution = time.Now().Add(interval)

	ticker := time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-ticker.C:
				task()
				scheduler.nextExecution = time.Now().Add(interval)
			case <-scheduler.stop:
				ticker.Stop()
				return
			}
		}
	}()

	return scheduler
}
