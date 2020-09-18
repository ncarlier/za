package outputs

import "github.com/ncarlier/trackr/pkg/events"

// Output writer
type Output interface {
	// Connect to the Output
	Connect() error
	// Close any connections to the Output
	Close() error
	// SendEvent sent event to the Output
	SendEvent(evt events.Event) error
}
