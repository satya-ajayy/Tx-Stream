package config

import "go-kafka/errors"

var DefaultConfig = []byte(`
logger:
	level: "info"

listen: ":8888"

prefix: "/ajayverse"

is_prod_mode: false

mongo:
  uri: "mongodb://localhost:27017"

kafka:
  uri: "localhost:6379"
  password: ""
`)

type Config struct {
	Logger     Logger `koanf:"logger"`
	Listen     string `koanf:"listen"`
	Prefix     string `koanf:"prefix"`
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
	URI string `koanf:"uri"`
	PWD string `koanf:"password"`
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
	if c.Kafka.URI == "" {
		ve.Add("redis.uri", "cannot be empty")
	}

	return ve.Err()
}
