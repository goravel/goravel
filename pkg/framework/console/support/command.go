package support

import "github.com/urfave/cli/v2"

type Command interface {
	Signature() string
	Description() string
	Flags() []cli.Flag
	Subcommands() []*cli.Command
	Handle(c *cli.Context) error
}
