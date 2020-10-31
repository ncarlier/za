package api

import (
	"net/http"

	"github.com/ncarlier/za/pkg/config"
	"github.com/ncarlier/za/pkg/middleware"
)

// HandlerFunc custom function handler
type HandlerFunc func(mux *http.ServeMux, conf *config.Config) http.Handler

// Route is the structure of an HTTP route definition
type Route struct {
	Path        string
	HandlerFunc HandlerFunc
	Middlewares []middleware.Middleware
}

func route(path string, handler HandlerFunc, middlewares ...middleware.Middleware) Route {
	return Route{
		Path:        path,
		HandlerFunc: handler,
		Middlewares: middlewares,
	}
}

// Routes is a list of Route
type Routes []Route

func routes(conf *config.Config) Routes {
	return Routes{
		route(
			"/collect",
			collectHandler,
			middleware.Methods("GET", "POST"),
		),
		route(
			"/badge/",
			badgeHandler,
			middleware.Methods("GET"),
		),
		route(
			"/za.js",
			fileHandler("za.js"),
			middleware.Methods("GET"),
		),
		route(
			"/za.min.js",
			fileHandler("za.min.js"),
			middleware.Methods("GET"),
		),
		route(
			"/varz",
			varzHandler,
			middleware.Methods("GET"),
		),
		route(
			"/healthz",
			healthzHandler,
			middleware.Methods("GET"),
		),
		route(
			"/",
			infoHandler,
			middleware.Methods("GET", "POST"),
		),
	}
}
