package main

import (
	"context"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

const userKey contextKey = "user"

func filterUser(handler http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		username, err := ExtractUsernameFromJWT(r.Header["Token"][0])
		if err != nil {
			http.Error(w, "Error extracting username from JWT", http.StatusInternalServerError)
		}

		collection := client.Database("amazon_db").Collection("users")

		var filter bson.M
		var user User

		if username != "" {
			filter = bson.M{"username": username}
			error_name := collection.FindOne(context.TODO(), filter).Decode(&user)
			if error_name != nil {
				http.Error(w, "Username is not registered", http.StatusBadRequest)
				return
			}
		}

		ctx := context.WithValue(r.Context(), userKey, user)
		handler.ServeHTTP(w, r.WithContext(ctx))
	}
}
