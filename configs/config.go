package configs

import "github.com/LiquidCats/graceful"

type Config struct {
	App     AppConfig `envconfig:"APP" yaml:"app"`
	HTTP    graceful.HttpConfig
	Metrics graceful.HttpConfig
	Redis   RedisConfig   `envconfig:"REDIS" yaml:"redis"`
	Workers WorkersConfig `envconfig:"WORKERS" yaml:"workers"`
}
