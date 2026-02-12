package bootstrap

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/uptraceapi"
	"github.com/urfave/cli/v3"
)

func RunServer(
	ctx context.Context,
	logger *slog.Logger,
	client *uptraceapi.Client,
	conf *appconf.Config,
	server *mcp.Server,
	cmd *cli.Command,
) error {
	if cmd.String("http") != "" {
		return runHTTPServer(ctx, logger, server, conf, cmd)
	}
	return runStdioServer(ctx, logger, server, cmd)
}

func runStdioServer(
	ctx context.Context,
	logger *slog.Logger,
	server *mcp.Server,
	cmd *cli.Command,
) error {
	debug := cmd.Bool("debug")

	logger.Info("starting MCP server (stdio)",
		slog.String("name", AppName),
		slog.String("version", AppVersion),
		slog.Bool("debug", debug),
	)

	var transport mcp.Transport = &mcp.StreamableClientTransport{}
	if debug {
		transport = &mcp.LoggingTransport{}
	}

	return server.Run(ctx, transport)
}

func runHTTPServer(
	ctx context.Context,
	logger *slog.Logger,
	mcpServer *mcp.Server,
	conf *appconf.Config,
	cmd *cli.Command,
) error {
	httpAddr := cmd.String("http")
	debug := cmd.Bool("debug")

	var handler http.Handler = mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return mcpServer
	}, nil)

	if debug {
		handler = loggingMiddleware(handler, conf.Logging.MaxBodySize)
	}

	logger.Info("starting MCP server (HTTP)",
		slog.String("name", AppName),
		slog.String("version", AppVersion),
		slog.String("addr", httpAddr),
		slog.Bool("debug", debug),
	)

	server := &http.Server{
		Addr:    httpAddr,
		Handler: handler,
	}

	go func() {
		<-ctx.Done()
		server.Shutdown(context.Background())
	}()

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
