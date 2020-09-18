package helper

import (
	"encoding/base64"
	"net/http"
)

// WriteBeacon write GIF beacon in HTTP response
func WriteBeacon(w http.ResponseWriter, trackingStatus string) {
	// Set tracking information header
	w.Header().Set("Tk", trackingStatus)
	// Set cache policy headers
	w.Header().Set("Expires", "Mon, 01 Jan 1990 00:00:00 GMT")
	w.Header().Set("Cache-Control", "no-store")
	w.Header().Set("Pragma", "no-cache")

	// Return 1x1px transparent GIF
	w.Header().Set("Content-Type", "image/gif")
	w.WriteHeader(http.StatusOK)
	b, _ := base64.StdEncoding.DecodeString("R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7")
	w.Write(b)
}
