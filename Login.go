package main

import (
	"context"
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.URL.Query().Get("username")
	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")

	expectedKeysToLogin := []string{"email", "password", "username"}

	for key := range r.Form {
		if !contains(expectedKeysToLogin, key) {
			http.Error(w, "Unexpected key in form data: "+key, http.StatusBadRequest)
			return
		}
	}

	if username == "" && email == "" {
		http.Error(w, "Either username or email is required", http.StatusBadRequest)
		return
	}

	if username != "" && email != "" {
		http.Error(w, "Enter either username or email", http.StatusBadRequest)
		return
	}

	if password == "" {
		http.Error(w, "Password field is missing", http.StatusBadRequest)
		return
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
	if email != "" {
		filter = bson.M{"email": email}
		error_email := collection.FindOne(context.Background(), filter).Decode(&user)
		if error_email != nil {
			http.Error(w, "Email is not registered", http.StatusBadRequest)
			return
		}
	}
	if password != user.Password {
		http.Error(w, "Incorrect password", http.StatusBadRequest)
		return
	}

	tokenString, err := GenerateJWT(username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message":  "Login successful",
		"token":    tokenString,
		"is_admin": user.IsAdmin,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error creating JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	// fmt.Println(user)
}
