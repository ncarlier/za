package events

import (
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mssola/user_agent"
	"github.com/ncarlier/trackr/pkg/geoip"
	"github.com/ncarlier/trackr/pkg/helper"
	"github.com/ncarlier/trackr/pkg/logger"
)

// PageView contains tracked attribute when a page is viewed
type PageView struct {
	TrackingID       string
	ClientIP         string
	CountryCode      string
	Protocol         string
	UserAgent        string
	Browser          string
	OS               string
	UserLanguage     string
	DocumentHostName string
	DocumentPath     string
	DocumentReferer  string
	IsNewVisitor     bool
	IsNewSession     bool
	Tags             map[string]string
	Timestamp        time.Time
}

// HostName returns document hostname without scheme
func (p PageView) HostName() string {
	result := strings.TrimPrefix(p.DocumentHostName, "http://")
	return strings.TrimPrefix(result, "https://")
}

// Type returns event type
func (p PageView) Type() string {
	return "pageview"
}

// TS returns timestamp
func (p PageView) TS() time.Time {
	return p.Timestamp
}

// FormattedTS returns formatted timestamp
func (p PageView) FormattedTS() string {
	return p.Timestamp.Format("02/Jan/2006 03:04:05")
}

// Labels returns page view labels
func (p PageView) Labels() Labels {
	labels := Labels{
		"tid":          p.TrackingID,
		"hostname":     p.DocumentHostName,
		"path":         p.DocumentPath,
		"isNewVisitor": strconv.FormatBool(p.IsNewVisitor),
		"country":      p.CountryCode,
	}
	// Add tags to labels
	for k, v := range p.Tags {
		labels[k] = v
	}

	return labels
}

func NewPageViewEvent(r *http.Request, tags map[string]string, geoipdb *geoip.DB) (Event, error) {
	q := r.URL.Query()
	ua := user_agent.New(r.UserAgent())
	browser, _ := ua.Browser()

	pageview := PageView{
		TrackingID:       q.Get("tid"),
		ClientIP:         helper.ParseClientIP(r),
		Protocol:         r.Proto,
		UserAgent:        ua.UA(),
		Browser:          browser,
		OS:               ua.OS(),
		UserLanguage:     q.Get("ul"),
		DocumentHostName: helper.ParseHostname(q.Get("dh")),
		DocumentPath:     helper.ParsePathname(q.Get("dp")),
		DocumentReferer:  q.Get("dr"),
		IsNewVisitor:     q.Get("nv") == "1",
		IsNewSession:     q.Get("ns") == "1",
		Tags:             tags,
		Timestamp:        time.Now(),
	}
	if geoipdb != nil {
		if ip := net.ParseIP(pageview.ClientIP); ip != nil {
			if cc, err := geoipdb.LookupCountry(ip); err != nil {
				logger.Warning.Printf("unable to retrieve IP country code: %v", err)
			} else {
				pageview.CountryCode = cc
			}
		}
	}
	return pageview, nil
}
