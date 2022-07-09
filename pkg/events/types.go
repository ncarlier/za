package events

import "time"

// Types of events
var Types = newEventTypes()

func newEventTypes() *eventTypes {
	return &eventTypes{
		Badge:     "badge", // event alias
		Event:     "event",
		Exception: "exception",
		PageView:  "pageview",
	}
}

type eventTypes struct {
	Badge     string
	Event     string
	Exception string
	PageView  string
}

func (t *eventTypes) IsValid(name string) bool {
	switch name {
	case
		t.Badge,
		t.Event,
		t.Exception,
		t.PageView:
		return true
	}
	return false
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
