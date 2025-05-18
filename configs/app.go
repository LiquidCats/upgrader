package configs

type AppConfig struct {
	Port string `envconfig:"PORT" yaml:"port" default:"8080"`
}
