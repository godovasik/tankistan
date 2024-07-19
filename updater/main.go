package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

//var userList = []string{"silly", "icy", "Tobby1", "roman73_ul", "Bahkunkul", "Kleo4", "marcus7004", "KPECTbI4"}

func main() {

	//if len(os.Args) < 2 {
	//	fmt.Println("No arguments provided")
	//	return
	//}
	//username := os.Args[1] // Replace with the desired username

	initMongoDB()
	database := client.Database("tankistan")
	timestampColl := database.Collection("user_data")
	userColl := database.Collection("users")
	defer closeMongoDB()

	//fmt.Println(listOfUsersFromDB(userColl))

	stop := make(chan bool)
	// Goroutine to listen for user input
	go func() {
		var input string
		fmt.Scanln(&input)
		stop <- true
	}()

	for {
		select {
		case <-stop:
			fmt.Println("Program stopped by user input")
			os.Exit(0)
		default:

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

			formattedTime := time.Now().Format("2006/01/02 15:04:05")
			fmt.Println(formattedTime, " Starting update process...")
			users, err := listOfUsersFromDB(userColl)
			if err != nil {
				fmt.Println("cant read users")
				return
			}
			for _, username := range users {
				err := sendReqAndUpdateUser(timestampColl, userColl, username)
				if err != nil {
					fmt.Println(err)
				}
				time.Sleep(5 * time.Second)
			}
			fmt.Println("im done!")
			time.Sleep(12 * time.Hour)
		}
	}

}
