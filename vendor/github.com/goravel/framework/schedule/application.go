package schedule

import (
	"context"
	"strings"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/goravel/framework/contracts/cache"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/log"
	"github.com/goravel/framework/contracts/schedule"
	"github.com/goravel/framework/support/carbon"
)

type Application struct {
	artisan console.Artisan
	cache   cache.Cache
	log     log.Log
	cron    *cron.Cron
	events  []schedule.Event
	debug   bool
}

func NewApplication(artisan console.Artisan, cache cache.Cache, log log.Log, debug bool) *Application {
	return &Application{
		artisan: artisan,
		cache:   cache,
		cron: cron.New(cron.WithParser(cron.NewParser(
			cron.SecondOptional|cron.Minute|cron.Hour|cron.Dom|cron.Month|cron.Dow|cron.Descriptor,
		)), cron.WithLogger(NewLogger(log, debug)), cron.WithLocation(carbon.Now().StdTime().Location())),
		log:   log,
		debug: debug,
	}
}

func (app *Application) Call(callback func()) schedule.Event {
	return NewCallbackEvent(callback)
}

func (app *Application) Command(command string) schedule.Event {
	return NewCommandEvent(command)
}

func (app *Application) Events() []schedule.Event {
	return app.events
}

func (app *Application) Register(events []schedule.Event) {
	app.addEvents(events)
}

func (app *Application) Run() {
	app.cron.Run()
}

func (app *Application) Shutdown(ctx ...context.Context) error {
	if len(ctx) == 0 {
		ctx = append(ctx, context.Background())
	}

	cronCtx := app.cron.Stop()

	for {
		select {
		case <-cronCtx.Done():
			return nil
		case <-ctx[0].Done():
			return ctx[0].Err()
		}
	}
}

func (app *Application) addEvents(events []schedule.Event) {
	for _, event := range events {
		chain := cron.NewChain(cron.Recover(NewLogger(app.log, app.debug)))
		if event.GetDelayIfStillRunning() {
			chain = cron.NewChain(cron.DelayIfStillRunning(NewLogger(app.log, app.debug)), cron.Recover(NewLogger(app.log, app.debug)))
		} else if event.GetSkipIfStillRunning() {
			chain = cron.NewChain(cron.SkipIfStillRunning(NewLogger(app.log, app.debug)), cron.Recover(NewLogger(app.log, app.debug)))
		}
		_, err := app.cron.AddJob(event.GetCron(), chain.Then(app.getJob(event)))

		if err != nil {
			app.log.Errorf("add schedule error: %v", err)
			continue
		}
		app.events = append(app.events, event)
	}
}

func (app *Application) getJob(event schedule.Event) cron.Job {
	return cron.FuncJob(func() {
		if event.IsOnOneServer() && event.GetName() != "" {
			keySuffix := carbon.Now().Format("Hi")
			if segments := strings.Split(event.GetCron(), " "); len(segments) == 6 {
				keySuffix = carbon.Now().Format("His")
			}
			if app.cache.Lock(event.GetName()+keySuffix, 1*time.Hour).Get() {
				app.runJob(event)
			}
		} else {
			app.runJob(event)
		}
	})
}

func (app *Application) runJob(event schedule.Event) {
	if event.GetCommand() != "" {
		if err := app.artisan.Call(event.GetCommand()); err != nil {
			app.log.Errorf("run %s command error: %v", event.GetCommand(), err)
		}
	} else {
		event.GetCallback()()
	}
}
