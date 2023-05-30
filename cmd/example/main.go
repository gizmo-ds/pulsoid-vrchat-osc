package main

import (
	"os"

	"github.com/gizmo-ds/go-cli-template/cmd/example/commands"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/urfave/cli/v2"
)

var AppVersion = "development"

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

func main() {
	_ = (&cli.App{
		Name:     "example",
		Version:  AppVersion,
		Suggest:  true,
		Commands: commands.Commands,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-color",
				Usage: "Disable color output",
			},
			&cli.IntFlag{
				Name: "log-level",
				Usage: `Set the log level
(0 = Debug, 1 = Info, 2 = Warn, 3 = Error, 4 = Fatal, 5 = Panic, other = NoLog)`,
				Value:       1,
				DefaultText: "Info",
			},
		},
		Before: func(c *cli.Context) error {
			zerolog.SetGlobalLevel(zerolog.Level(c.Int("log-level")))
			log.Logger = zerolog.New(zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: c.Bool("no-color"),
			}).With().Timestamp().Logger()
			return nil
		},
		ExitErrHandler: func(c *cli.Context, err error) {
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to run command")
			}
		},
	}).Run(os.Args)
}
