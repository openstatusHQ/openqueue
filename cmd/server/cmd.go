package server

import "github.com/urfave/cli/v3"

func Command() *cli.Command {
	cmd := &cli.Command{
		Name:   "server",
		Usage:  "Run queue server",
		Action: action,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "port",
				Aliases: []string{"p"},
				Value:   "8080",
				Usage:   "Port to listen on",
				Sources: cli.EnvVars("OPENSTATUS_PORT"),
			},
		},
	}
	return cmd
}
