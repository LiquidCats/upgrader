package configs

import "github.com/go-playground/sensitive"

type Redis struct {
	Address  string           `yaml:"address" envconfig:"ADDRESS"`
	Password sensitive.String `yaml:"password" envconfig:"PASSWORD"`
	DB       int              `yaml:"DB" envconfig:"DB"`
}
