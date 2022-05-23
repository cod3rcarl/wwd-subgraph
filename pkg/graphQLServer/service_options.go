package server

import (
	"github.com/cod3rcarl/wwd-subgraph/pkg/client"
)

type Option interface {
	apply(s *Server)
}

type optionFunc func(s *Server)

func (f optionFunc) apply(s *Server) {
	f(s)
}

func Options(opts ...Option) Option {
	return optionFunc(func(s *Server) {
		for _, opt := range opts {
			opt.apply(s)
		}
	})
}

func (s *Server) withOptions(opts ...Option) *Server {
	for _, opt := range opts {
		opt.apply(s)
	}

	return s
}

type wwdatabaseOption struct {
	client *client.Client
}

func (o wwdatabaseOption) apply(s *Server) {
	s.client = o.client
}
