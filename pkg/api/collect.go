package api

import (
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mssola/user_agent"
	"github.com/ncarlier/trackr/pkg/config"
	"github.com/ncarlier/trackr/pkg/geoip"
	"github.com/ncarlier/trackr/pkg/helper"
	"github.com/ncarlier/trackr/pkg/logger"
	"github.com/ncarlier/trackr/pkg/model"
	"github.com/ncarlier/trackr/pkg/outputs"
)

func collectHandler(conf *config.Config) http.Handler {
	outputs, err := outputs.NewOutputsManager(conf.Outputs)
	if err != nil {
		logger.Error.Fatalf("unable to initialize outputs manager: %s", err)
	}
	geoIPDatabase, err := geoip.New(conf.Global.GeoIPDatabase)
	if err != nil {
		logger.Error.Fatalf("unable to load geo IP database: %s", err)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !isAllowedToCollect(r) {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Decode User-Agent
		ua := user_agent.New(r.UserAgent())

		// Don't track Bot requests
		if ua.Bot() {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		browser, _ := ua.Browser()
		q := r.URL.Query()

		trackingID := q.Get("tid")
		if !conf.ValidateTrackingID(r.Referer(), trackingID) {
			logger.Debug.Printf("tracking ID %s doesn't match website origin: %s", trackingID, r.Referer())
			w.WriteHeader(http.StatusNoContent)
			return
		}

		pageview := model.PageView{
			TrackingID:       trackingID,
			ClientIP:         helper.ParseClientIP(r),
			Protocol:         r.Proto,
			UserAgent:        ua.UA(),
			Browser:          browser,
			OS:               ua.OS(),
			UserLanguage:     q.Get("ul"),
			DocumentHostName: parseHostname(q.Get("dh")),
			DocumentPath:     parsePathname(q.Get("dp")),
			DocumentReferer:  q.Get("dr"),
			IsNewVisitor:     q.Get("nv") == "1",
			IsNewSession:     q.Get("ns") == "1",
			Tags:             conf.Global.Tags,
			Timestamp:        time.Now(),
		}

		if geoIPDatabase != nil {
			if ip := net.ParseIP(pageview.ClientIP); ip != nil {
				pageview.CountryCode, err = geoIPDatabase.LookupCountry(ip)
				if err != nil {
					logger.Warning.Printf("unable to retrieve IP country code: %v", err)
				}
			}
		}

		// Send page view to outputs manager
		outputs.Send(pageview)

		// Set tracking information header
		w.Header().Set("Tk", "N")

		// Write GIF beacon as response
		helper.WriteBeacon(w)
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
