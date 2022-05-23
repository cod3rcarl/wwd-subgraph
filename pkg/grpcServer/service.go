package server

import (
	"fmt"
	"net"

	client "github.com/cod3rcarl/wwd-grpc/pkg/wwdatabase"
	pb "github.com/cod3rcarl/wwd-protorepo-wwdatabase/v1"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const ServiceName = "grpc-server"

type Service struct {
	config Config
	logger *zap.Logger
	server *grpc.Server
	serviceOptions
	pb.UnimplementedWwdatabaseServer
}

type serviceOptions struct {
	clientOption
}

func NewServer(cfg Config, logger *zap.Logger, opts ...Option) *Service {
	srv := grpc.NewServer()

	// register the gRPC server for reflection to expose available endpoints
	reflection.Register(srv)

	svc := &Service{
		config: cfg,
		logger: logger,
		server: srv,
	}

	pb.RegisterWwdatabaseServer(srv, svc)

	return svc.withOptions(opts...)
}

func WithClient(c client.ServiceInterface) Option {
	return clientOption{c}
}

func (s *Service) Close() {
	s.logger.Info("stopping", zap.String("service", "grpc-server"))
	s.server.GracefulStop()
}

func (s *Service) ListenAndServe() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.config.Host, s.config.Port))
	if err != nil {
		return errors.Errorf("failed to listen: %v", err)
	}

	s.logger.Info("starting grpc server", zap.String("address", lis.Addr().String()))
	if err2 := s.server.Serve(lis); err2 != nil {
		return errors.Errorf("failed to serve: %v", err2)
	}

	return nil
}
