package bootstrap

import (
	"context"
	"log/slog"
	"os"
	"strings"

	slogattrs "github.com/go-slog/otelslog"
	slogmulti "github.com/samber/slog-multi"
	"github.com/urfave/cli/v3"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/uptrace/mcp/appconf"
)

const (
	AppName    = "uptrace-mcp"
	AppVersion = "0.1.0"
)

func New(ctx context.Context, conf *appconf.Config, options ...fx.Option) *fx.App {
	return fx.New(
		fx.StartTimeout(conf.Service.StartTimeout),
		fx.StopTimeout(conf.Service.StopTimeout),

		fx.Supply(
			fx.Annotate(
				ctx,
				fx.As(new(context.Context)),
			),
			conf,
		),

		fx.Provide(
			NewSlog,
		),
		fx.WithLogger(func(logger *slog.Logger) fxevent.Logger {
			l := &fxevent.SlogLogger{Logger: logger}
			l.UseLogLevel(slog.LevelDebug)
			return l
		}),

		fx.Options(options...),
	)
}

func Run(ctx context.Context, cmd *cli.Command, options ...fx.Option) error {
	conf, err := appconf.Load(cmd.String("config"))
	if err != nil {
		return err
	}

	app := New(ctx, conf, options...)
	app.Run()
	return app.Err()
}

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
				//AddSource: true,
			}),
		),
		otelslog.NewHandler("uptrace", otelslog.WithSource(true)),
	))

	return LoggerResults{
		Logger: logger,
		Level:  level,
	}
}
