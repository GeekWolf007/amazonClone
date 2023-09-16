package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func AddProduct(w http.ResponseWriter, r *http.Request) {

	var requestBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusInternalServerError)
		return
	}

	if requestBody["product_name"] == nil || requestBody["product_price"] == nil || requestBody["product_description"] == nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	productName := requestBody["product_name"].(string)
	productPrice := requestBody["product_price"].(float64)
	productDescription := requestBody["product_description"].(string)

	collectionProduct := client.Database("amazon_db").Collection("products")
	user := r.Context().Value(userKey).(User)

	product := Product{
		ID:                 uuid.New().String(),
		ProductName:        productName,
		ProductPrice:       productPrice,
		ProductDescription: productDescription,
		CreatedBy:          user.Username,
	}

	_, err = collectionProduct.InsertOne(context.TODO(), product)
	if err != nil {
		http.Error(w, "Error inserting product", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Product added successfully"))
}
