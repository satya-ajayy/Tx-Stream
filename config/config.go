package config

import (
	// Local Packages
	errors "go-kafka/errors"
)

var DefaultConfig = []byte(`
logger:
  level: "info"

listen: ":8888"

is_prod_mode: false

mongo:
  uri: "mongodb://localhost:27017"

kafka:
  brokers: "localhost:9092"
  consume: true
  topic: "transactions-input"
  records_per_poll: 5
  consumer_name: "transactions-consumer"
`)

type Config struct {
	Logger     Logger `koanf:"logger"`
	Listen     string `koanf:"listen"`
	IsProdMode bool   `koanf:"is_prod_mode"`
	Mongo      Mongo  `koanf:"mongo"`
	Kafka      Kafka  `koanf:"kafka"`
}

type Logger struct {
	Level string `koanf:"level"`
}

type Mongo struct {
	URI string `koanf:"uri"`
}

type Kafka struct {
	Brokers        string `koanf:"brokers"`
	Consume        bool   `koanf:"consume"`
	Topic          string `koanf:"topic"`
	RecordsPerPoll int    `koanf:"records_per_poll"`
	ConsumerName   string `koanf:"consumer_name"`
}

// Validate validates the configuration
func (c *Config) Validate() error {
	ve := errors.ValidationErrs()

	if c.Listen == "" {
		ve.Add("listen", "cannot be empty")
	}
	if c.Logger.Level == "" {
		ve.Add("logger.level", "cannot be empty")
	}
	if c.Mongo.URI == "" {
		ve.Add("mongo.uri", "cannot be empty")
	}
	if c.Kafka.Brokers == "" {
		ve.Add("kafka.brokers", "cannot be empty")
	}

	return ve.Err()
}
