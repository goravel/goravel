package schedule

import (
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/goravel/framework/contracts/schedule"
)

type Event struct {
	callback            func()
	command             string
	cron                string
	name                string
	delayIfStillRunning bool
	onOneServer         bool
	skipIfStillRunning  bool
}

func NewCallbackEvent(callback func()) *Event {
	return &Event{
		callback: callback,
		cron:     "* * * * *",
	}
}

func NewCommandEvent(command string) *Event {
	return &Event{
		command: command,
		name:    command,
		cron:    "* * * * *",
	}
}

// At schedule the command at a given time.
func (r *Event) At(time string) schedule.Event {
	return r.DailyAt(time)
}

// Cron The Cron expression representing the event's frequency.
func (r *Event) Cron(expression string) schedule.Event {
	r.cron = expression

	return r
}

// Daily schedule the event to run daily.
func (r *Event) Daily() schedule.Event {
	event := r.Cron(r.spliceIntoPosition(1, "0"))

	return event.Cron(r.spliceIntoPosition(2, "0"))
}

// DailyAt schedule the event to run daily at a given time (10:00, 19:30, etc).
func (r *Event) DailyAt(time string) schedule.Event {
	segments := strings.Split(time, ":")
	event := r.Cron(r.spliceIntoPosition(2, segments[0]))

	if len(segments) == 2 {
		return event.Cron(r.spliceIntoPosition(1, segments[1]))
	} else {
		return event.Cron(r.spliceIntoPosition(1, "0"))
	}
}

// TwiceDaily schedule the event to run twice a day,default at 1:00 and 13:00.
func (r *Event) TwiceDaily(hours ...int) schedule.Event {
	first, second := 1, 13
	if len(hours) >= 2 {
		first, second = hours[0], hours[1]
	}

	return r.Cron(r.spliceIntoPosition(2, simplifyRanges(first, second)))
}

// DelayIfStillRunning do not allow the event to overlap each other.
func (r *Event) DelayIfStillRunning() schedule.Event {
	r.delayIfStillRunning = true

	return r
}

// EverySecond schedule the event to run every second.
func (r *Event) EverySecond() schedule.Event {
	return r.Cron(r.spliceIntoPosition(0, "*"))
}

// EveryTwoSeconds schedule the event to run every two seconds.
func (r *Event) EveryTwoSeconds() schedule.Event {
	return r.Cron(r.spliceIntoPosition(0, "*/2"))
}

// EveryFiveSeconds schedule the event to run every five seconds.
func (r *Event) EveryFiveSeconds() schedule.Event {
	return r.Cron(r.spliceIntoPosition(0, "*/5"))
}

// EveryTenSeconds schedule the event to run every ten seconds.
func (r *Event) EveryTenSeconds() schedule.Event {
	return r.Cron(r.spliceIntoPosition(0, "*/10"))
}

// EveryFifteenSeconds schedule the event to run every fifteen seconds.
func (r *Event) EveryFifteenSeconds() schedule.Event {
	return r.Cron(r.spliceIntoPosition(0, "*/15"))
}

// EveryTwentySeconds schedule the event to run every twenty seconds.
func (r *Event) EveryTwentySeconds() schedule.Event {
	return r.Cron(r.spliceIntoPosition(0, "*/20"))
}

// EveryThirtySeconds schedule the event to run every thirty seconds.
func (r *Event) EveryThirtySeconds() schedule.Event {
	return r.Cron(r.spliceIntoPosition(0, "*/30"))
}

// EveryMinute schedule the event to run every minute.
func (r *Event) EveryMinute() schedule.Event {
	return r.Cron(r.spliceIntoPosition(1, "*"))
}

// EveryTwoMinutes schedule the event to run every two minutes.
func (r *Event) EveryTwoMinutes() schedule.Event {
	return r.Cron(r.spliceIntoPosition(1, "*/2"))
}

// EveryThreeMinutes schedule the event to run every three minutes.
func (r *Event) EveryThreeMinutes() schedule.Event {
	return r.Cron(r.spliceIntoPosition(1, "*/3"))
}

