package configs

type Workers []WorkerConfig

type WorkerConfig struct {
	FromTopic   string `yaml:"from_topic"`
	ToWebsocket string `json:"to_websocket"`
}
