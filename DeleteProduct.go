package main

import (
	"context"
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

func DeleteProduct(w http.ResponseWriter, r *http.Request) {

	var err error

	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	expectedKeysToDeleteProduct := []string{"username", "login_token", "password", "product_id"}

	for key := range r.Form {
		if !contains(expectedKeysToDeleteProduct, key) {
			http.Error(w, "Unexpected key in form data: "+key, http.StatusBadRequest)
			return
		}
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

	productID := r.URL.Query().Get("product_id")

	requiredFields := map[string]string{
		"productID": productID,
	}

	for field, value := range requiredFields {
		if value == "" {
			http.Error(w, field+" is required", http.StatusBadRequest)
			return
		}
	}

	var filter_username bson.M
	var filter_product bson.M
	var user User
	var Product Product
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

	if productID != "" {
		filter_product = bson.M{"id": productID}
		error_product := collectionProduct.FindOne(context.Background(), filter_product).Decode(&Product)
		if error_product != nil {
			http.Error(w, "Product is not registered", http.StatusBadRequest)
			return
		}
	}
	if Product.CreatedBy != username {
		http.Error(w, "Product is not created by this user", http.StatusBadRequest)
		return
	}

	_, err = collectionProduct.DeleteOne(context.TODO(), filter_product)
	if err != nil {
		http.Error(w, "Error deleting product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product deleted successfully"))

}
