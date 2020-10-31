package api

import (
	"net/http"
	"path"
	"strings"

	"github.com/ncarlier/za/pkg/config"
)

func badgeHandler(mux *http.ServeMux, conf *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Rewrite-Original-URI", r.URL.RequestURI())
		tid := path.Base(r.URL.Path)
		if !strings.HasSuffix(tid, ".svg") {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		tid = strings.TrimSuffix(tid, ".svg")
		r.URL.Path = "/collect"
		q := r.URL.Query()
		q.Set("t", "badge")
		q.Set("tid", tid)
		r.URL.RawQuery = q.Encode()
		mux.ServeHTTP(w, r)
	})
}
