package console

import (
	"context"
	"io"
	"os"
	"slices"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/errors"
	"github.com/goravel/framework/support/color"
	"github.com/goravel/framework/support/env"
)

var (
	noANSI     bool
	noANSIFlag = &cli.BoolFlag{
		Name:        "no-ansi",
		Destination: &noANSI,
		HideDefault: true,
		Usage:       "Force disable ANSI output",
	}
)

type Application struct {
	commands   []console.Command
	name       string
	usage      string
	usageText  string
	useArtisan bool
	version    string
	writer     io.Writer
}

// NewApplication Create a new Artisan application.
// Will add artisan flag to the command if useArtisan is true.
func NewApplication(name, usage, usageText, version string, useArtisan bool) *Application {
	return &Application{
		name:       name,
		usage:      usage,
		usageText:  usageText,
		useArtisan: useArtisan,
		version:    version,
		writer:     os.Stdout,
	}
}

// Call Run an Artisan console command by name.
func (r *Application) Call(command string) error {
	if len(os.Args) == 0 {
		return nil
	}

	commands := []string{os.Args[0]}

	if r.useArtisan {
		commands = append(commands, "artisan")
	}

	return r.Run(append(commands, strings.Split(command, " ")...), false)
}

// CallAndExit Run an Artisan console command by name and exit.
func (r *Application) CallAndExit(command string) {
	if len(os.Args) == 0 {
		return
	}

	commands := []string{os.Args[0]}

	if r.useArtisan {
		commands = append(commands, "artisan")
	}

	_ = r.Run(append(commands, strings.Split(command, " ")...), true)
}

// Register commands to the application.
func (r *Application) Register(commands []console.Command) {
	r.commands = append(r.commands, commands...)
}

// Run a command. Args come from os.Args.
func (r *Application) Run(args []string, exitIfArtisan bool) error {
	if noANSI || env.IsNoANSI() || slices.Contains(args, "--no-ansi") {
		color.Disable()
	} else {
		color.Enable()
	}

	artisanIndex := -1
	if r.useArtisan {
		for i, arg := range args {
			if arg == "artisan" {
				artisanIndex = i
				break
			}
		}
	} else {
		artisanIndex = 0
	}

	if artisanIndex != -1 {
		command, err := r.command()
		if err != nil {
			return err
		}

		if artisanIndex+1 == len(args) {
			args = append(args, "list")
		}

		cliArgs := append([]string{args[0]}, args[artisanIndex+1:]...)
		if err := command.Run(context.Background(), cliArgs); err != nil {
			if exitIfArtisan {
				panic(err.Error())
			}

			return err
		}

		if exitIfArtisan {
			os.Exit(0)
		}
	}

	return nil
}

// SetCommands Set the commands for the application.
func (r *Application) SetCommands(commands []console.Command) {
	r.commands = commands
}

func (r *Application) command() (*cli.Command, error) {
	cliCommands, err := commandsToCliCommands(r.commands)
	if err != nil {
		return nil, err
	}

	command := &cli.Command{}
	command.CommandNotFound = commandNotFound
	command.Commands = cliCommands
	command.Flags = []cli.Flag{noANSIFlag}
	command.Name = r.name
	command.OnUsageError = onUsageError
	command.Usage = r.usage
	command.UsageText = r.usageText
	command.Version = r.version
	command.Writer = r.writer

	// There is a concurrency issue with urfave/cli v3 when help is not hidden.
	command.HideHelp = true

	return command, nil
}

func commandsToCliCommands(commands []console.Command) ([]*cli.Command, error) {
	cliCommands := make([]*cli.Command, len(commands))

	for i, item := range commands {
		arguments := item.Extend().Arguments
		cliArguments, err := argumentsToCliArgs(arguments)
		if err != nil {
			return nil, errors.ConsoleCommandRegisterFailed.Args(item.Signature(), err)
		}
		cliCommands[i] = &cli.Command{
			Name:  item.Signature(),
			Usage: item.Description(),
			Action: func(ctx context.Context, cmd *cli.Command) error {
				cliCtx := NewCliContext(cmd, arguments)
				if cliCtx.OptionBool("help") {
					return cli.ShowCommandHelp(ctx, cmd, cmd.Name)
				}

				return item.Handle(cliCtx)
			},
			Category:     item.Extend().Category,
			ArgsUsage:    item.Extend().ArgsUsage,
			Flags:        flagsToCliFlags(item.Extend().Flags),
			Arguments:    cliArguments,
			OnUsageError: onUsageError,
		}
	}

	return cliCommands, nil
}

