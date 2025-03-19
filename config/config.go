package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		Listen  `yaml:"listen"`
		Log     `yaml:"log"`
		Storage `yaml:"storage"`
	}

	// HTTP -.
	Listen struct {
		Host string `env-required:"true" yaml:"host" env:"LISTEN_HOST"`
		Port string `env-required:"true" yaml:"port" env:"LISTEN_PORT"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	// Storage -.
	Storage struct {
		// Type string
		Path string `yaml:"path" env:"STORAGE_PATH"`
	}
)

// NewConfig returns app config.
func NewConfig(path string) (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig(path, cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
