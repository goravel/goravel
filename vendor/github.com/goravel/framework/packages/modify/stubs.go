package modify

import (
	"strings"

	"github.com/goravel/framework/packages/paths"
	"github.com/goravel/framework/support/env"
)

func commands() string {
	content := `package DummyPackage

import "github.com/goravel/framework/contracts/console"

func Commands() []console.Command {
	return []console.Command{}
}
`

	return replaceDummyPackage(content)
}

func filters() string {
	content := `package DummyPackage

import "github.com/goravel/framework/contracts/validation"

func Filters() []validation.Filter {
	return []validation.Filter{}
}
`

	return replaceDummyPackage(content)
}

func migrations() string {
	content := `package DummyPackage

import "github.com/goravel/framework/contracts/database/schema"

func Migrations() []schema.Migration {
	return []schema.Migration{}
}
`

	return replaceDummyPackage(content)
}

func providers() string {
	content := `package DummyPackage

import "github.com/goravel/framework/contracts/foundation"

func Providers() []foundation.ServiceProvider {
	return []foundation.ServiceProvider{}
}
`

	return replaceDummyPackage(content)
}

func jobs() string {
	content := `package DummyPackage

import "github.com/goravel/framework/contracts/queue"

func Jobs() []queue.Job {
	return []queue.Job{}
}
`

	return replaceDummyPackage(content)
}

func seeders() string {
	content := `package DummyPackage

import "github.com/goravel/framework/contracts/database/seeder"

func Seeders() []seeder.Seeder {
	return []seeder.Seeder{}
}
`

	return replaceDummyPackage(content)
}

func rules() string {
	content := `package DummyPackage

import "github.com/goravel/framework/contracts/validation"

func Rules() []validation.Rule {
	return []validation.Rule{}
}
`

	return replaceDummyPackage(content)
}

func replaceDummyPackage(content string) string {
	bootstrapPackage := paths.NewPaths(env.MainPath()).Bootstrap().Package()

	return strings.ReplaceAll(content, "DummyPackage", bootstrapPackage)
}
