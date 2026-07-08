package middleware

import (
	"net/http"
	"strings"

	appjwt "myprojects/internal/jwt"

	"github.com/golang-jwt/jwt/v5"
)

func JWT(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			tokenString := authHeader[7:]

			claims := &appjwt.Claims{}

			token, err := jwt.ParseWithClaims(
				tokenString,
				claims,
				func(t *jwt.Token) (interface{}, error) {
					return []byte(secret), nil
				},
			)

			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
