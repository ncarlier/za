package api

import (
	"net/http"

	"github.com/ncarlier/trackr/pkg/config"
	"github.com/ncarlier/trackr/pkg/middleware"
)

var commonMiddlewares = []middleware.Middleware{
	middleware.Cors,
	middleware.Logger,
}

// NewRouter creates router with declared routes
func NewRouter(conf *config.Config) *http.ServeMux {
	router := http.NewServeMux()

	var middlewares = commonMiddlewares

	// Register HTTP routes...
	for _, route := range routes {
		handler := route.HandlerFunc(conf)
		for _, mw := range route.Middlewares {
			handler = mw(handler)
		}
		for _, mw := range middlewares {
			if route.Path == "/healthz" {
				continue
			}
			handler = mw(handler)
		}
		router.Handle(route.Path, handler)
	}

	return router
}
