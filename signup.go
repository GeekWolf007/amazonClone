package main

import (
	"context"
	"encoding/json"
	"net/http"
)

func Signup(w http.ResponseWriter, r *http.Request) {

	var requestBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusInternalServerError)
		return
	}

	mandatoryFields := []string{"email", "password", "username", "phone"}
	if err := checkMandatoryFields(mandatoryFields, requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email := requestBody["email"].(string)
	password := requestBody["password"].(string)
	username := requestBody["username"].(string)
	phone := requestBody["phone"].(string)
	isAdminValue, ok := requestBody["isAdmin"].(string)

	var isAdmin bool = false
	if ok && isAdminValue == "adminpass" {
		isAdmin = true
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

	if checkDuplicate("email", email, w, collection) ||
		checkDuplicate("phone", phone, w, collection) ||
		checkDuplicate("username", username, w, collection) {
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
