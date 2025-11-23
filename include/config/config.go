package config

import (
	"fmt"

	"github.com/RomanKovalev007/pull_request_service/include/repository"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port string `env:"PORT" env-default:"8080"`
	Timeout string `env:"TIMEOUT" env-default:"30s"`
	BaseURL string `env:"BASE_URL" env-default:"http://localhost:8080"`
	Migration_Path string `env:"MIGRATION_PATH" env-default:"file:///migrations"`

	repository.Config
}

func ParseConfigFromEnv() (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config from env: %w", err)
	}

	return cfg, nil
}