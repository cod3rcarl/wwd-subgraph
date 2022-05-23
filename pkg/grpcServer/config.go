package server

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	Host string `envconfig:"GRPC_HOST" required:"true"`
	Port string `envconfig:"GRPC_PORT" required:"true"`
}

func ReadConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, errors.Errorf("failed to parse config; error=%v", err)
	}

	return cfg, nil
}
