package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/gizmo-ds/pulsoid-vrchat-osc/cmd/cli/action"
	"github.com/gizmo-ds/pulsoid-vrchat-osc/internal/global"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/urfave/cli/v2"
)

func init() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
}

func main() {
	_ = (&cli.App{
		Name:    "pulsoid-vrchat-osc",
		Version: global.AppVersion,
		Suggest: true,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-color",
				Usage: "Disable color output",
			},
			&cli.StringFlag{
				Name: "log-level",
				Usage: fmt.Sprintf("Set the log level\n(%s)",
					strings.Join([]string{
						zerolog.LevelDebugValue, zerolog.LevelInfoValue, zerolog.LevelWarnValue,
						zerolog.LevelErrorValue, zerolog.LevelFatalValue, zerolog.LevelPanicValue,
					}, ", "),
				),
				Value: zerolog.LevelInfoValue,
			},
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to config file",
				Value:   "config.toml",
			},
		},
		Before: func(c *cli.Context) error {
			err := global.LoadConfig(c.String("config"))
			if err != nil {
				return err
			}

			logLevelStr := c.String("log-level")
			if !c.IsSet("log-level") && global.Config.Logger.Level != nil {
				// TODO: 兼容`v0.1.3`及之前的配置文件
				switch reflect.TypeOf(global.Config.Logger.Level).Kind() {
				case reflect.String:
					logLevelStr = global.Config.Logger.Level.(string)
				case reflect.Int64:
					if s := zerolog.Level(global.Config.Logger.Level.(int64)).String(); s != "" {
						logLevelStr = s
					}
				default:
					log.Fatal().Msg("Invalid log level")
				}
			}
			logLevel, err := zerolog.ParseLevel(logLevelStr)
			if err != nil {
				return err
			}
			zerolog.SetGlobalLevel(logLevel)

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
		Action: action.Action,
	}).Run(os.Args)
}
