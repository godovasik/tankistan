package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

var client *mongo.Client

func initMongoDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	password := os.Getenv("MONGO_PASS")
	mongoUrl := fmt.Sprintf("mongodb+srv://zamplt:%s@cluster0.3lwdyvn.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0", password)
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))
	if err != nil {
		log.Fatal(err)
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB")
}

func closeMongoDB() {
	if client != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := client.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
		log.Println("Disconnected from MongoDB")
	}
}

func insertDatastamp(collection *mongo.Collection, data Datastamp) error {
	doesEx, err := doesTimestampExist(collection, data)
	if err != nil {
		return err
	}
	if doesEx {
		fmt.Printf("timestamp for %s already exists\n", data.Name)
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err = collection.InsertOne(ctx, data)
		if err != nil {
			return err
		}
		fmt.Println(data.Name, " added to db")
	}
	return err
}

func doesTimestampExist(collection *mongo.Collection, data Datastamp) (bool, error) {
	filter := bson.M{"name": data.Name, "timestamp": data.Timestamp}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func doesUserExist(collection *mongo.Collection, name string) (bool, error) {
	filter := bson.M{"name": name}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func listAllUsers(collection *mongo.Collection) error {
	return nil
}

func updateUser(collection *mongo.Collection, data Datastamp) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	exists, err := doesUserExist(collection, data.Name)
	//fmt.Println("93line, exists:", exists)
	if exists {
		filter := bson.M{"name": data.Name}
		update := bson.M{
			"$set": bson.M{
				"lastupdate": time.Now().Truncate(time.Hour),
			},
		}
		_, err = collection.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return err
		}
		fmt.Println("User already exists, updated")
		return nil
	}
	user := User{data.Name, data.Timestamp, true}
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	fmt.Printf("User %s succesfully added to db\n", data.Name)
	return nil
}

func updateUserData(timestampColl, userColl *mongo.Collection, data Datastamp) error {

	err := insertDatastamp(timestampColl, data)
	if err != nil {
		return err
	}

	err = updateUser(userColl, data)
	return err
}
