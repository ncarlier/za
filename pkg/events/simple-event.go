package events

import (
	"encoding/base64"
	"encoding/json"
	"net"
	"net/http"
	"time"

	"github.com/mssola/user_agent"
	"github.com/ncarlier/trackr/pkg/geoip"
	"github.com/ncarlier/trackr/pkg/helper"
	"github.com/ncarlier/trackr/pkg/logger"
)

// SimpleEvent contains tracked attribute for a simple event
type SimpleEvent struct {
	BaseEvent
	Payload map[string]interface{} `json:"payload"`
}

// Type returns event type
func (se SimpleEvent) Type() string {
	return Types.Event
}

// TS returns timestamp
func (se SimpleEvent) TS() time.Time {
	return se.Timestamp
}

// FormattedTS returns formatted timestamp
func (se SimpleEvent) FormattedTS() string {
	return se.Timestamp.Format("02/Jan/2006 03:04:05")
}

// Labels returns exception labels
func (se SimpleEvent) Labels() Labels {
	labels := Labels{
		"tid":  se.TrackingID,
		"type": se.Type(),
	}
	// Add tags to labels
	for k, v := range se.Tags {
		labels[k] = v
	}

	return labels
}

// NewSimpleEvent create simple event from HTTP request
func NewSimpleEvent(r *http.Request, tags map[string]string, geoipdb *geoip.DB) (Event, error) {
	q := r.URL.Query()
	ua := user_agent.New(r.UserAgent())
	browser, _ := ua.Browser()

	d := q.Get("d")
	// TODO add support to JWT payload
	data, err := base64.StdEncoding.DecodeString(d)
	if err != nil {
		return nil, err
	}
	var objmap map[string]interface{}
	if err = json.Unmarshal(data, &objmap); err != nil {
		return nil, err
	}

	event := SimpleEvent{
		BaseEvent: BaseEvent{
			TrackingID: q.Get("tid"),
			ClientIP:   helper.ParseClientIP(r),
			UserAgent:  ua.UA(),
			Browser:    browser,
			OS:         ua.OS(),
			Tags:       tags,
			Timestamp:  time.Now(),
		},
		Payload: objmap,
	}
	if geoipdb != nil {
		if ip := net.ParseIP(event.ClientIP); ip != nil {
			if cc, err := geoipdb.LookupCountry(ip); err != nil {
				logger.Warning.Printf("unable to retrieve IP country code: %v", err)
			} else {
				event.CountryCode = cc
			}
		}
	}
	return event, nil
}
