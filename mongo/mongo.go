package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	*mongo.Client
}

func Connect(clusterURL string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Client().ApplyURI(clusterURL)
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	db := &DB{
		Client: client,
	}
	return db, nil
}

func CloseConn(db *DB) error {
	return db.Client.Disconnect(context.TODO())
}
