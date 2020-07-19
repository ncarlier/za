package model

import (
	"strings"
	"time"
)

// PageView contains tracked attribute when a page is viewed
type PageView struct {
	TrackingID       string
	ClientIP         string
	Protocol         string
	UserAgent        string
	DocumentHostName string
	DocumentPath     string
	DocumentReferrer string
	IsNewVisitor     bool
	IsNewSession     bool
	Timestamp        time.Time
}

// FormattedTS returns formatted timestamp
func (p PageView) FormattedTS() string {
	return p.Timestamp.Format("02/Jan/2006 03:04:05")
}

// HostName returns document hostname without scheme
func (p PageView) HostName() string {
	result := strings.TrimPrefix(p.DocumentHostName, "http://")
	return strings.TrimPrefix(result, "https://")
}
