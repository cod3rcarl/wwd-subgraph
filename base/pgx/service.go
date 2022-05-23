package pgx

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

var ErrConnectionString = errors.New(`failed to parse database config, use a db_dsn if db_url has special characters`)

type Service struct {
	Pool *pgxpool.Pool
	serviceOptions
}

type serviceOptions struct {
	loggerOption
}

func NewPgx(cfg Config, opts ...Option) (*Service, error) {
	svc := &Service{
		serviceOptions: serviceOptions{
			loggerOption: loggerOption{
				Logger: zap.NewNop(),
			},
		},
	}

	svc.withOptions(opts...)

	poolConfig, err := pgxpool.ParseConfig(cfg.DBConnection)
	if err != nil {
		return nil, ErrConnectionString
	}

	// increase log level to warn, otherwise every query is logged at info level
	poolConfig.ConnConfig.Logger = zapadapter.NewLogger(svc.Logger.WithOptions(zap.IncreaseLevel(zap.WarnLevel)))

	svc.Pool, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, errors.Errorf("failed to connect to database: %v", err)
	}

	return svc, nil
}

func WithLogger(l *zap.Logger) Option {
	return loggerOption{Logger: l}
}

func (s *Service) Stop() {
	s.Logger.Info("gracefully shutting down pgx service")
	s.Pool.Close()
}

func (s *Service) NotFound(err error) bool {
	return pgxscan.NotFound(err)
}
