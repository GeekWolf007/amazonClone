package main

import (
	"context"
	"fmt"
	"net/http"

	"gopkg.in/mgo.v2/bson"
)

func deleteUser(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(userKey).(User)
	collection := client.Database("amazon_db").Collection("users")

	deleteFilter := bson.M{"username": user.Username}
	_, deleteErr := collection.DeleteOne(context.TODO(), deleteFilter)

	if deleteErr != nil {
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("User %s deleted successfully", user.Username)))

}
