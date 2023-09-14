package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"gopkg.in/mgo.v2/bson"
)

func AddProduct(w http.ResponseWriter, r *http.Request) {

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

	expectedKeysToAddProduct := []string{"product_name", "product_price", "product_description"}

	for key := range r.Form {
		if !contains(expectedKeysToAddProduct, key) {
			http.Error(w, "Unexpected key in form data: "+key, http.StatusBadRequest)
			return
		}
	}

	productName := requestBody["product_name"].(string)
	productPrice := requestBody["product_price"].(float64)
	productDescription := requestBody["product_description"].(string)

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

	product := Product{
		ID:                 uuid.New().String(),
		ProductName:        productName,
		ProductPrice:       productPrice,
		ProductDescription: productDescription,
		CreatedBy:          username,
	}

	_, err = collectionProduct.InsertOne(context.TODO(), product)
	if err != nil {
		http.Error(w, "Error inserting product", http.StatusInternalServerError)
		return
	}

	fmt.Println(username)
	w.Write([]byte("Product added successfully"))
}
