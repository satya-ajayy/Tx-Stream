package redis

import (
	// Go Internal Packages
	"context"
	"encoding/json"
	"fmt"

	// Local Packages
	models "tx-stream/models"

	// External Packages
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type DeadLetterQueue struct {
	client   *redis.Client
	logger   *zap.Logger
	listName string
}

func NewDeadLetterQueue(client *redis.Client, logger *zap.Logger) *DeadLetterQueue {
	return &DeadLetterQueue{client: client, logger: logger, listName: "failed-transactions"}
}

// Send stores all failed records into Redis with the key as "tx:{transaction_id}"
func (r *DeadLetterQueue) Send(ctx context.Context, records []models.Record) error {
	if len(records) == 0 {
		return nil
	}

	successCount := 0
	for _, record := range records {
		jsonData, err := json.Marshal(record)
		if err != nil {
			r.logger.Error("failed to marshal record", zap.Error(err))
			continue
		}

		key := fmt.Sprintf("tx:%s", record.Key)
		err = r.client.Set(ctx, key, jsonData, 0).Err()
		if err != nil {
			r.logger.Error("failed to store record", zap.String("key", key), zap.Error(err))
			continue
		}
		successCount++
	}

	if successCount > 0 {
		r.logger.Info("successfully sent records", zap.Int("count", successCount))
	}

	return nil
}
