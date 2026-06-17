package middleware

import (
	"log"
	"net/http"
	"time"
)

type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &ResponseWriter{
			ResponseWriter: w,
			statusCode:     200,
		}

		next.ServeHTTP(rw, r)

		log.Printf(
			"%s, %s, %d, %s",
			r.Method,
			r.URL.Path,
			rw.statusCode,
			time.Since(start),
		)
	})
}
