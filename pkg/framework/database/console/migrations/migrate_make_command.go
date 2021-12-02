package migrations

import (
	"github.com/goravel/framework/database/migrations"
	"github.com/urfave/cli/v2"
	"log"
)

type MigrateMakeCommand struct {
}

func (receiver MigrateMakeCommand) Signature() string {
	return "make:migration"
}

func (receiver MigrateMakeCommand) Description() string {
	return "Create a new migration file"
}

func (receiver MigrateMakeCommand) Flags() []cli.Flag {
	var flags []cli.Flag

	return flags
}

func (receiver MigrateMakeCommand) Subcommands() []*cli.Command {
	var subcommands []*cli.Command

	return subcommands
}

func (receiver MigrateMakeCommand) Handle(c *cli.Context) error {
	name := c.Args().First()

	if name == "" {
		log.Fatalln(`Not enough arguments (missing: "name").`)
	}

	table, create := TableGuesser{}.Guess(name)

	migrations.MigrateCreator{}.Create(name, table, create)

	log.Printf("Created Migration: %s", name)

	return nil
}
