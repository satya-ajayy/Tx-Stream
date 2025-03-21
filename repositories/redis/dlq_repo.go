package redis

import (
	// Go Internal Packages
	"context"
	"encoding/json"
	"go.uber.org/zap"

	// Local Packages
	models "tx-stream/models"

	// External Packages
	"github.com/redis/go-redis/v9"
)

type DeadLetterQueue struct {
	Client   *redis.Client
	Logger   *zap.Logger
	ListName string
}

func NewDeadLetterQueue(client *redis.Client, logger *zap.Logger) *DeadLetterQueue {
	return &DeadLetterQueue{Client: client, Logger: logger, ListName: "failed-transactions"}
}

// Send pushes all failed records into the Redis list "failed-transactions"
func (r *DeadLetterQueue) Send(ctx context.Context, records []models.Record) error {
	if len(records) == 0 {
		return nil
	}

	var transactions []interface{}
	for _, record := range records {
		transaction, err := json.Marshal(record)
		if err != nil {
			r.Logger.Error("failed to marshal transaction", zap.Error(err))
			continue
		}
		transactions = append(transactions, transaction)
	}

	err := r.Client.LPush(ctx, r.ListName, transactions...).Err()
	if err != nil {
		return err
	}

	return nil
}
