package db

import (
	"context"
	"library_management_system/config/dbconfig"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var BooksCollection *mongo.Collection
var UsersCollection *mongo.Collection

// InitDB initializes the MongoDB client and collections.
func InitDB() {
	clientOptions := options.Client().ApplyURI(dbconfig.MongoURI)
	var err error
	dbContext := context.Background()
	Client, err = mongo.Connect(dbContext, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = Client.Ping(dbContext, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize collections
	BooksCollection = Client.Database(dbconfig.DatabaseName).Collection(dbconfig.BooksCollection)
	UsersCollection = Client.Database(dbconfig.DatabaseName).Collection(dbconfig.UsersCollection)

	// Create username index for UsersCollection
	_, err = UsersCollection.Indexes().CreateOne(dbContext, mongo.IndexModel{
		Keys:    bson.D{{Key: dbconfig.Username, Value: 1}}, // Index on 'username'
		Options: options.Index().SetUnique(true),            // Ensure unique values
	})
	if err != nil {
		log.Fatalf("Failed to create index: %v", err)
	}
}
