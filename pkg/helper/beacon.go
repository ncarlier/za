package helper

import (
	"encoding/base64"
	"net/http"
)

func writeBeaconHeaders(w http.ResponseWriter, trackingStatus string) {
	// Set tracking information header
	w.Header().Set("Tk", trackingStatus)
	// Set cache policy headers
	w.Header().Set("Expires", "Mon, 01 Jan 1990 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")
}

// WriteBadgeBeacon write SVG badge beacon in HTTP response
func WriteBadgeBeacon(w http.ResponseWriter, trackingStatus, code string) {
	// Headers
	writeBeaconHeaders(w, trackingStatus)

	badge, err := GetBadge(code)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	// Return SVG image
	w.Header().Set("Content-Type", "image/svg+xml")
	w.WriteHeader(http.StatusOK)
	w.Write(badge)
}

// WriteGifBeacon write GIF beacon in HTTP response
func WriteGifBeacon(w http.ResponseWriter, trackingStatus string) {
	// Headers
	writeBeaconHeaders(w, trackingStatus)

	// Return 1x1px transparent GIF
	w.Header().Set("Content-Type", "image/gif")
	w.WriteHeader(http.StatusOK)
	b, _ := base64.StdEncoding.DecodeString("R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7")
	w.Write(b)
}
