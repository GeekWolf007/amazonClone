package main

import (
	"context"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
)

func AddProduct(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	expectedKeysToAddProduct := []string{"username", "login_token", "password", "product_name", "product_price", "product_description"}

	for key := range r.Form {
		if !contains(expectedKeysToAddProduct, key) {
			http.Error(w, "Unexpected key in form data: "+key, http.StatusBadRequest)
			return
		}
	}

	username := r.FormValue("username")
	loginToken := r.FormValue("login_token")
	password := r.FormValue("password")
	productName := r.FormValue("product_name")
	productPrice := r.FormValue("product_price")
	productDescription := r.FormValue("product_description")

	requiredFields := map[string]string{
		"username":           username,
		"loginToken":         loginToken,
		"password":           password,
		"productName":        productName,
		"productPrice":       productPrice,
		"productDescription": productDescription,
	}

	for field, value := range requiredFields {
		if value == "" {
			http.Error(w, field+" is required", http.StatusBadRequest)
			return
		}
	}

	var filter_username bson.M
	var user User
	collectionUser := client.Database("amazon_db").Collection("users")
	collectionProduct := client.Database("amazon_db").Collection("products")

	if username != "" {
		filter_username = bson.M{"username": username}
		error_name := collectionUser.FindOne(context.TODO(), filter_username).Decode(&user)
		if error_name != nil {
			http.Error(w, "Username is not registered", http.StatusBadRequest)
			return
		}
	}

	if password != user.Password {
		http.Error(w, "Incorrect password", http.StatusBadRequest)
		return
	}

	if loginToken != user.LoginToken {
		http.Error(w, "Incorrect login token", http.StatusBadRequest)
		return
	}

	productPriceInt, err := strconv.Atoi(productPrice)
	if err != nil {
		http.Error(w, "Error converting product price to integer", http.StatusInternalServerError)
		return
	}

	product := Product{
		ID:                 uuid.New().String(),
		ProductName:        productName,
		ProductPrice:       productPriceInt,
		ProductDescription: productDescription,
	}

	_, err = collectionProduct.InsertOne(context.TODO(), product)
	if err != nil {
		http.Error(w, "Error inserting product", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Product added successfully"))
}
