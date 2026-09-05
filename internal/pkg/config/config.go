package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	PostgresDSN string   `env:"CONN_STRING,required"`
	RedisAddr   string   `env:"REDIS_CONN,required"`
	KafkaAddrs  []string `env:"KAFKA_CONN,required"`
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &cfg, nil
}
