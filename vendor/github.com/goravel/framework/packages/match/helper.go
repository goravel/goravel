package match

import (
	"go/token"
	"strconv"
	"strings"

	"github.com/dave/dst"

	"github.com/goravel/framework/contracts/packages/match"
)

// Config matches configuration additions in init functions.
// It searches for config.Add() or facades.Config().Add() calls with the specified key path.
// The key parameter uses dot notation to represent nested configuration keys.
//
// Example usage:
//
//	match.GoFile("config/custom.go").
//	    Find(match.Config("database.connections.mysql")).
//	    Modify(modifyDatabaseConfig).
//	    Apply()
//
// This matches:
//
//	func init() {
//	    config.Add("database", map[string]any{
//	        "connections": map[string]any{
//	            "mysql": map[string]any{
//	                "host": "localhost",
//	            },
//	        },
//	    })
//	}
func Config(key string) []match.GoNode {
	keys := strings.Split(key, ".")
	matchers := []match.GoNode{
		Func(Ident("init")),
		CallExpr(
			SelectorExpr(
				AnyOf(
					Ident("config"),
					CallExpr(
						SelectorExpr(
							Ident("facades"),
							Ident("Config"),
						),
						AnyNodes(),
					),
				),
				Ident("Add"),
			),
			GoNodes{
				BasicLit(strconv.Quote(keys[0])),
				AnyNode(),
			},
		),
	}

	for _, k := range keys[1:] {
		matchers = append(matchers, KeyValueExpr(BasicLit(strconv.Quote(k)), AnyNode()))
	}

	return matchers
}

// Commands matches the Commands() function that returns a slice of console commands.
// It looks for a function returning []console.Command composite literals.
//
// Example usage:
//
//	match.GoFile("providers/custom_provider.go").
//	    Find(match.Commands()).
//	    Modify(addCommandToList(&CustomCommand{})).
//	    Apply()
//
// This matches:
//
//	func (receiver *ServiceProvider) Commands() []console.Command {
//	    return []console.Command{
//	        &commands.MigrateCommand{},
//	        &commands.SeedCommand{},
//	    }
//	}
func Commands() []match.GoNode {
	return []match.GoNode{
		Func(Ident("Commands")),
		TypeOf(&dst.ReturnStmt{}),
		CompositeLit(
			ArrayType(
				SelectorExpr(
					Ident("console"),
					Ident("Command"),
				),
				AnyNode(),
			),
		),
	}
}

// Filters matches the Filters() function that returns a slice of validation filters.
// It looks for a function returning []validation.Filter composite literals.
//
// Example usage:
//
//	match.GoFile("providers/validation_provider.go").
//	    Find(match.Filters()).
//	    Modify(addFilterToList(&TrimFilter{})).
//	    Apply()
//
// This matches:
//
//	func (receiver *ValidationServiceProvider) Filters() []validation.Filter {
//	    return []validation.Filter{
//	        &filters.TrimFilter{},
//	        &filters.SanitizeFilter{},
//	    }
//	}
func Filters() []match.GoNode {
	return []match.GoNode{
		Func(Ident("Filters")),
		TypeOf(&dst.ReturnStmt{}),
		CompositeLit(
			ArrayType(
				SelectorExpr(
					Ident("validation"),
					Ident("Filter"),
				),
				AnyNode(),
			),
		),
	}
}

// Imports matches import declaration blocks in Go source files.
// It identifies GenDecl nodes with IMPORT token type.
//
// Example usage:
//
//	match.GoFile("main.go").
//	    Find(match.Imports()).
//	    Modify(addImportToBlock("github.com/goravel/framework")).
//	    Apply()
//
// This matches:
//
//	import (
//	    "fmt"
//	    "strings"
//	    "github.com/goravel/framework/facades"
//	)
func Imports() []match.GoNode {
	return []match.GoNode{
		GoNode{
			match: func(n dst.Node) bool {
				if block, ok := n.(*dst.GenDecl); ok {
					return block.Tok == token.IMPORT
				}

				return false
			},
		},
	}
}

