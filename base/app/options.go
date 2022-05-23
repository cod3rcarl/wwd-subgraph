package app

import (
	"go.uber.org/zap"
)

type Option interface {
	apply(s *App)
}

type optionFunc func(a *App)

func (f optionFunc) apply(a *App) {
	f(a)
}

func Options(opts ...Option) Option {
	return optionFunc(func(a *App) {
		for _, opt := range opts {
			opt.apply(a)
		}
	})
}

func (a *App) withOptions(opts ...Option) *App {
	for _, opt := range opts {
		opt.apply(a)
	}

	return a
}

type loggerOption struct {
	logger *zap.Logger
}

func (o loggerOption) apply(a *App) {
	a.logger = o.logger
}
