package events

import (
	"net"
	"net/http"
	"time"

	"github.com/mssola/user_agent"
	"github.com/ncarlier/za/pkg/geoip"
	"github.com/ncarlier/za/pkg/helper"
	"github.com/ncarlier/za/pkg/logger"
)

// BaseEvent contains common events fields
type BaseEvent struct {
	TrackingID  string            `json:"tid"`
	ClientIP    string            `json:"-"`
	CountryCode string            `json:"country,omitempty"`
	UserAgent   string            `json:"-"`
	Browser     string            `json:"browser"`
	OS          string            `json:"os"`
	Tags        map[string]string `json:"tags"`
	Timestamp   time.Time         `json:"timestamp"`
}

// TS returns timestamp
func (ev BaseEvent) TS() time.Time {
	return ev.Timestamp
}

// FormattedTS returns formatted timestamp
func (ev BaseEvent) FormattedTS() string {
	return ev.Timestamp.Format("02/Jan/2006 03:04:05")
}

// Labels returns labels
func (ev BaseEvent) Labels() Labels {
	labels := Labels{
		"tid": ev.TrackingID,
	}
	// Add tags to labels
	for k, v := range ev.Tags {
		labels[k] = v
	}
	return labels
}

// NewBaseEvent create new base event
func NewBaseEvent(r *http.Request, tags map[string]string, geoipdb *geoip.DB) BaseEvent {
	q := r.URL.Query()
	ua := user_agent.New(r.UserAgent())
	browser, _ := ua.Browser()
	clientIP := helper.ParseClientIP(r)
	var cc = ""
	if geoipdb != nil {
		if ip := net.ParseIP(clientIP); ip != nil {
			var err error
			if cc, err = geoipdb.LookupCountry(ip); err != nil {
				logger.Warning.Printf("unable to retrieve IP country code: %v", err)
			}
		}
	}
	return BaseEvent{
		TrackingID:  q.Get("tid"),
		ClientIP:    clientIP,
		CountryCode: cc,
		UserAgent:   ua.UA(),
		Browser:     browser,
		OS:          ua.OS(),
		Tags:        tags,
		Timestamp:   time.Now(),
	}
}
