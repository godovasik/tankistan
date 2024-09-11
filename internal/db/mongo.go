package db

import (
	"context"
	// "fmt"
	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	// "os"
	"tankistan/config"
	// "time"
)

var Client *mongo.Client

func Connect() {
	clientOptions := options.Client().ApplyURI(config.MongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal(err)
	}
	Client = client
	log.Println("Connected to MongoDB!")
}

func Disconnect() {
	if err := Client.Disconnect(context.Background()); err != nil {
		log.Fatal(err)
	}
}
