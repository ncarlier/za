package model

// Output writer
type Output interface {
	// Connect to the Output
	Connect() error
	// Close any connections to the Output
	Close() error
	// Send page view to the Output
	Send(view PageView) error
}
