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
		// Respect Do Not Track HTTP header
		if r.Header.Get("DNT") == "1" {
			// Write GIF beacon as response
			helper.WriteBeacon(w, "N")
			return
		}

		// Don't track prerendered pages
		if r.Header.Get("X-Moz") == "prefetch" || r.Header.Get("X-Purpose") == "preview" {
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

		// Validate HTTP request
		if !isValidRequest(r) {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		browser, _ := ua.Browser()
		q := r.URL.Query()

		trackingID := q.Get("tid")
		if !conf.ValidateTrackingID(r.Referer(), trackingID) {
			logger.Debug.Printf("tracking ID %s doesn't match website origin: %s", trackingID, r.Referer())
			helper.WriteBeacon(w, "N")
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
		outputs.SendPageView(pageview)

		// Write GIF beacon as response
		helper.WriteBeacon(w, "P")
	})
}

func isValidRequest(r *http.Request) bool {
	// Validate HTTP request
	q := r.URL.Query()
	tid := q.Get("tid")
	_, validType := model.EventTypes[q.Get("t")]
	return tid != "" && validType
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
