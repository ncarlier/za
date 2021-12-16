package helper

import (
	"net/http"
	"net/url"
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

// ParsePathname extract path with leading slash
func ParsePathname(s string) string {
	return "/" + strings.TrimLeft(s, "/")
}

// ParseHostname extract hostname from an URL
func ParseHostname(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		return ""
	}
	return u.Scheme + "://" + u.Host
}
