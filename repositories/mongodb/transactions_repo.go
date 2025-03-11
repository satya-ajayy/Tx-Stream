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
	client     *mongo.Client
	collection string
}

func NewTxRepository(client *mongo.Client) *TxRepository {
	return &TxRepository{client: client, collection: "transactions"}
}

// InsertTransaction inserts a new transaction into database
func (r *TxRepository) InsertTransaction(ctx context.Context, transaction models.Transaction) error {
	collection := r.client.Database("mybase").Collection(r.collection)
	_, err := collection.InsertOne(ctx, transaction)
	if err != nil {
		return err
	}
	return nil
}
