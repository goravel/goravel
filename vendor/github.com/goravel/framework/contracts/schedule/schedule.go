package schedule

import (
	"context"
)

type Schedule interface {
	// Call add a new callback event to the schedule.
	Call(callback func()) Event
	// Command adds a new Artisan command event to the schedule.
	Command(command string) Event
	// Events returns all registered schedule events.
	Events() []Event
	// Register schedules.
	Register(events []Event)
	// Run schedules.
	Run()
	// Shutdown schedules.
	Shutdown(ctx ...context.Context) error
}
