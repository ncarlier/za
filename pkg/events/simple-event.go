package events

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"time"
)

// SimpleEvent contains tracked attribute for a simple event
type SimpleEvent struct {
	BaseEvent
	Payload map[string]interface{} `json:"payload"`
}

// Type returns event type
func (se *SimpleEvent) Type() string {
	return Types.Event
}

// TS returns timestamp
func (se *SimpleEvent) TS() time.Time {
	return se.BaseEvent.TS()
}

// FormattedTS returns formatted timestamp
func (se *SimpleEvent) FormattedTS() string {
	return se.BaseEvent.FormattedTS()
}

// Labels returns exception labels
func (se *SimpleEvent) Labels() Labels {
	labels := se.BaseEvent.Labels()
	labels["type"] = se.Type()
	return labels
}

// ToMap convert event to map structure
func (se *SimpleEvent) ToMap() map[string]interface{} {
	result := se.BaseEvent.ToMap()
	result["payload"] = se.Payload

	return result
}

// NewSimpleEvent create simple event from HTTP request
func NewSimpleEvent(base *BaseEvent, r *http.Request) (Event, error) {
	q := r.Form

	var objmap map[string]interface{}
	d := q.Get("d")
	if d != "" {
		// TODO add support to JWT payload
		data, err := base64.StdEncoding.DecodeString(d)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(data, &objmap); err != nil {
			return nil, err
		}
	}

	event := SimpleEvent{
		BaseEvent: *base,
		Payload:   objmap,
	}
	return &event, nil
}
