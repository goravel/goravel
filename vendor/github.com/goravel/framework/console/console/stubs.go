package console

type Stubs struct {
}

func (r Stubs) Command() string {
	return `package DummyPackage

import (
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
)

type DummyCommand struct {
}

// Signature The name and signature of the console command.
func (r *DummyCommand) Signature() string {
	return "DummySignature"
}

// Description The console command description.
func (r *DummyCommand) Description() string {
	return "Command description"
}

// Extend The console command extend.
func (r *DummyCommand) Extend() command.Extend {
	return command.Extend{Category: "app"}
}

// Handle Execute the console command.
func (r *DummyCommand) Handle(ctx console.Context) error {
	
	return nil
}
`
}
