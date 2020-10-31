package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/mssola/user_agent"
	"github.com/ncarlier/za/pkg/config"
	"github.com/ncarlier/za/pkg/events"
	"github.com/ncarlier/za/pkg/geoip"
	"github.com/ncarlier/za/pkg/helper"
	"github.com/ncarlier/za/pkg/logger"
	"github.com/ncarlier/za/pkg/outputs"
)

func collectHandler(mux *http.ServeMux, conf *config.Config) http.Handler {
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
			helper.WriteGifBeacon(w, "N")
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
			logger.Debug.Printf("invalid request parameters: %v", r.URL.Query())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		q := r.URL.Query()

		trackingID := q.Get("tid")
		eventType := q.Get("t")
		tracker := conf.GetTracker(trackingID)
		if tracker == nil || (eventType != "badge" && !strings.HasPrefix(r.Referer(), tracker.Origin)) {
			logger.Debug.Printf("tracking ID %s doesn't match website origin: %s", trackingID, r.Referer())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var event events.Event
		switch eventType {
		case "pageview":
			event, err = events.NewPageViewEvent(r, conf.Global.Tags, geoIPDatabase)
		case "exception":
			event, err = events.NewExceptionEvent(r, conf.Global.Tags, geoIPDatabase)
		case "event", "badge":
			event, err = events.NewSimpleEvent(r, conf.Global.Tags, geoIPDatabase)
		default:
			err = errors.New("event type not yet implemented: " + eventType)
		}
		if err != nil {
			logger.Debug.Printf("error: unable to create event: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Send event to outputs manager
		outputs.SendEvent(event)

		if eventType == "badge" {
			// Write badge beacon as response
			helper.WriteBadgeBeacon(w, "P", tracker.Badge)
		} else {
			// Write GIF beacon as response
			helper.WriteGifBeacon(w, "P")
		}
	})
}

func isValidRequest(r *http.Request) bool {
	// Validate HTTP request
	q := r.URL.Query()
	tid := q.Get("tid")
	t := q.Get("t")
	return tid != "" && events.Types.IsValid(t)
}
