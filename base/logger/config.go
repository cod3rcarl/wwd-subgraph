package logger

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	Level string `envconfig:"LOGGER_LEVEL" default:"INFO"`
}

func ReadConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, errors.Errorf("failed to parse config; error=%v", err)
	}

	return cfg, nil
}
