package kafka

import (
	// Go Internal Packages
	"context"
	"errors"
	"fmt"

	// Local Packages
	models "go-kafka/models"

	// External Packages
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/plugin/kprom"
	"go.uber.org/zap"
)

type ConsumerConfig struct {
	Brokers        []string
	Name           string
	Topic          string
	RecordsPerPoll int
}

type Consumer struct {
	Client    *kgo.Client
	Config    *ConsumerConfig
	Processor TxProcessor
	Logger    *zap.Logger
}

type TxProcessor interface {
	ProcessRecords(records []models.Record) error
}

// NewTxConsumer creates a new consumer and starts a goroutine for each partition to consume
// the records that are fetched (PS: Must call Poll to start consuming the records)
func NewTxConsumer(conf *ConsumerConfig, processor TxProcessor, metrics *kprom.Metrics, logger *zap.Logger) (*Consumer, error) {
	c := &Consumer{Config: conf, Processor: processor, Logger: logger}

	opts := []kgo.Opt{
		kgo.SeedBrokers(conf.Brokers...), // Connects to Kafka brokers
		kgo.ConsumerGroup(conf.Name),     // Specifies the consumer group
		kgo.ConsumeTopics(conf.Topic),    // Specifies a single topic to consume
		kgo.WithHooks(metrics),           // Attaches monitoring hooks
		kgo.DisableAutoCommit(),          // Disables auto-commit
		kgo.BlockRebalanceOnPoll(),       // Blocks rebalancing until the poll loop is running
	}

	client, err := kgo.NewClient(opts...)
	if err != nil || client == nil {
		return nil, err
	}

	c.Client = client
	return c, nil
}

// Poll polls for records from the Kafka broker.
func (c *Consumer) Poll(ctx context.Context) error {
	defer c.Client.Close()

	consumerName := c.Config.Name
	recordsPerPoll := c.Config.RecordsPerPoll

	for {
		// Check if the context is canceled before polling
		if ctx.Err() != nil {
			c.Logger.Warn("Polling stopped: context canceled")
			return ctx.Err() // Exit gracefully
		}

		c.Logger.Info(fmt.Sprintf("%s: polling for records", consumerName))
		fetches := c.Client.PollRecords(ctx, recordsPerPoll)

		// Handle client shutdown
		if fetches.IsClientClosed() {
			return errors.New("kafka client closed")
		}

		// Handle context cancellation explicitly
		if errors.Is(fetches.Err0(), context.Canceled) {
			return errors.New("context got canceled")
		}

		// Preallocate records slice efficiently
		records := make([]models.Record, len(fetches.Records()))
		for idx, record := range fetches.Records() {
			records[idx] = models.Record{
				Key:   record.Key,
				Value: record.Value,
				Topic: record.Topic,
			}
		}

		// Process records and handle errors robustly
		if err := c.Processor.ProcessRecords(records); err != nil {
			c.Logger.Error("Failed to process records", zap.Error(err))
			continue // Don't exit on a single failure
		}

		// Commit processed records
		_ = c.Client.CommitRecords(ctx, fetches.Records()...)
	}
}
