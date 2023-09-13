package main

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var client *mongo.Client

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func isValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, _ := regexp.MatchString(regex, email)
	return match
}
func isValidPhoneNumber(phone string) bool {
	regex := `^\d{10}$`
	match, _ := regexp.MatchString(regex, phone)
	return match
}

func initDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		panic(err)
	}
}

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	expectedKeysToSignup := []string{"email", "password", "username", "phone", "isAdmin"}

	for key := range r.Form {
		if !contains(expectedKeysToSignup, key) {
			http.Error(w, "Unexpected key in form data: "+key, http.StatusBadRequest)
			return
		}
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	username := r.FormValue("username")
	phone := r.FormValue("phone")
	isAdminValue := r.FormValue("isAdmin")

	var isAdmin bool

	if isAdminValue == "adminpass" {
		isAdmin = true
	} else if isAdminValue == "" {
		isAdmin = false
	} else {
		http.Error(w, "Invalid admin password, creating normal user", http.StatusBadRequest)
		isAdmin = false
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
		Email:      email,
		Password:   password,
		Username:   username,
		Phone:      phone,
		IsLogged:   false,
		IsAdmin:    isAdmin,
		LoginToken: "",
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
