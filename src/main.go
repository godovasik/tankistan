package main

import (
	"fmt"
	"os"
	"time"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("No arguments provided")
		return
	}
	username := os.Args[1] // Replace with the desired username

	initMongoDB()
	database := client.Database("tankistan")
	timestampColl := database.Collection("user_data")
	userColl := database.Collection("users")
	defer closeMongoDB()

	c := make(chan os.Signal, 1)
	go func() {
		<-c
		fmt.Println("\nCleaning up...")
		closeMongoDB()
		os.Exit(0)
	}()

	for true {
		resp, err := sendRequest(username)
		if err != nil {
			fmt.Println("fuck you")
			return
		}

		data, err := parse(resp)
		if err != nil {
			fmt.Println("fuck you second time: ", err)
			return
		}
		resp.Body.Close()

		var compact Datastamp

		compact.Store(data)

		err = updateUserData(timestampColl, userColl, compact)
		if err != nil {
			fmt.Println(err)
		}

		time.Sleep(1 * time.Hour)
	}

}
