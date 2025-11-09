package server

import (
	"context"
	"strconv"

	"github.com/openstatushq/openqueue/pkg/config"
	"github.com/openstatushq/openqueue/pkg/server"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

func action(ctx context.Context, cmd *cli.Command) error {

	log.Ctx(ctx).Info().Msg("Starting OpenQueue server...")
	port := cmd.String("port")
	p, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	configPath := cmd.String("config")
	err = config.LoadConfigFile(ctx, configPath)

	cfg := config.GetConfig()

	if err != nil {
		return err
	}
	opts := server.Options{}
	opts.Port = p

	for _, q := range cfg.Queues {
		opts.Queues = append(opts.Queues, struct {
			Name string
			DB   string
		}{
			Name: q.Name,
			DB:   q.DB,
		})
	}

	err = server.NewServer(ctx, opts)
	if err != nil {
		return err
	}

	return nil
}
