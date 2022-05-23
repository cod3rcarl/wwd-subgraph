package app

import (
	"context"

	"go.uber.org/zap"
)

type App struct {
	//nolint:containedctx // we need the context to be called in HandleGracefulShutdown
	ctx    context.Context
	Cancel context.CancelFunc
	appOptions
}

type appOptions struct {
	loggerOption
}

func NewApp(opts ...Option) *App {
	app := &App{
		appOptions: appOptions{
			loggerOption: loggerOption{
				logger: zap.NewNop(),
			},
		},
	}

	app.ctx, app.Cancel = context.WithCancel(context.Background())

	return app.withOptions(opts...)
}

func WithLogger(l *zap.Logger) Option {
	return loggerOption{logger: l}
}
