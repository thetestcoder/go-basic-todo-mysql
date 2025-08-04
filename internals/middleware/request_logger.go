package middleware

import (
	"log"
	"net/http"
	"time"
)

func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			writer,
			http.StatusOK,
		}

		next.ServeHTTP(rw, request)

		log.Printf("Method: %s, URL: %s, Status: %d Took: %v", request.Method, request.URL, rw.status, time.Since(start))
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}
