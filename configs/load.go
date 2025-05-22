package configs

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

func Load(prefix string) (Config, error) {
	var cfg Config

	if err := envconfig.Process(prefix, &cfg); err != nil {
		return cfg, err
	}

	file, err := os.OpenFile(".app.cfg.yaml", os.O_RDONLY, 0666)
	if err != nil {
		return cfg, err
	}

	var workers Workers
	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&workers); err != nil {
		return cfg, err
	}

	cfg.Workers = workers

	return cfg, nil
}
