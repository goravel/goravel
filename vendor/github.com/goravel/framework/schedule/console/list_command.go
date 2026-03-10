package console

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/robfig/cron/v3"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/schedule"
	"github.com/goravel/framework/support/carbon"
	"github.com/goravel/framework/support/debug"
	"github.com/goravel/framework/support/str"
)

var cronParser = cron.NewParser(cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)

type List struct {
	schedule schedule.Schedule
}

func NewList(schedule schedule.Schedule) *List {
	return &List{
		schedule: schedule,
	}
}

// Signature The name and signature of the console command.
func (r *List) Signature() string {
	return "schedule:list"
}

// Description The console command description.
func (r *List) Description() string {
	return "List all scheduled tasks"
}

// Extend The console command extend.
func (r *List) Extend() command.Extend {
	return command.Extend{
		Category: "schedule",
	}
}

// Handle Execute the console command.
func (r *List) Handle(ctx console.Context) error {
	ctx.NewLine()
	events := r.schedule.Events()
	if len(events) == 0 {
		ctx.Warning("No scheduled tasks have been defined.")
		return nil
	}

	spacing := r.getCronExpressionSpacing(events)
	for _, event := range events {
		expression := event.GetCron()
		ctx.TwoColumnDetail(
			fmt.Sprintf("<fg=yellow>%s</>  %s", r.formatCronExpression(expression, spacing), r.getCommand(event)),
			fmt.Sprintf("<fg=7472a3>Next Due: %s</>", r.getNextDueDate(expression)),
		)
	}
	return nil
}

func (r *List) formatCronExpression(expression string, spacing []int) string {
	parts := strings.Fields(expression)
	padded := make([]string, len(spacing))

	// if parts length is less than spacing, prepend empty strings to ensure alignment
	parts = append(make([]string, len(spacing)-len(parts)), parts...)

	for i := range padded {
		val := parts[i]
		padded[i] = val + strings.Repeat(" ", max(spacing[i]-utf8.RuneCountInString(val), 0))
	}

	return strings.Join(padded, " ")
}

func (r *List) getCommand(event schedule.Event) string {
	// display artisan command signature
	if c := event.GetCommand(); c != "" {
		// highlight the parameters...
		return regexp.MustCompile(`(artisan [\w\-:]+) (.+)`).ReplaceAllString("artisan "+c, `$1 <fg=yellow>$2</>`)
	}

	// display name
	if name := event.GetName(); len(name) > 0 {
		return name
	}

	// display closure location
	file, line := r.getClosureLocation(event.GetCallback())
	file, _ = filepath.Rel(str.Of(file).Dirname(3).String(), file)

	return fmt.Sprintf("Closure at: %s:%d", file, line)
}

func (r *List) getClosureLocation(closure any) (file string, line int) {
	info := debug.GetFuncInfo(closure)

	return info.File, info.Line
}

func (r *List) getCronExpressionSpacing(events []schedule.Event) []int {
	// supports both six-field (second-level) and five-field (minute-level) cron expressions.
	spacing := make([]int, 6)

	for _, event := range events {
		parts := strings.Fields(event.GetCron())
		offset := len(spacing) - len(parts)
		for i, part := range parts {
			spacing[offset+i] = max(spacing[offset+i], utf8.RuneCountInString(part))
		}
	}

	// if the first field (seconds) is not used
	// return spacing for five fields (minute-level cron)
	if spacing[0] == 0 {
		return spacing[1:]
	}

	return spacing
}

func (r *List) getNextDueDate(cronSpec string) string {
	if sch, err := cronParser.Parse(cronSpec); err == nil {
		now := carbon.Now()
		if next := sch.Next(now.StdTime()); !next.IsZero() {
			return carbon.FromStdTime(next).DiffForHumans(now)
		}
	}

	return ""
}
