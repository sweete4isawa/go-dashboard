package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var JwtKey = []byte("secret")

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenStr := r.Header.Get("Authorization")

		if tokenStr == "" {
			http.Error(w, "No token", 401)
			return
		}

		tokenStr = strings.Replace(tokenStr, "Bearer ", "", 1)

		_, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return JwtKey, nil
		})

		if err != nil {
			http.Error(w, "Invalid token", 401)
			return
		}

		next(w, r)
	}
}
