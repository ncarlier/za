package config

import "regexp"

var badgeRe = regexp.MustCompile(`^\w+|\w+|#\d{6}$`)

// Tracker structure
type Tracker struct {
	Origin     string
	TrackingID string `toml:"tracking_id"`
	Badge      string `toml:"badge"`
}

func validateBadgeSyntaxe(badge string) bool {
	return badgeRe.MatchString(badge)
}
