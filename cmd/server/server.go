package server

import (
	"context"
	"strconv"

	"github.com/openstatushq/openqueue/pkg/server"
	"github.com/urfave/cli/v3"
)

func action(ctx context.Context, cmd *cli.Command) error {

	port := cmd.String("port")
	p, err := strconv.Atoi(port)
	if err != nil {
		return err
	}

	err = server.NewServer(p)
	if err != nil {
		return err
	}

	return nil
}
