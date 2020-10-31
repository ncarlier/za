package api

import (
	"net/http"

	"github.com/ncarlier/za/pkg/assets"
	"github.com/ncarlier/za/pkg/config"
)

func fileHandler(filename string) HandlerFunc {
	return func(mux *http.ServeMux, conf *config.Config) http.Handler {
		fs := http.FileServer(assets.GetFS())
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fs.ServeHTTP(w, r)
		})
	}
}
