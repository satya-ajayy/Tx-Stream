package main

import (
	// Go Internal Packages
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	// Local Packages
	config "go-kafka/config"
	kafka "go-kafka/kafka"
	mongodb "go-kafka/repositories/mongodb"
	txsvc "go-kafka/services/transactions"

	// External Packages
	"github.com/alecthomas/kingpin/v2"
	_ "github.com/jsternberg/zap-logfmt"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/rawbytes"
	"github.com/twmb/franz-go/plugin/kprom"
	"go.uber.org/zap"
)

// LoadSecrets Loads the secret variables and overrides the config
func LoadSecrets(k config.Config) config.Config {
	MongoURI := os.Getenv("MONGO_URI")
	if MongoURI != "" {
		k.Mongo.URI = MongoURI
	}

	KafkaBrokers := os.Getenv("KAFKA_BROKERS")
	if KafkaBrokers != "" {
		k.Kafka.Brokers = KafkaBrokers
	}

	IsProdMode := os.Getenv("IS_PROD_MODE")
	k.IsProdMode = IsProdMode == "true"
	return k
}

// LoadConfig loads the default configuration and overrides it with the config file
// specified by the path defined in the config flag
func LoadConfig() *koanf.Koanf {
	configPathMsg := "Path to the application config file"
	configPath := kingpin.Flag("config", configPathMsg).Short('c').Default("config.yml").String()

	kingpin.Parse()
	k := koanf.New(".")
	_ = k.Load(rawbytes.Provider(config.DefaultConfig), yaml.Parser())
	if *configPath != "" {
		_ = k.Load(file.Provider(*configPath), yaml.Parser())
	}
	return k
}

func main() {
	k := LoadConfig()
	appKonf := config.Config{}

	// Unmarshalling config into struct
	err := k.Unmarshal("", &appKonf)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Update and Validate config before starting the server
	updatedKonf := LoadSecrets(appKonf)
	if err = updatedKonf.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	if !updatedKonf.IsProdMode {
		k.Print()
	}

	cfg := zap.NewProductionConfig()
	cfg.Encoding = "logfmt"
	_ = cfg.Level.UnmarshalText([]byte(k.String("logger.level")))
	cfg.InitialFields = make(map[string]any)
	cfg.InitialFields["host"], _ = os.Hostname()
	cfg.InitialFields["service"] = "txconsumer"
	cfg.OutputPaths = []string{"stdout"}
	logger, _ := cfg.Build()
	defer func() {
		_ = logger.Sync()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Connect to mongodb
	mongoClient, err := mongodb.Connect(ctx, updatedKonf.Mongo.URI)
	if err != nil {
		logger.Fatal("cannot create mongo client", zap.Error(err))
	}

	txRepo := mongodb.NewTxRepository(mongoClient)
	txProcessor := txsvc.NewTxProcessor(logger, txRepo)

	metrics := kprom.NewMetrics("et")
	conf := &kafka.ConsumerConfig{
		Brokers:        []string{updatedKonf.Kafka.Brokers},
		Name:           updatedKonf.Kafka.ConsumerName,
		Topic:          updatedKonf.Kafka.Topic,
		RecordsPerPoll: 5,
	}

	txConsumer, err := kafka.NewTxConsumer(conf, txProcessor, metrics, logger)
	if err != nil {
		logger.Fatal("cannot create transactions consumer", zap.Error(err))
	}

	if txConsumer != nil {
		txConsumer.Poll(ctx)
	}
}
