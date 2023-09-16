package main

import (
	"encoding/json"
	"net/http"
)

func ViewCart(w http.ResponseWriter, r *http.Request) {

	user := r.Context().Value(userKey).(User)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user.Cart)
}
