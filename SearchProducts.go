package main

import (
	"context"
	"encoding/json"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

func SearchProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Search query parameter 'query' is required", http.StatusBadRequest)
		return
	}

	collection := client.Database("amazon_db").Collection("products")

	filter := bson.M{
		"$or": []bson.M{
			{"product_name": bson.M{"$regex": query, "$options": "i"}},
			{"product_description": bson.M{"$regex": query, "$options": "i"}},
		},
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		http.Error(w, "Error searching products", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var products []Product
	for cursor.Next(context.TODO()) {
		var product Product
		err := cursor.Decode(&product)
		if err != nil {
			http.Error(w, "Error decoding product", http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	json.NewEncoder(w).Encode(products)
}
