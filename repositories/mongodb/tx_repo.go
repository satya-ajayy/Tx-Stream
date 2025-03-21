package mongodb

import (
	// Go Internal Packages
	"context"

	// Local Packages
	models "go-kafka/models"

	// External Packages
	"go.mongodb.org/mongo-driver/mongo"
)

type TxRepository struct {
	Client     *mongo.Client
	Collection string
}

func NewTxRepository(client *mongo.Client) *TxRepository {
	return &TxRepository{Client: client, Collection: "transactions"}
}

// InsertTransaction inserts a single transaction into database
func (r *TxRepository) InsertTransaction(ctx context.Context, tx models.MongoTransaction) error {
	collection := r.Client.Database("mybase").Collection(r.Collection)
	_, err := collection.InsertOne(ctx, tx)
	if err != nil {
		return err
	}
	return nil
}

// InsertTransactions inserts a batch of transactions into database
func (r *TxRepository) InsertTransactions(ctx context.Context, txs []interface{}) error {
	collection := r.Client.Database("mybase").Collection(r.Collection)
	_, err := collection.InsertMany(ctx, txs)
	if err != nil {
		return err
	}
	return nil
}
