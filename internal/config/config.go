package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	HTTPServer HTTPServer `yaml:"http_server"`
	Log        Log        `yaml:"logger"`
}

type HTTPServer struct {
	Port        string        `yaml:"port" default:":8080"`
	Timeout     time.Duration `yaml:"timeout" default:"10s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" default:"60s"`
}

type Log struct {
	Level string `yaml:"level" default:"info"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
