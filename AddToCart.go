package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

func AddToCart(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	username := r.FormValue("username")
	login_token := r.FormValue("login_token")
	product_id := r.FormValue("product_id")

	expectedKeysToAddToCart := []string{"product_id", "login_token", "username"}

	for key := range r.Form {
		if !contains(expectedKeysToAddToCart, key) {
			http.Error(w, "Unexpected key in form data: "+key, http.StatusBadRequest)
			return
		}
	}

	if username == "" {
		http.Error(w, "Username field is missing", http.StatusBadRequest)
		return
	}
	if login_token == "" {
		http.Error(w, "Login token is missing", http.StatusBadRequest)
		return
	}
	if product_id == "" {
		http.Error(w, "Product id is missing", http.StatusBadRequest)
		return
	}

	collectionUser := client.Database("amazon_db").Collection("users")
	collectionProduct := client.Database("amazon_db").Collection("products")

	var filter_username bson.M
	var filter_product bson.M
	var user User
	var product Product

	if username != "" {
		filter_username = bson.M{"username": username}
		error_name := collectionUser.FindOne(context.TODO(), filter_username).Decode(&user)
		if error_name != nil {
			http.Error(w, "Username is not registered", http.StatusBadRequest)
			return
		}
	}
	if login_token != user.LoginToken {
		http.Error(w, "Incorrect login token", http.StatusBadRequest)
		return
	}

	if product_id != "" {
		filter_product = bson.M{"id": product_id}
		error_name := collectionProduct.FindOne(context.TODO(), filter_product).Decode(&product)
		if error_name != nil {
			fmt.Println(error_name)
			http.Error(w, "Product is not registered or ProductID is wrong", http.StatusBadRequest)
			return
		}
	}

	var found bool
	if len(user.Cart) != 0 {
		for i, item := range user.Cart {
			if item.Product.ID == product_id {
				found = true
				user.Cart[i].Quantity++
				break
			}
		}
		if !found {
			user.Cart = append(user.Cart, CartItem{Product: &product, Quantity: 1})
		}
	} else {
		user.Cart = append(user.Cart, CartItem{Product: &product, Quantity: 1})
	}

	update := bson.M{"$set": bson.M{"cart": user.Cart}}
	_, errUpdating := collectionUser.UpdateOne(context.TODO(), filter_username, update)
	if errUpdating != nil {
		http.Error(w, "Error updating user cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
