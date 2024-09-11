package config

import (
	"os"
)

var (
	MongoURI string
)

func LoadConfig() {
	MongoURI = os.Getenv("MONGO_URI")
}