// EveryFourMinutes schedule the event to run every four minutes.
func (r *Event) EveryFourMinutes() schedule.Event {
	return r.Cron(r.spliceIntoPosition(1, "*/4"))
}

// EveryFiveMinutes schedule the event to run every five minutes.
func (r *Event) EveryFiveMinutes() schedule.Event {
	return r.Cron(r.spliceIntoPosition(1, "*/5"))
}

// EveryTenMinutes schedule the event to run every ten minutes.
func (r *Event) EveryTenMinutes() schedule.Event {
	return r.Cron(r.spliceIntoPosition(1, "*/10"))
}

// EveryFifteenMinutes schedule the event to run every fifteen minutes.
func (r *Event) EveryFifteenMinutes() schedule.Event {
	return r.Cron(r.spliceIntoPosition(1, "*/15"))
}

// EveryThirtyMinutes schedule the event to run every thirty minutes.
func (r *Event) EveryThirtyMinutes() schedule.Event {
	return r.Cron(r.spliceIntoPosition(1, "*/30"))
}

// EveryTwoHours schedule the event to run every two hours.
func (r *Event) EveryTwoHours() schedule.Event {
	event := r.Cron(r.spliceIntoPosition(1, "0"))

	return event.Cron(r.spliceIntoPosition(2, "*/2"))
}

// EveryThreeHours schedule the event to run every three hours.
func (r *Event) EveryThreeHours() schedule.Event {
	event := r.Cron(r.spliceIntoPosition(1, "0"))

	return event.Cron(r.spliceIntoPosition(2, "*/3"))
}

// EveryFourHours schedule the event to run every four hours.
func (r *Event) EveryFourHours() schedule.Event {
	event := r.Cron(r.spliceIntoPosition(1, "0"))

	return event.Cron(r.spliceIntoPosition(2, "*/4"))
}

// EverySixHours schedule the event to run every six hours.
func (r *Event) EverySixHours() schedule.Event {
	event := r.Cron(r.spliceIntoPosition(1, "0"))

	return event.Cron(r.spliceIntoPosition(2, "*/6"))
}

func (r *Event) GetCron() string {
	return r.cron
}

func (r *Event) GetCommand() string {
	return r.command
}

func (r *Event) GetCallback() func() {
	return r.callback
}

func (r *Event) GetName() string {
	return r.name
}

func (r *Event) GetSkipIfStillRunning() bool {
	return r.skipIfStillRunning
}

func (r *Event) GetDelayIfStillRunning() bool {
	return r.delayIfStillRunning
}

// Hourly schedule the event to run hourly.
func (r *Event) Hourly() schedule.Event {
	return r.Cron(r.spliceIntoPosition(1, "0"))
}

// HourlyAt schedule the event to run hourly at a given offset in the hour.
func (r *Event) HourlyAt(offset []string) schedule.Event {
	return r.Cron(r.spliceIntoPosition(1, strings.Join(offset, ",")))
}

// Days schedule the event to run on specific days of the week.
func (r *Event) Days(days ...time.Weekday) schedule.Event {
	if len(days) == 0 {
		return r
	}

	return r.Cron(r.spliceIntoPosition(5, simplifyRanges(days...)))
}

// Weekdays schedule the event to run on weekdays (Monday to Friday).
func (r *Event) Weekdays() schedule.Event {
	return r.Cron(r.spliceIntoPosition(5, formatRange(time.Monday, time.Friday)))
}

// Weekends schedule the event to run on weekends (Saturday and Sunday).
func (r *Event) Weekends() schedule.Event {
	return r.Cron(r.spliceIntoPosition(5, simplifyRanges(time.Saturday, time.Sunday)))
}

// Mondays schedule the event to run on Mondays.
func (r *Event) Mondays() schedule.Event {
	return r.Cron(r.spliceIntoPosition(5, simplifyRanges(time.Monday)))
}

// Tuesdays schedule the event to run on Tuesdays.
func (r *Event) Tuesdays() schedule.Event {
	return r.Cron(r.spliceIntoPosition(5, simplifyRanges(time.Tuesday)))
}

// Wednesdays schedule the event to run on Wednesdays.
func (r *Event) Wednesdays() schedule.Event {
	return r.Cron(r.spliceIntoPosition(5, simplifyRanges(time.Wednesday)))
}

