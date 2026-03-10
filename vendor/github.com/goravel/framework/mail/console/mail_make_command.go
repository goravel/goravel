package console

import (
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/support"
	supportconsole "github.com/goravel/framework/support/console"
	"github.com/goravel/framework/support/file"
)

type MailMakeCommand struct {
}

func NewMailMakeCommand() *MailMakeCommand {
	return &MailMakeCommand{}
}

// Signature The name and signature of the console command.
func (r *MailMakeCommand) Signature() string {
	return "make:mail"
}

// Description The console command description.
func (r *MailMakeCommand) Description() string {
	return "Create a new mail class"
}

// Extend The console command extend.
func (r *MailMakeCommand) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Create the mail even if it already exists",
			},
		},
	}
}

// Handle Execute the console command.
func (r *MailMakeCommand) Handle(ctx console.Context) error {
	m, err := supportconsole.NewMake(ctx, "mail", ctx.Argument(0), support.Config.Paths.Mails)
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	if err := file.PutContent(m.GetFilePath(), r.populateStub(r.getStub(), m.GetPackageName(), m.GetStructName())); err != nil {
		ctx.Error(err.Error())
		return nil
	}

	ctx.Success("Mail created successfully")

	return nil
}

func (r *MailMakeCommand) getStub() string {
	return Stubs{}.Mail()
}

// populateStub Populate the place-holders in the command stub.
func (r *MailMakeCommand) populateStub(stub string, packageName, structName string) string {
	stub = strings.ReplaceAll(stub, "DummyMail", structName)
	stub = strings.ReplaceAll(stub, "DummyPackage", packageName)

	return stub
}
