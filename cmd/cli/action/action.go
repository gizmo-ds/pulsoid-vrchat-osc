package action

import (
	"github.com/urfave/cli/v2"
)

func Action(c *cli.Context) error {
	NewPulsoid().Start()
	return nil
}