func flagsToCliFlags(flags []command.Flag) []cli.Flag {
	var cliFlags []cli.Flag
	for _, flag := range flags {
		switch flag.Type() {
		case command.FlagTypeBool:
			flag := flag.(*command.BoolFlag)
			cliFlags = append(cliFlags, &cli.BoolFlag{
				Name:        flag.Name,
				Aliases:     flag.Aliases,
				HideDefault: flag.DisableDefaultText,
				Usage:       flag.Usage,
				Required:    flag.Required,
				Value:       flag.Value,
			})
		case command.FlagTypeFloat64:
			flag := flag.(*command.Float64Flag)
			cliFlags = append(cliFlags, &cli.FloatFlag{
				Name:     flag.Name,
				Aliases:  flag.Aliases,
				Usage:    flag.Usage,
				Required: flag.Required,
				Value:    flag.Value,
			})
		case command.FlagTypeFloat64Slice:
			flag := flag.(*command.Float64SliceFlag)
			cliFlags = append(cliFlags, &cli.FloatSliceFlag{
				Name:     flag.Name,
				Aliases:  flag.Aliases,
				Usage:    flag.Usage,
				Required: flag.Required,
				Value:    cli.NewFloatSlice(flag.Value...).Value(),
			})
		case command.FlagTypeInt:
			flag := flag.(*command.IntFlag)
			cliFlags = append(cliFlags, &cli.IntFlag{
				Name:     flag.Name,
				Aliases:  flag.Aliases,
				Usage:    flag.Usage,
				Required: flag.Required,
				Value:    flag.Value,
			})
		case command.FlagTypeIntSlice:
			flag := flag.(*command.IntSliceFlag)
			cliFlags = append(cliFlags, &cli.IntSliceFlag{
				Name:     flag.Name,
				Aliases:  flag.Aliases,
				Usage:    flag.Usage,
				Required: flag.Required,
				Value:    flag.Value,
			})
		case command.FlagTypeInt64:
			flag := flag.(*command.Int64Flag)
			cliFlags = append(cliFlags, &cli.Int64Flag{
				Name:     flag.Name,
				Aliases:  flag.Aliases,
				Usage:    flag.Usage,
				Required: flag.Required,
				Value:    flag.Value,
			})
		case command.FlagTypeInt64Slice:
			flag := flag.(*command.Int64SliceFlag)
			cliFlags = append(cliFlags, &cli.Int64SliceFlag{
				Name:     flag.Name,
				Aliases:  flag.Aliases,
				Usage:    flag.Usage,
				Required: flag.Required,
				Value:    flag.Value,
			})
		case command.FlagTypeString:
			flag := flag.(*command.StringFlag)
			cliFlags = append(cliFlags, &cli.StringFlag{
				Name:     flag.Name,
				Aliases:  flag.Aliases,
				Usage:    flag.Usage,
				Required: flag.Required,
				Value:    flag.Value,
			})
		case command.FlagTypeStringSlice:
			flag := flag.(*command.StringSliceFlag)
			cliFlags = append(cliFlags, &cli.StringSliceFlag{
				Name:     flag.Name,
				Aliases:  flag.Aliases,
				Usage:    flag.Usage,
				Required: flag.Required,
				Value:    flag.Value,
			})
		}
	}

	var (
		existHelp bool
		existH    bool
	)
	for _, flag := range cliFlags {
		names := flag.Names()
		if slices.Contains(names, "help") {
			existHelp = true
		}
		if slices.Contains(names, "h") {
			existH = true
		}
	}

	if !existHelp {
		helpFlag := &cli.BoolFlag{
			Name:        "help",
			Usage:       "Show help",
			HideDefault: true,
		}
		if !existH {
			helpFlag.Aliases = []string{"h"}
		}
		cliFlags = append(cliFlags, helpFlag)
	}

	cliFlags = append(cliFlags, noANSIFlag)

	return cliFlags
}

