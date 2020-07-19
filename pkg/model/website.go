package model

// WebSite structure
type WebSite struct {
	Origin     string
	TrackingID string `toml:"tracking_id"`
}
