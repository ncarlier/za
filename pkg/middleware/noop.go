package middleware

import (
	"net/http"
)

// Noop is a middleware to do nothing, yes!
func Noop(next http.Handler) http.Handler {
	return http.HandlerFunc(next.ServeHTTP)
}
