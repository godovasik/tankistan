package main

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

//var userList = []string{"silly", "icy", "Tobby1", "roman73_ul", "Bahkunkul", "Kleo4", "marcus7004", "KPECTbI4"}

var (
	client        *mongo.Client
	database      *mongo.Database
	timestampColl *mongo.Collection
	userColl      *mongo.Collection
)

func main() {

	if err := initMongoDB(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/newUser", makeNewUserHandler(timestampColl, userColl))
	port := ":8080"
	fmt.Printf("Server is listening on port %s\n", port)
	go func() {
		fmt.Println("Starting server...")
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Println("Server error:", err)
			return
		}
	}()

	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), 0, 10, 0, 0, now.Location()).Add(24 * time.Hour)
		duration := next.Sub(now)
		fmt.Println("sleepy joe:", duration)
		time.Sleep(duration)

		if err := ensureConnected(); err != nil {
			log.Printf("Failed to ensure connection: %v", err)
			continue
		}

		formattedTime := time.Now().Format("2006/01/02 15:04:05")
		fmt.Println(formattedTime, " Starting update process...")
		users, err := listOfUsersFromDB()
		if err != nil {
			fmt.Println("cant read users")
			return
		}
		for _, username := range users {
			err := sendReqAndUpdateUser(username)
			if err != nil {
				fmt.Println(err)
			}
			time.Sleep(3 * time.Second)
		}
		fmt.Println("im done!")
		//time.Sleep(1 * time.Minute)

	}
}
