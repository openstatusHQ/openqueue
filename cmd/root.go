package main

import (
	"context"
	"os"

	"github.com/openstatushq/openqueue/cmd/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

func main() {
	if err := os.Setenv("TZ", "UTC"); err != nil {
		panic(err)
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	zerolog.DefaultContextLogger = func() *zerolog.Logger {
		logger := log.With().Caller().Logger()
		return &logger
	}()
	app := &cli.Command{
		Name: "openqueue",
		Commands: []*cli.Command{
			server.Command(),
		},
	}
	ctx := context.Background()
	log.Ctx(ctx).Info().Msg("OpenQueue - Queue management system")
	if err := app.Run(ctx, os.Args); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Application error")
		os.Exit(1)
	}
}
