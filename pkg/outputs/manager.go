package outputs

import (
	"github.com/ncarlier/trackr/pkg/events"
	"github.com/ncarlier/trackr/pkg/logger"
)

const maxEntriesChanSize = 5000

// Manager used to handle outputs worker
type Manager struct {
	outputs []Output
}

// NewOutputsManager create an outputs manager
func NewOutputsManager(outputs []Output) (*Manager, error) {
	manager := Manager{
		outputs: make([]Output, len(outputs)),
	}

	for idx, output := range outputs {
		if err := output.Connect(); err != nil {
			logger.Error.Printf("unable to connect to the output writer: %v", err)
		} else {
			manager.outputs[idx] = output
		}
	}

	return &manager, nil
}

// SendEvent sent event to all outputs
func (m *Manager) SendEvent(event events.Event) {
	for _, out := range m.outputs {
		if err := out.SendEvent(event); err != nil {
			logger.Error.Println("unable to send event to the output:", err)
		}
	}
}

// Shutdown stop the manager
func (m *Manager) Shutdown() {
	for _, out := range m.outputs {
		if err := out.Close(); err != nil {
			logger.Error.Println("unable to close the output:", err)
		}
	}
}
