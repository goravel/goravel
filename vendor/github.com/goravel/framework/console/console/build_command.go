package console

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/process"
)

type BuildCommand struct {
	config  config.Config
	process process.Process
}

func NewBuildCommand(config config.Config, process process.Process) *BuildCommand {
	return &BuildCommand{
		config:  config,
		process: process,
	}
}

// Signature The name and signature of the console command.
func (r *BuildCommand) Signature() string {
	return "build"
}

// Description The console command description.
func (r *BuildCommand) Description() string {
	return "Build the application"
}

// Extend The console command extend.
func (r *BuildCommand) Extend() command.Extend {
	return command.Extend{
		Flags: []command.Flag{
			&command.StringFlag{
				Name:    "arch",
				Aliases: []string{"a"},
				Usage:   "Target architecture",
				Value:   "amd64",
			},
			&command.StringFlag{
				Name:    "os",
				Aliases: []string{"o"},
				Usage:   "Target os",
			},
			&command.BoolFlag{
				Name:               "static",
				Aliases:            []string{"s"},
				Value:              false,
				Usage:              "Static compilation",
				DisableDefaultText: true,
			},
			&command.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Usage:   "Output binary name",
			},
		},
	}
}

// Handle Execute the console command.
func (r *BuildCommand) Handle(ctx console.Context) error {
	var err error
	if r.config.GetString("app.env") == "production" {
		ctx.Warning("**************************************")
		ctx.Warning("*     Application In Production!     *")
		ctx.Warning("**************************************")

		if !ctx.Confirm("Do you really wish to run this command?") {
			ctx.Warning("Command cancelled!")
			return nil
		}
	}

	os := ctx.Option("os")
	if os == "" {
		if os, err = ctx.Choice("Select target os", []console.Choice{
			{Key: "Linux", Value: "linux"},
			{Key: "Windows", Value: "windows"},
			{Key: "Darwin", Value: "darwin"},
		}, console.ChoiceOption{Default: runtime.GOOS}); err != nil {
			ctx.Error(fmt.Sprintf("Select target os error: %v", err))
			return nil
		}
	}

	if res := r.process.Env(map[string]string{
		"CGO_ENABLED": "0",
		"GOOS":        os,
		"GOARCH":      ctx.Option("arch"),
	}).WithSpinner("Building...").Run(generateCommand(ctx.Option("name"), ctx.OptionBool("static"))); res.Failed() {
		ctx.Error(res.Error().Error())
		return nil
	}

	ctx.Info("Built successfully.")

	return nil
}

func generateCommand(name string, static bool) string {
	args := []string{"go", "build"}

	if static {
		args = append(args, "-ldflags", `"-s -w -extldflags -static"`)
	}

	if name != "" {
		args = append(args, "-o", name)
	}

	args = append(args, ".")

	return strings.Join(args, " ")
}
