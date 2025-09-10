package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ConnectDB establishes a connection to MongoDB.
func ConnectDB() (*mongo.Client, error) {
	mongoURI := os.Getenv("MONGO_URI")
	// if mongoURI == "" {
	// 	mongoURI = "mongodb://localhost:27017"
	// 	log.Printf("MONGO_URI not set, using default: %s", mongoURI)
	// }

	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	return client, nil
}

// GetCollection returns a specific MongoDB collection.
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		dbName = "donation_platform"
		log.Printf("MONGO_DB_NAME not set, using default: %s", dbName)
	}
	return client.Database(dbName).Collection(collectionName)
}
