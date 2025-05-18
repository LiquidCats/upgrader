package configs

type Config struct {
	App    AppConfig `envconfig:"APP" yaml:"app"`
	Redis  Redis     `envconfig:"REDIS" yaml:"redis"`
	Chains Chains    `envconfig:"CHAINS" yaml:"chains"`
}
