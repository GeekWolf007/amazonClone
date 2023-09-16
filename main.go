package main

import (
	"log"
	"net/http"
)

// func ClearCollection(collection *mongo.Collection) error {
// 	_, err := collection.DeleteMany(context.Background(), bson.M{})
// 	return err
// }
// func ClearHandler(collection *mongo.Collection) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		if r.Method != "GET" {
// 			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 			return
// 		}

// 		err := ClearCollection(collection)
// 		if err != nil {
// 			http.Error(w, "Error clearing collection: "+err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("Collection cleared successfully"))
// 	}
// }

func main() {
	initDB()

	mux := http.NewServeMux()

	// collection := client.Database("amazon_db").Collection("users")
	// mux.HandleFunc("/clear", ClearHandler(collection))

	mux.HandleFunc("/user/signup", Signup)
	mux.HandleFunc("/user/login", Login)
	mux.HandleFunc("/user/delete", IsAuthorized("DELETE", filterUser(deleteUser)))
	mux.HandleFunc("/allusers", IsAuthorized("POST", filterUser(ShowAllUsers)))
	mux.HandleFunc("/products/viewall", ViewProducts)
	mux.HandleFunc("/products/add", IsAuthorized("POST", filterUser(AddProduct)))
	mux.HandleFunc("/products/delete", IsAuthorized("DELETE", filterProduct(filterUser(DeleteProduct))))
	mux.HandleFunc("/cart", IsAuthorized("GET", filterUser(ViewCart)))
	mux.HandleFunc("/cart/add", IsAuthorized("POST", filterProduct(filterUser(AddToCart))))
	mux.HandleFunc("/cart/remove", IsAuthorized("DELETE", filterProduct(filterUser(RemoveFromCart))))

	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
