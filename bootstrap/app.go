package bootstrap

import (
	"context"
	"log/slog"

	"github.com/uptrace/mcp/appconf"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

const (
	AppName    = "uptrace-mcp"
	AppVersion = "0.2.0"
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
