package main

import (
	// "log"
	// "net/http"
	"tankistan/config"
	"tankistan/internal/db"
)

func main() {
	config.LoadConfig()

	db.Connect()
	defer db.Disconnect()

}
