package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"

	"github.com/uptrace/mcp/bootstrap"
	"github.com/uptrace/mcp/tools"
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
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return bootstrap.Run(
				ctx,
				cmd,
				fx.Invoke(runServer),
			)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func runServer(ctx context.Context, logger *slog.Logger) error {
	server := newServer()
	tools.Register(server)

	logger.Info("starting MCP server", slog.String("name", bootstrap.AppName), slog.String("version", bootstrap.AppVersion))

	return server.Run(ctx, mcp.NewStdioTransport())
}

func newServer() *mcp.Server {
	return mcp.NewServer(&mcp.Implementation{
		Name:    bootstrap.AppName,
		Version: bootstrap.AppVersion,
	}, nil)
}
