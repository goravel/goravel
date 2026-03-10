package console

import (
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/packages"
	"github.com/goravel/framework/support"
	supportconsole "github.com/goravel/framework/support/console"
	"github.com/goravel/framework/support/file"
	"github.com/goravel/framework/support/str"
)

type TestMakeCommand struct {
}

func NewTestMakeCommand() *TestMakeCommand {
	return &TestMakeCommand{}
}

// Signature The name and signature of the console command.
func (r *TestMakeCommand) Signature() string {
	return "make:test"
}

// Description The console command description.
func (r *TestMakeCommand) Description() string {
	return "Create a new test class"
}

// Extend The console command extend.
func (r *TestMakeCommand) Extend() command.Extend {
	return command.Extend{
		Category: "make",
		Flags: []command.Flag{
			&command.BoolFlag{
				Name:    "force",
				Aliases: []string{"f"},
				Usage:   "Create the test even if it already exists",
			},
		},
	}
}

// Handle Execute the console command.
func (r *TestMakeCommand) Handle(ctx console.Context) error {
	filePath := ctx.Argument(0)
	m, err := supportconsole.NewMake(ctx, "test", filePath, support.Config.Paths.Tests)
	if err != nil {
		ctx.Error(err.Error())
		return nil
	}

	stub := r.getStub()

	var testCase, testsImport string
	if str.Of(filePath).Contains("/", "\\") {
		testCase = fmt.Sprintf("%s.TestCase", packages.Paths().Tests().Package())
		testsImport = fmt.Sprintf(`"%s"`, packages.Paths().Tests().Import())
	} else {
		testCase = "TestCase"
	}

	if err := file.PutContent(m.GetFilePath(), r.populateStub(stub, m.GetPackageName(), m.GetStructName(), testsImport, testCase)); err != nil {
		ctx.Error(err.Error())
		return nil
	}

	ctx.Success("Test created successfully")

	return nil
}

func (r *TestMakeCommand) getStub() string {
	return Stubs{}.Test()
}

// populateStub Populate the place-holders in the command stub.
func (r *TestMakeCommand) populateStub(stub, packageName, structName, testsImport, testCase string) string {
	stub = strings.ReplaceAll(stub, "DummyTestSuite", structName+"Suite")
	stub = strings.ReplaceAll(stub, "DummyPackage", packageName)
	stub = strings.ReplaceAll(stub, "DummyTestImport", testsImport)
	stub = strings.ReplaceAll(stub, "DummyTestCase", testCase)

	return stub
}
