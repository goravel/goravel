package console

import (
	"strings"

	"github.com/goravel/framework/support/env"
	"github.com/goravel/framework/support/str"
)

type Stubs struct {
}

func (r Stubs) Test() string {
	return `package DummyPackage

import (
	"testing"

	"github.com/stretchr/testify/suite"
	
	DummyTestImport
)

type DummyTestSuite struct {
	suite.Suite
	DummyTestCase
}

func TestDummyTestSuite(t *testing.T) {
	suite.Run(t, new(DummyTestSuite))
}

// SetupTest will run before each test in the suite.
func (s *DummyTestSuite) SetupTest() {
}

// TearDownTest will run after each test in the suite.
func (s *DummyTestSuite) TearDownTest() {
}

func (s *DummyTestSuite) TestIndex() {
	// TODO
}
`
}

func (r Stubs) ServiceProvider() string {
	return `package DummyPackage

import (
	"github.com/goravel/framework/contracts/binding"
	"github.com/goravel/framework/contracts/foundation"
)

type DummyServiceProvider struct{}

// Relationship provides the service provider's bindings, their dependencies, and the services they provide for.
// It's optional if the service provider doesn't depend on or provide any other services.
func (r *DummyServiceProvider) Relationship() binding.Relationship {
	return binding.Relationship{
		Bindings:     []string{},
		Dependencies: []string{},
		ProvideFor:   []string{},
	}
}
// Register service bindings here
func (r *DummyServiceProvider) Register(app foundation.Application) {
	// Example:
	// app.Singleton("example", func(app foundation.Application) (any, error) {
	//     return &ExampleService{}, nil
	// })
}

// Boot performs post-registration booting of services.
// It will be called after all service providers have been registered.
func (r *DummyServiceProvider) Boot(app foundation.Application) {
}
`
}

type PackageMakeCommandStubs struct {
	main string
	pkg  string
	root string
	name string
}

func NewPackageMakeCommandStubs(pkg, root string) *PackageMakeCommandStubs {
	return &PackageMakeCommandStubs{main: env.MainPath(), pkg: pkg, root: root, name: packageName(pkg)}
}

func (r PackageMakeCommandStubs) Readme() string {
	content := `# DummyName
`

	return strings.ReplaceAll(content, "DummyName", r.name)
}

func (r PackageMakeCommandStubs) ServiceProvider() string {
	content := `package DummyName

import (
	"github.com/goravel/framework/contracts/binding"
	"github.com/goravel/framework/contracts/foundation"
)

const Binding = "DummyPackage"

var App foundation.Application

type ServiceProvider struct {
}

// Relationship returns the relationship of the service provider.
func (r *ServiceProvider) Relationship() binding.Relationship {
	return binding.Relationship{
		Bindings: []string{},
		Dependencies: []string{},
		ProvideFor: []string{},
	}
}

// Register registers the service provider.
func (r *ServiceProvider) Register(app foundation.Application) {
	App = app

	app.Bind(Binding, func(app foundation.Application) (any, error) {
		return &DummyCamelName{}, nil
	})
}

// Boot boots the service provider, will be called after all service providers are registered.
func (r *ServiceProvider) Boot(app foundation.Application) {

}
`

	content = strings.ReplaceAll(content, "DummyPackage", r.pkg)
	content = strings.ReplaceAll(content, "DummyName", r.name)
	content = strings.ReplaceAll(content, "DummyCamelName", str.Of(r.name).Studly().String())

	return content
}

func (r PackageMakeCommandStubs) Main() string {
	content := `package DummyName

type DummyCamelName struct {}
`

	content = strings.ReplaceAll(content, "DummyName", r.name)
	content = strings.ReplaceAll(content, "DummyCamelName", str.Of(r.name).Studly().String())

	return content
}

