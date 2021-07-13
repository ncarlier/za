package api

import (
	"github.com/ncarlier/za/pkg/config"
	"github.com/ncarlier/za/pkg/middleware"
)

func routes(conf *config.Config) Routes {
	get := middleware.Methods("GET")
	getAndPost := middleware.Methods("GET", "POST")
	gzip := middleware.Gzip
	cors := middleware.Cors("*")
	logger := middleware.Logger
	return Routes{
		route(
			"/collect",
			collectHandler,
			getAndPost,
			gzip,
			cors,
			logger,
		),
		route(
			"/badge/",
			badgeHandler,
			get,
			gzip,
			cors,
		),
		route(
			"/za.js",
			fileHandler("za.js"),
			get,
			gzip,
			cors,
		),
		route(
			"/za.min.js",
			fileHandler("za.min.js"),
			middleware.Methods("GET"),
			get,
			gzip,
			cors,
		),
		route(
			"/varz",
			varzHandler,
			get,
			gzip,
			cors,
			logger,
		),
		route(
			"/healthz",
			healthzHandler,
			get,
			gzip,
			cors,
		),
		route(
			"/",
			infoHandler,
			getAndPost,
			gzip,
			cors,
			logger,
		),
	}
}
