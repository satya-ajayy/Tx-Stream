package mongodb

import (
	// Go Internal Packages
	"context"

	// External Packages
	"go.mongodb.org/mongo-driver/mongo"
)

type TxRepository struct {
	client     *mongo.Client
	collection string
}

func NewTxRepository(client *mongo.Client) *TxRepository {
	return &TxRepository{client: client, collection: "transactions"}
}

// InsertTransactions inserts a batch of transactions into database
func (r *TxRepository) InsertTransactions(ctx context.Context, txs []interface{}) error {
	collection := r.client.Database("mybase").Collection(r.collection)
	_, err := collection.InsertMany(ctx, txs)
	if err != nil {
		return err
	}
	return nil
}
