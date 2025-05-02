package mongodb

import (
	// Go Internal Packages
	"context"

	// Local Packages
	models "tx-stream/models"

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

// InsertTransaction inserts a single transaction into the database
func (r *TxRepository) InsertTransaction(ctx context.Context, tx models.MongoTransaction) error {
	collection := r.client.Database("mybase").Collection(r.collection)
	_, err := collection.InsertOne(ctx, tx)
	if err != nil {
		return err
	}
	return nil
}

// InsertTransactions inserts a batch of transactions into the database
func (r *TxRepository) InsertTransactions(ctx context.Context, txs []interface{}) error {
	collection := r.client.Database("mybase").Collection(r.collection)
	_, err := collection.InsertMany(ctx, txs)
	if err != nil {
		return err
	}
	return nil
}
