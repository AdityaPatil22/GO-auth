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

var MongoClient *mongo.Client

func ConnectDB() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("❌ MONGO_URI not set in environment")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("❌ Failed to connect to MongoDB:", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ MongoDB ping failed:", err)
	}

	MongoClient = client
	fmt.Println("✅ Connected to MongoDB!")
}
