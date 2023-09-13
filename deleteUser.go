package main

import (
	"context"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

func deleteUser(w http.ResponseWriter, r *http.Request) {

	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")
	login_token := r.FormValue("login_token")

	expectedKeysToDelete := []string{"email", "password", "username", "login_token"}

	for key := range r.Form {
		if !contains(expectedKeysToDelete, key) {
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

	if login_token == "" {
		http.Error(w, "Login token is missing", http.StatusBadRequest)
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

	if login_token != user.LoginToken {
		http.Error(w, "Incorrect login token", http.StatusBadRequest)
		return
	}

	fmt.Println("User to be deleted: ", user)
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
