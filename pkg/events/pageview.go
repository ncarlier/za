package events

import (
	"net/http"
	"strings"
	"time"

	"github.com/ncarlier/za/pkg/helper"
)

// PageView contains tracked attribute when a page is viewed
type PageView struct {
	BaseEvent
	Protocol         string `json:"-"`
	UserLanguage     string `json:"language"`
	DocumentHostName string `json:"hostname"`
	DocumentPath     string `json:"path"`
	DocumentReferer  string `json:"referer"`
	IsNewVisitor     bool   `json:"new_visitor"`
	IsNewSession     bool   `json:"new_session"`
	TimeOnPage       int    `json:"top"`
}

// HostName returns document hostname without scheme
func (p PageView) HostName() string {
	result := strings.TrimPrefix(p.DocumentHostName, "http://")
	return strings.TrimPrefix(result, "https://")
}

// Type returns event type
func (p PageView) Type() string {
	return Types.PageView
}

// TS returns timestamp
func (p PageView) TS() time.Time {
	return p.BaseEvent.TS()
}

// FormattedTS returns formatted timestamp
func (p PageView) FormattedTS() string {
	return p.BaseEvent.FormattedTS()
}

// Labels returns page view labels
func (p PageView) Labels() Labels {
	labels := p.BaseEvent.Labels()
	labels["type"] = p.Type()
	return labels
}

// NewPageViewEvent create page view event from HTTP request
func NewPageViewEvent(base BaseEvent, r *http.Request) (Event, error) {
	q := r.Form

	pageview := PageView{
		BaseEvent:        base,
		Protocol:         r.Proto,
		UserLanguage:     q.Get("ul"),
		DocumentHostName: helper.ParseHostname(q.Get("dh")),
		DocumentPath:     helper.ParsePathname(q.Get("dp")),
		DocumentReferer:  q.Get("dr"),
		IsNewVisitor:     q.Get("nv") == "1",
		IsNewSession:     q.Get("ns") == "1",
		TimeOnPage:       helper.ParseInt(q.Get("top"), 0),
	}
	return pageview, nil
}
