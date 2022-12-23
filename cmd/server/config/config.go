package config

import (
	"errors"
	"fmt"
	"time"

	"workshop/pkg/logger"

	"github.com/ardanlabs/conf/v3"
)

type Help string

func (h Help) String() string {
	return string(h)
}

type Config struct {
	GracefullTimeout time.Duration `conf:"default:30s"`

	Log  logger.Config
	HTTP HTTP
	GRPC GRPC
	DB   Postgres
}

type Postgres struct {
	DSN string `conf:"required,env:DATABASE_DSN"`

	MaxIdleConns    int           `conf:"default:1"`
	MaxOpenConns    int           `conf:"default:5"`
	ConnMaxLifetime time.Duration `conf:"default:5s"`
	ConnMaxIdleTime time.Duration `conf:"default:5s"`
}

type HTTP struct {
	Addr         string        `conf:"default::8080"`
	ReadTimeout  time.Duration `conf:"default:1s"`
	WriteTimeout time.Duration `conf:"default:1s"`
	IdleTimeout  time.Duration `conf:"default:5s"`
}

type GRPC struct {
	Addr string `conf:"default::50051"`
}

func New() (Config, Help, error) {
	cfg := Config{}

	if help, err := conf.Parse("", &cfg); err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			return Config{}, Help(help), err
		}
		return Config{}, "", fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, "", nil
}
