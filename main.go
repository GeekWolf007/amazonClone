package main

import (
	"log"
	"net/http"
)

func main() {
	initDB()

	mux := http.NewServeMux()

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
