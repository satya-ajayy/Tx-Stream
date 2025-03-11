package transactions

import (
	// Go Internal Packages
	"context"
	"encoding/json"
	"fmt"

	models "go-kafka/models"

	// External Packages
	"go.uber.org/zap"
)

type TxRepository interface {
	InsertTransaction(ctx context.Context, transaction models.Transaction) error
}

type TxProcessor struct {
	Logger *zap.Logger
	TxRepo TxRepository
}

func NewTxProcessor(logger *zap.Logger, txRepo TxRepository) *TxProcessor {
	return &TxProcessor{TxRepo: txRepo, Logger: logger}
}

func (p *TxProcessor) ProcessRecords(records []models.Record) error {
	var tx models.Transaction
	for _, record := range records {
		err := json.Unmarshal(record.Value, &tx)
		if err != nil {
			p.Logger.Error("failed to unmarshal transaction", zap.Error(err))
			return err
		}
		println(fmt.Sprintf("Received Transaction %f", tx.Amount))
	}
	return nil
}
