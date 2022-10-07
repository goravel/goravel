package commands

import (
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/facades"
)

type Test struct {
}

//Signature The name and signature of the console command.
func (receiver *Test) Signature() string {
	return "command"
}

//Description The console command description.
func (receiver *Test) Description() string {
	return "test command"
}

//Extend The console command extend.
func (receiver *Test) Extend() command.Extend {
	return command.Extend{
		Flags: []command.Flag{
			{
				Name:    "name",
				Value:   "World",
				Aliases: []string{"n"},
				Usage:   "Name",
			},
			{
				Name:    "age",
				Value:   "18",
				Aliases: []string{"a"},
				Usage:   "Age",
			},
		},
	}
}

//Handle Execute the console command.
func (receiver *Test) Handle(ctx console.Context) error {
	fmt.Printf("Run test command success, argument_0: %s, argument_1: %s, option_name: %s, option_age: %s, arguments: %s\n",
		ctx.Argument(0), ctx.Argument(1), ctx.Option("name"), ctx.Option("age"), strings.Join(ctx.Arguments(), ","))
	facades.Log.Info("test")

	return nil
}
