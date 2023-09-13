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

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	username := r.URL.Query().Get("username")
	login_token := r.URL.Query().Get("login_token")

	expectedKeysToAddToCart := []string{"login_token", "username"}

	for key := range r.Form {
		if !contains(expectedKeysToAddToCart, key) {
			http.Error(w, "Unexpected key in form data: "+key, http.StatusBadRequest)
			return
		}
	}

	if username == "" {
		http.Error(w, "Username field is missing", http.StatusBadRequest)
		return
	}
	if login_token == "" {
		http.Error(w, "Login token is missing", http.StatusBadRequest)
		return
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
	if login_token != user.LoginToken {
		http.Error(w, "Incorrect login token", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.Cart)
}
