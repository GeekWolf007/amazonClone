package main

import (
	"context"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

const productKey contextKey = "product"

func filterProduct(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var Product Product
		collectionProduct := client.Database("amazon_db").Collection("products")

		productID := r.URL.Query().Get("product_id")

		if productID == "" {
			http.Error(w, "Product ID is not provided", http.StatusBadRequest)
			return
		}

		if productID != "" {
			error_product := collectionProduct.FindOne(context.Background(), bson.M{"id": productID}).Decode(&Product)
			if error_product != nil {
				http.Error(w, "Product is not registered", http.StatusBadRequest)
				return
			}
		}

		ctx := context.WithValue(r.Context(), productKey, Product)
		handler.ServeHTTP(w, r.WithContext(ctx))
	}
}
