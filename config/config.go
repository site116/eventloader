package config


// Config represents the configuration for the Kafka producer and worker pool.
type Config struct {
	Broker `yaml:"broker"`
	Pool   `yaml:"pool"`
}

// Broker represents the configuration for the Kafka broker.
type Broker struct {
	Topic string   `yaml:"topic"`
	Seeds []string `yaml:"seeds"`
	SASL  struct {
		Enabled  bool   `yaml:"enabled"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"sasl"`
}

// Pool represents the configuration for the worker pool. 
type Pool struct {
	NumWorkers       int   `yaml:"num_workers"`
	MessagesPerBatch int32 `yaml:"messagesPerBatch"`
	ErrorThreshold   int32 `yaml:"errorThreshold"`
	Batches          int32 `yaml:"batches"`
}
