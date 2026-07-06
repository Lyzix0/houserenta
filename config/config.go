package config

import (
	"fmt"
	"net/url"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type (
	Config struct {
		Logger logger
		PG     pg
	}
	pg struct {
		User     string        `envconfig:"POSTGRES_USER" required:"true"`
		Password string        `envconfig:"POSTGRES_PASSWORD" required:"true"`
		DB       string        `envconfig:"POSTGRES_DB" required:"true"`
		Host     string        `envconfig:"POSTGRES_HOST" required:"true"`
		Port     string        `envconfig:"POSTGRES_PORT" default:"5432"`
		Timeout  time.Duration `envconfig:"POSTGRES_TIMEOUT" default:"10s"`
		PoolMax  int           `envconfig:"PG_POOL_MAX" default:"10"`

		URL string `ignored:"true"`
	}

	logger struct {
		Level  string `envconfig:"LOGGER_LEVEL" default:"DEBUG"`
		Folder string `envconfig:"LOGGER_FOLDER" required:"true"`
	}
)

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	config.PG.URL = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&connect_timeout=%d",
		url.QueryEscape(config.PG.User),
		url.QueryEscape(config.PG.Password),
		config.PG.Host,
		config.PG.Port,
		config.PG.DB,
		int(config.PG.Timeout.Seconds()),
	)

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get Logger config: %w", err)
		panic(err)
	}

	return config
}
