package server

import (
	ww "github.com/cod3rcarl/wwd-grpc/pkg/wwdatabase"
)

type Option interface {
	apply(s *Service)
}

type optionFunc func(s *Service)

func (f optionFunc) apply(s *Service) {
	f(s)
}

func Options(opts ...Option) Option {
	return optionFunc(func(s *Service) {
		for _, opt := range opts {
			opt.apply(s)
		}
	})
}

func (s *Service) withOptions(opts ...Option) *Service {
	for _, opt := range opts {
		opt.apply(s)
	}

	return s
}

type clientOption struct {
	wwdatabase ww.ServiceInterface
}

func (o clientOption) apply(s *Service) {
	s.wwdatabase = o.wwdatabase
}
