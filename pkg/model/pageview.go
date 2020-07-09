package model

import (
	"time"
)

// PageView contains tracked attribute when a page is viewed
type PageView struct {
	TrackingID       string
	DocumentHostName string
	DocumentPath     string
	DocumentReferrer string
	IsNewVisitor     bool
	IsNewSession     bool
	Timestamp        time.Time
}
