package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

type Config struct {
	Port         int           `env:"PORT" envDefault:"5800" validate:"min=1,max=65535"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"30s"`
	ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"10s"`
	IdleTimeout  time.Duration `env:"IDLE_TIMEOUT" envDefault:"60s"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	if err := validator.New().Struct(&cfg); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	return &cfg, nil
}
