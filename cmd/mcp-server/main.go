package main

import (
	"bytes"
	"context"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"

	"github.com/uptrace/mcp/appconf"
	"github.com/uptrace/mcp/bootstrap"
	"github.com/uptrace/mcp/tools"
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
					return runServer(ctx, logger, client, conf, cmd)
				}),
			)
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func runServer(
	ctx context.Context,
	logger *slog.Logger,
	client *uptraceapi.Client,
	conf *appconf.Config,
	cmd *cli.Command,
) error {
	if cmd.String("http") != "" {
		return runHTTPServer(ctx, logger, client, conf, cmd)
	}
	return runStdioServer(ctx, logger, client, conf, cmd)
}

func runStdioServer(
	ctx context.Context,
	logger *slog.Logger,
	client *uptraceapi.Client,
	conf *appconf.Config,
	cmd *cli.Command,
) error {
	server := newServer(client, conf)
	debug := cmd.Bool("debug")

	logger.Info("starting MCP server (stdio)",
		slog.String("name", bootstrap.AppName),
		slog.String("version", bootstrap.AppVersion),
		slog.Bool("debug", debug),
	)

	var transport mcp.Transport = mcp.NewStdioTransport()
	if debug {
		transport = mcp.NewLoggingTransport(transport, os.Stderr)
	}

	return server.Run(ctx, transport)
}

func runHTTPServer(
	ctx context.Context,
	logger *slog.Logger,
	client *uptraceapi.Client,
	conf *appconf.Config,
	cmd *cli.Command,
) error {
	httpAddr := cmd.String("http")
	debug := cmd.Bool("debug")

	var handler http.Handler = mcp.NewStreamableHTTPHandler(func(r *http.Request) *mcp.Server {
		return newServer(client, conf)
	}, nil)

	if debug {
		handler = loggingMiddleware(logger, handler)
	}

	logger.Info("starting MCP server (HTTP)",
		slog.String("name", bootstrap.AppName),
		slog.String("version", bootstrap.AppVersion),
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

func newServer(client *uptraceapi.Client, conf *appconf.Config) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    bootstrap.AppName,
		Version: bootstrap.AppVersion,
	}, nil)
	tools.Register(server, client, conf)
	return server
}

func loggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewReader(body))

		logger.Debug("HTTP request",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.String("body", string(body)),
		)

		next.ServeHTTP(w, r)
	})
}
