package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io"
	"net/http"
)

type ResponseWrapper struct {
	Response     Response `json:"response"`
	ResponseType string   `json:"responseType"`
}

type Response struct {
	Name              string   `json:"name"`
	HasPremium        bool     `json:"hasPremium"`
	GearScore         int      `json:"GearScore"`
	Deaths            int      `json:"deaths"`
	Kills             int      `json:"kills"`
	EarnedCrystals    int      `json:"earnedCrystals"`
	CaughtGolds       int      `json:"caughtGolds"`
	TurretsPlayed     []Item   `json:"turretsPlayed"`
	HullsPlayed       []Item   `json:"hullsPlayed"`
	DronesPlayed      []Item   `json:"dronesPlayed"`
	PaintsPlayed      []Item   `json:"paintsPlayed"`
	ResistanceModules []Item   `json:"resistanceModules"`
	ModesPlayed       []Item   `json:"modesPlayed"`
	PreviousRating    Rating   `json:"previousRating"`
	Rating            Rating   `json:"rating"`
	Rank              int      `json:"rank"`
	Score             int      `json:"score"`
	ScoreBase         int      `json:"scoreBase"`
	ScoreNext         int      `json:"scoreNext"`
	SuppliesUsage     []Supply `json:"suppliesUsage"`
	//Presents        []Present`json:"presents"`
}

type Item struct {
	Grade       int    `json:"grade"`
	Id          int    `json:"id"`
	Name        string `json:"name"`
	ScoreEarned int    `json:"scoreEarned"`
	TimePlayed  int    `json:"timePlayed"`
	Type        string `json:"type"` // for modes
	//ImageUrl  string `json:"imageUrl"`
}

type Rating struct {
	Crystals   Info `json:"crystals"`
	Efficiency Info `json:"efficiency"`
	Golds      Info `json:"golds"`
	Score      Info `json:"score"`
}

type Info struct {
	Position int `json:"position"`
	Value    int `json:"value"`
}

type Supply struct {
	Name   string `json:"name"`
	Usages int    `json:"usages"`
}

func main() {
	username := "silly" // Replace with the desired username
	url := fmt.Sprintf("https://ratings.tankionline.com/api/eu/profile/?user=%s&lang=en", username)

	// Send GET request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	var data ResponseWrapper
	// Parse JSON into struct
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	fmt.Printf("%d\n", data.Response.Score)
}
