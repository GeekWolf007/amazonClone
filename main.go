package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

func ClearCollection(w http.ResponseWriter, r *http.Request) {
	collection := client.Database("amazon_db").Collection("users")
	_, err := collection.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "Error clearing collection: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "Collection cleared successfully")
}

func main() {
	initDB()

	mux := http.NewServeMux()

	mux.HandleFunc("/user/signup", Signup)
	mux.HandleFunc("/user/login", Login)
	mux.HandleFunc("/user/delete", deleteUser)
	mux.HandleFunc("/allusers", ShowAllUsers)
	mux.HandleFunc("/products/viewall", ViewProducts)
	mux.HandleFunc("/products/add", AddProduct)
	mux.HandleFunc("/products/search", SearchProducts)
	mux.HandleFunc("/cart", ViewCart)
	mux.HandleFunc("/cart/add", AddToCart)
	mux.HandleFunc("/cart/remove", RemoveFromCart)

	mux.HandleFunc("/clear", ClearCollection)

	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
