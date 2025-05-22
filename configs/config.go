package configs

import "github.com/LiquidCats/graceful"

type Config struct {
	App     AppConfig           `envconfig:"APP" yaml:"app"`
	HTTP    graceful.HttpConfig `envconfig:"HTTP" yaml:"http"`
	Redis   RedisConfig         `envconfig:"REDIS" yaml:"redis"`
	Workers WorkersConfig       `envconfig:"WORKERS" yaml:"workers"`
}
