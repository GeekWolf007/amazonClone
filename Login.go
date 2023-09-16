package main

import (
	"context"
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

func Login(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var requestBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusInternalServerError)
		return
	}

	password := requestBody["password"].(string)
	username := requestBody["username"].(string)

	if username == "" {
		http.Error(w, "Username is missing!", http.StatusBadRequest)
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
