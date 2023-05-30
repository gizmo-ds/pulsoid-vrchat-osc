package commands

import (
	"github.com/gizmo-ds/go-cli-template/pkg/hello"
	"github.com/urfave/cli/v2"
)

func init() {
	Commands = append(Commands, HelloCommand)
}

var HelloCommand = &cli.Command{
	Name:      "hello",
	Usage:     "say hello",
	ArgsUsage: "[name]",
	Action: func(c *cli.Context) error {
		println(hello.SayHello(c.Args().First()))
		return nil
	},
}
