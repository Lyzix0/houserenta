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
		HTTP   http
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

	http struct {
		AllowedOrigins []string `envconfig:"ALLOWED_ORIGINS" default:"http://localhost:3000,http://localhost:5050"`
	}
)

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig: %w", err)
	}

	if len(config.HTTP.AllowedOrigins) == 0 {
		return Config{}, fmt.Errorf(
			"ALLOWED_ORIGINS must not be empty: an empty list is treated by fiber/cors as allowing all origins, " +
				"which is incompatible with session cookies (CORS AllowCredentials) — list explicit origins instead",
		)
	}

	for _, origin := range config.HTTP.AllowedOrigins {
		if origin == "*" {
			return Config{}, fmt.Errorf(
				"ALLOWED_ORIGINS must not contain '*': the API authenticates via session cookies (CORS AllowCredentials), " +
					"which browsers and fiber/cors both refuse to combine with a wildcard origin — list explicit origins instead",
			)
		}
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
