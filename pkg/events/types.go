package events

import "time"

// Types of events
var Types = newEventTypes()

func newEventTypes() *eventTypes {
	return &eventTypes{
		Event:     "event",
		PageView:  "pageview",
		Exception: "exception",
	}
}

type eventTypes struct {
	Event     string
	PageView  string
	Exception string
}

func (t *eventTypes) IsValid(name string) bool {
	return name == t.Event || name == t.Exception || name == t.PageView
}

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

// Event is the generic interface for a tracking event
type Event interface {
	// Type returns event type
	Type() string
	// TS returns event timestamp
	TS() time.Time
	// FormattedTS returns event formated timestamp
	FormattedTS() string
	// Labels return event labels
	Labels() Labels
}
