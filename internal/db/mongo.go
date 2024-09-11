package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	// "os"
	"tankistan/config"
	"tankistan/internal/models"
	"time"
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

func GetCollection(databaseName, collectionName string) *mongo.Collection {
	return Client.Database(databaseName).Collection(collectionName)
}

func InsertDatastamp(data models.Datastamp) error {
	timestampColl := GetCollection("tankistan", "user_data")
	doesEx, err := DoesTimestampExist(data)
	if err != nil {
		return err
	}
	if doesEx {
		fmt.Printf("timestamp for %s already exists\n", data.Name)
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_, err = timestampColl.InsertOne(ctx, data)
		if err != nil {
			return err
		}
		fmt.Println("datastamp for", data.Name, "updated")
	}
	return err
}

func DoesTimestampExist(data models.Datastamp) (bool, error) {
	collection := GetCollection("tankistan", "user_data")
	filter := bson.M{"name": data.Name, "timestamp": data.Timestamp}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func DoesUserExist(name string) (bool, error) {
	collection := GetCollection("tankistan", "users")
	filter := bson.M{"name": name}
	count, err := collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func UpdateUser(name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	userColl := GetCollection("tankistan", "users")
	exists, err := DoesUserExist(name)
	//fmt.Println("93line, exists:", exists)

	timeNowTrunc := time.Now().Truncate(24 * time.Hour)

	if exists {
		filter := bson.M{"name": name}
		update := bson.M{
			"$set": bson.M{
				"lastupdate": timeNowTrunc,
			},
		}
		_, err = userColl.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return err
		}
		//fmt.Println(name, "already exists, lastupdate updated")
		return nil
	}
	user := models.User{Name: name, LastUpdate: timeNowTrunc, KeepTrack: true}
	_, err = userColl.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	fmt.Printf("User %s succesfully added to db\n", name)
	return nil
}

func UpdateUserData(data models.Datastamp) error {
	err := UpdateUser(data.Name)
	if err != nil {
		return err
	}
	err = InsertDatastamp(data)
	if err != nil {
		return err
	}
	return nil
}

func ListOfUsersFromDB() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	projection := bson.M{"name": 1, "_id": 0}
	userColl := GetCollection("tankistan", "users")

	cursor, err := userColl.Find(ctx, bson.M{}, options.Find().SetProjection(projection))
	if err != nil {
		fmt.Println("amogus")
	}
	defer cursor.Close(context.TODO())

	var users []string

	for cursor.Next(context.TODO()) {
		var result struct {
			Name string `bson:"name"`
		}
		err := cursor.Decode(&result)
		if err != nil {
			return nil, err
		}
		users = append(users, result.Name)
	}

	if err := cursor.Err(); err != nil {
		log.Fatal(err)
	}
	return users, nil
}

// need to add sendReqAndUpdateUser functon
