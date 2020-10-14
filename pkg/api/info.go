package api

import (
	"encoding/json"
	"net/http"

	"github.com/ncarlier/za/pkg/config"
	"github.com/ncarlier/za/pkg/version"
)

// Info API informations model structure.
type Info struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func infoHandler(conf *config.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		info := Info{
			Name:    "ZerØ Analytics",
			Version: version.Version,
		}
		data, err := json.Marshal(info)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	})
}