func argumentsToCliArgs(args []command.Argument) ([]cli.Argument, error) {
	len := len(args)
	if len == 0 {
		return nil, nil
	}
	cliArgs := make([]cli.Argument, 0, len)
	previousIsRequired := true
	for _, v := range args {
		if v.GetMin() != 0 && !previousIsRequired {
			return nil, errors.ConsoleCommandRequiredArgumentWrongOrder.Args(v.GetName())
		}
		if v.GetMin() != 0 {
			previousIsRequired = true
		} else {
			previousIsRequired = false
		}
		switch arg := v.(type) {
		case *command.ArgumentFloat32:
			cliArgs = append(cliArgs, &cli.Float32Args{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentFloat64:
			cliArgs = append(cliArgs, &cli.Float64Args{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentInt:
			cliArgs = append(cliArgs, &cli.IntArgs{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentInt8:
			cliArgs = append(cliArgs, &cli.Int8Args{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentInt16:
			cliArgs = append(cliArgs, &cli.Int16Args{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentInt32:
			cliArgs = append(cliArgs, &cli.Int32Args{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentInt64:
			cliArgs = append(cliArgs, &cli.Int64Args{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentString:
			cliArgs = append(cliArgs, &cli.StringArgs{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentTimestamp:
			cliArgs = append(cliArgs, &cli.TimestampArgs{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
				Config: cli.TimestampConfig{
					Layouts: arg.Layouts,
				},
			})
		case *command.ArgumentUint:
			cliArgs = append(cliArgs, &cli.UintArgs{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentUint8:
			cliArgs = append(cliArgs, &cli.Uint8Args{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentUint16:
			cliArgs = append(cliArgs, &cli.Uint16Args{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentUint32:
			cliArgs = append(cliArgs, &cli.Uint32Args{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentUint64:
			cliArgs = append(cliArgs, &cli.Uint64Args{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})

		case *command.ArgumentFloat32Slice:
			cliArgs = append(cliArgs, &cli.Float32Args{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentFloat64Slice:
			cliArgs = append(cliArgs, &cli.Float64Args{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentIntSlice:
			cliArgs = append(cliArgs, &cli.IntArgs{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentInt8Slice:
			cliArgs = append(cliArgs, &cli.Int8Args{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentInt16Slice:
			cliArgs = append(cliArgs, &cli.Int16Args{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentInt32Slice:
			cliArgs = append(cliArgs, &cli.Int32Args{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentInt64Slice:
			cliArgs = append(cliArgs, &cli.Int64Args{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentStringSlice:
			cliArgs = append(cliArgs, &cli.StringArgs{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentTimestampSlice:
			cliArgs = append(cliArgs, &cli.TimestampArgs{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
				Config: cli.TimestampConfig{
					Layouts: arg.Layouts,
				},
			})
		case *command.ArgumentUintSlice:
			cliArgs = append(cliArgs, &cli.UintArgs{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentUint8Slice:
			cliArgs = append(cliArgs, &cli.Uint8Args{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentUint16Slice:
			cliArgs = append(cliArgs, &cli.Uint16Args{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentUint32Slice:
			cliArgs = append(cliArgs, &cli.Uint32Args{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		case *command.ArgumentUint64Slice:
			cliArgs = append(cliArgs, &cli.Uint64Args{
				Name:      arg.Name,
				UsageText: arg.Usage,
				Min:       arg.GetMin(),
				Max:       arg.GetMax(),
			})
		default:
			return nil, errors.ConsoleCommandArgumentUnknownType.Args(arg, arg)
		}
	}
	return cliArgs, nil
}