// Jobs matches the Jobs() function that returns a slice of queue jobs.
// It looks for a function returning []queue.Job composite literals.
//
// Example usage:
//
//	match.GoFile("providers/queue_provider.go").
//	    Find(match.Jobs()).
//	    Modify(addJobToList(&EmailJob{})).
//	    Apply()
//
// This matches:
//
//	func (receiver *QueueServiceProvider) Jobs() []queue.Job {
//	    return []queue.Job{
//	        &jobs.EmailJob{},
//	        &jobs.ProcessDataJob{},
//	    }
//	}
func Jobs() []match.GoNode {
	return []match.GoNode{
		Func(Ident("Jobs")),
		TypeOf(&dst.ReturnStmt{}),
		CompositeLit(
			ArrayType(
				SelectorExpr(
					Ident("queue"),
					Ident("Job"),
				),
				AnyNode(),
			),
		),
	}
}

// Migrations matches the Migrations() function that returns a slice of database migrations.
// It looks for a function returning []schema.Migration composite literals.
//
// Example usage:
//
//	match.GoFile("providers/migration_provider.go").
//	    Find(match.Migrations()).
//	    Modify(addMigrationToList(&CreateUsersTable{})).
//	    Apply()
//
// This matches:
//
//	func (receiver *MigrationServiceProvider) Migrations() []schema.Migration {
//	    return []schema.Migration{
//	        &migrations.CreateUsersTable{},
//	        &migrations.CreatePostsTable{},
//	    }
//	}
func Migrations() []match.GoNode {
	return []match.GoNode{
		Func(Ident("Migrations")),
		TypeOf(&dst.ReturnStmt{}),
		CompositeLit(
			ArrayType(
				SelectorExpr(
					Ident("schema"),
					Ident("Migration"),
				),
				AnyNode(),
			),
		),
	}
}

// Providers matches the Providers() function that returns a slice of service providers.
// It looks for a function returning []foundation.ServiceProvider composite literals.
// This is the recommended way to register service providers.
//
// Example usage:
//
//	match.GoFile("config/app.go").
//	    Find(match.Providers()).
//	    Modify(addProviderToList(&CustomServiceProvider{})).
//	    Apply()
//
// This matches:
//
//	func Providers() []foundation.ServiceProvider {
//	    return []foundation.ServiceProvider{
//	        &auth.ServiceProvider{},
//	        &cache.ServiceProvider{},
//	    }
//	}
func Providers() []match.GoNode {
	return []match.GoNode{
		Func(
			Ident("Providers"),
		),
		CompositeLit(
			ArrayType(
				SelectorExpr(
					Ident("foundation"),
					Ident("ServiceProvider"),
				),
				AnyNode(),
			),
		),
	}
}

// Rules matches the Rules() function that returns a slice of validation rules.
// It looks for a function returning []validation.Rule composite literals.
//
// Example usage:
//
//	match.GoFile("providers/validation_provider.go").
//	    Find(match.Rules()).
//	    Modify(addRuleToList(&CustomRule{})).
//	    Apply()
//
// This matches:
//
//	func (receiver *ValidationServiceProvider) Rules() []validation.Rule {
//	    return []validation.Rule{
//	        &rules.EmailRule{},
//	        &rules.PhoneRule{},
//	    }
//	}
func Rules() []match.GoNode {
	return []match.GoNode{
		Func(Ident("Rules")),
		TypeOf(&dst.ReturnStmt{}),
		CompositeLit(
			ArrayType(
				SelectorExpr(
					Ident("validation"),
					Ident("Rule"),
				),
				AnyNode(),
			),
		),
	}
}

// ProvidersInConfig represents the old logic of registering service providers inside the `config/app.go` file.
// If you are using the new code structure that registers service providers via the Providers() function in the
// bootstrap folder, you should use the Providers function instead.
func ProvidersInConfig() []match.GoNode {
	return []match.GoNode{
		Func(Ident("init")),
		CallExpr(
			SelectorExpr(
				Ident("config"),
				Ident("Add"),
			),
			GoNodes{
				BasicLit(strconv.Quote("app")),
				AnyNode(),
			},
		),
		KeyValueExpr(BasicLit(strconv.Quote("providers")), AnyNode()),
		CompositeLit(
			ArrayType(
				SelectorExpr(
					Ident("foundation"),
					Ident("ServiceProvider"),
				),
				AnyNode(),
			),
		),
	}
}

// RegisterFunc matches the Register() method in service providers.
// This is used to register services into the service container during the boot process.
//
// Example usage:
//
//	match.GoFile("providers/custom_provider.go").
//	    Find(match.RegisterFunc()).
//	    Modify(addServiceRegistration).
//	    Apply()
//
// This matches:
//
//	func (receiver *ServiceProvider) Register(app foundation.Application) {
//	    app.Bind("custom", func() (any, error) {
//	        return &CustomService{}, nil
//	    })
//	}
func RegisterFunc() []match.GoNode {
	return []match.GoNode{Func(Ident("Register"))}
}

