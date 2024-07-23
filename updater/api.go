package main

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
)

type Username struct {
	Username string `json:"username"`
}

func makeNewUserHandler(timestampColl, userColl *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var user Username
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Received new user: %s", user.Username)
		err = updateUser(user.Username)
	}
}
