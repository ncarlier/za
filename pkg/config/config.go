package config

import (
	"github.com/ncarlier/za/pkg/outputs"
	"github.com/ncarlier/za/pkg/usage"
)

// Config is the root of the configuration
type Config struct {
	Log      LogConfig
	HTTP     HTTPConfig
	GeoIP    GeoIPConfig `toml:"geo-ip"`
	Global   GlobalConfig
	Trackers []TrackerConfig
	Outputs  []outputs.Output
}

// LogConfig for log configuration section
type LogConfig struct {
	Level   string
	Format  string
	Modules string
}

// HTTPConfig for HTTP configuration section
type HTTPConfig struct {
	ListenAddr string
}

// GeoIPConfig for GeoIP configuration section
type GeoIPConfig struct {
	Database string
}

// GlobalConfig for global configuration section
type GlobalConfig struct {
	Tags map[string]string
}

// TrackerConfig for tracker configurtion section
type TrackerConfig struct {
	Origin       string
	TrackingID   string `toml:"tracking_id"`
	Badge        string
	RateLimiting map[string]interface{}
	RateLimiter  usage.RateLimiter
}

// NewConfig create new configuration
func NewConfig() *Config {
	c := &Config{
		Log: LogConfig{
			Level:  "info",
			Format: "json",
		},
		HTTP: HTTPConfig{
			ListenAddr: ":8080",
		},
		Global: GlobalConfig{
			Tags: make(map[string]string),
		},
		Trackers: make([]TrackerConfig, 0),
		Outputs:  make([]outputs.Output, 0),
	}
	return c
}

// GetTracker retrive tracker configuration
func (c *Config) GetTracker(trackingID string) *TrackerConfig {
	for _, tracker := range c.Trackers {
		if tracker.TrackingID == trackingID {
			return &tracker
		}
	}
	return nil
}
