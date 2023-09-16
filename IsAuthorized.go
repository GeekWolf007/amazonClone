package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type Error struct {
	Message string `json:"message"`
}

func SetError(err Error, message string) Error {
	err.Message = message
	return err
}

type contextKey string

func IsAuthorized(method string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != method {
			var err Error
			err = SetError(err, "Method not allowed")
			json.NewEncoder(w).Encode(err)
			return
		}

		if r.Header["Token"] == nil {
			var err Error
			err = SetError(err, "No Token Found")
			json.NewEncoder(w).Encode(err)
			return
		}

		var mySigningKey = []byte("secretkey")

		token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("there was an error in parsing")
			}
			return mySigningKey, nil
		})

		if err != nil {
			var err Error
			err = SetError(err, "Your Token has been expired")
			json.NewEncoder(w).Encode(err)
			return
		}

		if token.Valid {
			handler.ServeHTTP(w, r)
			return
		}

		var error Error
		error = SetError(error, "Not Authorized")
		json.NewEncoder(w).Encode(error)
	}
}
