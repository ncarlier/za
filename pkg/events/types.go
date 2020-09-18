package events

import "time"

// EventTypes event types
var EventTypes = map[string]bool{
	"pageview":  true,
	"exception": true,
	"event":     true,
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
