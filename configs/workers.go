package configs

type WorkersConfig []WorkerConfig

type WorkerConfig struct {
	FromTopic   string `yaml:"from_topic"`
	ToWebsocket string `yaml:"to_websocket"`
}
