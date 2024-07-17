package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No arguments provided")
		return
	}
	username := os.Args[1] // Replace with the desired username

	resp, err := sendRequest(username)
	if err != nil {
		fmt.Println("fuck you")
		return
	}
	defer resp.Body.Close()

	data, err := parse(resp)
	if err != nil {
		fmt.Println("fuck you second time: ", err)
		return
	}

	var compact Datastamp

	compact.Store(data)

	compact.Print()

}
