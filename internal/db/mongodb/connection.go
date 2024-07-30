package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/Build-D-An-Ki-n-Truc/auth/internal/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client is the MongoDB client object
var Client *mongo.Client

// AdminColl is the collection object for the admin collection
var AdminColl *mongo.Collection

// UserColl is the collection object for the user collection
var UserColl *mongo.Collection

// BrandColl is the collection object for the brand collection
var BrandColl *mongo.Collection

// Initialize a connection to MongoDB
func InitializeMongoDBClient() error {
	// Set up a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure the context is canceled after the function returns

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.DbUrl))
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Test the connection
	var result bson.M
	if err := client.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	// Set the global client object
	Client = client
	AdminColl = Client.Database("BuildUserDB").Collection("admin")
	UserColl = Client.Database("BuildUserDB").Collection("user")
	BrandColl = Client.Database("BuildUserDB").Collection("brand")

	return nil

}

// Disconnect from MongoDB
func DisconnectMongoDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel() // Ensure the context is canceled after the function returns

	// Disconnect the MongoDB client
	err := Client.Disconnect(ctx)

	if err != nil {
		return fmt.Errorf("failed to disconnect from MongoDB: %v", err)
	}

	return nil
}
