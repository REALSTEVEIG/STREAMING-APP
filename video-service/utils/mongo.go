package utils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB() (*mongo.Client, error) {
	// Get MongoDB URI from environment variables (default to localhost)
	mongoURI := GetEnv("MONGO_URI", "mongodb://localhost:27017")

	// Create a context with a timeout for connecting to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create client options using the provided URI
	clientOptions := options.Client().ApplyURI(mongoURI)

	// Connect to MongoDB using the context and client options
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client, nil
}
