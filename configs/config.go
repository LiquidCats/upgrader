package configs

type Config struct {
	App     AppConfig `envconfig:"APP" yaml:"app"`
	Redis   Redis     `envconfig:"REDIS" yaml:"redis"`
	Workers Workers   `envconfig:"WORKERS" yaml:"workers"`
}
