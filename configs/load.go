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

	decoder := yaml.NewDecoder(file)
	if err = decoder.Decode(&cfg); err != nil { //nolint:musttag
		return cfg, err
	}

	return cfg, nil
}
