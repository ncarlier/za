package events

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/mssola/user_agent"
	"github.com/ncarlier/za/pkg/geoip"
	"github.com/ncarlier/za/pkg/helper"
	"github.com/ncarlier/za/pkg/logger"
)

// Exception contains tracked attribute when an error is catched
type Exception struct {
	BaseEvent
	Message string `json:"msg"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	URL     string `json:"url"`
	Error   string `json:"error"`
}

// Type returns event type
func (ex Exception) Type() string {
	return Types.Exception
}

// TS returns timestamp
func (ex Exception) TS() time.Time {
	return ex.Timestamp
}

// FormattedTS returns formatted timestamp
func (ex Exception) FormattedTS() string {
	return ex.Timestamp.Format("02/Jan/2006 03:04:05")
}

// Labels returns exception labels
func (ex Exception) Labels() Labels {
	labels := Labels{
		"tid":  ex.TrackingID,
		"type": ex.Type(),
	}
	// Add tags to labels
	for k, v := range ex.Tags {
		labels[k] = v
	}

	return labels
}

// NewExceptionEvent create exception event from HTTP request
func NewExceptionEvent(r *http.Request, tags map[string]string, geoipdb *geoip.DB) (Event, error) {
	q := r.URL.Query()
	ua := user_agent.New(r.UserAgent())
	browser, _ := ua.Browser()

	line, err := strconv.Atoi(q.Get("exl"))
	if err != nil {
		return nil, err
	}
	column, err := strconv.Atoi(q.Get("exc"))
	if err != nil {
		return nil, err
	}

	exception := Exception{
		BaseEvent: BaseEvent{
			TrackingID: q.Get("tid"),
			ClientIP:   helper.ParseClientIP(r),
			UserAgent:  ua.UA(),
			Browser:    browser,
			OS:         ua.OS(),
			Tags:       tags,
			Timestamp:  time.Now(),
		},
		Message: q.Get("exm"),
		Line:    line,
		Column:  column,
		URL:     q.Get("exu"),
		Error:   q.Get("exe"),
	}
	if geoipdb != nil {
		if ip := net.ParseIP(exception.ClientIP); ip != nil {
			if cc, err := geoipdb.LookupCountry(ip); err != nil {
				logger.Warning.Printf("unable to retrieve IP country code: %v", err)
			} else {
				exception.CountryCode = cc
			}
		}
	}
	return exception, nil
}