// Thursdays schedule the event to run on Thursdays.
func (r *Event) Thursdays() schedule.Event {
	return r.Cron(r.spliceIntoPosition(5, simplifyRanges(time.Thursday)))
}

// Fridays schedule the event to run on Fridays.
func (r *Event) Fridays() schedule.Event {
	return r.Cron(r.spliceIntoPosition(5, simplifyRanges(time.Friday)))
}

// Saturdays schedule the event to run on Saturdays.
func (r *Event) Saturdays() schedule.Event {
	return r.Cron(r.spliceIntoPosition(5, simplifyRanges(time.Saturday)))
}

// Sundays schedule the event to run on Sundays.
func (r *Event) Sundays() schedule.Event {
	return r.Cron(r.spliceIntoPosition(5, simplifyRanges(time.Sunday)))
}

// Weekly schedule the event to run weekly.
func (r *Event) Weekly() schedule.Event {
	return r.Cron(r.spliceIntoPosition(1, "0")).
		Cron(r.spliceIntoPosition(2, "0")).
		Cron(r.spliceIntoPosition(5, "0"))
}

// Monthly schedule the event to run monthly.
func (r *Event) Monthly() schedule.Event {
	return r.Cron(r.spliceIntoPosition(1, "0")).
		Cron(r.spliceIntoPosition(2, "0")).
		Cron(r.spliceIntoPosition(3, "1"))
}

// Quarterly schedule the event to run quarterly.
func (r *Event) Quarterly() schedule.Event {
	return r.Cron(r.spliceIntoPosition(1, "0")).
		Cron(r.spliceIntoPosition(2, "0")).
		Cron(r.spliceIntoPosition(3, "1")).
		Cron(r.spliceIntoPosition(4, "1-12/3"))
}

// Yearly schedule the event to run yearly.
func (r *Event) Yearly() schedule.Event {
	return r.Cron(r.spliceIntoPosition(1, "0")).
		Cron(r.spliceIntoPosition(2, "0")).
		Cron(r.spliceIntoPosition(3, "1")).
		Cron(r.spliceIntoPosition(4, "1"))
}

func (r *Event) IsOnOneServer() bool {
	return r.onOneServer
}

func (r *Event) Name(name string) schedule.Event {
	r.name = name

	return r
}

func (r *Event) OnOneServer() schedule.Event {
	r.onOneServer = true

	return r
}

// SkipIfStillRunning do not allow the event to overlap each other.
func (r *Event) SkipIfStillRunning() schedule.Event {
	r.skipIfStillRunning = true

	return r
}

// spliceIntoPosition splice the given value into the given position of the expression.
func (r *Event) spliceIntoPosition(position int, value string) string {
	segments := strings.Split(r.GetCron(), " ")

	if position == 0 {
		// if position is 0, it indicates a second-level cron expression.
		segments = append([]string{value}, segments...)

		return strings.Join(segments, " ")
	}

	segments[position-1] = value

	return strings.Join(segments, " ")
}

// simplifyRanges simplifies a list of integers into a string representation of ranges.
// For example, [1, 2, 3, 5, 6, 8] becomes "1-3,5-6,8".
func simplifyRanges[T ~int | ~int8 | ~int16 | ~int32 | ~int64](nums ...T) string {
	if len(nums) == 0 {
		return ""
	}

	slices.Sort(nums)
	nums = slices.Compact(nums)

	var (
		result []string
		start  = nums[0]
		end    = nums[0]
	)

	for i := 1; i < len(nums); i++ {
		if nums[i] == end+1 {
			end = nums[i]
		} else {
			result = append(result, formatRange(start, end))
			start = nums[i]
			end = nums[i]
		}
	}
	result = append(result, formatRange(start, end))

	return strings.Join(result, ",")
}

func formatRange[T ~int | ~int8 | ~int16 | ~int32 | ~int64](start, end T) string {
	if start == end {
		return strconv.FormatInt(int64(start), 10)
	}
	return strconv.FormatInt(int64(start), 10) + "-" + strconv.FormatInt(int64(end), 10)
}
