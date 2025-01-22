package events

import (
	"net/http"
	"strconv"
	"time"
)

// Exception contains tracked attribute when an error is triggered
type Exception struct {
	BaseEvent
	Message string `json:"msg"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	URL     string `json:"url"`
	Error   string `json:"error"`
}

// Type returns event type
func (ex *Exception) Type() string {
	return Types.Exception
}

// TS returns timestamp
func (ex *Exception) TS() time.Time {
	return ex.BaseEvent.TS()
}

// FormattedTS returns formatted timestamp
func (ex *Exception) FormattedTS() string {
	return ex.BaseEvent.FormattedTS()
}

// Labels returns exception labels
func (ex *Exception) Labels() Labels {
	labels := ex.BaseEvent.Labels()
	labels["type"] = ex.Type()
	return labels
}

// ToMap convert event to map structure
func (ex *Exception) ToMap() map[string]interface{} {
	result := ex.BaseEvent.ToMap()
	result["msg"] = ex.Message
	result["line"] = ex.Line
	result["column"] = ex.Column
	result["url"] = ex.URL
	result["error"] = ex.Error
	return result
}

// NewExceptionEvent create exception event from HTTP request
func NewExceptionEvent(base *BaseEvent, r *http.Request) (*Exception, error) {
	q := r.Form

	line, err := strconv.Atoi(q.Get("exl"))
	if err != nil {
		return nil, err
	}
	column, err := strconv.Atoi(q.Get("exc"))
	if err != nil {
		return nil, err
	}

	exception := Exception{
		BaseEvent: *base,
		Message:   q.Get("exm"),
		Line:      line,
		Column:    column,
		URL:       q.Get("exu"),
		Error:     q.Get("exe"),
	}
	return &exception, nil
}
