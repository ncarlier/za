package model

import (
	"strconv"
	"strings"
	"time"
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
	DocumentHostName string
	DocumentPath     string
	DocumentReferer  string
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

// Labels returns page view labels
func (p PageView) Labels() Labels {
	return Labels{
		"tid":          p.TrackingID,
		"hostname":     p.DocumentHostName,
		"path":         p.DocumentPath,
		"isNewVisitor": strconv.FormatBool(p.IsNewVisitor),
	}
}
