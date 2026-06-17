package middleware

import (
	"log"
	"net/http"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			value := recover()

			if value != nil {
				http.Error(w, "Internal error", http.StatusInternalServerError)
				log.Println(value)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
