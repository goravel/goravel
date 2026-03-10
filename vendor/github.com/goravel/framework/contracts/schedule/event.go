package schedule

import "time"

type Event interface {
	// At schedule the event to run at the specified time.
	At(time string) Event
	// Cron schedule the event using the given Cron expression.
	Cron(expression string) Event
	// Daily schedule the event to run daily.
	Daily() Event
	// DailyAt schedule the event to run daily at a given time (10:00, 19:30, etc).
	DailyAt(time string) Event
	// TwiceDaily schedule the event to run twice a day,default at 1:00 and 13:00.
	TwiceDaily(hours ...int) Event
	// DelayIfStillRunning if the event is still running, the event will be delayed.
	DelayIfStillRunning() Event
	// EverySecond schedule the event to run every second.
	EverySecond() Event
	// EveryTwoSeconds schedule the event to run every two seconds.
	EveryTwoSeconds() Event
	// EveryFiveSeconds schedule the event to run every five seconds.
	EveryFiveSeconds() Event
	// EveryTenSeconds schedule the event to run every ten seconds.
	EveryTenSeconds() Event
	// EveryFifteenSeconds schedule the event to run every fifteen seconds.
	EveryFifteenSeconds() Event
	// EveryTwentySeconds schedule the event to run every twenty seconds.
	EveryTwentySeconds() Event
	// EveryThirtySeconds schedule the event to run every thirty seconds.
	EveryThirtySeconds() Event
	// EveryMinute schedule the event to run every minute.
	EveryMinute() Event
	// EveryTwoMinutes schedule the event to run every two minutes.
	EveryTwoMinutes() Event
	// EveryThreeMinutes schedule the event to run every three minutes.
	EveryThreeMinutes() Event
	// EveryFourMinutes schedule the event to run every four minutes.
	EveryFourMinutes() Event
	// EveryFiveMinutes schedule the event to run every five minutes.
	EveryFiveMinutes() Event
	// EveryTenMinutes schedule the event to run every ten minutes.
	EveryTenMinutes() Event
	// EveryFifteenMinutes schedule the event to run every fifteen minutes.
	EveryFifteenMinutes() Event
	// EveryThirtyMinutes schedule the event to run every thirty minutes.
	EveryThirtyMinutes() Event
	// EveryTwoHours schedule the event to run every two hours.
	EveryTwoHours() Event
	// EveryThreeHours schedule the event to run every three hours.
	EveryThreeHours() Event
	// EveryFourHours schedule the event to run every four hours.
	EveryFourHours() Event
	// EverySixHours schedule the event to run every six hours.
	EverySixHours() Event
	// GetCron get cron expression.
	GetCron() string
	// GetCommand get the command.
	GetCommand() string
	// GetCallback get callback.
	GetCallback() func()
	// GetName get name.
	GetName() string
	// GetSkipIfStillRunning get skipIfStillRunning bool.
	GetSkipIfStillRunning() bool
	// GetDelayIfStillRunning get delayIfStillRunning bool.
	GetDelayIfStillRunning() bool
	// Hourly schedule the event to run hourly.
	Hourly() Event
	// HourlyAt schedule the event to run hourly at a given offset in the hour.
	HourlyAt(offset []string) Event
	// Days schedule the event to run on specific days of the week.
	Days(days ...time.Weekday) Event
	// Weekdays schedule the event to run on weekdays (Monday to Friday).
	Weekdays() Event
	// Weekends schedule the event to run on weekends (Saturday and Sunday).
	Weekends() Event
	// Mondays schedule the event to run on Mondays.
	Mondays() Event
	// Tuesdays schedule the event to run on Tuesdays.
	Tuesdays() Event
	// Wednesdays schedule the event to run on Wednesdays.
	Wednesdays() Event
	// Thursdays schedule the event to run on Thursdays.
	Thursdays() Event
	// Fridays schedule the event to run on Fridays.
	Fridays() Event
	// Saturdays schedule the event to run on Saturdays.
	Saturdays() Event
	// Sundays schedule the event to run on Sundays.
	Sundays() Event
	// Weekly schedule the event to run weekly.
	Weekly() Event
	// Monthly schedule the event to run monthly.
	Monthly() Event
	// Quarterly schedule the event to run quarterly.
	Quarterly() Event
	// Yearly schedule the event to run yearly.
	Yearly() Event
	// IsOnOneServer get isOnOneServer bool.
	IsOnOneServer() bool
	// Name set the event name.
	Name(name string) Event
	// OnOneServer only allow the event to run on one server for each cron expression.
	OnOneServer() Event
	// SkipIfStillRunning if the event is still running, the event will be skipped.
	SkipIfStillRunning() Event
}
