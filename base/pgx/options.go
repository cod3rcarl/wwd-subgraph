package pgx

import "go.uber.org/zap"

type Option interface {
	apply(s *Service)
}

func (s *Service) withOptions(opts ...Option) *Service {
	for _, opt := range opts {
		opt.apply(s)
	}

	return s
}

type loggerOption struct {
	Logger *zap.Logger
}

func (o loggerOption) apply(s *Service) {
	s.Logger = o.Logger
}
