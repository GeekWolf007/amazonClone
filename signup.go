package main

import (
	"context"
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func Signup(w http.ResponseWriter, r *http.Request) {
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

	expectedKeysToSignup := []string{"email", "password", "username", "phone"}

	for key := range r.Form {
		if !contains(expectedKeysToSignup, key) {
			http.Error(w, "Unexpected key in form data: "+key, http.StatusBadRequest)
			return
		}
	}

	email := requestBody["email"].(string)
	password := requestBody["password"].(string)
	username := requestBody["username"].(string)
	phone := requestBody["phone"].(string)
	isAdminValue, ok := requestBody["isAdmin"].(string)

	var isAdmin bool

	if !ok || (isAdminValue != "adminpass" && isAdminValue != "") {
		http.Error(w, "Admin Password missing or not correct , creating normal user!", http.StatusBadRequest)
		isAdmin = false
	} else {
		isAdmin = true
	}

	if email == "" || password == "" || username == "" || phone == "" {
		http.Error(w, "Email, password, phone and username are required", http.StatusBadRequest)
		return
	}

	if !isValidEmail(email) {
		http.Error(w, "Invalid email address", http.StatusBadRequest)
		return
	}
	if !isValidPhoneNumber(phone) {
		http.Error(w, "Invalid phone number format. Please enter 10 digits.", http.StatusBadRequest)
		return
	}

	collection := client.Database("amazon_db").Collection("users")

	filter := bson.M{"email": email}
	var existingUser map[string]interface{}
	err = collection.FindOne(context.Background(), filter).Decode(&existingUser)
	if err == nil {
		http.Error(w, "Email is already registered", http.StatusBadRequest)
		return
	} else if err != mongo.ErrNoDocuments {
		http.Error(w, "Error checking email availability: "+err.Error(), http.StatusInternalServerError)
		return
	}

	filter = bson.M{"phone": phone}
	err = collection.FindOne(context.Background(), filter).Decode(&existingUser)
	if err == nil {
		http.Error(w, "Phone number is already registered", http.StatusBadRequest)
		return
	} else if err != mongo.ErrNoDocuments {
		http.Error(w, "Error checking phone number availability: "+err.Error(), http.StatusInternalServerError)
		return
	}

	filter = bson.M{"username": username}
	err = collection.FindOne(context.Background(), filter).Decode(&existingUser)
	if err == nil {
		http.Error(w, "username number is already registered", http.StatusBadRequest)
		return
	} else if err != mongo.ErrNoDocuments {
		http.Error(w, "Error checking username number availability: "+err.Error(), http.StatusInternalServerError)
		return
	}

	user := User{
		Email:    email,
		Password: password,
		Username: username,
		Phone:    phone,
		IsAdmin:  isAdmin,
	}

	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(w, "Error inserting data into database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Message string `json:"message"`
		User    User   `json:"user"`
	}{
		Message: "Signup successful!",
		User:    user,
	}

	jsonUser, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Error encoding user to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonUser)

}
