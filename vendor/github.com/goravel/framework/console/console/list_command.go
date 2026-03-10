package console

import (
	"github.com/urfave/cli/v3"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

type ListCommand struct{}

func NewListCommand() *ListCommand {
	return &ListCommand{}
}

// Signature The name and signature of the console command.
func (r *ListCommand) Signature() string {
	return "list"
}

// Description The console command description.
func (r *ListCommand) Description() string {
	return "List commands"
}

// Extend The console command extend.
func (r *ListCommand) Extend() command.Extend {
	return command.Extend{}
}

// Handle Execute the console command.
func (r *ListCommand) Handle(ctx console.Context) error {
	return cli.ShowAppHelp(ctx.Instance())
}
