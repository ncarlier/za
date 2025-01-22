package api

import (
	"net/http"
	"path"

	"github.com/ncarlier/za/pkg/config"
)

func pingHandler(mux *http.ServeMux, conf *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Rewrite-Original-URI", r.URL.RequestURI())
		tid := path.Base(r.URL.Path)
		r.URL.Path = "/collect"
		q := r.URL.Query()
		q.Set("t", "ping")
		q.Set("tid", tid)
		r.URL.RawQuery = q.Encode()
		mux.ServeHTTP(w, r)
	})
}
