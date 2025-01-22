package api

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/mssola/user_agent"
	"github.com/ncarlier/za/pkg/config"
	"github.com/ncarlier/za/pkg/events"
	"github.com/ncarlier/za/pkg/geoip"
	"github.com/ncarlier/za/pkg/helper"
	"github.com/ncarlier/za/pkg/outputs"
)

func collectHandler(mux *http.ServeMux, conf *config.Config) http.Handler {
	outputsMgr, err := outputs.NewOutputsManager(conf.Outputs)
	if err != nil {
		slog.Error("unable to initialize outputs manager", "error", err)
		os.Exit(1)
	}
	geoIPDatabase, err := geoip.New(conf.GeoIP.Database)
	if err != nil {
		slog.Error("unable to load geo IP database", "error", err)
		os.Exit(1)
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Respect Do Not Track HTTP header
		if handleDoNotTrackRequest(r, w) {
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

		// Parse HTTP request
		trackingID, eventType, err := parseRequest(r)
		if err != nil {
			slog.Debug("invalid request", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Validate tracker
		tracker := conf.GetTracker(trackingID)
		if tracker == nil {
			slog.Debug("tracker not found", "tid", trackingID)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Validate origin
		if eventType != "badge" && !tracker.Match(r) {
			slog.Debug("tracking ID doesn't match website origin", "tid", trackingID, "referer", r.Referer())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Apply usage limitation
		if tracker.RateLimiter != nil {
			_, _, _, ok, err := tracker.RateLimiter.Take(r.Context(), "global")
			if err != nil {
				slog.Error("unable to apply usage limitation", "error", err)
			} else if !ok {
				slog.Debug("rate limiting activated", "tid", trackingID)
				w.WriteHeader(http.StatusTooManyRequests)
				return
			}
		}

		// Create base event
		base := events.NewBaseEvent(r, conf.Global.Tags, geoIPDatabase)

		// Specialize event
		var event events.Event
		switch eventType {
		case "pageview":
			event, err = events.NewPageViewEvent(base, r)
		case "exception":
			event, err = events.NewExceptionEvent(base, r)
		case "ping":
			event, err = events.NewPingEvent(base, r)
		case "event", "badge":
			event, err = events.NewSimpleEvent(base, r)
		default:
			err = errors.New("event type not yet implemented: " + eventType)
		}
		if err != nil {
			slog.Debug("unable to create event", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Send event to outputs manager
		outputsMgr.SendEvent(event)

		// Build response according to the event and the request
		buildResponse(r, w, tracker)
	})
}

func buildResponse(r *http.Request, w http.ResponseWriter, tracker *config.TrackerConfig) {
	values := r.Form
	if values.Get("t") == "badge" {
		// Write badge beacon as response
		helper.WriteBadgeBeacon(w, "P", tracker.Badge)
	} else if r.Method == http.MethodPost {
		// Write empty response
		w.WriteHeader(http.StatusAccepted)
	} else {
		// Write GIF beacon as response
		helper.WriteGifBeacon(w, "P")
	}
}

func handleDoNotTrackRequest(r *http.Request, w http.ResponseWriter) bool {
	if r.Header.Get("DNT") == "1" {
		if r.Method == http.MethodPost {
			// Write empty response
			w.WriteHeader(http.StatusNoContent)
		} else {
			// Write GIF beacon as response
			helper.WriteGifBeacon(w, "N")
		}
		return true
	}
	return false
}

func parseRequest(r *http.Request) (trackingID, eventType string, err error) {
	r.ParseForm()
	values := r.Form
	trackingID = values.Get("tid")
	eventType = values.Get("t")
	if trackingID == "" {
		err = errors.New("tracking ID not provided")
	}
	if !events.Types.IsValid(eventType) {
		err = fmt.Errorf("invalid event type: %s", eventType)
	}
	return
}
