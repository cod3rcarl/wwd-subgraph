package pgx

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	// # Example DSN
	// user=jack password=secret host=pg.example.com port=5432 dbname=mydb sslmode=verify-ca pool_max_conns=10
	//
	// # Example URL
	// postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10
	DBConnection string `envconfig:"DB_CONNECTION" required:"true"`
}

func ReadConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, errors.Errorf("failed to parse config; error=%v", err)
	}

	return cfg, nil
}
