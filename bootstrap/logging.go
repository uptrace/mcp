package bootstrap

import (
	"log/slog"
	"os"
	"strings"

	slogattrs "github.com/go-slog/otelslog"
	slogmulti "github.com/samber/slog-multi"
	"github.com/uptrace/mcp/appconf"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.uber.org/fx"
)

type LoggerResults struct {
	fx.Out

	Logger *slog.Logger
	Level  *slog.LevelVar
}

func NewSlog(conf *appconf.Config) LoggerResults {
	level := new(slog.LevelVar)

	switch strings.ToLower(conf.Logging.Level) {
	case "debug":
		level.Set(slog.LevelDebug)
	case "info":
		level.Set(slog.LevelInfo)
	case "warn":
		level.Set(slog.LevelWarn)
	case "error":
		level.Set(slog.LevelError)
	default:
		level.Set(slog.LevelInfo)
	}

	logger := slog.New(slogmulti.Fanout(
		slogattrs.NewHandler(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
				Level: level,
			}),
		),
		otelslog.NewHandler("uptrace", otelslog.WithSource(true)),
	))

	return LoggerResults{
		Logger: logger,
		Level:  level,
	}
}
