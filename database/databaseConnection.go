package database

import (
	"context"
	"fmt"
	"log"
	"time"
	
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client = DBinstance()

func DBinstance() *mongo.Client {
	MongoDB := "mongodb://localhost:27017" // local db, otherwise paste provided URI; Depending on config you may write "user@password" before the localhost
	fmt.Println(MongoDB)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MongoDB))

	if err != nil {
		log.Fatal(err)
	}

	defer cancel()

	fmt.Println("Connected to MongoDB")
	
	return client
}

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("restaurant").Collection(collectionName)

	return collection
}