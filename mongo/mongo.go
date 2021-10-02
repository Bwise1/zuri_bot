package mongo

import (
	"context"
	"log"
	"os"
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
		log.Printf("Error connecting to db: %v", err)
		return nil, err

	}
	db := &DB{
		Client: client,
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("cannot ping server %v", err)
	}
	log.Println("Database connected successfully")
	return db, nil
}

func CloseConn(db *DB) error {
	return db.Client.Disconnect(context.TODO())
}

func (db *DB) GetCollection(name string) *mongo.Collection {
	dbName := os.Getenv("DB_NAME")
	return db.Client.Database(dbName).Collection(name)
}
