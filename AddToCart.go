package main

import (
	"context"
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

func AddToCart(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(userKey).(User)
	product := r.Context().Value(productKey).(Product)

	collectionUser := client.Database("amazon_db").Collection("users")

	var found bool
	if len(user.Cart) != 0 {
		for i, item := range user.Cart {
			if item.Product.ID == product.ID {
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
	_, errUpdating := collectionUser.UpdateOne(context.TODO(), bson.M{"username": user.Username}, update)
	if errUpdating != nil {
		http.Error(w, "Error updating user cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
