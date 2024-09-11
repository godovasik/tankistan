package utils

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
	"tankistan/internal/models"
)

func parse(resp *http.Response) (models.ResponseWrapper, error) {
	var data models.ResponseWrapper
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
