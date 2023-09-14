package main

import (
	"context"
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

func ViewCart(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
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

	var filter_username bson.M
	var user User
	collectionUser := client.Database("amazon_db").Collection("users")

	if username != "" {
		filter_username = bson.M{"username": username}
		error_name := collectionUser.FindOne(context.TODO(), filter_username).Decode(&user)
		if error_name != nil {
			http.Error(w, "Username is not registered", http.StatusBadRequest)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.Cart)
}
