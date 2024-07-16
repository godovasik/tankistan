package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"os"
)

func parse(resp *http.Response) (ResponseWrapper, error) {
	var data ResponseWrapper
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	err = json.Unmarshal(body, &data)

	return data, err
}

func sendRequest(username string) (*http.Response, error) {
	url := fmt.Sprintf("https://ratings.tankionline.com/api/eu/profile/?user=%s&lang=en", username)

	// Send GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making request:", err)
		return nil, err
	}

	return resp, err
}

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

	fmt.Printf("username: %s\ngs: %d\nk/d = %.2f\n",
		data.Response.Name, data.Response.GearScore, float32(data.Response.Kills)/float32(data.Response.Deaths))

}