// BootFunc matches the Boot() method in service providers.
// This is used to perform actions after all services are registered.
//
// Example usage:
//
//	match.GoFile("providers/custom_provider.go").
//	    Find(match.BootFunc()).
//	    Modify(addBootLogic).
//	    Apply()
//
// This matches:
//
//	func (receiver *ServiceProvider) Boot(app foundation.Application) {
//	    facades.Route().Get("/health", healthController.Check)
//	    facades.Event().Register(events)
//	}
func BootFunc() []match.GoNode {
	return []match.GoNode{Func(Ident("Boot"))}
}

// Seeders matches the Seeders() function that returns a slice of database seeders.
// It looks for a function returning []seeder.Seeder composite literals.
//
// Example usage:
//
//	match.GoFile("providers/seeder_provider.go").
//	    Find(match.Seeders()).
//	    Modify(addSeederToList(&UserSeeder{})).
//	    Apply()
//
// This matches:
//
//	func (receiver *SeederServiceProvider) Seeders() []seeder.Seeder {
//	    return []seeder.Seeder{
//	        &seeders.UserSeeder{},
//	        &seeders.PostSeeder{},
//	    }
//	}
func Seeders() []match.GoNode {
	return []match.GoNode{
		Func(Ident("Seeders")),
		TypeOf(&dst.ReturnStmt{}),
		CompositeLit(
			ArrayType(
				SelectorExpr(
					Ident("seeder"),
					Ident("Seeder"),
				),
				AnyNode(),
			),
		),
	}
}

// ValidationRules matches the rules() function that returns validation rules.
// It looks for a lowercase rules() function returning []validation.Rule composite literals.
// This is typically used in validation service providers for internal rule registration.
//
// Example usage:
//
//	match.GoFile("validation/setup.go").
//	    Find(match.ValidationRules()).
//	    Modify(addValidationRule(&EmailRule{})).
//	    Apply()
//
// This matches:
//
//	func rules() []validation.Rule {
//	    return []validation.Rule{
//	        &rules.Required{},
//	        &rules.Email{},
//	    }
//	}
func ValidationRules() []match.GoNode {
	return []match.GoNode{
		Func(Ident("rules")),
		TypeOf(&dst.ReturnStmt{}),
		CompositeLit(
			ArrayType(
				SelectorExpr(
					Ident("validation"),
					Ident("Rule"),
				),
				AnyNode(),
			),
		),
	}
}

// ValidationFilters matches the filters() function that returns validation filters.
// It looks for a lowercase filters() function returning []validation.Filter composite literals.
// This is typically used in validation service providers for internal filter registration.
//
// Example usage:
//
//	match.GoFile("validation/setup.go").
//	    Find(match.ValidationFilters()).
//	    Modify(addValidationFilter(&TrimFilter{})).
//	    Apply()
//
// This matches:
//
//	func filters() []validation.Filter {
//	    return []validation.Filter{
//	        &filters.Trim{},
//	        &filters.Lowercase{},
//	    }
//	}
func ValidationFilters() []match.GoNode {
	return []match.GoNode{
		Func(Ident("filters")),
		TypeOf(&dst.ReturnStmt{}),
		CompositeLit(
			ArrayType(
				SelectorExpr(
					Ident("validation"),
					Ident("Filter"),
				),
				AnyNode(),
			),
		),
	}
}

// FoundationSetup matches the Boot function containing foundation.Setup() chain calls.
// It matches both patterns:
//   - foundation.Setup().WithConfig(...).Start()
//   - foundation.Setup().WithMiddleware(...).WithConfig(...).Start()
//   - foundation.Setup().WithCommand(...).Start()
//   - return foundation.Setup().WithConfig(...).Start()
//
// Example usage:
//
//	GoFile("bootstrap/app.go").
//	    Find(match.FoundationSetup()).
//	    Modify(foundationSetupMiddleware(middleware)).
//	    Apply()
func FoundationSetup() []match.GoNode {
	return []match.GoNode{
		Func(Ident("Boot")),
		AnyOf(
			TypeOf(&dst.ExprStmt{}),
			TypeOf(&dst.ReturnStmt{}),
		),
	}
}
