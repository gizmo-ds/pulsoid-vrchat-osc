package main

import (
	"os"

	"github.com/gizmo-ds/pulsoid-vrchat-osc/cmd/cli/action"
	"github.com/gizmo-ds/pulsoid-vrchat-osc/internal/global"
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
		Name:    "pulsoid-vrchat-osc",
		Version: AppVersion,
		Suggest: true,
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
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to config file",
				Value:   "config.toml",
			},
		},
		Before: func(c *cli.Context) error {
			log.Logger = zerolog.New(zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: c.Bool("no-color"),
			}).With().Timestamp().Logger()
			err := global.LoadConfig(c.String("config"))
			if err != nil {
				return err
			}
			logLevel := c.Int("log-level")
			if global.Config.Logger.Level != nil {
				logLevel = *global.Config.Logger.Level
			}
			zerolog.SetGlobalLevel(zerolog.Level(logLevel))
			return nil
		},
		ExitErrHandler: func(c *cli.Context, err error) {
			if err != nil {
				log.Fatal().Err(err).Msg("Failed to run command")
			}
		},
		Action: action.Action,
	}).Run(os.Args)
}
