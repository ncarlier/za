package model

// Tracker structure
type Tracker struct {
	Origin     string
	TrackingID string `toml:"tracking_id"`
}
