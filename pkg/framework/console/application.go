package console

import (
	"github.com/goravel/framework/console/support"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

var cliInstance *cli.App

func init() {
	cliInstance = cli.NewApp()
}

type Application struct {
}

func (app *Application) Init() {
	args := os.Args

	if len(args) > 2 {
		if args[1] == "artisan" {
			cliApp := app.Instance()
			var cliArgs []string
			cliArgs = append(cliArgs, args[0])

			for i := 2; i < len(args); i++ {
				cliArgs = append(cliArgs, args[i])
			}

			if err := cliApp.Run(cliArgs); err != nil {
				log.Fatalln(err.Error())
			}

			os.Exit(0)
		}
	}
}

func (app *Application) Instance() *cli.App {
	return cliInstance
}

func (app *Application) Register(commands []support.Command) {
	for _, command := range commands {
		command := command
		cliCommand := cli.Command{
			Name:  command.Signature(),
			Usage: command.Description(),
			Action: func(c *cli.Context) error {
				return command.Handle(c)
			},
		}

		if len(command.Flags()) > 0 {
			cliCommand.Flags = command.Flags()
		}

		if len(command.Subcommands()) > 0 {
			cliCommand.Subcommands = command.Subcommands()
		}

		cliInstance.Commands = append(cliInstance.Commands, &cliCommand)
	}
}
