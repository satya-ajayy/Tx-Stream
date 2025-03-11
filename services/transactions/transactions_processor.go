package transactions

import (
	// Go Internal Packages
	"context"

	// Local Packages
	helpers "go-kafka/helpers"
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

func (t *TxProcessor) ProcessRecords(records []models.Record) error {
	helpers.PrintStruct(records)
	return nil
}
