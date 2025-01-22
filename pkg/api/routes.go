package api

import (
	"net/http"
	"strings"

	"github.com/ncarlier/za/pkg/config"
	"github.com/ncarlier/za/pkg/middleware"
)

func routes(conf *config.Config) Routes {
	middlewares := middleware.Middlewares{
		middleware.Gzip,
		middleware.Cors("*"),
	}
	logger := middleware.Noop
	if strings.Contains(conf.Log.Modules, "http") {
		logger = middleware.Logger
	}
	return Routes{
		route(
			"/collect",
			collectHandler,
			middlewares.UseBefore(logger).UseAfter(
				middleware.Methods(http.MethodGet, http.MethodPost),
			)...,
		),
		route(
			"/badge/",
			badgeHandler,
			middlewares.UseBefore(middleware.Methods(http.MethodGet))...,
		),
		route(
			"/ping/",
			pingHandler,
			middlewares.UseBefore(middleware.Methods(http.MethodPost))...,
		),
		route(
			"/za.js",
			fileHandler,
			middlewares.UseBefore(middleware.Methods(http.MethodGet))...,
		),
		route(
			"/za.min.js",
			fileHandler,
			middlewares.UseBefore(middleware.Methods(http.MethodGet))...,
		),
		route(
			"/varz",
			varzHandler,
			middlewares.UseBefore(middleware.Methods(http.MethodGet))...,
		),
		route(
			"/healthz",
			healthzHandler,
			middlewares.UseBefore(middleware.Methods(http.MethodGet))...,
		),
		route(
			"/",
			infoHandler,
			middlewares.UseBefore(logger).UseAfter(
				middleware.Methods(http.MethodGet),
			)...,
		),
	}
}
