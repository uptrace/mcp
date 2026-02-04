package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"
	"go.uber.org/fx"

	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/bootstrap"
	"github.com/uptrace/mcp/uptraceapi"
)

func main() {
	cmd := &cli.Command{
		Name:  "mcp-server",
		Usage: "Uptrace MCP server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "Path to config file",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "http",
				Usage:   "HTTP address to listen on (e.g., :8080). If not set, uses stdio transport.",
				Aliases: []string{"H"},
			},
			&cli.BoolFlag{
				Name:    "debug",
				Usage:   "Enable debug mode (logs all JSON-RPC messages)",
				Aliases: []string{"d"},
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return bootstrap.Run(
				ctx,
				cmd,
				fx.Provide(bootstrap.NewUptraceClient),
				fx.Invoke(func(
					ctx context.Context,
					logger *slog.Logger,
					client *uptraceapi.Client,
					conf *appconf.Config,
				) error {
					return bootstrap.RunServer(ctx, logger, client, conf, cmd)
				}),
			)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
