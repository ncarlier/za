package api

import (
	"net/http"

	"github.com/ncarlier/trackr/pkg/config"
	"github.com/ncarlier/trackr/pkg/middleware"
)

// HandlerFunc custom function handler
type HandlerFunc func(conf *config.Config) http.Handler

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

var routes = Routes{
	route("/", infoHandler, middleware.Methods("GET", "POST")),
	route("/collect", collectHandler, middleware.Methods("GET", "POST")),
	route("/trackr.js", fileHandler("trackr.js"), middleware.Methods("GET")),
	route("/trackr.min.js", fileHandler("trackr.min.js"), middleware.Methods("GET")),
	route("/healthz", healthzHandler, middleware.Methods("GET")),
	route("/varz", varzHandler, middleware.Methods("GET")),
}
