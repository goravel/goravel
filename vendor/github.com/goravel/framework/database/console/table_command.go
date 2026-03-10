package console

import (
	"fmt"
	"strings"

	"github.com/goravel/framework/contracts/config"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/console/command"
	"github.com/goravel/framework/contracts/database/driver"
	"github.com/goravel/framework/contracts/database/schema"
)

type TableCommand struct {
	config config.Config
	schema schema.Schema
}

func NewTableCommand(config config.Config, schema schema.Schema) *TableCommand {
	return &TableCommand{
		config: config,
		schema: schema,
	}
}

// Signature The name and signature of the console command.
func (r *TableCommand) Signature() string {
	return "db:table"
}

// Description The console command description.
func (r *TableCommand) Description() string {
	return "Display information about the given database table"
}

// Extend The console command extend.
func (r *TableCommand) Extend() command.Extend {
	return command.Extend{
		Category:  "db",
		ArgsUsage: "  [--] [<table>]",
		Flags: []command.Flag{
			&command.StringFlag{
				Name:    "database",
				Aliases: []string{"d"},
				Usage:   "The database connection",
			},
		},
	}
}

// Handle Execute the console command.
func (r *TableCommand) Handle(ctx console.Context) error {
	ctx.NewLine()
	r.schema = r.schema.Connection(ctx.Option("database"))
	table := ctx.Argument(0)
	tables, err := r.schema.GetTables()
	if err != nil {
		ctx.Error(fmt.Sprintf("Failed to get tables: %s", err.Error()))
		return nil
	}
	if len(table) == 0 {
		table, err = ctx.Choice("Which table would you like to inspect?", func() (choices []console.Choice) {
			for i := range tables {
				choices = append(choices, console.Choice{
					Key:   tables[i].Name,
					Value: tables[i].Name,
				})
			}
			return
		}())
		if err != nil {
			ctx.Line(err.Error())
			return nil
		}
	}
	for i := range tables {
		if tables[i].Name == table {
			r.display(ctx, tables[i])
			return nil
		}
	}
	if len(table) > 0 {
		ctx.Warning(fmt.Sprintf("Table '%s' doesn't exist.", table))
		ctx.NewLine()
	}
	return nil
}

func (r *TableCommand) display(ctx console.Context, table driver.Table) {
	columns, err := r.schema.GetColumns(table.Name)
	if err != nil {
		ctx.Error(fmt.Sprintf("Failed to get columns: %s", err.Error()))
		return
	}
	indexes, err := r.schema.GetIndexes(table.Name)
	if err != nil {
		ctx.Error(fmt.Sprintf("Failed to get indexes: %s", err.Error()))
		return
	}
	foreignKeys, err := r.schema.GetForeignKeys(table.Name)
	if err != nil {
		ctx.Error(fmt.Sprintf("Failed to get foreign keys: %s", err.Error()))
		return
	}
	name := table.Name
	if len(table.Schema) > 0 {
		name = table.Schema + "." + table.Name
	}
	ctx.TwoColumnDetail(fmt.Sprintf("<fg=green;op=bold>%s</>", name), fmt.Sprintf("<fg=gray>%s</>", table.Comment))
	ctx.TwoColumnDetail("Columns", fmt.Sprintf("%d", len(columns)))
	ctx.TwoColumnDetail("Size", fmt.Sprintf("%.3f MB", float64(table.Size)/1024/1024))
	if len(table.Engine) > 0 {
		ctx.TwoColumnDetail("Engine", table.Engine)
	}
	if len(table.Collation) > 0 {
		ctx.TwoColumnDetail("Collation", table.Collation)
	}
	if len(columns) > 0 {
		ctx.NewLine()
		ctx.TwoColumnDetail("<fg=green;op=bold>Column</>", "Type")
		for i := range columns {
			var (
				key        = columns[i].Name
				value      = columns[i].Type
				attributes []string
			)
			if columns[i].Autoincrement {
				attributes = append(attributes, "autoincrement")
			}
			attributes = append(attributes, columns[i].TypeName)
			if columns[i].Nullable {
				attributes = append(attributes, "nullable")
			}
			if len(columns[i].Collation) > 0 {
				attributes = append(attributes, columns[i].Collation)
			}
			key = fmt.Sprintf("%s <fg=gray>%s</>", key, strings.Join(attributes, ", "))
			if columns[i].Default != "" {
				value = fmt.Sprintf("<fg=gray>%s</> %s", columns[i].Default, value)
			}
			ctx.TwoColumnDetail(key, value)
		}
	}
	if len(indexes) > 0 {
		ctx.NewLine()
		ctx.TwoColumnDetail("<fg=green;op=bold>Index</>", "")
		for i := range indexes {
			attributes := []string{indexes[i].Type}
			if len(indexes[i].Columns) > 1 {
				attributes = append(attributes, "compound")
			}
			if indexes[i].Unique {
				attributes = append(attributes, "unique")
			}
			if indexes[i].Primary {
				attributes = append(attributes, "primary")
			}
			ctx.TwoColumnDetail(fmt.Sprintf("%s <fg=gray>%s</>", indexes[i].Name, strings.Join(indexes[i].Columns, ", ")), strings.Join(attributes, ", "))
		}
	}
	if len(foreignKeys) > 0 {
		ctx.NewLine()
		ctx.TwoColumnDetail("<fg=green;op=bold>Foreign Key</>", "On Update / On Delete")
		for i := range foreignKeys {
			key := fmt.Sprintf("%s <fg=gray>%s references %s on %s</>",
				foreignKeys[i].Name,
				strings.Join(foreignKeys[i].Columns, ", "),
				strings.Join(foreignKeys[i].ForeignColumns, ", "),
				foreignKeys[i].ForeignTable)
			ctx.TwoColumnDetail(key, fmt.Sprintf("%s / %s", foreignKeys[i].OnUpdate, foreignKeys[i].OnDelete))
		}
	}
	ctx.NewLine()
}
