package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("secret")

func Login(w http.ResponseWriter, r *http.Request) {
	var user map[string]string
	json.NewDecoder(r.Body).Decode(&user)

	if user["username"] == "admin" && user["password"] == "admin" {

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user["username"],
			"exp":      time.Now().Add(time.Hour * 1).Unix(),
		})

		tokenString, _ := token.SignedString(jwtKey)

		json.NewEncoder(w).Encode(map[string]string{
			"token": tokenString,
		})
		return
	}

	http.Error(w, "Unauthorized", http.StatusUnauthorized)
}
