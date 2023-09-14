package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

func deleteUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header["Token"] == nil {
		var err Error
		err = SetError(err, "No Token Found")
		json.NewEncoder(w).Encode(err)
		return
	}

	token := r.Header["Token"]
	username, err := ExtractUsernameFromJWT(token[0])
	if err != nil {
		http.Error(w, "Error extracting username from JWT", http.StatusInternalServerError)
	}

	collection := client.Database("amazon_db").Collection("users")

	var filter bson.M
	var user User

	if username != "" {
		filter = bson.M{"username": username}
		error_name := collection.FindOne(context.TODO(), filter).Decode(&user)
		if error_name != nil {
			http.Error(w, "Username is not registered", http.StatusBadRequest)
			return
		}
	}

	deleteFilter := bson.M{"username": user.Username}
	_, deleteErr := collection.DeleteOne(context.TODO(), deleteFilter)

	if deleteErr != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("User %s deleted successfully", user.Username)))

}
