package server

import (
	"context"
	"net/http"

	"github.com/cod3rcarl/wwd-subgraph/pkg/client"
	graphresolvers "github.com/cod3rcarl/wwd-subgraph/pkg/graphQLServer/graph/resolvers"

	graphgenerated "github.com/cod3rcarl/wwd-subgraph/pkg/graphQLServer/graph/generated"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

type Config struct {
	Port               string   `envconfig:"PORT" default:"8181"`
	CORSAllowedOrigins []string `envconfig:"CORS_ALLOWED_ORIGINS" default:"*"`
	CORSAllowedHeaders []string `envconfig:"CORS_ALLOWED_HEADERS" default:"Authorization,Content-Type"`
	CORSAllowedMethods []string `envconfig:"CORS_ALLOWED_METHODS" default:"GET,POST,HEAD,OPTIONS"`
}

type Server struct {
	config  Config
	logger  *zap.Logger
	handler http.Handler

	serviceOptions
}

type serviceOptions struct {
	wwdatabaseOption
}

func NewServer(cfg Config, logger *zap.Logger, opts ...Option) (*Server, error) {
	svr := &Server{
		config: cfg,
		logger: logger,
	}

	svr.withOptions(opts...)

	svr.handler = svr.createHandler()

	return svr, nil
}

func (s *Server) createHandler() http.Handler {
	resolver := &graphresolvers.Resolver{
		Logger: s.logger,
		Server: s.wwdatabaseOption.client,
	}
	srv := handler.NewDefaultServer(graphgenerated.NewExecutableSchema(graphgenerated.Config{Resolvers: resolver}))

	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		s.logger.Error("server panic", zap.Any("err", err))

		return errors.Errorf("server panic")
	})

	return s.newRouter(srv, s.config)
}

func (s *Server) Start() error {
	s.logger.Info("starting server")

	if err := http.ListenAndServe(":"+s.config.Port, s.handler); err != nil {
		return errors.Errorf("server error: %v", err)
	}

	return nil
}

func WithWWDatabase(p *client.Client) Option {
	return wwdatabaseOption{p}
}
