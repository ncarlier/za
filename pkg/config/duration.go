package config

import (
	"time"
)

// Duration custom configuration type
type Duration struct {
	time.Duration
}

// UnmarshalTOML unmarshal TOML bytes to Duration
func (d *Duration) UnmarshalTOML(data []byte) (err error) {
	s := string(data)
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	d.Duration, err = time.ParseDuration(s)
	return
}
