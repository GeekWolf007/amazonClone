package main

import (
	"context"
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

func ShowAllUsers(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
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

	if !user.IsAdmin {
		http.Error(w, "You are not an admin", http.StatusBadRequest)
		return
	}

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var users []User
	if err := cursor.All(context.TODO(), &users); err != nil {
		http.Error(w, "Error decoding users", http.StatusInternalServerError)
		return
	}

	jsonUsers, err := json.Marshal(users)
	if err != nil {
		http.Error(w, "Error encoding users to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonUsers)
}
