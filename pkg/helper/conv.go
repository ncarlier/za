package helper

import "strconv"

// ParseIntWithFallback convert string to int and use the fallback if impossible
func ParseInt(s string, fallback int) int {
	if s == "" {
		return fallback
	}
	result, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}
	return result
}
