package common

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

// ConnectDB initializes and returns a MongoDB connection
func ConnectDB(uri, dbName string) *mongo.Database {
	clientOptions := options.Client().ApplyURI(uri)

	// Set a timeout for the connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the database to ensure connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("MongoDB Ping failed: %v", err)
	}

	log.Println("Connected to MongoDB successfully!")
	DB = client.Database(dbName) // Set the global variable for reuse
	return DB
}

func GetCollection(table string) (*mongo.Collection, context.Context) {
	collection := DB.Collection(table)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// Ensure the context gets canceled when no longer needed
	go func() {
		<-ctx.Done()
		cancel()
	}()

	return collection, ctx
}
