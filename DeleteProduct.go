package main

import (
	"context"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

func DeleteProduct(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(userKey).(User)
	product := r.Context().Value(productKey).(Product)

	collectionProduct := client.Database("amazon_db").Collection("products")

	var err error

	if product.CreatedBy != user.Username {
		http.Error(w, "Product is not created by this user", http.StatusBadRequest)
		return
	}

	_, err = collectionProduct.DeleteOne(context.TODO(), bson.M{"id": product.ID})
	if err != nil {
		http.Error(w, "Error deleting product", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product deleted successfully"))

}
