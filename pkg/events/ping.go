package events

import (
	"net/http"
)

// PingEvent contains tracked attribute for a ping event
type PingEvent struct {
	SimpleEvent
	From string
	To   string
}

// Type returns event type
func (pe *PingEvent) Type() string {
	return Types.Ping
}

// ToMap convert event to map structure
func (pe *PingEvent) ToMap() map[string]interface{} {
	result := pe.SimpleEvent.ToMap()
	result["type"] = pe.Type()
	result["from"] = pe.From
	result["to"] = pe.To

	return result
}

// NewPingEvent create ping event from HTTP request
func NewPingEvent(base *BaseEvent, r *http.Request) (*PingEvent, error) {
	simple, err := NewSimpleEvent(base, r)
	if err != nil {
		return nil, err
	}

	event := PingEvent{
		SimpleEvent: *simple,
		From:        r.Header.Get("Ping-From"),
		To:          r.Header.Get("Ping-To"),
	}
	return &event, nil
}
