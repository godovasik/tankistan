package main

import (
	"fmt"
	"os"
	"time"
)

var userList = []string{"silly", "icy", "Tobby1", "roman73_ul", "Bahkunkul", "Kleo4", "marcus7004", "KPECTbI4"}

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
			formattedTime := time.Now().Format("2006/01/02 15:04:05")
			fmt.Println(formattedTime, " Starting update process...")
			for _, username := range userList {

				err := sendReqAndUpdateUser(timestampColl, userColl, username)
				if err != nil {
					fmt.Println(err)
				}
				time.Sleep(5 * time.Second)
			}
			fmt.Println("im done!")
			time.Sleep(3 * time.Hour)
		}
	}

}
