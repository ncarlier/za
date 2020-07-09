package api

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ncarlier/trackr/pkg/config"
	"github.com/ncarlier/trackr/pkg/helper"
	"github.com/ncarlier/trackr/pkg/logger"
	"github.com/ncarlier/trackr/pkg/model"
	"github.com/ncarlier/trackr/pkg/output"
)

func collectHandler(conf *config.Config) http.Handler {
	writer, err := output.NewPrometheusOutputWriter("http://0.0.0.0:9090")
	if err != nil {
		logger.Error.Fatalf("unable to creat metric writer: %s", err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isAllowedToCollect(r) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// TODO match tracking ID with referer HTTP header

		q := r.URL.Query()
		// TODO collect more data:
		// - country code from request IP
		// - device and browser from User-Agent
		pageview := model.PageView{
			TrackingID:       q.Get("tid"),
			DocumentHostName: parseHostname(q.Get("dh")),
			DocumentPath:     parsePathname(q.Get("dp")),
			DocumentReferrer: q.Get("dr"),
			IsNewVisitor:     q.Get("nv") == "1",
			IsNewSession:     q.Get("ns") == "1",
			Timestamp:        time.Now(),
		}

		// Write metric
		writer.Write(pageview)

		// Set tracking information header
		w.Header().Set("Tk", "N")

		// Set cache policy headers
		w.Header().Set("Expires", "Mon, 01 Jan 1990 00:00:00 GMT")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")

		// Return 1x1px transparent GIF
		w.Header().Set("Content-Type", "image/gif")
		w.WriteHeader(http.StatusOK)
		b, _ := base64.StdEncoding.DecodeString("R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7")
		w.Write(b)
	})
}

func isAllowedToCollect(r *http.Request) bool {
	// Apply Do Not Track HTTP header
	if r.Header.Get("DNT") == "1" {
		return false
	}

	// Don't track prerendered pages
	if r.Header.Get("X-Moz") == "prefetch" || r.Header.Get("X-Purpose") == "preview" {
		return false
	}

	// Don't track Bot requests
	if helper.IsBotUserAgent(r.UserAgent()) {
		return false
	}

	// Validate HTTP request
	requiredQueryVars := []string{"tid", "t", "dh", "dp"}
	q := r.URL.Query()
	for _, k := range requiredQueryVars {
		if q.Get(k) == "" {
			return false
		}
	}

	return true
}

func parsePathname(p string) string {
	return "/" + strings.TrimLeft(p, "/")
}

func parseHostname(r string) string {
	u, err := url.Parse(r)
	if err != nil {
		return ""
	}
	return u.Scheme + "://" + u.Host
}
