package main

import (
	"fmt"
	"log"

	basePgx "github.com/cod3rcarl/wwd-grpc/base/pgx"
	grpcServer "github.com/cod3rcarl/wwd-grpc/pkg/server"
	"github.com/cod3rcarl/wwd-grpc/pkg/storage"
	wwd "github.com/cod3rcarl/wwd-grpc/pkg/wwdatabase"
	baseApp "github.com/cod3rcarl/wwd-subgraph/base/app"
	baseLogger "github.com/cod3rcarl/wwd-subgraph/base/logger"
	wwdatabaseClient "github.com/cod3rcarl/wwd-subgraph/pkg/client"
	server "github.com/cod3rcarl/wwd-subgraph/pkg/graphQLServer"
	"github.com/kelseyhightower/envconfig"
	"go.uber.org/zap"
)

func main() {
	logger := createLogger()
	storage := createStorage(logger)
	client := createWWDatabaseClient(logger)
	app := baseApp.NewApp(
		baseApp.WithLogger(logger),
	)

	grpcSrv := CreateGrpcServer(logger, storage)

	graphqlSrv := createGraphQLServer(logger, client)

	go func() {
		if err := grpcSrv.ListenAndServe(); err != nil {
			logger.Error("error running grpc server server", zap.Error(err))

			// a fatal log would call os.Exit(1), so cancel the context to trigger graceful shutdown
			app.Cancel()
		}
	}()

	go func() {
		if err := graphqlSrv.Start(); err != nil {
			logger.Error("error running GraphQL server", zap.Error(err))

			// a fatal log would call os.Exit(1), so cancel the context to trigger graceful shutdown
			app.Cancel()
		}
	}()

	serviceShutdownOrder := []baseApp.Stopper{
		grpcSrv,
		client,
	}

	app.HandleGracefulShutdown(serviceShutdownOrder)
}

func createLogger() *zap.Logger {
	cfg, err := baseLogger.ReadConfig()
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to read logger config, error=%v", err))
	}

	logger, err := baseLogger.NewLogger(cfg)
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to create logger, error=%v", err))
	}

	return logger
}

func createStorage(logger *zap.Logger) *storage.Service {
	cfg, err := basePgx.ReadConfig()
	if err != nil {
		logger.Fatal(
			"failed to read config",
			zap.String("service", "storage"),
			zap.Error(err),
		)
	}

	base, err := basePgx.NewPgx(cfg, basePgx.WithLogger(logger))
	if err != nil {
		logger.Fatal(
			"failed to create service",
			zap.String("service", "storage"),
			zap.Error(err),
		)
	}

	return storage.NewStorage(storage.WithBase(base))
}

func createWWDatabaseClient(logger *zap.Logger) *wwdatabaseClient.Client {
	cfg, err := wwdatabaseClient.ReadConfig()
	if err != nil {
		logger.Fatal(
			"failed to read config",
			zap.String("service", "wwdatabase"),
			zap.Error(err),
		)
	}
	client, err := wwdatabaseClient.NewClient(logger, cfg)
	if err != nil {
		logger.Fatal(
			"failed to create client",
			zap.String("service", "wwdatabase client"),
			zap.Error(err),
		)
	}

	return client
}

func createGraphQLServer(
	logger *zap.Logger,
	client *wwdatabaseClient.Client,
) *server.Server {
	var cfg server.Config
	if err := envconfig.Process("", &cfg); err != nil {
		logger.Fatal(
			"failed to read config",
			zap.String("service", "graphql-server"),
			zap.Error(err),
		)
	}

	svc, err := server.NewServer(
		cfg,
		logger,
		server.WithWWDatabase(client),
	)
	if err != nil {
		logger.Fatal(
			"failed to create service",
			zap.String("service", "graphql-server"),
			zap.Error(err),
		)
	}

	return svc
}

func CreateGrpcServer(
	l *zap.Logger,
	s *storage.Service,
) *grpcServer.Service {
	cfg, err := grpcServer.ReadConfig()
	if err != nil {
		l.Fatal(
			"failed to read config",
			zap.String("service", "grpc-server"),
			zap.Error(err),
		)
	}

	wwdService := wwd.NewService(
		wwd.WithLogger(l),
		wwd.WithStorage(s),
	)

	return grpcServer.NewServer(
		cfg,
		l,
		grpcServer.WithClient(wwdService),
	)
}
