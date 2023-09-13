package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
)

func ViewProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	collection := client.Database("amazon_db").Collection("products")

	count, err := collection.CountDocuments(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	product_id1 := uuid.New().String()
	fmt.Println(product_id1)
	product_id2 := uuid.New().String()

	if count == 0 {
		defaultProducts := []interface{}{
			Product{
				ID:                 product_id1,
				ProductName:        "Default Product 1",
				ProductPrice:       100,
				ProductDescription: "Description for Default Product 1",
				CreatedBy:          "defaultProductCreator",
			},
			Product{
				ID:                 product_id2,
				ProductName:        "Default Product 2",
				ProductPrice:       200,
				ProductDescription: "Description for Default Product 2",
				CreatedBy:          "defaultProductCreator",
			},
		}

		_, err := collection.InsertMany(context.TODO(), defaultProducts)
		if err != nil {
			log.Fatal(err)
		}
	}

	cursor, err := collection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	var products []Product
	if err := cursor.All(context.TODO(), &products); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
