package middleware

import (
	"net/http"
)

// Middleware function definition
type Middleware func(next http.Handler) http.Handler
