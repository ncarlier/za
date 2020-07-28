package helper

import (
	"net/http"
	"strings"
)

// ParseClientIP extract client IP from HTTP request
func ParseClientIP(r *http.Request) string {
	clientIP := r.Header.Get("X-Forwarded-For")
	if clientIP == "" {
		clientIP = r.RemoteAddr
	}
	if comma := strings.Index(clientIP, ","); comma != -1 {
		clientIP = clientIP[0:comma]
	}
	if colon := strings.LastIndex(clientIP, ":"); colon != -1 {
		clientIP = clientIP[:colon]
	}

	return clientIP
}
