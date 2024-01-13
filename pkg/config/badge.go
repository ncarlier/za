package config

import "regexp"

var badgeRe = regexp.MustCompile(`^\w+|\w+|#\d{6}$`)

func validateBadgeSyntaxe(badge string) bool {
	return badgeRe.MatchString(badge)
}