func (r PackageMakeCommandStubs) Config() string {
	content := `package DummyPackage

import (
	"DummyFacadesImport"
)

func init() {
	config := DummyFacadesPackage.Config()
	config.Add("DummyName", map[string]any{
		
	})
}
`

	file := `package main

import "strings"

func config(configPackage string, facadesImport, facadesPackage string) string {
	content := DummyContent
	content = strings.ReplaceAll(content, "DummyPackage", configPackage)
	content = strings.ReplaceAll(content, "DummyFacadesImport", facadesImport)
	content = strings.ReplaceAll(content, "DummyFacadesPackage", facadesPackage)

	return content
}
`

	content = strings.ReplaceAll(content, "DummyName", r.name)
	content = strings.ReplaceAll(file, "DummyContent", "`"+content+"`")

	return content
}

func (r PackageMakeCommandStubs) OldConfig() string {
	content := `package config

import (
	"github.com/goravel/framework/facades"
)

func init() {
	config := facades.Config()
	config.Add("DummyName", map[string]any{
		
	})
}
`

	return strings.ReplaceAll(content, "DummyName", r.name)
}

func (r PackageMakeCommandStubs) Contracts() string {
	content := `package contracts

type DummyCamelName interface {}
`

	return strings.ReplaceAll(content, "DummyCamelName", str.Of(r.name).Studly().String())
}

func (r PackageMakeCommandStubs) Facades() string {
	content := `package facades

import (
	"log"

	"DummyMain/DummyRoot"
	"DummyMain/DummyRoot/contracts"
)

func DummyCamelName() contracts.DummyCamelName {
	instance, err := DummyName.App.Make(DummyName.Binding)
	if err != nil {
		log.Println(err)
		return nil
	}

	return instance.(contracts.DummyCamelName)
}
`

	content = strings.ReplaceAll(content, "DummyMain", r.main)
	content = strings.ReplaceAll(content, "DummyRoot", r.root)
	content = strings.ReplaceAll(content, "DummyName", r.name)
	content = strings.ReplaceAll(content, "DummyCamelName", str.Of(r.name).Studly().String())

	return content
}

func (r PackageMakeCommandStubs) Setup() string {
	content := `
package main

import (
	"os"

	"github.com/goravel/framework/packages"
	"github.com/goravel/framework/packages/modify"
	"github.com/goravel/framework/support/path"
)

func main() {
	setup := packages.Setup(os.Args)
	serviceProvider := "&DummyName.ServiceProvider{}"
	moduleImport := setup.Paths().Module().Import()
	configPath := path.Config("DummyName.go")

	setup.Install(
		// Register the service provider
		modify.RegisterProvider(moduleImport, serviceProvider),

		// Add config
		modify.File(configPath).Overwrite(config(setup.Paths().Config().Package(), setup.Paths().Facades().Import(), setup.Paths().Facades().Package())),
	).Uninstall(
		// Remove config/cache.go
		modify.File(configPath).Remove(),

		// Unregister the service provider
		modify.UnregisterProvider(moduleImport, serviceProvider),
	).Execute()
}
`
	content = strings.ReplaceAll(content, "DummyName", r.name)

	return content
}

func (r PackageMakeCommandStubs) OldSetup() string {
	content := `package main

import (
	"os"

	"github.com/goravel/framework/packages"
	"github.com/goravel/framework/packages/match"
	"github.com/goravel/framework/packages/modify"
	"github.com/goravel/framework/support/path"
)

func main() {
	setup := packages.Setup(os.Args)
	setup.Install(
		modify.GoFile(path.Config("app.go")).
			Find(match.Imports()).Modify(modify.AddImport(setup.Paths().Module().Import())).
			Find(match.ProvidersInConfig()).Modify(modify.Register("&DummyName.ServiceProvider{}")),
	).Uninstall(
		modify.GoFile(path.Config("app.go")).
			Find(match.ProvidersInConfig()).Modify(modify.Unregister("&DummyName.ServiceProvider{}")).
			Find(match.Imports()).Modify(modify.RemoveImport(setup.Paths().Module().Import())),
	).Execute()
}
`
	content = strings.ReplaceAll(content, "DummyName", r.name)

	return content
}
