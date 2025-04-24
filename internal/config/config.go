package config

import (
	"log"

	"github.com/caarlos0/env/v6"
	"github.com/literaen/simple_project/pkg/config"
)

type Config struct {
	DB config.DB_CREDS

	REDIS config.REDIS_CREDS

	PORT string `env:"PORT"`

	GRPC_Port string `env:"GRPC_PORT"`

	TASK_SERVICE_HOST string `env:"TASK_SERVICE_HOST"`
	TASK_SERVICE_PORT string `env:"TASK_SERVICE_PORT"`

	KAFKA_BROKERS []string `env:"KAFKA_BROKERS"`
}

func ProvideDBCreds(cfg *Config) *config.DB_CREDS {
	return &cfg.DB
}

func ProvideRedisCreds(cfg *Config) *config.REDIS_CREDS {
	return &cfg.REDIS
}

func LoadEnv() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Failed to load env variables: %v", err)
	}
	return cfg
}
