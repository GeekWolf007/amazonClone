package main

import (
	"context"
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

func ShowAllUsers(w http.ResponseWriter, r *http.Request) {

	collection := client.Database("amazon_db").Collection("users")

	user := r.Context().Value(userKey).(User)

	if !user.IsAdmin {
		http.Error(w, "You are not an admin", http.StatusBadRequest)
		return
	}

	users_list, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Error fetching users", http.StatusInternalServerError)
		return
	}
	defer users_list.Close(context.TODO())

	var users []User
	if err := users_list.All(context.TODO(), &users); err != nil {
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
