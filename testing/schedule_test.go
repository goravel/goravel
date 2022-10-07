package testing

import (
	"context"
	"fmt"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/file"
	"github.com/stretchr/testify/assert"
	"goravel/bootstrap"
	"io/ioutil"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestSchedule(t *testing.T) {
	kernel := HttpKernelSchedule{}
	kernel.Create()

	bootstrap.Boot()

	second, _ := strconv.Atoi(time.Now().Format("05"))
	// Make sure run 3 times
	ctx, _ := context.WithTimeout(context.Background(), time.Duration(120+6+60-second)*time.Second)
	go func(ctx context.Context) {
		facades.Schedule.Run()

		for {
			select {
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	time.Sleep(time.Duration(120+5+60-second) * time.Second)
	log := fmt.Sprintf("storage/logs/goravel-%s.log", time.Now().Format("2006-01-02"))
	assert.True(t, file.Exist(log))
	data, err := ioutil.ReadFile(log)
	assert.Nil(t, err)
	assert.Equal(t, 3, strings.Count(string(data), "schedule closure immediately"))
	assert.Equal(t, 3, strings.Count(string(data), "Run test command success, argument_0: argument0, argument_1: argument1, option_name: Goravel, option_age: 18, arguments: argument0,argument1"))
	assert.Equal(t, 2, strings.Count(string(data), "schedule closure DelayIfStillRunning"))
	assert.Equal(t, 1, strings.Count(string(data), "schedule closure SkipIfStillRunning"))
	//assert.True(t, file.Remove("./storage"))
}

type HttpKernelSchedule struct {
}

func (r *HttpKernelSchedule) stub() string {
	return `package console

import (
    "time"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/schedule"
	"github.com/goravel/framework/facades"

	"goravel/app/console/commands"
)

type Kernel struct {
}

func (kernel *Kernel) Schedule() []schedule.Event {
	return []schedule.Event{
		facades.Schedule.Call(func() {
			facades.Log.Info("schedule closure immediately")
		}).EveryMinute(),
		facades.Schedule.Call(func() {
			time.Sleep(61 * time.Second)
			facades.Log.Info("schedule closure DelayIfStillRunning")
		}).EveryMinute().DelayIfStillRunning(),
		facades.Schedule.Call(func() {
			time.Sleep(61 * time.Second)
			facades.Log.Info("schedule closure SkipIfStillRunning")
		}).EveryMinute().SkipIfStillRunning(),
		facades.Schedule.Command("test --name Goravel argument0 argument1").EveryMinute(),
	}
}

func (kernel *Kernel) Commands() []console.Command {
	return []console.Command{
		&commands.Test{},
	}
}
`
}

func (r *HttpKernelSchedule) Create() {
	path := "../app/console/kernel.go"
	file.Remove(path)
	file.Create(path, r.stub())
}
