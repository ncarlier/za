package model

// Output writer
type Output interface {
	// Connect to the Output
	Connect() error
	// Close any connections to the Output
	Close() error
	// SendPageView page view to the Output
	SendPageView(view PageView) error
}
