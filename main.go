package main

import (
	"context"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

func ClearCollection(collection *mongo.Collection) error {
	_, err := collection.DeleteMany(context.Background(), bson.M{})
	return err
}
func ClearHandler(collection *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		err := ClearCollection(collection)
		if err != nil {
			http.Error(w, "Error clearing collection: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Collection cleared successfully"))
	}
}

func main() {
	initDB()

	mux := http.NewServeMux()

	collection := client.Database("amazon_db").Collection("users")

	mux.HandleFunc("/clear", ClearHandler(collection))

	mux.HandleFunc("/user/signup", Signup)
	mux.HandleFunc("/user/login", Login)
	mux.HandleFunc("/user/delete", IsAuthorized(deleteUser))
	mux.HandleFunc("/allusers", IsAuthorized(ShowAllUsers))
	mux.HandleFunc("/products/viewall", ViewProducts)
	mux.HandleFunc("/products/add", IsAuthorized(AddProduct))
	mux.HandleFunc("/products/delete", IsAuthorized(DeleteProduct))
	mux.HandleFunc("/cart", IsAuthorized(ViewCart))
	mux.HandleFunc("/cart/add", IsAuthorized(AddToCart))
	mux.HandleFunc("/cart/remove", IsAuthorized(RemoveFromCart))

	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
