package api

import (
	"errors"
	"fmt"
	"net/http"

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
			logger.Debug.Printf("invalid request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// Validate tracker
		tracker := conf.GetTracker(trackingID)
		if tracker == nil || (eventType != "badge" && !helper.Match(tracker.Origin, r.Referer())) {
			logger.Debug.Printf("tracking ID %s doesn't match website origin: %s", trackingID, r.Referer())
			w.WriteHeader(http.StatusBadRequest)
			return
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
		case "event", "badge":
			event, err = events.NewSimpleEvent(base, r)
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

		// Build response according to the event and the request
		buildResponse(r, w, tracker, event)
	})
}

func buildResponse(r *http.Request, w http.ResponseWriter, tracker *config.Tracker, event events.Event) {
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

func parseRequest(r *http.Request) (trackingID string, eventType string, err error) {
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
