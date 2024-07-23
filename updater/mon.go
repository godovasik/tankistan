package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"time"
)

func initMongoDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	password := os.Getenv("MONGO_PASS")
	mongoUrl := fmt.Sprintf("mongodb+srv://zamplt:%s@cluster0.3lwdyvn.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0", password)
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoUrl))
	if err != nil {
		return err
	}

	// Ping the database to verify connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}

	database = client.Database("tankistan")
	timestampColl = database.Collection("user_data")
	userColl = database.Collection("users")

	log.Println("Connected to MongoDB")
	return nil
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

func ensureConnected() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Printf("Lost connection to MongoDB. Attempting to reconnect...")
		if err := initMongoDB(); err != nil {
			return fmt.Errorf("failed to reconnect: %v", err)
		}
		log.Println("Reconnected successfully")
	}

	return nil
}

func insertDatastamp(data Datastamp) error {
	doesEx, err := doesTimestampExist(timestampColl, data)
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

// this refreshes user even all is same, will fix later
func updateUser(data Datastamp) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	exists, err := doesUserExist(userColl, data.Name)
	//fmt.Println("93line, exists:", exists)
	if exists {
		filter := bson.M{"name": data.Name}
		update := bson.M{
			"$set": bson.M{
				"lastupdate": time.Now().Truncate(24 * time.Hour),
			},
		}
		_, err = userColl.UpdateOne(context.TODO(), filter, update)
		if err != nil {
			return err
		}
		//fmt.Println(data.Name, "already exists, lastupdate updated")
		return nil
	}
	user := User{data.Name, data.Timestamp, true}
	_, err = userColl.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	fmt.Printf("User %s succesfully added to db\n", data.Name)
	return nil
}

func updateUserData(data Datastamp) error {
	err := updateUser(data)
	if err != nil {
		return err
	}
	err = insertDatastamp(data)
	if err != nil {
		return err
	}
	return nil
}

func sendReqAndUpdateUser(username string) error {
	resp, err := sendRequest(username)
	if err != nil {
		fmt.Println("fuck you")
		return err
	}

	data, err := parse(resp)
	if err != nil {
		fmt.Println("fuck you second time: ", err)
		return err
	}
	resp.Body.Close()

	var compact Datastamp

	compact.Store(data)

	err = updateUserData(compact)
	return err

}

func listOfUsersFromDB() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	projection := bson.M{"name": 1, "_id": 0}

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
