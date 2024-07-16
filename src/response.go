package main

import "time"

var (
	HULLS []string = ["Wasp", "Hopper", "Hornet", "Viking"]
	
)

type ResponseWrapper struct {
	Response     Response `json:"response"`
	ResponseType string   `json:"responseType"`
}

type Response struct {
	Name              string   `json:"name"`       // !
	HasPremium        bool     `json:"hasPremium"` //!
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

// --- data to save in db ---

type Datastamp struct {
	Timestamp      time.Time
	Name           string
	Rank           int
	Kills          int
	Deaths         int
	EarnedCrystals int
	GearScore             int
	Hulls          []Thing
	Turrets        []Thing
	Drones         []Thing
	Supplies       []Supply
}

type Thing struct {
	name        string
	scoreEarned int
	timePlayed  int
}

func (d *Datastamp) Store(data ResponseWrapper) {
	r := data.Response
	d.Name, d.Rank, d.Kills, d.Deaths, d.EarnedCrystals, d.GearScore =
		r.Name, r.Rank, r.Kills, r.Deaths, r.EarnedCrystals, r.GearScore
}
