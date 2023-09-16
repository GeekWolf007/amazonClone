package main

import (
	"context"
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

func RemoveFromCart(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(userKey).(User)
	Product := r.Context().Value(productKey).(Product)

	collectionUser := client.Database("amazon_db").Collection("users")

	var product CartItem
	var found bool = false

	if len(user.Cart) == 0 {
		http.Error(w, "Cart is empty", http.StatusBadRequest)
		return
	} else {
		for _, item := range user.Cart {
			if item.Product.ID == Product.ID {
				product = item
				found = true
			}
		}
	}

	if !found {
		http.Error(w, "Product is not in cart", http.StatusBadRequest)
		return
	}

	if product.Quantity == 1 {

		updatedCart := []CartItem{}
		for _, item := range user.Cart {
			if item.Product.ID != Product.ID {
				updatedCart = append(updatedCart, item)
			}
		}
		user.Cart = updatedCart
	} else {

		for i, item := range user.Cart {
			if item.Product.ID == Product.ID {
				user.Cart[i].Quantity--
				break
			}
		}
	}

	_, updateErr := collectionUser.UpdateOne(
		context.TODO(),
		bson.M{"username": user.Username},
		bson.M{"$set": bson.M{"cart": user.Cart}},
	)
	if updateErr != nil {
		http.Error(w, "Error updating user cart", http.StatusInternalServerError)
		return
	}

	userData, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "Error converting user data to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	w.Write(userData)

}
